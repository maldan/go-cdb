package cdb_goson

import (
	"errors"
	"github.com/maldan/go-cdb/cdb_goson/goson"
)

func (d *DataTable[T]) Insert(v T) {
	bytes := goson.Marshal(v)

	bytes = wrap(bytes)

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
