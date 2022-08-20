package main

import (
	"os"
	"testing"

	"github.com/odas0r/zet/cmd/assert"
)

func TestHistory(t *testing.T) {
	history := &History{Root: "/tmp/foo"}
	config := &Config{Root: "/tmp/foo"}

	t.Run("can initialize history", func(t *testing.T) {
		err := config.Init()
		assert.Equal(t, err, nil, "config.Init should not fail")

		err = history.Init()
		assert.Equal(t, err, nil, "history.Init should not fail")
	})

	t.Run("can append zettel into history", func(t *testing.T) {
		zettel := &Zettel{Title: "Some random title"}

		err := zettel.New(config)
		assert.Equal(t, err, nil, "zettel.New should not fail")

		history.Insert(zettel.Path)
		assert.Equal(t, err, nil, "history.Insert should not fail")

		lines, _ := ReadLines(history.Path)
		assert.Equal(t, lines[0], zettel.Path, "zettel path must be in history file")
	})

	t.Run("can append zettel into history (2)", func(t *testing.T) {
		zettelOne := &Zettel{ID: 123, Title: "Some random title"}
		err := zettelOne.New(config)
		assert.Equal(t, err, nil, "zettelOne.New should not fail")

		zettelTwo := &Zettel{ID: 1234, Title: "Some random title"}
		err = zettelTwo.New(config)
		assert.Equal(t, err, nil, "zettelTwo.New should not fail")

		err = history.Insert(zettelOne.Path)
		assert.Equal(t, err, nil, "history.Insert should not fail")

		err = history.Insert(zettelTwo.Path)
		assert.Equal(t, err, nil, "history.Insert should not fail")

		lines, _ := ReadLines(history.Path)
		assert.Equal(t, lines[1], zettelOne.Path, "zettelOne path must be in history file")
		assert.Equal(t, lines[2], zettelTwo.Path, "zettelTwo path must be in history file")

		err = history.Insert(zettelOne.Path)
		assert.Equal(t, err, nil, "history.Insert should not fail")

		lines, _ = ReadLines(history.Path)
		assert.Equal(t, lines[1], zettelTwo.Path, "zettelTwo path must be in history file")
		assert.Equal(t, lines[2], zettelOne.Path, "zettelOne path must be in history file")
	})

	t.Run("can delete zettel from history", func(t *testing.T) {
		zettel := &Zettel{Path: "/tmp/foo/fleet/some-random-title.1234.md"}

		err := zettel.Read(config)
		assert.Equal(t, err, nil, "zettel.Read should not fail")

		err = history.Delete(zettel.Path)
		assert.Equal(t, err, nil, "history.Delete should not fail")

		lines, _ := ReadLines(history.Path)
		for _, line := range lines {
			if line == zettel.Path {
				t.Errorf("zettel was not removed from history")
			}
		}
	})

	// cleanup
	err := os.RemoveAll("/tmp/foo")
	assert.Equal(t, err, nil, "os.RemoveAll should not fail")
}
