package main

import (
	"github.com/maldan/go-cmhp/cmhp_byte"
	"reflect"
)

const _hSize = 4
const _hTotal = 1

func Pack[T any](v T) []byte {
	size := _hSize + _hTotal
	total := uint8(0)
	sizeOfEachPart := make([]int, 0)

	typeOf := reflect.TypeOf(v)
	valueOf := reflect.ValueOf(v)
	for i := 0; i < typeOf.NumField(); i++ {
		size += 4

		if typeOf.Field(i).Type.Kind() == reflect.String {
			size += 4
			size += valueOf.Field(i).Len()
			sizeOfEachPart = append(sizeOfEachPart, 4+valueOf.Field(i).Len())
		}

		total += 1
	}

	// Allocate buffer
	out := make([]byte, size)

	// Set size
	cmhp_byte.From32ToBuffer(&size, &out, 0)

	// Set number of fields
	out[4] = total

	// Set offsets
	firstOffset := int(_hSize + _hTotal + (4 * total))
	for i := 0; i < typeOf.NumField(); i++ {
		cmhp_byte.From32ToBuffer(&firstOffset, &out, 5+(i*4))
		firstOffset += sizeOfEachPart[i]
	}

	// Copy data
	firstOffset = int(_hSize + _hTotal + (4 * total))
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Type.Kind() == reflect.String {
			// Copy size
			partSize := sizeOfEachPart[i] - 4
			cmhp_byte.From32ToBuffer(&partSize, &out, firstOffset)
			// Copy data
			copy(out[firstOffset+4:], valueOf.Field(i).Interface().(string))

			// Jump to next cell
			firstOffset += 4
			firstOffset += valueOf.Field(i).Len()
		}
	}

	return out
}
