package configurations

import (
	mvc_middlewares "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/mvc/middlewares"
	_ "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host/docs"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_mvc"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const AppName = "6022-AgentIdentity-Registry"

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
	Controllers    []registry_mvc.ControllerInterface `group:"controllers"`
}

// useRestApi registers middleware and routes on the Fiber app using controllers resolved from the dig container.
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
		controller.RegisterRoutes(p.App)
	}
}
