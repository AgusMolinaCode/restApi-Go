package routes

import (
	"net/http"
	"time"

	"github.com/AgusMolinaCode/restApi-Go.git/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getEvents(c *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func getEventByID(c *gin.Context) {
	id := c.Param("id")
	event, err := models.GetEventByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event", "details": err.Error()})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	c.JSON(http.StatusOK, event)
}

func createEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el user_id del token JWT
	userID, _ := c.Get("userID")

	event.ID = uuid.New().String()
	event.UserID = userID.(string) // Asignar el user_id del token al evento
	event.CreatedAt = time.Now().Format(time.RFC3339)
	event.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := event.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully"})
}

func updateEventByID(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var updatedEvent models.Event
	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := models.GetEventByID(id)
	if err != nil || event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this event"})
		return
	}

	updatedEvent.UpdatedAt = time.Now().Format(time.RFC3339)

	err = models.UpdateEventByID(id, updatedEvent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

func deleteEventByID(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	event, err := models.GetEventByID(id)
	if err != nil || event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this event"})
		return
	}

	err = models.DeleteEventByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func registerForEvent(c *gin.Context) {
	// Implementa la lógica para registrar un usuario en un evento
	c.JSON(http.StatusOK, gin.H{"message": "User registered for event"})
}

func cancelRegistration(c *gin.Context) {
	// Implementa la lógica para cancelar el registro de un usuario en un evento
	c.JSON(http.StatusOK, gin.H{"message": "Registration cancelled"})
}
