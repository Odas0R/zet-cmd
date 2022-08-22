package main

import (
	"strings"
	"testing"
	"time"

	"github.com/odas0r/zet/cmd/assert"
)

func TestEditor(t *testing.T) {
	t.Run("GrepLinksById, gives the correct ids", func(t *testing.T) {
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

		output, err := GrepLinksById(z3.ID)
		assert.Equal(t, err, nil, "GrepLinksById should not fail")

		path1 := strings.Split(output[0], ":")[1]
		path2 := strings.Split(output[1], ":")[1]
		path3 := strings.Split(output[2], ":")[1]

    outputStr := strings.Join([]string{path1, path2, path3}, " ")

		assert.Equal(t, strings.Contains(outputStr, z1.Path), true, "find-links should give z1.Path")
		assert.Equal(t,  strings.Contains(outputStr, z2.Path), true, "find-links should give z2.Path")
		assert.Equal(t, strings.Contains(outputStr, z4.Path), true, "find-links should give z4.Path")
	})
}
