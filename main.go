package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver = "mysql"
)

// DBConn creates DB Connection object
func DBConn() (db *sql.DB) {
	// DB Connection parameters (MySQL)
	dbSource := os.Getenv("VERLOOP_DSN")

	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		panic(err.Error())
	}

	return db
}

func paragraphOps(db *sql.DB, sentenceID int) (paragraphID int, err error) {
	// Get the unfinished paragraph
	var startParagraph int
	paragraphStmt, err := db.Query("select IFNULL(paragraph_id, 0), IFNULL(start_sentence, 0) from paragraph where start_sentence is not NULL and end_sentence is NULL;")
	if err != nil {
		log.Println(err.Error())
		return paragraphID, err
	}

	if paragraphStmt.Next() {
		err = paragraphStmt.Scan(&paragraphID, &startParagraph)
		if err != nil {
			log.Println(err.Error())
			return paragraphID, err
		}
	}

	// Get the max value
	if paragraphID == 0 {
		maxParagraphStmt, err := db.Query("select IFNULL(max(paragraph_id), 0) from paragraph;")
		if err != nil {
			log.Println(err.Error())
			return paragraphID, err
		}
		if maxParagraphStmt.Next() {
			err = maxParagraphStmt.Scan(&paragraphID)
			if err != nil {
				log.Println(err.Error())
				return paragraphID, err
			}
		}
		// Create next sentence
		paragraphID++
	}

	// Check the size of the paragraph and update/create paragraph
	updateParagraph, err := db.Prepare("update paragraph set end_sentence = ? where paragraph_id = ?")
	updateParagraph.Close()
	if err != nil {
		log.Println(err.Error())
		return paragraphID, err
	}

	addParagraph, err := db.Prepare("insert into paragraph (paragraph_id, start_sentence) values (?, ?)")
	defer addParagraph.Close()
	if err != nil {
		log.Println(err.Error())
		return paragraphID, err
	}

	if startParagraph != 0 && ((sentenceID+1)-startParagraph) == 10 {
		_, err = updateParagraph.Exec(sentenceID, paragraphID)
		if err != nil {
			log.Println(err.Error())
			return paragraphID, err
		}
	}

	if startParagraph == 0 {
		_, err = addParagraph.Exec(paragraphID, sentenceID)
		if err != nil {
			log.Println(err.Error())
			return paragraphID, err
		}
	}

	return paragraphID, nil
}

// All the operations for sentence
func sentenceOps(db *sql.DB, word string) (sentenceID int, err error) {
	// Get the unfinished sentence id
	sentenceStmt, err := db.Query("select IFNULL(max(sentence_id), 0) from sentence group by sentence_id having count(word) < 15 order by sentence_id desc;")
	if err != nil {
		log.Println(err.Error())
		return sentenceID, err
	}
	//var sentenceID int
	if sentenceStmt.Next() {
		err = sentenceStmt.Scan(&sentenceID)
		if err != nil {
			log.Println(err.Error())
			return sentenceID, err
		}
	}

	// Get the max value
	if sentenceID == 0 {
		maxSentenceStmt, err := db.Query("select IFNULL(max(sentence_id), 0) from sentence;")
		if err != nil {
			log.Println(err.Error())
			return sentenceID, err
		}
		if maxSentenceStmt.Next() {
			err = maxSentenceStmt.Scan(&sentenceID)
			if err != nil {
				log.Println(err.Error())
				return sentenceID, err
			}
		}
		// Create next sentence
		sentenceID++
	}
	addSentence, err := db.Prepare("insert into sentence values (?, ?)")
	defer addSentence.Close()
	if err != nil {
		log.Println(err.Error())
		return sentenceID, err
	}

	_, err = addSentence.Exec(sentenceID, word)
	if err != nil {
		log.Println(err.Error())
		return sentenceID, err
	}

	return sentenceID, nil
}

// AddWord adds word to the story
func AddWord(w http.ResponseWriter, r *http.Request) {
	var reqWord map[string]string

	// DB Connection
	db := DBConn()
	defer db.Close()

	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&reqWord)
		if err != nil {
			json.NewEncoder(w).Encode(`{'error': 'Error in decoding JSON'}`)
			return
		}

		sentenceID, err := sentenceOps(db, reqWord["word"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(`{'error': 'Internal server error'}`)
			return
		}

		_, err = paragraphOps(db, sentenceID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(`{'error': 'Internal server error'}`)
			return
		}

	} else {
		fmt.Fprintf(w, "Method not supported")
	}

}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/add", AddWord)
	fmt.Println("Serving on :9000")
	log.Fatal(http.ListenAndServe(":9000", mux))
}
