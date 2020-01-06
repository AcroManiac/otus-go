package tasks

import (
	"errors"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	// Make randomization
	rand.Seed(time.Now().UnixNano())
}

func emulateActivity(taskNumber, delaySec int, error bool) {
	time.Sleep(time.Duration(delaySec) * time.Second)
	if !error {
		log.Printf("Task %d fired in %d seconds", taskNumber, delaySec)
	}
}

var tasksNoErrors = []func() error{
	func() error {
		emulateActivity(1, rand.Intn(5), false)
		return nil
	},
	func() error {
		emulateActivity(2, rand.Intn(5), false)
		return nil
	},
	func() error {
		emulateActivity(3, rand.Intn(5), false)
		return nil
	},
	func() error {
		emulateActivity(4, rand.Intn(5), false)
		return nil
	},
}

var tasksWithErrors = []func() error{
	func() error {
		emulateActivity(1, 1, true)
		return errors.New("error in task 1")
	},
	func() error {
		emulateActivity(2, 1, true)
		return errors.New("error in task 2")
	},
	func() error {
		emulateActivity(3, 1, true)
		return errors.New("error in task 3")
	},
	func() error {
		emulateActivity(4, 2, false)
		return nil
	},
	func() error {
		emulateActivity(5, 3, false)
		return nil
	},
}

const (
	concurrentTaskNumber = 4
	errorsLimit          = 2
)

func TestRun(t *testing.T) {
	log.Println("************ Test 1 ************")
	err := Run(tasksNoErrors, concurrentTaskNumber, errorsLimit)
	assert.Nil(t, err, "There should be no errors in this test")

	log.Println()
	log.Println("************ Test 2 ************")
	err = Run(tasksWithErrors, concurrentTaskNumber, errorsLimit)
	assert.NotNil(t, err, "There should be error in this test")
}
