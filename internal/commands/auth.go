package commands

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc"
	"github.com/reconmap/cli/internal/terminal"
	"github.com/reconmap/shared-lib/pkg/api"
	"github.com/reconmap/shared-lib/pkg/configuration"
	"golang.org/x/oauth2"
)

type IDTokenClaim struct {
	Email string `json:"email"`
}

func Login() error {
	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}

	provider, err := oidc.NewProvider(oauth2.NoContext, config.AuthUrl+"/realms/reconmap")
	if err != nil {
		panic(err)
	}

	oauthConfig := oauth2.Config{
		ClientID:    "web-client",
		RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
		Endpoint:    provider.Endpoint(),
		Scopes:      []string{oidc.ScopeOpenID, "email"},
	}

	var stateSeed uint64
	binary.Read(rand.Reader, binary.LittleEndian, &stateSeed)
	state := fmt.Sprintf("%x", stateSeed)

	authCodeURL := oauthConfig.AuthCodeURL(state)
	fmt.Printf("Open %s\n", authCodeURL)
	fmt.Println()

	fmt.Printf("Enter authorization code: ")
	var code string
	if _, err := fmt.Scanln(&code); err != nil {
		panic(err)
	}

	token, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		panic(err)
	}

	err = api.SaveSessionToken(token.AccessToken)

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		panic("id_token is missing")
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: oauthConfig.ClientID})
	idToken, err := verifier.Verify(oauth2.NoContext, rawIDToken)
	if err != nil {
		panic(err)
	}

	idTokenClaim := IDTokenClaim{}
	if err := idToken.Claims(&idTokenClaim); err != nil {
		panic(err)
	}

	var apiUrl string = config.ApiUrl + "/users/login"

	formData := map[string]string{}
	jsonData, err := json.Marshal(formData)

	client := &http.Client{}
	req, err := api.NewRmapRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	api.AddBearerToken(req)
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

	if err == nil {
		terminal.PrintGreenTick()
		fmt.Printf(" Successfully logged in as '%s'\n", idTokenClaim.Email)
	}

	return err
}

func Logout() error {
	if _, err := api.ReadSessionToken(); err != nil {
		return errors.New("There is no active user session")
	}

	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}
	var apiUrl string = config.ApiUrl + "/users/logout"

	client := &http.Client{}
	req, err := api.NewRmapRequest("POST", apiUrl, nil)
	if err != nil {
		return err
	}

	if err = api.AddBearerToken(req); err != nil {
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

	configPath, err := api.GetSessionTokenPath()
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
