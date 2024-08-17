package main

import (
	"github.com/gin-gonic/gin"
	"mes-gin/config"
	"mes-gin/routes"
)

// @title Gin MES
// @version 0.0.1
// @description This is a sample server for a mes application.
// @termsOfService http://swagger.io/terms/

// @contact.name Joseph Wang
// @contact.url https://www.hfwong.com
// @contact.email 13510870118@gmail.com

// @host localhost:8080
// @BasePath /

func main() {
	// Load the application configuration
	err := config.Load()
	if err != nil {
		panic("Failed to load configuration")
	}

	// Initialize the database connection
	config.InitDB()

	// Initialize the redis connection
	config.InitRedis()

	// Create a new Gin router
	router := gin.Default()

	// Set up the application routes
	routes.SetupRoutes(router)

	// Start the server
	router.Run(config.AppConfig.ServerPort)
}
