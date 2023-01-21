package engine

import (
	"errors"
	"fmt"
	"github.com/maldan/go-cdb/pack"
	"github.com/maldan/go-cdb/util"
	"io"
	"os"
	"sync"
	"unsafe"
)

type Storage[T IEngineComparable] struct {
	Buffer []StorageOperation[T]

	bufferMu  sync.Mutex
	dataTable *os.File
}

const OpAdd = 1
const OpDelete = 2
const OpUpdate = 3
const OpSetInfo = 4

const _hStart = 2
const _hLen = 4
const _hEnd = 2

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
	bytes := pack.Pack[T](v)
	bodyLength := uint32(len(bytes))

	// Record data
	recordData := make([]byte, _hStart+_hLen+len(bytes)+_hEnd)

	// Start of record
	recordData[0] = 0x12
	recordData[1] = 0x34

	// End of record
	recordData[_hStart+_hLen+int(bodyLength)] = 0x56
	recordData[_hStart+_hLen+int(bodyLength)+1] = 0x78

	// Copy record length
	util.Copy40(&recordData, _hStart, unsafe.Pointer(&bodyLength))

	// Copy body
	copy(recordData[_hStart+_hLen:], bytes)

	return recordData
}

func UnpackRecord[T IEngineComparable](bytes *[]byte) (T, int64, error) {
	// Impossible to parse
	if len(*bytes) < _hStart+_hLen {
		return *new(T), 0, errors.New("not full")
	}

	// Check header
	if !((*bytes)[0] == 0x12 && (*bytes)[1] == 0x34) {
		return *new(T), 1, errors.New("non record")
	}

	// Read length
	length := uint32(0)
	util.Read32(bytes, unsafe.Pointer(&length), 2)

	// Check bounds
	if len(*bytes) < _hStart+_hLen+int(length)+_hEnd {
		return *new(T), 0, errors.New("not full")
	}

	// Check end
	if !((*bytes)[2+4+length] == 0x56 && (*bytes)[2+4+length+1] == 0x78) {
		return *new(T), 1, errors.New("non end")
	}

	r := (*bytes)[_hStart+_hLen:]
	out := pack.Unpack[T](&r)

	return out, int64(_hStart + _hLen + length + _hEnd), nil
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

func (r *Storage[T]) Load(populate func(T)) {
	h, err := os.OpenFile("sas.body", os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}

	// Buff
	offset := int64(0)
	packageOffset := int64(0)
	fileInfo, _ := h.Stat()

	for {
		buffer := make([]byte, 65536)
		_, err := h.ReadAt(buffer, offset)
		packageOffset = 0

		// Read packaged
		for {
			page := buffer[packageOffset:]
			p, n, e2 := UnpackRecord[T](&page)

			offset += n
			packageOffset += n
			if e2 != nil {
				// fmt.Printf("%v\n", e2)
				break
			}
			populate(p)
			// fmt.Printf("%+v\n", p)
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

/*func (r *Storage[T]) LoadInfo() {
	h, err := os.OpenFile("sas.info", os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}
	h.Seek(0, 0)

	// Read to buff
	buff := make([]byte, 8)
	h.Read(buff)
	fmt.Printf("%v\n", buff)

	autoIncrement := uint32(255)
	length := uint32(255)

	util.Read32(&buff, unsafe.Pointer(&autoIncrement), 0)
	util.Read32(&buff, unsafe.Pointer(&length), 4)
}

func (r *Storage[T]) LoadHeader() {
	r.dataTable, _ = os.OpenFile("sas.bin", os.O_RDONLY, 0777)

	// r.dataTable.Seek(0, 0)
}*/

/*func (r *Storage[T]) WriteInfo() {
	h, _ := os.OpenFile("sas.info", os.O_WRONLY|os.O_CREATE, 0777)

	// Work with buffer
	for i := 0; i < len(r.Buffer); i++ {
		if r.Buffer[i].Type == OpSetInfo {
			header := make([]byte, 4+4)
			util.Copy32(&header, 0, unsafe.Pointer(&r.Buffer[i].AutoIncrement))
			util.Copy32(&header, 4, unsafe.Pointer(&r.Buffer[i].Length))
			h.Seek(0, 0)
			h.Write(header)
		}
	}

	h.Close()
}
*/
/*func (r *Storage[T]) WriteHeader() {
	// Add
	h, err := os.OpenFile("sas.header", os.O_APPEND|os.O_CREATE, 0777)

	if err != nil {
		panic(err)
	}
	for i := 0; i < len(r.Buffer); i++ {
		if r.Buffer[i].Type == OpAdd {
			record := make([]byte, 1+5)

			record[0] = 0
			// util.Copy40(&record, 1, unsafe.Pointer(&r.Buffer[i].RecordOffset))
			util.Copy40(&record, 1, unsafe.Pointer(&r.Buffer[i].RecordLength))

			// fmt.Printf("OFF: %v\n", record)
			h.Write(record)
		}
	}
	h.Close()

	// Change
	h, err = os.OpenFile("sas.header", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(r.Buffer); i++ {
		if r.Buffer[i].Type == OpDelete || r.Buffer[i].Type == OpUpdate {
			offset := r.Buffer[i].Position * (5 + 5 + 1)
			h.Seek(int64(offset), 0)
			h.Write([]byte{1}) // mark as deleted
		}
	}
	h.Close()

	// isDeleted = 1 byte
	// offset - 5 byte
	// length - 5 byte
}*/
