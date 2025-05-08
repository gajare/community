package main

import (
	"log"
	"os"

	"procore-accident-logs/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Gin router
	router := gin.Default()
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	// Configure CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", frontendURL)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Routes
	router.POST("/api/auth/token", handlers.GetAuthToken)
	router.GET("/api/accident-logs", handlers.GetAccidentLogs)
	router.GET("/api/accident-logs/filter", handlers.GetFilteredAccidentLogs)
	router.GET("/api/accident-type-logs/filter", handlers.GetAccidentTypeLogs)
	router.GET("/api/accident-logs/:id", handlers.GetAccidentLogDetails)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	router.Run(":" + port)
}
