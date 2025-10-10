package models

import (
	"time"

	"github.com/lilbonekit/event-management-svc/db"
)

type Event struct {
	ID          int64     `json:"id" example:"1"`
	Name        string    `json:"name" binding:"required" example:"Go Meetup"`
	Description string    `json:"description" binding:"required" example:"Monthly developer meetup"`
	Location    string    `json:"location" binding:"required" example:"Kyiv, UNIT.City"`
	DateTime    time.Time `json:"date_time" binding:"required" example:"2025-10-08T18:30:00Z"`
	UserID      int64     `json:"user_id" example:"42"`
}

func (e *Event) Save() error {
	query := `INSERT INTO events (name, description, location, datetime, user_id) 
	VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		panic("Could not prepare statement.")
	}
	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	e.ID = id
	return nil
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var event Event

	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (e *Event) Update() error {
	query := `
	UPDATE events 
	SET name=?, description=?, location=?, datetime=? WHERE id=?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	return err
}

func (e *Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID)
	return err
}

func (e *Event) Register(userID int64) error {
	query := `INSERT INTO registrations (user_id, event_id) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(userID, e.ID)
	return err
}

func (e *Event) CancelRegistration(userID int64) error {
	query := `DELETE FROM registrations WHERE user_id = ? AND event_id = ?`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(userID, e.ID)
	return err
}
