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
	// POST method to add word to story
	mux.HandleFunc("/add", story.AddToStory)
	fmt.Println("Serving on :9000")
	log.Fatal(http.ListenAndServe(":9000", mux))
}
