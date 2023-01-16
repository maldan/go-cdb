package cdb_flat

import (
	"fmt"
	"github.com/maldan/go-cdb/util"
	"reflect"
	"sync"
	"time"
)

type idAble interface {
	comparable
	GetId() uint
}

type DataBase[T idAble] struct {
	Name          string
	ChunkAmount   uint
	AutoIncrement uint
	IndexList     []string
	ShowLogs      bool

	length uint

	rawDataList  []T
	recordList   []Record[T]
	indexList    map[string][]uint
	saveTaskList []uint

	dataMu  sync.RWMutex
	storeMu sync.Mutex
}

type Record[T idAble] struct {
	chunkId   uint
	position  uint
	isDeleted bool
	isChanged bool
}

type SearchResult[T idAble] struct {
	dataBase *DataBase[T]
	Count    uint
	IsFound  bool
	List     []Record[T]
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

func (d *DataBase[T]) Init() *DataBase[T] {
	d.indexList = make(map[string][]uint, 0)

	// Read auto increment
	info, err := util.ReadText(d.Name + "/counter")
	if err == nil {
		d.AutoIncrement = util.StrToUInt(info)
	}

	// Init chunks
	for i := 0; i < int(d.ChunkAmount); i++ {
		t := time.Now()
		chunk, err := util.ReadJson[[]T](d.Name + "/chunk_" + fmt.Sprintf("%v", i) + ".json")
		if err == nil {
			for _, v := range chunk {
				d.Add(v)
			}
		}
		if d.ShowLogs {
			fmt.Printf("Load chunk [%v:%v] - %v records | %v\n", d.Name, i, len(chunk), time.Since(t))
		}
	}

	return d
}

func (d *DataBase[T]) EnableAutoSave() *DataBase[T] {
	go func() {
		for {
			d.dataMu.Lock()
			list := make([]uint, 0)
			if len(d.saveTaskList) > 0 {
				t := d.saveTaskList[0]
				d.saveTaskList = d.saveTaskList[1:]
				list = append(list, t)
			}
			d.dataMu.Unlock()

			for _, id := range list {
				d.SaveChunk(id)
			}

			time.Sleep(time.Second)
		}
	}()
	return d
}

func (d *DataBase[T]) SaveChunk(id uint) {
	d.dataMu.RLock()

	// List to out
	outList := make([]T, 0, d.Length()/d.ChunkAmount)

	// Check records
	for i := 0; i < len(d.recordList); i++ {
		if d.recordList[i].isDeleted {
			continue
		}
		if d.recordList[i].chunkId != id {
			continue
		}

		d.recordList[i].isChanged = false
		outList = append(outList, d.rawDataList[i])
	}

	d.dataMu.RUnlock()

	// Store
	d.storeMu.Lock()
	fmt.Printf("%v\n", "Saved")
	util.WriteJson(fmt.Sprintf("%v/chunk_%v.json", d.Name, id), &outList)
	d.storeMu.Unlock()
}

func (d *DataBase[T]) Add(v T) {
	d.dataMu.Lock()
	defer d.dataMu.Unlock()

	lastIndex := len(d.rawDataList)
	d.rawDataList = append(d.rawDataList, v)
	d.recordList = append(d.recordList, Record[T]{
		chunkId:  d.IdToChunk(v.GetId()),
		position: uint(lastIndex),
	})
	d.length += 1

	d.AddIndex(uint(lastIndex))
}

func (d *DataBase[T]) AddIndex(position uint) {
	for _, index := range d.IndexList {
		strIndex := ""

		f := reflect.ValueOf(&d.rawDataList[position]).Elem().FieldByName(index)
		mapIndex := reflect.ValueOf(f).Interface()
		strIndex = fmt.Sprintf("%s:%v", index, mapIndex)

		d.indexList[strIndex] = append(d.indexList[strIndex], position)
	}

	// Id index
	idIndex := fmt.Sprintf("__id:%v", d.rawDataList[position].GetId())
	d.indexList[idIndex] = append(d.indexList[idIndex], position)
}

func (d *DataBase[T]) FindBy(where func(*T) bool) SearchResult[T] {
	d.dataMu.RLock()
	defer d.dataMu.RUnlock()

	out := make([]Record[T], 0)
	for i := 0; i < len(d.rawDataList); i++ {
		if d.recordList[i].isDeleted {
			continue
		}
		if where(&d.rawDataList[i]) {
			out = append(out, d.recordList[i])
		}
	}

	return SearchResult[T]{
		dataBase: d,
		Count:    uint(len(out)),
		IsFound:  len(out) > 0,
		List:     out,
	}
}

func (d *DataBase[T]) FindById(id uint) SearchResult[T] {
	d.dataMu.RLock()
	defer d.dataMu.RUnlock()

	var fIndex = fmt.Sprintf("__id:%v", id)
	var out = make([]Record[T], 0)
	for _, position := range d.indexList[fIndex] {
		if d.recordList[position].isDeleted {
			continue
		}
		out = append(out, d.recordList[position])
	}
	return SearchResult[T]{
		dataBase: d,
		Count:    uint(len(out)),
		IsFound:  len(out) > 0,
		List:     out,
	}
}

func (d *DataBase[T]) FindByIndex(indexName string, indexValue any) SearchResult[T] {
	d.dataMu.RLock()
	defer d.dataMu.RUnlock()

	var fIndex = fmt.Sprintf("%v:%v", indexName, indexValue)
	var out = make([]Record[T], 0)

	for _, position := range d.indexList[fIndex] {
		if d.recordList[position].isDeleted {
			continue
		}
		out = append(out, d.recordList[position])
	}
	return SearchResult[T]{
		dataBase: d,
		Count:    uint(len(out)),
		IsFound:  len(out) > 0,
		List:     out,
	}
}

func (d *DataBase[T]) Update(r Record[T], v T) {
	d.dataMu.Lock()
	defer d.dataMu.Unlock()
	d.rawDataList[r.position] = v

	// Update record info
	d.recordList[r.position].isChanged = true

	// Add chunk to save
	d.saveTaskList = append(d.saveTaskList, d.IdToChunk(v.GetId()))
	d.saveTaskList = util.Unique(d.saveTaskList)
}

func (d *DataBase[T]) Delete(r Record[T]) {
	d.dataMu.Lock()
	defer d.dataMu.Unlock()

	d.length -= 1

	// Update record info
	d.recordList[r.position].isChanged = true
	d.recordList[r.position].isDeleted = true

	// Add chunk to save
	d.saveTaskList = append(d.saveTaskList, d.IdToChunk(d.rawDataList[r.position].GetId()))
	d.saveTaskList = util.Unique(d.saveTaskList)
}

func (d *DataBase[T]) Length() uint {
	return d.length
}

// GenerateId generate thread safe autoincrement id
func (d *DataBase[T]) GenerateId() uint {
	var out uint
	d.dataMu.Lock()
	d.AutoIncrement += 1
	out = d.AutoIncrement
	util.WriteText(d.Name+"/counter", fmt.Sprintf("%v", d.AutoIncrement))
	d.dataMu.Unlock()

	return out
}

func (d *DataBase[T]) IdToChunk(id uint) uint {
	return id % d.ChunkAmount
}
