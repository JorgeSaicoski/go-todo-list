package main

import (
	"net/http"

	"github.com/JorgeSaicoski/go-todo-list/internal/api"
	"github.com/JorgeSaicoski/go-todo-list/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDatabase()

	router := gin.Default()
	router.LoadHTMLGlob("public/*")

	router.GET("/tasks", api.GetTasks)
	router.POST("/task", api.CreateTask)
	router.PATCH("/task/update/:id", api.UpdateTask)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Tasks",
		})
	})

	router.Run(":8000")
}
