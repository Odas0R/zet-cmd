package main

import (
	"bufio"
	"log"
	"os"
)

type Zettel struct {
	ID       int64
	Name     string
	Type     string
	Title    string
	Path     string
	FileName string
	Tags     []string
	Links    []string
	Lines    []string
}

func (z Zettel) Create() (Zettel, error) {
	// config := &Config{}
	//
	// // parse the template
	// tmpl, err := template.ParseFiles(fmt.Sprintf("%s/zet.tmpl.md"))
	// if err != nil {
	// 	return z, err
	// }
	//
	// // create the zettel file
	// f, err := os.Create(z.Path)
	// if ez != nil {
	// 	return z, err
	// }
	//
	// // put the given title to the zettel
	// err = tmpl.Execute(f, z)
	// if err != nil {
	// 	return z, err
	// }
	// f.Close()
	//
	// // Set the lines of the file
	// z.Lines = getLines(z.Path)
	//
	return z, nil
}

func (z Zettel) Open() {
	// cmd := exec.Command(open, z.Path)
	// cmd.Start()
}

// -------------------- private methods -----------------------

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
