package slack_bridge_external_mvc

import (
	"encoding/json"
	"errors"

	common_errors "github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	"github.com/6022-labs/agentic-orchestrator/channels-and-integrations/slack-bridge/src/slack_bridge_external_mvc/requests"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const (
	ADD_AGENT_TO_CHANNEL_COMMAND_ENDPOINT           = "/command-handler/add-agent-to-channel"
	LIST_AGENTS_IN_CHANNEL_COMMAND_ENDPOINT         = "/command-handler/list-agents-in-channel"
	REMOVE_AGENT_FROM_CHANNEL_COMMAND_ENDPOINT      = "/command-handler/remove-agent-from-channel"
	VERIFY_WALLET_OWNERSHIP_COMMAND_ENDPOINT        = "/command-handler/verify-wallet-ownership"
	PAIR_SLACK_USER_TO_HUMAN_AGENT_COMMAND_ENDPOINT = "/command-handler/pair-slack-user-to-human-agent"
)

// Slack only shows a slash-command reply body on HTTP 200, so command outcomes are always 200 + ephemeral text.
const (
	FAILED_TO_HANDLE_COMMAND_REPLY     = "Failed to handle command"
	INVALID_AGENT_COMMAND_FORMAT_REPLY = "Invalid command format, expected: [collection address] [token ID]"
	INVALID_PAIR_COMMAND_FORMAT_REPLY  = "Invalid command format, expected: " +
		"[slack user ID] [human agent collection address] [human agent collection token ID]"
)

type CommandHandlerController struct {
	logger                    *zap.Logger
	addAgentToChannel         command_handlers.AddAgentToChannelInterface
	listAgentsInChannel       command_handlers.ListAgentInChannelInterface
	verifyWalletOwnership     command_handlers.VerifyWalletOwnershipInterface
	removeAgentFromChannel    command_handlers.RemoveAgentFromChannelInterface
	pairSlackUserToHumanAgent command_handlers.PairSlackUserToHumanAgentInterface
}

func NewCommandHandlerController(
	logger *zap.Logger,
	addAgentToChannel command_handlers.AddAgentToChannelInterface,
	listAgentsInChannel command_handlers.ListAgentInChannelInterface,
	verifyWalletOwnership command_handlers.VerifyWalletOwnershipInterface,
	removeAgentFromChannel command_handlers.RemoveAgentFromChannelInterface,
	pairSlackUserToHumanAgent command_handlers.PairSlackUserToHumanAgentInterface,
) *CommandHandlerController {
	return &CommandHandlerController{
		logger:                    logger,
		addAgentToChannel:         addAgentToChannel,
		listAgentsInChannel:       listAgentsInChannel,
		verifyWalletOwnership:     verifyWalletOwnership,
		removeAgentFromChannel:    removeAgentFromChannel,
		pairSlackUserToHumanAgent: pairSlackUserToHumanAgent,
	}
}

func (controller *CommandHandlerController) InitRoutes(app fiber.Router) {
	app.Post(ADD_AGENT_TO_CHANNEL_COMMAND_ENDPOINT, controller.HandleAddAgentToChannelCommand)
	app.Post(LIST_AGENTS_IN_CHANNEL_COMMAND_ENDPOINT, controller.HandleListAgentsInChannelCommand)
	app.Post(VERIFY_WALLET_OWNERSHIP_COMMAND_ENDPOINT, controller.HandleVerifyWalletOwnershipCommand)
	app.Post(REMOVE_AGENT_FROM_CHANNEL_COMMAND_ENDPOINT, controller.HandleRemoveAgentFromChannelCommand)
	app.Post(PAIR_SLACK_USER_TO_HUMAN_AGENT_COMMAND_ENDPOINT, controller.HandlePairSlackUserToHumanAgentCommand)
}

