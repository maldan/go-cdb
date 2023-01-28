package cdb_proto

import (
	"errors"
	"github.com/edsrzf/mmap-go"
	"github.com/maldan/go-cdb/cdb_proto/pack"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (d *DataTable[T]) open() {
	// Check if file exists
	if _, err := os.Stat(d.Name); errors.Is(err, fs.ErrNotExist) {
		// Create path for file
		err = os.MkdirAll(filepath.Dir(d.Name), 0777)
		if err != nil {
			panic(err)
		}

		// Init file, because 0 length file fails with memory mapping
		if err = ioutil.WriteFile(d.Name, pack.EmptyHeader[T](), 0777); err != nil {
			panic(err)
		}
	}

	// Open file
	f, err := os.OpenFile(d.Name, os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}

	d.file = f
}

func (d *DataTable[T]) remap() {
	// Unmap previous
	if d.mem != nil {
		err := d.mem.Unmap()
		if err != nil {
			panic(err)
		}
	}

	// Map new
	mem, err := mmap.Map(d.file, mmap.RDWR, 0)
	if err != nil {
		panic(err)
	}
	d.mem = mem
}

func (d *DataTable[T]) Insert(v T) {
	bytes := pack.Pack(v)

	// Get file size
	stat, err := d.file.Stat()
	if err != nil {
		panic(err)
	}

	// Write at end of file
	n, err := d.file.WriteAt(bytes, stat.Size())
	if err != nil {
		panic(err)
	}
	if n != len(bytes) {
		panic(errors.New("incomplete writing"))
	}
}
