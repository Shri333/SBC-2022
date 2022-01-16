package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func validateMutationBody(c *gin.Context, body interface{}) bool {
	var object, name string
	err := c.Bind(body)
	switch value := body.(type) {
	case *ItemMutationBody:
		object = "item"
		name = value.Name
	case *GroupMutationBody:
		object = "group"
		name = value.Name
	}

	if err != nil {
		var response string
		switch e := err.(type) {
		case validator.ValidationErrors:
			field := strings.ToLower(e[0].Field())
			if e[0].Tag() == "gt" {
				response = fmt.Sprintf("error: %s %s must be greater than zero", object, field)
			} else {
				response = fmt.Sprintf("error: %s %s is missing", object, field)
			}
		case *json.UnmarshalTypeError:
			response = fmt.Sprintf("error: %s %s must be a %s", object, e.Field, e.Type)
		default:
			response = fmt.Sprintf("error: %s", err.Error())
		}

		c.String(http.StatusBadRequest, response)
		return false
	}

	if len(strings.TrimSpace(name)) == 0 {
		response := fmt.Sprintf("error: %s name cannot be empty", object)
		c.String(http.StatusBadRequest, response)
		return false
	}

	return true
}

func handleItemMutationResult(c *gin.Context, db *sql.DB, body *ItemMutationBody, err error, id, status int) {
	if err != nil {
		if message := err.Error(); strings.HasPrefix(message, "UNIQUE") {
			response := fmt.Sprintf("error: item with name '%s' already exists", body.Name)
			c.String(http.StatusBadRequest, response)
		} else {
			response := fmt.Sprintf("error: no group with id '%d' exists", *body.GroupId)
			c.String(http.StatusNotFound, response)
		}
	} else {
		item := Item{id, body.Name, *body.Count, nil}
		if body.GroupId != nil {
			var groupName string
			row := db.QueryRowContext(c, `SELECT name FROM groups WHERE id = ?`, body.GroupId)
			row.Scan(&groupName)
			item.Parent = &Group{*body.GroupId, groupName, nil}
		}

		c.JSON(status, item)
	}
}

func handleGroupMutationResult(c *gin.Context, db *sql.DB, body *GroupMutationBody, err error, id, status int) {
	if err != nil {
		response := fmt.Sprintf("error: group with name '%s' already exists", body.Name)
		c.String(http.StatusBadRequest, response)
	} else {
		items := []Item{}
		rows, _ := db.QueryContext(c, `
			SELECT items.id, items.name, items.count FROM items
			JOIN groups ON items.groupId = groups.id
			WHERE groups.id = ?
		`, id)

		for rows.Next() {
			item := Item{}
			rows.Scan(&item.Id, &item.Name, &item.Count)
			items = append(items, item)
		}

		c.JSON(status, Group{id, body.Name, items})
	}
}
