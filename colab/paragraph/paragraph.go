package paragraph

import (
	"log"

	"CollaborativeStory/database"
)

// AddToParagraph adds sentence to a paragraph
func AddToParagraph(sentenceID int) (paragraphID int, err error) {
	// DB Connection
	db := database.DBConn()
	defer db.Close()

	// Get the unfinished paragraph
	var startParagraph int
	paragraphStmt, err := db.Query("select IFNULL(paragraph_id, 0), IFNULL(start_sentence, 0) from paragraph where start_sentence is not NULL and end_sentence is NULL;")
	if err != nil {
		log.Println(err.Error())
		return paragraphID, err
	}
	defer paragraphStmt.Close()
	if paragraphStmt.Next() {
		err = paragraphStmt.Scan(&paragraphID, &startParagraph)
		if err != nil {
			log.Println(err.Error())
			return paragraphID, err
		}
	}

	// Get the max value
	if paragraphID == 0 {
		lastParagraphStmt, err := db.Query("select IFNULL(max(paragraph_id), 0) from paragraph;")
		if err != nil {
			log.Println(err.Error())
			return paragraphID, err
		}
		defer lastParagraphStmt.Close()
		if lastParagraphStmt.Next() {
			err = lastParagraphStmt.Scan(&paragraphID)
			if err != nil {
				log.Println(err.Error())
				return paragraphID, err
			}
		}
		// Create next paragraph
		paragraphID++
	}

	// Check the size of the paragraph and update/create paragraph
	updateParagraph, err := db.Prepare("update paragraph set end_sentence = ?, updated_at = current_timestamp() where paragraph_id = ?")
	if err != nil {
		log.Println(err.Error())
		return paragraphID, err
	}
	defer updateParagraph.Close()

	addParagraph, err := db.Prepare("insert into paragraph (paragraph_id, start_sentence, created_at) values (?, ?, current_timestamp())")
	if err != nil {
		log.Println(err.Error())
		return paragraphID, err
	}
	defer addParagraph.Close()

	if startParagraph != 0 && (sentenceID-startParagraph) == 10 {
		_, err = updateParagraph.Exec(sentenceID, paragraphID)
		if err != nil {
			log.Println(err.Error())
			return paragraphID, err
		}
	}

	if startParagraph == 0 {
		_, err = addParagraph.Exec(paragraphID, sentenceID)
		if err != nil {
			log.Println(err.Error())
			return paragraphID, err
		}
	}

	return paragraphID, nil
}
