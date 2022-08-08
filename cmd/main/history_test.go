package main

import (
	"os"
	"strings"
	"testing"

	"github.com/samber/lo"
)

func TestHistory(t *testing.T) {
	h := &History{Root: "/tmp/foo"}
	c := &Config{Root: "/tmp/foo"}
	if err := c.Init(); err != nil {
		t.Errorf("error: failed to initialize config %V", err)
	}

	t.Run("can initialize history", func(t *testing.T) {
		if err := h.Init(); err != nil {
			t.Errorf("error: failed to initialize history %V", err)
		}
	})

	t.Run("can append zettel into history", func(t *testing.T) {
		zettel := &Zettel{Title: "Some random title"}
		if err := zettel.New(c); err != nil {
			t.Errorf("error: failed to create zettel %V", err)
		}

		if err := h.Insert(zettel.Path); err != nil {
			t.Errorf("error: failed to insert zettel on history %V", err)
		}

		lines, _ := ReadLines(h.Path)
		linesJoined := strings.Join(lines, " ")

		AssertStringContainsSubstringsNoOrder(t, linesJoined, []string{zettel.Path})
	})

	t.Run("can append zettel into history (2)", func(t *testing.T) {
		zettelOne := &Zettel{ID: 123, Title: "Some random title"}
		if err := zettelOne.New(c); err != nil {
			t.Errorf("error: failed to create zettel %V", err)
		}

		zettelTwo := &Zettel{ID: 1234, Title: "Some random title"}
		if err := zettelTwo.New(c); err != nil {
			t.Errorf("error: failed to create zettel %V", err)
		}

		if err := h.Insert(zettelOne.Path); err != nil {
			t.Errorf("error: failed to insert zettel on history %V", err)
		}
		if err := h.Insert(zettelTwo.Path); err != nil {
			t.Errorf("error: failed to insert zettel on history %V", err)
		}

		lines, _ := ReadLines(h.Path)
		linesJoined := strings.Join(lines, " ")

		AssertStringContainsSubstringsInOrder(t, linesJoined, []string{zettelOne.Path, zettelTwo.Path})

		if err := h.Insert(zettelOne.Path); err != nil {
			t.Errorf("error: failed to insert zettel on history %V", err)
		}

		lines, _ = ReadLines(h.Path)
		linesJoined = strings.Join(lines, " ")

		AssertStringContainsSubstringsInOrder(t, linesJoined, []string{zettelTwo.Path, zettelOne.Path})
	})

	t.Run("can delete zettel from history", func(t *testing.T) {
		zettel := &Zettel{Path: "/tmp/foo/fleet/some-random-title.1234.md"}

		if err := zettel.Read(c); err != nil {
			t.Errorf("error: failed to read zettel %V", err)
		}

		if err := h.Delete(zettel.Path); err != nil {
			t.Errorf("error: failed to insert zettel on history %V", err)
		}

		lines, _ := ReadLines(h.Path)

		if hasLine := lo.Contains(lines, zettel.Path); hasLine {
			t.Errorf("error: path was not removed from history")
		}
	})

	// cleanup
	err := os.RemoveAll("/tmp/foo")
	if err != nil {
		t.Errorf("failed to cleanup")
	}
}
