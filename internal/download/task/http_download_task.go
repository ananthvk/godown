package task

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ananthvk/godown/internal/download/storage"
)

type HTTPDownloadTask struct {
	Url           string
	WriterFactory storage.WriterFactory
}

func (h *HTTPDownloadTask) Execute() {
	fmt.Printf("Downloading %s\n", h.Url)
	resp, err := http.Get(h.Url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while getting %q: %v\n", h.Url, err)
		return
	}
	defer resp.Body.Close()
	v := resp.Header["Content-Disposition"]
	fmt.Println(v[0])

	dest, err := h.WriterFactory.CreateStream("file.pdf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while downloading %q: %v\n", h.Url, err)
		return
	}
	defer dest.Close()

	b, err := io.Copy(dest, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while saving %q: %v\n", h.Url, err)
		return
	}
	fmt.Printf("\n%d bytes received\n", b)
}
