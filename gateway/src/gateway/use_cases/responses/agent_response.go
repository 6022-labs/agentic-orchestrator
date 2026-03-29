package responses

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/models/enums"
	"github.com/google/uuid"
)

type AgentResponse struct {
	ID uuid.UUID `json:"id"`

	Name          string            `json:"name"`
	Role          enums.AgentRole   `json:"role"`
	Owner         *string           `json:"owner"`
	Creator       string            `json:"creator"`
	Images        map[string]string `json:"images"`
	CloneOf       *uint64           `json:"cloneOf"`
	Attributes    map[string]string `json:"attributes"`
	EnsDomain     string            `json:"ensDomain"`
	IsInternal    bool              `json:"isInternal"`
	WalletAddress string            `json:"walletAddress"`

	LLMSetupStatus      *string `json:"llmSetupStatus"`
	DatabaseSetupStatus *string `json:"databaseSetupStatus"`

	IsSetupCompleted *bool `json:"isSetupCompleted"`

	ProposalId        *uint64 `json:"proposalId"`
	CollectionTokenId *uint64 `json:"collectionTokenId"`
	CollectionAddress string  `json:"collectionAddress"`
}

func AgentResponseFromModel(
	agent *models.Agent,
	llmConnectionSetupStatus, databaseConnectionSetupStatus *string,
	llmConnectionSetupIsCompleted, databaseConnectionSetupIsCompleted bool,
) *AgentResponse {
	if agent == nil {
		return nil
	}

	var isSetupCompletedPtr *bool = nil
	if agent.IsInternal {
		isSetupCompleted := false
		if agent.Role != enums.AgentRoleHuman {
			isSetupCompleted = llmConnectionSetupIsCompleted && databaseConnectionSetupIsCompleted
		} else {
			isSetupCompleted = databaseConnectionSetupIsCompleted
		}

		isSetupCompletedPtr = &isSetupCompleted
	}

	return &AgentResponse{
		ID: agent.ID,

		Name:          agent.Name,
		Role:          agent.Role,
		Owner:         agent.Owner,
		Images:        agent.Images,
		Creator:       agent.Creator,
		CloneOf:       agent.CloneOf,
		EnsDomain:     agent.EnsDomain,
		IsInternal:    agent.IsInternal,
		Attributes:    agent.Attributes,
		WalletAddress: agent.WalletAddress,

		LLMSetupStatus:      llmConnectionSetupStatus,
		DatabaseSetupStatus: databaseConnectionSetupStatus,

		IsSetupCompleted: isSetupCompletedPtr,

		ProposalId:        agent.ProposalId,
		CollectionTokenId: agent.CollectionTokenId,
		CollectionAddress: agent.CollectionAddress,
	}
}
