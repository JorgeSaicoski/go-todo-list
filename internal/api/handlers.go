package handlers

import (
        "net/http"

        "github.com/gin-gonic/gin"
        "github.com/JorgeSaicoski/go-todo-list/database"
        "github.com/JorgeSaicoski/go-todo-list/models"
)

func GetTasks(c *gin.Context) {
        var tasks []models.Task
        database.DB.Find(&tasks)
        c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
        var task models.Task
        if err := c.BindJSON(&task); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return   

        }

        database.DB.Create(&task)
        c.JSON(http.StatusOK, task)   

}