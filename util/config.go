package util

import (
	"encoding/json"
	"log"
	"os"
)

type Connection struct {
	Host string
	Port string
}

type DbfCatalog struct {
	Name    string
	DbfPath string
}

type Config struct {
	DbfCatalogs []DbfCatalog
	Connection  Connection
}

func GetConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal("Не могу открыть config: ", err)
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := new(Config)
	err = decoder.Decode(&config)
	return config, err
}
