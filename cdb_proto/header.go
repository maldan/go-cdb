package cdb_proto

/**
Header struct

[P R O T 1 2 3 4] - file id
[1] - version

[0 0 0 0 0 0 0 0] - auto increment
[0 0 0 0 0 0 0 0] - total records

[0] - amount of parts for parallel search
[0 0 0 0 0 0 0 0] * 64 - offsets for parallel search

[...] - struct info
*/

/**
Struct info
[0] - total fields
[
	[0] - id/index
	[0] - type
	[0 0 0 0] - max length/capacity
	[0] - name length
	[...] - name
] * totalFields
*/

func (d *DataTable[T]) emptyHeader() []byte {
	bytes := make([]byte, headerSize)

	// File id
	copy(bytes, "PROT1234")

	// Version
	bytes[8] = 1

	return bytes
}
