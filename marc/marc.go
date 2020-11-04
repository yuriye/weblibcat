package marc

import (
	"github.com/LindsayBradford/go-dbf/godbf"
	"os"
	"sort"
	"strconv"
	"strings"
)

type DirectoryEntry struct {
	Tag    string
	Data   string
	offset int
}

func (entry DirectoryEntry) NewDirectoryEntry(data string) DirectoryEntry {
	entry.Data = data
	entry.Tag = Substr(data, 0, 3)
	return entry
}

func (entry *DirectoryEntry) setEntryData(data string) {
	entry.Data = data
	entry.Tag = Substr(data, 0, 3)
	entry.offset, _ = strconv.Atoi(Substr(data, 7, 5))
}

type Field struct {
	Tag  string
	Data string
}

func (record *Record) appendField(entry *DirectoryEntry) {
	start, _ := strconv.Atoi(Substr(entry.Data, 7, 5))
	start += record.base
	len, _ := strconv.Atoi(Substr(entry.Data, 3, 4))
	fieldData := Substr(record.Data, start, len)
	field := Field{Data: fieldData, Tag: entry.Tag}
	record.Fields = append(record.Fields, field)
}

func (record *Record) fillFields() {
	fieldsArr := strings.Split(record.Data, "\x1e")
	for i, v := range record.Entries {
		field := Field{Data: fieldsArr[i+1], Tag: v.Tag}
		record.Fields = append(record.Fields, field)
	}
}

type Record struct {
	Data        string
	Entries     []DirectoryEntry
	Fields      []Field
	base        int
	fieldsCount int
}

func NewMarcRecord(data string) *Record {
	record := new(Record)
	record.Data = data
	record.base, _ = strconv.Atoi(Substr(record.Data, 12, 5))
	record.fieldsCount = (record.base - 25) / 12
	record.FillEnries()
	sort.SliceStable(record.Entries, func(i, j int) bool {
		return record.Entries[i].offset < record.Entries[j].offset
	})
	record.fillFields()
	return record
}

func Substr(input string, start int, length int) string {
	asRunes := []rune(input)
	if start >= len(asRunes) {
		return ""
	}
	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}
	return string(asRunes[start : start+length])
}

func (record *Record) getBase() int {
	base, _ := strconv.Atoi(Substr(record.Data, 12, 5))
	return base
}

func (record *Record) GetDirectory() string {
	dir := Substr(record.Data, 24, record.getBase()-25)
	return dir
}

func (record *Record) FillEnries() {
	record.Entries = []DirectoryEntry{}
	base := record.getBase()
	fieldsCount := (base - 25) / 12
	offset := 24
	for index := 0; index < fieldsCount; index++ {
		entryData := Substr(record.Data, offset+index*12, 12)
		entry := DirectoryEntry{}
		entry.setEntryData(entryData)
		record.Entries = append(record.Entries, entry)
	}
}

type Subfield struct {
	Tag     string
	Content string
}

type BinField struct {
	Tag        string
	Indicators [2]string
	Subfields  []Subfield
}

type BinRecord struct {
	ID     string
	Fields []BinField
}

func (record *BinRecord) GetISBN() string {
	for _, field := range record.Fields {
		if field.Tag == "020" {
			for _, subfield := range field.Subfields {
				if subfield.Tag == "a" {
					return subfield.Content
				}
			}
		}
	}
	return ""
}

func (record *BinRecord) GetAuthor() string {
	for _, field := range record.Fields {
		if field.Tag == "100" {
			for _, subfield := range field.Subfields {
				if subfield.Tag == "a" {
					return subfield.Content
				}
			}
		}
	}
	return ""
}

func (record *BinRecord) GetTitle() string {
	for _, field := range record.Fields {
		if field.Tag == "245" {
			for _, subfield := range field.Subfields {
				if subfield.Tag == "a" {
					return strings.Trim(subfield.Content, " ")
				}
			}
		}
	}
	return ""
}

func MakeBinField(field *Field) (*BinField, error) {
	binField := new(BinField)
	if Substr(field.Tag, 0, 2) == "00" {
		binField.Tag = field.Tag
		data := strings.Split(Substr(field.Data, 0, len(field.Data)), "\x1f")[0]
		data = strings.Trim(data, " ")
		subfield := Subfield{Tag: "", Content: data}
		binField.Subfields = append(binField.Subfields, subfield)
		return binField, nil
	}

	binField.Tag = field.Tag
	binField.Indicators[0] = Substr(field.Data, 0, 1)
	binField.Indicators[1] = Substr(field.Data, 1, 1)
	subs := strings.Split(Substr(field.Data, 2, len(field.Data)-2), "\x1f")
	for _, s := range subs {
		subfield := Subfield{Tag: Substr(s, 0, 1), Content: Substr(s, 1, len(s)-1)}
		if subfield.Tag == "" && subfield.Content == "" {
			continue
		}
		binField.Subfields = append(binField.Subfields, subfield)
	}
	return binField, nil
}

func makeBinRecord(record Record) BinRecord {
	binRecord := BinRecord{}
	for _, field := range record.Fields {
		binField, err := MakeBinField(&field)
		if err != nil {
			continue
		}
		if binField.Tag == "001" {
			binRecord.ID = binField.Subfields[0].Content
		}
		binRecord.Fields = append(binRecord.Fields, *binField)
	}
	return binRecord
}

func (binRecord BinRecord) String() string {
	result := ""
	for _, field := range binRecord.Fields {
		if field.Tag == "001" ||
			field.Tag == "005" ||
			field.Tag == "008" ||
			field.Tag == "520" ||
			field.Tag == "650" ||
			field.Tag == "773" ||
			field.Tag == "856" ||
			field.Tag == "952" {
			continue
		}
		result += "   "
		for _, subfield := range field.Subfields {
			result += subfield.Content + " "
		}
	}
	return result
}

type Index struct {
	Name  string
	Items map[string][]string
}

type Catalog struct {
	Name    string
	Records map[string]BinRecord
	Indexes map[string]Index
}

func CreateCatalogFromDBFNew(name string, dirName string) *Catalog {
	var err error
	catalog := Catalog{}
	catalog.Name = name
	catalog.Records = make(map[string]BinRecord)
	catalog.Indexes = make(map[string]Index)

	for mc := 2; mc < 21; mc++ {
		fileName := "MC" + strconv.Itoa(mc) + ".DBF"
		dbfName := dirName + fileName

		_, err = os.Stat(dbfName)
		if os.IsNotExist(err) {
			continue
		}

		dbfTable, err := godbf.NewFromFile(dbfName, "CP866")
		if err != nil {
			continue
		}

		for i := 0; i < dbfTable.NumberOfRecords(); i++ {
			marcRec := ""
			for fieldNumber := 1; fieldNumber <= mc; fieldNumber++ {
				part, _ := dbfTable.FieldValueByName(i, "MF"+strconv.Itoa(fieldNumber))
				marcRec += part
			}
			binRecord := makeBinRecord(*NewMarcRecord(marcRec))
			catalog.Records[binRecord.ID] = binRecord
		}
	}
	return &catalog
}
