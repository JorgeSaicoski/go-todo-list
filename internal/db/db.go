package db

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Task struct {
	gorm.Model
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"` // expecting "pending", "in-progress", or "completed"
	DueDate     *time.Time `json:"dueDate,omitempty"`
}

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=yourpassword dbname=taskdb port=5432 sslmode=disable TimeZone=UTC"
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&Task{})
}
