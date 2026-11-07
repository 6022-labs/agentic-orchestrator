package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/clients"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_http_agent_swarm_conversation_broker/mappings"
	"github.com/google/uuid"
)

const jsonContentType = "application/json"

type ConversationBrokerRequester struct {
	client clients.ClientWithResponsesInterface
}

func NewConversationBrokerRequester(client clients.ClientWithResponsesInterface) *ConversationBrokerRequester {
	return &ConversationBrokerRequester{client: client}
}

func (r *ConversationBrokerRequester) unexpectedStatus(statusCode int, body []byte) error {
	return fmt.Errorf("unexpected status code: %d, response: %s", statusCode, string(body))
}

func (r *ConversationBrokerRequester) ListConversations(
	ctx context.Context,
	page, elementsPerPage *uint,
) (*responses.ListBrokerConversationsResponse, error) {
	params := &clients.GetConversationsParams{}
	if page != nil {
		value := int(*page)
		params.Page = &value
	}
	if elementsPerPage != nil {
		value := int(*elementsPerPage)
		params.ElementsPerPage = &value
	}

	rsp, err := r.client.GetConversationsWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if rsp.JSON200 == nil {
		return nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}

	return (*mappings.ExternalListConversationsResponse)(rsp.JSON200).ToUseCaseResponse(), nil
}

func (r *ConversationBrokerRequester) GetConversation(
	ctx context.Context,
	conversationId uuid.UUID,
) (*responses.BrokerConversationResponse, error) {
	rsp, err := r.client.GetConversationsIdWithResponse(ctx, conversationId.String())
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode() == http.StatusNotFound {
		return nil, errors.NewNotFoundError("conversation_not_found", "conversation not found")
	}
	if rsp.JSON200 == nil {
		return nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}

	return (*mappings.ExternalConversationResponse)(rsp.JSON200).ToUseCaseResponse(), nil
}

func (r *ConversationBrokerRequester) GetConversationMessages(
	ctx context.Context,
	conversationId uuid.UUID,
) ([]*responses.BrokerMessageResponse, error) {
	rsp, err := r.client.GetConversationsIdMessagesWithResponse(ctx, conversationId.String())
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode() == http.StatusNotFound {
		return nil, errors.NewNotFoundError("conversation_not_found", "conversation not found")
	}
	if rsp.JSON200 == nil {
		return nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}

	messages := make([]*responses.BrokerMessageResponse, len(*rsp.JSON200))
	for i := range *rsp.JSON200 {
		external := mappings.ExternalMessageResponse((*rsp.JSON200)[i])
		messages[i] = external.ToUseCaseResponse()
	}
	return messages, nil
}

func (r *ConversationBrokerRequester) ListConversationParticipants(
	ctx context.Context,
	conversationId uuid.UUID,
) ([]*responses.BrokerParticipantResponse, error) {
	rsp, err := r.client.GetConversationsIdParticipantsWithResponse(ctx, conversationId.String())
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode() == http.StatusNotFound {
		return nil, errors.NewNotFoundError("conversation_not_found", "conversation not found")
	}
	if rsp.JSON200 == nil {
		return nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}

	participants := make([]*responses.BrokerParticipantResponse, len(*rsp.JSON200))
	for i := range *rsp.JSON200 {
		external := mappings.ExternalParticipantResponse((*rsp.JSON200)[i])
		participants[i] = external.ToUseCaseResponse()
	}
	return participants, nil
}

func (r *ConversationBrokerRequester) ForceConcludeConversation(ctx context.Context, conversationId uuid.UUID) error {
	rsp, err := r.client.PostConversationsIdForceConcludeWithResponse(ctx, conversationId.String())
	if err != nil {
		return err
	}
	if rsp.StatusCode() == http.StatusNotFound {
		return errors.NewNotFoundError("conversation_not_found", "conversation not found")
	}
	if rsp.StatusCode() != http.StatusNoContent {
		return r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}
	return nil
}

func (r *ConversationBrokerRequester) GetAgentPromptTemplate(
	ctx context.Context,
	collectionAddress string,
	collectionTokenId uint64,
	key string,
) (*responses.AgentPromptTemplateResponse, error) {
	rsp, err := r.client.GetAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyWithResponse(
		ctx,
		collectionAddress,
		strconv.FormatUint(collectionTokenId, 10),
		key,
	)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode() == http.StatusNotFound {
		return nil, errors.NewNotFoundError("prompt_template_not_found", "prompt template not found")
	}
	if rsp.JSON200 == nil {
		return nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}

	return (*mappings.ExternalPromptTemplateResponse)(rsp.JSON200).ToUseCaseResponse(), nil
}

func (r *ConversationBrokerRequester) SaveAgentPromptTemplate(
	ctx context.Context,
	collectionAddress string,
	collectionTokenId uint64,
	key, content string,
) (*responses.AgentPromptTemplateResponse, error) {
	body, err := json.Marshal(map[string]string{"content": content})
	if err != nil {
		return nil, err
	}

	rsp, err := r.client.PostAgentsCollectionAddressCollectionTokenIdPromptTemplatesKeyWithBodyWithResponse(
		ctx,
		collectionAddress,
		strconv.FormatUint(collectionTokenId, 10),
		key,
		jsonContentType,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode() == http.StatusBadRequest {
		// A 400 is either a ValidationError {field,message} or a BadRequestError {code,message}.
		var body struct {
			Field   string `json:"field"`
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if json.Unmarshal(rsp.Body, &body) != nil || body.Message == "" {
			return nil, errors.NewValidationError("request", string(rsp.Body))
		}
		if body.Field == "" || body.Code != "" {
			return nil, errors.NewValidationError("request", body.Message)
		}
		return nil, errors.NewValidationError(body.Field, body.Message)
	}
	if rsp.StatusCode() == http.StatusNotFound {
		return nil, errors.NewNotFoundError("agent_not_found", "agent not found")
	}
	if rsp.JSON200 == nil {
		return nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}

	return (*mappings.ExternalPromptTemplateResponse)(rsp.JSON200).ToUseCaseResponse(), nil
}

func (r *ConversationBrokerRequester) GetAgentCard(
	ctx context.Context,
	ensDomain string,
) (*responses.AgentCardResponse, error) {
	rsp, err := r.client.GetRemoteAgentsEnsDomainCardWithResponse(ctx, ensDomain)
	if err != nil {
		return nil, err
	}
	if rsp.JSON200 == nil {
		return nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}

	return (*mappings.ExternalAgentCardResponse)(rsp.JSON200).ToUseCaseResponse(), nil
}

func (r *ConversationBrokerRequester) GetAgentUrl(
	ctx context.Context,
	ensDomain string,
) (*responses.AgentUrlResponse, error) {
	rsp, err := r.client.GetRemoteAgentsEnsDomainUrlWithResponse(ctx, ensDomain)
	if err != nil {
		return nil, err
	}
	if rsp.JSON200 == nil {
		return nil, r.unexpectedStatus(rsp.StatusCode(), rsp.Body)
	}

	return (*mappings.ExternalAgentUrlResponse)(rsp.JSON200).ToUseCaseResponse(), nil
}
