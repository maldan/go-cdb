package cdb_goson

import (
	"github.com/maldan/go-cdb/cdb_goson/goson"
)

type Record struct {
	offset int
	size   int
}

type SearchResult[T any] struct {
	IsFound bool
	Count   int
	Result  []Record
	table   *DataTable[T]
}

func (s *SearchResult[T]) Unpack() []T {
	out := make([]T, 0)

	for i := 0; i < len(s.Result); i++ {
		r := s.Result[i]

		realData := unwrap(s.table.mem[r.offset : r.offset+r.size])
		v := goson.Unmarshall[T](realData, s.table.Header.IdToName)
		out = append(out, v)
	}

	return out
}

func (s *Record) Delete() bool {
	return true
}

func (s *Record) Update(field string, value any) bool {
	return true
}
