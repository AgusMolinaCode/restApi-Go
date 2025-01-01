package routes

import (
	"github.com/AgusMolinaCode/restApi-Go.git/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/events", getEvents)
	router.GET("/events/:id", getEventByID)
	router.GET("/events/:id/registrations", getRegistrationsByEventID)
	router.GET("/tags", getAllTags)

	protected := router.Group("/", middleware.AuthMiddleware())
	{
		protected.POST("/events", createEvent)
		protected.PUT("/events/:id", updateEventByID)
		protected.DELETE("/events/:id", deleteEventByID)
		protected.POST("/events/:id/register", registerForEvent)
		protected.DELETE("/events/:id/register", cancelRegistration)
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
