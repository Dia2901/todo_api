package db

import (
	"database/sql"
	"gojek/web-server-gin/pkg/config"
	"gojek/web-server-gin/pkg/handleError"
)

var DB *sql.DB

func SetupDB() {
	connStr := "postgres://postgres:" + config.DBpass + "@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	handleError.Check(err)
	DB = db
	createTable()
}
