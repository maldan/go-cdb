package pack

import (
	"encoding/binary"
	"github.com/maldan/go-cmhp/cmhp_byte"
	"reflect"
)

const _hStart = 2
const _hSize = 4
const _hTotal = 1
const _hEnd = 2

func ReadHeader(m []byte, offset int, fieldOffsetIndex int) (int, int, int) {
	size := int(binary.LittleEndian.Uint32(m[offset:]))

	localOffset := offset + 5 + fieldOffsetIndex

	fieldOffset := int(binary.LittleEndian.Uint32(m[localOffset:]))
	fieldLen := int(binary.LittleEndian.Uint32(m[localOffset+4:]))

	return size, fieldLen, offset + fieldOffset
}

func Pack[T any](v T) []byte {
	size := _hStart + _hSize + _hTotal + _hEnd

	// Get types
	typeOf := reflect.TypeOf(v)
	valueOf := reflect.ValueOf(v)

	// Total amount of fields
	total := uint8(typeOf.NumField())

	// Calculate size of data
	sizeOfEachPart := make([]int, 0)
	for i := 0; i < typeOf.NumField(); i++ {
		size += 4 // Size of field

		if typeOf.Field(i).Type.Kind() == reflect.String {
			size += 4 // Length of string
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
	cmhp_byte.From32ToBuffer(&size, &out, _hStart)

	// Set number of fields
	out[_hStart+_hSize] = total

	// Set offsets and size
	firstOffset := int(_hStart + _hSize + _hTotal + (8 * total))
	for i := 0; i < typeOf.NumField(); i++ {
		cmhp_byte.From32ToBuffer(&firstOffset, &out, _hStart+_hSize+_hTotal+(i*8))
		cmhp_byte.From32ToBuffer(&sizeOfEachPart[i], &out, _hStart+_hSize+_hTotal+(i*8)+4)
		firstOffset += sizeOfEachPart[i]
	}

	// Copy data
	firstOffset = int(_hStart + _hSize + _hTotal + (8 * total))
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
