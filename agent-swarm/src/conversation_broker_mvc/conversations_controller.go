package conversation_broker_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	REGISTER_CONVERSATION_ENDPOINT           = "/conversations"
	LIST_CONVERSATIONS_ENDPOINT              = "/conversations"
	GET_CONVERSATION_ENDPOINT                = "/conversations/:id"
	UPDATE_CONVERSATION_DESCRIPTION_ENDPOINT = "/conversations/:id/description"
	FORCE_CONCLUDE_ENDPOINT                  = "/conversations/:id/force-conclude"
	ATTACH_PARTICIPANT_ENDPOINT              = "/conversations/:id/participants"
	LIST_CONVERSATION_PARTICIPANTS_ENDPOINT  = "/conversations/:id/participants"
	DETACH_PARTICIPANT_ENDPOINT              = "/conversations/:id/participants/:participantId"
	INGEST_MESSAGE_ENDPOINT                  = "/conversations/:id/messages"
	GET_CONVERSATION_MESSAGES_ENDPOINT       = "/conversations/:id/messages"
)

type ConversationsController struct {
	logger                        *zap.Logger
	registerConversation          use_cases.RegisterConversationInterface
	getConversation               use_cases.GetConversationInterface
	listConversations             use_cases.ListConversationsInterface
	getConversationMessages       use_cases.GetConversationMessagesInterface
	listConversationParticipants  use_cases.ListConversationParticipantsInterface
	updateConversationDescription use_cases.UpdateConversationDescriptionInterface
	forceConclude                 use_cases.ForceConcludeInterface
	attachParticipant             use_cases.AttachParticipantInterface
	detachParticipant             use_cases.DetachParticipantInterface
	ingestMessage                 use_cases.IngestMessageInterface
}

func NewConversationsController(
	logger *zap.Logger,
	registerConversation use_cases.RegisterConversationInterface,
	getConversation use_cases.GetConversationInterface,
	listConversations use_cases.ListConversationsInterface,
	getConversationMessages use_cases.GetConversationMessagesInterface,
	listConversationParticipants use_cases.ListConversationParticipantsInterface,
	updateConversationDescription use_cases.UpdateConversationDescriptionInterface,
	forceConclude use_cases.ForceConcludeInterface,
	attachParticipant use_cases.AttachParticipantInterface,
	detachParticipant use_cases.DetachParticipantInterface,
	ingestMessage use_cases.IngestMessageInterface,
) *ConversationsController {
	return &ConversationsController{
		logger:                        logger,
		registerConversation:          registerConversation,
		getConversation:               getConversation,
		listConversations:             listConversations,
		getConversationMessages:       getConversationMessages,
		listConversationParticipants:  listConversationParticipants,
		updateConversationDescription: updateConversationDescription,
		forceConclude:                 forceConclude,
		attachParticipant:             attachParticipant,
		detachParticipant:             detachParticipant,
		ingestMessage:                 ingestMessage,
	}
}

func (c *ConversationsController) InitRoutes(app fiber.Router) {
	app.Post(REGISTER_CONVERSATION_ENDPOINT, c.Register)
	app.Get(LIST_CONVERSATIONS_ENDPOINT, c.List)
	app.Get(GET_CONVERSATION_ENDPOINT, c.Get)
	app.Patch(UPDATE_CONVERSATION_DESCRIPTION_ENDPOINT, c.UpdateDescription)
	app.Post(FORCE_CONCLUDE_ENDPOINT, c.ForceConclude)
	app.Post(ATTACH_PARTICIPANT_ENDPOINT, c.Attach)
	app.Get(LIST_CONVERSATION_PARTICIPANTS_ENDPOINT, c.ListParticipants)
	app.Delete(DETACH_PARTICIPANT_ENDPOINT, c.Detach)
	app.Post(INGEST_MESSAGE_ENDPOINT, c.Ingest)
	app.Get(GET_CONVERSATION_MESSAGES_ENDPOINT, c.GetMessages)
}

// @Summary List Conversations
// @Description Returns conversations sorted by most recently updated, paginated.
// @Tags Conversations
// @Produce json
// @Param page query int false "Page (1-based)"
// @Param elementsPerPage query int false "Elements per page (max 100)"
// @Success 200 {object} responses.ListConversationsResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /conversations [get]
func (c *ConversationsController) List(ctx *fiber.Ctx) error {
	req := new(requests.ListConversationsRequest)
	if err := ctx.QueryParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	resp, err := c.listConversations.List(req)
	if err != nil {
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.JSON(resp)
}

// @Summary Get Conversation Messages
// @Description Returns every message of the conversation in chronological order.
// @Tags Conversations
// @Produce json
// @Param id path string true "Conversation ID"
// @Success 200 {array} responses.ConversationMessageResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /conversations/{id}/messages [get]
func (c *ConversationsController) GetMessages(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}
	messages, err := c.getConversationMessages.Get(&requests.ConversationIdRequest{ConversationID: id})
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.JSON(messages)
}

