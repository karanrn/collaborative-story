package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"CollaborativeStory/colab/story"
)

func main() {
	// HTTP multiplexer/router
	mux := http.NewServeMux()

	mux.HandleFunc("/add", story.AddToStory)     // POST method to add word to story
	mux.HandleFunc("/stories", story.GetStories) // GET method to list/get all stories
	fmt.Println("Serving on :9000")
	log.Fatal(http.ListenAndServe(":9000", mux))
}
