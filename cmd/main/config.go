package main

import (
	"errors"
	"path"
)

var config = &Config{}

type Sub struct {
	Templates string
	Fleet     string
	Permanent string
	Journal   string
	Assets    string
}

type Scripts struct {
	Fzf       string
	FzfMulti  string
	Open      string
	FindLinks string
	Query     string
	Clear     string
	Ripgrep   string
}

type Config struct {
	Path    string
	Scripts Scripts
	Sub     Sub
}

func (c *Config) Init(configPath string) error {
	if configPath == "" {
		return errors.New("error: config path cannot be an empty string")
	}

	c.Path = configPath

	// set paths of config
	c.Sub.Fleet = path.Join(c.Path, "fleet")
	c.Sub.Permanent = path.Join(c.Path, "permanent")
	c.Sub.Templates = path.Join(c.Path, "templates")
	c.Sub.Journal = path.Join(c.Path, "journal")
	c.Sub.Assets = path.Join(c.Path, "assets")

	// set scripts paths
	zetExecutablePath := "/home/odas0r/github.com/odas0r/zet-cmd"

	c.Scripts.Query = path.Join(zetExecutablePath, "scripts/query")
	c.Scripts.Fzf = path.Join(zetExecutablePath, "scripts/fzf")
	c.Scripts.FzfMulti = path.Join(zetExecutablePath, "scripts/fzf-multi")
	c.Scripts.Open = path.Join(zetExecutablePath, "scripts/open")
	c.Scripts.Clear = path.Join(zetExecutablePath, "scripts/clear")
	c.Scripts.FindLinks = path.Join(zetExecutablePath, "scripts/find-links")
	c.Scripts.Ripgrep = path.Join(zetExecutablePath, "scripts/ripgrep")

	if err := c.setupLayout(); err != nil {
		return err
	}

	return nil
}

func (c *Config) setupLayout() error {
	// create zet/
	if err := Mkdir(c.Path); err != nil {
		return err
	}

	// create templates/
	if err := Mkdir(c.Sub.Templates); err != nil {
		return err
	}

	// create templates/journal.tmpl.md
	if err := NewFile(journalTmpl, path.Join(c.Sub.Templates, "journal.tmpl.md")); err != nil {
		return err
	}

	// create templates/zet.tmpl.md
	if err := NewFile(zetTmpl, path.Join(c.Sub.Templates, "zet.tmpl.md")); err != nil {
		return err
	}

	// create assets/
	if err := Mkdir(c.Sub.Assets); err != nil {
		return err
	}

	// create fleet/
	if err := Mkdir(c.Sub.Fleet); err != nil {
		return err
	}

	// create permanent/
	if err := Mkdir(c.Sub.Permanent); err != nil {
		return err
	}

	// create journal/
	if err := Mkdir(c.Sub.Journal); err != nil {
		return err
	}

	return nil
}
