package download

import (
	"fmt"
	"net/url"

	"os"

	"github.com/ananthvk/godown/internal/download/task"
)

type Downloader struct {
}

func (d *Downloader) Download(urlString string) {
	if !IsUrl(urlString) {
		fmt.Fprintf(os.Stderr, "Invalid URL %q\n", urlString)
		return
	}
	url, _ := url.Parse(urlString)
	var t task.Task

	switch url.Scheme {
	case "http", "https":
		t = &task.HTTPDownloadTask{Url: urlString}
	default:
		fmt.Fprintf(os.Stderr, "URL scheme '%s' not supported\n", url.Scheme)
		return
	}
	t.Execute()
}
