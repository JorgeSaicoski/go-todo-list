package db

import (
	"fmt"
	"os"
	"time"

	"github.com/JorgeSaicoski/pgconnect"
	"gorm.io/gorm"
)

var DB *pgconnect.DB

type Task struct {
	gorm.Model
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"` // expecting "pending", "in-progress", or "completed"
	DueDate     *time.Time `json:"dueDate,omitempty"`
}

func ConnectDatabase() {
	// Create config using environment variables
	config := pgconnect.DefaultConfig()

	// Override with environment variables if they exist
	config.Host = getEnv("POSTGRES_HOST", config.Host)
	config.Port = getEnv("POSTGRES_PORT", config.Port)
	config.User = getEnv("POSTGRES_USER", config.User)
	config.Password = getEnv("POSTGRES_PASSWORD", config.Password)
	config.DatabaseName = getEnv("POSTGRES_DB", "taskdb")
	config.SSLMode = getEnv("POSTGRES_SSLMODE", config.SSLMode)
	config.TimeZone = getEnv("POSTGRES_TIMEZONE", config.TimeZone)

	// Retry loop for database connection
	var err error
	maxRetries := 3
	retryDelay := 30 * time.Second

	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("Attempting to connect to database (attempt %d of %d)\n", attempt, maxRetries)

		DB, err = pgconnect.New(config)
		if err == nil {
			fmt.Println("Successfully connected to database")
			break
		}

		fmt.Printf("Failed to connect to database: %v\n", err)

		if attempt < maxRetries {
			fmt.Printf("Retrying in %v...\n", retryDelay)
			time.Sleep(retryDelay)
		} else {
			panic("Failed to connect to database after maximum retry attempts")
		}
	}

	// Auto migrate the Task model
	DB.AutoMigrate(&Task{})
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
