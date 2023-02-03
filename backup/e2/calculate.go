package main

import (
	"encoding/binary"
)

var indexMap = make([]uint32, 0)

func calculateIndex(m []byte) {
	offset := 0

	for {
		size := int(binary.LittleEndian.Uint32(m[offset:]))
		indexMap = append(indexMap, uint32(offset))
		offset += size
		if offset >= len(m) {
			break
		}
	}
	// cmhp_file.Write("offset", indexMap)
	// fmt.Printf("%v\n", len(indexMap))
}
