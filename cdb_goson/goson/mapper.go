package goson

import (
	"encoding/binary"
	"github.com/maldan/go-cdb/cdb_goson/core"
	"reflect"
	"unsafe"
)

type ValueMapper[T any] struct {
	Container T

	MapOffset map[string]unsafe.Pointer
}

func NewMapper[T any]() *ValueMapper[T] {
	mapper := ValueMapper[T]{
		Container: *new(T),
		MapOffset: map[string]unsafe.Pointer{},
	}

	typeOf := reflect.TypeOf(mapper.Container)
	start := unsafe.Pointer(&mapper.Container)

	for i := 0; i < typeOf.NumField(); i++ {
		mapper.MapOffset[typeOf.Field(i).Name] = unsafe.Add(start, typeOf.Field(i).Offset)
	}

	// cmhp_print.Print(mapper.MapOffset)

	return &mapper
}

func (v *ValueMapper[T]) Map(bytes []byte, fieldList []string, isSet bool) int {
	offset := 0

	for i := 0; i < len(fieldList); i++ {
		searchField := fieldList[i]

		switch bytes[0] {
		case core.TypeStruct:
			// Type
			offset += 1

			// Size
			offset += 2
			fieldSize := int(binary.LittleEndian.Uint16(bytes[1:]))

			// Amount
			amount := int(bytes[offset])
			offset += 1

			for j := 0; j < amount; j++ {
				// Get field name
				fieldLength := int(bytes[offset])
				offset += 1
				fieldName := string(bytes[offset : offset+fieldLength])
				offset += len(fieldName)

				// Go next
				offset += v.Map(bytes[offset:], []string{fieldName}, fieldName == searchField)
			}

			return fieldSize
		case core.TypeString:
			// Field size
			fieldSize := int(binary.LittleEndian.Uint16(bytes[1:]))

			if isSet {
				ptr := v.MapOffset[searchField]
				*(*string)(ptr) = string(bytes[3 : 3+fieldSize])
			}

			return 1 + 2 + fieldSize
		case core.TypeTime:
			// Field size
			fieldSize := int(bytes[1])

			return 1 + 1 + fieldSize
		case core.Type32:
			return 1 + 4
		case core.TypeSlice:
			// Field size
			fieldSize := int(binary.LittleEndian.Uint16(bytes[1:]))

			return 1 + 2 + 2 + fieldSize
		default:
			// log.Fatalf("unknown type %v", bytes[0])
		}
	}

	return 0
}
