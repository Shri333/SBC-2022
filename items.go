package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ItemMutationBody struct {
	Name    string `binding:"required"`
	Count   *int   `binding:"required,gt=0"`
	GroupId *int
}

func getItems(c *gin.Context) {
	value, _ := c.Get("db")
	db := value.(*sql.DB)
	rows, _ := db.QueryContext(c, `
		SELECT items.*, groups.name AS groupName
		FROM items LEFT JOIN groups ON items.groupId = groups.id
	`)

	items := []Item{}
	for rows.Next() {
		var (
			item      Item
			groupId   sql.NullInt32
			groupName sql.NullString
		)

		rows.Scan(&item.Id, &item.Name, &item.Count, &groupId, &groupName)
		if groupId.Valid {
			item.Parent = &Group{int(groupId.Int32), groupName.String, nil}
		}

		items = append(items, item)
	}

	c.JSON(http.StatusOK, items)
}

func getItemById(c *gin.Context) {
	value, _ := c.Get("db")
	db := value.(*sql.DB)
	idstr := c.Param("id")
	row := db.QueryRowContext(c, `
		SELECT items.*, groups.name AS groupName
		FROM items LEFT JOIN groups ON items.groupId = groups.id
		WHERE items.id = ?
	`, idstr)

	var (
		item      Item
		groupId   sql.NullInt32
		groupName sql.NullString
	)

	if row.Scan(&item.Id, &item.Name, &item.Count, &groupId, &groupName) != nil {
		message := fmt.Sprintf("error: no item with id '%s' exists", idstr)
		c.String(http.StatusNotFound, message)
		return
	}

	if groupId.Valid {
		item.Parent = &Group{int(groupId.Int32), groupName.String, nil}
	}

	c.JSON(http.StatusOK, item)
}

func postItem(c *gin.Context) {
	value, _ := c.Get("db")
	db := value.(*sql.DB)
	body := &ItemMutationBody{}
	if !validateMutationBody(c, body) {
		return
	}

	res, err := db.ExecContext(c, `
		INSERT INTO items (name, count, groupId) VALUES (?, ?, ?)
	`, body.Name, body.Count, body.GroupId)

	id := 0
	if res != nil {
		id1, _ := res.LastInsertId()
		id = int(id1)
	}

	handleItemMutationResult(c, db, body, err, id, http.StatusCreated)
}

func putItem(c *gin.Context) {
	value, _ := c.Get("db")
	db := value.(*sql.DB)
	idstr := c.Param("id")
	body := &ItemMutationBody{}
	if !validateMutationBody(c, body) {
		return
	}

	id, row := 0, db.QueryRowContext(c, "SELECT id FROM items WHERE id = ?", idstr)
	if row.Scan(&id) != nil {
		message := fmt.Sprintf("error: no item with id '%s' exists", idstr)
		c.String(http.StatusNotFound, message)
	} else {
		_, err := db.ExecContext(c, `
			UPDATE items SET name = ?, count = ?, groupId = ? WHERE id = ?
		`, body.Name, body.Count, body.GroupId, id)
		handleItemMutationResult(c, db, body, err, id, http.StatusOK)
	}
}

func deleteItem(c *gin.Context) {
	value, _ := c.Get("db")
	db := value.(*sql.DB)
	idstr := c.Param("id")
	res, _ := db.ExecContext(c, "DELETE FROM items WHERE id = ?", idstr)

	if rows, _ := res.RowsAffected(); rows == 0 {
		message := fmt.Sprintf("error: no item with id '%s' exists", idstr)
		c.String(http.StatusNotFound, message)
		return
	}

	c.Status(http.StatusOK)
}
