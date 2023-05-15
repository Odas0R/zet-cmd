package fs

import (
	"testing"

	"github.com/odas0r/zet/pkg/assert"
)

func TestGrep(t *testing.T) {
	// Create folder for testing
	t.Run("can grep recursively through directories", func(t *testing.T) {
		Mkdir("/tmp/test_grep")

		WriteToFile("some text\n", "/tmp/test_grep/file.txt")
		WriteToFile("another text\n", "/tmp/test_grep/file.txt")

		results, _ := Grep("another text", []string{"/tmp/test_grep"}, []*Result{})
		assert.Equal(t, len(results), 1, "len(result) should be 1")
		assert.Equal(t, results[0].Path, "/tmp/test_grep/file.txt", "first result must be on /tmp/test_grep/file.txt")
		assert.Equal(t, results[0].LineNr, 2, "first result must be on line 1")

		Mkdir("/tmp/test_grep/1")
		Mkdir("/tmp/test_grep/1/2")
		Mkdir("/tmp/test_grep/1/2/3")

		WriteToFile("some text\n", "/tmp/test_grep/file.txt")
		WriteToFile("some text\n", "/tmp/test_grep/1/file.txt")
		WriteToFile("some text\n", "/tmp/test_grep/1/2/file.txt")
		WriteToFile("some text\n", "/tmp/test_grep/1/2/3/file.txt")

		results, _ = Grep("some text", []string{"/tmp/test_grep"}, []*Result{})

		assert.Equal(t, len(results), 5, "len(result) should be 5")
		assert.Equal(t, results[3].Path, "/tmp/test_grep/file.txt", "result[3] should be file.txt")
		assert.Equal(t, results[2].Path, "/tmp/test_grep/1/file.txt", "result[2] should be /1/file.txt")
		assert.Equal(t, results[1].Path, "/tmp/test_grep/1/2/file.txt", "result[1] should be /1/2/file.txt")
		assert.Equal(t, results[0].Path, "/tmp/test_grep/1/2/3/file.txt", "result[0] should be /1/2/3/file.txt")

		WriteToFile("key text\n", "/tmp/test_grep/file1.txt")

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
