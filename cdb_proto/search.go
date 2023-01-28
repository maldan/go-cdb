package cdb_proto

import (
	"fmt"
	"reflect"
)

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
		size, fieldLen, startData := readHeader(m, offset, fieldOffsetIndex)
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
