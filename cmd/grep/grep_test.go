package grep

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/odas0r/zet/cmd/assert"
)

func TestGrep(t *testing.T) {
	// Create folder for testing
	t.Run("can grep recursively through directories", func(t *testing.T) {
		Mkdir("/tmp/test_grep")

		Cat("some text\n", "/tmp/test_grep/file.txt")
		Cat("another text\n", "/tmp/test_grep/file.txt")

		results, _ := Grep("another text", []string{"/tmp/test_grep"}, []*Result{})
		assert.Equal(t, len(results), 1, "len(result) should be 1")
		assert.Equal(t, results[0].Path, "/tmp/test_grep/file.txt", "first result must be on /tmp/test_grep/file.txt")
		assert.Equal(t, results[0].LineNr, 2, "first result must be on line 1")

		Mkdir("/tmp/test_grep/1")
		Mkdir("/tmp/test_grep/1/2")
		Mkdir("/tmp/test_grep/1/2/3")

		Cat("some text\n", "/tmp/test_grep/file.txt")
		Cat("some text\n", "/tmp/test_grep/1/file.txt")
		Cat("some text\n", "/tmp/test_grep/1/2/file.txt")
		Cat("some text\n", "/tmp/test_grep/1/2/3/file.txt")

		results, _ = Grep("some text", []string{"/tmp/test_grep"}, []*Result{})

		assert.Equal(t, len(results), 5, "len(result) should be 5")
		assert.Equal(t, results[3].Path, "/tmp/test_grep/file.txt", "result[3] should be file.txt")
		assert.Equal(t, results[2].Path, "/tmp/test_grep/1/file.txt", "result[2] should be /1/file.txt")
		assert.Equal(t, results[1].Path, "/tmp/test_grep/1/2/file.txt", "result[1] should be /1/2/file.txt")
		assert.Equal(t, results[0].Path, "/tmp/test_grep/1/2/3/file.txt", "result[0] should be /1/2/3/file.txt")

		Cat("key text\n", "/tmp/test_grep/file1.txt")

		results, _ = Grep("text", []string{"/tmp/test_grep"}, []*Result{})
		assert.Equal(t, len(results), 7, "len(result) should be 7")

		// Cleanup
		Remove("/tmp/test_grep/file.txt")
		Remove("/tmp/test_grep/file1.txt")
		Remove("/tmp/test_grep/1/file.txt")
		Remove("/tmp/test_grep/1/2/file.txt")
		Remove("/tmp/test_grep/1/2/3/file.txt")
	})
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

func Remove(path string) {
	if err := os.Remove(path); err != nil {
		log.Fatalf("error: failed to remove the file %v", err)
	}
}

func PrintLine(t *testing.T, path string) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}

	scanner := bufio.NewScanner(file)
	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}

	t.Log("\n" + strings.Join(lines, "\n"))
}
