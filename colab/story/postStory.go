package story

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"CollaborativeStory/colab/models"
	"CollaborativeStory/database"
)

// PostStory creates and updates story
func PostStory(w http.ResponseWriter, r *http.Request) {
	var reqWord map[string]string

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
	story, err := database.GetLatestStory()
	if err != nil {
		log.Println(err.Error())
	}

	var storyResp models.PostResponse
	if story.Title == "" {
		// Add title word to the new story
		err = database.AddStory(story.ID+1, word, true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
			return
		}

		storyResp.ID = story.ID + 1
		storyResp.Title = word
		storyResp.CurrentSentence = ""
	} else {
		if story.Title != "" && len(strings.Split(story.Title, " ")) < 2 {
			// Update title of the story
			err = database.AddStory(story.ID, word, false)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
				return
			}

			storyResp.ID = story.ID
			storyResp.Title = story.Title + " " + word
			storyResp.CurrentSentence = ""
		} else {
			// Add word to sentence of the story
			sentenceID, err := database.AddToSentence(word)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
				return
			}

			// Add/Update paragraph of the story
			paragraphID, err := database.AddToParagraph(sentenceID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
				return
			}

			if story.StartParagraph == 0 {
				// Start a new story
				err = database.UpdateStoryParagraph(story.ID, paragraphID, false)
				if err != nil {
					log.Println(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
					return
				}
			}

			if story.StartParagraph != 0 && (paragraphID-story.StartParagraph) == 7 {
				// End the story
				err = database.UpdateStoryParagraph(story.ID, paragraphID, true)
				if err != nil {
					log.Println(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
					return
				}
			}

			// Update the story timestamp (updated_at)
			err = database.UpdateStoryTimestamp(story.ID)
			if err != nil {
				log.Println(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
				return
			}
			storyResp.ID = story.ID
			storyResp.Title = story.Title
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

}
