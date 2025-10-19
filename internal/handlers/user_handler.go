package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lilbonekit/event-management-svc/internal/models"
	"github.com/lilbonekit/event-management-svc/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler { return &UserHandler{svc: s} }

func (h *UserHandler) SignUp(c *gin.Context) {
	var in models.User
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := h.svc.Register(ctx, models.NewUser{Email: in.Email, Password: in.Password})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyInUse):
			c.JSON(http.StatusConflict, gin.H{"message": "Email already in use"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user_id": id})
}

func (h *UserHandler) SignIn(c *gin.Context) {
	var in models.User
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	tok, id, err := h.svc.SignIn(ctx, in.Email, in.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound),
			errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User signed in", "token": tok, "user_id": id})
}
