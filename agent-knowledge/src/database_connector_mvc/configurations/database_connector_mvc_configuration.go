package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/src/database_connector_mvc"
	"go.uber.org/dig"
)

func AddOrchestratorAgentDatabaseConnectorMvcConfiguration(container *dig.Container) {
	// Controllers
	err := container.Provide(
		database_connector_mvc.NewDatabaseConnectionController,
		dig.Group("controllers"),
		dig.As(new(database_connector_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		database_connector_mvc.NewAgentPromptTemplateController,
		dig.Group("controllers"),
		dig.As(new(database_connector_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		database_connector_mvc.NewAgentContextController,
		dig.Group("controllers"),
		dig.As(new(database_connector_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		database_connector_mvc.NewAgentContextMessageController,
		dig.Group("controllers"),
		dig.As(new(database_connector_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		database_connector_mvc.NewAgentContextConversationMappingController,
		dig.Group("controllers"),
		dig.As(new(database_connector_mvc.ControllerInterface)),
	)
	if err != nil {
		panic(err)
	}
}
