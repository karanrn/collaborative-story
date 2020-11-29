package story

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"CollaborativeStory/database"
)

/*
type storySentence struct {
	words []string
}
*/

type storyParagraph struct {
	ID            int      `json:"-"`
	StartSentence int      `json:"-"`
	EndSentence   int      `json:"-"`
	Sentences     []string `json:"sentences"`
}

type detailedStory struct {
	ID             int              `json:"id"`
	Title          string           `json:"title"`
	StartParagraph int              `json:"-"`
	EndParagraph   int              `json:"-"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
	Paragraphs     []storyParagraph `json:"paragraphs"`
}

// GetStory gets the specific story basis story_id
func GetStory(w http.ResponseWriter, r *http.Request) {
	storyID := mux.Vars(r)["id"]

	var resStory detailedStory

	db := database.DBConn()
	defer db.Close()

	// Get story details from story table
	storyStmt, err := db.Query(fmt.Sprintf("Select story_id, title, ifnull(start_paragraph, 0), ifnull(end_paragraph, 0), created_at, updated_at from story where story_id = %s", storyID))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
		return
	}
	defer storyStmt.Close()

	var createTs, updateTs time.Time
	if storyStmt.Next() {
		err = storyStmt.Scan(&resStory.ID, &resStory.Title, &resStory.StartParagraph, &resStory.EndParagraph, &createTs, &updateTs)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
			return
		}
		// Converting timestamp to TZ format
		resStory.CreatedAt = createTs.Format(time.RFC3339Nano)
		resStory.UpdatedAt = updateTs.Format(time.RFC3339Nano)
	}

	if resStory.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(`{'error': 'story does not exist'}`)
		return
	}

	// Get paragraph details from paragraph table
	if resStory.StartParagraph > 0 {
		var paraQuery string
		// Consideration for unfinished story
		if resStory.EndParagraph == 0 {
			paraQuery = fmt.Sprintf("Select paragraph_id, start_sentence, ifnull(end_sentence, 0) from paragraph where paragraph_id >= %d", resStory.StartParagraph)
		} else {
			paraQuery = fmt.Sprintf("Select paragraph_id, start_sentence, end_sentence from paragraph where paragraph_id >= %d and paragraph_id < %d", resStory.StartParagraph, resStory.EndParagraph)
		}

		paragraphStmt, err := db.Query(paraQuery)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
			return
		}
		defer paragraphStmt.Close()

		var tmpParagraph storyParagraph
		var sentenceQuery string
		for paragraphStmt.Next() {
			err = paragraphStmt.Scan(&tmpParagraph.ID, &tmpParagraph.StartSentence, &tmpParagraph.EndSentence)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
				return
			}

			// Consideration for unfinished paragraph
			if tmpParagraph.EndSentence == 0 {
				sentenceQuery = fmt.Sprintf("Select word from sentence where sentence_id >= %d", tmpParagraph.StartSentence)
			} else {
				sentenceQuery = fmt.Sprintf("Select word from sentence where sentence_id >= %d and sentence_id < %d", tmpParagraph.StartSentence, tmpParagraph.EndSentence)
			}
			// Get sentences of the story
			sentenceStmt, err := db.Query(sentenceQuery)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
				return
			}
			defer paragraphStmt.Close()

			var word string
			for sentenceStmt.Next() {
				err = sentenceStmt.Scan(&word)
				if err != nil {
					log.Println(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
					return
				}

				tmpParagraph.Sentences = append(tmpParagraph.Sentences, word)
			}

			resStory.Paragraphs = append(resStory.Paragraphs, tmpParagraph)
		}
	}

	result, err := json.Marshal(&resStory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(string(result))
}
