package story

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"CollaborativeStory/database"
)

/*
type storySentence struct {
	words []string
}
*/

type storyParagraph struct {
	ID            string   `json:"-"`
	StartSentence string   `json:"-"`
	EndSentence   string   `json:"-"`
	Sentences     []string `json:"sentences"`
}

type detailedStory struct {
	ID             string           `json:"id"`
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

	if storyStmt.Next() {
		err = storyStmt.Scan(&resStory.ID, &resStory.Title, &resStory.StartParagraph, &resStory.EndParagraph, &resStory.CreatedAt, &resStory.UpdatedAt)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
			return
		}
	}

	// Get paragraph details from paragraph table
	paragraphStmt, err := db.Query()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
		return
	}
	defer paragraphStmt.Close()

	var tmpParagraph storyParagraph
	for paragraphStmt.Next() {
		err = paragraphStmt.Scan(&tmpParagraph.ID, &tmpParagraph.StartSentence, &tmpParagraph.EndSentence)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
			return
		}

		// Get sentences of the story
		sentenceStmt, err := db.Query(fmt.Sprintf("Select word from sentence where sentence_id >= %d and sentence_id < %d", resStory.StartParagraph, resStory.EndParagraph))
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

	result, err := json.Marshal(&resStory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(string(result))
}
