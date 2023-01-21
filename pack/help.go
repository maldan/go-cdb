package pack

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

func __write8(b *[]byte, p int, v uint8) int {
	(*b)[p] = v
	return 1
}

func __write16(b *[]byte, p int, v uint32) int {
	__write8(b, p, byte(v&0xff))
	__write8(b, p+1, byte(v>>8))
	return 2
}

func __write32(b *[]byte, p int, v uint32) int {
	__write8(b, p, byte(v&0xff))
	__write8(b, p+1, byte((v>>8)&0xff))
	__write8(b, p+2, byte((v>>16)&0xff))
	__write8(b, p+3, byte(v>>24))
	return 4
}

func __read16(b *[]byte, p int) uint16 {
	return uint16(int((*b)[p]) + int((*b)[p+1])*256)
}

func __read32(b *[]byte, p int) uint32 {
	return uint32(int((*b)[p]) + int((*b)[p+1])*256 + int((*b)[p+2])*65536 + int((*b)[p+3])*16777216)
}

func __copy8(b *[]byte, p int, from uintptr, offset uintptr) int {
	(*b)[p] = *(*uint8)(unsafe.Pointer(from + offset))
	return 1
}

func __copy16(b *[]byte, p int, start uintptr, offset uintptr) int {
	__copy8(b, p, start, offset)
	__copy8(b, p+1, start, offset+1)
	return 2
}

func __copy32(b *[]byte, p int, start uintptr, offset uintptr) int {
	//bb := *(*[4]byte)(unsafe.Pointer(start + offset))
	//copy((*b)[p+1:], bb[:])
	__copy8(b, p, start, offset)
	__copy8(b, p+1, start, offset+1)
	__copy8(b, p+2, start, offset+2)
	__copy8(b, p+3, start, offset+3)
	return 4
}

func __copySmallString(b *[]byte, p int, from uintptr) int {
	bb := *(*[]byte)(unsafe.Pointer(from))
	ln := len(bb)
	__write8(b, p, uint8(ln))

	/*for i := 0; i < ln; i++ {
		__write8(b, p+1+i, bb[i])
	}*/
	copy((*b)[p+1:], bb)
	return 1 + ln
}

func __copyBigString(b *[]byte, p int, from uintptr) int {
	bb := *(*[]byte)(unsafe.Pointer(from))
	ln := len(bb)
	__write32(b, p, uint32(ln))

	/*for i := 0; i < ln; i++ {
		__write8(b, p+4+i, bb[i])
	}*/

	copy((*b)[p+4:], bb)

	return 4 + ln
}

func gas() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)
	var nativeEndian binary.ByteOrder

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		nativeEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		nativeEndian = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
	fmt.Printf("%v\n", nativeEndian)
}
