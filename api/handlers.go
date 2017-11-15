package api

import (

	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aalempijevic/communityteaminterview/repository"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"github.com/aalempijevic/communityteaminterview/batch"

	"github.com/tevino/abool"
)

//CommentRepo is used to retrieve comment data
var database *sql.DB

//BatchRunning is an atomic bool used to make sure we do not start two simultaneous batches
//we check set it from a goroutine running the batch job and check it from requests
var batchRunning *abool.AtomicBool

//Init initializes the repositories
func Init(db *sql.DB) {
	database = db
	batchRunning = abool.New()
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}


//WordsByTag gives a response containing words appearing with the specified tag sorted by highest frequency
func WordsByTag(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	tagId, err := strconv.Atoi(vars["tag"])
	if err != nil {
		log.Fatalf("Was not able to convert tag parameter to an int: %s", err)
	}

	words, err := repository.NewWordRepo(database).GetWordsByTag(tagId)
	if err != nil {
		log.Fatalf("Error while trying to get wordsByTag from repo: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(words)
	if err != nil {
		log.Fatalf("Unable to encode words as JSON object: %s", err)
	}
}

//CommentsByWord gives a response containing all comments with text that matches the word sorted by relevance
func CommentsByWord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	word := vars["word"]
	skip, _ := strconv.Atoi(vars["skip"])
	limit, _ := strconv.Atoi(vars["limit"])

	comments, err := repository.NewCommentRepo(database).GetCommentsByWord(word, skip, limit)
	if err != nil {
		log.Fatalf("Unable to retrieve comments: %s", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		log.Fatal(err)
	}
}



//ProcessWords executes our word processing batch against the specified threshold value
func ProcessWords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if batchRunning.IsSet() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Batch already running."))
	} else {
		vars := mux.Vars(r)

		thresholdStr := vars["threshold"]
		threshold, err := strconv.ParseFloat(thresholdStr, 64)
		if err != nil {
			log.Fatalf("Exception parsing float from threshold parameter: %s", vars["threshold"])
		}
		go func() {
			batchRunning.Set()
			batch.TruncateAndProcessWords(*repository.NewWordRepo(database), threshold)
			batchRunning.UnSet()
		}()
		if err != nil {
			log.Fatalf("Unable to retrieve comments: %s", err)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Processing request"))
	}
}