// @Summary List Conversation Participants
// @Description Returns every participant attached to the conversation.
// @Tags Conversations
// @Produce json
// @Param id path string true "Conversation ID"
// @Success 200 {array} responses.ParticipantResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /conversations/{id}/participants [get]
func (c *ConversationsController) ListParticipants(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}
	participants, err := c.listConversationParticipants.List(&requests.ConversationIdRequest{ConversationID: id})
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.JSON(participants)
}

// @Summary Update Conversation Description
// @Description Overwrites the description of an existing conversation. Use null to clear it.
// @Tags Conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Param request body requests.UpdateConversationDescriptionRequest true "Update Conversation Description Request"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} errors.ValidationError
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /conversations/{id}/description [patch]
func (c *ConversationsController) UpdateDescription(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}
	req := new(requests.UpdateConversationDescriptionRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	req.ConversationID = id
	if err := c.updateConversationDescription.Update(req); err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

// @Summary Register Conversation
// @Description Creates a conversation attached to a bridge. Optionally attaches participants in the same call.
// @Tags Conversations
// @Accept json
// @Produce json
// @Param request body requests.RegisterConversationRequest true "Register Conversation Request"
// @Success 201 {object} responses.ConversationResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /conversations [post]
func (c *ConversationsController) Register(ctx *fiber.Ctx) error {
	req := new(requests.RegisterConversationRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	conv, err := c.registerConversation.Register(req)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.Status(fiber.StatusCreated).JSON(conv)
}

// @Summary Get Conversation
// @Description Returns a conversation by id.
// @Tags Conversations
// @Produce json
// @Param id path string true "Conversation ID"
// @Success 200 {object} responses.ConversationResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /conversations/{id} [get]
func (c *ConversationsController) Get(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}
	conv, err := c.getConversation.Get(&requests.ConversationIdRequest{ConversationID: id})
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.JSON(conv)
}

// @Summary Force Conclude Conversation
// @Description Marks the conversation as concluded immediately, bypassing the swarm timers.
// @Tags Conversations
// @Param id path string true "Conversation ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /conversations/{id}/force-conclude [post]
func (c *ConversationsController) ForceConclude(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}
	if err := c.forceConclude.Conclude(&requests.ConversationIdRequest{ConversationID: id}); err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

// @Summary Attach Participant to Conversation
// @Description Adds a participant (agent or user) to the conversation.
// @Tags Conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Param request body requests.ConversationParticipantRequest true "Attach Participant Request"
// @Success 201 {string} string "Created"
// @Failure 400 {object} errors.ValidationError
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /conversations/{id}/participants [post]
func (c *ConversationsController) Attach(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}
	req := new(requests.ConversationParticipantRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	req.ConversationID = id
	if err := c.attachParticipant.Attach(req); err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.SendStatus(fiber.StatusCreated)
}

// @Summary Detach Participant from Conversation
// @Description Removes a participant from the conversation.
// @Tags Conversations
// @Param id path string true "Conversation ID"
// @Param participantId path string true "Participant ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /conversations/{id}/participants/{participantId} [delete]
func (c *ConversationsController) Detach(ctx *fiber.Ctx) error {
	convID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}
	partID, err := uuid.Parse(ctx.Params("participantId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid participant id")
	}
	req := &requests.ConversationParticipantRequest{ConversationID: convID, ParticipantID: partID}
	if err := c.detachParticipant.Detach(req); err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

// @Summary Ingest Message
// @Description Persists a message authored by a participant in the conversation. Broker assigns the timestamp.
// @Tags Conversations
// @Accept json
// @Produce json
// @Param id path string true "Conversation ID"
// @Param request body requests.IngestMessageRequest true "Ingest Message Request"
// @Success 201 {object} responses.ConversationMessageResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 404 {object} errors.NotFoundError
// @Failure 500 {string} string "Internal Server Error"
// @Router /conversations/{id}/messages [post]
func (c *ConversationsController) Ingest(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}
	req := new(requests.IngestMessageRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	req.ConversationID = id
	msg, err := c.ingestMessage.Ingest(req)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if notFoundErr, ok := err.(*errors.NotFoundError); ok {
			return ctx.Status(fiber.StatusNotFound).JSON(notFoundErr)
		}
		c.logger.Error("internal error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return ctx.Status(fiber.StatusCreated).JSON(msg)
}
