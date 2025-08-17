package task

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type HTTPDownloadTask struct {
	Url string
}

func (h *HTTPDownloadTask) Execute() {
	fmt.Printf("Downloading %s\n", h.Url)
	resp, err := http.Get(h.Url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while getting %q: %v\n", h.Url, err)
		return
	}
	defer resp.Body.Close()
	b, err := io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while downloading %q: %v\n", h.Url, err)
		return
	}
	fmt.Printf("\n%d bytes received\n", b)
}
