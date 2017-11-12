package config

import (
	"encoding/json"
	"os"
)

//DatabaseConfig represents the configuration needed to open a database connection
type DatabaseConfig struct {
	DriverName			string 	`json:"driver"`
	ConnectionString	string	`json:"connection_string"`
}

//LoadDatabaseConfig loads configuration for database connection from a file in JSON format
func LoadDatabaseConfig(path string) (DatabaseConfig, error) {
	file, err := os.Open(path)
	decoder := json.NewDecoder(file)
	appConfig := DatabaseConfig{}
	err = decoder.Decode(&appConfig)

	return appConfig, err
}
