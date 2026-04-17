package services

import (
	"fmt"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_conversation_broker/models"
	"github.com/google/uuid"
)

type ConversationBrokerRequester struct {
	client clients.ConversationBrokerClientInterface
}

func NewConversationBrokerRequester(client clients.ConversationBrokerClientInterface) *ConversationBrokerRequester {
	return &ConversationBrokerRequester{client: client}
}

func (r *ConversationBrokerRequester) RegisterAgent(collectionAddress string, collectionTokenId uint64) (uuid.UUID, error) {
	resp, err := r.client.RegisterAgent(&models.ExternalRegisterAgentRequest{
		CollectionAddress: collectionAddress,
		CollectionTokenId: fmt.Sprintf("%d", collectionTokenId),
	})
	if err != nil {
		return uuid.Nil, err
	}
	return resp.ID, nil
}

func (r *ConversationBrokerRequester) RegisterUser(displayName string) (uuid.UUID, error) {
	resp, err := r.client.RegisterUser(&models.ExternalRegisterUserRequest{DisplayName: displayName})
	if err != nil {
		return uuid.Nil, err
	}
	return resp.ID, nil
}

func (r *ConversationBrokerRequester) RegisterConversation(bridgeId uuid.UUID, description *string, participantIds []uuid.UUID) (uuid.UUID, error) {
	resp, err := r.client.RegisterConversation(&models.ExternalRegisterConversationRequest{
		BridgeID:       bridgeId,
		Description:    description,
		ParticipantIDs: participantIds,
	})
	if err != nil {
		return uuid.Nil, err
	}
	return resp.ID, nil
}

func (r *ConversationBrokerRequester) IngestMessage(conversationId uuid.UUID, participantId uuid.UUID, content string) (uuid.UUID, error) {
	resp, err := r.client.IngestMessage(conversationId, &models.ExternalIngestMessageRequest{
		ParticipantID: participantId,
		Content:       content,
	})
	if err != nil {
		return uuid.Nil, err
	}
	return resp.ID, nil
}

func (r *ConversationBrokerRequester) ForceConcludeConversation(conversationId uuid.UUID) error {
	return r.client.ForceConcludeConversation(conversationId)
}

func (r *ConversationBrokerRequester) AttachParticipant(conversationId uuid.UUID, participantId uuid.UUID) error {
	return r.client.AttachParticipant(conversationId, &models.ExternalConversationParticipantRequest{ParticipantID: participantId})
}

func (r *ConversationBrokerRequester) DetachParticipant(conversationId uuid.UUID, participantId uuid.UUID) error {
	return r.client.DetachParticipant(conversationId, participantId)
}
