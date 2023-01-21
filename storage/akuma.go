package storage

import "github.com/maldan/go-cdb/pack"

type AkumaStorage[T comparable] struct {
}

func (r AkumaStorage[T]) SaveContent(s string, bytes []byte) {
	//TODO implement me
	panic("implement me")
}

func (r AkumaStorage[T]) AppendContent(s string, bytes []byte) {
	//TODO implement me
	panic("implement me")
}

func (r AkumaStorage[T]) ReadContent(s string) {
	//TODO implement me
	panic("implement me")
}

func (r AkumaStorage[T]) ReadInfo(path string, populate func(*T)) {
	t := []byte{1}
	pack.Unpack[T](&t)
	populate(new(T))
}
