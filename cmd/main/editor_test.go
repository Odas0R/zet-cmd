package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/odas0r/zet/cmd/assert"
)

func TestEditor(t *testing.T) {
	t.Run("Grep, gives the correct ids", func(t *testing.T) {
		z1 := &Zettel{ID: time.Now().UnixNano(), Title: "Lorem ipsum dolor"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Lorem ipsum dolor sit"}
		z3 := &Zettel{ID: 12345, Title: "Lorem ipsum dolor sit amet,"}
		z4 := &Zettel{ID: time.Now().UnixNano(), Title: "Lorem ipsum dolor sit amet, consetetur"}

		z1.New()
		z2.New()
		z3.New()
		z4.New()

		fmt.Printf("z4.Lines: %v\n", z4.Lines)

		assert.Equal(t, z1.Link(z3), nil, "linking should not fail")
		assert.Equal(t, z2.Link(z3), nil, "linking should not fail")
		assert.Equal(t, z3.Link(z3), nil, "linking should not fail")

		results, ok := Grep("12345")
		assert.Equal(t, ok, true, "grep returns values")

		outputStr := strings.Join([]string{results[0].Path, results[1].Path, results[2].Path}, " ")

		assert.Equal(t, strings.Contains(outputStr, z1.Path), true, "find-links should give z1.Path")
		assert.Equal(t, strings.Contains(outputStr, z2.Path), true, "find-links should give z2.Path")
		assert.Equal(t, strings.Contains(outputStr, z4.Path), true, "find-links should give z4.Path")
	})
}
