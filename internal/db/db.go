package db

import (
	"context"
	"database/sql"
	"dekamond-task/internal/config"
	"dekamond-task/internal/models"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
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
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		mobile VARCHAR(11) NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := d.conn.Exec(query)
	return err
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (d *DB) CreateUser(ctx context.Context, mobile string) error {
	res, err := d.conn.ExecContext(ctx, `INSERT INTO users (mobile , created_at) VALUES ($1 , $2)`, mobile, time.Now())
	if err != nil {
		return err
	}
	inserted, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if inserted <= 0 {
		return errors.New("nothing affected")
	}
	return nil
}

func (d *DB) GetUser(ctx context.Context, mobile string) (*models.User, error) {
	rows, err := d.conn.QueryContext(ctx, `SELECT (id , mobile , created_at) FROM users WHERE mobile=$1`, mobile)
	if err != nil {
		return nil, err
	}
	var user models.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Mobile, &user.CreatedAt); err != nil {
			return nil, err
		}
	}
	return &user, nil
}
func (d *DB) GetUsers(ctx context.Context, page, limit int) ([]*models.User, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	rows, err := d.conn.QueryContext(ctx, `
		SELECT id, mobile, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Mobile, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
