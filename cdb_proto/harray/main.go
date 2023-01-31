package harray

import (
	"encoding/binary"
	"fmt"
)

/**
bucket struct
[0] - type. Bucket or data page
[0] - is full
[0 0 0 0] - block hash
[0 0 0 0 0 0 0 0] - next free block
[
	[0 0 0 0 0 0 0 0] - 0 cell. Contains offset to block
	[0 0 0 0 0 0 0 0] - offset to data
] * bucketSize
*/

/**
data block
[0] - type
[
	[0 0 0 0] - offset to data start
] * bucketSize
[
	[0 0 0 0] - length
	[0 0 0 0] - capacity
	[...] - data
] * bucketSize
*/

const OffsetToCellStart = 1 + 1 + 4 + 8

type HashArray struct {
	BucketSize int
	Memory     []byte
	Size       int
}

func New() *HashArray {
	return &HashArray{
		Memory:     make([]byte, 1024),
		BucketSize: 8,
		Size:       0,
	}
}

func (h *HashArray) CreateBucket() []byte {
	start := h.Size
	binary.LittleEndian.PutUint32(h.Memory[start+1+1:], 0x11223344)
	bucket := h.Memory[start : start+OffsetToCellStart+h.BucketSize*8]
	h.Size += len(bucket)
	return bucket
}

func (h *HashArray) CreateDataBlock() (int, []byte) {
	start := h.Size
	h.Memory[start] = 1
	h.Size += 256
	return start, h.Memory[start : start+h.Size]
}

func (h *HashArray) GetFreeBucket() []byte {
	if h.Memory[1] == 0 {
		return h.Memory[0 : OffsetToCellStart+h.BucketSize*8]
	}
	panic("gas")
	// blockOffset := binary.LittleEndian.Uint64(h.Memory[index*8:])
	// return h.Memory[0 : OffsetToCellStart+h.BucketSize*8]
}

func (h *HashArray) GetFreeDataBlock() []byte {
	// blockOffset := binary.LittleEndian.Uint64(h.Memory[index*8:])
	return h.Memory[0 : OffsetToCellStart+h.BucketSize*8]
}

func (h *HashArray) Add(key string, value int) {
	index := Hash(key) % h.BucketSize
	bb := h.CreateBucket()
	startBlock, dataBlock := h.CreateDataBlock()
	fmt.Printf("%v\n", index)
	binary.LittleEndian.PutUint64(bb[OffsetToCellStart+index*8:], uint64(startBlock))
	binary.LittleEndian.PutUint32(dataBlock[1+index*4:], 32)
}

func Hash(key string) int {
	start := int(key[0])
	x := len(key) << 32

	for i := 0; i < len(key); i++ {
		// b := key[i] << 4
		// fmt.Printf("%v\n", b)
		x += int(key[i]) << 32
	}
	return x ^ start
}
