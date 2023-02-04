package goson

import (
	"encoding/binary"
	"github.com/maldan/go-cdb/cdb_goson/core"
	"log"
	"reflect"
	"unsafe"
)

type ValueMapper[T any] struct {
	Container T
	NameToId  core.NameToId
	MapOffset []unsafe.Pointer
}

func NewMapper[T any](nameToId core.NameToId) *ValueMapper[T] {
	mapper := ValueMapper[T]{
		Container: *new(T),
		NameToId:  nameToId,
		MapOffset: make([]unsafe.Pointer, 255),
	}

	typeOf := reflect.TypeOf(mapper.Container)
	start := unsafe.Pointer(&mapper.Container)

	for i := 0; i < typeOf.NumField(); i++ {
		fieldId, ok := nameToId[typeOf.Field(i).Name]
		if !ok {
			log.Fatalf("field id %v not found", typeOf.Field(i).Name)
		}
		mapper.MapOffset[fieldId] = unsafe.Add(start, typeOf.Field(i).Offset)
	}

	return &mapper
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

func applyType[T any](v *ValueMapper[T], bytes []byte, offset int, fieldName uint8) {
	off := v.MapOffset[fieldName]

	switch bytes[offset] {
	case core.TypeString:
		fieldSize := int(binary.LittleEndian.Uint16(bytes[offset+1:]))
		bts := *(*reflect.SliceHeader)(unsafe.Pointer(&bytes))

		hh := (*reflect.StringHeader)(off)
		hh.Data = bts.Data + uintptr(offset) + 2 + 1
		hh.Len = fieldSize
	default:
		break
	}
}

func handleStruct[T any](v *ValueMapper[T], bytes []byte, offset int, searchField uint8) int {
	// Type
	offset += 1

	// Size
	offset += 2
	size := int(binary.LittleEndian.Uint16(bytes[1:]))

	// Field amount
	amount := int(bytes[offset])
	offset += 1

	for i := 0; i < amount; i++ {
		// Read field id
		fieldId := bytes[offset]
		offset += 1

		// Field matches
		if fieldId == searchField {
			applyType(v, bytes, offset, fieldId)
			return size
		}

		// Go next
		fieldSize := typeSize(bytes[offset:])
		offset += fieldSize
	}

	return size
}

func (v *ValueMapper[T]) Map(bytes []byte, fieldList []uint8) {
	offset := 0

	for i := 0; i < len(fieldList); i++ {
		searchField := fieldList[i]

		if bytes[0] == core.TypeStruct {
			offset += handleStruct[T](v, bytes, offset, searchField)
		}
	}
}
