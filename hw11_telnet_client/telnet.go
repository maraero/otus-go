package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type client struct {
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	cancel  context.CancelFunc
}

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer, cancel context.CancelFunc) TelnetClient {
	return &client{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
		cancel:  cancel}
}

func (c *client) Connect() error {
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	c.conn = conn
	return nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Send() error {
	defer c.cancel()
	_, err := io.Copy(c.conn, c.in)
	if err != nil {
		return fmt.Errorf("sending error: %w", err)
	}
	fmt.Fprint(os.Stdout, "EOF\n")
	return nil
}

func (c *client) Receive() error {
	defer c.cancel()
	_, err := io.Copy(c.out, c.conn)
	if err != nil {
		return fmt.Errorf("receiving error: %w", err)
	}
	fmt.Fprintln(os.Stderr, "connection closed")
	return nil
}
