package cdb_lite

import (
	"errors"
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_byte"
	"io"
	"os"
	"sync"
)

type Storage[T IEngineComparable] struct {
	Buffer []StorageOperation[T]

	bufferMu  sync.Mutex
	dataTable *os.File
}

const OpAdd = 1
const OpDelete = 2
const OpUpdate = 3

const _hStart = 2  // 0x1234
const _hStatus = 1 // is removed
const _hLen = 3    // len of package, max 16 mb
const _hEnd = 2    // 0x5678

type StorageOperation[T IEngineComparable] struct {
	Type     uint8
	Position uint32
	Data     T
}

func (r *Storage[T]) AddToBuffer(o StorageOperation[T]) {
	r.bufferMu.Lock()
	defer r.bufferMu.Unlock()
	r.Buffer = append(r.Buffer, o)
}

func PackRecord[T IEngineComparable](v *T) []byte {
	// Pack struct
	bytes := cmhp_byte.Pack[T](v)
	bodyLength := uint32(len(bytes))

	// Record data
	recordData := make([]byte, _hStart+_hStatus+_hLen+len(bytes)+_hEnd)

	// Start of record
	recordData[0] = 0x12
	recordData[1] = 0x34

	// Status
	recordData[2] = 0 // non removed

	// End of record
	recordData[_hStart+_hStatus+_hLen+int(bodyLength)] = 0x56
	recordData[_hStart+_hStatus+_hLen+int(bodyLength)+1] = 0x78

	// Copy record length
	cmhp_byte.From24ToBuffer(&bodyLength, &recordData, _hStart+_hStatus)

	// Copy body
	copy(recordData[_hStart+_hStatus+_hLen:], bytes)

	return recordData
}

func UnpackRecord[T IEngineComparable](bytes *[]byte) (T, int64, error) {
	// Impossible to parse, not enough data
	if len(*bytes) < _hStart+_hStatus+_hLen {
		return *new(T), 0, errors.New("not full")
	}

	// Check header
	if !((*bytes)[0] == 0x12 && (*bytes)[1] == 0x34) {
		return *new(T), 1, errors.New("non record")
	}

	// Read length
	length := cmhp_byte.Read24FromBuffer(bytes, _hStart+_hStatus)

	// Check bounds
	if len(*bytes) < _hStart+_hStatus+_hLen+int(length)+_hEnd {
		return *new(T), 0, errors.New("not full")
	}

	// Check end
	if !((*bytes)[_hStart+_hStatus+_hLen+int(length)] == 0x56 && (*bytes)[_hStart+_hStatus+_hLen+int(length)+1] == 0x78) {
		return *new(T), 1, errors.New("non end")
	}

	// Read all bytes and unpack
	r := (*bytes)[_hStart+_hStatus+_hLen:]
	out := cmhp_byte.Unpack[T](&r)

	return out, int64(_hStart + _hStatus + _hLen + length + _hEnd), nil
}

func (r *Storage[T]) Flush() {
	r.bufferMu.Lock()
	defer r.bufferMu.Unlock()

	h, err := os.OpenFile("sas.body", os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	// Work with buffer
	for i := 0; i < len(r.Buffer); i++ {
		if r.Buffer[i].Type == OpAdd || r.Buffer[i].Type == OpUpdate {
			record := PackRecord(&r.Buffer[i].Data)

			// Some other error
			n, err := h.Write(record)
			if err != nil {
				panic(err)
			}

			// Fuck!
			if n < len(record) {
				panic(errors.New(fmt.Sprintf("write less bytes than excepted. Need [%v], wrote [%v]", len(record), n)))
			}
		}
	}

	h.Close()

	// Clear buffer
	r.Buffer = make([]StorageOperation[T], 0)
}

func (r *Storage[T]) Load(populate func(*T)) {
	h, err := os.OpenFile("sas.body", os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}

	// Buff
	offset := int64(0)
	packageOffset := int64(0)
	fileInfo, _ := h.Stat()
	increasePage := 0

	for {
		buffer := make([]byte, 65536+increasePage)
		_, err := h.ReadAt(buffer, offset)
		packageOffset = 0

		// Read packaged
		for {
			page := buffer[packageOffset:]
			p, n, e2 := UnpackRecord[T](&page)

			offset += n
			packageOffset += n
			if e2 != nil {
				if e2.Error() == "not full" {
					increasePage += 1024
				}
				break
			}
			populate(&p)
		}

		if offset >= fileInfo.Size() {
			break
		}

		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}
	}
}
