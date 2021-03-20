package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	ApiUrl string `json:"api-url"`
}

func SaveConfig(config Config) (string, error) {
	home, err := os.UserHomeDir()

	path := filepath.Join(home, ".reconmap")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	jsondata, _ := json.MarshalIndent(config, "", " ")

	filepath := filepath.Join(path, "config.json")
	err = ioutil.WriteFile(filepath, jsondata, 400)

	return filepath, err
}

func ReadConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	path := filepath.Join(home, ".reconmap", "config.json")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	config := Config{}
	err = json.Unmarshal(bytes, &config)

	return &config
}
