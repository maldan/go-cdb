package cdb_lite

import (
	"github.com/maldan/go-cmhp/cmhp_byte"
)

func (d *DataEngine[T]) Init() *DataEngine[T] {
	d.indexList = make(map[string][]uint32, 0)

	return d
}

func (d *DataEngine[T]) Flush() {
	d.storage.Flush()
}

func (d *DataEngine[T]) Load() {
	d.storage.Load(func(v *T) {
		d.loadRecord(v)
	})
}

func New[T IEngineComparable](name string, result []string) *DataEngine[T] {
	n := DataEngine[T]{
		Name:              name,
		SearchFieldByList: result,
	}
	cmhp_byte.Pack[T](new(T)) // cache type + check
	return &n
}
