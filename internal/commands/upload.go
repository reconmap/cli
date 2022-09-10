package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/reconmap/cli/internal/terminal"
	"github.com/reconmap/shared-lib/pkg/api"
	"github.com/reconmap/shared-lib/pkg/configuration"
)

func UploadResults(command *api.Command, taskId int) error {
	return UploadCommandOutputUsingFileName(command, taskId)
}

func UploadCommandOutputUsingFileName(command *api.Command, taskId int) error {
	if len(strings.TrimSpace(command.OutputFileName)) == 0 {
		return errors.New("The command has not defined an output filename. Nothing has been uploaded to the server.")
	}

	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}
	var remoteURL string = config.ApiUrl + "/commands/outputs"

	var client *http.Client = &http.Client{}
	err = Upload(client, remoteURL, command.OutputFileName, command.ID, taskId)
	return err
}

func Upload(client *http.Client, url string, outputFileName string, commandId int, taskId int) (err error) {

	if _, err := os.Stat(outputFileName); os.IsNotExist(err) {
		return fmt.Errorf("Output file '%s' could not be found", outputFileName)
	}

	file, err := os.Open(filepath.Clean(outputFileName))
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("resultFile", filepath.Base(outputFileName))
	_, err = io.Copy(part, file)

	if err = writer.WriteField("commandId", strconv.Itoa(commandId)); err != nil {
		return
	}
	if taskId != 0 {
		if err = writer.WriteField("taskId", strconv.Itoa(taskId)); err != nil {
			return
		}
	}

	if err = writer.Close(); err != nil {
		return
	}

	req, err := api.NewRmapRequest("POST", url, body)
	if err != nil {
		return
	}

	err = api.AddBearerToken(req)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	terminal.PrintYellowDot()
	fmt.Printf(" Uploading command output '%s' to the server.\n", outputFileName)
	res, err := client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode == http.StatusUnauthorized {
		err = fmt.Errorf("your session has expired. Please login again")
	}
	terminal.PrintGreenTick()
	fmt.Printf(" Done\n")

	return
}
