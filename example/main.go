package main

import (
	_ "net/http/pprof"
)

type Test2 struct {
	Id        uint32 `json:"id"`
	FirstName string `json:"first_name" sql:"firstName" len:"32"`
}

type Test struct {
	// gorm.Model

	Id    uint32 `json:"id"`
	A     string `json:"a" len:"32"`
	GasId int    `json:"gasId"`

	FirstName string `json:"first_name" sql:"firstName" len:"32"`
	/*LastName              string `json:"last_name" len:"32"`
	DateOfBirth           string `json:"date_of_birth" len:"32"`
	Age                   int8   `json:"age"`
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
	UserId                int    `json:"user_id"`
	DocumentId            int    `json:"document_id"`*/

	// ReviewId              int    `json:"review_id"`
	// ReviewEnabled         bool   `json:"review_enabled"`
	//Created time.Time `json:"created"`
}

func (t Test) GetId() uint32 {
	return uint32(t.Id)
}

/*func main() {
	c := cdb_lite.New[Test]("D:/go_lib/cdb/example/sas", nil)
	c.Init()
	tt := time.Now()
	c.Load()
	fmt.Printf("Full Load: %v\n", time.Since(tt))
	fmt.Printf("%v\n", c.Length())

	// 121 mb string + 248 struct - 369
	// 247 mb - buffer of 1 mln structs
	//

	// fmt.Printf("%v\n", unsafe.Sizeof(cdb_lite.Record{}))
	for i := 0; i < 3; i++ {
		tt = time.Now()
		r := c.FindBy(func(v *Test) bool { return v.Id == 32768 || v.Id == 750000 })
		fmt.Printf("Search: %v\n", time.Since(tt))
		fmt.Printf("%v\n", r)
	}

	for {
		time.Sleep(time.Second * 1)
	}
}
*/
