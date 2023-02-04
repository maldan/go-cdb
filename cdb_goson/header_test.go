package cdb_goson_test

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"github.com/maldan/go-cdb/cdb_proto/pack"
	"reflect"
	"testing"
)

func TestHeader(t *testing.T) {
	typeOf := reflect.TypeOf(Test{})
	bytes := pack.EmptyHeader[Test]()
	offset := 0

	if string(bytes[0:8]) != "PROT1234" {
		t.Fatalf("doesn't have file id")
	}
	offset += 8

	if bytes[offset] == 0 {
		t.Fatalf("version can't be null")
	}
	offset += 1

	// Check struct info
	offset += 8
	offset += 8
	if bytes[offset] != uint8(typeOf.NumField()) {
		t.Fatalf("amout of field mismatched")
	}

	totalFields := int(bytes[offset])
	offset += 1
	for i := 0; i < totalFields; i++ {
		id := bytes[offset]
		offset += 1

		fieldType := bytes[offset]
		offset += 1

		maxCapacity := binary.LittleEndian.Uint32(bytes[offset:])
		offset += 4

		nameLength := int(bytes[offset])
		offset += 1

		fieldName := string(bytes[offset : offset+nameLength])
		offset += nameLength

		// Check field info
		f, ok := typeOf.FieldByName(fieldName)
		if !ok {
			panic("field not found")
		}
		if f.Tag.Get("id") != fmt.Sprintf("%v", id) {
			t.Fatalf("incorrect id")
		}
		if maxCapacity > 0 {
			if f.Tag.Get("len") != fmt.Sprintf("%v", maxCapacity) {
				t.Fatalf("incorrect len")
			}
		}

		// Check types
		if f.Type.Kind() == reflect.String {
			if fieldType != core.TString {
				t.Fatalf("incorrect type")
			}
		}
	}
}
