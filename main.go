package main

import (
	"log"
	"time"

	"abc.com/calc/db"
	"abc.com/calc/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	db.InitDB()

	// Create a Gin router
	server := gin.Default()

	// Add CORS middleware
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Update with your frontend's URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register routes
	routes.RegisterRoutes(server)

	// Log a message to confirm server startup
	log.Println("Server is starting on port 8080...")

	// Start the server
	err := server.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
