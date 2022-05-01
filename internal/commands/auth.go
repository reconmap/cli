package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

	formData := map[string]string{
		"username": username,
		"password": password,
	}
	jsonData, err := json.Marshal(formData)

	client := &http.Client{}
	req, err := httputils.NewRmapRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {

		if response.StatusCode == http.StatusForbidden || response.StatusCode == http.StatusUnauthorized {
			return errors.New("Invalid credentials")
		}

		if response.StatusCode == http.StatusMethodNotAllowed {
			return errors.New(fmt.Sprintf("Method POST not allowed for %s. Please make sure you are pointing to the API url and not the frontend one.", apiUrl))
		}

		return errors.New(fmt.Sprintf("Server returned code %d", response.StatusCode))
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return errors.New("Unable to read response from server")
	}

	var loginResponse api.LoginResponse

	if err = json.Unmarshal([]byte(body), &loginResponse); err != nil {
		return err
	}

	err = httputils.SaveSessionToken(loginResponse.AccessToken)
	if err == nil {
		terminal.PrintGreenTick()
		fmt.Printf(" Successfully logged in as '%s'\n", username)
	}

	return err
}

func Logout() error {
	if _, err := httputils.ReadSessionToken(); err != nil {
		return errors.New("There is no active user session")
	}

	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}
	var apiUrl string = config.ApiUrl + "/users/logout"

	client := &http.Client{}
	req, err := httputils.NewRmapRequest("POST", apiUrl, nil)
	if err != nil {
		return err
	}

	if err = httputils.AddBearerToken(req); err != nil {
		return err
	}

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("Response error received from the server")
	}

	defer response.Body.Close()

	configPath, err := httputils.GetSessionTokenPath()
	if _, err := os.Stat(configPath); err == nil {
		err = os.Remove(configPath)
		if err != nil {
			log.Println("Unable to remove file")
		}
	} else if errors.Is(err, os.ErrNotExist) {
		log.Println("warning: Session file does not exist")
	}

	terminal.PrintGreenTick()
	fmt.Printf(" Successfully logged out from the server\n")

	return err
}
