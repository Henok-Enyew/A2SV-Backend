package router

import (
	"task5/controllers"
	"task5/data"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	taskService := data.NewTaskService()
	taskController := controllers.NewTaskController(taskService)

	api := r.Group("/api/v1")
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("", taskController.GetAllTasks)
			tasks.GET("/:id", taskController.GetTaskByID)
			tasks.POST("", taskController.CreateTask)
			tasks.PUT("/:id", taskController.UpdateTask)
			tasks.DELETE("/:id", taskController.DeleteTask)
		}
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Task Management API",
			"version": "1.0.0",
		})
	})

	return r
}

