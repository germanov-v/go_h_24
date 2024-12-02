package main

import (
	"fmt"
	"github.com/germanov-v/go_h_24/hw01_hello_otus/stringutil"
)

func main() {
	strStart := "Hello, OTUS!"
	strResult := stringutil.Reverse(strStart)
	fmt.Println(strResult)
}
