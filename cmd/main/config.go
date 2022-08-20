package main

import (
	"errors"
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
	Ripgrep   string
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
	config.Scripts.Query = path.Join(root, "scripts/query")
	config.Scripts.Fzf = path.Join(root, "scripts/fzf")
	config.Scripts.Open = path.Join(root, "scripts/open")
	config.Scripts.FindLinks = path.Join(root, "scripts/find-links")
	config.Scripts.Ripgrep = path.Join(root, "scripts/ripgrep")

	return nil
}
