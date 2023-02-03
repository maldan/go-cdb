package goson_test

import (
	"encoding/json"
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto/dson"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_print"
	"testing"
	"time"
	"unsafe"
)

type Record struct {
	Name  string
	Type  string
	Zone  string
	Gavno Gavno
}

type Gavno struct {
	Name int
	Type int
	Zone int
	Has  string
}

type Test struct {
	// Locked
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
	Balance  int       `json:"balance"`
	Created  time.Time `json:"created"`

	// Stripe
	StripeCustomerId     string `json:"stripeCustomerId"`
	StripeSubscriptionId string `json:"stripeSubscriptionId"`

	// Locked
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Record Record

	// X
	RecordList []Gavno
}

func TestMap(t *testing.T) {
	bytes := dson.Pack(Test{
		Email:   "sasageo",
		Role:    "123",
		Balance: 1,
		Record: Record{
			Name: "X", Type: "Y",
			Gavno: Gavno{Name: 1, Type: 1},
		},
		//Created: time.Now(),
		RecordList: []Gavno{
			{Name: 1, Type: 2, Zone: 3},
			{Name: 4, Type: 5, Zone: 6},
			{Name: 7, Type: 8, Zone: 9, Has: "XXX"},
		},
	})

	mapper := dson.NewMapper[Test]()
	mapper.Map(bytes[1:], []string{"Email"}, false)
	cmhp_print.Print(mapper.Container)
}

func TestMapSpeed(b *testing.T) {
	bytes := dson.Pack(Test{
		Email:   "sasageo",
		Role:    "123",
		Balance: 1,
		Record: Record{
			Name: "X", Type: "Y",
			Gavno: Gavno{Name: 1, Type: 1},
		},
		Created: time.Now(),
		RecordList: []Gavno{
			{Name: 1, Type: 2, Zone: 3},
			{Name: 4, Type: 5, Zone: 6},
			{Name: 7, Type: 8, Zone: 9, Has: "XXX"},
		},
	})

	mapper := dson.NewMapper[Test]()

	tt := time.Now()
	for i := 0; i < 1_000_000; i++ {
		mapper.Map(bytes[1:], []string{"Email"}, false)
	}
	fmt.Printf("%v\n", time.Since(tt))
}

func TestB(t *testing.T) {
	bytes := dson.Pack(map[string]any{
		"a": 1,
	})
	fmt.Printf("%v\n", bytes)
}

func TestPack(t *testing.T) {
	bytes := dson.Pack(Test{
		Email:   "sasageo",
		Balance: 1,
		Role:    "123",
		Record: Record{
			Name: "X", Type: "Y",
			Gavno: Gavno{Name: 1, Type: 1},
		},
		//Created: time.Now(),
		RecordList: []Gavno{
			{Name: 1, Type: 2, Zone: 3},
			{Name: 4, Type: 5, Zone: 6},
			{Name: 7, Type: 8, Zone: 9, Has: "XXX"},
		},
	})

	cmhp_file.Write("a.bin", bytes)
}

func TestA(t *testing.T) {
	x := 0
	for i := 0; i < 1024; i++ {
		bytes := dson.Pack(Test{
			Email:   "sasageo",
			Balance: 1,
			Role:    "123",
			Record: Record{
				Name: "X", Type: "Y",
				Gavno: Gavno{Name: 1, Type: 1},
			},
			//Created: time.Now(),
			RecordList: []Gavno{
				{Name: 1, Type: 2, Zone: 3},
				{Name: 4, Type: 5, Zone: 6},
				{Name: 7, Type: 8, Zone: 9, Has: "XXX"},
			},
		})

		tt := Test{}
		dson.UnpackX(bytes, unsafe.Pointer(&tt), tt)
		x += len(tt.Role)
	}
	fmt.Printf("'%v'\n", x)
	/*	cmhp_print.Print(tt)
		fmt.Printf("Time: %v\n", tt.Created)

		cmhp_file.Write("tt.json", tt)*/
}

func BenchmarkPack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes := dson.Pack(Test{
			Email:   "sasageo",
			Balance: 1,
			Role:    "123",
			Record: Record{
				Name: "X", Type: "Y",
				Gavno: Gavno{Name: 1, Type: 1},
			},
			Created: time.Now(),
			RecordList: []Gavno{
				{Name: 1, Type: 2, Zone: 3},
				{Name: 4, Type: 5, Zone: 6},
				{Name: 7, Type: 8, Zone: 9, Has: "XXX"},
			},
		})
		b.SetBytes(int64(len(bytes)))
	}
}

func BenchmarkMarshall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes, _ := json.Marshal(Test{
			Email:   "sasageo",
			Balance: 1,
			Role:    "123",
			Record: Record{
				Name: "X", Type: "Y",
				Gavno: Gavno{Name: 1, Type: 1},
			},
			Created: time.Now(),
			RecordList: []Gavno{
				{Name: 1, Type: 2, Zone: 3},
				{Name: 4, Type: 5, Zone: 6},
				{Name: 7, Type: 8, Zone: 9, Has: "XXX"},
			},
		})
		b.SetBytes(int64(len(bytes)))
	}
}

func BenchmarkZ(b *testing.B) {
	bytes, _ := json.Marshal(Test{
		Email:   "sasageo",
		Balance: 1,
		Role:    "123",
		Record: Record{
			Name: "X", Type: "Y",
			Gavno: Gavno{Name: 1, Type: 1},
		},
		Created: time.Now(),
		RecordList: []Gavno{
			{Name: 1, Type: 2, Zone: 3},
			{Name: 4, Type: 5, Zone: 6},
			{Name: 7, Type: 8, Zone: 9, Has: "XXX"},
		},
	})

	x := 0
	for i := 0; i < b.N; i++ {
		tt := Test{}
		json.Unmarshal(bytes, &tt)
		x = tt.Balance
	}
	fmt.Printf("Time: %v\n", x)
}

func BenchmarkX(b *testing.B) {
	bytes := dson.Pack(Test{
		Email:   "sasageo",
		Balance: 1,
		Role:    "123",
		Record: Record{
			Name: "X", Type: "Y",
			Gavno: Gavno{Name: 1, Type: 1},
		},
		Created: time.Now(),
		RecordList: []Gavno{
			{Name: 1, Type: 2, Zone: 3},
			{Name: 4, Type: 5, Zone: 6},
			{Name: 7, Type: 8, Zone: 9, Has: "XXX"},
		},
	})
	x := 0
	for i := 0; i < b.N; i++ {
		tt := Test{}
		dson.UnpackX(bytes, unsafe.Pointer(&tt), tt)
		x = tt.Balance
	}
	fmt.Printf("Time: %v\n", x)
}
