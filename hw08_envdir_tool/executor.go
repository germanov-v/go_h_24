package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	//	"strings"
)

type EnvironmentItems map[string]*string

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
