package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	fmt.Fprintln(os.Stderr, "arguments was not set")
	fmt.Fprintln(os.Stdout, "arguments was not set")
	fmt.Println(strings.SplitN("dat=val=test", "=", 3))
	// Place your code here.
}
