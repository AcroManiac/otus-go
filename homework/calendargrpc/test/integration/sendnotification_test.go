package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AcroManiac/otus-go/homework/calendargrpc/internal/infrastructure/database"
	"github.com/pkg/errors"
)

type sendNotificationTest struct{}

func (t *sendNotificationTest) connectionToPostgreSQLServiceWithDSN(arg1 string) error {

	// Create cancel context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create and initialize connection pool to database
	db = database.NewDatabaseConnectionDSN(arg1)
	if err := db.Init(ctx); err != nil {
		return errors.Wrap(err, "unable to connect to database")
	}

	return nil
}

func (t *sendNotificationTest) iWaitForMinute(arg1 int) error {
	log.Printf("Waiting for %d minute...", arg1)
	time.Sleep(time.Duration(arg1) * time.Minute)
	return nil
}

func (t *sendNotificationTest) selectionFromTableShouldReturnNotice(
	arg1 string, arg2 int) error {

	// Create context for query execution
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get connection from pool
	conn, err := db.Get(ctx)
	if err != nil {
		return errors.Wrap(err, "couldn't get connection from pool")
	}
	defer conn.Release()

	// Select records from database
	var count int
	queryText := fmt.Sprintf("SELECT count(*) FROM %s WHERE id = '%s'", arg1, eventID)
	row := conn.QueryRow(ctx, queryText)
	if err = row.Scan(&count); err != nil {
		return errors.Wrap(err, "error selecting records from database")
	}
	if count != arg2 {
		return errors.New("wrong number of notices in database")
	}

	return nil
}
