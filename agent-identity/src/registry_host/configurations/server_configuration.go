package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/mvc"
	mvc_middlewares "github.com/6022-labs/agentic-orchestrator/agent-identity/src/common/mvc/middlewares"
	_ "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host/docs"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_host/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_mvc"
	fiberotel "github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/contrib/v3/swaggo"
	fiberzap "github.com/gofiber/contrib/v3/zap"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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

func newHttpServer(hostSettings *settings.HostSettings, errorWriter *mvc.ErrorWriter) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:      AppName,
		BodyLimit:    hostSettings.BodyLimitBytes,
		ErrorHandler: errorWriter.Write,
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
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// Spans only: ApiRequestMetrics owns http.server.*.
	p.App.Use(fiberotel.New(fiberotel.WithoutMetrics(true)))
	p.App.Use(p.RequestMetrics.Handle)

	p.App.Use(fiberzap.New(fiberzap.Config{
		Logger: p.Logger,
	}))

	p.App.Get("/swagger/*", swaggo.New())

	for _, controller := range p.Controllers {
		controller.RegisterRoutes(p.App)
	}
}
