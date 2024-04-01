package fs

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultPerms = 0600
)

func Create(path string) error {
	myfile, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := myfile.Close(); err != nil {
		return err
	}

	return nil
}

// Read retrieves content from a file
func Read(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ReadLines returns all lines from a file
func ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	// close file descriptor
	if err := file.Close(); err != nil {
		return nil, err
	}

	return lines, nil
}

// Write writes content to a file. It will append to the file if it already
// exists and create it if it doesn't.
func Write(path string, text string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func Mkdir(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// List returns all files in a directory, same as `ls`
func List(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var paths []string
	for _, file := range files {
		fullPath := filepath.Join(dir, file.Name())
		paths = append(paths, fullPath)
	}
	return paths
}

func Remove(path string) error {
	return os.Remove(path)
}

// rm -r
func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

func InsertLine(path, newLine string) error {
	lines, err := ReadLines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for _, line := range lines {
		fileContent += line
		fileContent += "\n"
	}

	fileContent += newLine
	fileContent += "\n"

	return os.WriteFile(path, []byte(fileContent), 0644)
}

func InsertLineAtIndex(path, newLine string, index int) error {
	lines, err := ReadLines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += newLine
			fileContent += "\n"
		}
		fileContent += line
		fileContent += "\n"
	}

	return os.WriteFile(path, []byte(fileContent), 0644)
}

func Input(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return text
}

func InputConfirm(prompt string) bool {
	answer := Input(prompt + " (y/n): ")
	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer == "y" || answer == "yes" {
		return true
	}
	return false
}

// Open opens a file with the default open command from the user system
func Open(path string) error {
	cmd := fmt.Sprintf("open %s", path)

	if err := Exec(cmd); err != nil {
		return err
	}

	return nil
}

// Editor opens a file with the default $EDITOR from the user system
func Editor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return errors.New("error: $EDITOR is not set")
	}

	cmd := fmt.Sprintf("%s %s", editor, path)

	if err := Exec(cmd); err != nil {
		return err
	}

	return nil
}
