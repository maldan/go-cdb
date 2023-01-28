package cdb_lite

import (
	"github.com/maldan/go-cmhp/cmhp_byte"
	"os"
	"runtime"
	"unsafe"
)

func (d *DataEngine[T]) Init() *DataEngine[T] {
	d.indexList = make(map[string][]uint32, 0)

	// Ability
	h, err := os.OpenFile(d.Name+".blob", os.O_RDONLY, 0777)
	defer h.Close()
	if err == nil {
		fileInfo, _ := h.Stat()
		// fmt.Printf("XX: %v\n", fileInfo.Size()/int64(unsafe.Sizeof(*new(T))))
		probablyRecordsAmount := fileInfo.Size() / int64(unsafe.Sizeof(*new(T)))
		d.rawDataList = make([]T, 0, probablyRecordsAmount)
		d.recordList = make([]Record, 0, probablyRecordsAmount)
	}

	return d
}

func (d *DataEngine[T]) Flush() {
	d.storage.Flush()
}

func (d *DataEngine[T]) Load() {
	d.storage.Load(func(v T) {
		d.loadRecord(v)
	})
	runtime.GC()
}

func New[T IEngineComparable](name string, result []string) *DataEngine[T] {
	n := DataEngine[T]{
		Name:              name,
		SearchFieldByList: result,
		storage: Storage[T]{
			name: name,
		},
	}
	cmhp_byte.Pack[T](new(T)) // cache type + check
	return &n
}
