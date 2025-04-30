package api

import (
	"fmt"
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
	fmt.Printf("UpdateTask: Processing request for task ID: %s\n", id)

	var task db.Task

	// Find by ID using repository
	if err := TaskRepository.FindByID(id, &task); err != nil {
		fmt.Printf("UpdateTask: Error finding task with ID %s: %v\n", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	fmt.Printf("UpdateTask: Found task: %+v\n", task)

	// Get update data from JSON
	updateData := make(map[string]interface{})
	if err := c.BindJSON(&updateData); err != nil {
		fmt.Printf("UpdateTask: Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("UpdateTask: Received update data: %+v\n", updateData)

	// Validate fields
	fmt.Println("UpdateTask: Validating and applying fields...")
	if title, ok := updateData["title"].(string); ok && title != "" {
		newTitle := title
		fmt.Printf("UpdateTask: Updating title from to '%s'\n", title)
		task.Title = newTitle
	}

	if description, ok := updateData["description"].(string); ok {
		newDesc := description
		fmt.Printf("UpdateTask: Updating description to '%s'\n", newDesc)
		task.Description = &description
	}

	if status, ok := updateData["status"].(string); ok {
		newStatus := status
		fmt.Printf("UpdateTask: Updating status from '%s'\n", newStatus)
		task.Status = &status
	}

	if dueDate, ok := updateData["dueDate"].(*time.Time); ok {
		newDueDate := dueDate.String()
		fmt.Printf("UpdateTask: Updating due date from to '%s'\n", newDueDate)
		task.DueDate = dueDate
	}

	// Update using repository
	fmt.Printf("UpdateTask: Saving updated task: %+v\n", task)
	if err := TaskRepository.Update(&task); err != nil {
		fmt.Printf("UpdateTask: Error updating task: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("UpdateTask: Task successfully updated")
	c.JSON(http.StatusOK, task)
}
