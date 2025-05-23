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

// Handler encapsulates all the task-related API handlers
type TaskHandler struct {
	repo *pgconnect.Repository[db.Task]
}

// NewTaskHandler creates and returns a new TaskHandler instance
func NewTaskHandler(database *pgconnect.DB) *TaskHandler {
	return &TaskHandler{
		repo: pgconnect.NewRepository[db.Task](database),
	}
}

// paginateUserTasks handles paginated queries for user-specific tasks with optional conditions
func (h *TaskHandler) paginateUserTasks(c *gin.Context, additionalCondition string, args ...interface{}) {
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
	fmt.Println(pageSize)

	if pageSizeParam := c.Query("pageSize"); pageSizeParam != "" {
		if pageSizeVal, err := strconv.Atoi(pageSizeParam); err == nil && pageSizeVal > 0 {
			pageSize = pageSizeVal
		}
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	fmt.Println(userID)

	// Build the where condition
	whereCondition := "ownerId = ?"
	whereArgs := []interface{}{userID}

	if additionalCondition != "" {
		whereCondition += " AND " + additionalCondition
		whereArgs = append(whereArgs, args...)
	}
	fmt.Println("whereArgs")
	fmt.Println(whereArgs)

	// Get user's tasks with pagination
	if err := h.repo.PaginateWhere(&tasks, page, pageSize, whereCondition, whereArgs...); err != nil {
		fmt.Println("error in the paginate where")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the total count
	var count int64
	if err := h.repo.Count(&count, whereCondition, whereArgs...); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the paginated tasks
	c.JSON(http.StatusOK, gin.H{
		"tasks":      tasks,
		"total":      count,
		"page":       page,
		"pageSize":   pageSize,
		"totalPages": int(math.Ceil(float64(count) / float64(pageSize))),
	})
}

func (h *TaskHandler) GetTasksPaginated(c *gin.Context) {
	fmt.Println("paginated tasks______________________________________")

	h.paginateUserTasks(c, "", nil)
}

func (h *TaskHandler) GetNonCompletedTasksPaginated(c *gin.Context) {
	fmt.Println("non completed tasks______________________________________")

	h.paginateUserTasks(c, "status != ?", "completed")
}

func (h *TaskHandler) GetCompletedTasksPaginated(c *gin.Context) {
	fmt.Println("completed tasks______________________________________")
	h.paginateUserTasks(c, "status = ?", "completed")
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task db.Task

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user ID from the context (set by AuthMiddleware)
	userID, exists := c.Get("userID")
	fmt.Println("userID")
	fmt.Println(userID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Set the owner ID
	task.OwnerID = userID.(string)

	if err := h.repo.Create(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(task)
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var task db.Task

	// Find by ID using repository
	if err := h.repo.FindByID(id, &task); err != nil {
		fmt.Printf("UpdateTask: Error finding task with ID %s: %v\n", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Verify the user is the owner
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if task.OwnerID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this task"})
		return
	}

	// Get update data from JSON
	updateData := make(map[string]interface{})
	if err := c.BindJSON(&updateData); err != nil {
		fmt.Printf("UpdateTask: Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate fields
	if title, ok := updateData["title"].(string); ok && title != "" {
		newTitle := title
		task.Title = newTitle
	}

	if description, ok := updateData["description"].(string); ok {
		task.Description = &description
	}

	if status, ok := updateData["status"].(string); ok {
		task.Status = &status
	}

	if dueDate, ok := updateData["dueDate"].(*time.Time); ok {
		task.DueDate = dueDate
	}

	// Update using repository
	if err := h.repo.Update(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteSelectedTasks(c *gin.Context) {
	var taskIDs []string
	if err := c.BindJSON(&taskIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, id := range taskIDs {
		var task db.Task
		if err := h.repo.FindByID(id, &task); err != nil {
			fmt.Println("Skip tasks that don't exist, while deleting")
			continue // Skip tasks that don't exist
		}
		if err := h.repo.Delete(&task); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tasks deleted successfully"})
}

func (h *TaskHandler) DeleteAllCompletedTasks(c *gin.Context) {
	if err := h.repo.DeleteWhere("status = ?", "completed"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All completed tasks deleted successfully"})
}

func (h *TaskHandler) DeleteAllNonCompletedTasks(c *gin.Context) {
	if err := h.repo.DeleteWhere("status != ?", "completed"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All non-completed tasks deleted successfully"})
}
