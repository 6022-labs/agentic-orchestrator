package configurations

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_sc_sync_mvc"
	"go.uber.org/dig"
)

func AddOrchestratorScSyncMvcConfiguration(container *dig.Container) {
	// Controller
	err := container.Provide(
		orchestrator_sc_sync_mvc.NewUseCaseController,
		dig.As(new(orchestrator_sc_sync_mvc.ControllerInterface)),
		dig.Group("controllers"),
	)
	if err != nil {
		panic(err)
	}
}
