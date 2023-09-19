package api

import (
	"github.com/gin-gonic/gin"
)

type successResponse struct {
	Message string `json:"message" example:"success"`
}

type errorResponse struct {
	Error string `json:"error" example:"error message"`
}

func writeSuccessResponse(c *gin.Context, statusCode int, msg string) {
	c.JSON(statusCode, successResponse{Message: msg})
}

func writeErrorResponse(c *gin.Context, statusCode int, msg string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{Error: msg})
}
