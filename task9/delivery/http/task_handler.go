package http

import (
	"net/http"
	"task9/domain"
	"task9/usecase"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskUseCase *usecase.TaskUseCase
}

func NewTaskHandler(taskUseCase *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{taskUseCase: taskUseCase}
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
	var reqDTO CreateTaskRequest

	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	req := domain.CreateTaskRequest{
		Title:       reqDTO.Title,
		Description: reqDTO.Description,
		DueDate:     reqDTO.DueDate,
		Status:      reqDTO.Status,
	}

	task, err := h.taskUseCase.CreateTask(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid status" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{
			"status":  "error",
			"message": err.Error(),
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

	var reqDTO UpdateTaskRequest
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	req := domain.UpdateTaskRequest{
		Title:       reqDTO.Title,
		Description: reqDTO.Description,
		DueDate:     reqDTO.DueDate,
		Status:      reqDTO.Status,
	}

	task, err := h.taskUseCase.UpdateTask(id, req)
	if err != nil {
		statusCode := http.StatusNotFound
		if err.Error() == "invalid task ID format" || err.Error() == "invalid status" {
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

