package services

import (
	"context"

	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	cb_errors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/settings"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// facilitatorConcludeMarker in a facilitator completion (or the message that
// triggered it) signals the conversation is ready to conclude.
const facilitatorConcludeMarker = ":end:"

type ConversationProcessor struct {
	logger                          *zap.Logger
	processingSettings              *settings.ProcessingSettings
	conversationConclude            interfaces.ConversationConcludeInterface
	bridgeCallbackRequester         interfaces.BridgeCallbackRequesterInterface
	remoteAgentCardResolver         interfaces.RemoteAgentCardResolverInterface
	remoteAgentCompletionsRequester interfaces.RemoteAgentCompletionsRequesterInterface
	conversationPromptBuilder       interfaces.ConversationPromptBuilderInterface
	participantRoleResolver         interfaces.ParticipantRoleResolverInterface
	bridgeRepository                repositories.BridgeRepositoryInterface
	conversationRepository          repositories.ConversationRepositoryInterface
	conversationMessageRepository   repositories.ConversationMessageRepositoryInterface
	participantRepository           repositories.ParticipantRepositoryInterface
	emojiDecoder                    interfaces.EmojiDecoderInterface

	// One lock per conversation. In-process only — swap for a Postgres
	// advisory lock if the broker ever runs as multiple replicas.
	conversationLocks sync.Map
}

func NewConversationProcessor(
	logger *zap.Logger,
	processingSettings *settings.ProcessingSettings,
	conversationConclude interfaces.ConversationConcludeInterface,
	bridgeCallbackRequester interfaces.BridgeCallbackRequesterInterface,
	remoteAgentCardResolver interfaces.RemoteAgentCardResolverInterface,
	remoteAgentCompletionsRequester interfaces.RemoteAgentCompletionsRequesterInterface,
	conversationPromptBuilder interfaces.ConversationPromptBuilderInterface,
	participantRoleResolver interfaces.ParticipantRoleResolverInterface,
	bridgeRepository repositories.BridgeRepositoryInterface,
	conversationRepository repositories.ConversationRepositoryInterface,
	conversationMessageRepository repositories.ConversationMessageRepositoryInterface,
	participantRepository repositories.ParticipantRepositoryInterface,
	emojiDecoder interfaces.EmojiDecoderInterface,
) *ConversationProcessor {
	return &ConversationProcessor{
		logger:                          logger,
		processingSettings:              processingSettings,
		conversationConclude:            conversationConclude,
		bridgeCallbackRequester:         bridgeCallbackRequester,
		remoteAgentCardResolver:         remoteAgentCardResolver,
		remoteAgentCompletionsRequester: remoteAgentCompletionsRequester,
		conversationPromptBuilder:       conversationPromptBuilder,
		participantRoleResolver:         participantRoleResolver,
		bridgeRepository:                bridgeRepository,
		conversationRepository:          conversationRepository,
		conversationMessageRepository:   conversationMessageRepository,
		participantRepository:           participantRepository,
		emojiDecoder:                    emojiDecoder,
	}
}

func (p *ConversationProcessor) ProcessConversation(
	ctx context.Context,
	participant *models.Participant,
	conversation *models.Conversation,
) error {
	logFields := p.processingLogFields(participant, conversation)

	if p.shouldForceConcludeConversation(conversation) {
		p.logger.Warn("Force concluding conversation", logFields...)
		if err := p.conversationConclude.Conclude(ctx, conversation.ID); err != nil {
			return fmt.Errorf("failed to force conclude conversation: %w", err)
		}
		return nil
	}

	role, err := p.participantRoleResolver.Resolve(ctx, participant)
	if err != nil {
		return fmt.Errorf("failed to resolve participant role: %w", err)
	}

	if role != nil && *role == enums.ParticipantRoleHuman {
		// Humans have no runtime to invoke.
		return nil
	}

	bridge, err := p.bridgeRepository.GetById(ctx, conversation.BridgeID)
	if err != nil {
		return fmt.Errorf("failed to fetch bridge: %w", err)
	}
	if bridge == nil {
		p.logger.Warn("Bridge not found for conversation", logFields...)
		if err := p.conversationConclude.Conclude(ctx, conversation.ID); err != nil {
			return fmt.Errorf("failed to conclude orphan conversation: %w", err)
		}
		return nil
	}

	messages, err := p.conversationMessageRepository.ListProcessableForParticipantAndConversation(ctx,
		participant,
		role,
		conversation.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch processable messages: %w", err)
	}
	if len(messages) == 0 {
		p.logger.Info("No processable messages for participant in conversation", logFields...)
		return nil
	}

	participants, err := p.participantRepository.ListByConversationId(ctx, conversation.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch conversation participants: %w", err)
	}

	if participant.EnsDomain == nil {
		return fmt.Errorf("agent participant %s has no ENS domain", participant.ID)
	}

	card, err := p.remoteAgentCardResolver.ResolveCard(ctx, *participant.EnsDomain)
	if err != nil {
		return fmt.Errorf("failed to resolve remote agent card: %w", err)
	}

	turn, err := p.conversationPromptBuilder.BuildInput(ctx,
		string(bridge.Kind),
		participant,
		role,
		conversation,
		messages,
		participants,
	)
	if err != nil {
		return fmt.Errorf("failed to build turn: %w", err)
	}

	completions, err := p.remoteAgentCompletionsRequester.RequestCompletions(ctx, card, turn)
	var paymentRequired *cb_errors.PaymentRequiredError
	if errors.As(err, &paymentRequired) {
		p.logger.Info("Agent node requires payment, skipping turn", logFields...)
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to request turn: %w", err)
	}
	if len(completions) == 0 {
		return fmt.Errorf("no completions returned for conversation %s", conversation.ID)
	}

	completions = p.demojizeAll(completions)

	isFinal := false
	if role != nil && *role == enums.ParticipantRoleFacilitator {
		completions, isFinal, err = p.concludeByFacilitatorIfNeeded(ctx,
			participant, conversation, card, turn, completions, role, messages, participants,
		)
		if err != nil {
			return err
		}
	}

	if err := p.sendCompletionsToBridge(ctx, bridge, conversation, participant, completions, isFinal); err != nil {
		return fmt.Errorf("failed to send completions to bridge: %w", err)
	}

	return nil
}

// concludeByFacilitatorIfNeeded reproduces the facilitator wrap-up: when the
// broker timer elapsed or the conclude marker appears in a fresh completion or
// the latest inbound message, the facilitator gets one extra conclude turn.
func (p *ConversationProcessor) concludeByFacilitatorIfNeeded(
	ctx context.Context,
	participant *models.Participant,
	conversation *models.Conversation,
	card *models.RemoteAgentCard,
	turn *models.RemoteAgentInput,
	completions []string,
	role *enums.ParticipantRole,
	messages []*models.ConversationMessage,
	participants []*models.Participant,
) ([]string, bool, error) {
	readyToConclude := p.shouldForceConcludeByFacilitator(conversation) ||
		p.carriesConcludeMarker(completions, messages)
	if !readyToConclude {
		return completions, false, nil
	}

	concludeTurn, err := p.conversationPromptBuilder.BuildConcludeInput(ctx,
		turn,
		completions,
		participant,
		role,
		conversation,
		messages,
		participants,
	)
	if err != nil {
		return nil, false, fmt.Errorf("failed to build conclude turn: %w", err)
	}

	concludeCompletions, err := p.remoteAgentCompletionsRequester.RequestCompletions(ctx, card, concludeTurn)
	var paymentRequired *cb_errors.PaymentRequiredError
	if errors.As(err, &paymentRequired) {
		p.logger.Info("Agent node requires payment for conclude turn, concluding without wrap-up",
			p.processingLogFields(participant, conversation)...,
		)
		return completions, true, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("failed to request conclude turn: %w", err)
	}

	return append(completions, p.demojizeAll(concludeCompletions)...), true, nil
}

func (p *ConversationProcessor) sendCompletionsToBridge(
	ctx context.Context,
	bridge *models.Bridge,
	conversation *models.Conversation,
	participant *models.Participant,
	completions []string,
	isFinal bool,
) error {
	mu := p.lockForConversation(conversation.ID)
	mu.Lock()
	defer mu.Unlock()

	// Reload under lock to catch a concurrent conclude.
	reloaded, err := p.conversationRepository.GetById(ctx, conversation.ID)
	if err != nil {
		return fmt.Errorf("failed to reload conversation: %w", err)
	}
	if reloaded == nil {
		return fmt.Errorf("conversation %s vanished before sending completions", conversation.ID)
	}
	if reloaded.Status == enums.ConversationStatusConcluded {
		p.logger.Warn("Conversation already concluded, skipping completions",
			p.processingLogFields(participant, conversation)...,
		)
		return nil
	}

	for _, completion := range completions {
		if err := p.bridgeCallbackRequester.PostCompletion(
			ctx,
			bridge,
			conversation,
			participant,
			completion,
		); err != nil {
			return fmt.Errorf("failed to post completion to bridge: %w", err)
		}
	}

	if isFinal {
		p.logger.Info("Conversation is concluded by facilitator",
			p.processingLogFields(participant, conversation)...,
		)
		if err := p.conversationConclude.Conclude(ctx, conversation.ID); err != nil {
			return fmt.Errorf("failed to conclude conversation: %w", err)
		}
	}

	return nil
}

// carriesConcludeMarker mirrors the old facilitator strategy: the marker
// counts whether it comes from the facilitator's fresh completions or the
// latest inbound platform message (a user saying ":end:" ends the swarm).
func (p *ConversationProcessor) carriesConcludeMarker(
	completions []string,
	messages []*models.ConversationMessage,
) bool {
	for _, completion := range completions {
		if strings.Contains(strings.ToLower(completion), facilitatorConcludeMarker) {
			return true
		}
	}
	if len(messages) > 0 {
		return strings.Contains(strings.ToLower(messages[len(messages)-1].Content), facilitatorConcludeMarker)
	}
	return false
}

func (p *ConversationProcessor) demojizeAll(completions []string) []string {
	out := make([]string, len(completions))
	for i, completion := range completions {
		out[i] = p.emojiDecoder.DemojizeContent(completion)
	}
	return out
}

func (p *ConversationProcessor) lockForConversation(id uuid.UUID) *sync.Mutex {
	raw, _ := p.conversationLocks.LoadOrStore(id, &sync.Mutex{})
	return raw.(*sync.Mutex)
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

func (p *ConversationProcessor) processingLogFields(
	participant *models.Participant,
	conversation *models.Conversation,
) []zap.Field {
	fields := []zap.Field{
		zap.String("conversationId", conversation.ID.String()),
		zap.String("bridgeId", conversation.BridgeID.String()),
		zap.String("participantId", participant.ID.String()),
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
