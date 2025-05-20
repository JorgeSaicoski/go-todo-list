package api

import (
	"strings"
	"time"

	"github.com/JorgeSaicoski/pgconnect"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// TaskRouter handles routing for task-related endpoints
type TaskRouter struct {
	handler *TaskHandler
	router  *gin.Engine
}

// RouterConfig holds configuration for the router
type RouterConfig struct {
	AllowedOrigins   string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
	TemplatesGlob    string
}

// DefaultRouterConfig returns a default router configuration
func DefaultRouterConfig() RouterConfig {
	return RouterConfig{
		AllowedOrigins:   "http://localhost:3000",
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
		TemplatesGlob:    "public/*",
	}
}

// NewTaskRouter creates a new TaskRouter instance with full configuration
func NewTaskRouter(database *pgconnect.DB, config RouterConfig) *TaskRouter {
	router := gin.Default()

	// Configure CORS middleware
	origins := strings.Split(config.AllowedOrigins, ",")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     config.AllowedMethods,
		AllowHeaders:     config.AllowedHeaders,
		ExposeHeaders:    config.ExposeHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           time.Duration(config.MaxAge) * time.Second,
	}))

	// Load HTML templates if specified
	if config.TemplatesGlob != "" {
		router.LoadHTMLGlob(config.TemplatesGlob)
	}

	// Create the task handler
	taskHandler := NewTaskHandler(database)

	return &TaskRouter{
		handler: taskHandler,
		router:  router,
	}
}

// RegisterRoutes sets up all task-related routes
func (tr *TaskRouter) RegisterRoutes() {
	// Tasks endpoints
	tasksGroup := tr.router.Group("/tasks")
	{
		// Get all tasks with pagination
		tasksGroup.GET("", tr.handler.GetTasksPaginated)

		// Create a new task
		tasksGroup.POST("", tr.handler.CreateTask)

		// Update a task
		tasksGroup.PATCH("/update/:id", tr.handler.UpdateTask)

		// Get active (non-completed) tasks
		tasksGroup.GET("/active", tr.handler.GetNonCompletedTasksPaginated)

		// Get completed tasks
		tasksGroup.GET("/completed", tr.handler.GetCompletedTasksPaginated)

		// Delete tasks endpoints
		tasksGroup.POST("/delete", tr.handler.DeleteSelectedTasks)
		tasksGroup.DELETE("/completed", tr.handler.DeleteAllCompletedTasks)
		tasksGroup.DELETE("/active", tr.handler.DeleteAllNonCompletedTasks)
	}

	// Home route
	tr.router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"title": "Tasks",
		})
	})
}

// GetRouter returns the configured gin router
func (tr *TaskRouter) GetRouter() *gin.Engine {
	return tr.router
}

// Run starts the HTTP server
func (tr *TaskRouter) Run(addr string) error {
	return tr.router.Run(addr)
}
