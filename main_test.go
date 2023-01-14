package cdb_test

import (
	"fmt"
	"github.com/maldan/go-cdb"
	"testing"
	"time"
)

type Test struct {
	Id    int    `json:"id"`
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

type Test2 struct {
	Id int            `json:"id"`
	Ma map[string]any `json:"ma"`
}

func (t Test) GetId() any {
	return t.Id
}

func (t Test2) GetId() any {
	return t.Id
}

func TestAdd(t *testing.T) {
	testChunk := cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.Init()

	// Add
	for i := 0; i < 10000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}
	if testChunk.TotalElements() != 10000 {
		t.Errorf("Element amount not match")
	}

	// Test if changed
	if !testChunk.ChunkList[9].IsChanged {
		t.Errorf("Chunk must be not changed")
	}
}

func TestDelete(t *testing.T) {
	testChunk := cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.Init()

	// Add
	for i := 0; i < 10000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}

	// Delete
	testChunk.GetChunkByHash(1).DeleteBy(func(t *Test) bool {
		return t.Id == 1
	})
	if testChunk.TotalElements() != 10000-1 {
		t.Errorf("Delete not working")
	}

	// Delete in all
	testChunk.GetChunkByHash(555).DeleteBy(func(t *Test) bool {
		return t.Id == 555
	})
	if testChunk.TotalElements() != 10000-2 {
		t.Errorf("Delete not working")
	}

	// Find
	_, ok := testChunk.GetChunkByHash(1).Find(func(t *Test) bool {
		return t.Id == 1
	})
	if ok {
		t.Errorf("Delete not working")
	}
	_, ok = testChunk.GetChunkByHash(555).Find(func(t *Test) bool {
		return t.Id == 555
	})
	if ok {
		t.Errorf("Delete not working")
	}
	_, ok = testChunk.Find(func(t *Test) bool {
		return t.Id == 1
	})
	if ok {
		t.Errorf("Delete not working")
	}
	_, ok = testChunk.Find(func(t *Test) bool {
		return t.Id == 555
	})
	if ok {
		t.Errorf("Delete not working")
	}
}

func TestUpdate(t *testing.T) {
	testChunk := cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.Init()

	// Add
	testChunk.Add(Test{Id: 1})
	testChunk.Add(Test{Id: 2})
	testChunk.Add(Test{Id: 3})
	if testChunk.TotalElements() != 3 {
		t.Errorf("Fuck you")
	}

	// Update
	testChunk.Replace(Test{Id: 1, A: "gas"})
	if testChunk.TotalElements() != 3 {
		t.Errorf("Fuck you")
	}

	// Test if changed
	if testChunk.ChunkList[1].List[0].A != "gas" {
		t.Errorf("Update not working")
	}
}

func TestFind(t *testing.T) {
	testChunk := cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.Init()

	// Add
	for i := 0; i < 1000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}
	if testChunk.TotalElements() != 1000 {
		t.Errorf("Fuck you")
	}

	// Find in chunk
	v, ok := testChunk.GetChunkByHash(432).Find(func(t *Test) bool {
		return t.Id == 432
	})
	if !ok {
		t.Errorf("Value not found")
	}

	// Test if changed
	if v.Id != 432 {
		t.Errorf("Find not working")
	}

	// Find in all
	v, ok = testChunk.Find(func(t *Test) bool {
		return t.Id == 768
	})
	if !ok {
		t.Errorf("Value not found")
	}

	// Test if changed
	if v.Id != 768 {
		t.Errorf("Find not working")
	}
}

func TestFilter(t *testing.T) {
	testChunk := cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.Init()

	// Add
	for i := 0; i < 3; i++ {
		testChunk.Add(Test{Id: i + 1})
	}
	if testChunk.TotalElements() != 3 {
		t.Errorf("Fuck you")
	}

	// Find in chunk
	list := testChunk.FindMany(func(t *Test) bool {
		return t.Id > 2
	})
	if len(list) != 1 {
		t.Errorf(fmt.Sprintf("Fuck you %v", len(list)))
	}
}

