package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func Time() (time.Time, error) {
	return ntp.Time("0.beevik-ntp.pool.ntp.org")
}

func main() {
	if tm, err := Time(); nil == err {
		fmt.Println("Current time is", tm)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "An error occurred: %s\n", err.Error())
		os.Exit(1)
	}
}
