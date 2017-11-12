package api

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/aalempijevic/communityteaminterview/repository"
	"github.com/gorilla/mux"
	"strconv"
	"log"
	"database/sql"
)

var commentRepo *repository.CommentRepo
var wordsRepo *repository.WordRepo

func Init(db *sql.DB) {
	commentRepo = repository.NewCommentRepo(db)
	wordsRepo = repository.NewWordRepo(db)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}


func WordsByTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	tagId, err := strconv.Atoi(vars["tag"])
	if err != nil {
		log.Fatal(err)
	}

	words := wordsRepo.GetWordsByTag(tagId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(words)
	if err != nil {
		log.Fatal(err)
	}
}

func CommentsByWord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	word := vars["word"]
	skip, _ := strconv.Atoi(vars["skip"])
	limit, _ := strconv.Atoi(vars["limit"])

	comments := commentRepo.GetCommentsByWord(word, skip, limit)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(comments)
	if err != nil {
		log.Fatal(err)
	}
}
