package storage

import "io"

type WriterFactory interface {
	CreateStream(fileName string) (string, io.WriteCloser, error)
}
