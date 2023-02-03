package cdb_proto

type IData interface {
	Marshall(any)
	Unmarshall([]byte) any
}
