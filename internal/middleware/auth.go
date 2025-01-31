package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AgusMolinaCode/restApi-Go.git/internal/models"
	"github.com/AgusMolinaCode/restApi-Go.git/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func Signup(c *gin.Context) {
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

func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserByEmail(loginData.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := models.VerifyPassword(user.Password, loginData.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generar token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func GetUserByID(c *gin.Context) {
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

func GetAllUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func UpdateUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	if id != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to update this user"})
		return
	}

	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener el usuario actual para preservar los valores existentes
	existingUser, err := models.GetUserByID(id)
	if err != nil || existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Validar que whatsapp y email no estén vacíos
	if updatedUser.Whatsapp == "" {
		updatedUser.Whatsapp = existingUser.Whatsapp
	}
	if updatedUser.Email == "" {
		updatedUser.Email = existingUser.Email
	}

	err = models.UpdateUserByID(id, updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	if id != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to delete this user"})
		return
	}

	err := models.DeleteUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func ForgotPassword(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := models.SetResetToken(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token", "details": err.Error()})
		return
	}

	// Construir el enlace de restablecimiento de contraseña
	resetLink := fmt.Sprintf("https://restapi-go-production.up.railway.app/reset-password?token=%s", token)

	// Enviar el correo electrónico
	subject := "Password Reset Request"
	body := fmt.Sprintf("Click the link to reset your password: %s", resetLink)
	if err := utils.SendEmail(request.Email, subject, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset token sent to email"})
}

func ResetPassword(c *gin.Context) {
	var request struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := models.VerifyResetToken(request.Token)
	if err != nil || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	err = models.UpdatePassword(userID, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
