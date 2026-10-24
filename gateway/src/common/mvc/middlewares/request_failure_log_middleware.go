package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// RequestFailureLogMiddleware logs every 4xx (except 404) and 5xx response with its body.
type RequestFailureLogMiddleware struct {
	logger *zap.Logger
}

func NewRequestFailureLogMiddleware(logger *zap.Logger) *RequestFailureLogMiddleware {
	return &RequestFailureLogMiddleware{logger: logger}
}

func (m *RequestFailureLogMiddleware) Handle(ctx fiber.Ctx) error {
	err := ctx.Next()
	status := ctx.Response().StatusCode()

	if status < fiber.StatusBadRequest || status == fiber.StatusNotFound {
		return err
	}
	fields := []zap.Field{
		zap.String("method", ctx.Method()),
		zap.String("path", ctx.OriginalURL()),
		zap.Int("status", status),
		zap.ByteString("response", ctx.Response().Body()),
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	if status >= fiber.StatusInternalServerError {
		m.logger.Error("gateway request failed", fields...)
	} else {
		m.logger.Warn("gateway request rejected", fields...)
	}
	return err
}
