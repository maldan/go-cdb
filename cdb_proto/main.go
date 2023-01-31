package cdb_proto

import (
	"github.com/edsrzf/mmap-go"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"github.com/maldan/go-cdb/cdb_proto/pack"
	"os"
)

type DataTable[T any] struct {
	mem        mmap.MMap
	file       *os.File
	structInfo core.StructInfo

	// Table config
	RecordSizeInBytes uint8

	Name string
}

func (d *DataTable[T]) Debug__GetMem() mmap.MMap {
	return d.mem
}
func (d *DataTable[T]) Debug__GetStructInfo() core.StructInfo {
	return d.structInfo
}

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

func (s SearchResult[T]) Unpack() []T {
	out := make([]T, 0)

	for i := 0; i < len(s.Result); i++ {
		r := s.Result[i]
		v, _ := pack.Unpack[T](s.table.structInfo, s.table.mem[r.offset:r.offset+r.size])
		out = append(out, v)
	}

	return out
}

func New[T any](name string) *DataTable[T] {
	d := DataTable[T]{Name: name}
	d.structInfo.FieldNameToId = map[string]int{}
	d.open()
	d.remap()
	d.readHeader()
	return &d
}
