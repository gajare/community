package main

import (
	"log"
	"os"
	"procore-call-logs/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file only if running outside of Kubernetes (i.e., locally)
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found. Continuing with system environment variables.")
		}
	}

	// Read environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	frontendURL := os.Getenv("ALLOWED_ORIGINS")
	if frontendURL == "" {
		frontendURL = "http://localhost:3002"
	}

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		allowedOrigins := map[string]bool{
			"http://localhost:3002":       true,
			"http://localhost:3000":       true,
			"http://call-logs-frontend":   true,
		}
	
		origin := c.Request.Header.Get("Origin")
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin") // Prevent caching issues
		}
	
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
	router.GET("/api/call_logs", handlers.GetcallLogs)
	router.GET("/api/call_logs/filter", handlers.GetFilteredCallLogs)
	router.POST("/api/call_logs", handlers.CreateCallLog)
	router.PUT("/api/call_logs/:id", handlers.UpdateCallLog)
	router.DELETE("/api/call_logs/:id", handlers.DeleteCallLog)
	router.GET("/api/call_logs/:id", handlers.GetcallLogDetails)

	// Start server
	log.Printf("Starting server on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
