package db

import (
	"database/sql"
	"dekamond-task/internal/config"
	"fmt"
)

type DB struct {
	cfg  *config.Config
	conn *sql.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	conn, err := sql.Open("sqlite3", cfg.Database.Url)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{conn: conn}, nil
}
func (d *DB) InitTables() error {
	query := `
		CREATE DATABASE IF NOT EXISTS dekamond_task;
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			mobile VARCHAR(11) NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	if _, err := d.conn.Exec(query); err != nil {
		return err
	}
	return nil
}
func (db *DB) Close() error {
	return db.conn.Close()
}
