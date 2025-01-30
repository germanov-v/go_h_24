package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "tmp/1.txt", "file to read from")
	flag.StringVar(&to, "to", "tmp/1.txt", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	// Place your code here.

	if from == "" || to == "" {
		log.Fatal("empty from or to")
	}

	err := Copy(from, to, offset, limit)
	if err != nil {
		//panic
		log.Fatalf("copy error (main): %v\n", err)
	}

	fmt.Println("Success!")
}
