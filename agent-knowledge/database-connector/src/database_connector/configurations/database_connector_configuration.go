package configurations

import (
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/services"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases"
	"go.uber.org/dig"
)

func AddDatabaseConnectorConfiguration(container *dig.Container) {
	// Use cases
	err := container.Provide(
		use_cases.NewCreateDatabaseConnection,
		dig.As(new(use_cases.CreateDatabaseConnectionInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewRemoveDatabaseConnection,
		dig.As(new(use_cases.RemoveDatabaseConnectionInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewAddMessagesToAgentContext,
		dig.As(new(use_cases.AddMessagesToAgentContextInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewCreateAgentContext,
		dig.As(new(use_cases.CreateAgentContextInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetAgentContexts,
		dig.As(new(use_cases.GetAgentContextsInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetAgentContextByKey,
		dig.As(new(use_cases.GetAgentContextByKeyInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetAgentContextMessages,
		dig.As(new(use_cases.GetAgentContextMessagesInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewSaveAgentPromptTemplate,
		dig.As(new(use_cases.SaveAgentPromptTemplateInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetAgentPromptTemplate,
		dig.As(new(use_cases.GetAgentPromptTemplateInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewSetupAgentDatabase,
		dig.As(new(use_cases.SetupAgentDatabaseInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetLinkedAgentContextOfConversation,
		dig.As(new(use_cases.GetLinkedAgentContextOfConversationInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewLinkConversationToAgentContext,
		dig.As(new(use_cases.LinkConversationToAgentContextInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		use_cases.NewGetDatabaseConnectionSetupStatus,
		dig.As(new(use_cases.GetDatabaseConnectionSetupStatusInterface)),
	)
	if err != nil {
		panic(err)
	}

	// Services
	err = container.Provide(
		services.NewAgentDatabaseCredentialsStorage,
		dig.As(new(services.AgentDatabaseCredentialsStorageInterface)),
	)
	if err != nil {
		panic(err)
	}

	err = container.Provide(
		services.NewAgentDatabaseCredentialsGenerator,
		dig.As(new(services.AgentDatabaseCredentialsGeneratorInterface)),
	)
	if err != nil {
		panic(err)
	}
}
