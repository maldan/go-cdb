package core

const HeaderSize = 2048

// RecordStart is 1 bytes 0x12 of header for each record
const RecordStart = 1

// RecordSize is size of each record
const RecordSize = 4

// RecordFlags amount of fields, is deleted
const RecordFlags = 1

// RecordEnd is 1 bytes 0x34 of header for each record
const RecordEnd = 1

const RecordStartMark = 0x12
const RecordEndMark = 0x34

const MaskDeleted = 0b1000_0000
const MaskTotalFields = 0b0011_1111

type StructInfo struct {
	FieldCount    int
	FieldNameToId map[string]int
	FieldType     []int
	FieldName     []string
	FieldOffset   []uintptr
}
