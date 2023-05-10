package zettel

import (
	"testing"
	"time"

	"github.com/odas0r/zet/cmd/assert"
)

func TestHistory(t *testing.T) {
	t.Run("can append zettel into history", func(t *testing.T) {
		z := &Zettel{Title: "title", ID: time.Now().UTC().UnixNano()}
    ZettelNew(t, z)

		HistoryInsert(t, z)

		assert.Equal(t, history.Lines[0], z.Path, "zettel path must be in history file")

		HistoryInsert(t, z)

		assert.Equal(t, len(history.Lines), 1, "length of history must be equal to 1")
		assert.Equal(t, history.Lines[0], z.Path, "zettel path must be in history file")

		// Clean-Up
    HistoryClear(t)
	})

	t.Run("can append zettel into history (2)", func(t *testing.T) {
		z1 := &Zettel{Title: "Some random title", ID: time.Now().UTC().UnixNano()}
		z2 := &Zettel{Title: "Some random title", ID: time.Now().UTC().UnixNano()}

    ZettelNew(t, z1)
    ZettelNew(t, z2)

    HistoryInsert(t, z1)
    HistoryInsert(t, z2)

		assert.Equal(t, history.Lines[1], z1.Path, "zettelOne path must be in history file")
		assert.Equal(t, history.Lines[0], z2.Path, "zettelTwo path must be in history file")

    HistoryInsert(t, z1)

		assert.Equal(t, len(history.Lines), 2, "zettelTwo path must be in history file")
		assert.Equal(t, history.Lines[1], z2.Path, "zettelTwo path must be in history file")
		assert.Equal(t, history.Lines[0], z1.Path, "zettelOne path must be in history file")

		// Clean-Up
    HistoryClear(t)
	})

	t.Run("can delete zettel from history", func(t *testing.T) {
		z1 := &Zettel{Title: "Some random title", ID: time.Now().UTC().UnixNano()}
		z2 := &Zettel{Title: "Some random title", ID: time.Now().UTC().UnixNano()}
		z3 := &Zettel{Title: "Some random title", ID: time.Now().UTC().UnixNano()}

		ZettelNew(t, z1)
		ZettelNew(t, z2)
		ZettelNew(t, z3)

		HistoryInsert(t, z1)
		HistoryInsert(t, z2)
		HistoryInsert(t, z3)

    HistoryDelete(t, z2)

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
		HistoryClear(t)

		maxZettels := 50
		for i := 0; i < maxZettels; i++ {
			zettel := &Zettel{Title: "title", ID: time.Now().UTC().UnixNano()}
			ZettelNew(t, zettel)
			HistoryInsert(t, zettel)
		}

		zettel := &Zettel{Title: "title", ID: time.Now().UTC().UnixNano()}
		ZettelNew(t, zettel)

		err := history.Insert(zettel)
		assert.Equal(t, err.Error(), "error: history cannot have more than 50 zettels", "history.Insert should fail")
	})

	t.Run("can clear history", func(t *testing.T) {
		HistoryClear(t)

		lines, _ := ReadLines(history.Path)

		assert.Equal(t, len(lines), 0, "history should be empty")
		assert.Equal(t, len(history.Lines), 0, "history should be empty")
	})
}

func HistoryInsert(t *testing.T, z *Zettel) {
	assert.Equal(t, history.Insert(z), nil, "should be able to insert zettel into history")
}

func HistoryDelete(t *testing.T, z *Zettel) {
	assert.Equal(t, history.Delete(z), nil, "should be able to delete zettel")
}

func HistoryClear(t *testing.T) {
	assert.Equal(t, history.Clear(), nil, "should be able to clear history")
}
