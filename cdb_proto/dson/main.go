package dson

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_file"
	"reflect"
	"time"
	"unsafe"
)

func Pack[T any](v T) []byte {
	ir := IR{}
	BuildIR(&ir, v)
	cmhp_file.Write("aa", ir.Build())
	return ir.Build()
}

func Unpack(bytes []byte, ptr unsafe.Pointer, typeHint any) int {
	offset := 0
	tp := bytes[offset]
	offset += 1

	if tp == TypeStruct {
		// Size
		offset += 2

		// Amount
		amount := int(bytes[offset])
		offset += 1

		typeOf := reflect.TypeOf(typeHint)

		for i := 0; i < amount; i++ {
			fieldLen := int(bytes[offset])
			offset += 1
			fieldName := string(bytes[offset : offset+fieldLen])
			offset += fieldLen

			field, _ := typeOf.FieldByName(fieldName)

			if field.Type.Kind() == reflect.Slice {
				offset += Unpack(bytes[offset:], unsafe.Add(ptr, field.Offset), reflect.ValueOf(typeHint).FieldByName(fieldName).Interface())
			} else if field.Type.Kind() == reflect.Struct {
				// fmt.Printf("%v\n", reflect.ValueOf(typeHint).FieldByName(fieldName).Interface())
				offset += Unpack(bytes[offset:], unsafe.Add(ptr, field.Offset), reflect.ValueOf(typeHint).FieldByName(fieldName).Interface())
			} else {
				offset += Unpack(bytes[offset:], unsafe.Add(ptr, field.Offset), typeHint)
			}
		}
	}

	if tp == TypeSlice {
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
			offset += Unpack(bytes[offset:], unsafe.Pointer(elemSlice.Index(i).Addr().Pointer()), typeHint)
			arr[i] = elemSlice.Index(i).Interface()
		}

		g := elemSlice.Pointer()

		*(*reflect.SliceHeader)(ptr) = reflect.SliceHeader{
			Data: g,
			Len:  amount,
			Cap:  amount,
		}
	}

	if tp == TypeString {
		size := int(binary.LittleEndian.Uint16(bytes[offset:]))
		offset += 2
		blob := bytes[offset : offset+size]
		offset += size
		*(*string)(ptr) = string(blob)
	}

	if tp == TypeTime {
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

	if tp == TypeI32 {
		*(*int)(ptr) = int(binary.LittleEndian.Uint32(bytes[offset:]))
		offset += 4
	}

	return offset
}
