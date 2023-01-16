package cdb_flat_test

import (
	"github.com/maldan/go-cdb/cdb_flat"
	"testing"
	"time"
)

type Test struct {
	Id    uint   `json:"id"`
	A     string `json:"a"`
	GasId int    `json:"gasId"`

	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	DateOfBirth           string `json:"date_of_birth"`
	Age                   int    `json:"age"`
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
	UserId                int    `json:"user_id"`
	DocumentId            int    `json:"document_id"`
	ReviewId              int    `json:"review_id"`
	ReviewEnabled         bool   `json:"review_enabled"`

	Created time.Time `json:"created"`
}

func (t Test) GetId() uint {
	return t.Id
}

func TestDelete(t *testing.T) {
	c := cdb_flat.DataBase[Test]{ChunkAmount: 10, IndexList: []string{"GasId"}}
	c.Init()

	c.Add(Test{Id: 1, A: "X", GasId: 1})
	c.Add(Test{Id: 2, A: "Y", GasId: 1})
	c.Add(Test{Id: 3, A: "Z", GasId: 2})
	for i := 0; i < 10; i++ {
		c.Add(Test{Id: uint(4 + i), A: "U", GasId: 3})
	}

	if c.Length() != 13 {
		t.Errorf("Incorrect size %v", c.Length())
	}

	// Find and delete
	r2 := c.FindBy(func(x *Test) bool { return x.Id == 2 })
	c.Delete(r2.List[0])

	// Check delete
	r3 := c.FindBy(func(x *Test) bool { return true })
	if r3.Count == 3 {
		t.Errorf("Delete not working")
	}
}

func TestFindByIndex(t *testing.T) {
	c := cdb_flat.DataBase[Test]{ChunkAmount: 10, IndexList: []string{"GasId"}}
	c.Init()

	c.Add(Test{Id: 1, A: "X", GasId: 1})
	c.Add(Test{Id: 2, A: "Y", GasId: 1})
	c.Add(Test{Id: 3, A: "Z", GasId: 2})
	for i := 0; i < 10; i++ {
		c.Add(Test{Id: uint(4 + i), A: "U", GasId: 3})
	}

	if c.Length() != 13 {
		t.Errorf("Incorrect size %v", c.Length())
	}

	r := c.FindByIndex("GasId", 1)
	if r.Count != 2 {
		t.Errorf("Index not working")
	}
	u := r.UnpackList()
	if u[0].A != "X" {
		t.Errorf("Index not working")
	}
	if u[1].A != "Y" {
		t.Errorf("Index not working")
	}

	// Search
	r2 := c.FindBy(func(x *Test) bool { return x.Id == 2 })
	c.Delete(r2.List[0])

	// Test index again
	r2 = c.FindByIndex("GasId", 1)
	if r2.Count != 1 {
		t.Errorf("Index not working")
	}
}

func TestFindById(t *testing.T) {
	c := cdb_flat.DataBase[Test]{ChunkAmount: 10, IndexList: []string{"GasId"}}
	c.Init()

	c.Add(Test{Id: 1, A: "X", GasId: 1})
	c.Add(Test{Id: 2, A: "Y", GasId: 1})
	c.Add(Test{Id: 3, A: "Z", GasId: 2})
	for i := 0; i < 10; i++ {
		c.Add(Test{Id: uint(4 + i), A: "U", GasId: 3})
	}

	if c.Length() != 13 {
		t.Errorf("Incorrect size %v", c.Length())
	}

	r := c.FindById(3)
	if r.Count != 1 {
		t.Errorf("Index not working")
	}
	u := r.UnpackList()
	if u[0].A != "Z" {
		t.Errorf("Index not working")
	}

	// Search
	r2 := c.FindBy(func(x *Test) bool { return x.Id == 2 })
	c.Delete(r2.List[0])

	// Test index again
	r2 = c.FindById(2)
	if r2.IsFound {
		t.Errorf("Index not working")
	}
}
