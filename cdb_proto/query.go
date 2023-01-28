package cdb_proto

import (
	"encoding/binary"
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"github.com/maldan/go-cdb/cdb_proto/pack"
	"github.com/maldan/go-cdb/cdb_proto/parse"
)

func (d *DataTable[T]) Query(query string) {
	q, _ := parse.Query[T](query)

	d.Select(q)
}

func OpStrCmp(left func() string, right func() string) bool {
	a := left()
	b := right()
	return a == b
}

func OpAnd(left func() any, right func() any) bool {
	a := left()
	b := right()
	return a == b
}

type E2 struct {
	Mem   []byte
	Table []byte
	Op    string

	LeftId    int
	LeftSize  int
	LeftType  int
	LeftValue any
	LeftE2    *E2

	RightId    int
	RightSize  int
	RightType  int
	RightValue any
	RightE2    *E2
}

func (e E2) Do() {

}

func pop[T any](s *[]T) T {
	v := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return v
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

func (d *DataTable[T]) Select(query parse.QueryInfo) {
	offset := core.HeaderSize

	stack := passStack{
		stack:     make([]bool, 32),
		byteStack: make([][]byte, 32),
	}

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
		/*size, fieldLen, startData := pack.ReadHeader(d.mem, offset, fieldOffsetIndex)
		fieldData := d.mem[startData : startData+fieldLen]
		isFound := comparator(strAsBytes, fieldData)*/

		if isFound {
			fmt.Printf("Fount at %v\n", offset)
			break
		}

		offset += size
		if offset >= len(d.mem) {
			break
		}
	}
}
