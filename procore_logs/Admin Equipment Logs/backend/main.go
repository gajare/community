package main

import (
	"log"
	"os"
	"procore-call-logs/handlers"

	// "procore-equipment_logs/handlers"

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
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
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
	router.GET("/api/equipment_logs", handlers.GetcallLogs)
	router.GET("/api/equipment_logs/filter", handlers.GetFilteredCallLogs)
	router.POST("/api/equipment_logs", handlers.CreateCallLog)
	router.PUT("/api/equipment_logs/:id", handlers.UpdateCallLog)
	router.DELETE("/api/equipment_logs/:id", handlers.DeleteCallLog)
	router.GET("/api/equipment_logs/:id", handlers.GetcallLogDetails)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	router.Run(":" + port)
}
