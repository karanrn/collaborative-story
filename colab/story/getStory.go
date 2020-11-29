package story

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type storyParagraph struct {
	Sentences []string `json:"sentences"`
}

type detStory struct {
	ID         string           `json:"id"`
	Title      string           `json:"title"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
	Paragraphs []storyParagraph `json:"paragraphs"`
}

// GetStory gets the specific story basis story_id
func GetStory(w http.ResponseWriter, r *http.Request) {
	storyID := mux.Vars(r)["id"]

	fmt.Println(storyID)
}
