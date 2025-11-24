package router

import (
	"task6/controllers"
	"task6/data"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	taskService := data.NewTaskService(data.TaskCollection)
	taskController := controllers.NewTaskController(taskService)

	r.GET("/tasks", taskController.GetAllTasks)
	r.GET("/tasks/:id", taskController.GetTaskByID)
	r.POST("/tasks", taskController.CreateTask)
	r.PUT("/tasks/:id", taskController.UpdateTask)
	r.DELETE("/tasks/:id", taskController.DeleteTask)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Task Management API",
			"version": "1.0.0",
		})
	})

	return r
}

