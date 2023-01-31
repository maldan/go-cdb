package pack

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"reflect"
	"unsafe"
)

/**
Record struct

[12] - start of struct.
[0 0 0 0] - size of struct. Make configurable
[0] - flags of struct
0      - is removed
0	   - reserved
000000 - up to 64 fields

// look up table
[
	[0 0 0 0] - offset to field data. Make configurable
] * totalFields

// data table
[
	[...] - field data
] * totalFields

[34] - end of struct
*/

/**
[12]  	  - start of struct.
[0 0 0 0] - size
[0]       - flag
[
  [0]     - id
  [...]   - data
] * totalFields
[34]      - end of struct
*/

/*const _hStart = 2
const _hSize = 4
const _hFlags = 4
const _hEnd = 2*/

/*func ReadHeader(m []byte, offset int, fieldOffsetIndex int) (int, int, int) {
	// Skip header
	// offset += _hStart

	size := int(binary.LittleEndian.Uint32(m[offset+_hStart:]))

	localOffset := offset + _hStart + 5 + fieldOffsetIndex

	fieldOffset := int(binary.LittleEndian.Uint32(m[localOffset:]))
	fieldLen := int(binary.LittleEndian.Uint32(m[localOffset+4:]))

	return size, fieldLen, offset + fieldOffset
}*/

func ReadOffsetTable(startOfRecord []byte) []byte {
	// size := int(binary.LittleEndian.Uint32(startOfRecord[core.RecordStart:]))
	total := int(startOfRecord[core.RecordStart+core.RecordSize])
	return startOfRecord[core.RecordStart+core.RecordSize+core.RecordFlags : core.RecordStart+core.RecordSize+core.RecordFlags+total*core.RecordLenOff]
}

func ReadHeader2(m []byte, offset int) (int, []byte) {
	size := int(binary.LittleEndian.Uint32(m[offset+core.RecordStart:]))
	total := int(m[offset+core.RecordStart+core.RecordSize]) & core.MaskTotalFields
	return size, m[offset+core.RecordStart+core.RecordSize+core.RecordFlags : offset+core.RecordStart+core.RecordSize+core.RecordFlags+total*8]
}

func Unpack[T any](st core.StructInfo, bytes []byte) (T, error) {
	out := *new(T)
	offset := 0

	// Check header
	if bytes[0] != core.RecordStartMark {
		return out, errors.New("unknown header")
	}
	offset += core.RecordStart

	// Read size
	size := binary.LittleEndian.Uint32(bytes[offset:])
	offset += core.RecordSize
	fmt.Printf("%v\n", size)

	// Read total
	total := int(bytes[offset]) & core.MaskTotalFields
	offset += core.RecordFlags

	// Start of struct
	start := unsafe.Pointer(&out)

	// Read table
	for i := 0; i < total; i++ {
		/*fieldOff := int(binary.LittleEndian.Uint32(bytes[offset:]))
		fieldLen := int(binary.LittleEndian.Uint32(bytes[offset+4:]))
		fieldType := st.FieldType[i]

		// Read data from table
		if fieldType == core.TString {
			blob := bytes[fieldOff : fieldOff+fieldLen]
			*(*[]byte)(unsafe.Add(start, st.FieldOffset[i])) = make([]byte, len(blob))
			copy(*(*[]byte)(unsafe.Add(start, st.FieldOffset[i])), blob)
		}
		offset += core.RecordLenOff*/

		// Id
		// fieldId := bytes[offset]
		// offset += 1

		fieldType := st.FieldType[i]

		// Len
		fieldLen := int(bytes[offset])
		offset += 1

		// Copy
		if fieldType == core.TString {
			blob := bytes[offset : offset+fieldLen]
			*(*[]byte)(unsafe.Add(start, st.FieldOffset[i])) = make([]byte, len(blob))
			copy(*(*[]byte)(unsafe.Add(start, st.FieldOffset[i])), blob)
			offset += fieldLen
		}
	}

	return out, nil
}

func Pack[T any](v T) []byte {
	size := core.RecordStart + core.RecordSize + core.RecordFlags + core.RecordEnd

	// Get types
	typeOf := reflect.TypeOf(v)
	valueOf := reflect.ValueOf(v)

	// Total amount of fields
	total := typeOf.NumField()

	// Calculate size of data
	sizeOfEachPart := make([]int, 0)
	for i := 0; i < typeOf.NumField(); i++ {
		// size += core.RecordLenOff // Offset to field data + size of field data
		// size += 1 // field id
		if typeOf.Field(i).Type.Kind() == reflect.String {
			size += 1 // length of string
			size += valueOf.Field(i).Len()
			sizeOfEachPart = append(sizeOfEachPart, valueOf.Field(i).Len())
		}
	}

	// Allocate buffer
	out := make([]byte, size)

	// Set header
	out[0] = core.RecordStartMark

	// Set end
	out[size-1] = core.RecordEndMark

	// Set size
	binary.LittleEndian.PutUint32(out[core.RecordStart:], uint32(size))

	// Set flags
	out[core.RecordStart+core.RecordSize] = uint8(total)

	// Short name
	ssf := core.RecordStart + core.RecordSize + core.RecordFlags

	// Write fields
	startOffset := ssf
	for i := 0; i < typeOf.NumField(); i++ {
		//out[startOffset] = uint8(cmhp_convert.StrToInt(typeOf.Field(i).Tag.Get("id")))
		//startOffset += 1

		// Set str length
		binary.LittleEndian.PutUint32(out[startOffset:], uint32(sizeOfEachPart[i]))
		startOffset += 1

		// Copy content
		if typeOf.Field(i).Type.Kind() == reflect.String {
			// Copy data
			copy(out[startOffset:], valueOf.Field(i).Interface().(string))

			// Jump to next cell
			startOffset += valueOf.Field(i).Len()
		}
	}

	return out
}
