package database

import (
	"database/sql"
	"os"
	"strings"

	"CollaborativeStory/colab/models"
)

const (
	// mysql database driver
	dbDriver = "mysql"
)

// StoryDB is used to commnunicate and use database object
type StoryDB struct {
	db *sql.DB
}

// New Creates new StoryDB
func New() StoryDB {
	return StoryDB{}
}

// DB defines interfaces of the methods for database
type DB interface {
	FetchStory(storyID string) (models.DetailedStory, error)
	FetchStories(sort string, order string, offset int64, limit int64) ([]models.Story, error)
	FetchParagraphs(start int, end int, isComplete bool) ([]models.Paragraph, error)
	FetchSentences(start int, end int, isComplete bool) ([]string, error)
	AddStory(storyID int, word string, isNew bool) error
	AddToSentence(word string) (sentenceID int, err error)
	AddToParagraph(sentenceID int) (paragraphID int, err error)
	GetLatestStory() (models.Story, error)
	UpdateStoryParagraph(storyID int, paragraphID int, isEnd bool) error
	UpdateStoryTimestamp(storyID int) error
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
