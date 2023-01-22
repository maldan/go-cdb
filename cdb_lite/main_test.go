package cdb_lite_test

import (
	"fmt"
	"github.com/maldan/go-cdb/cdb_lite"
	"testing"
	"time"
)

type Test struct {
	Id uint32 `json:"id"`
	A  string `json:"a" len:"32"`
	//GasId int    `json:"gasId"`

	FirstName   string `json:"first_name" len:"32"`
	LastName    string `json:"last_name" len:"32"`
	DateOfBirth string `json:"date_of_birth" len:"32"`
	//Age                   int    `json:"age"`
	Address               string `json:"address" len:"32"`
	City                  string `json:"city" len:"32"`
	State                 string `json:"state" len:"32"`
	ZipCode               string `json:"zip_code" len:"32"`
	Email                 string `json:"email" len:"32"`
	Phone                 string `json:"phone" len:"32"`
	EmergencyPersonName   string `json:"emergency_person_name" len:"32"`
	EmergencyPersonPhone  string `json:"emergency_person_phone" len:"32"`
	ParentOrLegalGuardian string `json:"parent_or_legal_guardian" len:"32"`
	Date                  string `json:"date" len:"32"`
	//UserId                int    `json:"user_id"`
	//DocumentId            int    `json:"document_id"`
	//ReviewId              int    `json:"review_id"`
	//ReviewEnabled         bool   `json:"review_enabled"`

	//Created time.Time `json:"created"`
}

func (t Test) GetId() uint32 {
	return t.Id
}

/*func TestB(t *testing.T) {
	list := make([]Test, 0)
	for i := 0; i < 1_000_000; i++ {
		list = append(list, Test{
			Id:          uint32(i),
			Phone:       "+79961156919",
			FirstName:   "Roman",
			LastName:    "Moldovan",
			DateOfBirth: "1992-08-28",
			Address:     "lox",
			City:        "Zvenigovo", State: "Marii-El", ZipCode: "425061",
			Email:                "blackwanted@yandex.ru",
			Date:                 "2000-01-01 00:00:01",
			EmergencyPersonPhone: "+79961156919",
			EmergencyPersonName:  "Ganvarr",
		})
	}

	tt := time.Now()
	util.WriteJson("sas.json", &list)
	fmt.Printf("Write: %v\n", time.Since(tt))
}

func TestRead(t *testing.T) {
	tt := time.Now()
	list, _ := util.ReadJson[[]Test]("sas.json")
	fmt.Printf("Read: %v\n", time.Since(tt))
	fmt.Printf("%v\n", list[0])
}
*/

// []string{"Id", "FirstName", "LastName", "Date", "Phone", "Address", "City"}

func TestFindByQuery(t *testing.T) {
	c := cdb_lite.New[Test]("", []string{"Id"})
	c.Init()
	tt := time.Now()
	c.Load()
	fmt.Printf("Full Load: %v\n", time.Since(tt))
	fmt.Printf("%v\n", c.Length())

	for i := 0; i < 3; i++ {
		tt = time.Now()
		r := c.FindByQuery("Id==32768 || Id==750000", nil, 0, 1_000_000)
		fmt.Printf("Search: %v\n", time.Since(tt))
		fmt.Printf("%v\n", r)
	}
}

func TestFindByQueryParallel(t *testing.T) {
	c := cdb_lite.New[Test]("", nil)
	c.Init()
	tt := time.Now()
	c.Load()
	fmt.Printf("Full Load: %v\n", time.Since(tt))
	fmt.Printf("%v\n", c.Length())

	for i := 0; i < 3; i++ {
		tt = time.Now()
		r := c.FindByQueryParallel("Id==32768 || Id==750000")
		fmt.Printf("Search: %v\n", time.Since(tt))
		fmt.Printf("%v\n", r)
	}
}

func TestFindBy(t *testing.T) {
	c := cdb_lite.New[Test]("", nil)
	c.Init()
	tt := time.Now()
	c.Load()
	fmt.Printf("Full Load: %v\n", time.Since(tt))
	fmt.Printf("%v\n", c.Length())

	for i := 0; i < 3; i++ {
		tt = time.Now()
		r := c.FindBy(func(v *Test) bool { return v.Id == 32768 || v.Id == 750000 })
		fmt.Printf("Search: %v\n", time.Since(tt))
		fmt.Printf("%v\n", r)
	}
}

func TestMyWrite(t *testing.T) {
	c := cdb_lite.New[Test]("", nil)
	c.Init()

	for i := 0; i < 1_000_000; i++ {
		id := c.GenerateId()
		c.Add(&Test{
			Id:                   id,
			Phone:                "+79961156919",
			FirstName:            "Roman",
			LastName:             "Moldovan",
			DateOfBirth:          "1992-08-28",
			Address:              "lox",
			City:                 "Zvenigovo",
			State:                "Marii-El",
			ZipCode:              "425061",
			Email:                "blackwanted@yandex.ru",
			Date:                 "2000-01-01 00:00:01",
			EmergencyPersonPhone: "+79961156919",
			EmergencyPersonName:  "Ganvarr",
		})
	}
	c.Flush()
}

func TestMyRead(t *testing.T) {
	c := cdb_lite.New[Test]("", nil)
	c.Init()

	tt := time.Now()
	c.Load()
	fmt.Printf("Read: %v\n", time.Since(tt))

	/*tt := time.Now()
	for i := 0; i < 4; i++ {
		id := c.GenerateId()
		c.Add(Test{
			Id:          id,
			Phone:       "+79961156919",
			FirstName:   "Roman",
			LastName:    "Moldovan",
			DateOfBirth: "1992-08-28",
			Address:     "lox",
			City:        "Zvenigovo", State: "Marii-El", ZipCode: "425061",
			Email:                "blackwanted@yandex.ru",
			Date:                 "2000-01-01 00:00:01",
			EmergencyPersonPhone: "+79961156919",
			EmergencyPersonName:  "Ganvarr",
		})
		//r := c.FindById(id)
		//c.Delete(r.List[0])
	}
	fmt.Printf("Ops: %v\n", time.Since(tt))

	tt = time.Now()
	c.Flush()
	fmt.Printf("Total: %v\n", time.Since(tt))*/
}
