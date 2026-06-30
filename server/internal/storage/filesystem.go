package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type FSStorage struct {
	basePath string
}

func NewFSStorage(path string) *FSStorage {
	os.MkdirAll(path, 0755)
	return &FSStorage{basePath: path}
}

func (s *FSStorage) Save(ctx context.Context, id string, r io.Reader) (string, error) {
	filePath := filepath.Join(s.basePath, id+".jpg")

	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	return filePath, err
}

func (s *FSStorage) Delete(ctx context.Context, path string) error {
	return os.Remove(path)
}
