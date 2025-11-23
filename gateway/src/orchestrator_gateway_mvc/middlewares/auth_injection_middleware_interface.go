package middlewares

import "github.com/gofiber/fiber/v2"

type AuthInjectionMiddlewareInterface interface {
	Handle(c *fiber.Ctx) error
}
