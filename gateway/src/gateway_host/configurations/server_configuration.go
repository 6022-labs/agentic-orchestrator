package configurations

import (
	"time"

	mvc_middlewares "github.com/6022-labs/agentic-orchestrator/gateway/src/common/mvc/middlewares"
	gateway_settings "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/settings"
	_ "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_host/docs"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_host/settings"
	orchestrator_mvc "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_mvc"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_mvc/middlewares"
	fiberotel "github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/contrib/v3/swaggo"
	fiberzap "github.com/gofiber/contrib/v3/zap"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const AppName = "6022-Orchestrator-Gateway"

func ConfigureServer(container *dig.Container) {
	container.Provide(newHttpServer)

	err := container.Invoke(registerRoutes)
	if err != nil {
		panic(err)
	}
}

func newHttpServer(hostSettings *settings.HostSettings) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:   AppName,
		BodyLimit: hostSettings.BodyLimitBytes,
	})
}

type registerRoutesParams struct {
	dig.In

	App                               *fiber.App
	Logger                            *zap.Logger
	RateLimitSettings                 *gateway_settings.RateLimitSettings
	AuthenticationMandatoryMiddleware *middlewares.AuthenticationMandatoryMiddleware
	RequestMetrics                    *mvc_middlewares.ApiRequestMetricsMiddleware
	Controllers                       []orchestrator_mvc.ControllerInterface         `group:"controllers"`
	AuthInjectionMiddlewares          []middlewares.AuthInjectionMiddlewareInterface `group:"auth-injection-middlewares"`
}

// RegisterRoutes hooks up the routes and uses Fx to create new controller instances per request
func registerRoutes(p registerRoutesParams) {
	p.App.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	// Spans only: ApiRequestMetrics owns http.server.*.
	p.App.Use(fiberotel.Middleware(fiberotel.WithoutMetrics(true)))
	p.App.Use(p.RequestMetrics.Handle)

	p.App.Use(fiberzap.New(fiberzap.Config{
		Logger: p.Logger,
	}))

	p.App.Use(limiter.New(limiter.Config{
		Max:        p.RateLimitSettings.RequestsPerSecond,
		Expiration: 1 * time.Second,
	}))

	p.App.Get("/swagger/*", swaggo.New())

	for _, controller := range p.Controllers {
		controller.InitUnauthenticatedRoutes(p.App)
	}

	for _, authInjectionMiddleware := range p.AuthInjectionMiddlewares {
		p.App.Use(authInjectionMiddleware.Handle)
	}

	p.App.Use(p.AuthenticationMandatoryMiddleware.Handle)

	for _, controller := range p.Controllers {
		controller.InitAuthenticatedRoutes(p.App)
	}
}
