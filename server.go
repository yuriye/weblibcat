package main

import (
	"./dbf2marc"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"runtime"
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
	ISBN   string `json:"ISBN"`
	Author string `json:"autor"`
	Title  string `json:"title"`
}

type CatalogItems []CatalogItem

var catalogItems []CatalogItem

func init() {
	catalogItems = CatalogItems{
		CatalogItem{ISBN: "1", Author: "Foo", Title: "Bar"},
		CatalogItem{ISBN: "2", Author: "Baz", Title: "Qux"},
	}
}

func find(w http.ResponseWriter, r *http.Request) {
	catalogItem := CatalogItem{}
	err := json.NewDecoder(r.Body).Decode(&catalogItem)
	if err != nil {
		log.Print("error occurred while decoding catalogItem data :: ", err)
		return
	}
	log.Println(catalogItem.ISBN, catalogItem.Author, catalogItem.Title)
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

func main() {
	cats, err := dbf2marc.GetCats("config.json")
	if err != nil {
		log.Panicf("Config file error:", err)
		return
	}
	for key, value := range cats {
		log.Println(key, len(value.Records))
	}
	LogMemUsage()

	muxRouter := mux.NewRouter().StrictSlash(true)
	router := AddRoutes(muxRouter)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
	log.Println("Listenning at " + CONN_HOST + ":" + CONN_PORT)
	err = http.ListenAndServe(CONN_HOST+":"+CONN_PORT, router)
	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}

}

func LogMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Println(fmt.Sprintf("Alloc = %v MiB", bToMb(m.Alloc)) +
		fmt.Sprintf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc)) +
		fmt.Sprintf("\tSys = %v MiB", bToMb(m.Sys)) +
		fmt.Sprintf("\tNumGC = %v", m.NumGC))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
