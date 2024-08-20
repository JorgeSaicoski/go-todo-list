package api

import (
        "net/http"
        "fmt"
        "github.com/gin-gonic/gin"
        "github.com/JorgeSaicoski/go-todo-list/internal/db"
)

func GetTasks(c *gin.Context) {
        var tasks []db.Task
        fmt.Printf("chego")	
        db.DB.Find(&tasks)
        c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
        var task db.Task
        
        if err := c.BindJSON(&task); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return  
        }
	fmt.Printf("chego")	

        db.DB.Create(&task)
        c.JSON(http.StatusOK, task)  
}