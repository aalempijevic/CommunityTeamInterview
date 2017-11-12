package batch_test

import (
	"testing"
	//"github.com/stretchr/testify/assert"
	"github.com/aalempijevic/communityteaminterview/batch"
	"fmt"
	"github.com/aalempijevic/communityteaminterview/config"
	"log"
	"database/sql"
	"github.com/aalempijevic/communityteaminterview/repository"
)

var commentText = `Lorem ipsum dolor sit amet, consectetur adipiscing elit,
		sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
		Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
		aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate
		velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat
		non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

//
//func TestProcessComment(t *testing.T) {
//	text := "comment comment poorly #@^ A. poorly; !punctuated !@#%comment"
//	expected := batch.WordFrequencies{
//		batch.WordFrequency{Word: "comment", Frequency:3},
//		batch.WordFrequency{Word: "poorly", Frequency: 2},
//		batch.WordFrequency{Word:"punctuated", Frequency: 1}}
//	result := batch.GetWordFrequencies(text)
//	assert.Equal(t, expected, result)
//}
//
//func BenchmarkProcessComment( b *testing.B) {
//	for n := 0; n < b.N; n++ {
//		batch.GetWordFrequencies(commentText)
//	}
//}


func TestTruncateAndProcessWords(t *testing.T) {
	fmt.Println("load config")
	config, err := config.LoadDatabaseConfig("../appconfig.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Open conn")
	database, err := sql.Open(config.DriverName, config.ConnectionString)
	defer database.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("get sentences")

	repo := repository.NewWordRepo(database)
	batch.TruncateAndProcessWords(*repo, .65, false)
}