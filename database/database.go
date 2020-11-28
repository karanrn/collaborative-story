package database

import (
	"database/sql"
	"os"
)

const (
	// mysql database driver
	dbDriver = "mysql"
)

// DBConn creates DB Connection object
func DBConn() (db *sql.DB) {
	// DB Connection parameters (MySQL)
	dbSource := os.Getenv("VERLOOP_DSN")

	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		panic(err.Error())
	}

	return db
}
