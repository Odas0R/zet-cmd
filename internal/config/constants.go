package config

import (
	"github.com/odas0r/zet/pkg/fs"
)

type Config struct {
	Root          string
	FleetRoot     string
	PermanentRoot string
}

func New(root string) *Config {
	cfg := &Config{
		Root:          root,
		FleetRoot:     root + "/fleet",
		PermanentRoot: root + "/permanent",
	}

	if err := cfg.createRoot(); err != nil {
		panic(err)
	}

	return cfg
}

func (c *Config) createRoot() error {

	if err := fs.Mkdir(c.Root); err != nil {
		return err
	}

	if err := fs.Mkdir(c.FleetRoot); err != nil {
		return err
	}

	if err := fs.Mkdir(c.PermanentRoot); err != nil {
		return err
	}

	return nil
}
