package database

import (
	"time"
)

// AddToParagraph adds sentence to a paragraph
func (s StoryDB) AddToParagraph(sentenceID int) (paragraphID int, err error) {

	// Get the unfinished paragraph
	var startParagraph int
	paragraphStmt, err := s.db.Query("select IFNULL(paragraph_id, 0), IFNULL(start_sentence, 0) from paragraph where start_sentence is not NULL and end_sentence is NULL;")
	if err != nil {
		return paragraphID, err
	}
	defer paragraphStmt.Close()
	if paragraphStmt.Next() {
		err = paragraphStmt.Scan(&paragraphID, &startParagraph)
		if err != nil {
			return paragraphID, err
		}
	}

	// Get the max value
	if paragraphID == 0 {
		lastParagraphStmt, err := s.db.Query("select IFNULL(max(paragraph_id), 0) from paragraph;")
		if err != nil {
			return paragraphID, err
		}
		defer lastParagraphStmt.Close()
		if lastParagraphStmt.Next() {
			err = lastParagraphStmt.Scan(&paragraphID)
			if err != nil {
				return paragraphID, err
			}
		}
		// Create next paragraph
		paragraphID++
	}

	// Check the size of the paragraph and update/create paragraph
	updateParagraph, err := s.db.Prepare("update paragraph set end_sentence = ?, updated_at = ? where paragraph_id = ?")
	if err != nil {
		return paragraphID, err
	}
	defer updateParagraph.Close()

	addParagraph, err := s.db.Prepare("insert into paragraph (paragraph_id, start_sentence, created_at) values (?, ?, ?)")
	if err != nil {
		return paragraphID, err
	}
	defer addParagraph.Close()

	if startParagraph != 0 && (sentenceID-startParagraph) == 10 {
		_, err = updateParagraph.Exec(sentenceID, time.Now().In(time.UTC), paragraphID)
		if err != nil {
			return paragraphID, err
		}
	}

	if startParagraph == 0 {
		_, err = addParagraph.Exec(paragraphID, sentenceID, time.Now().In(time.UTC))
		if err != nil {
			return paragraphID, err
		}
	}

	return paragraphID, nil
}
