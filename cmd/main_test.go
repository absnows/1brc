package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	files := find("./../data/test_cases", ".txt")
	for _, test := range files {
		t.Run(test, func(t *testing.T) {
			output := evaluate(test + ".txt")
			assert.Equal(t, readFile(test+".out"), "{"+output+"}\n")
		})
	}
}

func readFile(input string) string {
	fileContent, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}

	return string(fileContent)
}

func find(root, extenstion string) []string {
	var a []string
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(d.Name()) == extenstion {
			a = append(a, path[:len(path)-len(extenstion)])
		}

		return nil
	})

	return a
}
