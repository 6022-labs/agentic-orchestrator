package orchestrator_gateway_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

const (
	LIST_CONVERSATIONS_ENDPOINT             = "/conversations"
	GET_CONVERSATION_ENDPOINT               = "/conversations/:id"
	GET_CONVERSATION_MESSAGES_ENDPOINT      = "/conversations/:id/messages"
	LIST_CONVERSATION_PARTICIPANTS_ENDPOINT = "/conversations/:id/participants"
	FORCE_CONCLUDE_CONVERSATION_ENDPOINT    = "/conversations/:id/force-conclude"
)

type ConversationController struct {
	logger                       *zap.Logger
	listConversations            use_cases.ListConversationsInterface
	getConversation              use_cases.GetConversationInterface
	getConversationMessages      use_cases.GetConversationMessagesInterface
	listConversationParticipants use_cases.ListConversationParticipantsInterface
	forceConcludeConversation    use_cases.ForceConcludeConversationInterface
}

func NewConversationController(
	logger *zap.Logger,
	listConversations use_cases.ListConversationsInterface,
	getConversation use_cases.GetConversationInterface,
	getConversationMessages use_cases.GetConversationMessagesInterface,
	listConversationParticipants use_cases.ListConversationParticipantsInterface,
	forceConcludeConversation use_cases.ForceConcludeConversationInterface,
) *ConversationController {
	return &ConversationController{
		logger:                       logger,
		listConversations:            listConversations,
		getConversation:              getConversation,
		getConversationMessages:      getConversationMessages,
		listConversationParticipants: listConversationParticipants,
		forceConcludeConversation:    forceConcludeConversation,
	}
}

func (controller *ConversationController) InitAuthenticatedRoutes(app fiber.Router) {
	app.Get(LIST_CONVERSATIONS_ENDPOINT, controller.List)
	app.Get(GET_CONVERSATION_ENDPOINT, controller.Get)
	app.Get(GET_CONVERSATION_MESSAGES_ENDPOINT, controller.GetMessages)
	app.Get(LIST_CONVERSATION_PARTICIPANTS_ENDPOINT, controller.ListParticipants)
	app.Post(FORCE_CONCLUDE_CONVERSATION_ENDPOINT, controller.ForceConclude)
}

func (controller *ConversationController) InitUnauthenticatedRoutes(app fiber.Router) {
	// Just to satisfy the interface
}

// List godoc
// @Summary      List conversations
// @Tags         conversation
// @Produce      json
// @Param        page query int false "Page"
// @Param        elementsPerPage query int false "Elements per page"
// @Success      200 {object} responses.ListBrokerConversationsResponse
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /conversations [get]
func (controller *ConversationController) List(c *fiber.Ctx) error {
	var page, elementsPerPage *uint
	if v := c.QueryInt("page", 0); v > 0 {
		page = lo.ToPtr(uint(v))
	}
	if v := c.QueryInt("elementsPerPage", 0); v > 0 {
		elementsPerPage = lo.ToPtr(uint(v))
	}

	response, err := controller.listConversations.List(c.UserContext(), page, elementsPerPage)
	if err != nil {
		controller.logger.Error("Failed to list conversations", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(response)
}

// Get godoc
// @Summary      Get conversation
// @Tags         conversation
// @Produce      json
// @Param        id path string true "Conversation ID"
// @Success      200 {object} responses.BrokerConversationResponse
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /conversations/{id} [get]
func (controller *ConversationController) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}

	response, err := controller.getConversation.Get(c.UserContext(), id)
	if err != nil {
		controller.logger.Error("Failed to get conversation", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(response)
}

// GetMessages godoc
// @Summary      Get conversation messages
// @Tags         conversation
// @Produce      json
// @Param        id path string true "Conversation ID"
// @Success      200 {array} responses.BrokerMessageResponse
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /conversations/{id}/messages [get]
func (controller *ConversationController) GetMessages(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}

	response, err := controller.getConversationMessages.Get(c.UserContext(), id)
	if err != nil {
		controller.logger.Error("Failed to get conversation messages", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(response)
}

// ListParticipants godoc
// @Summary      List conversation participants
// @Tags         conversation
// @Produce      json
// @Param        id path string true "Conversation ID"
// @Success      200 {array} responses.BrokerParticipantResponse
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /conversations/{id}/participants [get]
func (controller *ConversationController) ListParticipants(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}

	response, err := controller.listConversationParticipants.List(c.UserContext(), id)
	if err != nil {
		controller.logger.Error("Failed to list conversation participants", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(response)
}

// ForceConclude godoc
// @Summary      Force conclude a conversation
// @Tags         conversation
// @Param        id path string true "Conversation ID"
// @Success      204 {string} string "No Content"
// @Failure      500 {string} string "Internal Server Error"
// @Security     BearerAuth
// @Router       /conversations/{id}/force-conclude [post]
func (controller *ConversationController) ForceConclude(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid conversation id")
	}

	if err := controller.forceConcludeConversation.Conclude(c.UserContext(), id); err != nil {
		controller.logger.Error("Failed to force conclude conversation", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
