package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/reconmap/cli/internal/configuration"
	"github.com/reconmap/cli/internal/httputils"
)

func GetCommandById(id int) (*Command, error) {
	config, err := configuration.ReadConfig()
	if err != nil {
		return nil, err
	}
	var apiUrl string = config.ApiUrl + "/commands/" + strconv.Itoa(id)

	client := &http.Client{}
	req, err := httputils.NewRmapRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err = httputils.AddBearerToken(req); err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error from server: " + string(response.Status))
	}

	if err != nil {
		return nil, errors.New("Unable to read response from server")
	}

	var command *Command = &Command{}

	if err = json.Unmarshal([]byte(body), command); err != nil {
		return command, err
	}

	return command, nil
}

func GetCommandsByKeywords(keywords string) (*Commands, error) {
	config, err := configuration.ReadConfig()
	if err != nil {
		return nil, err
	}
	var apiUrl string = config.ApiUrl + "/commands?keywords=" + keywords

	client := &http.Client{}
	req, err := httputils.NewRmapRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err = httputils.AddBearerToken(req); err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error from server: " + string(response.Status))
	}

	if err != nil {
		return nil, errors.New("Unable to read response from server")
	}

	var commands *Commands = &Commands{}

	if err = json.Unmarshal(body, commands); err != nil {
		return commands, err
	}

	return commands, nil
}

func GetTasksByKeywords(keywords string) (*Tasks, error) {
	config, err := configuration.ReadConfig()
	if err != nil {
		return nil, err
	}
	var apiUrl string = config.ApiUrl + "/tasks?keywords=" + keywords

	client := &http.Client{}
	req, err := httputils.NewRmapRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err = httputils.AddBearerToken(req); err != nil {
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("error from server: " + string(response.Status))
	}

	if err != nil {
		return nil, errors.New("Unable to read response from server")
	}

	var tasks *Tasks = &Tasks{}

	if err = json.Unmarshal(body, tasks); err != nil {
		return tasks, err
	}

	return tasks, nil
}
