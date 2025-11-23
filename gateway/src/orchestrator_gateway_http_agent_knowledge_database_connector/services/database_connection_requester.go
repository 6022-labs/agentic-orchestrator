package services

import (
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway/services/interfaces"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_agent_knowledge_database_connector/clients"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator_gateway_http_agent_knowledge_database_connector/models"
	"github.com/google/uuid"
)

type DatabaseConnectionRequester struct {
	client clients.AgentDatabaseConnectorClientInterface
}

func NewDatabaseConnectionRequester(
	client clients.AgentDatabaseConnectorClientInterface,
) *DatabaseConnectionRequester {
	return &DatabaseConnectionRequester{
		client: client,
	}
}

func (requester *DatabaseConnectionRequester) Create(request *interfaces.CreateDatabaseConnectionRequest) (*interfaces.CreateDatabaseConnectionResult, error) {
	externalResponse, err := requester.client.CreateDatabaseConnection(&models.ExternalCreateDatabaseConnectionRequest{
		Host:    request.Host,
		Port:    request.Port,
		Driver:  request.Driver,
		AgentId: request.AgentId,
	})
	if err != nil {
		return nil, err
	}

	return externalResponse.ToCreateDatabaseConnectionResult(), nil
}

func (requester *DatabaseConnectionRequester) Remove(agentId uuid.UUID) error {
	err := requester.client.RemoveDatabaseConnection(&models.ExternalAgentIdRequest{
		AgentId: agentId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (requester *DatabaseConnectionRequester) SetupAgentDatabase(agentId uuid.UUID, copyFromAgentId *uuid.UUID) error {
	err := requester.client.SetupAgentDatabase(&models.ExternalSetupAgentDatabaseRequest{
		AgentId:         agentId,
		CopyFromAgentId: copyFromAgentId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (requester *DatabaseConnectionRequester) GetSetupStatus(agentId uuid.UUID) (*interfaces.DatabaseConnectionStatus, error) {
	response, err := requester.client.GetDatabaseConnectionSetupStatus(&models.ExternalAgentIdRequest{
		AgentId: agentId,
	})
	if err != nil {
		return nil, err
	}

	return &interfaces.DatabaseConnectionStatus{
		IsCompleted: response.IsCompleted,
		Status:      response.SetupStatus,
	}, nil
}
