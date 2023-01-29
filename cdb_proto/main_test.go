package cdb_proto_test

import (
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto"
	"github.com/maldan/go-cmhp/cmhp_print"
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
		table.Insert(Test{
			FirstName: fmt.Sprintf("LO %08d", i),
			LastName:  fmt.Sprintf("EB %08d", i),
			Phone:     fmt.Sprintf("DE %08d", i),
		})
	}
	fmt.Printf("%v\n", time.Since(tt))
}

func TestSimpleQuery(t *testing.T) {
	table := cdb_proto.New[Test]("../db/test")

	tt := time.Now()
	rs := table.Query("SELECT * FROM table WHERE FirstName == 'LO 00999999'")
	fmt.Printf("T1: %v\n", time.Since(tt))

	oo := rs.Unpack()
	cmhp_print.Print(oo)

	tt = time.Now()
	table.Query("SELECT * FROM table WHERE FirstName == 'LO 00999999'")
	fmt.Printf("T2: %v\n", time.Since(tt))
}
