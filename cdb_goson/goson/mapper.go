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

	return &mapper
}

func strCmp(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func typeSize(bytes []byte) int {
	switch bytes[0] {
	case core.TypeString:
		// Field size
		fieldSize := int(binary.LittleEndian.Uint16(bytes[1:]))
		return 1 + 2 + fieldSize
	case core.Type32:
		return 4
	default:
		return 0
	}
}

func applyType[T any](v *ValueMapper[T], bytes []byte, fieldName string) {
	off := v.MapOffset[fieldName]
	switch bytes[0] {
	case core.TypeString:
		fieldSize := int(binary.LittleEndian.Uint16(bytes[1:]))
		// fmt.Printf("%v\n", string(bytes[2+1:2+1+fieldSize]))
		*(*string)(off) = string(bytes[2+1 : 2+1+fieldSize])
	default:
		break
	}
}

func handleStruct[T any](v *ValueMapper[T], bytes []byte, offset int, searchField string) int {
	// Type
	offset += 1

	// Size
	offset += 2
	size := int(binary.LittleEndian.Uint16(bytes[1:]))

	// Field amount
	amount := int(bytes[offset])
	offset += 1

	for i := 0; i < amount; i++ {
		// Field type
		// fieldType := int(bytes[offset])
		// offset += 1

		// Get field name
		fieldLength := int(bytes[offset])
		offset += 1
		fieldName := string(bytes[offset : offset+fieldLength])
		offset += len(fieldName)

		// Field matches
		if fieldName == searchField {
			applyType(v, bytes[offset:], fieldName)
			return size
		}

		// Go next
		fieldSize := typeSize(bytes[offset:])
		offset += fieldSize
	}

	return size
}

func (v *ValueMapper[T]) Map(bytes []byte, fieldList []string) {
	offset := 0

	for i := 0; i < len(fieldList); i++ {
		searchField := fieldList[i]

		if bytes[0] == core.TypeStruct {
			offset += handleStruct[T](v, bytes, offset, searchField)
		}
		/*switch bytes[0] {
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
		}*/
	}
}
