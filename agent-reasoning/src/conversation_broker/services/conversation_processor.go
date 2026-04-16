package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/settings"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ConversationProcessor struct {
	logger                        *zap.Logger
	processingSettings            *settings.ProcessingSettings
	conversationConclude          interfaces.ConversationConcludeInterface
	bridgeCallbackRequester       interfaces.BridgeCallbackRequesterInterface
	agentCompletionsRequester     interfaces.AgentCompletionsRequesterInterface
	bridgeRepository              repositories.BridgeRepositoryInterface
	conversationRepository        repositories.ConversationRepositoryInterface
	conversationMessageRepository repositories.ConversationMessageRepositoryInterface
	participantRepository         repositories.ParticipantRepositoryInterface
}

func NewConversationProcessor(
	logger *zap.Logger,
	processingSettings *settings.ProcessingSettings,
	conversationConclude interfaces.ConversationConcludeInterface,
	bridgeCallbackRequester interfaces.BridgeCallbackRequesterInterface,
	agentCompletionsRequester interfaces.AgentCompletionsRequesterInterface,
	bridgeRepository repositories.BridgeRepositoryInterface,
	conversationRepository repositories.ConversationRepositoryInterface,
	conversationMessageRepository repositories.ConversationMessageRepositoryInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
) *ConversationProcessor {
	return &ConversationProcessor{
		logger:                        logger,
		processingSettings:            processingSettings,
		conversationConclude:          conversationConclude,
		bridgeCallbackRequester:       bridgeCallbackRequester,
		agentCompletionsRequester:     agentCompletionsRequester,
		bridgeRepository:              bridgeRepository,
		conversationRepository:        conversationRepository,
		conversationMessageRepository: conversationMessageRepository,
		participantRepository:         participantRepository,
	}
}

func (p *ConversationProcessor) ProcessConversation(participant *models.Participant, conversation *models.Conversation) error {
	logFields := p.processingLogFields(participant, conversation)

	if p.shouldForceConcludeConversation(conversation) {
		p.logger.Warn("Force concluding conversation", logFields...)
		if err := p.conversationConclude.Conclude(conversation.ID); err != nil {
			return fmt.Errorf("failed to force conclude conversation: %w", err)
		}
		return nil
	}

	bridge, err := p.bridgeRepository.GetById(conversation.BridgeID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		p.logger.Warn("Bridge not found for conversation", logFields...)
		if err := p.conversationConclude.Conclude(conversation.ID); err != nil {
			return fmt.Errorf("failed to conclude orphan conversation: %w", err)
		}
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to fetch bridge: %w", err)
	}

	messages, err := p.conversationMessageRepository.ListProcessableForParticipantAndConversation(participant, conversation.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch processable messages: %w", err)
	}
	if len(messages) == 0 {
		p.logger.Info("No processable messages for participant in conversation", logFields...)
		return nil
	}

	participants, err := p.participantRepository.ListByConversationId(conversation.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch conversation participants: %w", err)
	}

	forceCompletion := false
	if participant.Role != nil && *participant.Role == enums.ParticipantRoleFacilitator && p.shouldForceConcludeByFacilitator(conversation) {
		forceCompletion = true
	}

	completions, isFinal, err := p.agentCompletionsRequester.Completions(
		string(bridge.Kind),
		participant,
		conversation,
		messages,
		participants,
		forceCompletion,
	)
	if err != nil {
		return fmt.Errorf("failed to get completions: %w", err)
	}

	if participant.Role != nil && *participant.Role == enums.ParticipantRoleHuman {
		// Humans never produce completions.
		return nil
	}

	if len(completions) == 0 {
		return fmt.Errorf("no completions returned for conversation %s", conversation.ID)
	}

	if err := p.sendCompletionsToBridge(bridge, conversation, participant, completions); err != nil {
		return fmt.Errorf("failed to send completions to bridge: %w", err)
	}

	if isFinal {
		p.logger.Info("Conversation is concluded by compute", logFields...)
		if err := p.conversationConclude.Conclude(conversation.ID); err != nil {
			return fmt.Errorf("failed to conclude conversation: %w", err)
		}
	}

	return nil
}

func (p *ConversationProcessor) sendCompletionsToBridge(
	bridge *models.Bridge,
	conversation *models.Conversation,
	participant *models.Participant,
	completions []*string,
) error {
	// Reload to ensure not concluded by another tick in the meantime.
	reloaded, err := p.conversationRepository.GetById(conversation.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("conversation %s vanished before sending completions", conversation.ID)
	}
	if err != nil {
		return fmt.Errorf("failed to reload conversation: %w", err)
	}
	if reloaded.Status == enums.ConversationStatusConcluded {
		p.logger.Warn("Conversation already concluded, skipping completions",
			p.processingLogFields(participant, conversation)...,
		)
		return nil
	}

	for _, completion := range completions {
		if err := p.bridgeCallbackRequester.PostCompletion(bridge, conversation, participant, *completion); err != nil {
			return fmt.Errorf("failed to post completion to bridge: %w", err)
		}
	}
	return nil
}

func (p *ConversationProcessor) shouldForceConcludeByFacilitator(conversation *models.Conversation) bool {
	return p.shouldForceConcludeAfter(conversation, p.processingSettings.ConversationForceConcludeByFacilitatorAfter)
}

func (p *ConversationProcessor) shouldForceConcludeConversation(conversation *models.Conversation) bool {
	return p.shouldForceConcludeAfter(conversation, p.processingSettings.ConversationForceConcludeAfter)
}

func (p *ConversationProcessor) shouldForceConcludeAfter(conversation *models.Conversation, delaySeconds uint64) bool {
	if delaySeconds == 0 {
		return false
	}
	concludeAfter := time.Duration(delaySeconds) * time.Second

	referenceTime := conversation.CreatedAt
	if conversation.UnconcludedAt != nil {
		referenceTime = *conversation.UnconcludedAt
	}
	if referenceTime.IsZero() {
		return false
	}
	return referenceTime.Add(concludeAfter).Before(time.Now())
}

func (p *ConversationProcessor) processingLogFields(participant *models.Participant, conversation *models.Conversation) []zap.Field {
	fields := []zap.Field{
		zap.String("conversationId", conversation.ID.String()),
		zap.String("bridgeId", conversation.BridgeID.String()),
		zap.String("participantId", participant.ID.String()),
	}
	if participant.Role != nil {
		fields = append(fields, zap.String("participantRole", string(*participant.Role)))
	}
	if participant.EnsDomain != nil {
		fields = append(fields, zap.String("participantEnsDomain", *participant.EnsDomain))
	}
	if participant.CollectionAddress != nil {
		fields = append(fields, zap.String("participantCollectionAddress", *participant.CollectionAddress))
	}
	if participant.CollectionTokenId != nil {
		fields = append(fields, zap.Uint64("participantCollectionTokenId", *participant.CollectionTokenId))
	}
	return fields
}
