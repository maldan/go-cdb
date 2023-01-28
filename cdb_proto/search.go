package cdb_proto

import (
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto/pack"
	"reflect"
)

// Query("SELECT * FROM table WHERE FirstName == 'Roman' AND LastName != 'Lox'")
// Query("SELECT Id, Phone FROM table WHERE FirstName LIKE '%Lox'")
// Query("UPDATE table SET FirstName='Gavno' WHERE FirstName == 'Roman' LIMIT 1")
// Query("INSERT INTO table ")

// Example Find("FirstName == $0 && LastName != $0", 0)
// Update({ Set: "", Values: []any{} }, "")

func (d *DataTable[T]) Find(m []byte, field string, v string) {
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
}
