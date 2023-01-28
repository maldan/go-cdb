package cdb_proto_test

import (
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto"
	"testing"
	"time"
)

type Test struct {
	FirstName string `json:"firstName" id:"0"`
	LastName  string `json:"lastName" id:"1" len:"32"`
	Phone     string `json:"phone" id:"2" len:"64"`
}

func TestMyWrite(t *testing.T) {
	table := cdb_proto.New[Test]("../db/test")

	tt := time.Now()
	for i := 0; i < 1_000_000; i++ {
		table.Insert(Test{FirstName: fmt.Sprintf("%08d", i), LastName: "B", Phone: "C"})
	}
	fmt.Printf("%v\n", time.Since(tt))
}

func TestSimpleQuery(t *testing.T) {
	table := cdb_proto.New[Test]("../db/test")

	tt := time.Now()
	table.Query("SELECT * FROM table WHERE (FirstName == '00999999') AND (LastName == 'B')")
	fmt.Printf("%v\n", time.Since(tt))
}
