package main

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func initDB(db *sql.DB) {
	db.Exec(`
		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
	`)

	db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			count INTEGER NOT NULL,
			groupId INTEGER,
			FOREIGN KEY (groupId) REFERENCES groups(id) ON DELETE CASCADE
		);
	`)
}

func main() {
	db, _ := sql.Open("sqlite3", "data.db?_foreign_keys=on")
	initDB(db)

	gin.SetMode(gin.ReleaseMode) // change to gin.DebugMode to run in "debug" mode
	engine := gin.New()
	engine.SetTrustedProxies(nil)
	engine.Use(cors.Default())
	engine.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	items := engine.Group("/api/items")
	items.GET("", getItems)
	items.GET(":id", getItemById)
	items.POST("", postItem)
	items.PUT(":id", putItem)
	items.DELETE(":id", deleteItem)

	groups := engine.Group("/api/groups")
	groups.GET("", getGroups)
	groups.GET(":id", getGroupById)
	groups.POST("", postGroup)
	groups.PUT(":id", putGroup)
	groups.DELETE(":id", deleteGroup)

	engine.Static("/static", "./static/static")
	engine.NoRoute(func(c *gin.Context) { c.File("./static/index.html") })
	engine.Run()
}
