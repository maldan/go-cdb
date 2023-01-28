package cdb_proto_test

import (
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto"
	"testing"
)

type Test struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

func TestMyWrite(t *testing.T) {
	table := cdb_proto.New[Test]("../db/test")

	for i := 0; i < 10; i++ {
		table.Insert(Test{FirstName: "A", LastName: "B", Phone: "C"})
	}
	fmt.Printf("%v", &table)
}
