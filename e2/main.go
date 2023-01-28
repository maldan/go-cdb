package main

import (
	"encoding/binary"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_print"
	"os"
	"reflect"
	"strconv"
	"time"
)

type Test struct {
	FirstName string
	LastName  string
	Phone     string
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

func Find(m []byte, field string, v string) {
	offset := 0

	strAsBytes := []byte(v)
	comparator := StrCmp

	// Find field index
	typeOf := reflect.TypeOf(Test{})

	fieldOffsetIndex := 0
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Name == field {
			fieldOffsetIndex = i * 4
			break
		}
	}

	for {
		size := binary.LittleEndian.Uint32(m[offset:])

		localOffset := offset + 5 + fieldOffsetIndex

		fieldOffset := binary.LittleEndian.Uint32(m[localOffset:])
		localOffset = offset + int(fieldOffset)

		fieldLen := binary.LittleEndian.Uint32(m[localOffset:])
		fieldData := m[localOffset+4 : localOffset+4+int(fieldLen)]

		isFound := comparator(strAsBytes, fieldData)

		if isFound {
			fmt.Printf("%v %v\n", "GAS!", string(fieldData))
			break
		}

		offset += int(size)
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
	//mainX()
	//generate()
	//return
	/*generate()
	return*/
	// generate()
	// return
	// defer rec()

	f, _ := os.OpenFile("./file", os.O_RDWR, 0644)
	defer f.Close()

	mem, _ := mmap.Map(f, mmap.RDWR, 0)
	defer mem.Unmap()
	fmt.Println(len(mem))

	t := time.Now()
	/*offset := 0
	cycles := 0
	for i := 0; i < len(mem)-16; i++ {
		if string(mem[offset:offset+16]) == "sasageo 99999999" {
			fmt.Printf("%v %v\n", "GAS!", string(mem[offset:offset+16]))
			break
		}
		offset += 16
		cycles += 1
	}*/
	Find(mem, "Phone", "zTOO2Ot4okhs")
	fmt.Printf("Time: %v\n", time.Since(t))

	for {
		time.Sleep(time.Millisecond * 500)
	}
}
