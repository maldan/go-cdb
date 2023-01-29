package cdb_proto

import (
	"encoding/binary"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"github.com/maldan/go-cdb/cdb_proto/pack"
	"github.com/maldan/go-cdb/cdb_proto/parse"
)

func (d *DataTable[T]) Query(query string) SearchResult[T] {
	q, _ := parse.Query[T](query)

	return d.Select(q)
}

type passStack struct {
	stack     []bool
	byteStack [][]byte
}

func CheckExpression[T any](d *DataTable[T], offset int, table []byte, ops []parse.TokenType, st passStack) bool {
	if len(ops) == 0 {
		return true
	}

	stackCounter := 0

	stack := st.stack
	byteStack := st.byteStack

	for i := 0; i < len(ops); i++ {
		v := ops[i]

		switch v.Type {
		case core.TokenIdentifier:
			vOff := int(binary.LittleEndian.Uint32(table[v.TableOffset:]))
			vLen := int(binary.LittleEndian.Uint32(table[v.TableOffset+4:]))
			byteStack[stackCounter] = d.mem[offset+vOff : offset+vOff+vLen]
			stackCounter += 1
			break
		case core.TokenString:
			byteStack[stackCounter] = v.B
			stackCounter += 1
			break
		case core.TokenOp:
			switch v.Token {
			case "==":
				stackCounter -= 1
				b := byteStack[stackCounter]
				stackCounter -= 1
				a := byteStack[stackCounter]

				stack[stackCounter] = StrCmp(a, b)
				stackCounter += 1
				break
			case "&&":
				stackCounter -= 1
				b := stack[stackCounter]
				stackCounter -= 1
				a := stack[stackCounter]

				stack[stackCounter] = a == b
				stackCounter += 1
				break
			}
			break
		}
	}

	return stack[0]
}

func (d *DataTable[T]) Select(query parse.QueryInfo) SearchResult[T] {
	offset := core.HeaderSize

	stack := passStack{
		stack:     make([]bool, 32),
		byteStack: make([][]byte, 32),
	}

	// Return
	searchResult := SearchResult[T]{}

	for i := 0; i < len(query.WhereCondition); i++ {
		if query.WhereCondition[i].Type == core.TokenIdentifier {
			id, ok := d.structInfo.FieldNameToId[query.WhereCondition[i].Token]
			if !ok {
				panic("DA!")
			}
			query.WhereCondition[i].TableOffset = id * 8
		}
	}

	for {
		size, offTable := pack.ReadHeader2(d.mem, offset)

		isFound := CheckExpression(d, offset, offTable, query.WhereCondition, stack)

		// fmt.Printf("%v\n", isFound)
		/*
			size, fieldLen, startData := pack.ReadHeader(d.mem, offset, fieldOffsetIndex)
			fieldData := d.mem[startData : startData+fieldLen]
			isFound := comparator(strAsBytes, fieldData)
		*/

		if isFound {
			searchResult.table = d
			searchResult.Result = append(searchResult.Result, Record{offset: offset, size: size})
			break
		}

		offset += size
		if offset >= len(d.mem) {
			break
		}
	}

	return searchResult
}
