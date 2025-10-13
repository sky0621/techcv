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

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
