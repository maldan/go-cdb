package goson

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-cdb/cdb_goson/core"
	"strings"

	"reflect"
	"time"
)

type IR struct {
	Type    int
	Name    []byte
	Content []byte
	List    []*IR
}

func (r *IR) Len() int {
	outSize := 0

	// Name
	if len(r.Name) > 0 {
		outSize += 1
		outSize += len(r.Name)
	}

	// Type
	outSize += 1

	switch r.Type {
	case core.TypeStruct:
		// Size
		outSize += 2

		// Amount of elements
		outSize += 1

		for i := 0; i < len(r.List); i++ {
			outSize += r.List[i].Len()
		}
		break
	case core.TypeSlice:
		// Size
		outSize += 2

		// Amount of elements
		outSize += 2

		for i := 0; i < len(r.List); i++ {
			outSize += r.List[i].Len()
		}
		break
	case core.TypeString:
		outSize += 2
		outSize += len(r.Content)
		break
	case core.TypeTime:
		outSize += 1
		outSize += len(r.Content)
		break
	case core.Type32:
		outSize += 4
		break
	default:
		panic("unknown type " + fmt.Sprintf("%v", r.Type))
	}

	return outSize
}

func (r *IR) Build() []byte {
	s := make([]byte, 0, r.Len())

	// Name
	if len(r.Name) > 0 {
		s = append(s, uint8(len(r.Name)))
		s = append(s, r.Name...)
	}

	// Type
	s = append(s, uint8(r.Type))

	switch r.Type {
	case core.TypeStruct:
		// Len of struct
		l := r.Len()
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// Amount of elements
		l = len(r.List)
		s = append(s, uint8(l))

		// Elements
		for i := 0; i < len(r.List); i++ {
			s = append(s, r.List[i].Build()...)
		}
		break
	case core.TypeSlice:
		// Len of struct
		l := r.Len()
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// Amount of elements
		l = len(r.List)
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// Elements
		for i := 0; i < len(r.List); i++ {
			s = append(s, r.List[i].Build()...)
		}
		break
	case core.TypeString:
		// Content length
		l := len(r.Content)
		s = append(s, uint8(l))
		s = append(s, uint8(l>>8))

		// Content
		s = append(s, r.Content...)
		break
	case core.TypeTime:
		// Content length
		l := len(r.Content)
		s = append(s, uint8(l))

		// Content
		s = append(s, r.Content...)
		break
	case core.Type32:
		// Content
		s = append(s, r.Content...)
		break
	default:
		panic("unknown type " + fmt.Sprintf("%v", r.Type))
	}

	return s
}

func BuildIR(ir *IR, v any) {
	valueOf := reflect.ValueOf(v)
	typeOf := reflect.TypeOf(v)

	if typeOf.Kind() == reflect.Struct {
		if typeOf.Name() == "Time" {
			ir.Type = core.TypeTime
			ir.Content = []byte(valueOf.Interface().(time.Time).Format("2006-01-02T15:04:05.999-07:00"))
		} else {
			ir.Type = core.TypeStruct
			for i := 0; i < typeOf.NumField(); i++ {
				tr := IR{
					Name: []byte(typeOf.Field(i).Name),
				}
				ir.List = append(ir.List, &tr)
				BuildIR(&tr, valueOf.Field(i).Interface())
			}
		}
	}

	if typeOf.Kind() == reflect.Slice {
		ir.Type = core.TypeSlice

		for i := 0; i < valueOf.Len(); i++ {
			tr := IR{}
			ir.List = append(ir.List, &tr)
			BuildIR(&tr, valueOf.Index(i).Interface())
		}
	}

	if typeOf.Kind() == reflect.String {
		ir.Type = core.TypeString
		ir.Content = []byte(valueOf.String())
	}

	if typeOf.Kind() == reflect.Int {
		b := []byte{0, 0, 0, 0}
		binary.LittleEndian.PutUint32(b, uint32(valueOf.Int()))
		ir.Content = b
		ir.Type = core.Type32
	}
}

func NameToId(v any) map[string]int {
	out := map[string]int{}
	nameToId(v, &out)
	return out
}

func nameToId(v any, out *map[string]int) {
	typeOf := reflect.TypeOf(v)

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)

		// Skip private
		if string(field.Name[0]) == strings.ToLower(string(field.Name[0])) {
			continue
		}

		_, ok := (*out)[field.Name]
		if !ok {
			(*out)[field.Name] = len(*out)
		}

		if field.Type.Kind() == reflect.Struct {
			nameToId(reflect.New(field.Type).Elem().Interface(), out)
		}
		if field.Type.Kind() == reflect.Slice {
			nameToId(reflect.New(field.Type.Elem()).Elem().Interface(), out)
		}
	}
}
