package tasks

import (
	"errors"
	"log"
	"runtime"
)

func Run(tasks []func() error, N int, M int) error {

	// Set number of concurrently working goroutines
	runtime.GOMAXPROCS(N)

	// Create working channels
	errChan := make(chan error, M)   // buffered error channel
	doneChan := make(chan string, N) // buffered success channel

	defer func() {
		close(errChan)
		close(doneChan)
		log.Println("Run is exited")
	}()

	// Running tasks
	for _, task := range tasks {
		go func(t func() error) {
			err := t()
			if err != nil {
				// Make error handling - send error in channel
				errChan <- err
				return
			}
			doneChan <- "done"
		}(task)
	}

	var (
		errCounter  = 0
		doneCounter = 0
	)

	// Reading messages from both channels
	for {
		select {
		case state := <-doneChan:
			// Checking success state from channel
			if state == "done" {
				doneCounter++
				if doneCounter == len(tasks) {
					log.Printf("%d tasks were executed successfully", len(tasks))
					return nil
				}
			}
		case err := <-errChan:
			// Read errors from channel
			if err != nil {
				log.Printf("Task return error: %s", err.Error())
				errCounter++
				if errCounter == M {
					return errors.New("error limit is elapsed")
				}
			}
		}
	}

}
