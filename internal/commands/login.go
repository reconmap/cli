package commands

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/reconmap/cli/internal/httputils"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func Login(username string, password string) (string, error) {
	var apiUrl string = "https://api.reconmap.org/users/login"
	apiUrl = "http://localhost:8080/users/login"

	formData := url.Values{
		"username": {username},
		"password": {password},
	}

	client := &http.Client{}
	req, err := httputils.NewRmapRequest("POST", apiUrl, strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {

		if response.StatusCode == http.StatusForbidden {
			return "", errors.New("Invalid credentials")
		}
		if response.StatusCode == http.StatusUnauthorized {
			return "", errors.New("Invalid crednetials")
		}

		return "", errors.New("Response error received from the server")
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", errors.New("Unable to read response from server")
	}

	var loginResponse LoginResponse

	json.Unmarshal([]byte(body), &loginResponse)

	err = ioutil.WriteFile("rmap-session", []byte(loginResponse.AccessToken), 0644)

	return "Successful login", nil
}
