package engine

import (
	"sync"
)

type IEngineComparable interface {
	comparable
	GetId() uint32
}

type DataEngine[T IEngineComparable] struct {
	Name          string
	AutoIncrement uint32
	IndexList     []string
	ShowLogs      bool

	length uint32

	rawDataList  []T
	recordList   []Record
	indexList    map[string][]uint32
	saveTaskList []uint32

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
