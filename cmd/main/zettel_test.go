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
		err := zettel.New(c)
		ExpectNoError(t, err, fmt.Sprintf("error: failed to create a zettel %v", zettel))
	})

	t.Run("zettel has correct metadata", func(t *testing.T) {
		zettel := &Zettel{Title: "This is a foo title"}
		err := zettel.New(c)
		ExpectNoError(t, err, "error: failed to create a zettel")

		AssertStringEquals(t, "fleet", zettel.Type)
		AssertStringContainsSubstringsNoOrder(t, zettel.FileName, []string{"this-is-a-foo-title"})
		AssertStringContainsSubstringsNoOrder(t, zettel.Path, []string{"/tmp/foo", "this-is-a-foo-title"})
	})

	t.Run("zettel file exists", func(t *testing.T) {
		zettel := &Zettel{Title: "This is a foo title"}
		zettel.New(c)

		zettelExists := FileExists(zettel.Path)
		if !zettelExists {
			t.Errorf("error: zettel does not exist on path %s", zettel.Path)
		}
	})

	t.Run("got the correct lines from the zettel template file", func(t *testing.T) {
		zettel := &Zettel{Title: "This is a foo title"}
		zettel.New(c)

		AssertStringEquals(t, fmt.Sprintf("# %s", zettel.Title), zettel.Lines[0])
		AssertStringEquals(t, "#example", zettel.Lines[len(zettel.Lines)-1])
	})

	t.Run("can read a zettel on a given path", func(t *testing.T) {
		zettelStub := &Zettel{Title: "This is a foo title"}
		zettelStub.New(c)

		zettel := &Zettel{Path: zettelStub.Path}

		err := zettel.Read()
		ExpectNoError(t, err, "error: failed to read zettel")

		AssertStringEquals(t, "This is a foo title", zettel.Title)
		AssertStringEquals(t, "this-is-a-foo-title", zettel.Slug)
		AssertIntEquals(t, int(zettelStub.ID), int(zettel.ID))
		AssertStringEquals(t, zettelStub.FileName, zettel.FileName)
		AssertStringArraysEqualNoOrder(t, []string{"#example"}, zettel.Tags)
	})

	t.Run("can link a zettel and read his links", func(t *testing.T) {
		zettelOne := &Zettel{Title: "This is a foo title"}
		zettelOne.New(c)

		zettelTwo := &Zettel{Title: "This is a another title"}
		zettelTwo.New(c)

		err := zettelOne.Link(zettelTwo)
		ExpectNoError(t, err, "error: failed to link zettel")
		AssertStringArraysEqualNoOrder(t, []string{zettelTwo.Path}, zettelOne.Links)

		err = zettelOne.Read()
		ExpectNoError(t, err, "error: failed to read zettel")
		AssertStringArraysEqualNoOrder(t, []string{zettelTwo.Path}, zettelOne.Links)
		AssertStringContainsNoneOfTheSubstrings(t, strings.Join(zettelOne.Lines, ""), []string{zettelOne.Path})
	})

	t.Run("cant link same zettel twice", func(t *testing.T) {
		zettelOne := &Zettel{Title: "This is a foo title"}
		zettelOne.New(c)

		zettelTwo := &Zettel{Title: "This is a another title"}
		zettelTwo.New(c)

		err := zettelOne.Link(zettelTwo)
		ExpectNoError(t, err, "error: failed to link zettel")

		err = zettelOne.Link(zettelTwo)
		ExpectError(t, err, "should have failed, can't link zettels twice")
	})

	t.Run("find-links gives files with the id of the link", func(t *testing.T) {
		zettelOne := &Zettel{ID: 1324, Title: "This is a example title"}
		zettelOne.New(c)

		zettelTwo := &Zettel{ID: 1341, Title: "This is a as title"}
		zettelTwo.New(c)

		zettelThree := &Zettel{ID: 1343, Title: "This is a cw title"}
		zettelThree.New(c)

		zettelTwo.Link(zettelOne)
		zettelThree.Link(zettelOne)

		script, _ := filepath.Abs("../../scripts/find-links")
		output, _ := exec.Command(script, "1324", c.Sub.Fleet, c.Sub.Permanent).Output()

		if len(output) > 0 {
			links := strings.Split(bytes.NewBuffer(output).String(), "\n")
			links = links[:len(links)-1]

			AssertStringContainsSubstringsNoOrder(t, links[0], []string{zettelTwo.Path})
			AssertStringContainsSubstringsNoOrder(t, links[1], []string{zettelThree.Path})
		} else {
			t.Errorf("error: find-links gave 0 results, should have given 2")
		}

	})

	t.Run("can repair zettel links (1)", func(t *testing.T) {
		zettelOne := &Zettel{ID: 1224, Title: "This is a foo title"}
		zettelOne.New(c)

		zettelTwo := &Zettel{ID: 1241, Title: "This is a boo title"}
		zettelTwo.New(c)

		zettelThree := &Zettel{ID: 1243, Title: "This is a bye title"}
		zettelThree.New(c)

		if err := zettelTwo.Link(zettelOne); err != nil {
			t.Errorf("error: failed to link zettel")
		}
		if err := zettelThree.Link(zettelOne); err != nil {
			t.Errorf("error: failed to link zettel")
		}

		// modify the title
		zettelOne.Lines = lo.ReplaceAll(zettelOne.Lines, zettelOne.Lines[0], "# foo bar")
		if err := zettelOne.Write(); err != nil {
			t.Errorf("error: failed to write zettel")
		}

		// repair zettel
		if err := zettelOne.Repair(c); err != nil {
			t.Errorf("error: failed to repair zettel")
		}

		AssertStringEquals(t, "foo bar", zettelOne.Title)
		AssertStringEquals(t, "# foo bar", zettelOne.Lines[0])
		AssertStringEquals(t, "foo-bar", zettelOne.Slug)
		AssertStringEquals(t, "fleet", zettelOne.Type)
		AssertStringEquals(t, fmt.Sprintf("foo-bar.%d.md", zettelOne.ID), zettelOne.FileName)
		AssertStringEquals(t, fmt.Sprintf("%s/foo-bar.%d.md", c.Sub.Fleet, zettelOne.ID), zettelOne.Path)

		zettelTwo.Read()
		zettelThree.Read()

		link := fmt.Sprintf("- [%s](%s)", zettelOne.Title, zettelOne.Path)

		// Check if links are present on the files
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelTwo.Lines, " "), []string{link})
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelThree.Lines, " "), []string{link})
	})

	t.Run("can repair zettel links (2)", func(t *testing.T) {
		zettelOne := &Zettel{Path: fmt.Sprintf("%s/%s", c.Sub.Fleet, "foo-bar.1224.md")}
		zettelTwo := &Zettel{Path: fmt.Sprintf("%s/%s", c.Sub.Fleet, "this-is-a-boo-title.1241.md")}
		zettelThree := &Zettel{Path: fmt.Sprintf("%s/%s", c.Sub.Fleet, "this-is-a-bye-title.1243.md")}

		if err := zettelOne.Link(zettelTwo); err != nil {
			t.Errorf("error: failed to link zettel")
		}
		if err := zettelOne.Link(zettelThree); err != nil {
			t.Errorf("error: failed to link zettel")
		}

		// Modify zettelTwo
		zettelTwo.Lines = lo.ReplaceAll(zettelTwo.Lines, zettelTwo.Lines[0], "# changed title")
		if err := zettelTwo.Write(); err != nil {
			t.Errorf("error: failed to write zettel")
		}
		if err := zettelTwo.Repair(c); err != nil {
			t.Errorf("error: failed to repair zettel")
		}

		link := fmt.Sprintf("- [%s](%s)", zettelTwo.Title, zettelTwo.Path)
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelOne.Lines, " "), []string{link})

		// Modify zettelThree
		zettelThree.Lines = lo.ReplaceAll(zettelThree.Lines, zettelThree.Lines[0], "# pretty different")
		if err := zettelThree.Write(); err != nil {
			t.Errorf("error: failed to write zettel")
		}
		if err := zettelThree.Repair(c); err != nil {
			t.Errorf("error: failed to repair zettel")
		}

		link = fmt.Sprintf("- [%s](%s)", zettelThree.Title, zettelThree.Path)
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelOne.Lines, " "), []string{link})
	})

	// cleanup
	err := os.RemoveAll("/tmp/foo")
	if err != nil {
		t.Errorf("error: failed to cleanup")
	}
}
