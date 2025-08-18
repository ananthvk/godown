package download

import (
	"context"
	"log/slog"
	"net/url"
	"sync"

	"github.com/ananthvk/godown/internal/download/storage"
	"github.com/ananthvk/godown/internal/download/task"
)

type Downloader struct {
	writerFactory    storage.WriterFactory
	wg               sync.WaitGroup
	ignoreInvalidURL bool
}

func NewDownloader(basePath string, ignoreInvalidURL bool) *Downloader {
	downloader := Downloader{}
	downloader.writerFactory = &storage.FSWriterFactory{BasePath: basePath}
	downloader.ignoreInvalidURL = ignoreInvalidURL
	return &downloader
}

func (d *Downloader) Download(ctx context.Context, urlString string) {
	if !d.ignoreInvalidURL && !IsUrl(urlString) {
		slog.Error("invalid url", "url", urlString)
		return
	}
	url, _ := url.Parse(urlString)
	var t task.Task

	switch url.Scheme {
	case "http", "https":
		t = &task.HTTPDownloadTask{Url: urlString, WriterFactory: d.writerFactory}
	default:
		if d.ignoreInvalidURL && url.Scheme == "" {
			// If the url is invalid because it lacks a URL scheme, try adding a default http:// scheme
			t = &task.HTTPDownloadTask{Url: "http://" + urlString, WriterFactory: d.writerFactory}
		} else {
			slog.Error("unsupported url scheme", "scheme", url.Scheme)
			return
		}
	}

	d.wg.Add(1)
	go func() {
		t.Execute(ctx)
		d.wg.Done()
	}()
}

func (d *Downloader) Wait() {
	d.wg.Wait()
}
