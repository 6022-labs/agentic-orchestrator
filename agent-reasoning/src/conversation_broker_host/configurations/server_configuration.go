package configurations

import (
	_ "github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_host/docs"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_host/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker_mvc"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const appName = "6022-AgentReasoning-ConversationBroker"

func ConfigureServer(container *dig.Container) {
	container.Provide(newHttpServer)

	err := container.Invoke(useRestApi)
	if err != nil {
		panic(err)
	}
}

func newHttpServer(hostSettings *settings.HostSettings) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:   appName,
		Network:   "tcp",
		BodyLimit: hostSettings.BodyLimitBytes,
	})
}

type useRestApiParams struct {
	dig.In

	App         *fiber.App
	Logger      *zap.Logger
	Controllers []conversation_broker_mvc.ControllerInterface `group:"controllers"`
}

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
