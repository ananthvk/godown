package task

import (
	"io"
	"log/slog"
	"mime"
	"net/http"
	"path"
	"path/filepath"

	"github.com/ananthvk/godown/internal/download/storage"
)

const defaultFileName = "download"

type HTTPDownloadTask struct {
	Url           string
	WriterFactory storage.WriterFactory
}

func (h *HTTPDownloadTask) Execute() {
	slog.Info("starting download", slog.String("url", h.Url))
	resp, err := http.Get(h.Url)
	if err != nil {
		slog.Error("starting download", slog.String("url", h.Url), "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		slog.Info("http request sent", "status", resp.Status)
	} else {
		slog.Error("http request sent", "status", resp.Status)
		return
	}

	fileName := getFileName(resp)

	dest, err := h.WriterFactory.CreateStream(fileName)
	if err != nil {
		slog.Error("failed to create write stream", "url", h.Url, "filename", fileName, "err", err)
		return
	}
	defer dest.Close()

	b, err := io.Copy(dest, resp.Body)
	if err != nil {
		slog.Error("failed to save response", "url", h.Url, "filename", fileName, "err", err)
		return
	}
	slog.Info("finished download", "url", h.Url, "filename", fileName, "bytes", b)
}

func getFileName(resp *http.Response) string {
	fileName := getFileNameFromHeader(resp)
	if fileName == "" {
		fileName = getFileNameFromURL(resp)
		slog.Info("resolved filename from url", "filename", fileName)
	} else {
		slog.Info("resolved filename from http header", "filename", fileName)
	}
	return fileName
}

func getFileNameFromHeader(resp *http.Response) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition == "" {
		return ""
	}
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return ""
	}
	if fileName, ok := params["filename"]; ok {
		return fileName
	}
	return ""
}

func getFileNameFromURL(resp *http.Response) string {
	u := resp.Request.URL
	fileName := path.Base(u.Path)
	if fileName == "" || fileName == "." || fileName == "/" || fileName == ".." {
		fileName = defaultFileName
	}

	// If the file name does not have an extension, try to get the extension from the Content-Type header
	if filepath.Ext(fileName) == "" {
		fileName += getFileExt(resp)
	}

	return fileName
}

func getFileExt(resp *http.Response) string {
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return ""
	}
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return ""
	}
	extensions, err := mime.ExtensionsByType(mediatype)
	if err == nil && len(extensions) > 0 {
		return extensions[0]
	}
	return ""
}
