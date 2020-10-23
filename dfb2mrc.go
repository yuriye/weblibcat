package main

import (
	"./marc"
	"fmt"
	"runtime"
)

func main() {
	var directory string

	directory = "D:\\data\\ec5_base\\BOOK\\"
	catalogBooks := marc.CreateCatalogFromDBF("Книги", directory)
	//for _, rec := range catalog.Records {
	//	println(rec.String())
	//}
	println(len(catalogBooks.Records))
	PrintMemUsage()

	directory = "D:\\data\\ec5_base\\BIBS\\"
	catalogBiblio := marc.CreateCatalogFromDBF("Библиография", directory)
	println(len(catalogBiblio.Records))
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
