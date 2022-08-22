package main

import (
	"bufio"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
)

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, fs.ErrNotExist)
}

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

func AppendLine(path, str string, index int) error {
	lines, err := ReadLines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}

	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}

func MatchSubstring(startS string, endS string, str string) (result string, found bool) {
	s := strings.Index(str, startS)
	if s == -1 {
		return result, false
	}
	newS := str[s+len(startS):]
	e := strings.Index(newS, endS)
	if e == -1 {
		return result, false
	}
	result = newS[:e]
	return result, true
}

