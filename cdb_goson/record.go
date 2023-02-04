package cdb_goson

import (
	"fmt"
	"github.com/maldan/go-cdb/cdb_goson/core"
	"github.com/maldan/go-cdb/cdb_goson/goson"
	"reflect"
)

type Record[T any] struct {
	offset int
	size   int
	table  *DataTable[T]
}

type SearchResult[T any] struct {
	IsFound bool
	Count   int
	Result  []Record[T]
	table   *DataTable[T]
}

func (s *SearchResult[T]) Unpack() []T {
	out := make([]T, 0)

	for i := 0; i < len(s.Result); i++ {
		r := s.Result[i]

		//realData := unwrap(s.table.mem[r.offset : r.offset+r.size])
		//v := goson.Unmarshall[T](realData, s.table.Header.IdToName)
		out = append(out, r.Unpack())
	}

	return out
}

func (s *Record[T]) Unpack() T {
	realData := unwrap(s.table.mem[s.offset : s.offset+s.size])
	return goson.Unmarshall[T](realData, s.table.Header.IdToName)
}

func (s *Record[T]) Delete() bool {
	if s.table.mem[s.offset] != core.RecordStartMark {
		return false
	}

	// Read flags
	b := []byte{0}
	s.table.file.ReadAt(b, int64(s.offset+core.RecordStart+core.RecordSize))
	fmt.Printf("FB: %v\n", b[0])
	b[0] |= core.MaskDeleted
	fmt.Printf("FAPL: %v\n", b[0])

	// Write back
	s.table.file.WriteAt(b, int64(s.offset+core.RecordStart+core.RecordSize))

	return true
}

func (s *Record[T]) Update(fields map[string]any) bool {
	// Unpack current
	unpack := s.Unpack()

	// Delete old
	s.Delete()

	// Fill fields
	valueOf := reflect.ValueOf(&unpack).Elem()
	for name, value := range fields {
		field := valueOf.FieldByName(name)
		if field.Kind() == reflect.Invalid {
			continue
		}
		if !field.CanSet() {
			continue
		}
		switch value.(type) {
		case string:
			field.SetString(value.(string))
			break
		default:
			break
		}
	}

	// Create new
	s.table.Insert(unpack)

	return true
}
