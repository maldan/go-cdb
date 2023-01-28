package pack

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"reflect"
	"strconv"
)

/**
Header struct

[P R O T 1 2 3 4] - file id
[1] - version

[0 0 0 0 0 0 0 0] - auto increment
[0 0 0 0 0 0 0 0] - total records

[...] - struct info
*/

/**
Struct info
[0] - total fields
[
	[0] - id/index
	[0] - type
	[0 0 0 0] - slice max length/capacity
	[0] - name length
	[...] - name
] * totalFields
*/

func EmptyHeader[T any]() []byte {
	bytes := make([]byte, core.HeaderSize)

	offset := 0

	// File id
	copy(bytes, "PROT1234")

	// Version
	offset += 8
	bytes[offset] = 1
	offset += 1

	// AI and Total
	offset += 8
	offset += 8

	// Struct info
	typeOf := reflect.TypeOf(*new(T))

	// Num of fields
	bytes[offset] = uint8(typeOf.NumField())
	offset += 1

	// Fields
	for i := 0; i < typeOf.NumField(); i++ {
		// Set id
		if typeOf.Field(i).Tag.Get("id") == "" {
			panic(fmt.Sprintf(
				"field %v doesn't have id tag",
				typeOf.Field(i).Name,
			))
		}
		id, err := strconv.Atoi(typeOf.Field(i).Tag.Get("id"))
		if err != nil {
			panic(err)
		}
		bytes[offset] = uint8(id)
		offset += 1

		// Set type
		switch typeOf.Field(i).Type.Kind() {
		case reflect.Bool:
			bytes[offset] = core.TBool
			break
		case reflect.Uint8, reflect.Int8:
			bytes[offset] = core.T8
			break
		case reflect.String:
			bytes[offset] = core.TString
			break
		default:
			panic(fmt.Sprintf(
				"unsuported type %v for field %v",
				typeOf.Field(i).Type.Name(),
				typeOf.Field(i).Name,
			))
		}
		offset += 1

		// Set slice capacity
		if typeOf.Field(i).Tag.Get("len") != "" {
			maxLen, err := strconv.Atoi(typeOf.Field(i).Tag.Get("len"))
			if err != nil {
				panic(err)
			}
			binary.LittleEndian.PutUint32(bytes[offset:], uint32(maxLen))
			offset += 4
		} else {
			binary.LittleEndian.PutUint32(bytes[offset:], 0)
			offset += 4
		}

		// Set name
		bytes[offset] = uint8(len(typeOf.Field(i).Name))
		offset += 1
		copy(bytes[offset:], typeOf.Field(i).Name)
		offset += len(typeOf.Field(i).Name)
	}

	return bytes
}
