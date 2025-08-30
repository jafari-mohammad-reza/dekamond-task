package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Mobile    string    `json:"mobile"`
	CreatedAt time.Time `json:"created_at"`
}
