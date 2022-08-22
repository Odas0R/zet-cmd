package main

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/odas0r/zet/cmd/assert"
)

func TestEditor(t *testing.T) {
	t.Run("Grep, gives the correct ids", func(t *testing.T) {
		z1 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z3 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z4 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}

		z1.New()
		z2.New()
		z3.New()
		z4.New()

		z1.Link(z3)
		z2.Link(z3)
		z4.Link(z3)

		results, ok := Grep(strconv.FormatInt(z3.ID, 10))
		assert.Equal(t, ok, true, "grep returns values")

		outputStr := strings.Join([]string{results[0].Path, results[1].Path, results[2].Path}, " ")

		assert.Equal(t, strings.Contains(outputStr, z1.Path), true, "find-links should give z1.Path")
		assert.Equal(t, strings.Contains(outputStr, z2.Path), true, "find-links should give z2.Path")
		assert.Equal(t, strings.Contains(outputStr, z4.Path), true, "find-links should give z4.Path")
	})
}
