package download

import (
	"log/slog"
	"net/url"
	"sync"

	"github.com/ananthvk/godown/internal/download/storage"
	"github.com/ananthvk/godown/internal/download/task"
)

type Downloader struct {
	writerFactory storage.WriterFactory
	wg            sync.WaitGroup
}

func NewDownloader() *Downloader {
	downloader := Downloader{}
	downloader.writerFactory = &storage.FSWriterFactory{}
	return &downloader
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
		t = &task.HTTPDownloadTask{Url: urlString, WriterFactory: d.writerFactory}
	default:
		slog.Error("unsupported url scheme", "scheme", url.Scheme)
		return
	}

	d.wg.Add(1)
	go func() {
		t.Execute()
		d.wg.Done()
	}()
}

func (d *Downloader) Wait() {
	d.wg.Wait()
}
