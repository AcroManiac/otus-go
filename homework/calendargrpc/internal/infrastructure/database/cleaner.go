package database

import (
	"fmt"
	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/pkg/errors"
)

type Cleaner struct {
	conn   *Connection
	policy interfaces.RetentionPolicy
}

func NewDatabaseCleaner(conn *Connection, policy interfaces.RetentionPolicy) interfaces.Cleaner {
	return &Cleaner{
		conn:   conn,
		policy: policy,
	}
}

func (c *Cleaner) Clean() error {

	// Clean events
	err := ExecQuery(
		c.conn,
		fmt.Sprintf(
			"delete from events where start_time < now() - interval '%f hours'",
			c.policy.GetDuration().Hours()))
	if err != nil {
		return errors.Wrap(err, "failed deleting events table")
	}

	// Clean notices
	err = ExecQuery(
		c.conn,
		fmt.Sprintf(
			"delete from notices where send_time < now() - interval '%f hours'",
			c.policy.GetDuration().Hours()))
	if err != nil {
		return errors.Wrap(err, "failed deleting events table")
	}

	return nil
}
