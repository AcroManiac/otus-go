package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AcroManiac/otus-go/homework/hw-3-1/gotelnet"
)

var timeout string

func init() {
	flag.StringVar(&timeout, "timeout", "10s", "working timeout")
}

var Usage = func() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

// Test program with command
// go run main.go -timeout=10m 192.168.1.1 23
func main() {
	flag.Usage = Usage
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatalln("Specify host and port arguments")
	}
	host := flag.Arg(0)
	port := flag.Arg(1)

	// Convert timeout from string
	ctxTimeout, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatalf("Wring timeout: %s", timeout)
	}

	// Make running context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)

	// Create telnet client
	client := gotelnet.NewTelnetClient(host, port)
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Error connecting: %v", err)
	}

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := client.Receive(ctx); err != nil {
			log.Println(err)
			cancel()
		}
	}()
	go func() {
		if err := client.Send(ctx); err != nil {
			log.Println(err)
			cancel()
		}
	}()

	// Wait for interruption events
	select {
	case <-ctx.Done():
		log.Println("Program exited")
	case <-done:
		cancel()
		log.Println("User interrupted program. Bye!")
	}

	// Make connection shutdown
	if err := client.Close(); err != nil {
		log.Printf("Error while closing connection: %v", err)
	}
}
