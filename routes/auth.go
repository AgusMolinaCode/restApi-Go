package routes

import (
	"net/http"

	"github.com/AgusMolinaCode/restApi-Go.git/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = uuid.New().String()

	if err := user.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := models.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user", "details": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func getAllUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
