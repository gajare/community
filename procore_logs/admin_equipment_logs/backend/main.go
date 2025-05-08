package main

import (
	"equipment_logs/handlers"
	"log"
	"os"

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
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3001"
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
	router.GET("/api/equipment_logs", handlers.GetEquipmentLogs)
	router.GET("/api/equipment_logs/filter", handlers.GetFilteredEquipmentLogs)
	router.POST("/api/equipment_logs", handlers.CreateEquipmentLogs)
	router.PUT("/api/equipment_logs/:id", handlers.UpdateEquipmentLogs)
	router.DELETE("/api/equipment_logs/:id", handlers.DeleteEquipmentLogs)
	router.GET("/api/equipment_logs/:id", handlers.GetEquipmentLogsDetails)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	router.Run(":" + port)
}
