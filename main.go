package cdb

import (
	"fmt"
	"github.com/maldan/go-cdb/chunk"
	"github.com/maldan/go-cdb/core"
	"github.com/maldan/go-cdb/util"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type idAble interface {
	comparable
	GetId() int
}

type ChunkMaster[T idAble] struct {
	sync.Mutex
	Name          string
	Size          int
	AutoIncrement int
	ChunkList     []chunk.Chunk[T]
	IndexList     []string
	ShowLogs      bool
}

type ChunkMasterAutoIncrement struct {
	Counter int `json:"counter"`
}

func (m *ChunkMaster[T]) Init() *ChunkMaster[T] {
	m.Lock()
	defer m.Unlock()

	if m.Name == "" {
		panic("Chunk name not specified")
	}

	// Inner index for GetId()
	m.IndexList = append(m.IndexList, core.SystemIdField)

	// Read chunk info
	info, err := util.ReadJson[ChunkMasterAutoIncrement](m.Name + "/counter.json")
	if err == nil {
		m.AutoIncrement = info.Counter
	}

	// Init chunks
	t := time.Now()
	m.ChunkList = make([]chunk.Chunk[T], m.Size)
	loadTotal := 0
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].IndexList = m.IndexList
		m.ChunkList[i].Name = m.Name
		m.ChunkList[i].Id = i
		loadTotal += m.ChunkList[i].Load()
	}

	if m.ShowLogs {
		name := m.Name
		wd, _ := os.Getwd()
		name, _ = filepath.Rel(wd, m.Name)
		fmt.Printf("Load chunk [%v] - %v total | %v\n", name, loadTotal, time.Since(t))
	}

	return m
}

func (m *ChunkMaster[T]) EnableAutoSave() *ChunkMaster[T] {
	go func() {
		for {
			m.Save()
			time.Sleep(time.Second)
		}
	}()
	return m
}

func (m *ChunkMaster[T]) Save() {
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].Save()
	}
}

func (m *ChunkMaster[T]) GetChunkByHash(valueId any) *chunk.Chunk[T] {
	hash := m.Hash(valueId, m.Size)
	return &m.ChunkList[hash]
}

// GenerateId generate thread safe autoincrement id
func (m *ChunkMaster[T]) GenerateId() int {
	out := 0
	m.Lock()
	m.AutoIncrement += 1
	out = m.AutoIncrement
	info := ChunkMasterAutoIncrement{m.AutoIncrement}
	util.WriteJson(m.Name+"/counter.json", &info)
	m.Unlock()

	return out
}

func (m *ChunkMaster[T]) TotalElements() int {
	count := 0
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].RLock()
		count += len(m.ChunkList[i].List)
		m.ChunkList[i].RUnlock()
	}
	return count
}

// All Copy all values from list
/*func (m *ChunkMaster[T]) All() []T {
	out := make([]T, 0)

	m.ForEach(func(item T) bool {
		out = append(out, item)
		return true
	})
	return out
}*/

func (m *ChunkMaster[T]) Hash(x any, max int) int {
	switch x.(type) {
	case string:
		hash := 0
		str := x.(string)
		if str == "" {
			panic("empty string hash")
		}
		for i := 0; i < len(str); i++ {
			hash += int(str[i])
		}
		return hash % max
	case int:
		if x.(int) == 0 {
			panic("empty int hash")
		}
		return x.(int) % max
	}
	panic("unsupported hash type")
}
