package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type ClientParams struct {
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	cancel  context.CancelFunc
}

type client struct {
	ClientParams
	conn net.Conn
}

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(c ClientParams) TelnetClient {
	return &client{ClientParams: c, conn: nil}
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
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return fmt.Errorf("sending error: %w", err)
	}
	return nil
}

func (c *client) Receive() error {
	defer c.cancel()
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("receiving error: %w", err)
	}
	return nil
}
