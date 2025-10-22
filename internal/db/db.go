package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Path         string
	PingTimeout  time.Duration
	MaxOpenConns int
	MaxIdleConns int
}

func Open(cfg Config) (*sql.DB, error) {
	d, err := sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(cfg.MaxOpenConns)
	d.SetMaxIdleConns(cfg.MaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer cancel()
	if err := d.PingContext(ctx); err != nil {
		_ = d.Close()
		return nil, err
	}
	if err := createTables(d); err != nil {
		_ = d.Close()
		return nil, err
	}
	return d, nil
}

func createTables(d *sql.DB) error {
	if _, err := d.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`); err != nil {
		return err
	}

	if _, err := d.Exec(`
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime DATETIME NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`); err != nil {
		return err
	}

	if _, err := d.Exec(`
	CREATE TABLE IF NOT EXISTS registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		event_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (event_id) REFERENCES events(id),
		UNIQUE(user_id, event_id)
	);`); err != nil {
		return err
	}

	return nil
}
