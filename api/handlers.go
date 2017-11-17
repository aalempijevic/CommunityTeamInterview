package api

import (

	"database/sql"
	"encoding/json"
	"github.com/aalempijevic/communityteaminterview/repository"
	"github.com/gorilla/mux"
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

//WordsByTag gives a response containing words appearing with the specified tag sorted by highest frequency
func WordsByTag(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	tagId, err := strconv.Atoi(vars["tag"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Was not able to convert tag parameter to an int"))
		return
	}

	words, err := repository.NewWordRepo(database).GetWordsByTag(tagId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to get wordsByTag from repos"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(words)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to encode words"))
		return
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to get comments from repo"))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while trying to encode comments"))
		return
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
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Exception parsing float from threshold parameter"))
			return
		}
		go func() {
			batchRunning.Set()
			batch.TruncateAndProcessWords(*repository.NewWordRepo(database), threshold)
			batchRunning.UnSet()
		}()

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Processing request"))
	}
}
