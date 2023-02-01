package dson

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_file"
)

/*func CreateDynLength(length int) []byte {
	b := make([]byte, 0, 4)
	if length <= 0b0111_1111 {
		return []byte{uint8(length)}
	}
	if length <= 0b0111_1111_1111_1111 {
		z := (length) + (length >> 8)
		b1 := uint8(z) | 0b1000_0000
		b2 := uint8(z >> 8)
		return []byte{b1, b2}
	}
	return b
}

func ReadDynLength(b []byte) (int, int) {
	if b[0]&0b1000_0000 == 0 {
		return int(b[0]), 1
	}
	if b[0]&0b1000_0000 != 0 {
		b1 := int(b[0]) & 0b0111_1111
		b2 := int(b[1]) << 8

		return ((b1 | b2) >> 1) & 0b1111_1111_0111_1111, 2
	}
	return 0, 0
}*/

func Pack[T any](v T) {
	/*s := make([]byte, 0)
	S(&s, 0, v)
	cmhp_print.PrintBytesColored(s, 32, []cmhp_print.ColorRange{
		{0, 1, cmhp_print.BgRed},
	})
	fmt.Printf("%v\n", string(s))

	cmhp_file.Write("aa", s)
	cmhp_file.Write("bb", v)*/
	// x, _ := json.Marshal(v)
	// Traverser(x)

	ir := IR{}
	BuildIR(&ir, v)
	// cmhp_print.Print(ir)

	cmhp_file.Write("aa", ir.Build())

	Unpack(ir.Build())
}

func Unpack(bytes []byte) int {
	offset := 0
	tp := bytes[offset]
	offset += 1

	if tp == TypeStruct {
		// Size
		offset += 2

		amount := int(bytes[offset])
		offset += 1

		for i := 0; i < amount; i++ {
			fieldLen := int(bytes[offset])
			offset += 1
			fieldName := bytes[offset : offset+fieldLen]
			fmt.Printf("F: %v\n", string(fieldName))
			offset += fieldLen

			offset += Unpack(bytes[offset:])
		}
	}

	if tp == TypeSlice {
		// Size
		offset += 2

		amount := int(bytes[offset])
		offset += 2

		for i := 0; i < amount; i++ {
			offset += Unpack(bytes[offset:])
		}
	}

	if tp == TypeString {
		size := int(binary.LittleEndian.Uint16(bytes[offset:]))
		offset += 2
		blob := bytes[offset : offset+size]
		fmt.Printf("V: %v\n", string(blob))
		offset += size
	}

	if tp == TypeTime {
		size := int(bytes[offset])
		offset += 1
		blob := bytes[offset : offset+size]
		fmt.Printf("V: %v\n", string(blob))
		offset += size
	}

	if tp == TypeI32 {
		v := int(binary.LittleEndian.Uint32(bytes[offset:]))
		fmt.Printf("V: %v\n", v)
		offset += 4
	}

	return offset
}
