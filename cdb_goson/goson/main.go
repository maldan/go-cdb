package goson

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-cdb/cdb_goson/core"
	"log"
	"reflect"
	"time"
	"unsafe"
)

func Marshal[T any](v T, nameToId core.NameToId) []byte {
	ir := IR{}
	BuildIR(&ir, v, nameToId)

	return ir.Build()
}

func Unmarshall[T any](bytes []byte, idToName core.IdToName) T {
	t := new(T)
	unpack(bytes, unsafe.Pointer(t), *t, idToName)
	return *t
}

func unpack(bytes []byte, ptr unsafe.Pointer, typeHint any, idToName core.IdToName) int {
	offset := 0

	// Read type
	tp := bytes[offset]
	offset += 1

	if tp == core.TypeStruct {
		// Size
		offset += 2

		// Amount
		amount := int(bytes[offset])
		offset += 1

		typeOf := reflect.TypeOf(typeHint)

		for i := 0; i < amount; i++ {
			// field length
			//fieldLen := int(bytes[offset])
			//offset += 1

			// field name
			fieldName, ok := idToName[bytes[offset]]
			if !ok {
				log.Fatalf("field for id %v not found", bytes[offset])
			}
			offset += 1
			//fieldName := string(bytes[offset : offset+fieldLen])
			//offset += fieldLen

			field, _ := typeOf.FieldByName(fieldName)

			if field.Type.Kind() == reflect.Slice {
				offset += unpack(
					bytes[offset:],
					unsafe.Add(ptr, field.Offset),
					reflect.ValueOf(typeHint).FieldByName(fieldName).Interface(),
					idToName,
				)
			} else if field.Type.Kind() == reflect.Struct {
				offset += unpack(
					bytes[offset:],
					unsafe.Add(ptr, field.Offset),
					reflect.ValueOf(typeHint).FieldByName(fieldName).Interface(),
					idToName,
				)
			} else {
				offset += unpack(bytes[offset:], unsafe.Add(ptr, field.Offset), typeHint, idToName)
			}
		}
	}

	if tp == core.TypeSlice {
		// Size
		offset += 2

		// Amount
		amount := int(bytes[offset])
		offset += 2

		typeOf := reflect.TypeOf(typeHint).Elem()
		typeHint = reflect.New(typeOf).Elem().Interface()

		elemSlice := reflect.MakeSlice(reflect.SliceOf(typeOf), amount, amount)
		arr := make([]any, amount, amount)

		for i := 0; i < amount; i++ {
			offset += unpack(
				bytes[offset:],
				unsafe.Pointer(elemSlice.Index(i).Addr().Pointer()),
				typeHint,
				idToName,
			)
			arr[i] = elemSlice.Index(i).Interface()
		}

		g := elemSlice.Pointer()

		*(*reflect.SliceHeader)(ptr) = reflect.SliceHeader{
			Data: g,
			Len:  amount,
			Cap:  amount,
		}
	}

	if tp == core.TypeString {
		size := int(binary.LittleEndian.Uint16(bytes[offset:]))
		offset += 2
		blob := bytes[offset : offset+size]
		offset += size
		*(*string)(ptr) = string(blob)
	}

	if tp == core.TypeTime {
		size := int(bytes[offset])
		offset += 1
		blob := bytes[offset : offset+size]

		x, err := time.Parse("2006-01-02T15:04:05.999-07:00", string(blob))
		*(*time.Time)(ptr) = x
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		offset += size
	}

	if tp == core.Type32 {
		*(*int)(ptr) = int(binary.LittleEndian.Uint32(bytes[offset:]))
		offset += 4
	}

	return offset
}
