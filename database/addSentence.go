package database

import (
	"time"
)

// AddToSentence adds words to form a sentence
func (s StoryDB) AddToSentence(word string) (sentenceID int, err error) {

	// Get the unfinished sentence id
	sentenceStmt, err := s.db.Query("select IFNULL(max(sentence_id), 0) from sentence group by sentence_id having count(word) < 15 order by sentence_id desc;")
	if err != nil {
		return sentenceID, err
	}
	defer sentenceStmt.Close()
	if sentenceStmt.Next() {
		err = sentenceStmt.Scan(&sentenceID)
		if err != nil {
			return sentenceID, err
		}
	}

	// Get the max value
	if sentenceID == 0 {
		lastSentenceStmt, err := s.db.Query("select IFNULL(max(sentence_id), 0) from sentence;")
		if err != nil {
			return sentenceID, err
		}
		defer lastSentenceStmt.Close()
		if lastSentenceStmt.Next() {
			err = lastSentenceStmt.Scan(&sentenceID)
			if err != nil {
				return sentenceID, err
			}
		}
		// Create next sentence
		sentenceID++
	}
	addSentence, err := s.db.Prepare("insert into sentence values (?, ?, ?)")
	if err != nil {
		return sentenceID, err
	}
	defer addSentence.Close()

	_, err = addSentence.Exec(sentenceID, word, time.Now().In(time.UTC))
	if err != nil {
		return sentenceID, err
	}

	return sentenceID, nil
}
