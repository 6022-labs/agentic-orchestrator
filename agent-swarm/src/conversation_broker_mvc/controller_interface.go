package conversation_broker_mvc

import "github.com/gofiber/fiber/v3"

type ControllerInterface interface {
	InitRoutes(app fiber.Router)
}
