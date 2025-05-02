package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/JorgeSaicoski/go-todo-list/internal/api"
	"github.com/JorgeSaicoski/go-todo-list/internal/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDatabase()

	api.InitRepository(db.DB)

	router := gin.Default()

	// Get allowed origins from environment variable
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	origins := strings.Split(allowedOrigins, ",")

	// Configure CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	router.LoadHTMLGlob("public/*")

	router.GET("/tasks", api.GetTasks)
	router.POST("/task", api.CreateTask)
	router.PATCH("/task/update/:id", api.UpdateTask)
	router.GET("/tasks/active", api.GetNonCompletedTasksPaginated)
	router.GET("/tasks/completed", api.GetCompletedTasksPaginated)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Tasks",
		})
	})

	router.Run(":8000")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
