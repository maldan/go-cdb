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

/*func TestUnpack(t *testing.T) {
	v := Test{FirstName: "Roman", LastName: "Baran", Phone: "Oman"}
	bytes := pack.Pack(v)
	v2, _ := pack.Unpack(bytes)
	fmt.Printf("%v\n", v2)
}
*/

func TestPack(t *testing.T) {
	v := Test{FirstName: "Roman", LastName: "Baran", Phone: "Oman"}
	typeOf := reflect.TypeOf(Test{})
	valueOf := reflect.ValueOf(v)

	bytes := pack.Pack(v)

	// Check header
	if bytes[0] != core.RecordStartMark {
		t.Fatalf("%v", "Fuck")
	}

	// Check end
	if bytes[len(bytes)-1] != core.RecordEndMark {
		t.Fatalf("%v", "Fuck")
	}

	// Check total size
	totalSize := int(binary.LittleEndian.Uint32(bytes[core.RecordStart:]))
	if totalSize != len(bytes) {
		t.Fatalf("%v", "Fuck")
	}

	// Check correctness
	if bytes[totalSize-1] != core.RecordEndMark {
		t.Fatalf("%v", "Fuck")
	}

	// Check total fields
	if uint8(typeOf.NumField()) != bytes[core.RecordStart+core.RecordSize] {
		t.Fatalf("%v", "Fuck")
	}

	// Go to offset table
	offsetTable := core.RecordStart + core.RecordSize + core.RecordFlags

	// Read fields
	for i := 0; i < typeOf.NumField(); i++ {
		// Read ID
		// bytes[offsetTable]
		offsetTable += 1

		// Read str length
		fieldLen := int(binary.LittleEndian.Uint32(bytes[offsetTable:]))
		offsetTable += 4

		// Check data
		fieldData := bytes[offsetTable : offsetTable+fieldLen]
		if valueOf.Field(i).String() != string(fieldData) {
			t.Fatalf("%v", "Fuck")
		}
		offsetTable += fieldLen
	}
}
