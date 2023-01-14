package chunk

import "fmt"

// Contains value in chunk by [cond]
func (c *Chunk[T]) Contains(cond func(v *T) bool) bool {
	c.RLock()
	defer c.RUnlock()

	for i := 0; i < len(c.List); i++ {
		if cond(&c.List[i]) {
			return true
		}
	}
	return false
}

func (c *Chunk[T]) FindByIndex(indexName string, indexValue any) (T, bool) {
	c.RLock()
	defer c.RUnlock()

	strIndex := fmt.Sprintf("%s:%v", indexName, indexValue)
	for _, val := range c.indexStorage[strIndex] {
		return *val, true
	}

	return *new(T), false
}

func (c *Chunk[T]) FindManyByIndex(indexName string, indexValue any) []T {
	c.RLock()
	defer c.RUnlock()

	strIndex := fmt.Sprintf("%s:%v", indexName, indexValue)
	out := make([]T, 0)
	for _, val := range c.indexStorage[strIndex] {
		out = append(out, *val)
	}
	return out
}

// Find value in chunk by [cond]
func (c *Chunk[T]) Find(cond func(v *T) bool) (T, bool) {
	c.RLock()
	defer c.RUnlock()

	for i := 0; i < len(c.List); i++ {
		if cond(&c.List[i]) {
			return c.List[i], true
		}
	}

	return *new(T), false
}

// FindMany values in chunk by [cond]
func (c *Chunk[T]) FindMany(cond func(v *T) bool) []T {
	c.RLock()
	defer c.RUnlock()

	out := make([]T, 0)
	for i := 0; i < len(c.List); i++ {
		if cond(&c.List[i]) {
			out = append(out, c.List[i])
		}
	}

	return out
}
