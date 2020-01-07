package tasks

import (
	"errors"
	"log"
	"runtime"
	"sync"
)

var abortChan chan struct{}

func cancelled() bool {
	select {
	case <-abortChan:
		return true
	default:
		return false
	}
}

func Run(tasks []func() error, N int, M int) error {

	// Set number of concurrently working goroutines
	runtime.GOMAXPROCS(N)

	// Create working channels
	// (make channels buffered to prevent goroutines leaking)
	errChan := make(chan error, M)   // error channel
	doneChan := make(chan string, N) // success channel

	// Create channel to interrupt on errors
	abortChan = make(chan struct{})

	// Use wait group to synchronize writes to channels on function exit
	var wg sync.WaitGroup

	defer func() {
		wg.Wait() // Wait until all tasks exit
		close(errChan)
		close(doneChan)
		log.Println("Run() function is exited")
	}()

	// Running tasks in separate goroutines
	for _, task := range tasks {
		wg.Add(1) // increment wait group counter
		go func(t func() error) {
			defer wg.Done() // decrement wait group counter on goroutine exit

			// Check the goroutine state
			if cancelled() {
				return
			}

			// Calling the working task
			err := t()
			if err != nil {
				// Make error handling - send error in channel
				// (check if channel is valid)
				if !cancelled() {
					errChan <- err
				}
				return
			}
			if !cancelled() {
				doneChan <- "done"
			}
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
					return nil // No goroutines leaking with buffered channel
				}
			}
		case err := <-errChan:
			// Checking errors from channel
			if err != nil {
				log.Printf("Task return error: %s", err.Error())
				errCounter++
				if errCounter == M {
					// Send stop signal to all goroutines
					close(abortChan)
					return errors.New("error limit is elapsed")
				}
			}
		}
	}

}
