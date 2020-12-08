package story

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"CollaborativeStory/colab/models"
	"CollaborativeStory/database"
)

// GetStory gets the specific story basis story_id
func GetStory(s database.StoryDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storyID := mux.Vars(r)["id"]

		var resStory models.DetailedStory

		// Get story details from story table
		resStory, err := s.FetchStory(storyID)
		if err != nil {
			log.Printf("error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
			return
		}
		if resStory.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(`{'error': 'story does not exist'}`)
			return
		}

		// Get paragraph details from paragraph table
		if resStory.StartParagraph > 0 {
			var isParagraphComplete, isSentenceComplete bool
			// Consideration for unfinished story
			if resStory.EndParagraph == 0 {
				isParagraphComplete = false
			} else {
				isParagraphComplete = true
			}

			allParagraphs, err := s.FetchParagraphs(resStory.StartParagraph, resStory.EndParagraph, isParagraphComplete)
			if err != nil {
				log.Printf("error: %v", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
				return
			}

			for _, pg := range allParagraphs {
				if pg.EndSentence == 0 {
					isSentenceComplete = false
				} else {
					isSentenceComplete = true
				}

				pg.Sentences, err = s.FetchSentences(pg.StartSentence, pg.EndSentence, isSentenceComplete)
				if err != nil {
					log.Printf("error: %v", err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
					return
				}

				resStory.Paragraphs = append(resStory.Paragraphs, pg)
			}
		}

		result, err := json.Marshal(&resStory)
		if err != nil {
			log.Printf("error: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(string(result))
	}

}
