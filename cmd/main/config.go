package main

import (
	"errors"
	"path"
)

const ZET_PATH = "/home/odas0r/github.com/odas0r/zet-cmd"

type Sub struct {
	Templates string `json:"-"`
	Fleet     string `json:"-"`
	Permanent string `json:"-"`
	Journal   string `json:"-"`
	Assets    string `json:"-"`
}

type Config struct {
	Path string `json:"path"`
	Sub  Sub
}

func (c *Config) Init() (error) {
  if c.Path == "" {
    return errors.New("Config path cannot be an empty string")
  }

  if error := initFolderLayout(c); error != nil {
    return error
  }

  return nil
}

func initFolderLayout(config *Config) error {
	var (
		root      = config.Path
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
