package dson_test

import (
	"github.com/maldan/go-cdb/cdb_proto/dson"
	"testing"
	"time"
)

type Record struct {
	Name string
	Type string
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

	// X
	RecordList []Record
}

/*func TestDynLength(t *testing.T) {
	if len(dson.CreateDynLength(32)) != 1 {
		t.Fatalf("Fuck")
	}
	if len(dson.CreateDynLength(128)) != 2 {
		t.Fatalf("Fuck")
	}

	for i := 0; i < 128; i++ {
		if len(dson.CreateDynLength(i)) != 1 {
			t.Fatalf("Fuck")
		}
		bt := dson.CreateDynLength(i)
		tb, s := dson.ReadDynLength(bt)

		if s != 1 {
			t.Fatalf("Fuck size")
		}
		if i != tb {
			t.Fatalf("Fuck value")
		}
	}

	for i := 128; i < 65536; i++ {
		if len(dson.CreateDynLength(i)) != 2 {
			t.Fatalf("Fuck %v", i)
		}
		bt := dson.CreateDynLength(i)
		tb, s := dson.ReadDynLength(bt)

		if s != 2 {
			t.Fatalf("Fuck size")
		}
		if i != tb {
			t.Fatalf("Fuck value %v", i)
		}
	}
}*/

func TestA(t *testing.T) {
	dson.Pack(Test{
		Email:   "sasageo",
		Balance: 1,
		Role:    "123",
		Created: time.Now(),
		RecordList: []Record{
			{Name: "X", Type: "Y"},
		},
	})
}
