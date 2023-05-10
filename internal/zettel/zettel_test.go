package zettel

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/odas0r/zet/cmd/assert"
)

func TestZettel(t *testing.T) {
	t.Run("can get zettel id from a path", func(t *testing.T) {
		result, ok := MatchSubstring(".", ".", "/home/odas0r/zet/something/else/using-tags-on-your-zettelkasten.782490819082442123123.md")
		assert.Equal(t, ok, true, "found a match")
		assert.Equal(t, result, "782490819082442123123", "result is the zettel id")

		_, ok = MatchSubstring(".", ".", "/home/odas0r/zet/something/else/using-tags-on-your-zettelkasten782490819082442.md")
		assert.Equal(t, ok, false, "didn't find a match")

		_, ok = MatchSubstring(".", ".", "/home/odas0r/zet/something/else/using-tags-on-your-zettelkasten.782490819082442md")
		assert.Equal(t, ok, false, "didn't find a match")

		result, _ = MatchSubstring("(", ")", "- [sfasfasfasf](/home/odas0r/zet/something/else/using-tags-on-your-zettelkasten.782490819082442.md)")
		assert.Equal(t, result, "/home/odas0r/zet/something/else/using-tags-on-your-zettelkasten.782490819082442.md", "result is the same as given path")
	})

	t.Run("can create a zettel", func(t *testing.T) {
		z1 := &Zettel{ID: 1, Title: "Title example"}
		z2 := &Zettel{ID: 2, Title: "Title example"}
		z3 := &Zettel{ID: 3, Title: "Title example"}

		ZettelNew(t, z1)
		ZettelNew(t, z2)
		ZettelNew(t, z3)

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
		ZettelNew(t, zettel)

		zettel = &Zettel{Path: fmt.Sprintf("%s/%s", config.Sub.Fleet, "title-example.999.md")}
		ZettelRead(t, zettel)

		assert.Equal(t, zettel.ID, int64(999), "zettel.ID should be correct")
		assert.Equal(t, zettel.Type, "fleet", "zettel.Type should be correct")
		assert.Equal(t, zettel.Slug, "title-example", "zettel.Slug should be correct")
		assert.Equal(t, zettel.Path, "/tmp/foo/fleet/title-example.999.md", "zettel.Path should be correct")
		assert.Equal(t, fmt.Sprintf("# %s", zettel.Title), zettel.Lines[0], "zettel.Title should be on line 0")
	})

	t.Run("can validate a link", func(t *testing.T) {
		z1 := &Zettel{ID: time.Now().UnixNano(), Title: "Title Example"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Title Example"}
		z3 := &Zettel{ID: 997, Title: "Title Example"}

		ZettelNew(t, z1)
		ZettelNew(t, z2)
		ZettelNew(t, z3)

		validLinks := []string{
			fmt.Sprintf("- [%s](%s)", z1.Title, z1.Path),
			fmt.Sprintf("- [%s](%s)", z2.Title, z2.Path),
			"- [Title Example](/tmp/foo/fleet/title-example.997.md)",
		}

		for index, link := range validLinks {
			path, ok := ValidateLinkPath(link)
			assert.Equal(t, strings.Contains(validLinks[index], path), true, "path should not be empty")
			assert.Equal(t, ok, true, "link should be valid")
		}

		invalidLinks := []string{
			fmt.Sprintf("- [%s](%s)", z2.Title, "/random/path"),
			fmt.Sprintf("- [%s](%s", z2.Title, z2.Path),
			fmt.Sprintf("- [%s]%s)", z2.Title, z2.Path),
			fmt.Sprintf("- [%s]((((%s)", z2.Title, z2.Path),
			fmt.Sprintf("- [%s](((%s)))))", z2.Title, z2.Path),
		}

		for _, link := range invalidLinks {
			path, ok := ValidateLinkPath(link)
			assert.Equal(t, path, "", "link path should be empty")
			assert.Equal(t, ok, false, "link should be invalid")
		}
	})

	t.Run("can link a zettel and read his links", func(t *testing.T) {
		z1 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z3 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}

		ZettelNew(t, z1)
		ZettelNew(t, z2)
		ZettelNew(t, z3)

		//
		// Link 1: z1 --> z2
		//
		ZettelLink(t, z1, z2)

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
		ZettelLink(t, z1, z3)

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

		assert.Equal(t, z1.ReadLinks(), nil, "should be able to read zettel links")

		assert.Equal(t, strings.Join(z1.Lines, "\n"), strings.Join(expectedLines, "\n"), "file should have correct format")
		assert.Equal(t, len(z1.Links), 2, "z1 has 2 links")
		assert.Equal(t, z1.Links[0].Path, z2.Path, "z2 is linked to z1")
		assert.Equal(t, z1.Links[1].Path, z3.Path, "z3 is linked to z1")

		//
		// Link 2: z3 --> z2
		//
		ZettelLink(t, z3, z2)

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

		ZettelNew(t, z1)
		ZettelNew(t, z2)

		err := z1.Link(z1)
		assert.Equal(t, err.Error(), "error: cannot link the same file", "zettel cannot link himself")

		ZettelLink(t, z1, z2)

		err = z1.Link(z2)
		assert.Equal(t, err.Error(), "error: cannot have duplicated links", "zettel cannot have duplicated links")
	})

	t.Run("can repair zettel links", func(t *testing.T) {
		z1 := &Zettel{ID: 998, Title: "Title Example"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Title Example"}
		z3 := &Zettel{ID: time.Now().UnixNano(), Title: "Title Example"}

		ZettelNew(t, z1)
		ZettelNew(t, z2)
		ZettelNew(t, z3)

		ZettelLink(t, z3, z1)
		ZettelLink(t, z2, z1)

		// repair the zettel 1
		assert.Equal(t, z1.WriteLine(0, "# foo bar"), nil, "should be able to write line on index 0")

		err, ok := z1.Repair()
		assert.Equal(t, err, nil, "should be able to write line on index 0")
		assert.Equal(t, ok, true, "the filepath of the zettel was changed")

		assert.Equal(t, z1.Title, "foo bar", "zettelOne.Title should be correct")
		assert.Equal(t, z1.Lines[0], "# foo bar", "zettelOne.Lines[0] should be the new title")
		assert.Equal(t, z1.Slug, "foo-bar", "zettelOne.Slug should be correct")
		assert.Equal(t, z1.Type, "fleet", "zettelOne.Type should be correct")
		assert.Equal(t, z1.FileName, "foo-bar.998.md", "zettelOne.FileName should be correct")
		assert.Equal(t, z1.Path, "/tmp/foo/fleet/foo-bar.998.md", "zettelOne.Path should be correct")

		assert.Equal(t, z2.ReadLines(), nil, "should be able to read lines")
		assert.Equal(t, z3.ReadLines(), nil, "should be able to read lines")

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

		ZettelNew(t, z1)
		ZettelNew(t, z2)
		ZettelNew(t, z3)

		ZettelLink(t, z2, z1)
		ZettelLink(t, z3, z1)

		assert.Equal(t, z1.Permanent(), nil, "zettel should be able to merge into permanent")

		assert.Equal(t, z1.Type, "permanent", "z1.Type should be correct")
		assert.Equal(t, z1.Path, fmt.Sprintf("%s/title-example.997.md", config.Sub.Permanent), "z1.Path should be correct")

		ZettelRead(t, z2)
		ZettelRead(t, z3)

    assert.Equal(t, strings.Contains(strings.Join(z2.Lines, " "), z1.Path), true, "z2 should have the permanet z1 path")
    assert.Equal(t, strings.Contains(strings.Join(z3.Lines, " "), z1.Path), true, "z3 should have the permanet z1 path")

		assert.Equal(t, len(z2.Links), 1, "z2 zettel should have one link")
		assert.Equal(t, len(z3.Links), 1, "z2 zettel should have one link")

		z2.Links[0].ReadMetadata()
		z3.Links[0].ReadMetadata()

		assert.Equal(t, z2.Links[0].FileName == "title-example.997.md", true, "z1 is linked to z2")
		assert.Equal(t, z2.Links[0].Type, "permanent", "z2 link is type permanent")

		assert.Equal(t, z3.Links[0].FileName == "title-example.997.md", true, "z1 is linked to z3")
		assert.Equal(t, z3.Links[0].Type, "permanent", "z3 link is type permanent")
	})

	t.Run("can delete zettel without links", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			zettel := &Zettel{ID: time.Now().UnixNano(), Title: faker.New().Person().Title()}
			ZettelNew(t, zettel)
			err := zettel.WriteLine(1,
				fmt.Sprintf("%s%d",
					strings.Join(faker.New().Lorem().Sentences(10), "\n"),
					faker.New().RandomDigit(),
				))
			assert.Equal(t, err, nil, "should be able to write a line on index 1")
		}

		z1 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z1.New()

		err := z1.Delete()
		assert.Equal(t, err, nil, "delete should not fail")

		ok := z1.IsValid()
		assert.Equal(t, ok, false, "zettel must be invalid")
	})

	t.Run("can delete zettel with links", func(t *testing.T) {
		z1 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z2 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}
		z3 := &Zettel{ID: time.Now().UnixNano(), Title: "Title"}

		z1.New()
		z2.New()
		z3.New()

		z2.Link(z1)
		z3.Link(z1)

		assert.Equal(t, len(z2.Links), 1, "z2 should have one links")
		assert.Equal(t, len(z3.Links), 1, "z3 should have one links")

		//
		// Delete zettel with links
		//
		err := z1.Delete()
		assert.Equal(t, err, nil, "delete should not fail")

		ok := z1.IsValid()
		assert.Equal(t, ok, false, "zettel is not valid")

		z2.ReadLinks()
		z3.ReadLinks()

		assert.Equal(t, len(z2.Links), 0, "z2 should have no links")
		assert.Equal(t, len(z3.Links), 0, "z3 should have no links")
	})
}

func ZettelNew(t *testing.T, z *Zettel) {
	assert.Equal(t, z.New(), nil, "should be able to create zettel")
}

func ZettelRead(t *testing.T, z *Zettel) {
	assert.Equal(t, z.Read(), nil, "should be able to read zettel")
}

func ZettelLink(t *testing.T, z *Zettel, z1 *Zettel) {
	assert.Equal(t, z.Link(z1), nil, "should be able to link zettels")
}

func LogFile(t *testing.T, z *Zettel) {
	t.Log(
		strings.Join(z.Lines, "\n"),
	)
}
