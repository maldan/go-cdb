package cdb_proto

import (
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"github.com/maldan/go-cdb/cdb_proto/pack"
	"github.com/maldan/go-cdb/cdb_proto/parse"
)

func (d *DataTable[T]) Query(query string) {
	q, _ := parse.Query[T](query)

	d.Select(q)
}

func OpStrCmp(left func() string, right func() string) bool {
	a := left()
	b := right()
	return a == b
}

func OpAnd(left func() any, right func() any) bool {
	a := left()
	b := right()
	return a == b
}

func (d *DataTable[T]) Select(query parse.QueryInfo) {
	offset := core.HeaderSize

	comparator := StrCmp

	fieldOffsetIndex := 0
	/*for i := 0; i < query.TypeInfo.NumField(); i++ {
		if typeOf.Field(i).Name == field {
			fieldOffsetIndex = i * 8
			break
		}
	}*/
	// strAsBytes := []byte(query.Condition[0].Value.(string))
	strAsBytes := []byte("AA")

	for {
		size, fieldLen, startData := pack.ReadHeader(d.mem, offset, fieldOffsetIndex)
		fieldData := d.mem[startData : startData+fieldLen]

		// fmt.Printf("%v\n", string(fieldData))

		isFound := comparator(strAsBytes, fieldData)

		if isFound {
			fmt.Printf("%v %v\n", "GAS!", string(fieldData))
			break
		}

		offset += size
		if offset >= len(d.mem) {
			break
		}
	}
}
