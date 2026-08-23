package configurations

import (
	mvc_middlewares "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/mvc/middlewares"
	_ "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/docs"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_host/settings"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_internal_mvc"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const AppName = "6022-SlackBridge-Internal"

func ConfigureServer(container *dig.Container) {
	container.Provide(newHttpServer)

	err := container.Invoke(useRestApi)
	if err != nil {
		panic(err)
	}
}

func newHttpServer(hostSettings *settings.HostSettings) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:   AppName,
		Network:   "tcp",
		BodyLimit: hostSettings.BodyLimitBytes,
	})
}

type useRestApiParams struct {
	dig.In

	App            *fiber.App
	Logger         *zap.Logger
	RequestMetrics *mvc_middlewares.ApiRequestMetricsMiddleware
	Controllers    []slack_bridge_internal_mvc.ControllerInterface `group:"controllers"`
}

func useRestApi(p useRestApiParams) {
	p.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	// Spans only: ApiRequestMetrics owns http.server.*.
	p.App.Use(otelfiber.Middleware(otelfiber.WithoutMetrics(true)))
	p.App.Use(p.RequestMetrics.Handle)

	p.App.Use(fiberzap.New(fiberzap.Config{
		Logger: p.Logger,
	}))

	p.App.Get("/swagger/*", swagger.New())

	for _, controller := range p.Controllers {
		controller.InitRoutes(p.App)
	}
}
