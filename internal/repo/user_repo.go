package repo

import (
	"context"
	"database/sql"
)

type UserRepo struct{ db *sql.DB }

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, email, hashed string) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, password) VALUES (?, ?)`, email, hashed)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

type UserRow struct {
	ID       int64
	Email    string
	Password string
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (UserRow, error) {
	var u UserRow
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, password FROM users WHERE email = ?`, email).
		Scan(&u.ID, &u.Email, &u.Password)
	return u, err
}
