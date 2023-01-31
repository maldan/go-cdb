package cdb_proto

import (
	"encoding/binary"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"reflect"
	"unsafe"
)

type ValueMapper[T any] struct {
	Container T
	IdTable   []int
	TypeTable []int
	OutOffset []unsafe.Pointer
}

func (v *ValueMapper[T]) Map(structInfo core.StructInfo, fieldList []string) {
	typeOf := reflect.TypeOf(v.Container)
	start := unsafe.Pointer(&v.Container)

	// Fill id table
	for i := 0; i < len(fieldList); i++ {
		v.IdTable = append(v.IdTable, structInfo.FieldNameToId[fieldList[i]]*8)
		v.TypeTable = append(v.TypeTable, structInfo.FieldType[v.IdTable[i]/8])
		f, ok := typeOf.FieldByName(fieldList[i])
		if !ok {
			panic("Field " + fieldList[i] + " not found")
		}
		v.OutOffset = append(v.OutOffset, unsafe.Add(start, f.Offset))
	}
}

func (v *ValueMapper[T]) Map2(structInfo core.StructInfo, fieldList []string) {
	typeOf := reflect.TypeOf(v.Container)
	start := unsafe.Pointer(&v.Container)

	// Fill id table
	for i := 0; i < len(fieldList); i++ {
		v.IdTable = append(v.IdTable, structInfo.FieldNameToId[fieldList[i]])
		v.TypeTable = append(v.TypeTable, structInfo.FieldType[i])
		f, ok := typeOf.FieldByName(fieldList[i])
		if !ok {
			panic("Field " + fieldList[i] + " not found")
		}
		v.OutOffset = append(v.OutOffset, unsafe.Add(start, f.Offset))
	}
}

func (v *ValueMapper[T]) Fill(offset int, mem []byte, offTable []byte) {
	for i := 0; i < len(v.IdTable); i++ {
		vOff := int(binary.LittleEndian.Uint32(offTable[v.IdTable[i]:]))
		vLen := int(binary.LittleEndian.Uint32(offTable[v.IdTable[i]+4:]))

		if v.TypeTable[i] == core.TString {
			*(*[]byte)(v.OutOffset[i]) = mem[offset+vOff : offset+vOff+vLen]
		}
	}
}

func (v *ValueMapper[T]) Fill2(offset int, mem []byte) {
	for i := 0; i < len(v.IdTable); i++ {
		//id := mem[offset]
		//offset += 1

		if v.TypeTable[i] == core.TString {
			length := int(mem[offset])
			offset += 1

			*(*[]byte)(v.OutOffset[i]) = mem[offset : offset+length]
			offset += length
		}
	}
}
