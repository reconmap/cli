package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/reconmap/cli/internal/api"
	"github.com/reconmap/cli/internal/configuration"
	"github.com/reconmap/cli/internal/httputils"
	"github.com/reconmap/cli/internal/terminal"
)

func Login(username string, password string) error {
	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}
	var apiUrl string = config.ApiUrl + "/users/login"

	formData := url.Values{
		"username": {username},
		"password": {password},
	}

	client := &http.Client{}
	req, err := httputils.NewRmapRequest("POST", apiUrl, strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {

		if response.StatusCode == http.StatusForbidden || response.StatusCode == http.StatusUnauthorized {
			return errors.New("Invalid credentials")
		}

		return errors.New("Response error received from the server")
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return errors.New("Unable to read response from server")
	}

	var loginResponse api.LoginResponse

	json.Unmarshal([]byte(body), &loginResponse)

	err = httputils.SaveSessionToken(loginResponse.AccessToken)
	if err == nil {
		terminal.PrintGreenTick()
		fmt.Printf(" Successfully logged in as '%s'\n", username)
	}

	return err
}

func Logout() error {
	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}
	var apiUrl string = config.ApiUrl + "/users/logout"

	formData := url.Values{}

	client := &http.Client{}
	req, err := httputils.NewRmapRequest("POST", apiUrl, strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httputils.AddBearerToken(req)
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {

		if response.StatusCode == http.StatusForbidden || response.StatusCode == http.StatusUnauthorized {
			return errors.New("Invalid credentials")
		}

		return errors.New("Response error received from the server")
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return errors.New("Unable to read response from server")
	}

	var loginResponse api.LoginResponse

	json.Unmarshal([]byte(body), &loginResponse)

	configPath, err := httputils.GetSessionTokenPath()

	err = os.Remove(configPath)
	if err == nil {
		terminal.PrintGreenTick()
		fmt.Printf(" Successfully logged out from the server\n")
	}

	return err
}
