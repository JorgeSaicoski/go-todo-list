package main

import (
	"os"

	"github.com/JorgeSaicoski/go-todo-list/internal/api"
	"github.com/JorgeSaicoski/go-todo-list/internal/db"
)

func main() {
	// Connect to the database
	db.ConnectDatabase()

	// Get router config, possibly from environment variables
	config := api.DefaultRouterConfig()

	// Override with environment variables if needed
	if origins := getEnv("ALLOWED_ORIGINS", ""); origins != "" {
		config.AllowedOrigins = origins
	}

	// Create router with full configuration
	taskRouter := api.NewTaskRouter(db.DB, config)

	// Register all routes
	taskRouter.RegisterRoutes()

	// Start the server
	port := getEnv("PORT", "8000")
	taskRouter.Run(":" + port)
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
