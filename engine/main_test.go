package engine_test

import (
	"testing"
	"time"
)

type Test struct {
	Id    uint32 `json:"id"`
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

func (t Test) GetId() uint32 {
	return t.Id
}

func TestA(t *testing.T) {
	//c := cdb_core.DataBase[Test]{IndexList: []string{"GasId"}}
	//c.Init()
}
