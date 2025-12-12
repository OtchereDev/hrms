package frappe

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents a Frappe-compatible API response
// Frappe always returns responses in this format:
// { "message": <data>, "exc": <error>, "exc_type": <error_type> }
type Response struct {
	Message interface{} `json:"message"`
	Exc     interface{} `json:"exc"`
	ExcType string      `json:"exc_type,omitempty"`
}

// ValidationError represents a validation error in Frappe format
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorType constants matching Frappe's error types
const (
	ErrTypeValidation      = "ValidationError"
	ErrTypePermission      = "PermissionError"
	ErrTypeNotFound        = "DoesNotExistError"
	ErrTypeDuplicate       = "DuplicateEntryError"
	ErrTypeMandatory       = "MandatoryError"
	ErrTypeLink            = "LinkValidationError"
	ErrTypeValue           = "ValueError"
	ErrTypeInternal        = "InternalError"
	ErrTypeAuthentication  = "AuthenticationError"
	ErrTypeTimeout         = "TimeoutError"
)

// Success creates a successful Frappe response
func Success(data interface{}) *Response {
	return &Response{
		Message: data,
		Exc:     nil,
	}
}

// Error creates an error Frappe response
func Error(message string, excType string) *Response {
	return &Response{
		Message: nil,
		Exc:     message,
		ExcType: excType,
	}
}

// ErrorWithData creates an error response with additional data
func ErrorWithData(message string, excType string, data interface{}) *Response {
	return &Response{
		Message: data,
		Exc:     message,
		ExcType: excType,
	}
}

// ValidationErrors creates a validation error response
func ValidationErrors(errors []ValidationError) *Response {
	return &Response{
		Message: nil,
		Exc:     errors,
		ExcType: ErrTypeValidation,
	}
}

// SendSuccess sends a successful Frappe response
func SendSuccess(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Success(data))
}

// SendError sends an error Frappe response with appropriate status code
func SendError(c *fiber.Ctx, statusCode int, message string, excType string) error {
	return c.Status(statusCode).JSON(Error(message, excType))
}

// SendValidationError sends a validation error response
func SendValidationError(c *fiber.Ctx, errors []ValidationError) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(ValidationErrors(errors))
}

// SendNotFound sends a not found error
func SendNotFound(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusNotFound, message, ErrTypeNotFound)
}

// SendPermissionError sends a permission error
func SendPermissionError(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusForbidden, message, ErrTypePermission)
}

// SendAuthenticationError sends an authentication error
func SendAuthenticationError(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusUnauthorized, message, ErrTypeAuthentication)
}

// SendInternalError sends an internal server error
func SendInternalError(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusInternalServerError, message, ErrTypeInternal)
}

// SendDuplicateError sends a duplicate entry error
func SendDuplicateError(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusConflict, message, ErrTypeDuplicate)
}

// ListResponse represents a Frappe list response with pagination
type ListResponse struct {
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"total_count"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
}

// SendList sends a paginated list response
func SendList(c *fiber.Ctx, data interface{}, totalCount int64, page, pageSize int) error {
	return SendSuccess(c, ListResponse{
		Data:       data,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	})
}

// SendCreated sends a created response (201)
func SendCreated(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Success(data))
}

// SendBadRequest sends a bad request error (400)
func SendBadRequest(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusBadRequest, message, ErrTypeValidation)
}

