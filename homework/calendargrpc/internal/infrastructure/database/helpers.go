package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// eventsScanHelper keeps data for database scanning
type eventsScanHelper struct {
	Id          uuid.UUID      `db:"id"`
	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`
	Owner       string         `db:"owner"`
	StartTime   time.Time      `db:"start_time"`
	Duration    time.Duration  `db:"duration"`
	Notify      sql.NullString `db:"notify"`
}

// GetEventsQueryContext queries database db with query text queryText and timed context ctx
// Function returns slice of events or error if any
func GetEventsQueryContext(db *Connection, queryText string) ([]entities.Event, error) {
	var result []entities.Event

	// Create context for query execution
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
		var esh eventsScanHelper
		err := rows.Scan(
			&esh.Id, &esh.Title, &esh.Description, &esh.Owner,
			&esh.StartTime, &esh.Duration, &esh.Notify)
		if err != nil {
			return nil, errors.Wrap(err, "error occurred while scanning record data")
		}

		// Add data to result slice
		result = append(result, entities.Event{
			Id:        entities.IdType(esh.Id),
			Title:     esh.Title,
			StartTime: esh.StartTime,
			Duration:  esh.Duration,
			Description: func(ns *sql.NullString) *string {
				if ns.Valid {
					return &ns.String
				}
				return nil
			}(&esh.Description),
			Owner: esh.Owner,
			Notify: func(ns *sql.NullString) *time.Duration {
				if ns.Valid {
					dur, err := time.ParseDuration(ns.String)
					if err != nil {
						return nil
					}
					return &dur
				}
				return nil
			}(&esh.Notify),
		})
	}

	return result, nil
}

// ExecQuery gets connection from connection pool, creates timed context and
// executes database query with arguments
func ExecQuery(db *Connection, queryText string, arguments ...interface{}) error {

	// Create context for query execution
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get connection from pool
	conn, err := db.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "error getting database connection from pool")
	}
	defer conn.Release()

	// Execute database query with params
	_, err = conn.Exec(ctx, queryText, arguments...)
	if err != nil {
		return errors.Wrap(err, "failed executing database query")
	}

	return nil
}
