package pack

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"github.com/maldan/go-cmhp/cmhp_byte"
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
	if bytes[0] != 0x12 || bytes[1] != 0x34 {
		return out, errors.New("unknown header")
	}
	offset += core.RecordStart

	// @TODO Check footer

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
		fieldOff := int(binary.LittleEndian.Uint32(bytes[offset:]))
		fieldLen := int(binary.LittleEndian.Uint32(bytes[offset+4:]))
		fieldType := st.FieldType[i]

		/*fmt.Printf("OFF: %v\n", fieldOff)
		fmt.Printf("LN: %v\n", fieldLen)
		fmt.Printf("TT: %v\n", fieldType)*/

		// Read data from table
		if fieldType == core.TString {
			blob := bytes[fieldOff : fieldOff+fieldLen]
			*(*[]byte)(unsafe.Add(start, st.FieldOffset[i])) = make([]byte, len(blob))
			copy(*(*[]byte)(unsafe.Add(start, st.FieldOffset[i])), blob)
		}
		offset += core.RecordLenOff
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
		size += core.RecordLenOff // Offset to field data + size of field data

		if typeOf.Field(i).Type.Kind() == reflect.String {
			size += valueOf.Field(i).Len()
			sizeOfEachPart = append(sizeOfEachPart, valueOf.Field(i).Len())
		}
	}

	// Allocate buffer
	out := make([]byte, size)

	// Set header
	out[0] = 0x12
	out[1] = 0x34

	// Set end
	out[size-2] = 0x56
	out[size-1] = 0x78

	// Set size
	cmhp_byte.From32ToBuffer(&size, &out, core.RecordStart)

	// Set flags
	out[core.RecordStart+core.RecordSize] = uint8(total)

	// Short name
	ssf := core.RecordStart + core.RecordSize + core.RecordFlags

	// Set offsets and size
	firstOffset := ssf + (core.RecordLenOff * total)
	for i := 0; i < typeOf.NumField(); i++ {
		cmhp_byte.From32ToBuffer(&firstOffset, &out, ssf+(i*core.RecordLenOff))
		cmhp_byte.From32ToBuffer(&sizeOfEachPart[i], &out, ssf+(i*core.RecordLenOff)+4)
		firstOffset += sizeOfEachPart[i]
	}

	// Copy data
	firstOffset = ssf + (core.RecordLenOff * total)
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Type.Kind() == reflect.String {
			// Copy data
			copy(out[firstOffset:], valueOf.Field(i).Interface().(string))

			// Jump to next cell
			firstOffset += valueOf.Field(i).Len()
		}
	}

	return out
}
