package main

import (
	"flag"
	"log"

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

// To test main() run the command:
// go run main.go -from=./gocopy/gocopy.go -to=./copied.txt limit=-1 offset=0
func main() {
	flag.Parse()

	if err := gocopy.Copy(flagFrom, flagTo, flagLimit, flagOffset); err != nil {
		log.Fatalf("An error occurred while file copying: %s", err.Error())
	}

	log.Println("File copied successfully")
}
