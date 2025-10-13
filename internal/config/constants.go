package config

import (
	"github.com/odas0r/zet/pkg/fs"
)

type Config struct {
	Path string
}

func New(path string) *Config {
	cfg := &Config{
		Path: path,
	}

	if err := cfg.createRoot(); err != nil {
		panic(err)
	}

	return cfg
}

func (c *Config) createRoot() error {

	if err := fs.Mkdir(c.Path); err != nil {
		return err
	}

	return nil
}
