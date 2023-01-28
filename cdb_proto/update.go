package cdb_proto

func Update(m []byte, index int, field string, v string) {
	/*offset := int(indexMap[index])
	typeOf := reflect.TypeOf(Test{})

	fieldOffsetIndex := 0
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Name == field {
			fieldOffsetIndex = i * 8
			break
		}
	}

	_, _, startData := readHeader(m, offset, fieldOffsetIndex)
	newLen := len(v)
	cmhp_byte.From32ToBuffer(&newLen, &m, offset+_hSize+_hTotal+fieldOffsetIndex+4)
	copy(m[startData:], v)*/
}
