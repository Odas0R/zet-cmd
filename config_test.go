package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigInit(t *testing.T) {
	c := Config{}

	t.Run("can initialize config", func(t *testing.T) {
		// check if config can be initialized
		err := c.Init()
		if err != nil {
			t.Errorf("error: %q", err)
		}

		// check if file exists
		configFile := getFullPath(".zet.json")
		_, err = os.Stat(configFile)
		if os.IsNotExist(err) {
			t.Errorf("error: %q", err)
		}
	})

	t.Run("initialized config paths are correct", func(t *testing.T) {
		got := c.RootPath
		want := getFullPath("zet")

		if got != want {
			t.Errorf("got: %q, want: %q", got, want)
		}
	})
}

// ------------------ utilities ------------------

func getFullPath(path string) string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(dirname, path)
}
