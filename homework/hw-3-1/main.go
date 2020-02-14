package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
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
	ctx, _ /*cancel*/ = context.WithTimeout(ctx, ctxTimeout)

	// Create telnet client
	client := gotelnet.NewTelnetClient(host, port)
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
}
