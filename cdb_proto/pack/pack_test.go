package pack_test

import (
	"encoding/binary"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"github.com/maldan/go-cdb/cdb_proto/pack"
	"reflect"
	"testing"
)

type Test struct {
	FirstName string `json:"firstName" id:"0"`
	LastName  string `json:"lastName" id:"1"`
	Phone     string `json:"phone" id:"2"`
}

func TestPack(t *testing.T) {
	v := Test{FirstName: "Roman", LastName: "Baran", Phone: "Oman"}
	typeOf := reflect.TypeOf(Test{})
	valueOf := reflect.ValueOf(v)

	bytes := pack.Pack(v)

	// Check header
	if bytes[0] != 0x12 || bytes[1] != 0x34 {
		t.Fatalf("%v", "Fuck")
	}

	// Check end
	if bytes[len(bytes)-2] != 0x56 || bytes[len(bytes)-1] != 0x78 {
		t.Fatalf("%v", "Fuck")
	}

	// Check total size
	totalSize := int(binary.LittleEndian.Uint32(bytes[core.RecordStart:]))
	if totalSize != len(bytes) {
		t.Fatalf("%v", "Fuck")
	}

	// Check correctness
	if bytes[totalSize-2] != 0x56 || bytes[totalSize-1] != 0x78 {
		t.Fatalf("%v", "Fuck")
	}

	// Check total fields
	if uint8(typeOf.NumField()) != bytes[core.RecordStart+core.RecordSize] {
		t.Fatalf("%v", "Fuck")
	}

	// Go to offset table
	offsetTable := core.RecordStart + core.RecordSize + core.RecordFlags

	// Check fields len
	for i := 0; i < typeOf.NumField(); i++ {
		fieldLen := int(binary.LittleEndian.Uint32(bytes[offsetTable+i*core.RecordLenOff+4:]))
		if valueOf.Field(i).Len() != fieldLen {
			t.Fatalf("%v", "Fuck")
		}
	}

	// Check values
	for i := 0; i < typeOf.NumField(); i++ {
		fieldOffset := int(binary.LittleEndian.Uint32(bytes[offsetTable+i*core.RecordLenOff:]))
		fieldLen := int(binary.LittleEndian.Uint32(bytes[offsetTable+i*core.RecordLenOff+4:]))
		fieldData := bytes[fieldOffset : fieldOffset+fieldLen]

		if valueOf.Field(i).String() != string(fieldData) {
			t.Fatalf("%v", "Fuck")
		}
	}
}
