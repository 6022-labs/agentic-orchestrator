package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
)

// ConfigureDaemon starts the conversation processor daemon and hooks its
// graceful shutdown to the Fiber app lifecycle.
func ConfigureDaemon(container *dig.Container) {
	err := container.Invoke(func(app *fiber.App, daemon interfaces.ConversationProcessorDaemonInterface) {
		daemon.Start()
		app.Hooks().OnShutdown(func() error {
			daemon.Stop()
			return nil
		})
	})
	if err != nil {
		panic(err)
	}
}
