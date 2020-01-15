package tasks

import (
	"errors"
	"log"
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

	// Create working channels
	errChan := make(chan error, M+1) // error channel
	doneChan := make(chan string, N) // success channel

	// Create task channel
	taskChan := make(chan func() error, N)

	// Initialize abort channel
	abortChan = make(chan struct{})

	// Use wait group to synchronize writes to channels on function exit
	var wg sync.WaitGroup

	defer func() {
		close(taskChan)
		wg.Wait() // Wait until all tasks exit
		close(errChan)
		close(doneChan)
		log.Println("Run() function is exited")
	}()

	// Creating N reusable goroutines for task execution
	for i := 0; i < N; i++ {
		// Create goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Make infinite loop with channel blocking on reading
			for {
				// Wait for task and check channels for closed state
				task := <-taskChan
				if cancelled() || task == nil {
					log.Println("Goroutine exited on channel closing")
					return
				}

				// Execute task
				err := task()
				if err != nil {
					// Error handling
					if !cancelled() {
						errChan <- err
					}
				}
				// Successful execution
				if !cancelled() {
					doneChan <- "done"
				}
			}
		}()
	}

	// Send all tasks for execution in a separate goroutine
	// to prevent blocking on writing
	go func() {
		for _, task := range tasks {
			// Check the goroutine state
			if cancelled() {
				return
			}
			taskChan <- task
		}
	}()

	var (
		errCounter  int
		doneCounter int
	)

	for {
		// Reading messages from state channels
		select {
		case state := <-doneChan:
			// Checking success state from channel
			if state == "done" {
				doneCounter++
			}
		case err := <-errChan:
			// Checking errors from channel
			doneCounter++
			if err != nil {
				log.Printf("Task return error: %s", err.Error())
				errCounter++
				if errCounter == M {
					// Send stop signal to all goroutines
					close(abortChan)
					return errors.New("error limit is exceeded")
				}
			}
		}

		// Checking if all tasks completed
		if doneCounter == len(tasks) {
			log.Printf("%d tasks were executed successfully", doneCounter)
			// Send stop signal to all goroutines
			close(abortChan)
			return nil
		}
	}

}
