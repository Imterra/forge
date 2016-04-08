package objectstorage

import (
	"os"
	"path/filepath"
)

type Storage interface {
	HasObject(filename string, checksum [64]byte) bool
	GetFilePath(filename string, checksum [64]byte) string
	//	GetObject(filename string, checksum [64]byte) string
	//	StoreObject(filename string, data []byte) ([64]byte, error)
}

type FileStorage struct {
	Root string
}

func (s *FileStorage) HasObject(filename string, checksum [64]byte) bool {

	path := filepath.Join(s.Root, filename, GetHashString(checksum))
	_, err := os.Open(path)
	return err == nil
}

func (s *FileStorage) GetFilePath(filename string, checksum [64]byte) string {
	return filepath.Join(s.Root, filename, GetHashString(checksum))
}
