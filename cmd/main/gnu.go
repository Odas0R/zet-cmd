package main

import (
	"errors"
	"os"
)

// Creates a file if it doesnt exist
func CreateFile(text string, filePath string) error {
	if fileExists := FileExists(filePath); fileExists {
		return nil
	}

	if err := Cat(text, filePath); err != nil {
		return err
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

// Creates a directory on the given path
func Mkdir(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
