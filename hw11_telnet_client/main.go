package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?

	timeout := flag.Duration("timeout", 5*time.Second, "timeout")

	flag.Parse()
	args := flag.Args()
	if len(args) <= 1 {
		fmt.Println("Usage: go run main.go [-timeout 1s] host port")
		os.Exit(1)
	}

	host := args[0]
	port := args[1]

	addr := net.JoinHostPort(host, port)

	client := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		//
		_, err := fmt.Fprintf(os.Stderr, "Error connet: %v\n", err)
		if err != nil {
			return
		}
		//fmt.Printf("Error conection: %v\n", err)
		os.Exit(1)
	}

	defer client.Close()

	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan, os.Interrupt)

	//done := make(chan bool, 1)
	//done := make(chan interface{}, 1)
	done := make(chan struct{}, 1)

	go func() {
		err := client.Receive()
		// Надо stderr тут
		if err != nil {
			fmt.Fprintf(os.Stderr, "connection was closed by: %v\n", err)

		} else {
			//fmt.Printf(os.Stderr, "Received SIGINT %v\n")
			fmt.Fprintln(os.Stderr, "end of line detected")
		}
		done <- struct{}{}
	}()

	go func() {
		err := client.Send()
		if err != nil {
			//fmt.Printf("Error sending: %v\n", err)
			//_, err := fmt.Fprintf(os.Stderr, "Received SIGINT %v\n", err)
			fmt.Fprintf(os.Stderr, "Error sending: %v\n", err)
		}
		//done <- true
		done <- struct{}{}
	}()

	select {
	case <-sChan:
		fmt.Fprintln(os.Stderr, "SIGINT, exit")
	case <-done:
	}

	<-done
}
