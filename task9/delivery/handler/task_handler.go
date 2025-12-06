package handler

import (
	"net/http"
	"task9/delivery/http"
	"task9/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskUseCase *usecase.TaskUseCase
}

func NewTaskHandler(taskUseCase *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{
		taskUseCase: taskUseCase,
	}
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.taskUseCase.GetAllTasks()
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

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id := c.Param("id")

	task, err := h.taskUseCase.GetTaskByID(id)
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

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req http.CreateTaskRequest

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

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid due date format, use RFC3339",
		})
		return
	}

	task, err := h.taskUseCase.CreateTask(req.Title, req.Description, dueDate, req.Status)
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

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var req http.UpdateTaskRequest
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

	var dueDate interface{}
	if req.DueDate != "" {
		parsed, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "invalid due date format, use RFC3339",
			})
			return
		}
		dueDate = parsed
	}

	task, err := h.taskUseCase.UpdateTask(id, req.Title, req.Description, dueDate, req.Status)
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

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := h.taskUseCase.DeleteTask(id)
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


