package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

	t.Run("can fix a zettel", func(t *testing.T) {
		zettelOne := &Zettel{Title: "This is a foo title"}
		zettelOne.New(c)

		// modify the title
		lines := zettelOne.Lines
		lines[0] = "# foo bar"

		output := strings.Join(lines, "\n")
		err := ioutil.WriteFile("myfile", []byte(output), 0644)
		ExpectNoError(t, err, "error: failed to replace lines")

		// repair zettel
		err = zettelOne.Repair(c)
		ExpectNoError(t, err, "error: failed to repair zettel")
	})

	// cleanup
	err := os.RemoveAll("/tmp/foo")
	if err != nil {
		t.Errorf("error: failed to cleanup")
	}
}
