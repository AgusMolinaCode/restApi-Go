package routes

import (
	"github.com/AgusMolinaCode/restApi-Go.git/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/events", getEvents)
	router.GET("/events/:id", getEventByID)
	router.GET("/tags", getAllTags)
	router.GET("/events/by-tags", getEventsByTags)
	router.GET("/events/by-category", getEventsByCategory)
	router.GET("/events/by-date", getEventsByDate)
	router.GET("/events/categories", getAllCategories)
	router.GET("/events/by-name", getEventsByName)
	router.GET("/events/summaries", getEventSummaries)

	protected := router.Group("/", middleware.AuthMiddleware())
	{
		protected.POST("/events", createEvent)
		protected.PUT("/events/:id", updateEventByID)
		protected.DELETE("/events/:id", deleteEventByID)
		protected.POST("/events/:id/register", registerForEvent)
		protected.DELETE("/events/:id/register", cancelRegistration)
		protected.GET("/events/:id/registration", getRegistrationByEvent)
		protected.PUT("/users/:id", middleware.UpdateUserByID)
		protected.DELETE("/users/:id", middleware.DeleteUserByID)
	}

	router.POST("/signup", middleware.Signup)
	router.POST("/login", middleware.Login)
	router.POST("/forgot-password", middleware.ForgotPassword)
	router.POST("/reset-password", middleware.ResetPassword)
	router.GET("/users/:id", middleware.GetUserByID)
	router.GET("/users", middleware.GetAllUsers)
}
