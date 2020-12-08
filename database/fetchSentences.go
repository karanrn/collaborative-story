package database

import (
	"fmt"
)

// FetchSentences gets all the sentences of the paragraph
func (s StoryDB) FetchSentences(start int, end int, isComplete bool) ([]string, error) {
	var words []string
	var sentenceQuery string
	// Consideration for unfinished paragraph
	if !isComplete {
		sentenceQuery = fmt.Sprintf("Select word from sentence where sentence_id >= %d", start)
	} else {
		sentenceQuery = fmt.Sprintf("Select word from sentence where sentence_id >= %d and sentence_id < %d", start, end)
	}
	// Get sentences of the story
	sentenceStmt, err := s.db.Query(sentenceQuery)
	if err != nil {
		return nil, err
	}
	defer sentenceStmt.Close()

	var word string
	for sentenceStmt.Next() {
		err = sentenceStmt.Scan(&word)
		if err != nil {
			return nil, err
		}

		words = append(words, word)
	}

	return words, nil
}
