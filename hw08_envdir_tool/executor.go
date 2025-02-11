package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	//	"strings"
)

type EnvironmentItems map[string]*string

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

		// empty
		if len(contentFile) == 0 {
			data[key] = nil // указатель
			continue
		}

		//  new => 0x0A
		contentFile = bytes.ReplaceAll(contentFile, []byte{0}, []byte("\n"))

		rowIndexes := bytes.IndexByte(contentFile, '\n')

		var lineFirst []byte

		if rowIndexes == -1 {
			lineFirst = contentFile
		} else {
			lineFirst = contentFile[:rowIndexes]
		}

		resultStr := strings.TrimRight(string(lineFirst), "\t") // табуляция
		data[key] = &resultStr

	}
	return data, nil
}

// пробуем сплтить и потом конкатинация
func createNewEnvironment(currentEnv []string, envExt Environment) []string {
	resultEnv := make([]string, len(currentEnv)+len(envExt))

	for _, valEnv := range currentEnv {
		partItems := strings.SplitN(valEnv, "=", 2)
		keyEnv := partItems[0]

		if _, ok := envExt[keyEnv]; ok {
			continue // в новом окружении, уже есть переменная
		}

		resultEnv = append(resultEnv, keyEnv)
	}

	for ker, value := range envExt {
		//if value != nil {
		//	resultEnv = append(resultEnv, ker)
		//}
		if !value.NeedRemove {
			resultEnv = append(resultEnv, ker+"="+value.Value)
		}
	}

	return resultEnv
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.
	returnCode = 1
	if len(cmd) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "arguments was not set")

		return
	}

	currentEnvironment := os.Environ()
	newEnv := createNewEnvironment(currentEnvironment, env)

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = newEnv
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	err := command.Run()
	if err != nil {

		// An ExitError reports an unsuccessful exit by a command.
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		_, _ = fmt.Fprintln(os.Stderr, err)

		return
	}

	returnCode = 0
	return
}
