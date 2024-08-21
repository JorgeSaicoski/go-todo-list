package api

import (
        "net/http"
        "fmt"
        "github.com/gin-gonic/gin"
        "github.com/JorgeSaicoski/go-todo-list/internal/db"
)

func GetTasks(c *gin.Context) {
        var tasks []db.Task
        db.DB.Find(&tasks)
        c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
        var task db.Task
        
        if err := c.BindJSON(&task); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return  
        }

        db.DB.Create(&task)
        c.JSON(http.StatusOK, task)  
}

func UpdateTask(c *gin.Context){
        id := c.Param("id")
        var task db.Task

        if err := db.DB.First(&task, id).Error; err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
                return
        }

        var updateData := make(map[string]interface{})
        
        if updateData."Title" != ""{
                updateData["Title"] = updateData."Title"
        }
        if updateData."Content" != ""{
                updateData["Content"] = updateData."Content"
        }

        if len(updateData) > 0 {
                if err := db.DB.Model(&task).Updates(updateData).Error; err != nil{
                        c.JSON(http.StatusInternalServerError, gin.H {"error": err.Error()})
                        return
                }
        } else {
                c.JSON(http.StatusOK, gin.H{"message": "No changes detected"})
                return
        }
        c.JSON(http.StatusOK, task)
}

func FinishTask(c *gin.Context){
        id := c.Param("id")
        var task db.Task
        if err := db.DB.First(&task, id).Error; err != nil {
                c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
                return
        }
        
        task.Complete = true
        if err := db.DB.Save(&task).Error; err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }
        c.JSON(http.StatusOK, task)
}