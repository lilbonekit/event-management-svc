package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lilbonekit/event-management-svc/models"
)

// registerForEvent godoc
// @Summary      Register current user for an event
// @Description  Creates a registration record (requires JWT)
// @Tags         registrations
// @Security     BearerAuth
// @Param        id   path int true "Event ID"
// @Produce      json
// @Success      201  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /events/{id}/registrations [post]
func registerForEvent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID"})
		return
	}

	event, err := models.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}

	userID := c.GetInt64("userID")
	if err := event.Register(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not register for event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registered for event"})
}

// cancelRegistration godoc
// @Summary      Cancel current user's registration
// @Description  Deletes a registration record (requires JWT)
// @Tags         registrations
// @Security     BearerAuth
// @Param        id   path int true "Event ID"
// @Produce      json
// @Success      204  {string} string "No Content"
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /events/{id}/registrations [delete]
func cancelRegistration(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID"})
		return
	}

	event, err := models.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}

	userID := c.GetInt64("userID")
	if err := event.CancelRegistration(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not cancel registration"})
		return
	}

	c.Status(http.StatusNoContent)
}
