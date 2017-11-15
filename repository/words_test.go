package repository_test

import (
	"github.com/aalempijevic/communityteaminterview/repository"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestWordRepo_GetWordsByTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	tagId := 5
	rows := sqlmock.NewRows([]string{"word"}).
		AddRow("some").
		AddRow("words")

	mock.ExpectQuery("select w.word from words").WithArgs(tagId).WillReturnRows(rows)

	words, err := repository.NewWordRepo(db).GetWordsByTag(tagId)
	if err != nil {
		t.Errorf("Unexpected error while executing GetWordsByTag: %s", err)
	}

	// Make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}

	assert.Equal(t, []string{"some", "words"}, words)
}
