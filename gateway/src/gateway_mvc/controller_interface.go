package orchestrator_gateway_mvc

import "github.com/gofiber/fiber/v2"

type ControllerInterface interface {
	InitAuthenticatedRoutes(app fiber.Router)
	InitUnauthenticatedRoutes(app fiber.Router)
}
