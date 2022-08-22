package main

import (
	"testing"
	"time"

	"github.com/odas0r/zet/cmd/assert"
)

func TestHistory(t *testing.T) {
	t.Run("can append zettel into history", func(t *testing.T) {
		zettel := &Zettel{Title: "Some random title"}

		if err := zettel.New(); err != nil {
			assert.Equal(t, err, nil, "zettel.New should not fail")
		}

		if err := history.Insert(zettel); err != nil {
			assert.Equal(t, err, nil, "history.Insert should not fail")
		}

		if err := history.Read(); err != nil {
			assert.Equal(t, err, nil, "history.Read should not fail")
		}

		assert.Equal(t, history.Lines[0], zettel.Path, "zettel path must be in history file")
	})

	t.Run("can append zettel into history (2)", func(t *testing.T) {
		zettelOne := &Zettel{ID: 123, Title: "Some random title"}
		zettelTwo := &Zettel{ID: 1234, Title: "Some random title"}

		if err := zettelOne.New(); err != nil {
			assert.Equal(t, err, nil, "zettelOne.New should not fail")
		}
		if err := zettelTwo.New(); err != nil {
			assert.Equal(t, err, nil, "zettelTwo.New should not fail")
		}

		if err := history.Insert(zettelOne); err != nil {
			assert.Equal(t, err, nil, "history.Insert should not fail")
		}
		if err := history.Insert(zettelTwo); err != nil {
			assert.Equal(t, err, nil, "history.Insert should not fail")
		}

		if err := history.Read(); err != nil {
			assert.Equal(t, err, nil, "history.Read should not fail")
		}

		assert.Equal(t, history.Lines[1], zettelOne.Path, "zettelOne path must be in history file")
		assert.Equal(t, history.Lines[0], zettelTwo.Path, "zettelTwo path must be in history file")

		if err := history.Insert(zettelOne); err != nil {
			assert.Equal(t, err, nil, "history.Insert should not fail")
		}
		if err := history.Read(); err != nil {
			assert.Equal(t, err, nil, "history.Read should not fail")
		}

		assert.Equal(t, history.Lines[1], zettelTwo.Path, "zettelTwo path must be in history file")
		assert.Equal(t, history.Lines[0], zettelOne.Path, "zettelOne path must be in history file")
	})

	t.Run("can delete zettel from history", func(t *testing.T) {
		zettel := &Zettel{Path: "/tmp/foo/fleet/some-random-title.1234.md"}

		if err := zettel.Read(); err != nil {
			assert.Equal(t, err, nil, "zettel.Read should not fail")
		}

		if err := history.Delete(zettel); err != nil {
			assert.Equal(t, err, nil, "history.Delete should not fail")
		}

		if err := history.Read(); err != nil {
			assert.Equal(t, err, nil, "history.Read should not fail")
		}

		for _, line := range history.Lines {
			if line == zettel.Path {
				t.Errorf("zettel was not removed from history")
			}
		}
	})

	t.Run("history size must be of 50 zettels", func(t *testing.T) {
		// update history lines
		err := history.Read()
		assert.Equal(t, err, nil, "history.Read should not fail")

		maxZettels := 50 - len(history.Lines)

		for i := 0; i < maxZettels; i++ {
			id := time.Now().UTC().UnixNano()

			zettel := &Zettel{Title: "title", ID: id}

			if err := zettel.New(); err != nil {
				assert.Equal(t, err, nil, "zettel.New should not fail")
			}

			// add to history
			if err := history.Insert(zettel); err != nil {
				assert.Equal(t, err, nil, "history.Insert should not fail")
			}
		}

		id := time.Now().UTC().UnixNano()
		zettel := &Zettel{Title: "title", ID: id}

		if err := zettel.New(); err != nil {
			assert.Equal(t, err, nil, "zettel.New should not fail")
		}

		err = history.Insert(zettel)
		assert.Equal(t, err.Error(), "error: history cannot have more than 50 zettels", "history.Insert should fail")
	})

	t.Run("can clear history", func(t *testing.T) {
		if err := history.Clear(); err != nil {
			assert.Equal(t, err, nil, "history.Clear should not fail")
		}

		if err := history.Read(); err != nil {
			assert.Equal(t, err, nil, "history.Read should not fail")
		}

		assert.Equal(t, len(history.Lines), 0, "history should be empty")
	})
}
