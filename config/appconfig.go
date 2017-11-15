package config

import (
	"encoding/json"
	"os"
)

//DatabaseConfig represents the configuration needed to open a database connection
type DatabaseConfig struct {
	DriverName       string `json:"driver"`
	ConnectionString string `json:"connection_string"`
}

//LoadDatabaseConfig loads configuration for database connection from a file in JSON format
func LoadDatabaseConfig(path string) (DatabaseConfig, error) {
	file, err := os.Open(path)

	appConfig := DatabaseConfig{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&appConfig)

	return appConfig, err
}
