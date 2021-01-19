package commands

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func Login(username string, password string) (string, error) {
	var apiUrl string = "https://api.reconmap.org/users/login"
	apiUrl = "http://localhost:8080/users/login"
	response, err := http.PostForm(apiUrl, url.Values{
		"username": {username},
		"password": {password},
	})

	if err != nil {
		return "", err
	}

	if response.StatusCode == 403 {
		return "", errors.New("Invalid credentials")
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
