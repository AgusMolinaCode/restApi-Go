package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/events", getEvents)
	router.GET("/events/:id", getEventByID)
	router.GET("/events/:id/registrations", getRegistrationsByEventID)

	protected := router.Group("/", AuthMiddleware())
	{
		protected.POST("/events", createEvent)
		protected.PUT("/events/:id", updateEventByID)
		protected.DELETE("/events/:id", deleteEventByID)
		protected.POST("/events/:id/register", registerForEvent)
		protected.DELETE("/events/:id/register", cancelRegistration)
		protected.PUT("/users/:id", updateUserByID)
		protected.DELETE("/users/:id", deleteUserByID)
	}

	router.POST("/signup", signup)
	router.POST("/login", login)
	router.POST("/forgot-password", forgotPassword)
	router.POST("/reset-password", resetPassword)
	router.GET("/users/:id", getUserByID)
	router.GET("/users", getAllUsers)
}