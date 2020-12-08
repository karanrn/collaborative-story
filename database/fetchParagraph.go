package database

import (
	"fmt"

	"CollaborativeStory/colab/models"
)

// FetchParagraphs returns/fetches paragraphs of the story
func (s *StoryDB) FetchParagraphs(start int, end int, isComplete bool) ([]models.Paragraph, error) {
	var paraQuery string
	// Consideration for unfinished story
	if !isComplete {
		paraQuery = fmt.Sprintf("Select paragraph_id, start_sentence, ifnull(end_sentence, 0) from paragraph where paragraph_id >= %d", start)
	} else {
		paraQuery = fmt.Sprintf("Select paragraph_id, start_sentence, end_sentence from paragraph where paragraph_id >= %d and paragraph_id < %d", start, end)
	}

	paragraphStmt, err := s.db.Query(paraQuery)
	if err != nil {
		return nil, err
	}
	defer paragraphStmt.Close()

	var allParagraphs []models.Paragraph
	var tmpParagraph models.Paragraph
	for paragraphStmt.Next() {
		err = paragraphStmt.Scan(&tmpParagraph.ID, &tmpParagraph.StartSentence, &tmpParagraph.EndSentence)
		if err != nil {
			return nil, err
		}
		allParagraphs = append(allParagraphs, tmpParagraph)
	}

	return allParagraphs, nil
}
