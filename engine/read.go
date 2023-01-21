package engine

import "fmt"

func (d *DataEngine[T]) FindBy(where func(*T) bool) SearchResult[T] {
	d.dataMu.RLock()
	defer d.dataMu.RUnlock()

	out := make([]Record, 0)
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
		Count:    uint32(len(out)),
		IsFound:  len(out) > 0,
		List:     out,
	}
}

func (d *DataEngine[T]) FindById(id uint32) SearchResult[T] {
	d.dataMu.RLock()
	defer d.dataMu.RUnlock()

	var fIndex = fmt.Sprintf("__id:%v", id)
	var out = make([]Record, 0)
	for _, position := range d.indexList[fIndex] {
		if d.recordList[position].isDeleted {
			continue
		}
		out = append(out, d.recordList[position])
	}
	return SearchResult[T]{
		dataBase: d,
		Count:    uint32(len(out)),
		IsFound:  len(out) > 0,
		List:     out,
	}
}

func (d *DataEngine[T]) FindByIndex(indexName string, indexValue any) SearchResult[T] {
	d.dataMu.RLock()
	defer d.dataMu.RUnlock()

	var fIndex = fmt.Sprintf("%v:%v", indexName, indexValue)
	var out = make([]Record, 0)

	for _, position := range d.indexList[fIndex] {
		if d.recordList[position].isDeleted {
			continue
		}
		out = append(out, d.recordList[position])
	}
	return SearchResult[T]{
		dataBase: d,
		Count:    uint32(len(out)),
		IsFound:  len(out) > 0,
		List:     out,
	}
}

func (d *DataEngine[T]) Length() uint32 {
	return d.length
}
