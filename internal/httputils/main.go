package httputils

import (
	"io"
	"io/ioutil"
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

func AddBearerToken(req *http.Request) error {
	b, err := ioutil.ReadFile("rmap-session")
	if err != nil {
		return err
	}
	session := string(b)

	var bearer = "Bearer " + session
	req.Header.Add("Authorization", bearer)

	return nil
}
