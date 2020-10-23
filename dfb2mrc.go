package main

import (
	"./marc"
	"fmt"
	"github.com/LindsayBradford/go-dbf/godbf"
	"os"
	"runtime"
	"strconv"
)

func main() {
	dirName := "D:\\data\\ec5_base\\BOOK\\"
	if len(os.Args) > 1 {
		dirName = os.Args[1]
	}

	records := []marc.Record{}
	for mc := 2; mc < 11; mc++ {
		fileName := "MC" + strconv.Itoa(mc) + ".DBF"
		dbfTable, err := godbf.NewFromFile(dirName+fileName, "CP866")
		if err != nil {
			fmt.Println(err)
			fmt.Println("Ошибка открытия!")
			continue
		}

		for i := 0; i < dbfTable.NumberOfRecords(); i++ {
			marcRec := ""
			for fieldNumber := 1; fieldNumber <= mc; fieldNumber++ {
				part, _ := dbfTable.FieldValueByName(i, "MF"+strconv.Itoa(fieldNumber))
				marcRec += part
			}
			marcRecord := marc.NewMarcRecord(marcRec)
			records = append(records, *marcRecord)
		}
	}

	catalog := marc.CreateCatalog("Книги", &records)
	records = nil
	println(len(catalog.Records))
	//for key, val := range catalog.Records {
	//	println(key, marc.BinRecordToString(val))
	//}
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
