package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lilbonekit/event-management-svc/models"
)

// getEvents godoc
// @Summary      List all events
// @Description  Returns a list of events
// @Tags         events
// @Produce      json
// @Success      200 {array}  models.Event
// @Failure      500 {object} map[string]string
// @Router       /events [get]
func getEvents(c *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not retrieve events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

// getEvent godoc
// @Summary      Get an event by ID
// @Tags         events
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  models.Event
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /events/{id} [get]
func getEvent(c *gin.Context) {
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

	c.JSON(http.StatusOK, event)
}

// createEvent godoc
// @Summary      Create a new event
// @Description  Creates a new event (requires JWT)
// @Tags         events
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        event body     models.Event true "Event to create"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /events [post]
func createEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	event.UserID = c.GetInt64("userID")

	if err := event.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": event})
}

// updateEvent godoc
// @Summary      Update an existing event
// @Description  Updates an event by ID (requires JWT)
// @Tags         events
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path int           true "Event ID"
// @Param        event body models.Event  true "Updated event"
// @Success      200   {object} map[string]interface{}
// @Failure      400   {object} map[string]string
// @Failure      401   {object} map[string]string
// @Failure      403   {object} map[string]string
// @Failure      404   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Router       /events/{id} [put]
func updateEvent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID"})
		return
	}

	existingEvent, err := models.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}

	if c.GetInt64("userID") != existingEvent.UserID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have permission to update this event"})
		return
	}

	var updatedEvent models.Event
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	updatedEvent.ID = existingEvent.ID

	if err := updatedEvent.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated", "event": updatedEvent})
}

// deleteEvent godoc
// @Summary      Delete an event
// @Description  Deletes an event by ID (requires JWT)
// @Tags         events
// @Security     BearerAuth
// @Param        id   path int true "Event ID"
// @Produce      json
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      401  {object} map[string]string
// @Failure      403  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /events/{id} [delete]
func deleteEvent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid event ID"})
		return
	}

	existingEvent, err := models.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}

	if c.GetInt64("userID") != existingEvent.UserID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have permission to delete this event"})
		return
	}

	if err := existingEvent.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted"})
}
