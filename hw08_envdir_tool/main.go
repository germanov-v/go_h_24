package main

import (
	"fmt"
	"os"
)

func main() {
	//env, _ := ReadDir("./testdata/env")
	//code := RunCmd([]string{"/bin/bash", "./testdata/echo.sh", "arg=1", "arg2=2"}, env)
	//
	//os.Exit(code)
	//fmt.Fprintln(os.Stderr, "arguments was not set")
	//fmt.Fprintln(os.Stdout, "arguments was not set")
	//fmt.Println(strings.SplitN("dat=val=test", "=", 3))
	// Place your code here.

	envDir := os.Args[1]
	cmdArgs := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading envdir:", err)
		os.Exit(1)
	}

	code := RunCmd(cmdArgs, env)

	os.Exit(code)
}
