package dto

import (
	"errors"
	"github.com/19parwiz/order/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPError struct {
	Code    int
	Message string
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func FromError(err error) *HTTPError {
	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}

func HandleError(ctx *gin.Context, err error) {
	var statusCode int
	response := ErrorResponse{}

	switch {
	case errors.Is(err, domain.ErrOrderNotFound):
		statusCode = http.StatusNotFound
		response.Error = "Order not found"
	case errors.Is(err, domain.ErrInvalidOrderData):
		statusCode = http.StatusBadRequest
		response.Error = "Invalid order data"
	case errors.Is(err, domain.ErrIDGenerationFailed):
		statusCode = http.StatusInternalServerError
		response.Error = "Failed to generate order ID"
	// Add more cases as needed
	default:
		statusCode = http.StatusInternalServerError
		response.Error = "Internal server error"
		response.Details = err.Error() // Include details in development mode
	}

	ctx.JSON(statusCode, response)
}
