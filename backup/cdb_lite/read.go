package cdb_lite

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/maldan/go-cmhp/cmhp_slice"
	"runtime"
)

func (d *DataEngine[T]) FindBy(where func(*T) bool) SearchResult[T] {
	ch := make(chan SearchResult[T], 1)
	d.findBy(where, ch, 0, int(d.length))
	return <-ch
}

func (d *DataEngine[T]) FindByParallel(where func(*T) bool) SearchResult[T] {
	// Prepare out
	out := make([]Record, 0)

	// Prepare channels
	treads := runtime.NumCPU()
	chanRes := make(chan SearchResult[T], treads)
	sliceSize := int(d.length) / treads

	// Run search
	for i := 0; i < treads; i++ {
		go d.findBy(where, chanRes, i*sliceSize, sliceSize)
	}

	// Collect answers
	for i := 0; i < treads; i++ {
		rs := <-chanRes
		out = append(out, rs.List...)
	}

	return SearchResult[T]{
		dataBase: d,
		Count:    uint32(len(out)),
		IsFound:  len(out) > 0,
		List:     out,
	}
}

func (d *DataEngine[T]) findBy(where func(*T) bool, ch chan<- SearchResult[T], offset int, limit int) {
	d.dataMu.RLock()
	defer d.dataMu.RUnlock()

	// Make slices
	rawSlice := cmhp_slice.Paginate(d.rawDataList, offset, limit)
	recordSlice := cmhp_slice.Paginate(d.recordList, offset, limit)

	// Test
	out := make([]Record, 0)
	for i := 0; i < len(rawSlice); i++ {
		if recordSlice[i].isDeleted {
			continue
		}
		if where(&rawSlice[i]) {
			out = append(out, recordSlice[i])
		}
	}

	ch <- SearchResult[T]{
		dataBase: d,
		Count:    uint32(len(out)),
		IsFound:  len(out) > 0,
		List:     out,
	}
}

func (d *DataEngine[T]) FindByQueryParallel(query string) SearchResult[T] {
	// Prepare out
	out := make([]Record, 0)

	treads := runtime.NumCPU()
	fmt.Printf("%v\n", treads)
	chanRes := make(chan SearchResult[T], treads)
	ss := int(d.length) / treads
	for i := 0; i < treads; i++ {
		go d.FindByQuery(query, chanRes, i*ss, ss)
	}
	for i := 0; i < treads; i++ {
		rs := <-chanRes
		out = append(out, rs.List...)
	}

	return SearchResult[T]{
		dataBase: d,
		Count:    uint32(len(out)),
		IsFound:  len(out) > 0,
		List:     out,
	}
}

func (d *DataEngine[T]) FindByQuery(query string, ch chan<- SearchResult[T], offset int, limit int) SearchResult[T] {
	d.dataMu.RLock()
	defer d.dataMu.RUnlock()

	// Prepare out
	out := make([]Record, 0)

	// Prepare expression
	expression, err := govaluate.NewEvaluableExpression(query)
	if err != nil {
		panic(err)
	}

	rawSlice := cmhp_slice.Paginate(d.rawDataList, offset, limit)
	recordSlice := cmhp_slice.Paginate(d.recordList, offset, limit)
	mapSlice := cmhp_slice.Paginate(d.rawDataListAsMap, offset, limit)

	for i := 0; i < len(rawSlice); i++ {
		if recordSlice[i].isDeleted {
			continue
		}

		// Evaluate expression
		result, err := expression.Evaluate(mapSlice[i])
		if err != nil {
			panic(err)
		}

		if result.(bool) {
			out = append(out, recordSlice[i])
		}
	}

	if ch != nil {
		ch <- SearchResult[T]{
			dataBase: d,
			Count:    uint32(len(out)),
			IsFound:  len(out) > 0,
			List:     out,
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
