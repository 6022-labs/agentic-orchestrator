package mvc

import (
	"errors"

	common_errors "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/errors"
	"github.com/gofiber/fiber/v3"
)

// WriteError maps a use-case error to its HTTP status; an unknown error collapses to 500.
func WriteError(ctx fiber.Ctx, err error) error {
	var validationError *common_errors.ValidationError
	if errors.As(err, &validationError) {
		return ctx.Status(fiber.StatusBadRequest).JSON(validationError)
	}
	return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
}
