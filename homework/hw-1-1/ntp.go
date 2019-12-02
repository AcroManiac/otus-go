package main

import (
	"fmt"
	"syscall"
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
		fmt.Println("An error occurred:", err.Error())
		syscall.Exit(1)
	}
}
