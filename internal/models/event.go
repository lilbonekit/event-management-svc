package models

import "time"

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"date_time" binding:"required"` // RFC3339
	UserID      int64     `json:"user_id"`
}

type NewEvent struct {
	Name        string
	Description string
	Location    string
	DateTime    time.Time
	UserID      int64
}
