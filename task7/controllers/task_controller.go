package controllers

import (
	"net/http"
	"task7/data"
	"task7/models"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskService *data.TaskService
}

func NewTaskController(taskService *data.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

func (tc *TaskController) GetAllTasks(c *gin.Context) {
	tasks, err := tc.taskService.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to retrieve tasks",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   tasks,
		"count":  len(tasks),
	})
}

func (tc *TaskController) GetTaskByID(c *gin.Context) {
	id := c.Param("id")

	task, err := tc.taskService.GetTaskByID(id)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() == "invalid task ID format" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   task,
	})
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if req.Status != "" && req.Status != "pending" && req.Status != "in_progress" && req.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "status must be one of: pending, in_progress, completed",
		})
		return
	}

	task, err := tc.taskService.CreateTask(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to create task",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "task created successfully",
		"data":    task,
	})
}

func (tc *TaskController) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if req.Status != "" && req.Status != "pending" && req.Status != "in_progress" && req.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "status must be one of: pending, in_progress, completed",
		})
		return
	}

	task, err := tc.taskService.UpdateTask(id, req)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() == "invalid task ID format" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "task updated successfully",
		"data":    task,
	})
}

func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := tc.taskService.DeleteTask(id)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() == "invalid task ID format" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "task deleted successfully",
	})
}

