package main

import (
	"./marc"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
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
		//cats[catalog.Name] = marc.CreateCatalogFromDBF(catalog.Name, catalog.DbfPath)
		cats[catalog.Name] = marc.CreateCatalogFromDBFNew(catalog.Name, catalog.DbfPath)
	}
	return cats, nil
}

func main() {
	cats, err := GetCats("config.json")
	if err != nil {
		log.Panicf("Config file error:", err)
		return
	}
	for key, value := range cats {
		println(key, len(value.Records))
	}
	PrintMemUsage()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
