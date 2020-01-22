package main

import (
	"flag"
	"github.com/AcroManiac/otus-go/homework/hw-2-2/goenvdir"
	"log"
	"os"
)

// To test main() run the command:
// go run main.go ./dir printenv USER CITY PLANET
func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		log.Fatalln("Not enough arguments")
	}

	args := flag.Args()
	env, err := goenvdir.ReadDir(args[0])
	if err != nil {
		log.Fatalf("Error reading environment directory: %s", err.Error())
	}

	code := goenvdir.RunCmd(args[1:], env)
	if code != 0 {
		log.Printf("Error running external program. Exit code: %d", code)
	}
	os.Exit(code)
}
