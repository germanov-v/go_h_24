package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here

	data := make(Environment)

	items, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.IsDir() {
			continue
		}
		key := item.Name()

		if strings.Contains(key, "=") {
			continue
		}

		pathFile := filepath.Join(dir, item.Name())
		contentFile, err := os.ReadFile(pathFile)
		if err != nil {
			return nil, err
		}

		// empty
		if len(contentFile) == 0 {
			data[key] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		//  new => 0x0A
		if bytes.Contains(contentFile, []byte("\n")) {
			rowIndexes := bytes.IndexByte(contentFile, '\n')

			var lineFirst []byte

			if rowIndexes == -1 {
				lineFirst = contentFile
			} else {
				lineFirst = contentFile[:rowIndexes]
			}
			contentFile = lineFirst
		} else {
			contentFile = bytes.ReplaceAll(contentFile, []byte{0}, []byte("\n"))
		}

		// strings.TrimSpace(string(lineFirst)) //
		//resultStr := strings.TrimRight(string(lineFirst), " \t") // табуляция
		resultStr := strings.TrimRight(string(contentFile), " \t")
		data[key] = EnvValue{Value: resultStr, NeedRemove: false}

	}
	return data, nil
}
