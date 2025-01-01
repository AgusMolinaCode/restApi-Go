package routes

import (
	"net/http"
	"time"

	"github.com/AgusMolinaCode/restApi-Go.git/internal/models"
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

	// Verificar que el número de tags no exceda 3
	if len(event.Tags) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A maximum of 3 tags are allowed per event"})
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

	// Verificar que el número de tags no exceda 3
	if len(updatedEvent.Tags) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A maximum of 3 tags are allowed per event"})
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
	eventID := c.Param("id")
	userID, _ := c.Get("userID")

	// Verificar si el usuario ya está registrado
	exists, err := models.IsUserRegisteredForEvent(eventID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check registration", "details": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already registered for this event"})
		return
	}

	// Registrar al usuario en el evento
	registration := models.Registration{
		ID:        uuid.New().String(),
		EventID:   eventID,
		UserID:    userID.(string),
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := registration.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register for event", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered for event"})
}

func cancelRegistration(c *gin.Context) {
	eventID := c.Param("id")
	userID, _ := c.Get("userID")

	// Cancelar el registro del usuario en el evento
	err := models.DeleteRegistration(eventID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel registration", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration cancelled"})
}

func getRegistrationsByEventID(c *gin.Context) {
	eventID := c.Param("id")

	registrations, err := models.GetRegistrationsByEventID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve registrations", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, registrations)
}

func getAllTags(c *gin.Context) {
	tags, err := models.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tags", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tags)
}
