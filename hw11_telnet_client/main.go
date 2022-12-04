package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	host, port, timeout := parseArgs()
	addr := net.JoinHostPort(host, port)

	ctx, cancel := context.WithCancel(context.Background())

	clientParams := ClientParams{addr: addr, timeout: *timeout, in: os.Stdin, out: os.Stdout, cancel: cancel}
	client := NewTelnetClient(clientParams)
	if err := client.Connect(); err != nil {
		log.Fatalf("failed to connect to %v: %v", addr, err)
	}
	defer client.Close()

	go receive(client)
	go send(client)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	select {
	case <-sigCh:
		cancel()
	case <-ctx.Done():
		close(sigCh)
	}
}

func send(client TelnetClient) {
	if err := client.Send(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to send: %v", err)
		return
	}
}

func receive(client TelnetClient) {
	if err := client.Receive(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to receive: %v", err)
		return
	}
}

func parseArgs() (host string, port string, timeout *time.Duration) {
	timeout = flag.Duration("timeout", time.Second*10, "timeout")
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("not enough arguments, pass host and port")
	}
	host = flag.Arg(0)
	port = flag.Arg(1)
	return host, port, timeout
}
