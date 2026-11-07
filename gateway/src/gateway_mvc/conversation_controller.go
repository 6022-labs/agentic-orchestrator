package orchestrator_gateway_mvc

import (
	apperrors "github.com/6022-labs/agentic-orchestrator/gateway/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/gofiber/fiber/v3"
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
	errorWriter                  *mvc.ErrorWriter
	listConversations            use_cases.ListConversationsInterface
	getConversation              use_cases.GetConversationInterface
	getConversationMessages      use_cases.GetConversationMessagesInterface
	listConversationParticipants use_cases.ListConversationParticipantsInterface
	forceConcludeConversation    use_cases.ForceConcludeConversationInterface
}

func NewConversationController(
	logger *zap.Logger,
	errorWriter *mvc.ErrorWriter,
	listConversations use_cases.ListConversationsInterface,
	getConversation use_cases.GetConversationInterface,
	getConversationMessages use_cases.GetConversationMessagesInterface,
	listConversationParticipants use_cases.ListConversationParticipantsInterface,
	forceConcludeConversation use_cases.ForceConcludeConversationInterface,
) *ConversationController {
	return &ConversationController{
		logger:                       logger,
		errorWriter:                  errorWriter,
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
// @Failure      500 {object} errors.InternalError
// @Security     BearerAuth
// @Router       /conversations [get]
func (controller *ConversationController) List(c fiber.Ctx) error {
	var page, elementsPerPage *uint
	if v := fiber.Query[int](c, "page", 0); v > 0 {
		page = lo.ToPtr(uint(v))
	}
	if v := fiber.Query[int](c, "elementsPerPage", 0); v > 0 {
		elementsPerPage = lo.ToPtr(uint(v))
	}

	response, err := controller.listConversations.List(c.Context(), page, elementsPerPage)
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}

	return c.JSON(response)
}

// Get godoc
// @Summary      Get conversation
// @Description  Caller must control (own or be the evm address of) at least one participating agent
// @Tags         conversation
// @Produce      json
// @Param        id path string true "Conversation ID"
// @Success      200 {object} responses.BrokerConversationResponse
// @Failure      400 {object} errors.ValidationError "zero uuid; a malformed id answers BadRequestError invalid_params"
// @Failure      400 {object} errors.BadRequestError "invalid_params"
// @Failure      401 {object} errors.UnauthorizedError
// @Failure      403 {object} errors.ForbiddenError
// @Failure      404 {object} errors.NotFoundError
// @Failure      500 {object} errors.InternalError
// @Failure      502 {object} errors.UnavailableError
// @Security     BearerAuth
// @Router       /conversations/{id} [get]
func (controller *ConversationController) Get(c fiber.Ctx) error {
	request := new(requests.ConversationByIdRequest)
	if err := c.Bind().URI(request); err != nil {
		return controller.errorWriter.Write(
			c,
			apperrors.NewBadRequestError("invalid_params", "invalid conversation id"),
		)
	}

	response, err := controller.getConversation.Get(c.Context(), request)
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}

	return c.JSON(response)
}

// GetMessages godoc
// @Summary      Get conversation messages
// @Description  Caller must control (own or be the evm address of) at least one participating agent
// @Tags         conversation
// @Produce      json
// @Param        id path string true "Conversation ID"
// @Success      200 {array} responses.BrokerMessageResponse
// @Failure      400 {object} errors.ValidationError "zero uuid; a malformed id answers BadRequestError invalid_params"
// @Failure      400 {object} errors.BadRequestError "invalid_params"
// @Failure      401 {object} errors.UnauthorizedError
// @Failure      403 {object} errors.ForbiddenError
// @Failure      404 {object} errors.NotFoundError
// @Failure      500 {object} errors.InternalError
// @Failure      502 {object} errors.UnavailableError
// @Security     BearerAuth
// @Router       /conversations/{id}/messages [get]
func (controller *ConversationController) GetMessages(c fiber.Ctx) error {
	request := new(requests.ConversationByIdRequest)
	if err := c.Bind().URI(request); err != nil {
		return controller.errorWriter.Write(
			c,
			apperrors.NewBadRequestError("invalid_params", "invalid conversation id"),
		)
	}

	response, err := controller.getConversationMessages.Get(c.Context(), request)
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}

	return c.JSON(response)
}

// ListParticipants godoc
// @Summary      List conversation participants
// @Description  Caller must control (own or be the evm address of) at least one participating agent
// @Tags         conversation
// @Produce      json
// @Param        id path string true "Conversation ID"
// @Success      200 {array} responses.BrokerParticipantResponse
// @Failure      400 {object} errors.ValidationError "zero uuid; a malformed id answers BadRequestError invalid_params"
// @Failure      400 {object} errors.BadRequestError "invalid_params"
// @Failure      401 {object} errors.UnauthorizedError
// @Failure      403 {object} errors.ForbiddenError
// @Failure      404 {object} errors.NotFoundError
// @Failure      500 {object} errors.InternalError
// @Failure      502 {object} errors.UnavailableError
// @Security     BearerAuth
// @Router       /conversations/{id}/participants [get]
func (controller *ConversationController) ListParticipants(c fiber.Ctx) error {
	request := new(requests.ConversationByIdRequest)
	if err := c.Bind().URI(request); err != nil {
		return controller.errorWriter.Write(
			c,
			apperrors.NewBadRequestError("invalid_params", "invalid conversation id"),
		)
	}

	response, err := controller.listConversationParticipants.List(c.Context(), request)
	if err != nil {
		return controller.errorWriter.Write(c, err)
	}

	return c.JSON(response)
}

// ForceConclude godoc
// @Summary      Force conclude a conversation
// @Description  Caller must control (own or be the evm address of) at least one participating agent
// @Tags         conversation
// @Param        id path string true "Conversation ID"
// @Success      204 {string} string "No Content"
// @Failure      400 {object} errors.ValidationError "zero uuid; a malformed id answers BadRequestError invalid_params"
// @Failure      400 {object} errors.BadRequestError "invalid_params"
// @Failure      401 {object} errors.UnauthorizedError
// @Failure      403 {object} errors.ForbiddenError
// @Failure      404 {object} errors.NotFoundError
// @Failure      500 {object} errors.InternalError
// @Failure      502 {object} errors.UnavailableError
// @Security     BearerAuth
// @Router       /conversations/{id}/force-conclude [post]
func (controller *ConversationController) ForceConclude(c fiber.Ctx) error {
	request := new(requests.ConversationByIdRequest)
	if err := c.Bind().URI(request); err != nil {
		return controller.errorWriter.Write(
			c,
			apperrors.NewBadRequestError("invalid_params", "invalid conversation id"),
		)
	}

	if err := controller.forceConcludeConversation.Conclude(c.Context(), request); err != nil {
		return controller.errorWriter.Write(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
