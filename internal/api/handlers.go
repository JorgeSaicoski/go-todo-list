package api

import (
	"net/http"
	"time"

	"github.com/JorgeSaicoski/go-todo-list/internal/db"
	"github.com/JorgeSaicoski/pgconnect"
	"github.com/gin-gonic/gin"
)

// TaskRepository is a repository for Task models
var TaskRepository *pgconnect.Repository[db.Task]

// InitRepository initializes the task repository
func InitRepository(database *pgconnect.DB) {
	TaskRepository = pgconnect.NewRepository[db.Task](database)
}

func GetTasks(c *gin.Context) {
	var tasks []db.Task
	if err := TaskRepository.FindAll(&tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var task db.Task

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := TaskRepository.Create(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task db.Task

	// Find by ID using repository
	if err := TaskRepository.FindByID(id, &task); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Get update data from JSON
	updateData := make(map[string]interface{})
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate fields
	if title, ok := updateData["Title"].(string); ok && title != "" {
		task.Title = title
	}
	if description, ok := updateData["Description"].(*string); ok {
		task.Description = description
	}
	if status, ok := updateData["Status"].(*string); ok {
		task.Status = status
	}
	if dueDate, ok := updateData["DueDate"].(*time.Time); ok {
		task.DueDate = dueDate
	}

	// Update using repository
	if err := TaskRepository.Update(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}
