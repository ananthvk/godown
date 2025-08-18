package storage

import (
	"io"
	"os"
	"path"
)

type FSWriterFactory struct {
	BasePath string
}

func (f *FSWriterFactory) CreateStream(fileName string) (io.WriteCloser, error) {
	if len(f.BasePath) != 0 {
		err := os.MkdirAll(f.BasePath, 0755)
		if err != nil {
			return nil, err
		}
	}
	return os.Create(path.Join(f.BasePath, fileName))
}
