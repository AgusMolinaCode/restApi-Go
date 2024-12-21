package main

import (
	"github.com/AgusMolinaCode/restApi-Go.git/db"
	"github.com/AgusMolinaCode/restApi-Go.git/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
