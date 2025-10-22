package repo

import (
	"context"
	"database/sql"

	"github.com/lilbonekit/event-management-svc/internal/models"
)

type EventRepo struct{ db *sql.DB }

func NewEventRepo(db *sql.DB) *EventRepo { return &EventRepo{db: db} }

func (r *EventRepo) List(ctx context.Context) ([]models.Event, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, description, location, datetime, user_id FROM events`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserID); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *EventRepo) GetByID(ctx context.Context, id int64) (models.Event, error) {
	var e models.Event
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, description, location, datetime, user_id FROM events WHERE id=?`, id).
		Scan(&e.ID, &e.Name, &e.Description, &e.Location, &e.DateTime, &e.UserID)
	return e, err
}

func (r *EventRepo) Create(ctx context.Context, q DBTX, in models.NewEvent) (int64, error) {
	res, err := q.ExecContext(ctx,
		`INSERT INTO events (name, description, location, datetime, user_id) VALUES (?,?,?,?,?)`,
		in.Name, in.Description, in.Location, in.DateTime, in.UserID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *EventRepo) Update(ctx context.Context, q DBTX, e models.Event) error {
	_, err := q.ExecContext(ctx,
		`UPDATE events SET name=?, description=?, location=?, datetime=? WHERE id=? AND user_id=?`,
		e.Name, e.Description, e.Location, e.DateTime, e.ID, e.UserID,
	)
	return err
}

func (r *EventRepo) Delete(ctx context.Context, q DBTX, id, ownerID int64) (int64, error) {
	res, err := q.ExecContext(ctx, `DELETE FROM events WHERE id=? AND user_id=?`, id, ownerID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *EventRepo) Register(ctx context.Context, q DBTX, userID, eventID int64) error {
	_, err := q.ExecContext(ctx, `INSERT INTO registrations (user_id, event_id) VALUES (?, ?)`, userID, eventID)
	return err
}

func (r *EventRepo) CancelRegistration(ctx context.Context, q DBTX, userID, eventID int64) error {
	_, err := q.ExecContext(ctx, `DELETE FROM registrations WHERE user_id=? AND event_id=?`, userID, eventID)
	return err
}
