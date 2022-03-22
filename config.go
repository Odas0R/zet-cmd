package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
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

// ----------------------- utility -----------------------

// templates/
// assets/
// fleet/
// permanent/
// journal/
//  goals.md
//  habits.md
//  ideas.md
//  inspiration.md
//  todos.md
func setupLayout(c Config) error {
	var (
		root      = c.RootDir
		templates = path.Join(c.RootDir, "templates")
		assets    = path.Join(c.RootDir, "assets")
		permanent = path.Join(c.RootDir, "permanent")
		journal   = path.Join(c.RootDir, "permanent")
	)

	// templates setup
	if err := os.Mkdir(templates, 0755); err != nil {
		return err
	}
	// add two templates, journal and zettel

	if err := os.Mkdir(templates, 0755); err != nil {
		return err
	}

	return nil
}

// If the file doesn't exist, create it, or append to the file
func appendTextToFile(filepath string, text string) error {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
    return err
	}
	if _, err := f.WriteString(text); err != nil {
    return err
	}
	if err := f.Close(); err != nil {
    return err
	}

  return nil
}
