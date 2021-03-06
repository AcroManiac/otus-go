package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"net"
	"time"

	"github.com/ahamtat/otus-go/homework/calendargrpc/internal/infrastructure/logger"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Connection struct {
	connUri string
	pool    *pgxpool.Pool
}

func NewDatabaseConnection(user, password, host, database string, port int) *Connection {
	return NewDatabaseConnectionDSN(
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			user, password, host, port, database))
}

func NewDatabaseConnectionDSN(dsn string) *Connection {
	c := &Connection{connUri: dsn}
	return c
}

// Create and initialize connection pool to database
func (c *Connection) Init(ctx context.Context) error {
	cfg, err := pgxpool.ParseConfig(c.connUri)
	if err != nil {
		return errors.Wrap(err, "failed to parse postgres config")
	}

	cfg.MaxConns = 8
	cfg.ConnConfig.TLSConfig = nil
	cfg.ConnConfig.PreferSimpleProtocol = true
	cfg.ConnConfig.RuntimeParams["standard_conforming_strings"] = "on"
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		Timeout:   1 * time.Second,
		KeepAlive: 5 * time.Minute,
	}).DialContext

	if logger.GetLogger() != nil {
		cfg.ConnConfig.Logger = zapadapter.NewLogger(logger.GetLogger())
	}

	pool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to connect to postgres")
	}
	c.pool = pool

	return nil
}

// Get connection from connection pool
func (c *Connection) Get(ctx context.Context) (*pgxpool.Conn, error) {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed getting connection from pool")
	}
	return conn, nil
}
