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

// StoryDB is used to commnunicate and use database object
type StoryDB struct {
	db *sql.DB
}

// InitDB creates DB Connection object
func InitDB(s *StoryDB) error {
	var err error
	// DB Connection parameters (MySQL)
	dbSource := strings.TrimPrefix((os.Getenv("DATABASE_DSN")), "mysql://")

	// Adding parseTime to process/parse timestamp into time.Time
	s.db, err = sql.Open(dbDriver, dbSource+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	return s.db.Ping()
}
