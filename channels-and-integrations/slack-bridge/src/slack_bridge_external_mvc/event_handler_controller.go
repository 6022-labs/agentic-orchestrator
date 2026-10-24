package slack_bridge_external_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/event_handlers/requests"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const (
	EVENT_HANDLER_ENDPOINT = "/event-handler"
)

type EventHandlerController struct {
	logger                  *zap.Logger
	channelLeftHandler      event_handlers.ChannelLeftHandlerInterface
	messageCreatedHandler   event_handlers.MessageCreatedHandlerInterface
	messageChangedHandler   event_handlers.MessageChangedHandlerInterface
	messageDeletedHandler   event_handlers.MessageDeletedHandlerInterface
	channelCreatedHandler   event_handlers.ChannelCreatedHandlerInterface
	channelDeletedHandler   event_handlers.ChannelDeletedHandlerInterface
	channelArchiveHandler   event_handlers.ChannelArchiveHandlerInterface
	channelIdChangedHandler event_handlers.ChannelIdChangedHandlerInterface
}

func NewEventHandlerController(
	logger *zap.Logger,
	channelLeftHandler event_handlers.ChannelLeftHandlerInterface,
	messageCreatedHandler event_handlers.MessageCreatedHandlerInterface,
	messageChangedHandler event_handlers.MessageChangedHandlerInterface,
	messageDeletedHandler event_handlers.MessageDeletedHandlerInterface,
	channelCreatedHandler event_handlers.ChannelCreatedHandlerInterface,
	channelDeletedHandler event_handlers.ChannelDeletedHandlerInterface,
	channelArchiveHandler event_handlers.ChannelArchiveHandlerInterface,
	channelIdChangedHandler event_handlers.ChannelIdChangedHandlerInterface,
) *EventHandlerController {
	return &EventHandlerController{
		logger:                  logger,
		channelLeftHandler:      channelLeftHandler,
		messageCreatedHandler:   messageCreatedHandler,
		messageChangedHandler:   messageChangedHandler,
		messageDeletedHandler:   messageDeletedHandler,
		channelCreatedHandler:   channelCreatedHandler,
		channelDeletedHandler:   channelDeletedHandler,
		channelArchiveHandler:   channelArchiveHandler,
		channelIdChangedHandler: channelIdChangedHandler,
	}
}

func (controller *EventHandlerController) InitRoutes(app fiber.Router) {
	app.Post(EVENT_HANDLER_ENDPOINT, controller.HandleEventRequest)
}

