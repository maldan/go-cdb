package cdb_proto_test

import (
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto"
	"testing"
)

type Test struct {
	FirstName string `json:"firstName" id:"0"`
	LastName  string `json:"lastName" id:"1" len:"32"`
	Phone     string `json:"phone" id:"2" len:"64"`
}

func TestMyWrite(t *testing.T) {
	table := cdb_proto.New[Test]("../db/test")

	for i := 0; i < 10; i++ {
		table.Insert(Test{FirstName: "A", LastName: "B", Phone: "C"})
	}
	fmt.Printf("%v", &table)
}
