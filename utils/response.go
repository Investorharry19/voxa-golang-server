package utils

import "github.com/gofiber/fiber/v2"

// Standard API response structure
type APIResponse struct {
	Success   bool        `json:"success"`              // "success" or "error"
	Message   string      `json:"message,omitempty"`    // readable human message
	ErrorCode int         `json:"error_code,omitempty"` // internal or business error
	Data      interface{} `json:"data,omitempty"`       // payload (optional)
}

// Success response
func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success:   false,
		Message:   message,
		ErrorCode: statusCode,
	})
}
