package cdb_flat_test

import (
	"github.com/maldan/go-cdb/cdb_flat"
	"testing"
)

func BenchmarkAdd(b *testing.B) {
	c := cdb_flat.DataBase[Test]{ChunkAmount: 10}
	c.Init()

	for i := 0; i < b.N; i++ {
		c.Add(Test{Id: 1, A: "X"})
	}
}
