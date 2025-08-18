package download

import (
	"log/slog"
	"net/url"

	"github.com/ananthvk/godown/internal/download/storage"
	"github.com/ananthvk/godown/internal/download/task"
)

type Downloader struct {
}

func (d *Downloader) Download(urlString string) {
	if !IsUrl(urlString) {
		slog.Error("invalid url", "url", urlString)
		return
	}
	url, _ := url.Parse(urlString)
	var t task.Task

	switch url.Scheme {
	case "http", "https":
		t = &task.HTTPDownloadTask{Url: urlString, WriterFactory: &storage.FSWriterFactory{}}
	default:
		slog.Error("unsupported url scheme", "scheme", url.Scheme)
		return
	}
	t.Execute()
}
