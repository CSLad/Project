package database

import (
	"database/sql"
	"time"
)

type appdbimpl struct {
	c *sql.DB
}

type User struct {
	Username  string
	Following string
	Banned    string
}

type Image struct {
	ImageURL  string    `json:"imageurl"`
	Username  string    `json:"username"`
	Likes     int       `json:"likes"`
	Comments  string    `json:"comments"`
	CreatedAt time.Time `json:"created_at"`
}

// Example query function
func (db *appdbimpl) GetName() (string, error) {
	var name string
	err := db.c.QueryRow("SELECT name FROM example_table WHERE id=1").Scan(&name)
	return name, err
}
