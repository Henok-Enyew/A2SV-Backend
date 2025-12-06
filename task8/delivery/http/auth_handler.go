package http

import (
	"net/http"
	"task8/domain"
	"task8/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var reqDTO RegisterRequest

	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	req := domain.RegisterRequest{
		Username: reqDTO.Username,
		Password: reqDTO.Password,
	}

	user, err := h.authUseCase.Register(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "username already exists" {
			statusCode = http.StatusConflict
		} else if err.Error() == "password must be at least 6 characters" {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "user registered successfully",
		"data":    user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var reqDTO LoginRequest

	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	req := domain.LoginRequest{
		Username: reqDTO.Username,
		Password: reqDTO.Password,
	}

	token, user, err := h.authUseCase.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}

func (h *AuthHandler) PromoteUser(c *gin.Context) {
	var reqDTO PromoteRequest

	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err := h.authUseCase.PromoteUser(reqDTO.Username)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user promoted to admin successfully",
	})
}

