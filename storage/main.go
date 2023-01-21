package storage

type IStorage[T comparable] interface {
	SaveContent(string, []byte)
	AppendContent(string, []byte)

	ReadInfo(string, func(*T))
	ReadContent(string)
}
