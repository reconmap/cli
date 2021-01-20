package commands

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/reconmap/cli/internal/httputils"
)

func UploadResults() error {
	var client *http.Client = &http.Client{}
	var remoteURL string = "https://api.reconmap.org/tasks/results"
	remoteURL = "http://localhost:8080/tasks/results"

	err := Upload(client, remoteURL)
	return err
}

func Upload(client *http.Client, url string) (err error) {
	file, err := os.Open("report-20405-reconmap.org.txt")
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("resultFile", filepath.Base("report-20405-reconmap.org.txt"))
	_, err = io.Copy(part, file)

	writer.WriteField("taskId", "4")

	writer.Close()

	req, err := httputils.NewRmapRequest("POST", url, body)
	if err != nil {
		return
	}

	httputils.AddBearerToken(req)

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode == http.StatusUnauthorized {
		err = fmt.Errorf("your session has expired. Please login again")
	}
	return
}
