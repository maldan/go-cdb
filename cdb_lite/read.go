package cdb_lite

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"reflect"
)

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

func (d *DataEngine[T]) fastConvert(v *T) map[string]any {
	m := map[string]any{}

	typeOf := reflect.TypeOf(v).Elem()
	valueOf := reflect.ValueOf(v).Elem()
	for i := 0; i < len(d.SearchFieldByList); i++ {
		m[typeOf.Field(i).Name] = valueOf.FieldByName(d.SearchFieldByList[i]).Interface()
	}
	return m
}

func (d *DataEngine[T]) FindByQuery(query string) SearchResult[T] {
	d.dataMu.RLock()
	defer d.dataMu.RUnlock()

	// Prepare out
	out := make([]Record, 0)

	// Prepare expression
	expression, err := govaluate.NewEvaluableExpression(query)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(d.rawDataList); i++ {
		if d.recordList[i].isDeleted {
			continue
		}

		// Evaluate expression
		result, err := expression.Evaluate(d.rawDataListAsMap[i])
		if err != nil {
			panic(err)
		}

		if result.(bool) {
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
