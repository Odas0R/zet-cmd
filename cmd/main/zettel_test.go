package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/samber/lo"
)

func TestZettel(t *testing.T) {
	// initialize the config
	c := &Config{Path: "/tmp/foo"}
	c.Init()

	t.Run("can create a zettel", func(t *testing.T) {
		zettel := &Zettel{Title: "This is a foo title"}

		if err := zettel.New(c); err != nil {
			t.Errorf("error: %s", err)
		}
	})

	t.Run("zettel has correct metadata", func(t *testing.T) {
		zettel := &Zettel{Title: "This is a foo title"}

		if err := zettel.New(c); err != nil {
			t.Errorf("error: %s", err)
		}

		AssertStringEquals(t, "fleet", zettel.Type)
		AssertStringContainsSubstringsNoOrder(t, zettel.FileName, []string{"this-is-a-foo-title"})
		AssertStringContainsSubstringsNoOrder(t, zettel.Path, []string{"/tmp/foo", "this-is-a-foo-title"})
	})

	t.Run("zettel file exists", func(t *testing.T) {
		zettel := &Zettel{Title: "This is a foo title"}

		if err := zettel.New(c); err != nil {
			t.Errorf("error: %s", err)
		}

		if zettelExists := FileExists(zettel.Path); !zettelExists {
			t.Errorf("error: zettel does not exist on path %s", zettel.Path)
		}
	})

	t.Run("got the correct lines from the zettel template file", func(t *testing.T) {
		zettel := &Zettel{Title: "This is a foo title"}
		if err := zettel.New(c); err != nil {
			t.Errorf("error: %s", err)
		}

		AssertStringEquals(t, fmt.Sprintf("# %s", zettel.Title), zettel.Lines[0])
		AssertStringEquals(t, "#example", zettel.Lines[len(zettel.Lines)-1])
	})

	t.Run("can read a zettel on a given path", func(t *testing.T) {
		zettelStub := &Zettel{Title: "This is a foo title"}
		if err := zettelStub.New(c); err != nil {
			t.Errorf("error: %s", err)
		}

		zettel := &Zettel{Path: zettelStub.Path}
		if err := zettel.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}

		AssertStringEquals(t, "This is a foo title", zettel.Title)
		AssertStringEquals(t, "this-is-a-foo-title", zettel.Slug)
		AssertIntEquals(t, int(zettelStub.ID), int(zettel.ID))
		AssertStringEquals(t, zettelStub.FileName, zettel.FileName)
		AssertStringArraysEqualNoOrder(t, []string{"#example"}, zettel.Tags)
	})

	t.Run("can link a zettel and read his links", func(t *testing.T) {
		zettelOne := &Zettel{Title: "This is a foo title"}
		zettelTwo := &Zettel{Title: "This is a another title"}

		if err := zettelOne.New(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelTwo.New(c); err != nil {
			t.Errorf("error: %s", err)
		}

		if err := zettelOne.Link(zettelTwo); err != nil {
			t.Errorf("error: %s", err)
		}
    if err := zettelTwo.Link(zettelOne); err != nil {
			t.Errorf("error: %s", err)
    }

		if err := zettelOne.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelTwo.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}

		AssertStringArraysEqualNoOrder(t, []string{zettelTwo.Path}, zettelOne.Links)
		AssertStringArraysEqualNoOrder(t, []string{zettelOne.Path}, zettelTwo.Links)
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelOne.Lines, ""), []string{zettelTwo.Path})
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelTwo.Lines, ""), []string{zettelOne.Path})
	})

	t.Run("cant link same zettel twice", func(t *testing.T) {
		zettelOne := &Zettel{Title: "This is a foo title"}
		zettelTwo := &Zettel{Title: "This is a another title"}

		// create zettels
		if err := zettelOne.New(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelTwo.New(c); err != nil {
			t.Errorf("error: %s", err)
		}

		if err := zettelOne.Link(zettelTwo); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelOne.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}

		err := zettelOne.Link(zettelTwo)
		ExpectError(t, err, "should have failed, can't link zettels twice")
	})

	t.Run("find-links gives files with the id of the link", func(t *testing.T) {
		zettelOne := &Zettel{ID: 1324, Title: "This is a example title"}
		zettelTwo := &Zettel{ID: 1341, Title: "This is a as title"}
		zettelThree := &Zettel{ID: 1343, Title: "This is a cw title"}

		if err := zettelOne.New(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelTwo.New(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelThree.New(c); err != nil {
			t.Errorf("error: %s", err)
		}

		if err := zettelTwo.Link(zettelOne); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelThree.Link(zettelOne); err != nil {
			t.Errorf("error: %s", err)
		}

		script, err := filepath.Abs("../../scripts/find-links")
		if err != nil {
			t.Errorf("error: %s", err)
		}
		output, err := exec.Command(script, "1324", c.Sub.Fleet, c.Sub.Permanent).Output()
		if err != nil {
			t.Errorf("error: %s", err)
		}

		if len(output) > 0 {
			links := bytes.NewBuffer(output).String()

			AssertStringContainsSubstringsNoOrder(t, links, []string{zettelTwo.Path})
			AssertStringContainsSubstringsNoOrder(t, links, []string{zettelThree.Path})
		} else {
			t.Errorf("error: find-links gave 0 results")
		}

	})

	t.Run("can repair zettel links (1)", func(t *testing.T) {
		zettelOne := &Zettel{ID: 1224, Title: "This is a foo title"}
		zettelTwo := &Zettel{ID: 1241, Title: "This is a boo title"}
		zettelThree := &Zettel{ID: 1243, Title: "This is a bye title"}

		if err := zettelOne.New(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelTwo.New(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelThree.New(c); err != nil {
			t.Errorf("error: %s", err)
		}

		if err := zettelTwo.Link(zettelOne); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelThree.Link(zettelOne); err != nil {
			t.Errorf("error: %s", err)
		}

		// modify the title
		zettelOne.Lines = lo.ReplaceAll(zettelOne.Lines, zettelOne.Lines[0], "# foo bar")
		if err := zettelOne.Write(); err != nil {
			t.Errorf("error: %s", err)
		}

		// repair zettel
		if err := zettelOne.Repair(c); err != nil {
			t.Errorf("error: %s", err)
		}

		AssertStringEquals(t, "foo bar", zettelOne.Title)
		AssertStringEquals(t, "# foo bar", zettelOne.Lines[0])
		AssertStringEquals(t, "foo-bar", zettelOne.Slug)
		AssertStringEquals(t, "fleet", zettelOne.Type)
		AssertStringEquals(t, fmt.Sprintf("foo-bar.%d.md", zettelOne.ID), zettelOne.FileName)
		AssertStringEquals(t, fmt.Sprintf("%s/foo-bar.%d.md", c.Sub.Fleet, zettelOne.ID), zettelOne.Path)

		zettelTwo.Read(c)
		zettelThree.Read(c)

		link := fmt.Sprintf("- [%s](%s)", zettelOne.Title, zettelOne.Path)

		// Check if links are present on the files
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelTwo.Lines, " "), []string{link})
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelThree.Lines, " "), []string{link})
	})

	t.Run("can repair zettel links (2)", func(t *testing.T) {
		zettelOne := &Zettel{Path: fmt.Sprintf("%s/%s", c.Sub.Fleet, "foo-bar.1224.md")}
		zettelTwo := &Zettel{Path: fmt.Sprintf("%s/%s", c.Sub.Fleet, "this-is-a-boo-title.1241.md")}
		zettelThree := &Zettel{Path: fmt.Sprintf("%s/%s", c.Sub.Fleet, "this-is-a-bye-title.1243.md")}

		// Read Zettels
		if err := zettelOne.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelTwo.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelThree.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}

		// Link Zettels
		if err := zettelOne.Link(zettelTwo); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelOne.Link(zettelThree); err != nil {
			t.Errorf("error: %s", err)
		}

		// Read Zettels
		if err := zettelOne.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}

		if err := zettelTwo.WriteLine(0, "# changed title"); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelOne.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}

		if err := zettelTwo.Repair(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelOne.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}

		link := fmt.Sprintf("- [%s](%s)", zettelTwo.Title, zettelTwo.Path)
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelOne.Lines, " "), []string{link})

		// Modify zettelThree
		if err := zettelThree.WriteLine(0, "# pretty different"); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelThree.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelThree.Repair(c); err != nil {
			t.Errorf("error: %s", err)
		}

		// update the modified zettel
		if err := zettelOne.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}

		link = fmt.Sprintf("- [%s](%s)", zettelThree.Title, zettelThree.Path)
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelOne.Lines, " "), []string{link})
	})

	t.Run("can permanent zettel", func(t *testing.T) {
		zettelOne := &Zettel{ID: 1624, Title: "This is a foo title"}
		zettelTwo := &Zettel{ID: 1641, Title: "This is a boo title"}
		zettelThree := &Zettel{ID: 1643, Title: "This is a bye title"}

		if err := zettelOne.New(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelTwo.New(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelThree.New(c); err != nil {
			t.Errorf("error: %s", err)
		}

		if err := zettelTwo.Link(zettelOne); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelThree.Link(zettelOne); err != nil {
			t.Errorf("error: %s", err)
		}

		if err := zettelOne.Permanent(c); err != nil {
			t.Errorf("error: %s", err)
		}

		// Update
		if err := zettelOne.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelTwo.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}
		if err := zettelThree.Read(c); err != nil {
			t.Errorf("error: %s", err)
		}

		AssertStringEquals(t, fmt.Sprintf("%s/%s", c.Sub.Permanent, zettelOne.FileName), zettelOne.Path)
		link := fmt.Sprintf("- [%s](%s)", zettelOne.Title, zettelOne.Path)

		// Check if links are present on the files
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelTwo.Lines, " "), []string{link})
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelThree.Lines, " "), []string{link})
	})

	// cleanup
	err := os.RemoveAll("/tmp/foo")
	if err != nil {
		t.Errorf("error: failed to cleanup")
	}
}
