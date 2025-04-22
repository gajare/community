package main

import (
	"log"
	"os"
	"procore-call-logs/handlers"

	// "procore-call_logs/handlers"

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

	// Configure CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
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
	router.GET("/api/call_logs", handlers.GetAccidentLogs)
	router.GET("/api/call_logs/filter", handlers.GetFilteredAccidentLogs)
	router.POST("/api/call_logs", handlers.CreateAccidentLog)
	router.PUT("/api/call_logs/:id", handlers.UpdateAccidentLog)
	router.DELETE("/api/call_logs/:id", handlers.DeleteAccidentLog)
	router.GET("/api/call_logs/:id", handlers.GetAccidentLogDetails)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	router.Run(":" + port)
}
