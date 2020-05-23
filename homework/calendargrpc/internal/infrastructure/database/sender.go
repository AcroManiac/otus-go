package database

import (
	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/entities"
	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/pkg/errors"
)

// DatabaseSender structure
type DatabaseSender struct {
	conn *Connection
}

// NewLogSender constructs LogSender structure
func NewDatabaseSender(conn *Connection) interfaces.Sender {
	return &DatabaseSender{conn: conn}
}

// Send function writes data from notice to database.
// Returns error if any
func (s *DatabaseSender) Send(notice entities.Notice) error {
	err := ExecQuery(s.conn,
		"INSERT INTO notices (id, send_time) VALUES ($1, now())",
		notice.Id)
	if err != nil {
		return errors.Wrap(err, "failed inserting notice to database")
	}
	return nil
}
