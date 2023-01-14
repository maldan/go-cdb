package chunk

import (
	"fmt"
	"github.com/maldan/go-cdb/util"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

type idAble interface {
	GetId() any
}

type Chunk[T idAble] struct {
	sync.RWMutex
	IsLoad       bool
	IsInit       bool
	IsChanged    bool
	List         []T
	Name         string
	Id           int
	IndexList    []string
	indexStorage map[string][]*T
	ShowLogs     bool
}

func (c *Chunk[T]) BuildIndexMap() {
	c.indexStorage = make(map[string][]*T)

	// Build index
	for i := 0; i < len(c.List); i++ {
		c.AddToIndex(&c.List[i])
	}
}

func (c *Chunk[T]) AddToIndex(ref *T) {
	c.Lock()
	defer c.Unlock()

	for _, index := range c.IndexList {
		f := reflect.ValueOf(ref).Elem().FieldByName(index)
		mapIndex := reflect.ValueOf(f).Interface()
		strIndex := fmt.Sprintf("%s:%v", index, mapIndex)
		c.indexStorage[strIndex] = append(c.indexStorage[strIndex], ref)
	}
}

func (c *Chunk[T]) DeleteFromIndex(ref *T, where func(v *T) bool) {
	c.Lock()
	defer c.Unlock()

	for _, index := range c.IndexList {
		f := reflect.ValueOf(ref).Elem().FieldByName(index)
		mapIndex := reflect.ValueOf(f).Interface()
		strIndex := fmt.Sprintf("%s:%v", index, mapIndex)

		lenWas := len(c.indexStorage[strIndex])
		newList := make([]*T, 0, lenWas)
		for i := 0; i < lenWas; i++ {
			if !where(&c.List[i]) {
				newList = append(newList, c.indexStorage[strIndex][i])
			}
		}
		c.indexStorage[strIndex] = newList
	}
}

func (c *Chunk[T]) Save() {
	c.Lock()
	defer c.Unlock()

	if !c.IsChanged {
		return
	}

	// Write to disk
	t := time.Now()
	err := util.WriteJson(c.Name+"/chunk_"+fmt.Sprintf("%v", c.Id)+".json.tmp", &c.List)
	if err != nil {
		panic(err)
	}

	// Delete old
	util.DeleteFile(c.Name + "/chunk_" + fmt.Sprintf("%v", c.Id) + ".json")

	// Replace
	err = util.Rename(c.Name+"/chunk_"+fmt.Sprintf("%v", c.Id)+".json.tmp", c.Name+"/chunk_"+fmt.Sprintf("%v", c.Id)+".json")
	if err != nil {
		panic(err)
	}

	if c.ShowLogs {
		name := c.Name
		wd, _ := os.Getwd()
		name, _ = filepath.Rel(wd, c.Name)
		fmt.Printf("Save chunk [%v:%v] - %v records | %v\n", name, c.Id, len(c.List), time.Since(t))
	}
	c.IsChanged = false
}

func (c *Chunk[T]) Load() int {
	c.Lock()
	defer c.BuildIndexMap()
	defer c.Unlock()

	t := time.Now()
	chunk, err := util.ReadJson[[]T](c.Name + "/chunk_" + fmt.Sprintf("%v", c.Id) + ".json")
	if err != nil {
		c.List = make([]T, 0)
		c.IsInit = true
		if c.ShowLogs {
			fmt.Printf("Load chunk [%v:%v] - empty\n", c.Name, c.Id)
		}
		return 0
	}

	c.List = chunk
	c.IsLoad = true
	c.IsInit = true
	if c.ShowLogs {
		fmt.Printf("Load chunk [%v:%v] - %v records | %v\n", c.Name, c.Id, len(chunk), time.Since(t))
	}
	return len(chunk)
}

// DeleteBy values in chunk by condition [where]
func (c *Chunk[T]) DeleteBy(where func(v *T) bool) {
	c.Lock()

	// Filter values
	lenWas := len(c.List)
	newList := make([]T, 0, lenWas)
	deletedList := make([]T, 0, 10)
	for i := 0; i < lenWas; i++ {
		if !where(&c.List[i]) {
			newList = append(newList, c.List[i])
		} else {
			deletedList = append(deletedList, c.List[i])
		}
	}
	c.List = newList

	// Elements was deletes
	if lenWas != len(c.List) {
		c.IsChanged = true
	}

	c.Unlock()
}

func (c *Chunk[T]) Delete(v T) {
	id := v.GetId()
	c.DeleteBy(func(t *T) bool { return (*t).GetId() == id })
}
