package main

import (
	"fmt"
	hw04lrucache "github.com/go_h_24/hw04_lru_cache"
)

func main() {
	list := hw04lrucache.NewList()
	currentLen := list.Len()
	fmt.Printf("Length: %d\n", currentLen)
}
