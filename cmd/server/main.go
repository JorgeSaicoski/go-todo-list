package main

import (

        "github.com/gin-gonic/gin"
        "github.com/JorgeSaicoski/go-todo-list/internal/db"
        "github.com/JorgeSaicoski/go-todo-list/internal/api"
        "net/http"
)

func main() {
        db.ConnectDatabase()
        

        router := gin.Default()
        router.LoadHTMLGlob("public/*")

        router.GET("/tasks", api.GetTasks)
        router.POST("/tasks", api.CreateTask)
        router.GET("/", func(c *gin.Context) {
                c.HTML(http.StatusOK, "index.tmpl", gin.H{
                    "title": "Tasks",
                })
            })


        router.Run(":8000")
}
