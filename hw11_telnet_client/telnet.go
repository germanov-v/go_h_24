package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

// все из параметров
type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	//currentErrStd io.
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.

	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}

	//	return nil
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.

func (t *telnetClient) Connect() error {
	//conn, err := net.Dial("tcp", t.address)
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = conn

	//	fmt.Fprintf(os.Stdout, "Success connect to %s\n", t.address)
	fmt.Fprintf(os.Stderr, "Success connect to %s\n", t.address)
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

func (t *telnetClient) Send() error {
	//_, err := io.WriteString(t.conn, "\r\n")
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *telnetClient) Receive() error {
	//_,err := io.Copy(t.conn, t.in)
	//return err
	_, err := io.Copy(t.out, t.conn)
	//
	//if err != nil {
	//	return err
	//}

	//if err == io.EOF {
	//	return nil
	//}

	return err
}
