package cdb_proto

import (
	"encoding/binary"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"github.com/maldan/go-cdb/cdb_proto/pack"
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

func (d *DataTable[T]) ForEach(fn func(offset int) bool) {
	offset := core.HeaderSize

	for {
		size := int(binary.LittleEndian.Uint32(d.mem[offset+core.RecordStart:]))

		if !fn(offset) {
			break
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
	mapper.Map(d.structInfo, fieldList)

	/*cc := new(T)
	typeOf := reflect.TypeOf(cc).Elem()
	start := unsafe.Pointer(cc)

	idTable := make([]int, 0)
	typeTable := make([]int, 0)
	outOffset := make([]unsafe.Pointer, 0)

	// Fill id table
	for i := 0; i < len(fieldList); i++ {
		idTable = append(idTable, d.structInfo.FieldNameToId[fieldList[i]]*8)
		typeTable = append(typeTable, d.structInfo.FieldType[idTable[i]/8])
		f, ok := typeOf.FieldByName(fieldList[i])
		if !ok {
			panic("FNAF")
		}
		outOffset = append(outOffset, unsafe.Add(start, f.Offset))
	}*/

	d.ForEach(func(offset int) bool {
		size, offTable := pack.ReadHeader2(d.mem, offset)

		// Fill struct
		/*for i := 0; i < len(idTable); i++ {
			vOff := int(binary.LittleEndian.Uint32(offTable[idTable[i]:]))
			vLen := int(binary.LittleEndian.Uint32(offTable[idTable[i]+4:]))

			if typeTable[i] == core.TString {
				*(*[]byte)(outOffset[i]) = d.mem[offset+vOff : offset+vOff+vLen]
			}
		}*/

		mapper.Fill(offset, d.mem, offTable)

		if where(&mapper.Container) {
			searchResult.table = d
			searchResult.Result = append(searchResult.Result, Record{offset: offset, size: size})
			return false
		}

		return true
	})

	return searchResult
}
