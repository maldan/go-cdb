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

func TestMyWrite(t *testing.T) {
	table := cdb_proto.New[Test]("../db/test")

	tt := time.Now()
	for i := 0; i < 1_000_000; i++ {
		table.Insert(Test{
			FirstName: fmt.Sprintf("%08d", i),
			LastName:  "00000000",
			Phone:     "11111111",
			Sex:       "00000000",
			Rock:      "00000000",
			Gas:       "00000000",
			Yas:       "00000000",
			Taj:       "00000000",
			Mahal:     "00000000",
			Ebal:      "00000000",
			Sasal:     "00000000",
			Sasal2:    "XXXXXXXX",
		})
	}
	fmt.Printf("%v\n", time.Since(tt))
}

func TestSimpleQuery(t *testing.T) {
	table := cdb_proto.New[Test]("../db/test")

	tt := time.Now()
	rs := table.Query("SELECT * FROM table WHERE FirstName == '00999999' AND Phone == '1234567'")
	fmt.Printf("T1: %v\n", time.Since(tt))

	oo := rs.Unpack()
	cmhp_print.Print(oo)

	/*tt = time.Now()
	table.Query("SELECT * FROM table WHERE FirstName == '00999999'")
	fmt.Printf("T2: %v\n", time.Since(tt))*/
}

func TestCrazyQuery(t *testing.T) {
	table := cdb_proto.New[Test]("../db/test")

	tt := time.Now()
	rs := table.Find([]string{"FirstName"}, func(test *Test) bool {
		return test.FirstName == "00000000"
	})
	fmt.Printf("T1: %v\n", time.Since(tt))

	oo := rs.Unpack()
	cmhp_print.Print(oo)

	tt = time.Now()
	table.Find([]string{"FirstName", "Phone"}, func(test *Test) bool {
		return test.FirstName == "..."
	})
	fmt.Printf("T1: %v\n", time.Since(tt))
}

func TestDelete(t *testing.T) {
	table := cdb_proto.New[Test]("../db/test")

	table.Delete([]string{"FirstName"}, func(test *Test) bool {
		return test.FirstName == "00000000"
	})
}

func TestX3(t *testing.T) {
	a := 0
	tt := time.Now()
	for i := 0; i < 1_000_000; i++ {
		a = 0b0000_1111
		for j := 0; j < 64; j++ {
			a = a << 1
		}
	}
	fmt.Printf("T2: %v\n", time.Since(tt))
	fmt.Printf("%v\n", a)
}

func TestX2(t *testing.T) {
	type Gas struct {
		Id int
		S  []string
	}

	xyz := make([]Gas, 1_000_000)
	for i := 0; i < 1_000_000; i++ {
		xyz[i].Id = 1_000_000 - 1
	}

	a := 0
	tt := time.Now()
	for i := 0; i < 1_000_000; i++ {
		a += xyz[i].Id
	}
	fmt.Printf("T2: %v\n", time.Since(tt))
	fmt.Printf("%v\n", a)
	/*fmt.Printf("%v\n", xyz[0])

	tt := time.Now()
	sort.SliceStable(xyz, func(i, j int) bool {
		return xyz[i].Id > xyz[j].Id
	})
	fmt.Printf("T2: %v\n", time.Since(tt))

	fmt.Printf("%v\n", xyz[0])

	for {
		time.Sleep(time.Second)
	}*/
}
