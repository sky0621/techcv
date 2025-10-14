// Package mysql provides MySQL infrastructure components for the manager service.
package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config holds the connection parameters required to connect to MySQL.
type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Params   string
}

// DSN builds a MySQL DSN suitable for github.com/go-sql-driver/mysql.
func (c Config) DSN() string {
	params := c.Params
	if params == "" {
		params = "parseTime=true&loc=UTC&charset=utf8mb4"
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		params,
	)
}

// NewConnection opens a MySQL connection using the provided configuration and verifies connectivity.
func NewConnection(ctx context.Context, cfg Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, err
	}

	const (
		connMaxLifetime = 5 * time.Minute
		maxOpenConns    = 20
		maxIdleConns    = 10
		pingTimeout     = 5 * time.Second
	)

	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
