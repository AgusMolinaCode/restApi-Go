package routes

import (
	"net/http"
	"time"
	"github.com/AgusMolinaCode/restApi-Go.git/internal/models"
	"github.com/AgusMolinaCode/restApi-Go.git/internal/services"
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

	// Obtener la primera fecha disponible
	var firstAvailableDate string
	for date, dateTime := range event.DateTimes {
		if dateTime.Status == "disponibles" || dateTime.Status == "pocas unidades" {
			firstAvailableDate = date
			break
		}
	}

	if firstAvailableDate == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "No available dates"})
		return
	}

	// Calcular los días restantes para el primer evento disponible
	eventTime, err := time.Parse("02/01/2006", firstAvailableDate) // Asegúrate de que el formato sea correcto
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid event date", "details": err.Error()})
		return
	}
	daysUntilEvent := int(time.Until(eventTime).Hours() / 24)

	// Obtener el clima solo si faltan 7 días o menos
	if daysUntilEvent <= 7 {
		weather, err := services.GetWeather(event.Location.Lat, event.Location.Lng, firstAvailableDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve weather", "details": err.Error()})
			return
		}

		// Incluir la información del clima en la respuesta
		c.JSON(http.StatusOK, gin.H{
			"event":   event,
			"weather": weather,
		})
	} else {
		// Solo devolver los datos del evento
		c.JSON(http.StatusOK, gin.H{
			"event": event,
		})
	}
}

func createEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generar un ID dinámico para el evento
	event.ID = uuid.New().String()

	// Obtener el user_id del usuario autenticado
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	event.UserID = userID.(string)

	event.CreatedAt = time.Now().Format(time.RFC3339)
	event.UpdatedAt = event.CreatedAt

	if err := event.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully", "event": event})
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

	// Verificar que al menos una fecha y hora de inicio esté presente
	if len(updatedEvent.DateTimes) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one start date and time is required"})
		return
	}

	// Validar que si se proporciona un título de pago, también se proporcione un enlace, y viceversa
	for title, payment := range updatedEvent.PaymentLink {
		if title == "" || payment.Link == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Both payment title and link must be provided"})
			return
		}
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

	var registrationData struct {
		EventDate   string `json:"event_date" binding:"required"`
		PaymentLink string `json:"payment_link" binding:"required"`
	}

	if err := c.ShouldBindJSON(&registrationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
		ID:          uuid.New().String(),
		EventID:     eventID,
		UserID:      userID.(string),
		Whatsapp:    "+1234567890", // Obtener de la base de datos o del contexto
		CreatedAt:   time.Now().Format(time.RFC3339),
		EventDate:   registrationData.EventDate,
		PaymentLink: registrationData.PaymentLink,
	}

	if err := registration.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register for event", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered for event", "registration": registration})
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

func getAllTags(c *gin.Context) {
	tags, err := models.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tags", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tags)
}

func getRegistrationByEvent(c *gin.Context) {
	eventID := c.Param("id")
	userID, _ := c.Get("userID")

	// Verificar si el usuario es el creador del evento
	event, err := models.GetEventByID(eventID)
	if err != nil || event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.UserID == userID {
		// Si es el creador, devolver todos los registros
		registrations, err := models.GetRegistrationsByEventID(eventID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve registrations", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, registrations)
		return
	}

	// Si no es el creador, verificar si el usuario está registrado
	isRegistered, err := models.IsUserRegisteredForEvent(eventID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check registration", "details": err.Error()})
		return
	}
	if !isRegistered {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to view this registration"})
		return
	}

	// Obtener el registro del usuario
	registration, err := models.GetRegistrationByUserID(eventID, userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve registration", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, registration)
}

func getEventsByTags(c *gin.Context) {
	tags := c.QueryArray("tags")

	events, err := models.GetEventsByTags(tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events", "details": err.Error()})
		return
	}

	if len(events) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found for the specified tags"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func getEventsByCategory(c *gin.Context) {
	category := c.Query("category")

	events, err := models.GetEventsByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events", "details": err.Error()})
		return
	}

	if len(events) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found for the specified category"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func getEventsByDate(c *gin.Context) {
	date := c.Query("date")

	events, err := models.GetEventsByDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events", "details": err.Error()})
		return
	}

	if len(events) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found for the specified date"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func getAllCategories(c *gin.Context) {
	categories, err := models.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func getEventsByName(c *gin.Context) {
	name := c.Query("name")

	events, err := models.GetEventsByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events", "details": err.Error()})
		return
	}

	if len(events) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No events found for the specified name"})
		return
	}

	c.JSON(http.StatusOK, events)
}
