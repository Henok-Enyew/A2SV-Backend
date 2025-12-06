package delivery

import (
	"task9/delivery/http"
	"task9/delivery/middleware"
	"task9/infrastructure"
	"task9/repository"
	"task9/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	taskRepo := repository.NewTaskRepositoryMongo(infrastructure.TaskCollection)
	taskUseCase := usecase.NewTaskUseCase(taskRepo)
	taskHandler := http.NewTaskHandler(taskUseCase)

	userRepo := repository.NewUserRepositoryMongo(infrastructure.UserCollection)
	passwordHasher := infrastructure.NewBcryptHasher()
	tokenGenerator := infrastructure.NewJWTGenerator()
	authUseCase := usecase.NewAuthUseCase(userRepo, passwordHasher, tokenGenerator)
	authHandler := http.NewAuthHandler(authUseCase)

	authMiddleware := middleware.NewAuthMiddleware(tokenGenerator)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Task Management API",
			"version": "1.0.0",
		})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	protected := r.Group("/")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.GET("/tasks", taskHandler.GetAllTasks)
		protected.GET("/tasks/:id", taskHandler.GetTaskByID)

		admin := protected.Group("/")
		admin.Use(authMiddleware.RequireAdmin())
		{
			admin.POST("/tasks", taskHandler.CreateTask)
			admin.PUT("/tasks/:id", taskHandler.UpdateTask)
			admin.DELETE("/tasks/:id", taskHandler.DeleteTask)
			admin.POST("/promote", authHandler.PromoteUser)
		}
	}

	return r
}

