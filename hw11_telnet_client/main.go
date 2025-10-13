package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", time.Second*10, "connection timeout")
	flag.Parse()

	host := flag.Arg(0)
	port := flag.Arg(1)

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "connection error: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	done := make(chan struct{}, 2)

	go func() {
		defer close(done)
		if err := client.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "send error: %v\n", err)
		}
	}()

	go func() {
		defer close(done)
		if err := client.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "receive error: %v\n", err)
		}
	}()

	// Ждём либо сигнала прерывания, либо завершения одной из горутин
	select {
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "\ninterrupted")
	case <-done:
	}

	// Ждём завершения обеих горутин
	<-done
}
