package dson_test

import (
	"encoding/json"
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto/dson"
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

func TestB(t *testing.T) {
	bytes := dson.Pack(map[string]any{
		"a": 1,
	})
	fmt.Printf("%v\n", bytes)
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
		dson.Unpack(bytes, unsafe.Pointer(&tt), tt)
		x += len(tt.Role)
	}
	fmt.Printf("'%v'\n", x)
	/*	cmhp_print.Print(tt)
		fmt.Printf("Time: %v\n", tt.Created)

		cmhp_file.Write("tt.json", tt)*/
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
		dson.Unpack(bytes, unsafe.Pointer(&tt), tt)
		x = tt.Balance
	}
	fmt.Printf("Time: %v\n", x)
}
