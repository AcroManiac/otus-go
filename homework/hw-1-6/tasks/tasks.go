package tasks

import (
	"errors"
	"log"
	"sync"
)

var taskChan chan func() error

func cancelled() bool {
	select {
	case <-taskChan:
		return true
	default:
		return false
	}
}

func sendTasks(tasks []func() error, N int, currentIndex *int) {
	if cancelled() {
		return
	}

	stopIndex := *currentIndex + N
	for ; *currentIndex < stopIndex && *currentIndex < len(tasks); *currentIndex++ {
		taskChan <- tasks[*currentIndex]
		log.Printf("Task %d was sent for execution", *currentIndex)
	}
}

func Run(tasks []func() error, N int, M int) error {

	// Create working channels
	errChan := make(chan error, M)   // error channel
	doneChan := make(chan string, N) // success channel

	// Create task channel
	taskChan = make(chan func() error, N)

	// Use wait group to synchronize writes to channels on function exit
	var wg sync.WaitGroup

	defer func() {
		wg.Wait() // Wait until all tasks exit
		close(errChan)
		close(doneChan)
		log.Println("Run() function is exited")
	}()

	var (
		errCounter  int
		doneCounter int
		taskIndex   int
	)

	// Starting task sending pipeline
	sendTasks(tasks, N, &taskIndex)

	// Reading messages from both channels
	for {
		select {
		case task := <-taskChan:
			// Read task from channel and run it in a separate goroutine
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
			log.Println("Goroutine with task created")
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
					close(taskChan)
					return errors.New("error limit is exceeded")
				}
			}
		}

		// Checking if N goroutines exited
		if doneCounter == N {
			if taskIndex == len(tasks) {
				log.Printf("%d tasks were executed successfully", taskIndex)
				// Send stop signal to all goroutines
				close(taskChan)
				return nil
			}
			// Send next task bundle for execution
			sendTasks(tasks, N, &taskIndex)
			doneCounter = 0
		}
	}

}
