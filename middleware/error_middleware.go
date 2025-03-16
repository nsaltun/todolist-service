package middleware

import (
	"errors"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Message   string    `json:"message"`
	RequestID string    `json:"request_id,omitempty"`
	Status    int       `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
}

func ErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		if err == nil {
			return nil
		}

		// Get request details
		reqID := c.Get("X-Request-ID")
		path := c.Path()
		method := c.Method()
		stack := string(debug.Stack())

		// Determine status code based on error type
		status := getErrorStatusCode(err)

		// Log error with details
		logger.Error("request error",
			zap.Error(err),
			zap.String("stack", stack),
			zap.String("request_id", reqID),
			zap.String("path", path),
			zap.String("method", method),
			zap.Int("status", status),
		)

		// Return clean error message to client
		return c.Status(status).JSON(ErrorResponse{
			Message:   cleanErrorMessage(err),
			RequestID: reqID,
			Status:    status,
			Timestamp: time.Now(),
			Path:      path,
			Method:    method,
		})
	}
}

func getErrorStatusCode(err error) int {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return fiber.StatusNotFound
	case strings.Contains(err.Error(), "duplicate key"):
		return fiber.StatusConflict
	case strings.Contains(err.Error(), "validation"):
		return fiber.StatusBadRequest
	default:
		return fiber.StatusInternalServerError
	}
}

func cleanErrorMessage(err error) string {
	msg := err.Error()

	// Clean sensitive information
	msg = strings.ReplaceAll(msg, "pq:", "")
	msg = strings.TrimSpace(msg)

	return msg
}
