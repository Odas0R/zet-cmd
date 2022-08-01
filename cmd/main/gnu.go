package main

import (
	"errors"
	"os"
)

// Appends text to a file, if it doesn't exist, create it
func Cat(text string, filepath string) error {
	fileExists := FileExists(filepath)
	if fileExists {
		return nil
	}

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
