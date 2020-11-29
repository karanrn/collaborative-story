package sentence

import (
	"log"
	"time"

	"CollaborativeStory/database"
)

// AddToSentence adds words to form a sentence
func AddToSentence(word string) (sentenceID int, err error) {
	// DB Connection
	db := database.DBConn()
	defer db.Close()

	// Get the unfinished sentence id
	sentenceStmt, err := db.Query("select IFNULL(max(sentence_id), 0) from sentence group by sentence_id having count(word) < 15 order by sentence_id desc;")
	if err != nil {
		log.Println(err.Error())
		return sentenceID, err
	}
	defer sentenceStmt.Close()
	if sentenceStmt.Next() {
		err = sentenceStmt.Scan(&sentenceID)
		if err != nil {
			log.Println(err.Error())
			return sentenceID, err
		}
	}

	// Get the max value
	if sentenceID == 0 {
		lastSentenceStmt, err := db.Query("select IFNULL(max(sentence_id), 0) from sentence;")
		if err != nil {
			log.Println(err.Error())
			return sentenceID, err
		}
		defer lastSentenceStmt.Close()
		if lastSentenceStmt.Next() {
			err = lastSentenceStmt.Scan(&sentenceID)
			if err != nil {
				log.Println(err.Error())
				return sentenceID, err
			}
		}
		// Create next sentence
		sentenceID++
	}
	addSentence, err := db.Prepare("insert into sentence values (?, ?, ?)")
	if err != nil {
		log.Println(err.Error())
		return sentenceID, err
	}
	defer addSentence.Close()

	_, err = addSentence.Exec(sentenceID, word, time.Now())
	if err != nil {
		log.Println(err.Error())
		return sentenceID, err
	}

	return sentenceID, nil
}
