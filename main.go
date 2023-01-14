package cmhp_cdb

import (
	"fmt"
	"github.com/maldan/go-cmhp/cmhp_file"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type idAble interface {
	GetId() any
}

type ChunkMaster[T idAble] struct {
	sync.Mutex
	Name          string
	Size          int
	AutoIncrement int
	ChunkList     []Chunk[T]
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

	// Read chunk info
	info, err := cmhp_file.ReadGenericJSON[ChunkMasterAutoIncrement](m.Name + "/counter.json")
	if err == nil {
		m.AutoIncrement = info.Counter
	}

	// Init chunks
	t := time.Now()
	m.ChunkList = make([]Chunk[T], m.Size)
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

/*func (m *ChunkMaster[T]) SaveChanged() {
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].SaveIfChanged()
	}
}*/

/*func (m *ChunkMaster[T]) SaveChunkByHash(toHash any) {
	hash := m.Hash(toHash, m.Size)
	m.ChunkList[hash].Save()
}*/

// ForEach go over each value in all chunks
func (m *ChunkMaster[T]) ForEach(fn func(item T) bool) {
	status := true
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].RLock()
		for j := 0; j < len(m.ChunkList[i].List); j++ {
			if !fn(m.ChunkList[i].List[j]) {
				status = false
				break
			}
		}
		m.ChunkList[i].RUnlock()
		if !status {
			break
		}
	}
}

// GenerateId generate thread safe autoincrement id
func (m *ChunkMaster[T]) GenerateId() int {
	out := 0
	m.Lock()
	m.AutoIncrement += 1
	out = m.AutoIncrement
	info := ChunkMasterAutoIncrement{m.AutoIncrement}
	cmhp_file.Write(m.Name+"/counter.json", &info)
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

/*func (m *ChunkMaster[T]) StrHash(str string, max int) int {
	hash := 0
	for i := 0; i < len(str); i++ {
		hash += int(str[i])
	}
	return hash % max
}

func (m *ChunkMaster[T]) IntHash(num int, max int) int {
	return num % max
}*/

/*func (m *ChunkMaster[T]) BuildIndexMap() {
	m.IndexStorage = make(map[string][]*T)

	// Build index
	for _, index := range m.IndexList {
		for i := 0; i < m.Size; i++ {
			for j := 0; j < len(m.ChunkList[i].List); j++ {
				m.AddIndex(index, &m.ChunkList[i].List[j])
			}
		}
	}
}

func (m *ChunkMaster[T]) AddIndex(index string, ref *T) {
	m.Lock()
	defer m.Unlock()

	f := reflect.ValueOf(ref).Elem().FieldByName(index)
	mapIndex := reflect.ValueOf(f).Interface()
	strIndex := fmt.Sprintf("%v:%v", index, mapIndex)
	m.IndexStorage[strIndex] = append(m.IndexStorage[strIndex], ref)
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
