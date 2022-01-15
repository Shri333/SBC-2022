package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

type GroupMutationBody struct {
	Name string `binding:"required"`
}

func getGroups(c *gin.Context) {
	value, _ := c.Get("db")
	db := value.(*sql.DB)
	rows, _ := db.QueryContext(c, `
		SELECT items.*, groups.name AS groupName
		FROM items LEFT JOIN groups ON items.groupId = groups.id
	`)

	groupMap := make(map[int]*Group)
	for rows.Next() {
		var (
			item      Item
			groupId   sql.NullInt32
			groupName sql.NullString
		)

		rows.Scan(&item.Id, &item.Name, &item.Count, &groupId, &groupName)
		id := 0
		if groupId.Valid {
			id = int(groupId.Int32)
		}

		if value, ok := groupMap[id]; ok {
			value.Items = append(value.Items, item)
		} else {
			groupMap[id] = &Group{id, groupName.String, []Item{item}}
		}
	}

	groups := []*Group{}
	for _, value := range groupMap {
		groups = append(groups, value)
	}

	sort.Slice(groups, func(i, j int) bool { return groups[i].Id < groups[j].Id })
	c.JSON(http.StatusOK, groups)
}

func getGroupById(c *gin.Context) {
	value, _ := c.Get("db")
	db := value.(*sql.DB)
	idstr := c.Param("id")
	row := db.QueryRowContext(c, `SELECT id, name FROM groups WHERE id = ?`, idstr)
	group := Group{}

	if row.Scan(&group.Id, &group.Name) != nil {
		response := fmt.Sprintf("error: no group with id '%s' exists", idstr)
		c.String(http.StatusNotFound, response)
	} else {
		rows, _ := db.QueryContext(c, `
			SELECT items.id, items.name, items.count FROM items
			JOIN groups ON items.groupId = groups.id
			WHERE groups.id = ?
		`, idstr)

		for rows.Next() {
			item := Item{}
			rows.Scan(&item.Id, &item.Name, &item.Count)
			group.Items = append(group.Items, item)
		}

		c.JSON(http.StatusOK, group)
	}
}

func postGroup(c *gin.Context) {
	value, _ := c.Get("db")
	db := value.(*sql.DB)
	body := &GroupMutationBody{}
	if !validateMutationBody(c, body) {
		return
	}

	res, err := db.ExecContext(c, `INSERT INTO groups (name) VALUES (?)`, body.Name)
	id := 0
	if res != nil {
		id1, _ := res.LastInsertId()
		id = int(id1)
	}

	handleGroupMutationResult(c, db, body, err, id, http.StatusCreated)
}

func putGroup(c *gin.Context) {
	value, _ := c.Get("db")
	db := value.(*sql.DB)
	idstr := c.Param("id")
	body := &GroupMutationBody{}
	if !validateMutationBody(c, body) {
		return
	}

	id, row := 0, db.QueryRowContext(c, "SELECT id FROM groups WHERE id = ?", idstr)
	if row.Scan(&id) != nil {
		message := fmt.Sprintf("error: no group with id '%s' exists", idstr)
		c.String(http.StatusNotFound, message)
	} else {
		_, err := db.ExecContext(c, "UPDATE groups SET name = ? WHERE id = ?", body.Name, id)
		handleGroupMutationResult(c, db, body, err, id, http.StatusOK)
	}
}

func deleteGroup(c *gin.Context) {
	value, _ := c.Get("db")
	db := value.(*sql.DB)
	idstr := c.Param("id")
	res, _ := db.ExecContext(c, "DELETE FROM groups WHERE id = ?", idstr)

	if rows, _ := res.RowsAffected(); rows == 0 {
		message := fmt.Sprintf("error: no group with id '%s' exists", idstr)
		c.String(http.StatusNotFound, message)
		return
	}

	c.Status(http.StatusOK)
}
