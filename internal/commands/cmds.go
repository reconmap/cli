package commands

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/reconmap/cli/internal/httputils"
)

type Command struct {
	DockerImage   string `json:"docker_image"`
	ContainerArgs string `json:"container_args"`
}

func RunCommand(id int) error {

	var apiUrl string = "http://localhost:8080/commands/" + strconv.Itoa(id)

	client := &http.Client{}
	req, err := httputils.NewRmapRequest("GET", apiUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httputils.AddBearerToken(req)

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		return errors.New("error from server: " + string(body))
	}

	if err != nil {
		return errors.New("Unable to read response from server")
	}

	var command Command

	json.Unmarshal([]byte(body), &command)
	_, err = CreateNewContainer(command)

	return err
}
