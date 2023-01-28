package main

import (
	"encoding/binary"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"github.com/maldan/go-cmhp/cmhp_byte"
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_print"
	"os"
	"reflect"
	"strconv"
	"time"
)

type Test struct {
	FirstName string `id:"0"`
	LastName  string `id:"1"`
	Phone     string `id:"2"`
	Lox       int    `id:"3"`
}

func Lg(a int, b int) bool {
	return a > b
}

func StrCmp(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Update(m []byte, index int, field string, v string) {
	offset := int(indexMap[index])
	typeOf := reflect.TypeOf(Test{})

	fieldOffsetIndex := 0
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Name == field {
			fieldOffsetIndex = i * 8
			break
		}
	}

	_, _, startData := readHeader(m, offset, fieldOffsetIndex)
	newLen := len(v)
	cmhp_byte.From32ToBuffer(&newLen, &m, offset+_hSize+_hTotal+fieldOffsetIndex+4)
	copy(m[startData:], v)
}

func readHeader(m []byte, offset int, fieldOffsetIndex int) (int, int, int) {
	size := int(binary.LittleEndian.Uint32(m[offset:]))

	localOffset := offset + 5 + fieldOffsetIndex

	fieldOffset := int(binary.LittleEndian.Uint32(m[localOffset:]))
	fieldLen := int(binary.LittleEndian.Uint32(m[localOffset+4:]))

	return size, fieldLen, offset + fieldOffset
}

func Find(m []byte, field string, v string) {
	offset := 0

	strAsBytes := []byte(v)
	comparator := StrCmp

	// Find field index
	typeOf := reflect.TypeOf(Test{})

	fieldOffsetIndex := 0
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Name == field {
			fieldOffsetIndex = i * 8
			break
		}
	}

	for {
		size, fieldLen, startData := readHeader(m, offset, fieldOffsetIndex)
		fieldData := m[startData : startData+fieldLen]

		isFound := comparator(strAsBytes, fieldData)

		if isFound {
			fmt.Printf("%v %v\n", "GAS!", string(fieldData))
			break
		}

		offset += size
		if offset >= len(m) {
			break
		}
	}
}

func generate() {
	amt := 1_000_000
	buff := make([]byte, 0)
	// offset := 0
	for i := 0; i < amt; i++ {
		t := Test{
			FirstName: fmt.Sprintf("%08d", i),
			LastName:  "xox" + strconv.Itoa(i),
			Phone:     cmhp_crypto.UID(12),
		}

		// B
		bb := Pack(t)
		buff = append(buff, bb...)

		/*size := uint16(2 + 4 + len(t.FirstName) + 1 + len(t.LastName) + 1 + len(t.Phone) + 1)
		buff = append(buff, make([]byte, size)...)

		offset += cmhp_byte.From16ToBuffer(&size, &buff, offset) // size
		offset += cmhp_byte.From32ToBuffer(&i, &buff, offset)    // len
		// Data
		offset += cmhp_byte.FromString8ToBuffer(&t.FirstName, &buff, offset)
		offset += cmhp_byte.FromString8ToBuffer(&t.LastName, &buff, offset)
		offset += cmhp_byte.FromString8ToBuffer(&t.Phone, &buff, offset)*/
	}

	/*ptr := "sasageo 00000000"
	amt := 100_000_000
	buff := make([]byte, amt*len(ptr))
	for i := 0; i < amt; i++ {
		prt2 := fmt.Sprintf("sasageo %v00000000", i)
		copy(buff[i*len(ptr):], prt2[0:16])
	}*/

	// fmt.Printf("%v\n", string(buff))
	cmhp_file.Write("file", buff)
}

func rec() {
	recover()
}

func mainX() {
	// Find("FirstName", "Roman")
	bb := Pack(Test{FirstName: "ABC", LastName: "BBB", Phone: "CCC"})
	cmhp_print.PrintDebugBytes(bb, 4, 1, 4, 4, 4, 5, 5, 5)
	cmhp_file.Write("test", bb)
}

func main() {
	// mainX()
	// generate()
	// return
	/* generate()
	return */
	// generate()
	// return
	// defer rec()

	f, _ := os.OpenFile("./file", os.O_RDWR, 0777)
	defer f.Close()

	mem, _ := mmap.Map(f, mmap.RDWR, 0)
	defer mem.Unmap()

	calculateIndex(mem)

	t := time.Now()
	/* offset := 0
	cycles := 0
	for i := 0; i < len(mem)-16; i++ {
		if string(mem[offset:offset+16]) == "sasageo 99999999" {
			fmt.Printf("%v %v\n", "GAS!", string(mem[offset:offset+16]))
			break
		}
		offset += 16
		cycles += 1
	} */
	//Find(mem, "Phone", "bzu9AIR6KcH4")
	Find(mem, "FirstName", "AC3")

	// Update(mem, 500_000, "FirstName", "AC3")
	// mem.Flush()
	fmt.Printf("Time: %v\n", time.Since(t))

	for {
		time.Sleep(time.Millisecond * 500)
	}
}
