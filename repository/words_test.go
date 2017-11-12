package repository_test

import (
	"testing"

	"github.com/aalempijevic/communityteaminterview/config"
	"log"
	"database/sql"
	"github.com/aalempijevic/communityteaminterview/repository"
	"fmt"
)
//
//func TestWordRepo_StoreWords(t *testing.T) {
//	fmt.Println("load config")
//	config, err := config.LoadDatabaseConfig("../appconfig.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Open conn")
//	database, err := sql.Open(config.DriverName, config.ConnectionString)
//	defer database.Close()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println("write words")
//	words := []string {"kazoo", "banana"}
//	repository.NewWordRepo(database).StoreWords(words)
//}


func TestWordRepo_GetSentenceTags(t *testing.T) {
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
	sentences := repo.GetSentenceTags(.7)
	fmt.Println(len(sentences))
}

func TestWordRepo_GetWordsByTag(t *testing.T) {
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

	fmt.Println("get words")
	repo := repository.NewWordRepo(database)
	words := repo.GetWordsByTag(2)
	fmt.Println(words)

}

