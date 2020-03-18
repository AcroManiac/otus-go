package database

import (
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
)

type EventsCollector struct {
	conn *Connection
}

func NewDatabaseEventsCollector(conn *Connection) interfaces.EventsCollector {
	return &EventsCollector{conn: conn}
}

func (ec EventsCollector) GetEvents() ([]entities.Event, error) {
	// Select records from database
	events, err := GetEventsQueryContext(ec.conn,
		"SELECT id, title, description, owner, start_time, duration, notify "+
			"FROM events "+
			"WHERE start_time - notify <= now() "+
			"AND id NOT IN (SELECT id FROM notices)")
	if err != nil {
		return nil, err
	}
	return events, nil
}
