package main

import (
	"errors"
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

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	t.conn = conn
	return nil
}

func (t *telnetClient) Send() error {
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(os.Stderr, "Connection closed by client (Ctrl+D)")
			return err
		}
		return fmt.Errorf("send message error: %w", err)
	}
	return nil
}

func (t *telnetClient) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("receive message error: %w", err)
	}
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		if err != nil {
			return fmt.Errorf("close connection error: %w", err)
		}
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
