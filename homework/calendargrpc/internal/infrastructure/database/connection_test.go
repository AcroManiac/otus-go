package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	dbUser     = "dbuser"
	dbPassword = "En9NR2b869"
	dbHost     = "127.0.0.1"
	dbName     = "calendar"
	dbPort     = 5432
)

func CreateConnection(t *testing.T) *Connection {
	conn := NewDatabaseConnection(dbUser, dbPassword, dbHost, dbName, dbPort)
	require.NotNil(t, conn, "Database connection object should not be nil")
	err := conn.Init(context.Background())
	require.Nil(t, err, "Should be no error while initializing database connection")
	return conn
}

func TestConnection_Init(t *testing.T) {
	_ = CreateConnection(t)
}

func TestConnection_Get(t *testing.T) {
	// Test connection first
	conn := CreateConnection(t)
	pooledConn, err := conn.Get(context.Background())
	require.NotNil(t, pooledConn, "Pooled connection should not be nil")
	require.Nil(t, err, "Should be no error while getting connection from pool")
}
