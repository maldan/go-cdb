package util

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteJson(path string, v any) error {
	// Create path for file
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}

	// Write as json
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0777)

	return nil
}

func ReadJson[T any](path string) (T, error) {
	s := new(T)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return *s, err
	}
	err = json.Unmarshal(data, s)

	return *s, err
}

func Info(path string) (fs.FileInfo, error) {
	stat, err := os.Stat(path)
	return stat, err
}

func DeleteFile(path string) error {
	info, err := Info(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func Rename(src string, dst string) error {
	return os.Rename(src, dst)
}
