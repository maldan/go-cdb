package parse_test

import (
	"fmt"
	"github.com/maldan/go-cdb/cdb_proto/parse"
	"testing"
)

type Test struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
}

func TestParseSelect(t *testing.T) {
	q, err := parse.Query[Test]("SELECT * FROM table WHERE FirstName == 'Lox'")
	fmt.Printf("%+v\n", q)
	fmt.Printf("%v\n", err)
}
