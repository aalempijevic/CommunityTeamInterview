package repository_test

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/aalempijevic/communityteaminterview/model"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/aalempijevic/communityteaminterview/repository"
)



func TestCommentRepo_GetCommentsByWord(t *testing.T) {

	var rows = sqlmock.NewRows([]string{"id", "text"}).
		AddRow("1", "Test text in a test comment").
		AddRow("2", "test")

	var expectedResult = model.Comments{
		model.Comment{1,  "Test text in a test comment"},
		model.Comment{2, "test"},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error when opening a stub database connection: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT c.id, c.text FROM comments c").WithArgs("test", "test", 0, 1).WillReturnRows(rows)

	comments, err := repository.NewCommentRepo(db).GetCommentsByWord("test", 0, 1)
	if err != nil {
		t.Errorf("Unexpected error while executing GetCommentsByWord: %s", err)
	}

	assert.Equal(t, expectedResult, comments)
}

