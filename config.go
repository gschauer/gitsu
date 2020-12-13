package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	Users []user `json:"Users"`
}

var configPath string
var errNoUser = errors.New("no users")

// init determines the configPath.
func init() {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln("Cannot resolve config directory:", err)
	}
	configPath = filepath.Join(dir, "gitsu-go", "config.json")
}

// configDir resolves the config directory and creates it, if necessary.
func configDir() string {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		log.Fatalln("Cannot create config directory:", err)
	}
	return dir
}

// configFile opens the config file in read-write mode.
func configFile(flag int) *os.File {
	configDir()
	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_RDWR|flag, 0600)
	if err != nil {
		log.Fatalln("cannot access config file", configPath)
	}
	return file
}

// writeConfig saves the config into a file.
func writeConfig(config *config) error {
	file := configFile(os.O_TRUNC)
	defer file.Close()

	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

// readConfig loads the config from a file.
func readConfig() (*config, error) {
	config := &config{}
	if _, err := os.Stat(configPath); err != nil {
		return config, errNoUser
	}

	file := configFile(0)
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}
