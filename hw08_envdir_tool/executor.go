package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type EnvironmentItems map[string]string

func loadEnvironmentItemsByDir(path string) (EnvironmentItems, error) {
	data := make(EnvironmentItems)

	items, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.IsDir() {
			continue
		}
		key := item.Name()
		pathFile := filepath.Join(path, item.Name())
		contentFile, err := os.ReadFile(pathFile)
		if err != nil {
			return nil, err
		}

		data[key] = string(contentFile)

	}
	return data, nil
}

// пробуем сплтить и потом конкатинация
func createNewEnvironment(currentEnv []string, envExt Environment) []string {
	result := make([]string, len(currentEnv)+len(currentEnv))

	return result
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.

	if len(cmd) == 0 {
		_, err := fmt.Fprintln(os.Stderr, "arguments was not set")
		if err != nil {
			return -1 // ???
		}
		returnCode = 1
		return
	}

	currentEnvironment := os.Environ()
	newEnvironment := createNewEnvironment(currentEnvironment, env)

	return returnCode
}
