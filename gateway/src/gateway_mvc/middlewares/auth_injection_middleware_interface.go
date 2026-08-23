package middlewares

import "github.com/gofiber/fiber/v3"

type AuthInjectionMiddlewareInterface interface {
	Handle(c fiber.Ctx) error
}
