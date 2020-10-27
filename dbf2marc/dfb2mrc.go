package dbf2marc

import (
	"../marc"
	"encoding/json"
	"log"
	"os"
)

type DbfCatalog struct {
	Name    string
	DbfPath string
}

type Config struct {
	DbfCatalogs []DbfCatalog
}

func GetCats(configFile string) (map[string](*marc.Catalog), error) {
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal("Не могу открыть config: ", err)
		return nil, err
	}

	decoder := json.NewDecoder(file)
	config := new(Config)
	err = decoder.Decode(&config)
	cats := make(map[string](*marc.Catalog))
	for _, catalog := range config.DbfCatalogs {
		cats[catalog.Name] = marc.CreateCatalogFromDBFNew(catalog.Name, catalog.DbfPath)
	}
	return cats, nil
}
