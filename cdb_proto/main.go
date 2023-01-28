package cdb_proto

import (
	"github.com/edsrzf/mmap-go"
	"os"
)

type StructInfo struct {
	FieldCount    int
	FieldNameToId map[string]int
}

type DataTable[T any] struct {
	mem        mmap.MMap
	file       *os.File
	structInfo StructInfo

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
