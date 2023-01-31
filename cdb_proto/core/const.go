package core

const TBool = 1

const T8 = 2
const T16 = 3
const T32 = 4
const T64 = 5

const TF32 = 6
const TF64 = 7

const TString = 8

const TStruct = 9

const TSlice = 10

const HeaderSize = 1024

const TokenString = 0
const TokenOp = 1
const TokenIdentifier = 2
const TokenNumber = 3

// RecordStart is 1 bytes 0x12 of header for each record
const RecordStart = 1

// RecordSize is size of each record
const RecordSize = 4

// RecordFlags amount of fields, is deleted
const RecordFlags = 1

// RecordEnd is 1 bytes 0x34 of header for each record
const RecordEnd = 1

// RecordLenOff is 8 bytes size for offset and length for each field of struct
const RecordLenOff = 8

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
