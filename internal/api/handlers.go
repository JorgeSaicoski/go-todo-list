package api

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
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

func GetTasksPaginate(c *gin.Context) {
	var tasks []db.Task

	// Set default pagination values
	page := 1
	pageSize := 10

	// Get pagination data from query parameters
	if pageParam := c.Query("page"); pageParam != "" {
		if pageVal, err := strconv.Atoi(pageParam); err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	if pageSizeParam := c.Query("pageSize"); pageSizeParam != "" {
		if pageSizeVal, err := strconv.Atoi(pageSizeParam); err == nil && pageSizeVal > 0 {
			pageSize = pageSizeVal
		}
	}
	if err := TaskRepository.Paginate(&tasks, page, pageSize); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var count int64
	if err := TaskRepository.Count(&count, "status != ?", "completed"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return both the tasks and pagination information
	c.JSON(http.StatusOK, gin.H{
		"tasks":      tasks,
		"total":      count,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int(math.Ceil(float64(count) / float64(pageSize))),
	})

}

func GetNonCompletedTasksPaginated(c *gin.Context) {
	var tasks []db.Task

	// Set default pagination values
	page := 1
	pageSize := 10

	// Get pagination data from query parameters
	if pageParam := c.Query("page"); pageParam != "" {
		if pageVal, err := strconv.Atoi(pageParam); err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	if pageSizeParam := c.Query("pageSize"); pageSizeParam != "" {
		if pageSizeVal, err := strconv.Atoi(pageSizeParam); err == nil && pageSizeVal > 0 {
			pageSize = pageSizeVal
		}
	}

	// Using direct DB access to combine WHERE clause with pagination
	if err := TaskRepository.PaginateWhere(&tasks, page, pageSize, "status != ?", "completed"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var count int64
	if err := TaskRepository.Count(&count, "status != ?", "completed"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return both the tasks and pagination information
	c.JSON(http.StatusOK, gin.H{
		"tasks":      tasks,
		"total":      count,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int(math.Ceil(float64(count) / float64(pageSize))),
	})
}

func GetCompletedTasksPaginated(c *gin.Context) {
	var tasks []db.Task

	// Set default pagination values
	page := 1
	pageSize := 10 // Define the page size

	// Get pagination data from query parameters
	if pageParam := c.Query("page"); pageParam != "" {
		if pageVal, err := strconv.Atoi(pageParam); err == nil && pageVal > 0 {
			page = pageVal
		}
	}

	if pageSizeParam := c.Query("pageSize"); pageSizeParam != "" {
		if pageSizeVal, err := strconv.Atoi(pageSizeParam); err == nil && pageSizeVal > 0 {
			pageSize = pageSizeVal
		}
	}

	// Using direct DB access to combine WHERE clause with pagination
	if err := TaskRepository.PaginateWhere(&tasks, page, pageSize, "status = ?", "completed"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var count int64
	if err := TaskRepository.Count(&count, "status = ?", "completed"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return both the tasks and pagination information
	c.JSON(http.StatusOK, gin.H{
		"tasks":      tasks,
		"total":      count,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int(math.Ceil(float64(count) / float64(pageSize))),
	})
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

func DeleteSelectedTasks(c *gin.Context) {
	var taskIDs []string
	if err := c.BindJSON(&taskIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, id := range taskIDs {
		var task db.Task
		if err := TaskRepository.FindByID(id, &task); err != nil {
			continue // Skip tasks that don't exist
		}
		if err := TaskRepository.Delete(&task); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tasks deleted successfully"})
}

func DeleteAllCompletedTasks(c *gin.Context) {
	if err := TaskRepository.DeleteWhere("status = ?", "completed"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All completed tasks deleted successfully"})
}

func DeleteAllNonCompletedTasks(c *gin.Context) {
	if err := TaskRepository.DeleteWhere("status != ?", "completed"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All non-completed tasks deleted successfully"})
}
