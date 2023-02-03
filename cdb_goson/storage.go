package cdb_goson

import (
	"errors"
	"github.com/edsrzf/mmap-go"
	"github.com/maldan/go-cdb/cdb_goson/core"
	"github.com/maldan/go-cdb/cdb_proto/pack"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
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

	/*fmt.Printf("| AI: %v | Total: %v | Version %v |\n",
		binary.LittleEndian.Uint64(d.mem[9:]),
		binary.LittleEndian.Uint64(d.mem[9+8:]),
		d.mem[8],
	)*/

	// Calculate field offset
	typeOf := reflect.TypeOf(*new(T))
	d.structInfo.FieldOffset = make([]uintptr, 64)
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Tag.Get("id") != "" {
			id, err := strconv.Atoi(typeOf.Field(i).Tag.Get("id"))
			if err != nil {
				panic(err)
			}
			d.structInfo.FieldOffset[id] = typeOf.Field(i).Offset
		}
	}

	offset := 8 + 1 + 8 + 8
	amount := int(d.mem[offset])
	//fmt.Printf("Fields: %v\n", amount)
	offset += 1
	for i := 0; i < amount; i++ {
		fieldId := d.mem[offset]
		//fmt.Printf("\tId: %v", fieldId)
		offset += 1
		fieldType := d.mem[offset]
		//fmt.Printf("\tType: %v", fieldType)
		offset += 1
		//fmt.Printf("\t\tCapacity: %v", binary.LittleEndian.Uint32(d.mem[offset:]))
		offset += 4

		fieldLen := int(d.mem[offset])
		// fmt.Printf("\tLen: %v\n", fieldLen)
		offset += 1
		fieldName := string(d.mem[offset : offset+fieldLen])
		//fmt.Printf("\t\tName: %v\n", fieldName)
		offset += fieldLen

		// Set map
		d.structInfo.FieldNameToId[fieldName] = int(fieldId)
		d.structInfo.FieldName = append(d.structInfo.FieldName, fieldName)
		d.structInfo.FieldType = append(d.structInfo.FieldType, int(fieldType))
		// d.structInfo.FieldOffset = append(d.structInfo.FieldType, typeOf)
	}

	d.structInfo.FieldCount = amount
}

func unwrap(bytes []byte) []byte {
	if bytes[0] != 0x12 {
		panic("non package")
	}
	hh := core.RecordStart + core.RecordSize + core.RecordFlags
	return bytes[hh : len(bytes)-1]
}

func wrap(bytes []byte) []byte {
	fullSize := len(bytes) + core.RecordStart + core.RecordSize + core.RecordFlags + core.RecordEnd
	fullPackage := make([]byte, 0, fullSize)

	// Start
	fullPackage = append(fullPackage, 0x12)

	// Size
	fullPackage = append(fullPackage, uint8(fullSize))
	fullPackage = append(fullPackage, uint8(fullSize>>8))
	fullPackage = append(fullPackage, uint8(fullSize>>16))
	fullPackage = append(fullPackage, uint8(fullSize>>24))

	// Flags
	fullPackage = append(fullPackage, 0)

	// Body
	fullPackage = append(fullPackage, bytes...)

	// End
	fullPackage = append(fullPackage, 0x34)

	return fullPackage
}
