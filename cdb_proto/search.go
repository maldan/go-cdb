package cdb_proto

import (
	"encoding/binary"
	"github.com/maldan/go-cdb/cdb_proto/core"
)

// Query("SELECT * FROM table WHERE FirstName == 'Roman' AND LastName != 'Lox'")
// Query("SELECT Id, Phone FROM table WHERE FirstName LIKE '%Lox'")
// Query("UPDATE table SET FirstName='Gavno' WHERE FirstName == 'Roman' LIMIT 1")
// Query("INSERT INTO table ")

// Example Find("FirstName == $0 && LastName != $0", 0)
// Update({ Set: "", Values: []any{} }, "")

/*func (d *DataTable[T]) Find(m []byte, field string, v string) {
	offset := 0

	strAsBytes := []byte(v)
	comparator := StrCmp

	// Find field index
	typeOf := reflect.TypeOf(*new(T))

	fieldOffsetIndex := 0
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Name == field {
			fieldOffsetIndex = i * 8
			break
		}
	}

	for {
		size, fieldLen, startData := pack.ReadHeader(m, offset, fieldOffsetIndex)
		fieldData := m[startData : startData+fieldLen]

		isFound := comparator(strAsBytes, fieldData)

		if isFound {
			fmt.Printf("%v %v\n", "GAS!", string(fieldData))
			break
		}

		offset += size
		if offset >= len(m) {
			break
		}
	}
}*/

func (d *DataTable[T]) ForEach(fn func(offset int, size int) bool) {
	offset := core.HeaderSize

	for {
		size := int(binary.LittleEndian.Uint32(d.mem[offset+core.RecordStart:]))
		flags := int(d.mem[offset+core.RecordStart+core.RecordSize])

		// If field not deleted
		if flags&core.MaskDeleted != core.MaskDeleted {
			if !fn(offset, size) {
				break
			}
		}

		offset += size
		if offset >= len(d.mem) {
			break
		}
	}
}

func (d *DataTable[T]) Find(fieldList []string, where func(*T) bool) SearchResult[T] {
	// Return
	searchResult := SearchResult[T]{}

	mapper := ValueMapper[T]{}
	mapper.Map2(d.structInfo, fieldList)

	d.ForEach(func(offset int, size int) bool {
		mapper.Fill2(offset+core.RecordStart+core.RecordSize+core.RecordFlags, d.mem)

		if where(&mapper.Container) {
			searchResult.table = d
			searchResult.Result = append(searchResult.Result, Record{offset: offset, size: size})
			return false
		}

		return true
	})

	return searchResult
}