func TestIndex(t *testing.T) {
	testChunk := cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas", IndexList: []string{"GasId"}}
	testChunk.Init()

	testChunk.Add(Test{Id: 1, A: "X", GasId: 1})
	testChunk.Add(Test{Id: 2, A: "Y", GasId: 1})
	testChunk.Add(Test{Id: 3, A: "Z", GasId: 2})
	testChunk.Add(Test{Id: 4, A: "W", GasId: 3})

	l := testChunk.FindManyByIndex("GasId", 2)
	if len(l) == 0 {
		t.Errorf("Index not working")
	}
	for _, x := range l {
		if x.GasId != 2 {
			t.Errorf("Index not working")
		}
	}

	l = testChunk.FindManyByIndex("GasId", 1)
	if len(l) == 0 {
		t.Errorf("Index not working")
	}
	for _, x := range l {
		if x.GasId != 1 {
			t.Errorf("Index not working")
		}
	}
}

func TestIndexProblem(t *testing.T) {
	testChunk := cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas", IndexList: []string{"GasId"}}
	testChunk.Init()

	testChunk.Add(Test{Id: 1, A: "X", GasId: 1})
	testChunk.Add(Test{Id: 2, A: "Y", GasId: 1})
	testChunk.Add(Test{Id: 3, A: "Z", GasId: 2})
	testChunk.Add(Test{Id: 4, A: "W", GasId: 3})

	testChunk.GetChunkByHash(4).DeleteBy(func(t *Test) bool { return t.Id == 4 })

	fmt.Printf("%v\n", testChunk.GetChunkByHash(4).List)
	// fmt.Printf("%v\n", testChunk.FindManyByIndex("GasId", 3)[0])
}

func TestIndexReplaceProblem(t *testing.T) {
	testChunk := cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas", IndexList: []string{"GasId"}}
	testChunk.Init()

	testChunk.Add(Test{Id: 1, A: "X", GasId: 1})
	_, ok := testChunk.FindByIndex("GasId", 1)
	if !ok {
		t.Errorf("Index not working")
	}
	testChunk.Replace(Test{Id: 1, A: "Z", GasId: 1})
	l := testChunk.FindManyByIndex("GasId", 1)

	if len(l) != 1 {
		t.Errorf("Index not working")
	}

	if testChunk.TotalElements() != 1 {
		t.Errorf("Index not working")
	}

	// Check
	a, ok := testChunk.Find(func(x *Test) bool { return x.Id == 1 })
	if a.A != "Z" {
		t.Errorf("Index not working")
	}

	// Check
	a, ok = testChunk.FindByIndex("GasId", 1)
	if a.A != "Z" {
		t.Errorf("Index not working")
	}
}

func BenchmarkFind(b *testing.B) {
	testChunk := cdb.ChunkMaster[Test]{Size: 20, Name: "../a/gas"}
	testChunk.Init()

	for i := 0; i < 100000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}

	for i := 0; i < b.N; i++ {
		testChunk.GetChunkByHash(i + 1).Find(func(t *Test) bool {
			return t.Id == i
		})
	}
}

func BenchmarkFindByIndex(b *testing.B) {
	testChunk := cdb.ChunkMaster[Test]{Size: 20, Name: "../a/gas", IndexList: []string{"Id"}}
	testChunk.Init()

	for i := 0; i < 100000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}

	for i := 0; i < b.N; i++ {
		testChunk.FindByIndex("Id", i)
	}
}

func BenchmarkFindMany(b *testing.B) {
	testChunk := cdb.ChunkMaster[Test]{Size: 20, Name: "../a/gas"}
	testChunk.Init()

	for i := 0; i < 100000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}

	for i := 0; i < b.N; i++ {
		testChunk.FindMany(func(t *Test) bool {
			return t.Id == i
		})
	}
}

func BenchmarkFindManyByIndex(b *testing.B) {
	testChunk := cdb.ChunkMaster[Test]{Size: 20, Name: "../a/gas", IndexList: []string{"Id"}}
	testChunk.Init()

	for i := 0; i < 100000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}

	for i := 0; i < b.N; i++ {
		testChunk.FindManyByIndex("Id", i)
	}
}

func BenchmarkDelete(b *testing.B) {
	testChunk := cdb.ChunkMaster[Test]{Size: 20, Name: "../a/gas"}
	testChunk.Init()

	for i := 0; i < 100000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}

	for i := 0; i < b.N; i++ {
		testChunk.GetChunkByHash(i + 1).DeleteBy(func(t *Test) bool {
			return t.Id == i
		})
	}
}

/*func BenchmarkFastAdd(b *testing.B) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.Init()

	for i := 0; i < b.N; i++ {
		testChunk.Add(Test{Id: i})
	}
}*/
