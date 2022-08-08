package main

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	c := Config{Root: "/tmp/foo"}

	t.Run("can initialize config", func(t *testing.T) {
		err := c.Init()
		ExpectNoError(t, err, "Config.Init should not fail")
	})

	t.Run("config has the right values", func(t *testing.T) {
		AssertStringEquals(t, "/tmp/foo", c.Root)
		AssertStringEquals(t, "/tmp/foo/assets", c.Sub.Assets)
		AssertStringEquals(t, "/tmp/foo/journal", c.Sub.Journal)
		AssertStringEquals(t, "/tmp/foo/templates", c.Sub.Templates)
		AssertStringEquals(t, "/tmp/foo/permanent", c.Sub.Permanent)
		AssertStringEquals(t, "/tmp/foo/fleet", c.Sub.Fleet)
	})

	t.Run("config folders were created", func(t *testing.T) {
		if assets := FileExists("/tmp/foo/assets"); !assets {
			t.Errorf("assets folder does not exists")
		}

		if journal := FileExists("/tmp/foo/journal"); !journal {
			t.Errorf("journal folder does not exists")
		}

		if templates := FileExists("/tmp/foo/templates"); !templates {
			t.Errorf("templates folder does not exists")
		}

		if permanent := FileExists("/tmp/foo/permanent"); !permanent {
			t.Errorf("permanent folder does not exists")
		}

		if fleet := FileExists("/tmp/foo/fleet"); !fleet {
			t.Errorf("fleet folder does not exists")
		}
	})

	// cleanup
	err := os.RemoveAll("/tmp/foo")
	if err != nil {
		t.Errorf("failed to cleanup")
	}
}
