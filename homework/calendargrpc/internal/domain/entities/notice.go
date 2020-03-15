package entities

import "time"

// Time entity enqueued for sending
type Notice struct {
	Id    string    `json:"id"`
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
	Owner string    `json:"owner"`
}
