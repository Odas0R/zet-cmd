package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

const (
	configFile = ".zet.json"
	rootDir    = "zet"
)

type Config struct {
	RootPath string `json:"rootPath"`
	Editor   string `json:"editor"`

	// Auxiliary
	TemplatesPath string `json:"-"`
	FleetPath     string `json:"-"`
	PermanentPath string `json:"-"`
	JournalPath   string `json:"-"`
	AssetsPath    string `json:"-"`
}

func (c *Config) Init() error {
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

		c.RootPath = rootDirPath
		c.Editor = ""

		// append data to the config file
		file, _ := json.MarshalIndent(c, "", " ")
		_ = ioutil.WriteFile(configPath, file, 0644)

		return nil
	}

	// append data from the json file to the config struct
	err = json.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	// create the zettelkasten layout and inserts the auxiliary variables
	err = SetupLayout(c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) LookEditor() string {
	e := c.Editor
	if e != "" {
		return e
	}
	e = os.Getenv("EDITOR")
	if e != "" {
		return e
	}
	e = os.Getenv("VISUAL")
	if e != "" {
		return e
	}
	path, err := exec.LookPath("vim")
	if err != nil {
		return path
	}
	return ""
}

// TODO: try to port the script query and open to here
//
// 1. if there's tmux, check for fzf-tmux
// 2. if there's not tmux, check for fzf
// 3. if there's not fzf, throw error saying that they need to install
// fzf
func (c *Config) Edit() {
	editor := c.LookEditor()

	cmd := exec.Command(editor, c.RootPath)
	cmd.Start()
}

// Create the zettelkasten directory layout
func SetupLayout(c *Config) error {
	var (
		root      = c.RootPath
		templates = path.Join(root, "templates")
		assets    = path.Join(root, "assets")
		permanent = path.Join(root, "permanent")
		fleet     = path.Join(root, "fleet")
		journal   = path.Join(root, "journal")
	)

	// setup auxiliary paths
	c.FleetPath = fleet
	c.PermanentPath = permanent
	c.TemplatesPath = templates
	c.JournalPath = journal
	c.AssetsPath = assets

	// create zet/
	if err := Mkdir(root); err != nil {
		return err
	}

	// create templates/
	if err := Mkdir(templates); err != nil {
		return err
	}

	// create templates/journal.tmpl.md
	if err := Cat(journalTmpl, path.Join(templates, "journal.tmpl.md")); err != nil {
		return err
	}

	// create templates/zet.tmpl.md
	if err := Cat(zetTmpl, path.Join(templates, "zet.tmpl.md")); err != nil {
		return err
	}

	// create assets/
	if err := Mkdir(assets); err != nil {
		return err
	}

	// create fleet/
	if err := Mkdir(fleet); err != nil {
		return err
	}

	// create permanent/
	if err := Mkdir(permanent); err != nil {
		return err
	}

	// create journal/
	if err := Mkdir(journal); err != nil {
		return err
	}

	return nil
}
