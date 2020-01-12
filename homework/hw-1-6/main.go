package main

import (
	"github.com/AcroManiac/otus-go/homework/hw-1-6/tasks"
	"time"
)

func main() {
	tasks.Run([]func() error{
		func() error {
			println("first")
			time.Sleep(time.Second)
			return nil
		},
		func() error {
			println("second")
			time.Sleep(time.Second)
			return nil
		},
		func() error {
			println("third")
			time.Sleep(time.Second)
			return nil
		},
	}, 1, 1)
}
