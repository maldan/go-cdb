package util

import (
	"unsafe"
)

func Copy8(to *[]byte, offset int, from unsafe.Pointer) {
	hh := *(*[1]byte)(from)
	copy((*to)[offset:], hh[:])
}

func Copy16(to *[]byte, offset int, from unsafe.Pointer) {
	hh := *(*[2]byte)(from)
	copy((*to)[offset:], hh[:])
}

func Copy32(to *[]byte, offset int, from unsafe.Pointer) {
	hh := *(*[4]byte)(from)
	copy((*to)[offset:], hh[:])
}

func Copy40(to *[]byte, offset int, from unsafe.Pointer) {
	hh := *(*[8]byte)(from)
	copy((*to)[offset:], hh[:])
}

func Read32(from *[]byte, to unsafe.Pointer, fromOffset int) {
	_t := (*[4]byte)(to)
	copy(_t[:], (*from)[fromOffset:])
}
