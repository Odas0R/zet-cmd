package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/odas0r/zet/cmd/assert"
	"github.com/samber/lo"
)

func TestZettel(t *testing.T) {
	config := &Config{Root: "/tmp/foo"}
	history := &History{Root: "/tmp/foo"}

	// initialize config
	err := config.Init()
	assert.Equal(t, err, nil, "config.Init should not fail")

	// initialize history
	err = history.Init()
	assert.Equal(t, err, nil, "history.Init should not fail")

	t.Run("can create a zettel", func(t *testing.T) {
		zettel := &Zettel{ID: 1, Title: "Title example"}
		err := zettel.New(config)
		assert.Equal(t, err, nil, "zettel.New should not fail")

		zettel = &Zettel{ID: 2, Title: "Title example"}
		err = zettel.New(config)
		assert.Equal(t, err, nil, "zettel.New should not fail")

		zettel = &Zettel{ID: 3, Title: "Title example"}
		err = zettel.New(config)
		assert.Equal(t, err, nil, "zettel.New should not fail")
	})

	t.Run("zettel file exists", func(t *testing.T) {
		zettel := &Zettel{Path: "/tmp/foo/fleet/title-example.2.md"}

		zettelExists := FileExists(zettel.Path)
		assert.Equal(t, zettelExists, true, "zettel should exist")
	})

	t.Run("can read a zettel on a given path", func(t *testing.T) {
		zettel := &Zettel{Path: "/tmp/foo/fleet/title-example.2.md"}

		err := zettel.Read(config)
		assert.Equal(t, err, nil, "zettel.Read should not fail")

		assert.Equal(t, zettel.ID, int64(2), "zettel.ID should be correct")
		assert.Equal(t, zettel.Type, "fleet", "zettel.Type should be correct")
		assert.Equal(t, zettel.Slug, "title-example", "zettel.Slug should be correct")
		assert.Equal(t, zettel.Path, "/tmp/foo/fleet/title-example.2.md", "zettel.Path should be correct")
		assert.Equal(t, fmt.Sprintf("# %s", zettel.Title), zettel.Lines[0], "zettel.Title should be on line 0")
		assert.Equal(t, "#example", zettel.Lines[len(zettel.Lines)-1], "tag #example should be on the last line")
	})

	t.Run("can link a zettel and read his links", func(t *testing.T) {
		zettelOne := &Zettel{Path: "/tmp/foo/fleet/title-example.1.md"}
		zettelTwo := &Zettel{Path: "/tmp/foo/fleet/title-example.2.md"}
		zettelThree := &Zettel{Path: "/tmp/foo/fleet/title-example.3.md"}

		//
		// Read
		//

		err := zettelOne.Read(config)
		assert.Equal(t, err, nil, "zettelOne.Read should not fail")

		err = zettelTwo.Read(config)
		assert.Equal(t, err, nil, "zettelTwo.Read should not fail")

		err = zettelThree.Read(config)
		assert.Equal(t, err, nil, "zettelThree.Read should not fail")

		//
		// Link
		//

		err = zettelOne.Link(zettelTwo)
		assert.Equal(t, err, nil, "zettelOne should link to zettelTwo")

		err = zettelTwo.Link(zettelOne)
		assert.Equal(t, err, nil, "zettelTwo should link to zettelOne")

		err = zettelThree.Link(zettelOne)
		assert.Equal(t, err, nil, "zettelThree should link to zettelOne")

		//
		// Read
		//

		err = zettelOne.Read(config)
		assert.Equal(t, err, nil, "zettelOne.Read should not fail")

		err = zettelTwo.Read(config)
		assert.Equal(t, err, nil, "zettelTwo.Read should not fail")

		err = zettelThree.Read(config)
		assert.Equal(t, err, nil, "zettelThree.Read should not fail")

		containsLink := func(zettel *Zettel, zettelToLink *Zettel) bool {
			foundZettelLink := false
			for _, line := range zettel.Lines {
				hasLink := strings.Contains(line, zettelToLink.Path)
				if hasLink {
					foundZettelLink = hasLink
				}
			}
			return foundZettelLink
		}

		assert.Equal(t, zettelOne.Links[0], "/tmp/foo/fleet/title-example.2.md", "zettelOne contains zettelTwo link")
		assert.Equal(t, zettelTwo.Links[0], "/tmp/foo/fleet/title-example.1.md", "zettelTwo contains zettelOne link")
		assert.Equal(t, zettelThree.Links[0], "/tmp/foo/fleet/title-example.1.md", "zettelThree contains zettelOne link")
		assert.Equal(t, containsLink(zettelOne, zettelTwo), true, "zettelOne contains zettelTwo link on zettel.Lines")
		assert.Equal(t, containsLink(zettelTwo, zettelOne), true, "zettelTwo contains zettelOne link on zettel.Lines")
		assert.Equal(t, containsLink(zettelThree, zettelOne), true, "zettelThree contains zettelOne link on zettel.Lines")
	})

	t.Run("cant link same zettel twice", func(t *testing.T) {
		zettelOne := &Zettel{Path: "/tmp/foo/fleet/title-example.1.md"}
		zettelTwo := &Zettel{Path: "/tmp/foo/fleet/title-example.2.md"}

		//
		// Read
		//

		err := zettelOne.Read(config)
		assert.Equal(t, err, nil, "zettelOne.Read should not fail")

		err = zettelTwo.Read(config)
		assert.Equal(t, err, nil, "zettelTwo.Read should not fail")

		err = zettelOne.Link(zettelTwo)
		assert.NotEqual(t, err, nil, "zettelOne.Link should fail")
	})

	t.Run("find-links gives files with the id of the link", func(t *testing.T) {
		zettelOne := &Zettel{Path: "/tmp/foo/fleet/title-example.1.md"}
		zettelTwo := &Zettel{Path: "/tmp/foo/fleet/title-example.2.md"}

		//
		// Read
		//

		err := zettelOne.Read(config)
		assert.Equal(t, err, nil, "zettelOne.Read should not fail")

		err = zettelTwo.Read(config)
		assert.Equal(t, err, nil, "zettelTwo.Read should not fail")

		cmd := exec.Command("/bin/bash", config.Scripts.FindLinks, "2", config.Sub.Fleet, config.Sub.Permanent)

		output, err := cmd.Output()
		if err != nil {
			t.Errorf("error: %s", err)
		}

		if len(output) > 0 {
			links := bytes.NewBuffer(output).String()
			linkPath := strings.Split(links, ":")[1]

			assert.Equal(t, strings.TrimSpace(linkPath), "/tmp/foo/fleet/title-example.1.md", "find-links should give zettelOne.Path")
		} else {
			t.Errorf("error: find-links gave 0 results")
		}

	})

	t.Run("can repair zettel links (1)", func(t *testing.T) {
		zettelOne := &Zettel{Path: "/tmp/foo/fleet/title-example.1.md"}
		zettelTwo := &Zettel{Path: "/tmp/foo/fleet/title-example.2.md"}
		zettelThree := &Zettel{Path: "/tmp/foo/fleet/title-example.3.md"}

		//
		// Read
		//

		err := zettelOne.Read(config)
		assert.Equal(t, err, nil, "zettelOne.Read should not fail")
		err = zettelTwo.Read(config)
		assert.Equal(t, err, nil, "zettelTwo.Read should not fail")
		err = zettelThree.Read(config)
		assert.Equal(t, err, nil, "zettelThree.Read should not fail")

		// modify the title
		zettelOne.Lines = lo.ReplaceAll(zettelOne.Lines, zettelOne.Lines[0], "# foo bar")

		// Write
		err = zettelOne.Write()
		assert.Equal(t, err, nil, "zettelOne.Write should not fail")

		// Repair zettel
		err = zettelOne.Repair(config, history)
		assert.Equal(t, err, nil, "zettelOne.Repair should not fail")

		assert.Equal(t, zettelOne.Title, "foo bar", "zettelOne.Title should be correct")
		assert.Equal(t, zettelOne.Lines[0], "# foo bar", "zettelOne.Lines[0] should be correct")
		assert.Equal(t, zettelOne.Slug, "foo-bar", "zettelOne.Slug should be correct")
		assert.Equal(t, zettelOne.Type, "fleet", "zettelOne.Type should be correct")
		assert.Equal(t, zettelOne.FileName, "foo-bar.1.md", "zettelOne.FileName should be correct")
		assert.Equal(t, zettelOne.Path, "/tmp/foo/fleet/foo-bar.1.md", "zettelOne.Path should be correct")

		err = zettelTwo.Read(config)
		assert.Equal(t, err, nil, "zettelTwo.Read should not fail")

		err = zettelThree.Read(config)
		assert.Equal(t, err, nil, "zettelThree.Read should not fail")

		containsLink := func(zettel *Zettel, zettelToLink *Zettel) bool {
			foundZettelLink := false
			for _, line := range zettel.Lines {
				hasLink := strings.Contains(line, zettelToLink.Path)
				if hasLink {
					foundZettelLink = hasLink
				}
			}
			return foundZettelLink
		}

		assert.Equal(t, containsLink(zettelTwo, zettelOne), true, "zettelOne is linked to zettelTwo")
		assert.Equal(t, containsLink(zettelThree, zettelOne), true, "zettelOne is linked to zettelThree")
	})

	t.Run("can permanent zettel", func(t *testing.T) {
		zettelOne := &Zettel{Path: "/tmp/foo/fleet/foo-bar.1.md"}
		zettelTwo := &Zettel{Path: "/tmp/foo/fleet/title-example.2.md"}
		zettelThree := &Zettel{Path: "/tmp/foo/fleet/title-example.3.md"}

		//
		// Read
		//

		err := zettelOne.Read(config)
		assert.Equal(t, err, nil, "zettelOne.Read should not fail")
		err = zettelTwo.Read(config)
		assert.Equal(t, err, nil, "zettelTwo.Read should not fail")
		err = zettelThree.Read(config)
		assert.Equal(t, err, nil, "zettelThree.Read should not fail")

		err = zettelOne.Permanent(config, history)
		assert.Equal(t, err, nil, "zettelOne.Permanent should not fail")
    
		err = zettelOne.Read(config)
		assert.Equal(t, err, nil, "zettelOne.Read should not fail")
		err = zettelTwo.Read(config)
		assert.Equal(t, err, nil, "zettelTwo.Read should not fail")
		err = zettelThree.Read(config)
		assert.Equal(t, err, nil, "zettelThree.Read should not fail")

    assert.Equal(t, zettelOne.Type, "permanent", "zettelOne.Type should be correct")
    assert.Equal(t, zettelOne.Path, "/tmp/foo/permanent/foo-bar.1.md", "zettelOne.Path should be correct")
    assert.Equal(t, strings.Contains(zettelTwo.Links[0], "foo-bar.1.md"), true, "zettelOne is linked to zettelTwo")
    assert.Equal(t, strings.Contains(zettelThree.Links[0], "foo-bar.1.md"), true, "zettelOne is linked to zettelThree")
	})

	// cleanup
	err = os.RemoveAll("/tmp/foo")
	assert.Equal(t, err, nil, "os.RemoveAll should not fail")
}
