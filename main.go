package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// mysql database driver
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

// AddToParagraph adds sentence to a paragraph
func AddToParagraph(sentenceID int) (paragraphID int, err error) {
	// DB Connection
	db := DBConn()
	defer db.Close()

	// Get the unfinished paragraph
	var startParagraph int
	paragraphStmt, err := db.Query("select IFNULL(paragraph_id, 0), IFNULL(start_sentence, 0) from paragraph where start_sentence is not NULL and end_sentence is NULL;")
	if err != nil {
		log.Println(err.Error())
		return paragraphID, err
	}
	defer paragraphStmt.Close()
	if paragraphStmt.Next() {
		err = paragraphStmt.Scan(&paragraphID, &startParagraph)
		if err != nil {
			log.Println(err.Error())
			return paragraphID, err
		}
	}

	// Get the max value
	if paragraphID == 0 {
		lastParagraphStmt, err := db.Query("select IFNULL(max(paragraph_id), 0) from paragraph;")
		if err != nil {
			log.Println(err.Error())
			return paragraphID, err
		}
		defer lastParagraphStmt.Close()
		if lastParagraphStmt.Next() {
			err = lastParagraphStmt.Scan(&paragraphID)
			if err != nil {
				log.Println(err.Error())
				return paragraphID, err
			}
		}
		// Create next paragraph
		paragraphID++
	}

	// Check the size of the paragraph and update/create paragraph
	updateParagraph, err := db.Prepare("update paragraph set end_sentence = ? where paragraph_id = ?")
	if err != nil {
		log.Println(err.Error())
		return paragraphID, err
	}
	defer updateParagraph.Close()

	addParagraph, err := db.Prepare("insert into paragraph (paragraph_id, start_sentence) values (?, ?)")
	if err != nil {
		log.Println(err.Error())
		return paragraphID, err
	}
	defer addParagraph.Close()

	if startParagraph != 0 && (sentenceID-startParagraph) == 10 {
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

// AddToSentence adds words to form a sentence
func AddToSentence(word string) (sentenceID int, err error) {
	// DB Connection
	db := DBConn()
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
	addSentence, err := db.Prepare("insert into sentence values (?, ?)")
	if err != nil {
		log.Println(err.Error())
		return sentenceID, err
	}
	defer addSentence.Close()

	_, err = addSentence.Exec(sentenceID, word)
	if err != nil {
		log.Println(err.Error())
		return sentenceID, err
	}

	return sentenceID, nil
}

// AddToStory adds word to the story
func AddToStory(w http.ResponseWriter, r *http.Request) {
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

		// Get unfinished story
		var storyID int
		var title string
		var startParagraph int
		storyStmt, err := db.Query("select IFNULL(story_id, 0), title, IFNULL(start_paragraph, 0) from story where start_paragraph is not NULL and end_paragraph is NULL;")
		if err != nil {
			log.Println(err.Error())
		}
		defer storyStmt.Close()
		if storyStmt.Next() {
			err = storyStmt.Scan(&storyID, &title, &startParagraph)
			if err != nil {
				log.Println(err.Error())
			}
		}

		// Add new story
		addStoryStmt, err := db.Prepare("insert into story (story_id, title) values (?, ?)")
		if err != nil {
			log.Println(err.Error())
		}
		defer addStoryStmt.Close()

		// Update title
		updateTitleStmt, err := db.Prepare("update story set title = concat(title, \" \", ?) where story_id = ?")
		if err != nil {
			log.Println(err.Error())
		}
		defer updateTitleStmt.Close()

		// Update story
		updateStartStoryStmt, err := db.Prepare("update story set start_paragraph = ? where story_id = ?")
		if err != nil {
			log.Println(err.Error())
		}
		defer updateStartStoryStmt.Close()

		updateEndStoryStmt, err := db.Prepare("update story set end_paragraph = ? where story_id = ?")
		if err != nil {
			log.Println(err.Error())
		}
		defer updateEndStoryStmt.Close()

		// Get the max value
		if storyID == 0 {
			lastStoryStmt, err := db.Query("select IFNULL(story_id, 0), title from story where start_paragraph is null order by story_id desc limit 1")
			if err != nil {
				log.Println(err.Error())
			}
			defer lastStoryStmt.Close()
			if lastStoryStmt.Next() {
				err = lastStoryStmt.Scan(&storyID, &title)
				if err != nil {
					log.Println(err.Error())
				}
			}

			if storyID == 0 {
				maxStoryStmt, err := db.Query("select max(story_id) from story")
				if err != nil {
					log.Println(err.Error())
				}
				defer maxStoryStmt.Close()
				if maxStoryStmt.Next() {
					err = maxStoryStmt.Scan(&storyID)
					if err != nil {
						log.Println(err.Error())
					}
				}
			}

		}

		if title == "" {
			_, err = addStoryStmt.Exec(storyID+1, reqWord["word"])
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			if title != "" && len(strings.Split(title, " ")) < 2 {
				_, err = updateTitleStmt.Exec(reqWord["word"], storyID)
				if err != nil {
					log.Println(err.Error())
				}
			} else {
				sentenceID, err := AddToSentence(reqWord["word"])
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(`{'error': 'Internal server error'}`)
					return
				}

				paragraphID, err := AddToParagraph(sentenceID)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(`{'error': 'Internal server error'}`)
					return
				}

				if startParagraph == 0 {
					_, err = updateStartStoryStmt.Exec(paragraphID, storyID)
					if err != nil {
						log.Println(err.Error())
					}
				}

				if startParagraph != 0 && (paragraphID-startParagraph) == 7 {
					_, err = updateEndStoryStmt.Exec(paragraphID, storyID)
					if err != nil {
						log.Println(err.Error())
					}
				}
			}
		}

	} else {
		fmt.Fprintf(w, "Method not supported")
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/add", AddToStory)
	fmt.Println("Serving on :9000")
	log.Fatal(http.ListenAndServe(":9000", mux))
}
