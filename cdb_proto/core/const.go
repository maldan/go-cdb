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

// RecordStart is 2 bytes 0x1234 of header for each record
const RecordStart = 2
const RecordSize = 4
const RecordFlags = 1
const RecordEnd = 2

// RecordLenOff is 8 bytes size for offset and length for each field of struct
const RecordLenOff = 8

const MaskDeleted = 0b1000_0000
const MaskTotalFields = 0b0011_1111

type StructInfo struct {
	FieldCount    int
	FieldNameToId map[string]int
	FieldType     []int
	FieldName     []string
	FieldOffset   []uintptr
}
