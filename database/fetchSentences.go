package database

import (
	"fmt"
	"log"
)

// FetchSentences gets all the sentences of the paragraph
func FetchSentences(start int, end int, isComplete bool) ([]string, error) {
	var words []string
	var sentenceQuery string
	// Consideration for unfinished paragraph
	if !isComplete {
		sentenceQuery = fmt.Sprintf("Select word from sentence where sentence_id >= %d", start)
	} else {
		sentenceQuery = fmt.Sprintf("Select word from sentence where sentence_id >= %d and sentence_id < %d", start, end)
	}
	// Get sentences of the story
	sentenceStmt, err := db.Query(sentenceQuery)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer sentenceStmt.Close()

	var word string
	for sentenceStmt.Next() {
		err = sentenceStmt.Scan(&word)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		words = append(words, word)
	}

	return words, nil
}
