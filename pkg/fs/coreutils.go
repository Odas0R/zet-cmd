package fs

import (
	"bufio"
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

const (
	DefaultPerms = 0600
)

func Touch(path string) error {
	myfile, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := myfile.Close(); err != nil {
		return err
	}

	return nil
}

// CatContent retrieves content from a file
func CatContent(filePath string) (string, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// WriteToFile writes content to a file.
// It will append to the file if it already exists and create it if it doesn't.
func WriteToFile(text string, path string) error {
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

// Cat returns all lines from a file
func Cat(filePath string) ([]string, error) {
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

func Mkdir(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func Remove(path string) {
	if err := os.Remove(path); err != nil {
		log.Fatalf("error: failed to remove the file %v", err)
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

func InsertLine(path, newLine string, index int) error {
	lines, err := Cat(path)
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

	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}
