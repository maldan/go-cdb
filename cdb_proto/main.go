package cdb_proto

import (
	"github.com/edsrzf/mmap-go"
	"os"
)

const headerMaxParallelOffset = 64
const headerSize = 256

type DataTable[T any] struct {
	mem  mmap.MMap
	file *os.File

	Name string
}

func New[T any](name string) *DataTable[T] {
	d := DataTable[T]{Name: name}
	d.open()
	d.remap()
	return &d
}
