package story

import (
	"CollaborativeStory/database"
	"net/http"
)

// ColabStory uses database methods
type ColabStory struct {
	Database database.DB
}

// Service defines interface for story
type Service interface {
	PostStory(s database.StoryDB) http.HandlerFunc
	GetStories(s database.StoryDB) http.HandlerFunc
	GetStory(s database.StoryDB) http.HandlerFunc
}

// New factory function to create new ColabStory
func New(db database.DB) ColabStory {
	return ColabStory{Database: db}
}
