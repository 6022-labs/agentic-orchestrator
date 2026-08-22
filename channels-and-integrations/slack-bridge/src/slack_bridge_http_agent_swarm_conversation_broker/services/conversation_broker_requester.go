package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_http_agent_swarm_conversation_broker/helpers"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

const jsonContentType = "application/json"

type ConversationBrokerRequester struct {
	client clients.ClientWithResponsesInterface
}

func NewConversationBrokerRequester(client clients.ClientWithResponsesInterface) *ConversationBrokerRequester {
	return &ConversationBrokerRequester{client: client}
}

func (r *ConversationBrokerRequester) RegisterAgent(
	ctx context.Context,
	collectionAddress string,
	collectionTokenId uint64,
) (uuid.UUID, error) {
	body, err := json.Marshal(clients.RequestsRegisterAgentRequest{
		CollectionAddress: &collectionAddress,
		CollectionTokenId: lo.ToPtr(fmt.Sprintf("%d", collectionTokenId)),
	})
	if err != nil {
		return uuid.Nil, err
	}

	rsp, err := r.client.PostParticipantsAgentsWithBodyWithResponse(
		ctx,
		jsonContentType,
		bytes.NewReader(body),
	)
	if err != nil {
		return uuid.Nil, err
	}

	// The broker upserts agents: 200 if already known, 201 if newly created.
	payload := rsp.JSON201
	if payload == nil {
		payload = rsp.JSON200
	}
	if payload == nil {
		return uuid.Nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}
	return helpers.ParseUuid(payload.Id), nil
}

func (r *ConversationBrokerRequester) RegisterUser(ctx context.Context, displayName string) (uuid.UUID, error) {
	body, err := json.Marshal(clients.RequestsRegisterUserRequest{DisplayName: &displayName})
	if err != nil {
		return uuid.Nil, err
	}

	rsp, err := r.client.PostParticipantsUsersWithBodyWithResponse(
		ctx,
		jsonContentType,
		bytes.NewReader(body),
	)
	if err != nil {
		return uuid.Nil, err
	}
	if rsp.JSON201 == nil {
		return uuid.Nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}
	return helpers.ParseUuid(rsp.JSON201.Id), nil
}

func (r *ConversationBrokerRequester) RegisterConversation(
	ctx context.Context,
	bridgeId uuid.UUID,
	description *string,
	participantIds []uuid.UUID,
) (uuid.UUID, error) {
	ids := lo.Map(participantIds, func(id uuid.UUID, _ int) string { return id.String() })
	body, err := json.Marshal(clients.RequestsRegisterConversationRequest{
		BridgeId:       lo.ToPtr(bridgeId.String()),
		Description:    description,
		ParticipantIds: &ids,
	})
	if err != nil {
		return uuid.Nil, err
	}

	rsp, err := r.client.PostConversationsWithBodyWithResponse(
		ctx,
		jsonContentType,
		bytes.NewReader(body),
	)
	if err != nil {
		return uuid.Nil, err
	}
	if rsp.JSON201 == nil {
		return uuid.Nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}
	return helpers.ParseUuid(rsp.JSON201.Id), nil
}

func (r *ConversationBrokerRequester) IngestMessage(
	ctx context.Context,
	conversationId uuid.UUID,
	participantId uuid.UUID,
	content string,
) (uuid.UUID, error) {
	body, err := json.Marshal(clients.RequestsIngestMessageRequest{
		ParticipantId: lo.ToPtr(participantId.String()),
		Content:       &content,
	})
	if err != nil {
		return uuid.Nil, err
	}

	rsp, err := r.client.PostConversationsIdMessagesWithBodyWithResponse(
		ctx,
		conversationId.String(),
		jsonContentType,
		bytes.NewReader(body),
	)
	if err != nil {
		return uuid.Nil, err
	}
	if rsp.JSON201 == nil {
		return uuid.Nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}
	return helpers.ParseUuid(rsp.JSON201.Id), nil
}

func (r *ConversationBrokerRequester) ForceConcludeConversation(ctx context.Context, conversationId uuid.UUID) error {
	rsp, err := r.client.PostConversationsIdForceConcludeWithResponse(ctx, conversationId.String())
	if err != nil {
		return err
	}
	if rsp.StatusCode() != http.StatusNoContent {
		return r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}
	return nil
}

func (r *ConversationBrokerRequester) AttachParticipant(
	ctx context.Context,
	conversationId uuid.UUID,
	participantId uuid.UUID,
) error {
	body, err := json.Marshal(
		clients.RequestsConversationParticipantRequest{ParticipantId: lo.ToPtr(participantId.String())},
	)
	if err != nil {
		return err
	}

	rsp, err := r.client.PostConversationsIdParticipantsWithBodyWithResponse(
		ctx,
		conversationId.String(),
		jsonContentType,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	if rsp.StatusCode() != http.StatusCreated {
		return r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}
	return nil
}

func (r *ConversationBrokerRequester) DetachParticipant(
	ctx context.Context,
	conversationId uuid.UUID,
	participantId uuid.UUID,
) error {
	rsp, err := r.client.DeleteConversationsIdParticipantsParticipantIdWithResponse(
		ctx,
		conversationId.String(),
		participantId.String(),
	)
	if err != nil {
		return err
	}
	if rsp.StatusCode() != http.StatusNoContent {
		return r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}
	return nil
}

func (r *ConversationBrokerRequester) unexpectedStatus(statusCode int, body []byte) error {
	return fmt.Errorf("unexpected status code: %d, response: %s", statusCode, string(body))
}
