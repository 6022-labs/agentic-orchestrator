package configurations

import (
	_ "github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_host/docs"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_host/settings"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector_mvc"
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
		AppName: hostSettings.AppName,
		Network: "tcp",
	})
}

type useRestApiParams struct {
	dig.In

	App         *fiber.App
	Logger      *zap.Logger
	Controllers []llm_connector_mvc.ControllerInterface `group:"controllers"`
}

// UseRestApi hooks up the routes and uses Fx to create new controller instances per request
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
		controller.InitRoutes(p.App)
	}
}
