package storage

import "time"

type Event struct {
	time time.Time
	desc string
}
