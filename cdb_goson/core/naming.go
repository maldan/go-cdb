package core

import (
	"reflect"
	"strings"
)

/*func NameToId(v any) NameToId {
	out := NameToId{}
	nameToId(v, &out)
	return out
}

func NameList(v any) []string {
	out := NameToId{}
	nameToId(v, &out)
	out2 := make([]string, 0)
	for name, _ := range out {
		out2 = append(out2, name)
	}
	return out2
}

func IdToName(nameToId NameToId) IdToName {
	out := core.IdToName{}

	for k, v := range nameToId {
		out[v] = k
	}

	return out
}*/

type NameToId map[string]uint8
type IdToName map[uint8]string

func (n *NameToId) Add(name string) {
	(*n)[name] = uint8(len(*n)) + 1
}

func GetNameList(v any) []string {
	m := map[string]any{}
	l := make([]string, 0)

	// Get all fields names
	getFieldList(v, &l)
	for _, name := range l {
		m[name] = true
	}

	// Get unique list
	out := make([]string, 0)
	for name, _ := range m {
		out = append(out, name)
	}
	return out
}

func getFieldList(v any, out *[]string) {
	typeOf := reflect.TypeOf(v)

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)

		// Skip private
		if string(field.Name[0]) == strings.ToLower(string(field.Name[0])) {
			continue
		}

		*out = append(*out, field.Name)

		if field.Type.Kind() == reflect.Struct {
			getFieldList(reflect.New(field.Type).Elem().Interface(), out)
		}
		if field.Type.Kind() == reflect.Slice {
			getFieldList(reflect.New(field.Type.Elem()).Elem().Interface(), out)
		}
	}
}

/*func FromIdToName(nameToId NameToId) IdToName {
	out := core.IdToName{}

	for k, v := range nameToId {
		out[v] = k
	}

	return out
}*/

/*func nameToId(v any, out *NameToId) {
	typeOf := reflect.TypeOf(v)

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)

		// Skip private
		if string(field.Name[0]) == strings.ToLower(string(field.Name[0])) {
			continue
		}

		_, ok := (*out)[field.Name]
		if !ok {
			(*out)[field.Name] = uint8(len(*out)) + 1
		}

		if field.Type.Kind() == reflect.Struct {
			nameToId(reflect.New(field.Type).Elem().Interface(), out)
		}
		if field.Type.Kind() == reflect.Slice {
			nameToId(reflect.New(field.Type.Elem()).Elem().Interface(), out)
		}
	}
}
*/
