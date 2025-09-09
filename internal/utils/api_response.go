package utils

import (
	"time"
	"github.com/gofiber/fiber/v2"
)

// APIResponse represents a standardized API response format
type APIResponse struct {
	Status    string      `json:"status"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// PaginatedAPIResponse represents a paginated API response
type PaginatedAPIResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	Pagination Pagination  `json:"pagination,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	RequestID  string      `json:"request_id,omitempty"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Success creates a successful API response
func Success(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	response := APIResponse{
		Status:    "success",
		Message:   msg,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: c.Get("X-Request-ID", ""),
	}

	return c.JSON(response)
}

// SuccessWithPagination creates a successful paginated API response
func SuccessWithPagination(c *fiber.Ctx, data interface{}, pagination Pagination, message ...string) error {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	response := PaginatedAPIResponse{
		Status:     "success",
		Message:    msg,
		Data:       data,
		Pagination: pagination,
		Timestamp:  time.Now(),
		RequestID:  c.Get("X-Request-ID", ""),
	}

	return c.JSON(response)
}

// Error creates an error API response
func Error(c *fiber.Ctx, statusCode int, message string, err interface{}) error {
	response := APIResponse{
		Status:    "error",
		Message:   message,
		Error:     err,
		Timestamp: time.Now(),
		RequestID: c.Get("X-Request-ID", ""),
	}

	return c.Status(statusCode).JSON(response)
}

// BadRequest creates a bad request error response
func BadRequest(c *fiber.Ctx, message string, err interface{}) error {
	return Error(c, fiber.StatusBadRequest, message, err)
}

// Unauthorized creates an unauthorized error response
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message, nil)
}

// Forbidden creates a forbidden error response
func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message, nil)
}

// NotFound creates a not found error response
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, nil)
}

// InternalServerError creates an internal server error response
func InternalServerError(c *fiber.Ctx, message string, err interface{}) error {
	return Error(c, fiber.StatusInternalServerError, message, err)
}

// ValidationError creates a validation error response
func ValidationError(c *fiber.Ctx, errors interface{}) error {
	return Error(c, fiber.StatusUnprocessableEntity, "Validation failed", errors)
}

// Created creates a created response
func Created(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Resource created successfully"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	response := APIResponse{
		Status:    "success",
		Message:   msg,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: c.Get("X-Request-ID", ""),
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// NoContent creates a no content response
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, pageSize int, total int64) Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	hasNext := page < totalPages
	hasPrev := page > 1

	return Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}