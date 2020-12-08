package database

import (
	"fmt"
	"time"

	"CollaborativeStory/colab/models"
)

// FetchStories gets all the stories from the database
func (s StoryDB) FetchStories(sort string, order string, offset int64, limit int64) ([]models.Story, error) {
	var results []models.Story

	// Get all records
	storiesStmt, err := s.db.Query(fmt.Sprintf("Select story_id, title, created_at, updated_at from story order by %s %s limit %d, %d ", sort, order, offset, limit))
	if err != nil {
		return nil, err
	}
	defer storiesStmt.Close()

	for storiesStmt.Next() {
		var iStory models.Story
		var createTs, updateTs time.Time
		err = storiesStmt.Scan(&iStory.ID, &iStory.Title, &createTs, &updateTs)
		if err != nil {
			return nil, err
		}
		// Converting timestamps to TZ format (RFC3339Nano)
		iStory.CreatedAt = createTs.Format(time.RFC3339Nano)
		iStory.UpdatedAt = updateTs.Format(time.RFC3339Nano)
		results = append(results, iStory)
	}

	return results, nil
}
