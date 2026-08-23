package use_cases

import (
	"context"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
)

type ForceConcludeInterface interface {
	Conclude(ctx context.Context, request *requests.ConversationIdRequest) error
}

type ForceConclude struct {
	brokerMetrics        metrics_interfaces.BrokerMetricsInterface
	conversationConclude interfaces.ConversationConcludeInterface
}

func NewForceConclude(
	brokerMetrics metrics_interfaces.BrokerMetricsInterface,
	conversationConclude interfaces.ConversationConcludeInterface,
) *ForceConclude {
	return &ForceConclude{brokerMetrics: brokerMetrics, conversationConclude: conversationConclude}
}

func (uc *ForceConclude) Conclude(ctx context.Context, request *requests.ConversationIdRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	concluded, err := uc.conversationConclude.ConcludeIfExists(ctx, request.ConversationID)
	if err != nil {
		return err
	}

	if concluded {
		uc.brokerMetrics.RecordConversationConcluded(ctx, metrics_interfaces.ConclusionReasonForced)
	}
	return nil
}
