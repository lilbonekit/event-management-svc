package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lilbonekit/event-management-svc/internal/models"
	"github.com/lilbonekit/event-management-svc/internal/repo"
)

var (
	ErrNotFound      = errors.New("event not found")
	ErrNotOwner      = errors.New("forbidden: not owner")
	ErrAlreadyJoined = errors.New("already registered")
	ErrEventPassed   = errors.New("event already passed")
)

type EventService struct {
	db     *sql.DB
	events *repo.EventRepo
}

func NewEventService(db *sql.DB, ev *repo.EventRepo) *EventService {
	return &EventService{db: db, events: ev}
}

func (s *EventService) List(ctx context.Context) ([]models.Event, error) {
	return s.events.List(ctx)
}

func (s *EventService) Get(ctx context.Context, id int64) (models.Event, error) {
	e, err := s.events.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Event{}, ErrNotFound
		}
		return models.Event{}, err
	}
	return e, nil
}

func (s *EventService) Create(ctx context.Context, in models.NewEvent, callerID int64) (int64, error) {
	in.UserID = callerID
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	id, err := s.events.Create(ctx, tx, in)
	if err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (s *EventService) Update(ctx context.Context, e models.Event, callerID int64) error {
	cur, err := s.Get(ctx, e.ID)
	if err != nil {
		return err
	}
	if cur.UserID != callerID {
		return ErrNotOwner
	}

	e.UserID = callerID
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.events.Update(ctx, tx, e); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *EventService) Delete(ctx context.Context, id, callerID int64) error {
	cur, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if cur.UserID != callerID {
		return ErrNotOwner
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	aff, err := s.events.Delete(ctx, tx, id, callerID)
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrNotFound
	}
	return tx.Commit()
}

func (s *EventService) Register(ctx context.Context, userID, eventID int64) error {
	ev, err := s.Get(ctx, eventID)
	if err != nil {
		return err
	}

	if time.Now().UTC().After(ev.DateTime.UTC()) {
		return ErrEventPassed
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.events.Register(ctx, tx, userID, eventID); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrAlreadyJoined
		}
		return err
	}
	return tx.Commit()
}

func (s *EventService) CancelRegistration(ctx context.Context, userID, eventID int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.events.CancelRegistration(ctx, tx, userID, eventID); err != nil {
		return err
	}
	return tx.Commit()
}
