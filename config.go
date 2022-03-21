package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	configFile = ".zet"
	rootDir    = "zet"
)

type Config struct {
	RootDir string `json:"rootDir"`
}

func (c Config) Init() error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(dirname, configFile)
	rootDirPath := filepath.Join(dirname, rootDir)

	data, err := ioutil.ReadFile(configPath)

	// file does not exist
	if errors.Is(err, os.ErrNotExist) {
		// create a new config file
		_, err := os.Create(configPath)
		if err != nil {
			return err
		}

		// c.Root
		c.RootDir = rootDirPath

		// append data to the config file
		file, _ := json.MarshalIndent(c, "", " ")
		_ = ioutil.WriteFile(c.RootDir, file, 0644)

		return nil
	}

	// append data from the json file to the config struct
	err = json.Unmarshal(data, &c)
	if err != nil {
		return err
	}

  // TODO:
  // create the zettelkasten layout

	return nil
}

func (c Config) Edit() {
}
