package cdb_goson

import (
	"encoding/binary"
	"github.com/maldan/go-cdb/cdb_goson/goson"
	"github.com/maldan/go-cdb/cdb_proto/core"
)

func (d *DataTable[T]) ForEach(fn func(offset int, size int) bool) {
	offset := core.HeaderSize

	for {
		size := int(binary.LittleEndian.Uint32(d.mem[offset+core.RecordStart:]))
		flags := int(d.mem[offset+core.RecordStart+core.RecordSize])

		// If field not deleted
		if flags&core.MaskDeleted != core.MaskDeleted {
			if !fn(offset, size) {
				break
			}
		}

		offset += size
		if offset >= len(d.mem) {
			break
		}
	}
}

func (d *DataTable[T]) Find(fieldList []string, where func(*T) bool) SearchResult[T] {
	// Return
	searchResult := SearchResult[T]{}

	mapper := goson.NewMapper[T]()

	d.ForEach(func(offset int, size int) bool {
		mapper.Map(d.mem[offset+core.RecordStart+core.RecordSize+core.RecordFlags:], fieldList, false)

		if where(&mapper.Container) {
			searchResult.table = d
			searchResult.Result = append(searchResult.Result, Record{offset: offset, size: size})
			return false
		}

		return true
	})

	return searchResult
}
