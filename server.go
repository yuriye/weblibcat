package main

import (
	"./marc"
	"./util"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "4132"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"find",
		"POST",
		"/find",
		find,
	},
}

type CatalogItem struct {
	Author          string `json:"author"`
	Title           string `json:"title"`
	ISBN            string `json:"isbn"`
	BBK             string `json:"bbk"`
	PublishingPlace string `json:"publishing_place"`
}

type CatalogItems []CatalogItem

var catalogItems []CatalogItem

func findByISBN(isbn string, pCat *marc.Catalog) *[]marc.BinRecord {
	result := []marc.BinRecord{}
	for _, record := range pCat.Records {
		if strings.ReplaceAll(isbn, "-", "") !=
			strings.ReplaceAll(record.GetISBN(), "-", "") {
			continue
		}
		result = append(result, record)
	}
	return &result
}

func find(w http.ResponseWriter, r *http.Request) {
	catalogItem := CatalogItem{}
	err := json.NewDecoder(r.Body).Decode(&catalogItem)
	if err != nil {
		log.Print("error occurred while decoding catalogItem data :: ", err)
		return
	}
	log.Println(catalogItem.ISBN, catalogItem.Author, catalogItem.Title)
	var records *[]marc.BinRecord
	if catalogItem.ISBN != "" {
		records = findByISBN(catalogItem.ISBN, cats["Книги"])
	} else if catalogItem.Author != "" {

	} else {

	}
	catalogItems := []CatalogItem{}
	if records != nil {
		for _, binRecord := range *records {
			catalogItems = append(catalogItems,
				CatalogItem{
					ISBN:   binRecord.GetISBN(),
					Author: "auth",
					Title:  "title"})
		}
	}
	json.NewEncoder(w).Encode(catalogItems)
}

func AddRoutes(router *mux.Router) *mux.Router {
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

var cats map[string]*marc.Catalog

func GetCats(config *util.Config) (map[string]*marc.Catalog, error) {
	cats := make(map[string](*marc.Catalog))
	for _, catalog := range config.DbfCatalogs {
		cats[catalog.Name] = marc.CreateCatalogFromDBFNew(catalog.Name, catalog.DbfPath)
	}
	return cats, nil
}

func main() {
	var err error
	pConfig, err := util.GetConfig("config.json")
	cats, err := GetCats(pConfig)
	if err != nil {
		log.Panicf("Config file error:", err)
		return
	}
	for key, value := range cats {
		log.Println(key, len(value.Records))
		//util.PrintFieldsStatistis(*value)
	}

	//util.LogMemUsage()

	muxRouter := mux.NewRouter().StrictSlash(true)
	router := AddRoutes(muxRouter)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
	log.Println("Listenning at " + pConfig.Connection.Host + ":" + pConfig.Connection.Port)
	err = http.ListenAndServe(pConfig.Connection.Host+":"+pConfig.Connection.Port, router)
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}

}
