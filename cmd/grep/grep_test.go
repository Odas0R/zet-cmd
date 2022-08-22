package grep

import (
	"errors"
	"os"
	"testing"

	"github.com/odas0r/zet/cmd/assert"
)

func TestGrep(t *testing.T) {
	// Create folder for testing
	Mkdir("/tmp/test_grep")

	Cat("some text", "/tmp/test_grep/file.txt")
	Cat("another text", "/tmp/test_grep/file.txt")

	results, _ := Grep("another text", []string{"/tmp/test_grep"}, []*Result{})

	assert.Equal(t, results[0].FileName, "/tmp/test_grep/file.txt", "first result must be on /tmp/test_grep/file.txt")
	assert.Equal(t, results[0].LineNr, 1, "first result must be on line 1")

	Cat("some text", "/tmp/test_grep/file.txt")

	results, _ = Grep("some text", []string{"/tmp/test_grep"}, []*Result{})

	assert.Equal(t, len(results), 2, "Grep should find two results")
	assert.Equal(t, results[0].FileName, "/tmp/test_grep/file.txt", "Grep should find results")
	assert.Equal(t, results[0].LineNr, 1, "Grep should find results")
	assert.Equal(t, results[1].FileName, "/tmp/test_grep/file.txt", "Grep should find results")
	assert.Equal(t, results[1].LineNr, 2, "Grep should find results")

	// Cleanup
	os.RemoveAll("/tmp/test_grep")
}

func Mkdir(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func Cat(text string, filePath string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(text); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
