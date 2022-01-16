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
		SELECT items.id, items.name, items.count, groups.id AS groupId, groups.name AS groupName
		FROM items LEFT JOIN groups ON items.groupId = groups.id
		UNION
		SELECT items.id, items.name, items.count, groups.id AS groupId, groups.name AS groupName
		FROM groups LEFT JOIN items ON groups.id = items.groupId
	`)

	groupMap := map[int]*Group{0: {0, "", []Item{}}}
	for rows.Next() {
		var (
			id        sql.NullInt32
			name      sql.NullString
			count     sql.NullInt32
			groupId   sql.NullInt32
			groupName sql.NullString
		)

		rows.Scan(&id, &name, &count, &groupId, &groupName)
		if id.Valid {
			grId := 0
			if groupId.Valid {
				grId = int(groupId.Int32)
			}

			item := Item{int(id.Int32), name.String, int(count.Int32), nil}
			if value, ok := groupMap[grId]; ok {
				value.Items = append(value.Items, item)
			} else {
				groupMap[grId] = &Group{grId, groupName.String, []Item{item}}
			}
		} else {
			grId := int(groupId.Int32)
			groupMap[grId] = &Group{grId, groupName.String, []Item{}}
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
