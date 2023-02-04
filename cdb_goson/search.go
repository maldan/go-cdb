package cdb_goson

import (
	"encoding/binary"
	"github.com/maldan/go-cdb/cdb_goson/goson"
	"github.com/maldan/go-cdb/cdb_proto/core"
	"github.com/maldan/go-cmhp/cmhp_slice"
	"strings"
)

func (d *DataTable[T]) ForEach(fn func(offset int, size int) bool) {
	offset := core.HeaderSize

	for {
		// Read size and flags
		size := int(binary.LittleEndian.Uint32(d.mem[offset+core.RecordStart:]))
		flags := int(d.mem[offset+core.RecordStart+core.RecordSize])

		// If field not deleted
		if flags&core.MaskDeleted != core.MaskDeleted {
			if !fn(offset, size) {
				break
			}
		}

		// Go to next value
		offset += size
		if offset >= len(d.mem) {
			break
		}
	}
}

type ArgsFind[T any] struct {
	FieldList string
	Limit     int
	Where     func(*T) bool
}

func (d *DataTable[T]) FindBy(args ArgsFind[T]) SearchResult[T] {
	// Return
	searchResult := SearchResult[T]{}

	// Create mapper for capturing values from bytes
	mapper := goson.NewMapper[T](d.Header.NameToId)

	// Field list
	fieldList := cmhp_slice.Map(strings.Split(args.FieldList, ","), func(t string) string {
		return strings.Trim(t, " ")
	})
	fieldIdList := make([]uint8, 0)
	for _, v := range fieldList {
		fieldIdList = append(fieldIdList, d.Header.NameToId[v])
	}

	// Go through each record
	d.ForEach(func(offset int, size int) bool {
		mapper.Map(d.mem[offset+core.RecordStart+core.RecordSize+core.RecordFlags:], fieldIdList)

		// Collect values
		if args.Where(&mapper.Container) {
			searchResult.table = d
			searchResult.Result = append(searchResult.Result, Record{offset: offset, size: size})

			// Check limit
			if args.Limit > 0 && len(searchResult.Result) >= args.Limit {
				return false
			}
		}

		return true
	})

	return searchResult
}
