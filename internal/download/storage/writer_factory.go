package storage

import "io"

type WriterFactory interface {
	CreateStream(fileName string) (io.WriteCloser, error)
}
