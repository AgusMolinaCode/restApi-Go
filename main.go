package main

import (
	"log"

	"github.com/AgusMolinaCode/restApi-Go.git/pkg/database"
	"github.com/AgusMolinaCode/restApi-Go.git/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	database.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
