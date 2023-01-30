package main

import (
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto"
	"github.com/maldan/go-cmhp/cmhp_print"
	"time"
)

type Test struct {
	FirstName string `json:"firstName" id:"0"`
	LastName  string `json:"lastName" id:"1" len:"32"`
	Phone     string `json:"phone" id:"2" len:"64"`
	Sex       string `json:"sex" id:"3" len:"64"`
	Rock      string `json:"rock" id:"4" len:"64"`
	Gas       string `json:"gas" id:"5" len:"64"`
	Yas       string `json:"yas" id:"6" len:"64"`
	Taj       string `json:"taj" id:"7" len:"64"`
	Mahal     string `json:"mahal" id:"8" len:"64"`
	Ebal      string `json:"ebal" id:"9" len:"64"`
	Sasal     string `json:"sasal" id:"10" len:"64"`
	Sasal2    string `json:"sasal2" id:"11" len:"64"`
}

func f1() {
	table := cdb_proto.New[Test]("../db/test")

	tt := time.Now()
	rs := table.Query("SELECT * FROM table WHERE FirstName == '00999999'")
	fmt.Printf("T1: %v\n", time.Since(tt))

	oo := rs.Unpack()
	cmhp_print.Print(oo)

	tt = time.Now()
	table.Query("SELECT * FROM table WHERE FirstName == '00999999'")
	fmt.Printf("T2: %v\n", time.Since(tt))
}

func f2() {
	table := cdb_proto.New[Test]("../db/test")

	tt := time.Now()
	rs := table.CrazySelect([]string{"FirstName"}, func(test *Test) bool {
		return test.FirstName == "00999999"
	})
	fmt.Printf("T1: %v\n", time.Since(tt))

	oo := rs.Unpack()
	cmhp_print.Print(oo)

	tt = time.Now()
	table.CrazySelect([]string{"FirstName"}, func(test *Test) bool {
		return test.FirstName == "00999999"
	})
	fmt.Printf("T1: %v\n", time.Since(tt))
}

func main() {
	f1()
	f2()
}
