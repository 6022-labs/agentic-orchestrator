package slack_bridge_external_mvc

import "github.com/gofiber/fiber/v3"

type ControllerInterface interface {
	InitRoutes(app fiber.Router)
}
