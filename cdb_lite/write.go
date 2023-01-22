package cdb_lite

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_convert"
	"reflect"
)

func (d *DataEngine[T]) loadRecord(v *T) {
	d.dataMu.Lock()
	defer d.dataMu.Unlock()

	lastIndex := len(d.rawDataList)
	d.rawDataList = append(d.rawDataList, *v)
	d.recordList = append(d.recordList, Record{
		position: uint32(lastIndex),
	})
	d.length += 1

	// Add cache to query search
	if len(d.SearchFieldByList) > 0 {
		d.rawDataListAsMap = append(d.rawDataListAsMap, d.fastConvert(v))
	}

	d.AddIndex(uint32(lastIndex))
}

func (d *DataEngine[T]) Add(v *T) {
	d.loadRecord(v)

	// Add to buffer
	d.storage.AddToBuffer(StorageOperation[T]{
		Type: OpAdd,
		Data: *v,
	})
}

func (d *DataEngine[T]) AddIndex(position uint32) {
	for _, index := range d.IndexList {
		strIndex := ""

		f := reflect.ValueOf(&d.rawDataList[position]).Elem().FieldByName(index)
		mapIndex := reflect.ValueOf(f).Interface()
		strIndex = fmt.Sprintf("%s:%v", index, mapIndex)

		d.indexList[strIndex] = append(d.indexList[strIndex], position)
	}

	// Id index
	idIndex := "__id:" + cmhp_convert.IntToStr(int(d.rawDataList[position].GetId()))
	d.indexList[idIndex] = append(d.indexList[idIndex], position)
}

func (d *DataEngine[T]) Update(r Record, v *T) {
	d.dataMu.Lock()
	d.recordList[r.position].isDeleted = true
	d.storage.AddToBuffer(StorageOperation[T]{Type: OpUpdate, Position: r.position, Data: *v})
	d.dataMu.Unlock()

	d.Add(v)
}

func (d *DataEngine[T]) Delete(r Record) {
	d.dataMu.Lock()
	defer d.dataMu.Unlock()

	d.length -= 1

	// Update record info
	d.recordList[r.position].isDeleted = true

	d.storage.AddToBuffer(StorageOperation[T]{
		Type:     OpDelete,
		Position: r.position,
	})
}

// GenerateId generate thread safe autoincrement id
func (d *DataEngine[T]) GenerateId() uint32 {
	out := uint32(0)
	d.dataMu.Lock()
	d.AutoIncrement += 1
	out = d.AutoIncrement
	d.dataMu.Unlock()

	return out
}
