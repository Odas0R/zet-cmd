package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/odas0r/zet/cmd/assert"
)

func TestZettel(t *testing.T) {
	t.Run("can create a zettel", func(t *testing.T) {
		z1 := &Zettel{ID: 1, Title: "Title example"}
		z2 := &Zettel{ID: 2, Title: "Title example"}
		z3 := &Zettel{ID: 3, Title: "Title example"}

		z1.New()
		z2.New()
		z3.New()

		paths := []string{
			"/tmp/foo/fleet/title-example.1.md",
			"/tmp/foo/fleet/title-example.2.md",
			"/tmp/foo/fleet/title-example.3.md",
		}

		for _, path := range paths {
			zettelExists := FileExists(path)
			assert.Equal(t, zettelExists, true, "zettel should exist")
		}
	})

	t.Run("can read a zettel on a given path", func(t *testing.T) {
		zettel := &Zettel{ID: 999, Title: "Title Example"}
		zettel.New()

		zettel = &Zettel{Path: fmt.Sprintf("%s/%s", config.Sub.Fleet, "title-example.999.md")}

		zettel.Read()

		assert.Equal(t, zettel.ID, int64(999), "zettel.ID should be correct")
		assert.Equal(t, zettel.Type, "fleet", "zettel.Type should be correct")
		assert.Equal(t, zettel.Slug, "title-example", "zettel.Slug should be correct")
		assert.Equal(t, zettel.Path, "/tmp/foo/fleet/title-example.999.md", "zettel.Path should be correct")
		assert.Equal(t, fmt.Sprintf("# %s", zettel.Title), zettel.Lines[0], "zettel.Title should be on line 0")
	})

	t.Run("can validate a link", func(t *testing.T) {
		z1 := &Zettel{ID: time.Now().UnixNano(), Title: "Title Example"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Title Example"}

		z1.New()
		z2.New()

		//
		// Validating links
		//

		link := fmt.Sprintf("- [%s](%s)", z1.Title, z1.Path)

		path, ok := ValidateLinkPath(link)
		assert.Equal(t, path, z1.Path, "ValidateLink recieved the correct zettel")
		assert.Equal(t, ok, true, "ValidateLink should be valid")

		link = fmt.Sprintf("- [%s](%s)", z2.Title, z2.Path)

		path, ok = ValidateLinkPath(link)
		assert.Equal(t, path, z2.Path, "link is of a valid zettel")
		assert.Equal(t, ok, true, "link should be valid")

		invalidLinks := []string{
			fmt.Sprintf("- [%s](%s)", z2.Title, "/random/path"),
			fmt.Sprintf("- [%s](%s", z2.Title, z2.Path),
			fmt.Sprintf("- [%s]%s)", z2.Title, z2.Path),
			fmt.Sprintf("- [%s]((((%s)", z2.Title, z2.Path),
			fmt.Sprintf("- [%s](((%s)))))", z2.Title, z2.Path),
		}

		for _, link := range invalidLinks {
			path, ok = ValidateLinkPath(link)
			assert.Equal(t, path, "", "link path should be empty")
			assert.Equal(t, ok, false, "link should be invalid")
		}
	})

	t.Run("can link a zettel and read his links", func(t *testing.T) {
		z1 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z3 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}

		z1.New()
		z2.New()
		z3.New()

		//
		// Link 1: z1 --> z2
		//
		z1.Link(z2)

		expectedLines := []string{
			"# Title",
			"",
			"## Bibliography",
			"",
			"## Links",
			"",
			fmt.Sprintf("- [%s](%s)", z2.Title, z2.Path),
			"",
		}
		assert.Equal(t, strings.Join(z1.Lines, "\n"), strings.Join(expectedLines, "\n"), "file should have correct format")
		assert.Equal(t, len(z1.Links), 1, "z1 has 2 links")
		assert.Equal(t, z1.Links[0].Path, z2.Path, "z2 is linked to z1")

		//
		// Link 2: z1 --> z3
		//
		z1.Link(z3)

		expectedLines = []string{
			"# Title",
			"",
			"## Bibliography",
			"",
			"## Links",
			"",
			fmt.Sprintf("- [%s](%s)", z2.Title, z2.Path),
			fmt.Sprintf("- [%s](%s)", z3.Title, z3.Path),
			"",
		}

		assert.Equal(t, strings.Join(z1.Lines, "\n"), strings.Join(expectedLines, "\n"), "file should have correct format")
		assert.Equal(t, len(z1.Links), 2, "z1 has 2 links")
		assert.Equal(t, z1.Links[0].Path, z2.Path, "z2 is linked to z1")
		assert.Equal(t, z1.Links[1].Path, z3.Path, "z3 is linked to z1")

		z1.ReadLinks()

		assert.Equal(t, strings.Join(z1.Lines, "\n"), strings.Join(expectedLines, "\n"), "file should have correct format")
		assert.Equal(t, len(z1.Links), 2, "z1 has 2 links")
		assert.Equal(t, z1.Links[0].Path, z2.Path, "z2 is linked to z1")
		assert.Equal(t, z1.Links[1].Path, z3.Path, "z3 is linked to z1")

		//
		// Link 2: z3 --> z2
		//
		z3.Link(z2)

		expectedLines = []string{
			"# Title",
			"",
			"## Bibliography",
			"",
			"## Links",
			"",
			fmt.Sprintf("- [%s](%s)", z2.Title, z2.Path),
			"",
		}

		assert.Equal(t, strings.Join(z3.Lines, "\n"), strings.Join(expectedLines, "\n"), "file should have correct format")
		assert.Equal(t, len(z3.Links), 1, "z1 has 2 links")
		assert.Equal(t, z3.Links[0].Path, z2.Path, "z2 is linked to z3")
	})

	t.Run("cant link same zettel twice", func(t *testing.T) {
		z1 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}

		z1.New()
		z2.New()

		err := z1.Link(z1)
		assert.Equal(t, err.Error(), "error: cannot link the same file", "zettel cannot link himself")

		z1.Link(z2)
		err = z1.Link(z2)
		assert.Equal(t, err.Error(), "error: cannot have duplicated links", "zettel cannot have duplicated links")
	})

	t.Run("can repair zettel links", func(t *testing.T) {
		z1 := &Zettel{ID: 998, Title: "Title Example"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Title Example"}
		z3 := &Zettel{ID: time.Now().UnixNano(), Title: "Title Example"}

		z1.New()
		z2.New()
		z3.New()

		z3.Link(z1)
		z2.Link(z1)

		// repair the zettel 1
		z1.WriteLine(0, "# foo bar")
		z1.Repair()

		assert.Equal(t, z1.Title, "foo bar", "zettelOne.Title should be correct")
		assert.Equal(t, z1.Lines[0], "# foo bar", "zettelOne.Lines[0] should be the new title")
		assert.Equal(t, z1.Slug, "foo-bar", "zettelOne.Slug should be correct")
		assert.Equal(t, z1.Type, "fleet", "zettelOne.Type should be correct")
		assert.Equal(t, z1.FileName, "foo-bar.998.md", "zettelOne.FileName should be correct")
		assert.Equal(t, z1.Path, "/tmp/foo/fleet/foo-bar.998.md", "zettelOne.Path should be correct")

		z2.ReadLines()
		z3.ReadLines()

		expectedLines := []string{
			"# Title Example",
			"",
			"## Bibliography",
			"",
			"## Links",
			"",
			"- [foo bar](/tmp/foo/fleet/foo-bar.998.md)",
		}
		assert.Equal(t, strings.Join(z2.Lines, "\n"), strings.Join(expectedLines, "\n"), "zettel z2 should have correct data")
		assert.Equal(t, strings.Join(z3.Lines, "\n"), strings.Join(expectedLines, "\n"), "zettel z3 should have correct data")

	})

	t.Run("can permanent zettel", func(t *testing.T) {
		z1 := &Zettel{ID: 997, Title: "Title Example"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Title Example"}
		z3 := &Zettel{ID: time.Now().UnixNano(), Title: "Title Example"}

		z1.New()
		z2.New()
		z3.New()

		z2.Link(z1)
		z3.Link(z1)

		z1.Permanent()

		assert.Equal(t, z1.Type, "permanent", "z1.Type should be correct")
		assert.Equal(t, z1.Path, fmt.Sprintf("%s/title-example.997.md", config.Sub.Permanent), "z1.Path should be correct")

		z2.Read()
		z3.Read()

		z2.Links[0].ReadMetadata()
		z3.Links[0].ReadMetadata()

		assert.Equal(t, z2.Links[0].FileName == "title-example.997.md", true, "z1 is linked to z2")
		assert.Equal(t, z2.Links[0].Type, "permanent", "z2 link is type permanent")

		assert.Equal(t, z3.Links[0].FileName == "title-example.997.md", true, "z1 is linked to z3")
		assert.Equal(t, z3.Links[0].Type, "permanent", "z3 link is type permanent")
	})

	t.Run("can delete zettel", func(t *testing.T) {
		z1 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z3 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z4 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}

		z1.New()
		z2.New()
		z3.New()
		z4.New()

		z2.Link(z1)
		z3.Link(z1)
		z4.Link(z2)

		assert.Equal(t, len(z2.Links), 1, "z2 should have one links")
		assert.Equal(t, len(z3.Links), 1, "z3 should have one links")
		assert.Equal(t, len(z4.Links), 1, "z4 should have one links")

		//
		// Delete
		//
		z1.Delete()

		ok := z1.IsValid()
		assert.Equal(t, ok, false, "zettel is not valid")

		z2.ReadLinks()
		z3.ReadLinks()
		z4.ReadLinks()

		fmt.Printf("strings.Join(z2.Lines, \"\n\"): %v\n", strings.Join(z2.Lines, "\n"))

		assert.Equal(t, len(z2.Links), 0, "z2 should have no links")
		assert.Equal(t, len(z3.Links), 0, "z3 should have no links")
		assert.Equal(t, len(z4.Links), 1, "z4 should have no links")
	})
}
