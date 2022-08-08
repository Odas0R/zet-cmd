package main

import (
	"errors"
	"fmt"
	"path"
)

const (
	ZET_EXECUTABLE_PATH = "/home/odas0r/github.com/odas0r/zet-cmd"
)

type Sub struct {
	Templates string
	Fleet     string
	Permanent string
	Journal   string
	Assets    string
}

type Scripts struct {
	Fzf       string
	Open      string
	FindLinks string
	Query     string
}

type Config struct {
	Root    string
	Scripts Scripts
	Sub     Sub
}

func (c *Config) Init() error {
	if c.Root == "" {
		return errors.New("Config path cannot be an empty string")
	}

	if err := initFolderLayout(c); err != nil {
		return err
	}

  if err := initScripts(c); err != nil {
    return err
  }

	return nil
}

func initFolderLayout(config *Config) error {
	var (
		root      = config.Root
		templates = path.Join(root, "templates")
		assets    = path.Join(root, "assets")
		permanent = path.Join(root, "permanent")
		fleet     = path.Join(root, "fleet")
		journal   = path.Join(root, "journal")
	)

	// setup auxiliary paths
	config.Sub.Fleet = fleet
	config.Sub.Permanent = permanent
	config.Sub.Templates = templates
	config.Sub.Journal = journal
	config.Sub.Assets = assets

	// create zet/
	if err := Mkdir(root); err != nil {
		return err
	}

	// create templates/
	if err := Mkdir(templates); err != nil {
		return err
	}

	// create templates/journal.tmpl.md
	if err := CreateFile(journalTmpl, path.Join(templates, "journal.tmpl.md")); err != nil {
		return err
	}

	// create templates/zet.tmpl.md
	if err := CreateFile(zetTmpl, path.Join(templates, "zet.tmpl.md")); err != nil {
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

func initScripts(config *Config) error {
	var (
		root = ZET_EXECUTABLE_PATH
	)

	// setup auxiliary paths
	script := path.Join(root, "scripts/query")
	if exists := FileExists(script); !exists {
		return fmt.Errorf("error: script 'query' does not exist on %s", script)
	}
	config.Scripts.Query = script

	script = path.Join(root, "scripts/fzf")
	if exists := FileExists(script); !exists {
		return fmt.Errorf("error: script 'fzf' does not exist on %s", script)
	}
	config.Scripts.Fzf = script

	script = path.Join(root, "scripts/open")
	if exists := FileExists(script); !exists {
		return fmt.Errorf("error: script 'open' does not exist on %s", script)
	}
	config.Scripts.Open = script

	script = path.Join(root, "scripts/find-links")
	if exists := FileExists(script); !exists {
		return fmt.Errorf("error: script 'find-links' does not exist on %s", script)
	}
	config.Scripts.FindLinks = path.Join(root, "scripts/find-links")

	return nil
}
