package cdb_lite

import (
	"reflect"
	"sync"
)

type IEngineComparable interface {
	GetId() uint32
}

type DataEngine[T IEngineComparable] struct {
	Name              string
	AutoIncrement     uint32
	IndexList         []string
	SearchFieldByList []string
	ShowLogs          bool

	length uint32

	rawDataListAsMap []map[string]any
	rawDataList      []T
	recordList       []Record
	indexList        map[string][]uint32
	saveTaskList     []uint32

	dataMu  sync.RWMutex
	storeMu sync.Mutex

	storage Storage[T]
}

type Record struct {
	position  uint32
	isDeleted bool

	//offset uint32
	//length uint32
}

type SearchResult[T IEngineComparable] struct {
	dataBase *DataEngine[T]
	Count    uint32
	IsFound  bool
	List     []Record
}

func (r *SearchResult[T]) UnpackList() []T {
	r.dataBase.dataMu.RLock()
	defer r.dataBase.dataMu.RUnlock()

	var out = make([]T, 0)
	for _, x := range r.List {
		if x.isDeleted {
			continue
		}
		out = append(out, r.dataBase.rawDataList[x.position])
	}
	return out
}

type (
	TpInfo struct {
		Name           [][]byte
		Offset         []uintptr
		Type           []uint8
		MaxBytesLength []int
		FieldAmount    uint8

		MapOffset map[string]uintptr
		MapType   map[string]uint8
	}
)

var typeCache = map[reflect.Type]*TpInfo{}
