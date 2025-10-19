package service

import (
	"context"
	"errors"
	"strings"

	"github.com/lilbonekit/event-management-svc/internal/models"
	"github.com/lilbonekit/event-management-svc/internal/repo"
	"github.com/lilbonekit/event-management-svc/internal/utils"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyInUse  = errors.New("email already in use")
)

type UserService struct {
	repo *repo.UserRepo
}

func NewUserService(r *repo.UserRepo) *UserService { return &UserService{repo: r} }

func (s *UserService) Register(ctx context.Context, in models.NewUser) (int64, error) {
	hashed, err := utils.HashPassword(in.Password)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.Create(ctx, in.Email, hashed)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, ErrEmailAlreadyInUse
		}
		return 0, err
	}
	return id, nil
}

func (s *UserService) SignIn(ctx context.Context, email, plain string) (string, int64, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", 0, ErrUserNotFound
	}

	if !utils.CheckPasswordHash(plain, u.Password) {
		return "", 0, ErrInvalidCredentials
	}

	tok, err := utils.GenerateToken(u.Email, u.ID)
	if err != nil {
		return "", 0, err
	}

	return tok, u.ID, nil
}
