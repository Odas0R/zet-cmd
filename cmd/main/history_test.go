package main

import (
	"testing"
	"time"

	"github.com/odas0r/zet/cmd/assert"
)

func TestHistory(t *testing.T) {
	t.Run("can append zettel into history", func(t *testing.T) {
		zettel := &Zettel{Title: "title", ID: time.Now().UTC().UnixNano()}
		zettel.New()

		history.Insert(zettel)

		assert.Equal(t, history.Lines[0], zettel.Path, "zettel path must be in history file")

		history.Insert(zettel)
		assert.Equal(t, len(history.Lines), 1, "length of history must be equal to 1")
		assert.Equal(t, history.Lines[0], zettel.Path, "zettel path must be in history file")

		// Clean-Up
		history.Clear()
	})

	t.Run("can append zettel into history (2)", func(t *testing.T) {
		z1 := &Zettel{Title: "Some random title", ID: time.Now().UTC().UnixNano()}
		z2 := &Zettel{Title: "Some random title", ID: time.Now().UTC().UnixNano()}

		z1.New()
		z2.New()

		history.Insert(z1)
		history.Insert(z2)

		assert.Equal(t, history.Lines[1], z1.Path, "zettelOne path must be in history file")
		assert.Equal(t, history.Lines[0], z2.Path, "zettelTwo path must be in history file")

		history.Insert(z1)

		assert.Equal(t, len(history.Lines), 2, "zettelTwo path must be in history file")
		assert.Equal(t, history.Lines[1], z2.Path, "zettelTwo path must be in history file")
		assert.Equal(t, history.Lines[0], z1.Path, "zettelOne path must be in history file")

		// Clean-Up
		history.Clear()
	})

	t.Run("can delete zettel from history", func(t *testing.T) {
		z1 := &Zettel{Title: "Some random title", ID: time.Now().UTC().UnixNano()}
		z2 := &Zettel{Title: "Some random title", ID: time.Now().UTC().UnixNano()}
		z3 := &Zettel{Title: "Some random title", ID: time.Now().UTC().UnixNano()}

		z1.New()
		z2.New()
		z3.New()

		history.Insert(z1)
		history.Insert(z2)
		history.Insert(z3)

		history.Delete(z2)

		for _, line := range history.Lines {
			if line == z2.Path {
				t.Errorf("zettel was not removed from history")
			}
		}

		lines, _ := ReadLines(history.Path)

		for _, line := range lines {
			if line == z2.Path {
				t.Errorf("zettel was not removed from history")
			}
		}
	})

	t.Run("history size must be of 50 zettels", func(t *testing.T) {
		history.Clear()

		maxZettels := 50
		for i := 0; i < maxZettels; i++ {
			zettel := &Zettel{Title: "title", ID: time.Now().UTC().UnixNano()}
			zettel.New()
			history.Insert(zettel)
		}

		zettel := &Zettel{Title: "title", ID: time.Now().UTC().UnixNano()}
		zettel.New()

		err := history.Insert(zettel)
		assert.Equal(t, err.Error(), "error: history cannot have more than 50 zettels", "history.Insert should fail")
	})

	t.Run("can clear history", func(t *testing.T) {
		history.Clear()

		lines, _ := ReadLines(history.Path)

		assert.Equal(t, len(lines), 0, "history should be empty")
		assert.Equal(t, len(history.Lines), 0, "history should be empty")
	})
}
