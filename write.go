package cdb

// AddOrReplace value to chunk [toHash] and save it
func (m *ChunkMaster[T]) AddOrReplace(v T) {
	if !m.Replace(v) {
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

func (m *ChunkMaster[T]) Replace(val T) bool {
	// Hash
	hash := m.Hash(val.GetId(), m.Size)

	// Lock and update
	m.ChunkList[hash].Lock()
	isChanged := false
	id := val.GetId()
	replacedList := make([]*T, 0)
	for i := 0; i < len(m.ChunkList[hash].List); i++ {
		if m.ChunkList[hash].List[i].GetId() == id {
			old := m.ChunkList[hash].List[i]
			replacedList = append(replacedList, &old)
			m.ChunkList[hash].List[i] = val
			m.ChunkList[hash].IsChanged = true
			isChanged = true
			break
		}
	}
	m.ChunkList[hash].Unlock()

	// Delete old from index
	m.ChunkList[hash].DeleteFromIndex(replacedList)

	// Add new to index
	m.ChunkList[hash].AddToIndex(&val)

	return isChanged
}

/*// Replace value in chunk [toHash] by condition [where]
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
}*/

/*// DeleteInAll values in all chunks by condition [where]
func (m *ChunkMaster[T]) DeleteInAll(where func(v *T) bool) {
	for i := 0; i < m.Size; i++ {
		m.ChunkList[i].Lock()

		// Filter values
		lenWas := len(m.ChunkList[i].List)
		m.ChunkList[i].List = util.FilterSlice(m.ChunkList[i].List, func(i *T) bool {
			return !where(i)
		})

		// Elements was deletes
		if lenWas != len(m.ChunkList[i].List) {
			m.ChunkList[i].IsChanged = true
		}

		m.ChunkList[i].Unlock()
	}
}
*/
