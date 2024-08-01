package main

import (

        "github.com/gin-gonic/gin"
        "github.com/JorgeSaicoski/go-todo-list/internal/db"
        "github.com/JorgeSaicoski/go-todo-list/internal/api"
)

func main() {
        db.ConnectDatabase()
        

        router := gin.Default()

        router.GET("/tasks", api.GetTasks)
        router.POST("/tasks", api.CreateTask)
		router.GET("/", func(c *gin.Context) {
			c.File("public/index.html")
		})


        router.Run(":8080")
}