// HandleEventRequest godoc
// @Summary      Handle Slack event
// @Description  Receives and dispatches Slack event callbacks including messages, channel events, and URL verification
// @Tags         Event
// @Accept       json
// @Produce      plain
// @Param        request body map[string]interface{} true "Slack event payload"
// @Success      200 {string} string "OK"
// @Failure      400 {string} string "Bad Request"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /event-handler [post]
func (controller *EventHandlerController) HandleEventRequest(c fiber.Ctx) error {
	controller.logger.Info("Received event request")

	// Parse the body as a generic map to check the "type" field
	var bodyMap map[string]interface{}
	if err := c.Bind().Body(&bodyMap); err != nil {
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", bodyMap))

	requestTypeStr, ok := bodyMap["type"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("Missing or invalid type field")
	}

	if requestTypeStr == "url_verification" {
		return controller.handleUrlVerification(c, bodyMap)
	}

	if requestTypeStr == "event_callback" {
		eventBody, ok := bodyMap["event"].(map[string]interface{})
		if !ok {
			return c.Status(fiber.StatusBadRequest).SendString("Missing or invalid event payload")
		}

		eventType, ok := eventBody["type"].(string)
		if !ok {
			return c.Status(fiber.StatusBadRequest).SendString("Missing or invalid event type field")
		}

		if eventType == "message" {
			return controller.handleMessageEvent(c, bodyMap)
		}

		if eventType == "channel_created" {
			return controller.handleChannelCreatedEvent(c)
		}

		if eventType == "channel_deleted" {
			return controller.handleChannelDeletedEvent(c)
		}

		if eventType == "channel_archive" {
			return controller.handleChannelArchiveEvent(c)
		}

		if eventType == "channel_left" {
			return controller.handleChannelLeftEvent(c)
		}

		if eventType == "channel_id_changed" {
			return controller.handleChannelIdChangedEvent(c)
		}

		return c.Status(fiber.StatusBadRequest).SendString("Unsupported event type")
	}

	return c.Status(fiber.StatusBadRequest).SendString("Unsupported request type")
}

func (controller *EventHandlerController) handleUrlVerification(c fiber.Ctx, bodyMap map[string]interface{}) error {
	controller.logger.Info("Received URL verification request")

	if challenge, ok := bodyMap["challenge"].(string); ok {
		return c.SendString(challenge)
	}

	return c.Status(fiber.StatusBadRequest).SendString("Missing challenge field")
}

func (controller *EventHandlerController) handleMessageEvent(c fiber.Ctx, bodyMap map[string]interface{}) error {
	subtype, _ := bodyMap["event"].(map[string]interface{})["subtype"].(string)

	if subtype == "" || subtype == "bot_message" {
		return controller.handleMessageCreatedEvent(c)
	}

	if subtype == "message_changed" {
		return controller.handleMessageChangedEvent(c)
	}

	if subtype == "message_deleted" {
		return controller.handleMessageDeletedEvent(c)
	}

	controller.logger.Info("Received message event that is not supported", zap.String("subtype", subtype))
	return c.SendStatus(fiber.StatusOK)

}

func (controller *EventHandlerController) handleMessageCreatedEvent(c fiber.Ctx) error {
	request := new(requests.MessageCreatedRequest)
	if err := c.Bind().Body(request); err != nil {
		return err
	}

	err := controller.messageCreatedHandler.Handle(c.Context(), request)
	if err != nil {
		controller.logger.Error("Failed to handle event", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully handled message event")

	return c.SendStatus(fiber.StatusOK)
}

func (controller *EventHandlerController) handleMessageChangedEvent(c fiber.Ctx) error {
	request := new(requests.MessageChangedRequest)
	if err := c.Bind().Body(request); err != nil {
		return err
	}

	err := controller.messageChangedHandler.Handle(c.Context(), request)
	if err != nil {
		controller.logger.Error("Failed to handle event", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully handled message event")

	return c.SendStatus(fiber.StatusOK)
}

func (controller *EventHandlerController) handleMessageDeletedEvent(c fiber.Ctx) error {
	request := new(requests.MessageDeletedRequest)
	if err := c.Bind().Body(request); err != nil {
		return err
	}

	err := controller.messageDeletedHandler.Handle(c.Context(), request)
	if err != nil {
		controller.logger.Error("Failed to handle event", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully handled message event")

	return c.SendStatus(fiber.StatusOK)
}

func (controller *EventHandlerController) handleChannelCreatedEvent(c fiber.Ctx) error {
	request := new(requests.ChannelCreatedRequest)
	if err := c.Bind().Body(request); err != nil {
		return err
	}

	err := controller.channelCreatedHandler.Handle(c.Context(), request)
	if err != nil {
		controller.logger.Error("Failed to handle event", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully handled channel created event")

	return c.SendStatus(fiber.StatusOK)
}

func (controller *EventHandlerController) handleChannelDeletedEvent(c fiber.Ctx) error {
	request := new(requests.ChannelDeletedRequest)
	if err := c.Bind().Body(request); err != nil {
		return err
	}

	err := controller.channelDeletedHandler.Handle(c.Context(), request)
	if err != nil {
		controller.logger.Error("Failed to handle event", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully handled channel deleted event")

	return c.SendStatus(fiber.StatusOK)
}

func (controller *EventHandlerController) handleChannelArchiveEvent(c fiber.Ctx) error {
	request := new(requests.ChannelArchiveRequest)
	if err := c.Bind().Body(request); err != nil {
		return err
	}

	err := controller.channelArchiveHandler.Handle(c.Context(), request)
	if err != nil {
		controller.logger.Error("Failed to handle event", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully handled channel archive event")

	return c.SendStatus(fiber.StatusOK)
}

func (controller *EventHandlerController) handleChannelLeftEvent(c fiber.Ctx) error {
	request := new(requests.ChannelLeftRequest)
	if err := c.Bind().Body(request); err != nil {
		return err
	}

	err := controller.channelLeftHandler.Handle(c.Context(), request)
	if err != nil {
		controller.logger.Error("Failed to handle event", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully handled channel left event")

	return c.SendStatus(fiber.StatusOK)
}

func (controller *EventHandlerController) handleChannelIdChangedEvent(c fiber.Ctx) error {
	request := new(requests.ChannelIdChangedRequest)
	if err := c.Bind().Body(request); err != nil {
		return err
	}

	err := controller.channelIdChangedHandler.Handle(c.Context(), request)
	if err != nil {
		controller.logger.Error("Failed to handle event", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Successfully handled channel id changed event")

	return c.SendStatus(fiber.StatusOK)
}
