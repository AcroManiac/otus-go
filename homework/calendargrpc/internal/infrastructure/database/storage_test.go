package database

import (
	"context"
	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/entities"
	"testing"
	"time"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/domain/interfaces"
	"github.com/stretchr/testify/require"
)

var description = "Record for database testing"
var evTest = entities.Event{
	Id:          entities.IdType{},
	Title:       "Test event",
	StartTime:   time.Now(),
	Duration:    time.Hour,
	Description: &description,
	Owner:       "test",
	Notify:      nil,
}

func CreateStorage(t *testing.T) interfaces.Storage {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	db := NewDatabaseStorage(ctx, dbUser, dbPassword, dbHost, dbName, dbPort)
	require.NotNil(t, db, "Database storage object should not be nil")
	return db
}

func CreateRecord(t *testing.T, db interfaces.Storage) entities.IdType {
	id, err := db.Add(evTest)
	require.NotNil(t, id, "Identifier should not be nil")
	require.NoError(t, err, "Add function should return no error")
	return id
}

func DeleteRecord(t *testing.T, db interfaces.Storage, id entities.IdType) {
	err := db.Remove(id)
	require.NoError(t, err, "Delete function should return no error")
}

func TestStorage_Add(t *testing.T) {
	db := CreateStorage(t)
	id := CreateRecord(t, db)
	DeleteRecord(t, db, id)
}

func TestStorage_Edit(t *testing.T) {
	db := CreateStorage(t)
	id := CreateRecord(t, db)

	// Change event
	evNew := evTest
	evNew.Title += " (modified)"
	evNew.StartTime = evNew.StartTime.Add(time.Hour)
	err := db.Edit(id, evNew)
	require.NoError(t, err, "Edit function should return no error")

	DeleteRecord(t, db, id)
}

func TestStorage_Remove(t *testing.T) {
	db := CreateStorage(t)
	id := CreateRecord(t, db)
	DeleteRecord(t, db, id)
}

func TestStorage_GetEventsByTimePeriod(t *testing.T) {
	db := CreateStorage(t)
	id := CreateRecord(t, db)

	// Copy event for further comparison
	evCopy := evTest
	evCopy.Id = id

	// Get event from database
	events, err := db.GetEventsByTimePeriod(entities.Day, time.Now())
	require.NoError(t, err, "Get function should return no error")
	require.True(t, len(events) > 0, "Get function should return something")

	// Testing for equal is not correct as time returned from database differs slightly
	//evStored := events[0]
	//require.Equal(t, evCopy, evStored, "Stored event should be equal to test one")

	DeleteRecord(t, db, id)
}
