package database

import (
	"time"

	"CollaborativeStory/colab/models"
)

// AddStory creates a new story
func (s StoryDB) AddStory(storyID int, word string, isNew bool) error {
	// Add new story
	addStoryStmt, err := s.db.Prepare("insert into story (story_id, title, created_at) values (?, ?, ?)")
	if err != nil {
		return err
	}
	defer addStoryStmt.Close()

	// Update title
	updateTitleStmt, err := s.db.Prepare("update story set title = concat(title, \" \", ?), updated_at = ? where story_id = ?")
	if err != nil {
		return err
	}
	defer updateTitleStmt.Close()

	if isNew {
		_, err = addStoryStmt.Exec(storyID, word, time.Now().In(time.UTC))
		if err != nil {
			return err
		}
	} else {
		_, err = updateTitleStmt.Exec(word, time.Now().In(time.UTC), storyID)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetLatestStory gets latest story, unfinished or creates new one
func (s *StoryDB) GetLatestStory() (models.Story, error) {
	var story models.Story
	// Get unfinished story
	storyStmt, err := s.db.Query("select IFNULL(story_id, 0), title, IFNULL(start_paragraph, 0) from story where start_paragraph is not NULL and end_paragraph is NULL;")
	if err != nil {
		return models.Story{}, err
	}
	defer storyStmt.Close()
	if storyStmt.Next() {
		err = storyStmt.Scan(&story.ID, &story.Title, &story.StartParagraph)
		if err != nil {
			return models.Story{}, err
		}
	}

	// Get the max value
	if story.ID == 0 {
		// Get the latest story with only title in creation
		lastStoryStmt, err := s.db.Query("select IFNULL(story_id, 0), title from story where start_paragraph is null order by story_id desc limit 1")
		if err != nil {
			return models.Story{}, err
		}
		defer lastStoryStmt.Close()
		if lastStoryStmt.Next() {
			err = lastStoryStmt.Scan(&story.ID, &story.Title)
			if err != nil {
				return models.Story{}, err
			}
		}

		// Get the max value if it is brand new
		if story.ID == 0 {
			maxStoryStmt, err := s.db.Query("select ifnull(max(story_id), 0) from story")
			if err != nil {
				return models.Story{}, err
			}
			defer maxStoryStmt.Close()
			if maxStoryStmt.Next() {
				err = maxStoryStmt.Scan(&story.ID)
				if err != nil {
					return models.Story{}, err
				}
			}
		}

	}

	return story, nil
}

// UpdateStoryParagraph updates or completes story with paragraphs
func (s *StoryDB) UpdateStoryParagraph(storyID int, paragraphID int, isEnd bool) error {
	// Update story
	// Update start of story (paragraph)
	updateStartStoryStmt, err := s.db.Prepare("update story set start_paragraph = ?, updated_at = ? where story_id = ?")
	if err != nil {
		return err
	}
	defer updateStartStoryStmt.Close()

	// Update end of story (paragraph)
	updateEndStoryStmt, err := s.db.Prepare("update story set end_paragraph = ?, updated_at = ? where story_id = ?")
	if err != nil {
		return err
	}
	defer updateEndStoryStmt.Close()

	if !isEnd {
		_, err = updateStartStoryStmt.Exec(paragraphID, time.Now().In(time.UTC), storyID)
		if err != nil {
			return err
		}
	} else {
		_, err = updateEndStoryStmt.Exec(paragraphID, time.Now().In(time.UTC), storyID)
		if err != nil {
			return err
		}
	}

	return nil
}

// UpdateStoryTimestamp updates updated_at timestamp whenever a new word is added
func (s StoryDB) UpdateStoryTimestamp(storyID int) error {
	// Update last updated timestamp for the word added to story
	updateTimeStoryStmt, err := s.db.Prepare("update story set updated_at = ? where story_id = ?")
	if err != nil {
		return err
	}
	defer updateTimeStoryStmt.Close()

	_, err = updateTimeStoryStmt.Exec(time.Now().In(time.UTC), storyID)
	if err != nil {
		return err
	}

	return nil
}
