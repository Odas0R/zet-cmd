package fs

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Result struct {
	Path   string
	LineNr int
}

func findPattern(file *os.File, pattern string) ([]*Result, bool) {
	results := []*Result{}
	lineIdx := 1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			result := &Result{
				Path:   file.Name(),
				LineNr: lineIdx,
			}
			results = append(results, result)
		}
		lineIdx++
	}

	if err := file.Close(); err != nil {
		log.Fatalf("error: failed to close file socket %v",
			err)
	}

	return results, len(results) > 0
}

func find(pattern string, filename string) ([]*Result, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return []*Result{}, err
	}

	results, ok := findPattern(file, pattern)
	if ok {
		return results, nil
	}

	return []*Result{}, nil
}

func Grep(pattern string, dirnames []string, results []*Result) ([]*Result, error) {
	for _, dirname := range dirnames {
		files, err := ioutil.ReadDir(dirname)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		for _, file := range files {
			if file.IsDir() {
				r, err := Grep(pattern, []string{dirname + "/" + file.Name()}, results)
				if err != nil {
					log.Fatalf("error: %v", err)
				}

				results = append(results, r...)
			}

			filePath := dirname + "/" + file.Name()
			found, err := find(pattern, filePath)
			if err != nil {
				log.Fatalf("error: %v", err)
			}

			results = append(results, found...)
		}
	}

	return results, nil
}
