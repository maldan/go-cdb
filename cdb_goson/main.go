package cdb_goson

import (
	"github.com/edsrzf/mmap-go"
	"github.com/maldan/go-cdb/cdb_proto/core"
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

func New[T any](name string) *DataTable[T] {
	d := DataTable[T]{Name: name}
	d.structInfo.FieldNameToId = map[string]int{}
	d.open()
	d.remap()
	d.readHeader()
	return &d
}