// HandleAddAgentToChannelCommand godoc
// @Summary      Add agent to channel
// @Description  Handles the Slack command to add an agent to a channel by collection address and token ID
// @Tags         Command
// @Accept       x-www-form-urlencoded
// @Produce      plain
// @Param        channel_id formData string true "Slack channel ID"
// @Param        text formData string true "Command text: [collection address] [token ID]"
// @Param        user_id formData string true "Slack user ID"
// @Success      200 {string} string "Ephemeral reply: success, usage hint on malformed text, or failure notice"
// @Failure      400 {string} string "Invalid request body"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /command-handler/add-agent-to-channel [post]
func (controller *CommandHandlerController) HandleAddAgentToChannelCommand(c fiber.Ctx) error {
	controller.logger.Info("Received add agent command request")

	var request requests.CommandRequest
	if err := c.Bind().Body(&request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	controller.logger.Debug("Request", zap.Any("request", request))

	if err := controller.addAgentToChannel.Add(c.Context(), request.ToAddAgentToChannelRequest()); err != nil {
		if errors.Is(err, common_errors.ErrValidation) {
			controller.logger.Info("Invalid add agent command", zap.Error(err))
			return c.Status(fiber.StatusOK).SendString(INVALID_AGENT_COMMAND_FORMAT_REPLY)
		}
		controller.logger.Error("Failed to handle add agent command", zap.Error(err))
		return c.Status(fiber.StatusOK).SendString(FAILED_TO_HANDLE_COMMAND_REPLY)
	}

	return c.Status(fiber.StatusOK).SendString("Agent added to channel successfully")
}

// HandleListAgentsInChannelCommand godoc
// @Summary      List agents in channel
// @Description  Handles the Slack command to list all agents in a channel
// @Tags         Command
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        channel_id formData string true "Slack channel ID"
// @Param        text formData string false "Command text"
// @Param        user_id formData string true "Slack user ID"
// @Success      200 {string} string "JSON list of agents, or failure notice"
// @Failure      400 {string} string "Invalid request body"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /command-handler/list-agents-in-channel [post]
func (controller *CommandHandlerController) HandleListAgentsInChannelCommand(c fiber.Ctx) error {
	controller.logger.Info("Received list agents command request")

	var request requests.CommandRequest
	if err := c.Bind().Body(&request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	controller.logger.Debug("Request", zap.Any("request", request))

	result, err := controller.listAgentsInChannel.List(c.Context(), request.ToListAgentInChannelRequest())
	if err != nil {
		controller.logger.Error("Failed to handle list agents command", zap.Error(err))
		return c.Status(fiber.StatusOK).SendString(FAILED_TO_HANDLE_COMMAND_REPLY)
	}

	controller.logger.Debug("Response", zap.Any("result", result))

	jsonResult, err := json.Marshal(result)
	if err != nil {
		controller.logger.Error("Failed to marshal response", zap.Error(err))
		return c.Status(fiber.StatusOK).SendString(FAILED_TO_HANDLE_COMMAND_REPLY)
	}

	return c.Status(fiber.StatusOK).SendString(string(jsonResult))
}

// HandleRemoveAgentFromChannelCommand godoc
// @Summary      Remove agent from channel
// @Description  Handles the Slack command to remove an agent from a channel by collection address and token ID
// @Tags         Command
// @Accept       x-www-form-urlencoded
// @Produce      plain
// @Param        channel_id formData string true "Slack channel ID"
// @Param        text formData string true "Command text: [collection address] [token ID]"
// @Param        user_id formData string true "Slack user ID"
// @Success      200 {string} string "Ephemeral reply: success, usage hint on malformed text, or failure notice"
// @Failure      400 {string} string "Invalid request body"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /command-handler/remove-agent-from-channel [post]
func (controller *CommandHandlerController) HandleRemoveAgentFromChannelCommand(c fiber.Ctx) error {
	controller.logger.Info("Received remove agent command request")

	var request requests.CommandRequest
	if err := c.Bind().Body(&request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	controller.logger.Debug("Request", zap.Any("request", request))

	if err := controller.removeAgentFromChannel.Remove(
		c.Context(),
		request.ToRemoveAgentFromChannelRequest(),
	); err != nil {
		if errors.Is(err, common_errors.ErrValidation) {
			controller.logger.Info("Invalid remove agent command", zap.Error(err))
			return c.Status(fiber.StatusOK).SendString(INVALID_AGENT_COMMAND_FORMAT_REPLY)
		}
		controller.logger.Error("Failed to handle remove agent command", zap.Error(err))
		return c.Status(fiber.StatusOK).SendString(FAILED_TO_HANDLE_COMMAND_REPLY)
	}

	return c.Status(fiber.StatusOK).SendString("Agent removed from channel successfully")
}

// HandleVerifyWalletOwnershipCommand godoc
// @Summary      Verify wallet ownership
// @Description  Handles the Slack command to initiate wallet ownership verification for a user
// @Tags         Command
// @Accept       x-www-form-urlencoded
// @Produce      plain
// @Param        channel_id formData string true "Slack channel ID"
// @Param        text formData string false "Command text"
// @Param        user_id formData string true "Slack user ID"
// @Success      200 {string} string "Verification response, or failure notice"
// @Failure      400 {string} string "Invalid request body"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /command-handler/verify-wallet-ownership [post]
func (controller *CommandHandlerController) HandleVerifyWalletOwnershipCommand(c fiber.Ctx) error {
	controller.logger.Info("Received verify wallet ownership command request")

	var request requests.CommandRequest
	if err := c.Bind().Body(&request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	controller.logger.Debug("Request", zap.Any("request", request))

	resp, err := controller.verifyWalletOwnership.Verify(c.Context(), request.ToVerifyWalletOwnershipRequest())
	if err != nil {
		controller.logger.Error("Failed to verify wallet ownership", zap.Error(err))
		return c.Status(fiber.StatusOK).SendString(FAILED_TO_HANDLE_COMMAND_REPLY)
	}

	return c.Status(fiber.StatusOK).SendString(*resp)
}

// HandlePairSlackUserToHumanAgentCommand godoc
// @Summary      Pair Slack user to human agent
// @Description  Handles the Slack command to pair a Slack user with a human agent by collection address and token ID
// @Tags         Command
// @Accept       x-www-form-urlencoded
// @Produce      plain
// @Param        channel_id formData string true "Slack channel ID"
// @Param        text formData string true "Command text: [slack user ID] [collection address] [token ID]"
// @Param        user_id formData string true "Slack user ID"
// @Success      200 {string} string "Ephemeral reply: success, usage hint on malformed text, or failure notice"
// @Failure      400 {string} string "Invalid request body"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /command-handler/pair-slack-user-to-human-agent [post]
func (controller *CommandHandlerController) HandlePairSlackUserToHumanAgentCommand(c fiber.Ctx) error {
	controller.logger.Info("Received pair slack user to human agent command request")

	var request requests.CommandRequest
	if err := c.Bind().Body(&request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	controller.logger.Debug("Request", zap.Any("request", request))

	if err := controller.pairSlackUserToHumanAgent.Pair(
		c.Context(),
		request.ToPairSlackUserToHumanAgentRequest(),
	); err != nil {
		if errors.Is(err, common_errors.ErrValidation) {
			controller.logger.Info("Invalid pair slack user to human agent command", zap.Error(err))
			return c.Status(fiber.StatusOK).SendString(INVALID_PAIR_COMMAND_FORMAT_REPLY)
		}
		controller.logger.Error("Failed to handle pair slack user to human agent command", zap.Error(err))
		return c.Status(fiber.StatusOK).SendString(FAILED_TO_HANDLE_COMMAND_REPLY)
	}

	return c.Status(fiber.StatusOK).SendString("Slack user paired to human agent successfully")
}
