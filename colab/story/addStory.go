package story

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"CollaborativeStory/colab/paragraph"
	"CollaborativeStory/colab/sentence"
	"CollaborativeStory/database"
)

type response struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	CurrentSentence string `json:"current_sentence"`
}

// AddToStory adds word to the story
func AddToStory(w http.ResponseWriter, r *http.Request) {
	var reqWord map[string]string

	// DB Connection
	db := database.DBConn()
	defer db.Close()

	// Process if request is POST else reject
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&reqWord)
		if err != nil {
			json.NewEncoder(w).Encode(`{'error': 'error in decoding JSON'}`)
			return
		}

		word := strings.TrimSpace(reqWord["word"])
		// Check if multiple words are sent
		if word == "" || len(strings.Split(word, " ")) > 1 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(`{'error': 'multiple words sent'}`)
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
		addStoryStmt, err := db.Prepare("insert into story (story_id, title, created_at) values (?, ?, ?)")
		if err != nil {
			log.Println(err.Error())
		}
		defer addStoryStmt.Close()

		// Update title
		updateTitleStmt, err := db.Prepare("update story set title = concat(title, \" \", ?), updated_at = ? where story_id = ?")
		if err != nil {
			log.Println(err.Error())
		}
		defer updateTitleStmt.Close()

		// Update story
		// Update start of story (paragraph)
		updateStartStoryStmt, err := db.Prepare("update story set start_paragraph = ?, updated_at = ? where story_id = ?")
		if err != nil {
			log.Println(err.Error())
		}
		defer updateStartStoryStmt.Close()

		// Update end of story (paragraph)
		updateEndStoryStmt, err := db.Prepare("update story set end_paragraph = ?, updated_at = ? where story_id = ?")
		if err != nil {
			log.Println(err.Error())
		}
		defer updateEndStoryStmt.Close()

		// Update last updated timestamp for the word added to story
		updateTimeStoryStmt, err := db.Prepare("update story set updated_at = ?")
		if err != nil {
			log.Println(err.Error())
		}
		defer updateTimeStoryStmt.Close()

		// Get the max value
		if storyID == 0 {
			// Get the latest story with only title in creation
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

			// Get the max value if it is brand new
			if storyID == 0 {
				maxStoryStmt, err := db.Query("select ifnull(max(story_id), 0) from story")
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

		var storyResp response
		if title == "" {
			// Add title word to the new story
			_, err = addStoryStmt.Exec(storyID+1, word, time.Now().Format(time.RFC3339Nano))
			if err != nil {
				log.Println(err.Error())
			}

			storyResp.ID = storyID + 1
			storyResp.Title = word
			storyResp.CurrentSentence = ""
		} else {
			if title != "" && len(strings.Split(title, " ")) < 2 {
				// Update title of the story
				_, err = updateTitleStmt.Exec(word, time.Now().Format(time.RFC3339Nano), storyID)
				if err != nil {
					log.Println(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
					return
				}

				storyResp.ID = storyID
				storyResp.Title = title + " " + word
				storyResp.CurrentSentence = ""
			} else {
				// Add word to sentence of the story
				sentenceID, err := sentence.AddToSentence(word)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
					return
				}

				// Add/Update paragraph of the story
				paragraphID, err := paragraph.AddToParagraph(sentenceID)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
					return
				}

				if startParagraph == 0 {
					// Start a new story
					_, err = updateStartStoryStmt.Exec(paragraphID, time.Now().Format(time.RFC3339Nano), storyID)
					if err != nil {
						log.Println(err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
						return
					}
				}

				if startParagraph != 0 && (paragraphID-startParagraph) == 7 {
					// End the story
					_, err = updateEndStoryStmt.Exec(paragraphID, time.Now().Format(time.RFC3339Nano), storyID)
					if err != nil {
						log.Println(err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
						return
					}
				}

				// Update the story timestamp (updated_at)
				_, err = updateTimeStoryStmt.Exec(time.Now().Format(time.RFC3339Nano))
				if err != nil {
					log.Println(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
					return
				}
				storyResp.ID = storyID
				storyResp.Title = title
				storyResp.CurrentSentence = word
			}
		}

		w.WriteHeader(http.StatusOK)
		// Marshal the response
		resp, err := json.Marshal(storyResp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
			return
		}
		json.NewEncoder(w).Encode(string(resp))

	} else {
		// Reject all requests other than POST
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(`{'error': 'method not supported'}`)
	}
}
