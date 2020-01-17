package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/AcroManiac/otus-go/homework/hw-2-1/gocopy"
)

var (
	flagTo     string
	flagFrom   string
	flagLimit  int
	flagOffset int
)

func init() {
	flag.StringVar(&flagTo, "to", "", "destination file")
	flag.StringVar(&flagFrom, "from", "", "source file")
	flag.IntVar(&flagLimit, "limit", -1, "copying bytes number limit")
	flag.IntVar(&flagOffset, "offset", 0, "bytes offset from file start")
}

var Usage = func() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

// To test main() run the command:
// go run main.go -from ./gocopy/gocopy.go -to ./copied.txt limit -1 offset 0
func main() {
	flag.Usage = Usage
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := gocopy.Copy(flagFrom, flagTo, flagLimit, flagOffset); err != nil {
		log.Fatalf("An error occurred while file copying: %s", err.Error())
	}

	log.Println("File copied successfully")
}
