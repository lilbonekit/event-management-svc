package repo

import (
	"context"
	"database/sql"

	"github.com/lilbonekit/event-management-svc/internal/models"
)

type UserRepo struct{ db *sql.DB }

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, q DBTX, email, hashed string) (int64, error) {
	res, err := q.ExecContext(ctx, `INSERT INTO users (email, password) VALUES (?, ?)`, email, hashed)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx, `SELECT id, email, password FROM users WHERE email=?`, email).
		Scan(&u.ID, &u.Email, &u.Password)
	return u, err
}
