package configurations

import (
	_ "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_host/docs"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_host/settings"
	orchestrator_mvc "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_mvc"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_mvc/middlewares"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func ConfigureServer(container *dig.Container) {
	container.Provide(newHttpServer)

	err := container.Invoke(registerRoutes)
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

type registerRoutesParams struct {
	dig.In

	App                               *fiber.App
	Logger                            *zap.Logger
	AuthenticationMandatoryMiddleware *middlewares.AuthenticationMandatoryMiddleware
	Controllers                       []orchestrator_mvc.ControllerInterface         `group:"controllers"`
	AuthInjectionMiddlewares          []middlewares.AuthInjectionMiddlewareInterface `group:"auth-injection-middlewares"`
}

// RegisterRoutes hooks up the routes and uses Fx to create new controller instances per request
func registerRoutes(p registerRoutesParams) {
	p.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	p.App.Use(fiberzap.New(fiberzap.Config{
		Logger: p.Logger,
	}))

	p.App.Get("/swagger/*", swagger.New())

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
