package httputils

import (
	"io"
	"net/http"
	"runtime"
)

func NewRmapRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Rmap/1.0 ("+runtime.GOOS+")")
	return req, nil
}
