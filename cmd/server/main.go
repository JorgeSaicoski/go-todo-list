package main

import (
        "net/http"
        "github.com/gin-gonic/gin"
        "github.com/JorgeSaicoski/go-todo-list/db"
        "github.com/JorgeSaicoski/go-todo-list/api"
)

func main() {
        database.ConnectDatabase()

        router := gin.Default()

        router.GET("/tasks", handlers.GetTasks)
        router.POST("/tasks", handlers.CreateTask)
		router.GET("/", func(c *gin.Context) {
			c.File("public/index.html") // Replace "home.html" with your desired file
		})


        router.Run(":8080")
}
