package httputils

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/reconmap/cli/internal/build"
	"github.com/reconmap/cli/internal/configuration"
)

func NewRmapRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Rmap/"+build.BuildVersion+" ("+runtime.GOOS+")")
	return req, nil
}

func AddBearerToken(req *http.Request) error {
	jwtToken, err := ReadSessionToken()
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+jwtToken)

	return nil
}

func ReadSessionToken() (string, error) {
	reconmapConfigDir, err := configuration.GetReconmapConfigDirectory()
	if err != nil {
		return "", err
	}

	var configPath = filepath.Join(reconmapConfigDir, "session-token")

	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func SaveSessionToken(accessToken string) error {
	reconmapConfigDir, err := configuration.GetReconmapConfigDirectory()
	if err != nil {
		return err
	}

	var configPath = filepath.Join(reconmapConfigDir, "session-token")

	err = ioutil.WriteFile(configPath, []byte(accessToken), 0600)
	return err
}

func GetSessionTokenPath() (string, error) {
	reconmapConfigDir, err := configuration.GetReconmapConfigDirectory()
	if err != nil {
		return "", err
	}

	var configPath = filepath.Join(reconmapConfigDir, "session-token")
	return configPath, nil
}
