package main

import (
	"fmt"
	"github.com/germanov-v/go_h_24/hw01_hello_otus/stringutil"
	"golang.org/x/example/hello/reverse"
	"log"
)

func main() {
	// Place your code here.
	strStart := "Hello, OTUS!"
	strResultExReverse := reverse.String(strStart)
	strResultStrutil := stringutil.Reverse(strStart)

	if strResultExReverse != strResultStrutil {
		log.Fatalf("String is not equal: %s != %s\n", strResultExReverse, strResultStrutil)
	} else {
		fmt.Println(strResultExReverse)
	}

}
