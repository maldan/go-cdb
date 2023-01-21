package cdb_lite_test

import (
	"fmt"
	"github.com/maldan/go-cdb/cdb_lite"
	"testing"
	"time"
)

type Test struct {
	Id uint32 `json:"id"`
	A  string `json:"a"`
	//GasId int    `json:"gasId"`

	FirstName   string `json:"first_name" length:"32"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	//Age                   int    `json:"age"`
	Address               string `json:"address"`
	City                  string `json:"city"`
	State                 string `json:"state"`
	ZipCode               string `json:"zip_code"`
	Email                 string `json:"email"`
	Phone                 string `json:"phone"`
	EmergencyPersonName   string `json:"emergency_person_name"`
	EmergencyPersonPhone  string `json:"emergency_person_phone"`
	ParentOrLegalGuardian string `json:"parent_or_legal_guardian"`
	Date                  string `json:"date"`
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

func TestMyWrite(t *testing.T) {
	c := cdb_lite.New[Test]()
	c.Init()

	for i := 0; i < 1; i++ {
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
	}
	c.Flush()
}

func TestMyRead(t *testing.T) {
	c := cdb_lite.New[Test]()
	c.Init()

	tt := time.Now()
	c.Load(func(test Test) {
		fmt.Printf("%v\n", test)
	})
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
