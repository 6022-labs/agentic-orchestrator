package orchestrator_gateway_mvc

import "github.com/gofiber/fiber/v3"

const (
	PING_ENDPOINT = "/ping"
)

type PingController struct{}

func NewPingController() *PingController {
	return &PingController{}
}

func (controller *PingController) InitAuthenticatedRoutes(app fiber.Router) {
	// Just to satisfy the interface
}

func (controller *PingController) InitUnauthenticatedRoutes(app fiber.Router) {
	app.Get(PING_ENDPOINT, controller.Ping)
}

// Ping godoc
// @Summary      Health check endpoint
// @Description  Returns "pong" if the service is running
// @Tags         health
// @Produce      plain
// @Success      200 {string} string "pong"
// @Router       /ping [get]
func (controller *PingController) Ping(c fiber.Ctx) error {
	return c.SendString("pong")
}
