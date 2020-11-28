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

func sentence(db *sql.DB, word string) error {
	// Get the unfinished sentence id
	sentenceStmt, err := db.Query("select IFNULL(max(sentence_id), 0) from sentence group by sentence_id having count(word) < 15 order by sentence_id desc;")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	var sentenceID int
	if sentenceStmt.Next() {
		err = sentenceStmt.Scan(&sentenceID)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	// Get the max value
	if sentenceID == 0 {
		maxSentenceStmt, err := db.Query("select IFNULL(max(sentence_id), 0) from sentence;")
		if err != nil {
			log.Println(err.Error())
			return err
		}
		if maxSentenceStmt.Next() {
			err = maxSentenceStmt.Scan(&sentenceID)
			if err != nil {
				log.Println(err.Error())
				return err
			}
		}
		// Create next sentence
		sentenceID++
	}
	addSentence, err := db.Prepare("insert into sentence values (?, ?)")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, err = addSentence.Exec(sentenceID, word)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
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

		err = sentence(db, reqWord["word"])
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
