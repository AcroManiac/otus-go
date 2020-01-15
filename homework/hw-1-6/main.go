package main

import (
	"errors"
	"github.com/AcroManiac/otus-go/homework/hw-1-6/tasks"
	"time"
)

func main() {
	_ = tasks.Run([]func() error{
		func() error {
			println("first")
			time.Sleep(time.Second)
			return errors.New("")
		},
		func() error {
			println("second")
			time.Sleep(time.Second)
			return errors.New("")
		},
		func() error {
			println("third")
			time.Sleep(time.Second)
			return errors.New("")
		},
	}, 4, 1)
}
