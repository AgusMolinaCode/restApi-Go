package main

import (
	"log"

	"github.com/AgusMolinaCode/restApi-Go.git/internal/routes"
	"github.com/AgusMolinaCode/restApi-Go.git/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Intentar cargar .env, pero no fallar si no existe
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	database.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
