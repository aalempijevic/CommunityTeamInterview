package repository_test

import (
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/aalempijevic/communityteaminterview/repository"
	"github.com/aalempijevic/communityteaminterview/model"
	"github.com/stretchr/testify/assert"
)




func TestWordRepo_GetSentenceTags(t *testing.T) {
	var rows = sqlmock.NewRows([]string{"tagIds", "sentence"}).
		AddRow("1,3,5", "Sentence 1").
		AddRow("2,4", "Sentence 2")

	var expectedResults = model.Sentences{
		model.Sentence{TagIds: []int{1, 3, 5}, Text: "Sentence 1"},
		model.Sentence{TagIds: []int{2, 4}, Text: "Sentence 2"},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	threshold := .5
	mock.ExpectQuery("select tagIds, sentence from").WithArgs(threshold).WillReturnRows(rows)

	sentences, err := repository.NewWordRepo(db).GetSentenceTags(threshold)
	if err != nil {
		t.Errorf("Unexpected error while executing GetSentenceTags: %s", err)
	}

	assert.Equal(t, expectedResults, sentences)
}

