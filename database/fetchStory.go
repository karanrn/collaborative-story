package database

import (
	"fmt"
	"log"
	"time"

	"CollaborativeStory/colab/models"
)

// FetchStory retruns specific story requested
func FetchStory(storyID string) (models.DetailedStory, error) {

	var story models.DetailedStory
	// Get story details from story table
	storyStmt, err := db.Query(fmt.Sprintf("Select story_id, title, ifnull(start_paragraph, 0), ifnull(end_paragraph, 0), created_at, updated_at from story where story_id = %s", storyID))
	if err != nil {
		log.Println(err.Error())
		return models.DetailedStory{}, err
	}
	defer storyStmt.Close()

	var createTs, updateTs time.Time
	if storyStmt.Next() {
		err = storyStmt.Scan(&story.ID, &story.Title, &story.StartParagraph, &story.EndParagraph, &createTs, &updateTs)
		if err != nil {
			log.Println(err.Error())
			return models.DetailedStory{}, err
		}
		// Converting timestamp to TZ format
		story.CreatedAt = createTs.Format(time.RFC3339Nano)
		story.UpdatedAt = updateTs.Format(time.RFC3339Nano)
	}

	return story, nil
}
