package cmhp_cdb

import (
	"github.com/maldan/go-cmhp/cmhp_slice"
)

func (m *ChunkMaster[T]) GetChunkByHash(valueId any) *Chunk[T] {
	hash := m.Hash(valueId, m.Size)
	return &m.ChunkList[hash]
}

// Find value in chunk [valueId] by [cond]
/*func (m *ChunkMaster[T]) Find(valueId any, cond func(v T) bool) (T, bool) {
	hash := m.Hash(valueId, m.Size)
	m.ChunkList[hash].RLock()
	defer m.ChunkList[hash].RUnlock()

	for _, item := range m.ChunkList[hash].List {
		if cond(item) {
			return item, true
		}
	}
	return *new(T), false
}*/

/*// Contains value in chunk [valueId] by [cond]
func (m *ChunkMaster[T]) Contains(valueId any, cond func(v T) bool) bool {
	hash := m.Hash(valueId, m.Size)
	m.ChunkList[hash].RLock()
	defer m.ChunkList[hash].RUnlock()

	for _, item := range m.ChunkList[hash].List {
		if cond(item) {
			return true
		}
	}
	return false
}*/

// AddOrReplace value to chunk [toHash] and save it
func (m *ChunkMaster[T]) AddOrReplace(v T, where func(v T) bool) {
	if !m.Replace(v, where) {
		m.Add(v)
	}
}

// Add value to chunk [toHash]
func (m *ChunkMaster[T]) Add(v T) {
	// Hash
	hash := m.Hash(v.GetId(), m.Size)

	// Add
	m.ChunkList[hash].Lock()
	m.ChunkList[hash].List = append(m.ChunkList[hash].List, v)
	m.ChunkList[hash].IsChanged = true
	m.ChunkList[hash].Unlock()

	// Add to index map
	m.ChunkList[hash].AddToIndex(&v)
}

// Replace value in chunk [toHash] by condition [where]
func (m *ChunkMaster[T]) Replace(val T, where func(v T) bool) bool {
	// Hash
	hash := m.Hash(val.GetId(), m.Size)

	// Lock
	m.ChunkList[hash].Lock()
	defer m.ChunkList[hash].Unlock()

	// Change
	for i := 0; i < len(m.ChunkList[hash].List); i++ {
		if where(m.ChunkList[hash].List[i]) {
			m.ChunkList[hash].List[i] = val
			m.ChunkList[hash].IsChanged = true
			return true
		}
	}
	return false
}

// Delete values in chunk [toHash] by condition [where]
/*func (m *ChunkMaster[T]) Delete(valueId any, where func(v T) bool) {
	hash := m.Hash(valueId, m.Size)
	m.ChunkList[hash].Lock()
	defer m.ChunkList[hash].Unlock()

	// Filter values
	lenWas := len(m.ChunkList[hash].List)
	m.ChunkList[hash].List = cmhp_slice.Filter(m.ChunkList[hash].List, func(i T) bool {
		return !where(i)
	})

	// Elements was deletes
	if lenWas != len(m.ChunkList[hash].List) {
		m.ChunkList[hash].IsChanged = true
		m.ChunkList[hash].SaveWithoutLock()
	}
}*/

// DeleteInAll values in all chunks by condition [where]
func (m *ChunkMaster[T]) DeleteInAll(where func(v T) bool) {
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].Lock()

		// Filter values
		lenWas := len(m.ChunkList[i].List)
		m.ChunkList[i].List = cmhp_slice.Filter(m.ChunkList[i].List, func(i T) bool {
			return !where(i)
		})

		// Elements was deletes
		if lenWas != len(m.ChunkList[i].List) {
			m.ChunkList[i].IsChanged = true
		}

		m.ChunkList[i].Unlock()
	}
}

func (m *ChunkMaster[T]) FindInAll(fn func(v T) bool) (T, bool) {
	out := *new(T)
	isFound := false
	m.ForEach(func(item T) bool {
		if fn(item) {
			out = item
			isFound = true
			return false
		}
		return true
	})
	return out, isFound
}

func (m *ChunkMaster[T]) ContainsInAll(fn func(v T) bool) bool {
	isFound := false
	m.ForEach(func(item T) bool {
		if fn(item) {
			isFound = true
			return false
		}
		return true
	})
	return isFound
}

func (m *ChunkMaster[T]) FindMany(fn func(v T) bool) []T {
	out := make([]T, 0)

	m.ForEach(func(item T) bool {
		if fn(item) {
			out = append(out, item)
		}
		return true
	})
	return out
}

// All Copy all values from list
func (m *ChunkMaster[T]) All() []T {
	out := make([]T, 0)

	m.ForEach(func(item T) bool {
		out = append(out, item)
		return true
	})
	return out
}

func (m *ChunkMaster[T]) FindByIndex(indexName string, indexValue any) (T, bool) {
	for i := 0; i < len(m.ChunkList); i++ {
		v, ok := m.ChunkList[i].FindByIndex(indexName, indexValue)
		if ok {
			return v, ok
		}
	}
	return *new(T), false
}

func (m *ChunkMaster[T]) FindManyByIndex(indexName string, indexValue any) []T {
	out := make([]T, 0)

	for i := 0; i < len(m.ChunkList); i++ {
		l := m.ChunkList[i].FindManyByIndex(indexName, indexValue)
		out = append(out, l...)
	}
	return out
}
