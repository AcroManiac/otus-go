package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/pkg/errors"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/google/uuid"
)

type Storage struct {
	ctx  context.Context
	conn *Connection
}

func NewDatabaseStorage(ctx context.Context, user, password, host, database string, port int) interfaces.Storage {
	conn := NewDatabaseConnection(user, password, host, database, port)
	if err := conn.Init(ctx); err != nil {
		logger.Fatal("unable to connect to database", "error", err)
	}
	return &Storage{ctx: ctx, conn: conn}
}

// isExistTime is checking record with specified time and owner
func (s Storage) isExistTime(time time.Time, owner string) (entities.IdType, bool, error) {
	id := entities.IdType(uuid.UUID{})
	// Get connection from pool
	conn, err := s.conn.Get(s.ctx)
	if err != nil {
		return id, false, err
	}
	defer conn.Release()

	// Find record in database
	rows, err := conn.Query(s.ctx, "select id from events where start_time=$1 and owner=$2", time, owner)
	if err != nil {
		return id, false, errors.Wrap(err, "error finding record in database")
	}
	defer rows.Close()

	if rows.Next() {
		// Found specified time
		var storedId uuid.UUID
		if err = rows.Scan(&storedId); err == nil {
			id = entities.IdType(storedId)
		}
		return id, true, nil
	}

	// Time does not exist
	return id, false, nil
}

func ToInterval(d *time.Duration) interface{} {
	if d == nil {
		return nil
	}
	return fmt.Sprintf("%d ms", d.Milliseconds())
}

func (s *Storage) Add(ev entities.Event) (entities.IdType, error) {
	// Check if event time is occupied already
	id, ok, err := s.isExistTime(ev.StartTime, ev.Owner)
	if err != nil {
		return id, err
	}
	if ok {
		return id, entities.ErrTimeBusy
	}

	// Get connection from pool
	conn, err := s.conn.Get(s.ctx)
	if err != nil {
		return id, err
	}
	defer conn.Release()

	// Insert new record to database
	id = entities.IdType(uuid.New())
	_, err = conn.Exec(
		s.ctx,
		"insert into events(id, title, description, owner, start_time, duration, notify) "+
			"values ($1, $2, $3, $4, $5, $6, $7);",
		uuid.UUID(id), ev.Title, ev.Description, ev.Owner, ev.StartTime,
		ToInterval(&ev.Duration), ToInterval(ev.Notify),
	)
	if err != nil {
		return id, errors.Wrap(err, "error inserting new record to database")
	}

	return id, nil
}

func (s *Storage) Remove(id entities.IdType) error {
	// Get connection from pool
	conn, err := s.conn.Get(s.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// Delete record from database
	_, err = conn.Exec(s.ctx, "delete from events where id=$1", uuid.UUID(id))
	if err != nil {
		return errors.Wrap(err, "error deleting record from database")
	}

	return nil
}

func (s *Storage) Edit(id entities.IdType, ev entities.Event) error {
	// Check if event time is occupied already by other event
	existId, ok, err := s.isExistTime(ev.StartTime, ev.Owner)
	if err != nil {
		return err
	}
	if ok && existId != id {
		return entities.ErrTimeBusy
	}

	// Get connection from pool
	conn, err := s.conn.Get(s.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// Update record in database
	_, err = conn.Exec(s.ctx,
		"update events set title=$1, description=$2, owner=$3, start_time=$4, duration=$5, notify=$6 "+
			"where id=$7",
		ev.Title, ev.Description, ev.Owner, ev.StartTime,
		ToInterval(&ev.Duration), ToInterval(ev.Notify), uuid.UUID(id),
	)
	if err != nil {
		return errors.Wrap(err, "error updating record in database")
	}

	return nil
}

func (s Storage) GetEventsByTimePeriod(period entities.TimePeriod, t time.Time) ([]entities.Event, error) {
	var selected []entities.Event
	var startTime, stopTime time.Time

	// Calculate start and stop times for searching interval
	switch period {
	case entities.Day:
		startTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		stopTime = startTime.Add(24 * time.Hour)
	case entities.Month:
		startTime = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		stopTime = startTime.AddDate(0, 1, 0)
	case entities.Week:
		startTime = t.AddDate(0, 0, -int(t.Weekday()))
		stopTime = startTime.AddDate(0, 0, 7)
	}

	// Get connection from pool
	conn, err := s.conn.Get(s.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	// Select records from database
	rows, err := conn.Query(s.ctx,
		"select id, title, description, owner, start_time, duration, notify from events where start_time>=$1 and start_time<$2",
		startTime, stopTime)
	if err != nil {
		return nil, errors.Wrap(err, "error selecting records from database")
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id          uuid.UUID
			title       string
			description sql.NullString
			owner       string
			start_time  time.Time
			duration    time.Duration
			notify      sql.NullString
		)
		err := rows.Scan(
			&id, &title, &description, &owner,
			&start_time, &duration, &notify)
		if err != nil {
			return nil, errors.Wrap(err, "error occurred while scanning record data")
		}

		// Add data to result slice
		selected = append(selected, entities.Event{
			Id:        entities.IdType(id),
			Title:     title,
			StartTime: start_time,
			Duration:  duration,
			Description: func(ns *sql.NullString) *string {
				if ns.Valid {
					return &ns.String
				}
				return nil
			}(&description),
			Owner: owner,
			Notify: func(ns *sql.NullString) *time.Duration {
				if ns.Valid {
					dur, err := time.ParseDuration(ns.String)
					if err != nil {
						return nil
					}
					return &dur
				}
				return nil
			}(&description),
		})
	}

	return selected, nil
}
