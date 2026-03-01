package images_generator_mvc

import "github.com/gofiber/fiber/v2"

type ControllerInterface interface {
	RegisterRoutes(app fiber.Router)
}
