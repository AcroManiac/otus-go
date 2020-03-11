package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func GetEventsQueryContext(ctx context.Context, db *Connection, queryText string) ([]entities.Event, error) {
	var result []entities.Event

	// Get connection from pool
	conn, err := db.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	// Select records from database
	rows, err := conn.Query(ctx, queryText)
	if err != nil {
		return nil, errors.Wrap(err, "error selecting records from database")
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		var title string
		var description sql.NullString
		var owner string
		var startTime time.Time
		var duration time.Duration
		var notify sql.NullString
		err := rows.Scan(
			&id, &title, &description, &owner,
			&startTime, &duration, &notify)
		if err != nil {
			return nil, errors.Wrap(err, "error occurred while scanning record data")
		}

		// Add data to result slice
		result = append(result, entities.Event{
			Id:        entities.IdType(id),
			Title:     title,
			StartTime: startTime,
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

	return result, nil
}
