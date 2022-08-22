package main

import (
	"testing"

	"github.com/odas0r/zet/cmd/assert"
)

func TestConfig(t *testing.T) {
	t.Run("config has the right values", func(t *testing.T) {
		assert.Equal(t, config.Path, "/tmp/foo", "config.root must be correct")
		assert.Equal(t, config.Sub.Assets, "/tmp/foo/assets", "config.sub.assets must be correct")
		assert.Equal(t, config.Sub.Journal, "/tmp/foo/journal", "config.sub.journal must be correct")
		assert.Equal(t, config.Sub.Templates, "/tmp/foo/templates", "config.sub.templates must be correct")
		assert.Equal(t, config.Sub.Permanent, "/tmp/foo/permanent", "config.sub.permanent must be correct")
		assert.Equal(t, config.Sub.Fleet, "/tmp/foo/fleet", "config.sub.fleet must be correct")
	})

	t.Run("config folders were created", func(t *testing.T) {
		assetsExists := FileExists("/tmp/foo/assets")
		assert.Equal(t, assetsExists, true, "assets file must exist")

		journalExists := FileExists("/tmp/foo/journal")
		assert.Equal(t, journalExists, true, "journal file must exist")

		templatesExists := FileExists("/tmp/foo/templates")
		assert.Equal(t, templatesExists, true, "templates file must exist")

		permanentExists := FileExists("/tmp/foo/permanent")
		assert.Equal(t, permanentExists, true, "permanent file must exist")

		fleetExists := FileExists("/tmp/foo/fleet")
		assert.Equal(t, fleetExists, true, "fleet file must exist")
	})
}
