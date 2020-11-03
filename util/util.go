package util

import (
	"../marc"
	"fmt"
	"log"
	"runtime"
)

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

type fieldsStatistics map[string]int

func PrintFieldsStatistis(cat marc.Catalog) {
	statistics := fieldsStatistics{}
	for _, record := range cat.Records {
		for _, field := range record.Fields {
			statistics[field.Tag] += 1
			if field.Tag == "090" || field.Tag == "091" || field.Tag == "653" || field.Tag == "700" {
				print(field.Tag + ":")
				for _, subfield := range field.Subfields {
					println(subfield.Tag, subfield.Content)
				}
			}
		}
	}
	for key, item := range statistics {
		println(key, item)
	}
}
