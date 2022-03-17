package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

type Zettel struct {
	ID    int64
	Name  string
	Type  string
	Title string
	Tags  []string
	Links []string
	Lines []string
}

func New(filePath string) *Zettel {
	// zettelType := filepath.Dir(filePath)
  lines := getLines(filePath)

	return &Zettel{
		Name:  filePath,
		Lines: lines,
	}
}

func NewFromType(type "fleet" | "permanent") *Zettel {
	zettelType := filepath.Dir(filePath)
  lines := getLines(filePath)

	return &Zettel{
		Name:  filePath,
		Lines: lines,
	}
}


func getLines(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}
