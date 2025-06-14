package main

import (
	"log"
	"os"

	"tele-temp-mail/internal/config"
	"tele-temp-mail/internal/handlers"
	"tele-temp-mail/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize services
	emailService := services.NewEmailService(cfg)
	smtpService := services.NewSMTPService(cfg, emailService)

	// Start SMTP server in background
	go func() {
		if err := smtpService.Start(); err != nil {
			log.Fatalf("Failed to start SMTP server: %v", err)
		}
	}()

	// Setup HTTP router
	router := gin.Default()
	
	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Initialize handlers
	emailHandler := handlers.NewEmailHandler(emailService)

	// API routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", handlers.HealthCheck)
		v1.POST("/email", emailHandler.CreateEmail)
		v1.GET("/email/:id", emailHandler.GetEmail)
		v1.GET("/email/:id/messages", emailHandler.GetMessages)
		v1.DELETE("/email/:id", emailHandler.DeleteEmail)
	}

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on 0.0.0.0:%s", port)
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
