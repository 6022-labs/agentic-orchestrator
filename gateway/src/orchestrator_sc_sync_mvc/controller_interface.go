package orchestrator_sc_sync_mvc

import "github.com/gofiber/fiber/v2"

type ControllerInterface interface {
	InitRoutes(app fiber.Router)
}
