package database

import (
	"database/sql"
	"os"
	"strings"
)

const (
	// mysql database driver
	dbDriver = "mysql"
)

// DBConn creates DB Connection object
func DBConn() (db *sql.DB) {
	// DB Connection parameters (MySQL)
	dbSource := strings.TrimPrefix((os.Getenv("DATABASE_DSN")), "mysql://")

	// Adding parseTime to process/parse timestamp into time.Time
	db, err := sql.Open(dbDriver, dbSource+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	return db
}
