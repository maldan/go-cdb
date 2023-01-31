package harray_test

import (
	"fmt"
	"github.com/cespare/xxhash"
	"github.com/maldan/go-cdb/cdb_proto/harray"
	"github.com/maldan/go-cmhp/cmhp_print"
	"testing"
	"time"
)

func TestG(t *testing.T) {
	harr := harray.New()
	harr.Add("a", 1)
	cmhp_print.PrintBytesColored(harr.Memory, 32, []cmhp_print.ColorRange{
		{0, 2, cmhp_print.BgRed},
		{2, 4, cmhp_print.BgGreen},
	})
}

func TestCC(t *testing.T) {
	fmt.Printf("%v\n", harray.Hash("0"))
	fmt.Printf("%v\n", harray.Hash("1"))
	fmt.Printf("%v\n", harray.Hash("2"))
	fmt.Printf("%v\n", harray.Hash("ab"))
	fmt.Printf("%v\n", harray.Hash("ba"))
}

func TestX(t *testing.T) {
	x := uint64(0)
	tt := time.Now()
	for i := 0; i < 1_000_000; i++ {
		x += xxhash.Sum64String("ab5")
	}
	fmt.Printf("%v - %v\n", x, time.Since(tt))
}

func TestY(t *testing.T) {
	x := uint64(0)
	tt := time.Now()
	for i := 0; i < 1_000_000; i++ {
		x += uint64(harray.Hash("ab5"))
	}
	fmt.Printf("%v - %v\n", x, time.Since(tt))
}

func BenchmarkX(b *testing.B) {
	x := uint64(0)
	for i := 0; i < b.N; i++ {
		x += uint64(harray.Hash("ab5"))
	}
	fmt.Printf("%v\n", x)
}

func BenchmarkY(b *testing.B) {
	x := uint64(0)
	for i := 0; i < b.N; i++ {
		x += xxhash.Sum64String("ab5")
	}
	fmt.Printf("%v\n", x)
}
