package registry_mvc

import "github.com/gofiber/fiber/v3"

type ControllerInterface interface {
	RegisterRoutes(app fiber.Router)
}
