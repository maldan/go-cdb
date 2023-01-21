package cdb_lite

import (
	"github.com/maldan/go-cdb/engine"
	"github.com/maldan/go-cmhp/cmhp_byte"
)

func New[T engine.IEngineComparable]() *engine.DataEngine[T] {
	n := engine.DataEngine[T]{}
	cmhp_byte.Pack[T](new(T)) // cache type + check
	return &n
}
