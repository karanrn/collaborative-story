package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"CollaborativeStory/colab/story"
	"CollaborativeStory/database"
)

func main() {
	// Initialize database
	sDB := database.StoryDB{}
	database.InitDB(&sDB)
	// HTTP multiplexer/router
	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/add", story.PostStory(&sDB)).Methods("POST")               // POST method to add word to story
	router.HandleFunc("/stories", story.GetStories(&sDB)).Methods("GET")           // GET method to list/get all stories
	router.HandleFunc("/stories/{id:[0-9]+}", story.GetStory(&sDB)).Methods("GET") // GET to get specific story
	fmt.Println("Serving on :9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}
