package configurations

import (
	_ "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host/docs"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_mvc"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func ConfigureServer(container *dig.Container) {
	container.Provide(newHttpServer)

	err := container.Invoke(useRestApi)
	if err != nil {
		panic(err)
	}
}

func newHttpServer(hostSettings *settings.HostSettings) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:   hostSettings.AppName,
		Network:   "tcp",
		BodyLimit: hostSettings.BodyLimitBytes,
	})
}

type useRestApiParams struct {
	dig.In

	App         *fiber.App
	Logger      *zap.Logger
	Controllers []registry_mvc.ControllerInterface `group:"controllers"`
}

// RegisterRoutes hooks up the routes and uses Fx to create new controller instances per request
func useRestApi(p useRestApiParams) {
	p.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	p.App.Use(fiberzap.New(fiberzap.Config{
		Logger: p.Logger,
	}))

	p.App.Get("/swagger/*", swagger.New())

	for _, controller := range p.Controllers {
		controller.RegisterRoutes(p.App)
	}
}
