package storage

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type FSWriterFactory struct {
	BasePath string
	mu       sync.Mutex
}

func (f *FSWriterFactory) CreateStream(fileName string) (string, io.WriteCloser, error) {
	if len(f.BasePath) != 0 {
		err := os.MkdirAll(f.BasePath, 0755)
		if err != nil {
			return fileName, nil, err
		}
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	filePath := path.Join(f.BasePath, fileName)
	fileNameNew := ""

	ctr := 1
	// While the file exists, check for new names

	for {
		exists, err := doesFileExist(filePath)
		if err != nil {
			return "", nil, err
		}
		if !exists {
			fileName = fileNameNew
			break
		}
		// TODO: Optimization: Store the filename along with the counter in a map so that the next time the free file name can be
		// found easily
		slog.Info("file exists", "path", filePath)
		ext := filepath.Ext(fileName)
		fileNameWithoutExt := strings.TrimSuffix(fileName, ext)
		fileNameNew = fmt.Sprintf("%s (%d)%s", fileNameWithoutExt, ctr, ext)
		filePath = path.Join(f.BasePath, fileNameNew)
		ctr++
	}

	file, err := os.Create(filePath)
	return fileName, file, err
}

func doesFileExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
