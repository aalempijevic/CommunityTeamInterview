package main

import (
	"log"
	"net/http"
	"database/sql"
	"github.com/aalempijevic/communityteaminterview/api"
	"github.com/aalempijevic/communityteaminterview/config"
)

//Environment contains environment specific vars such as the database connection
type Environment struct {
	db *sql.DB
}

func main() {
	env := initEnvironment()
	defer env.db.Close()

	api.Init(env.db)
	router := api.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

//InitEnvironment will configure our database connection and any other environment specific vars
func initEnvironment() Environment {
	config, err := config.LoadDatabaseConfig("appconfig.json")
	if err != nil {
		log.Fatal(err)
	}

	database, err := sql.Open(config.DriverName, config.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	return Environment{db: database}
}