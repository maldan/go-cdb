package cdb_lite

import (
	"github.com/maldan/go-cdb/engine"
	"github.com/maldan/go-cdb/pack"
)

func New[T engine.IEngineComparable]() *engine.DataEngine[T] {
	n := engine.DataEngine[T]{}
	pack.Pack[T](new(T)) // cache type + check
	return &n
}
