package configurations

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_internal_mvc"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	Controllers []slack_bridge_internal_mvc.ControllerInterface `group:"controllers"`
}

func useRestApi(p useRestApiParams) {
	p.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	p.App.Use(fiberzap.New(fiberzap.Config{
		Logger: p.Logger,
	}))

	for _, controller := range p.Controllers {
		controller.InitRoutes(p.App)
	}
}
