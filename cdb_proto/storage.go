package cdb_proto

import (
	"encoding/binary"
	"errors"
	"fmt"
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

func (d *DataTable[T]) readHeader() {
	fileId := d.mem[0:8]
	if string(fileId) != "PROT1234" {
		panic("non db")
	}

	fmt.Printf("| AI: %v | Total: %v | Version %v |\n",
		binary.LittleEndian.Uint64(d.mem[9:]),
		binary.LittleEndian.Uint64(d.mem[9+8:]),
		d.mem[8],
	)

	offset := 8 + 1 + 8 + 8
	amount := int(d.mem[offset])
	fmt.Printf("Fields: %v\n", amount)
	offset += 1
	for i := 0; i < amount; i++ {
		fieldId := d.mem[offset]
		fmt.Printf("\tId: %v", fieldId)
		offset += 1
		fmt.Printf("\tType: %v", d.mem[offset])
		offset += 1
		fmt.Printf("\t\tCapacity: %v", binary.LittleEndian.Uint32(d.mem[offset:]))
		offset += 4

		fieldLen := int(d.mem[offset])
		// fmt.Printf("\tLen: %v\n", fieldLen)
		offset += 1
		fieldName := string(d.mem[offset : offset+fieldLen])
		fmt.Printf("\t\tName: %v\n", fieldName)
		offset += fieldLen

		// Set map
		d.structInfo.FieldNameToId[fieldName] = int(fieldId)
	}

	d.structInfo.FieldCount = amount
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
