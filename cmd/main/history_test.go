package main

import (
	"testing"

	"github.com/odas0r/zet/cmd/assert"
)

func TestHistory(t *testing.T) {
	t.Run("can append zettel into history", func(t *testing.T) {
		zettel := &Zettel{Title: "Some random title"}

		err := zettel.New()
		assert.Equal(t, err, nil, "zettel.New should not fail")

		history.Insert(zettel.Path)
		assert.Equal(t, err, nil, "history.Insert should not fail")

		err = history.Read()
		assert.Equal(t, err, nil, "history.Read should not fail")
		assert.Equal(t, history.Lines[0], zettel.Path, "zettel path must be in history file")
	})

	t.Run("can append zettel into history (2)", func(t *testing.T) {
		zettelOne := &Zettel{ID: 123, Title: "Some random title"}
		zettelTwo := &Zettel{ID: 1234, Title: "Some random title"}

		err := zettelOne.New()
		assert.Equal(t, err, nil, "zettelOne.New should not fail")

		err = zettelTwo.New()
		assert.Equal(t, err, nil, "zettelTwo.New should not fail")

		err = history.Insert(zettelOne.Path)
		assert.Equal(t, err, nil, "history.Insert should not fail")

		err = history.Insert(zettelTwo.Path)
		assert.Equal(t, err, nil, "history.Insert should not fail")

		err = history.Read()
		assert.Equal(t, err, nil, "history.Read should not fail")
		assert.Equal(t, history.Lines[1], zettelOne.Path, "zettelOne path must be in history file")
		assert.Equal(t, history.Lines[2], zettelTwo.Path, "zettelTwo path must be in history file")

		err = history.Insert(zettelOne.Path)
		assert.Equal(t, err, nil, "history.Insert should not fail")

		err = history.Read()
		assert.Equal(t, err, nil, "history.Read should not fail")
		assert.Equal(t, history.Lines[1], zettelTwo.Path, "zettelTwo path must be in history file")
		assert.Equal(t, history.Lines[2], zettelOne.Path, "zettelOne path must be in history file")
	})

	t.Run("can delete zettel from history", func(t *testing.T) {
		zettel := &Zettel{Path: "/tmp/foo/fleet/some-random-title.1234.md"}

		err := zettel.Read()
		assert.Equal(t, err, nil, "zettel.Read should not fail")

		err = history.Delete(zettel.Path)
		assert.Equal(t, err, nil, "history.Delete should not fail")

		err = history.Read()
		assert.Equal(t, err, nil, "history.Read should not fail")

		for _, line := range history.Lines {
			if line == zettel.Path {
				t.Errorf("zettel was not removed from history")
			}
		}
	})
}
