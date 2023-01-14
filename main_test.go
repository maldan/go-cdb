package cmhp_cdb_test

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_cdb"
	"testing"
)

type Test struct {
	Id    int    `json:"id"`
	A     string `json:"a"`
	GasId int    `json:"gasId"`
}

func (t Test) GetId() any {
	return t.Id
}

/*func TestM(t *testing.T) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.LoadAll()

	// Add
	testChunk.Add(Test{Id: 1}, 1)
	testChunk.Add(Test{Id: 2}, 2)
	testChunk.Add(Test{Id: 3}, 3)
	if testChunk.TotalElements() != 3 {
		t.Errorf("Fuck you")
	}

	// Delete
	testChunk.Delete(1, func(t Test) bool {
		return t.Id == 1
	})
	if testChunk.TotalElements() != 2 {
		t.Errorf("Fuck you")
	}

	// Test if changed
	if !testChunk.ChunkList[9].IsChanged {
		t.Errorf("Chunk must be changed")
	}
	if testChunk.ChunkList[3].IsChanged {
		t.Errorf("Chunk must not be changed")
	}

	// All
	if len(testChunk.All()) != 2 {
		t.Errorf("Fuck you")
	}
}*/

func TestAdd(t *testing.T) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
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
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.Init()

	// Add
	for i := 0; i < 10000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}

	// Delete
	testChunk.GetChunkByHash(1).Delete(func(t Test) bool {
		return t.Id == 1
	})
	if testChunk.TotalElements() != 10000-1 {
		t.Errorf("Delete not working")
	}

	// Delete in all
	testChunk.GetChunkByHash(555).Delete(func(t Test) bool {
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
	_, ok = testChunk.FindInAll(func(t Test) bool {
		return t.Id == 1
	})
	if ok {
		t.Errorf("Delete not working")
	}
	_, ok = testChunk.FindInAll(func(t Test) bool {
		return t.Id == 555
	})
	if ok {
		t.Errorf("Delete not working")
	}
}

func TestUpdate(t *testing.T) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.Init()

	// Add
	testChunk.Add(Test{Id: 1})
	testChunk.Add(Test{Id: 2})
	testChunk.Add(Test{Id: 3})
	if testChunk.TotalElements() != 3 {
		t.Errorf("Fuck you")
	}

	// Update
	testChunk.Replace(Test{Id: 1, A: "gas"}, func(t Test) bool {
		return t.Id == 1
	})
	if testChunk.TotalElements() != 3 {
		t.Errorf("Fuck you")
	}

	// Test if changed
	if testChunk.ChunkList[1].List[0].A != "gas" {
		t.Errorf("Update not working")
	}
}

func TestFind(t *testing.T) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
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
	v, ok = testChunk.FindInAll(func(t Test) bool {
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
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.Init()

	// Add
	for i := 0; i < 3; i++ {
		testChunk.Add(Test{Id: i + 1})
	}
	if testChunk.TotalElements() != 3 {
		t.Errorf("Fuck you")
	}

	// Find in chunk
	list := testChunk.FindMany(func(t Test) bool {
		return t.Id > 2
	})
	if len(list) != 1 {
		t.Errorf(fmt.Sprintf("Fuck you %v", len(list)))
	}
}

func TestIndex(t *testing.T) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas", IndexList: []string{"GasId"}}
	testChunk.Init()

	testChunk.Add(Test{1, "X", 1})
	testChunk.Add(Test{2, "Y", 1})
	testChunk.Add(Test{3, "Z", 2})
	testChunk.Add(Test{4, "W", 3})

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
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas", IndexList: []string{"GasId"}}
	testChunk.Init()

	testChunk.Add(Test{1, "X", 1})
	testChunk.Add(Test{2, "Y", 1})
	testChunk.Add(Test{3, "Z", 2})
	testChunk.Add(Test{4, "W", 3})

	testChunk.GetChunkByHash(4).Delete(func(t Test) bool { return t.Id == 4 })

	fmt.Printf("%v\n", testChunk.GetChunkByHash(4).List)
	fmt.Printf("%v\n", testChunk.FindManyByIndex("GasId", 3)[0])
}

/*
func BenchmarkSmart1(b *testing.B) {
	cm := cmhp_cdb.ChunkMaster[Test]{}
	for i := 0; i < b.N; i++ {
		cm.Hash(i, 20)
	}
}

func BenchmarkSmart2(b *testing.B) {
	cm := cmhp_cdb.ChunkMaster[Test]{}
	for i := 0; i < b.N; i++ {
		cm.Hash("123", 20)
	}
}*/

func BenchmarkFind(b *testing.B) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 20, Name: "../a/gas"}
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
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 20, Name: "../a/gas", IndexList: []string{"Id"}}
	testChunk.Init()

	for i := 0; i < 100000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}

	for i := 0; i < b.N; i++ {
		testChunk.FindByIndex("Id", i)
	}
}

func BenchmarkFindMany(b *testing.B) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 20, Name: "../a/gas"}
	testChunk.Init()

	for i := 0; i < 100000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}

	for i := 0; i < b.N; i++ {
		testChunk.FindMany(func(t Test) bool {
			return t.Id == i
		})
	}
}

func BenchmarkFindManyByIndex(b *testing.B) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 20, Name: "../a/gas", IndexList: []string{"Id"}}
	testChunk.Init()

	for i := 0; i < 100000; i++ {
		testChunk.Add(Test{Id: i + 1})
	}

	for i := 0; i < b.N; i++ {
		testChunk.FindManyByIndex("Id", i)
	}
}

/*func BenchmarkFastAdd(b *testing.B) {
	testChunk := cmhp_cdb.ChunkMaster[Test]{Size: 10, Name: "../a/gas"}
	testChunk.Init()

	for i := 0; i < b.N; i++ {
		testChunk.Add(Test{Id: i})
	}
}*/
