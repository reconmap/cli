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

const configFileName = "config.json"

func GetReconmapConfigDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".reconmap"), nil

}

func SaveConfig(config Config) (string, error) {
	var reconmapConfigDir, err = GetReconmapConfigDirectory()

	if _, err := os.Stat(reconmapConfigDir); os.IsNotExist(err) {
		os.MkdirAll(reconmapConfigDir, os.ModePerm)
	}

	jsondata, _ := json.MarshalIndent(config, "", " ")

	filepath := filepath.Join(reconmapConfigDir, configFileName)
	err = ioutil.WriteFile(filepath, jsondata, 400)

	return filepath, err
}

func ReadConfig() (*Config, error) {
	var reconmapConfigDir, err = GetReconmapConfigDirectory()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(reconmapConfigDir, configFileName)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	config := Config{}
	err = json.Unmarshal(bytes, &config)

	return &config, nil
}
