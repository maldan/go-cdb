package engine

func (d *DataEngine[T]) Init() *DataEngine[T] {
	d.indexList = make(map[string][]uint32, 0)

	return d
}

func (d *DataEngine[T]) Flush() {
	d.storage.Flush()
}

func (d *DataEngine[T]) Load(populate func(T)) {
	d.storage.Load(populate)
}
