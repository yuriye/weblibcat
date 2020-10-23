package marc

import (
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
	//binRecord := makeBinRecord(*record)
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

func MakeBinField(field *Field) (*BinField, error) {
	//var binField *BinField
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

func BinRecordToString(binRecord BinRecord) string {
	result := "\nID:" + binRecord.ID
	for _, field := range binRecord.Fields {
		result += "\n" + field.Tag + ": "
		result += "inds:"
		for _, ind := range field.Indicators {
			result += ind
		}
		for _, subfield := range field.Subfields {
			result += "\n\t\tsTag:" + subfield.Tag + " Content:" + subfield.Content
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

func CreateCatalog(name string, records *[]Record) *Catalog {
	catalog := Catalog{}
	catalog.Name = name
	catalog.Records = make(map[string]BinRecord)
	catalog.Indexes = make(map[string]Index)
	for _, record := range *records {
		binRecord := makeBinRecord(record)
		catalog.Records[binRecord.ID] = binRecord
	}
	return &catalog
}
