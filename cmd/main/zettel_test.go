package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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

	t.Run("can repair a zettel filename on title change", func(t *testing.T) {
		zettelOne := &Zettel{ID: 1224, Title: "This is a foo title"}
		zettelOne.New(c)

		zettelTwo := &Zettel{ID: 1241, Title: "This is a boo title"}
		zettelTwo.New(c)

		zettelThree := &Zettel{ID: 1243, Title: "This is a bye title"}
		zettelThree.New(c)

		zettelTwo.Link(zettelOne)
		zettelThree.Link(zettelOne)

		fmt.Println(zettelOne.Path, zettelTwo.Path, zettelThree.Path)

		// modify the title
		lines := zettelOne.Lines
		lines[0] = "# foo bar"
		linesWithEndLine := strings.Join(lines, "\n")
		ioutil.WriteFile(zettelOne.Path, []byte(linesWithEndLine), 0644)

		// repair zettel
		err := zettelOne.Repair(c)
		ExpectNoError(t, err, "error: failed to repair zettel")

		AssertStringEquals(t, "foo bar", zettelOne.Title)
		AssertStringEquals(t, "# foo bar", zettelOne.Lines[0])
		AssertStringEquals(t, "foo-bar", zettelOne.Slug)
		AssertStringEquals(t, "fleet", zettelOne.Type)
		AssertStringEquals(t, fmt.Sprintf("foo-bar.%d.md", zettelOne.ID), zettelOne.FileName)
		AssertStringEquals(t, fmt.Sprintf("%s/foo-bar.%d.md", c.Sub.Fleet, zettelOne.ID), zettelOne.Path)

		zettelTwo.Read()
		zettelThree.Read()

		// link := fmt.Sprintf("- [%s](%s)", zettelOne.Title, zettelOne.Path)
    // for _, line := range zettelTwo.Lines {
    //   fmt.Println(line)
    // }
    // for _, line := range zettelThree.Lines {
    //   fmt.Println(line)
    // }

		// Check if links are present on the files
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelTwo.Lines, " "), []string{zettelOne.Path})
		AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelThree.Lines, " "), []string{zettelOne.Path})

		// AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelTwo.Links, " "), []string{link})
		// AssertStringContainsSubstringsNoOrder(t, strings.Join(zettelThree.Links, " "), []string{link})
	})

	// cleanup
	// err := os.RemoveAll("/tmp/foo")
	// if err != nil {
	// 	t.Errorf("error: failed to cleanup")
	// }
}
