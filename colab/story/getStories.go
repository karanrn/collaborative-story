package story

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"CollaborativeStory/database"
	"CollaborativeStory/helper"
)

type story struct {
	ID        int    `json:"ID"`
	Title     string `json:"Title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Allowed values for sort and order
var allowedSortBy = []string{"title", "created_at", "updated_at"}
var allowedOrdering = []string{"asc", "desc"}

// GetStories lists all the stories from the database
func GetStories(w http.ResponseWriter, r *http.Request) {
	var results []story
	db := database.DBConn()
	defer db.Close()

	var limit, offset int64
	var sort, order string
	var err error

	// Get query parameters
	query := r.URL.Query()
	// Convert values to int for limit and offset
	if query.Get("limit") != "" {
		limit, err = strconv.ParseInt(query.Get("limit"), 10, 64)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(`{'error': 'limit is not an integer'}`)
			return
		}
	} else {
		// Default limit for result set
		limit = 100
	}

	if query.Get("offset") != "" {
		offset, err = strconv.ParseInt(query.Get("offset"), 10, 64)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(`{'error': 'offset is not an integer'}`)
			return
		}
	}

	if query.Get("sort") != "" {
		if helper.Contains(query.Get("sort"), allowedSortBy) {
			sort = query.Get("sort")
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(fmt.Sprintf("{'error': 'sort should be among these values %v'}", allowedSortBy))
			return
		}
	} else {
		// Default sorting on created_at
		sort = allowedSortBy[1]
	}

	if query.Get("order") != "" {
		if helper.Contains(query.Get("order"), allowedOrdering) {
			sort = query.Get("sort")
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(fmt.Sprintf("{'error': 'order should be among these values %v'}", allowedOrdering))
			return
		}
	} else {
		// Default ordering is ascending
		order = allowedOrdering[0]
	}

	// Get all records
	storiesStmt, err := db.Query(fmt.Sprintf("Select story_id, title, created_at, updated_at from story order by %s %s limit %d, %d ", sort, order, offset, limit))
	if err != nil {
		log.Println(err.Error())
	}
	defer storiesStmt.Close()

	for storiesStmt.Next() {
		var iStory story
		err = storiesStmt.Scan(&iStory.ID, &iStory.Title, &iStory.CreatedAt, &iStory.UpdatedAt)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
			return
		}

		results = append(results, iStory)
	}

	resp, err := json.Marshal(results)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(`{'error': 'internal server error'}`)
		return
	}
	json.NewEncoder(w).Encode(fmt.Sprintf("{'limit': %d, 'offset': %d, 'count': %d, 'results': %v }", limit, offset, len(results), string(resp)))

}
