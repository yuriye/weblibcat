package main

import (
	"./marc"
	"fmt"
	"github.com/LindsayBradford/go-dbf/godbf"
	"os"
	"sort"
	"strconv"
)

func main() {
	dirName := "D:\\data\\ec5_base\\BOOK\\"
	if len(os.Args) > 1 {
		dirName = os.Args[1]
	}
	statistic := make(map[string]uint)
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
			for _, entry := range marcRecord.Entries {
				statistic[entry.Tag]++
			}
		}
	}
	keys := []string{}
	for key, _ := range statistic {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Println(key, "=>", statistic[key])
	}
}
