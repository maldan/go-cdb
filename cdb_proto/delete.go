package cdb_proto

func (d *DataTable[T]) Delete(fieldList []string, where func(*T) bool) {
	/*mapper := ValueMapper[T]{}
	mapper.Map(d.structInfo, fieldList)

	d.ForEach(func(offset int) bool {
		// size, offTable := pack.ReadHeader2(d.mem, offset)
		offTable := pack.ReadOffsetTable(d.mem[offset:])

		mapper.Fill(offset, d.mem, offTable)

		if where(&mapper.Container) {
			// Read flags
			b := []byte{0}
			d.file.ReadAt(b, int64(offset+core.RecordStart+core.RecordSize))
			fmt.Printf("FB: %v\n", b[0])
			b[0] |= core.MaskDeleted
			fmt.Printf("FAPL: %v\n", b[0])

			// Write back
			d.file.WriteAt(b, int64(offset+core.RecordStart+core.RecordSize))

			// Check
			fmt.Printf("FA: %v\n", d.mem[offset+core.RecordStart+core.RecordSize])
		}

		return true
	})*/
}
