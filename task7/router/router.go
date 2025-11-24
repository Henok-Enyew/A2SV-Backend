package router

import (
	"task7/controllers"
	"task7/data"
	"task7/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	taskService := data.NewTaskService(data.TaskCollection)
	taskController := controllers.NewTaskController(taskService)

	userService := data.NewUserService(data.UserCollection)
	authController := controllers.NewAuthController(userService)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Task Management API",
			"version": "1.0.0",
		})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/tasks", taskController.GetAllTasks)
		protected.GET("/tasks/:id", taskController.GetTaskByID)

		admin := protected.Group("/")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/tasks", taskController.CreateTask)
			admin.PUT("/tasks/:id", taskController.UpdateTask)
			admin.DELETE("/tasks/:id", taskController.DeleteTask)
			admin.POST("/promote", authController.PromoteUser)
		}
	}

	return r
}

