package database

import (
	"context"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/logger"
)

type EventsCollector struct {
	ctx  context.Context
	conn *Connection
}

func NewDatabaseEventsCollector(ctx context.Context, user, password, host, database string, port int) interfaces.EventsCollector {
	conn := NewDatabaseConnection(user, password, host, database, port)
	if err := conn.Init(ctx); err != nil {
		logger.Fatal("unable to connect to database", "error", err)
	}
	return &EventsCollector{ctx: ctx, conn: conn}
}

func (ec EventsCollector) GetEvents() ([]entities.Event, error) {
	// Select records from database
	events, err := GetEventsQueryContext(ec.ctx, ec.conn,
		`select id, title, description, owner, start_time, duration, notify 
			from events where start_time - notify <= now()`)
	if err != nil {
		return nil, err
	}
	return events, nil
}
