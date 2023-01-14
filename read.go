package cdb

func (m *ChunkMaster[T]) HasBy(fn func(v *T) bool) bool {
	for i := 0; i < m.Size; i++ {
		_, ok := m.ChunkList[i].Find(fn)
		if ok {
			return true
		}
	}
	return false
}

func (m *ChunkMaster[T]) Find(fn func(v *T) bool) (T, bool) {
	for i := 0; i < m.Size; i++ {
		v, ok := m.ChunkList[i].Find(fn)
		if ok {
			return v, ok
		}
	}
	return *new(T), false
}

func (m *ChunkMaster[T]) FindMany(fn func(v *T) bool) []T {
	out := make([]T, 0)

	for i := 0; i < m.Size; i++ {
		values := m.ChunkList[i].FindMany(fn)
		out = append(out, values...)
	}

	return out
}

func (m *ChunkMaster[T]) FindByIndex(indexName string, indexValue any) (T, bool) {
	for i := 0; i < len(m.ChunkList); i++ {
		v, ok := m.ChunkList[i].FindByIndex(indexName, indexValue)
		if ok {
			//
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
