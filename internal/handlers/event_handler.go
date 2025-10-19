package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lilbonekit/event-management-svc/internal/models"
	"github.com/lilbonekit/event-management-svc/internal/service"
)

type EventHandler struct{ svc *service.EventService }

func NewEventHandler(s *service.EventService) *EventHandler { return &EventHandler{svc: s} }

// GET /events
func (h *EventHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	evs, err := h.svc.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve events"})
		return
	}
	c.JSON(http.StatusOK, evs)
}

// GET /events/:id
func (h *EventHandler) Get(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID"})
		return
	}
	e, err := h.svc.Get(ctx, id)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}
	c.JSON(http.StatusOK, e)
}

// POST /events  (JWT)
func (h *EventHandler) Create(c *gin.Context) {
	var in models.Event
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	id, err := h.svc.Create(ctx, models.NewEvent{
		Name: in.Name, Description: in.Description, Location: in.Location, DateTime: in.DateTime,
	}, c.GetInt64("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create event"})
		return
	}
	in.ID = id
	in.UserID = c.GetInt64("userID")
	c.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": in})
}

// PUT /events/:id  (JWT)
func (h *EventHandler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID"})
		return
	}
	var in models.Event
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}
	in.ID = id
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.Update(ctx, in, c.GetInt64("userID")); err != nil {
		switch err {
		case service.ErrNotOwner:
			c.JSON(http.StatusForbidden, gin.H{"message": "You do not have permission to update this event"})
		case service.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event updated", "event": in})
}

// DELETE /events/:id  (JWT)
func (h *EventHandler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID"})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.Delete(ctx, id, c.GetInt64("userID")); err != nil {
		switch err {
		case service.ErrNotOwner:
			c.JSON(http.StatusForbidden, gin.H{"message": "You do not have permission to delete this event"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete event"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event deleted"})
}

func parseID(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }

// POST /events/:id/registrations  (JWT)
func (h *EventHandler) Register(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID"})
		return
	}

	uid := c.GetInt64("userID")
	if uid <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.svc.Register(ctx, uid, id); err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		case errors.Is(err, service.ErrAlreadyJoined):
			c.JSON(http.StatusConflict, gin.H{"message": "Already registered"})
		case errors.Is(err, service.ErrEventPassed):
			c.JSON(http.StatusBadRequest, gin.H{"message": "Event already passed"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registered for event"})
}

// DELETE /events/:id/registrations  (JWT)
func (h *EventHandler) CancelRegistration(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID"})
		return
	}

	uid := c.GetInt64("userID")
	if uid <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.svc.CancelRegistration(ctx, uid, id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not cancel registration"})
		return
	}

	c.Status(http.StatusNoContent)
}
