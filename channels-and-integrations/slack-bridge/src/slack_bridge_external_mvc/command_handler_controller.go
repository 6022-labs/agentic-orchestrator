package slack_bridge_external_mvc

import (
	"encoding/json"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external/command_handlers"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/channels-and-integrations/slack-bridge/src/slack_bridge_external_mvc/requests"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	ADD_AGENT_TO_CHANNEL_COMMAND_ENDPOINT           = "/command-handler/add-agent-to-channel"
	LIST_AGENTS_IN_CHANNEL_COMMAND_ENDPOINT         = "/command-handler/list-agents-in-channel"
	REMOVE_AGENT_FROM_CHANNEL_COMMAND_ENDPOINT      = "/command-handler/remove-agent-from-channel"
	VERIFY_WALLET_OWNERSHIP_COMMAND_ENDPOINT        = "/command-handler/verify-wallet-ownership"
	PAIR_SLACK_USER_TO_HUMAN_AGENT_COMMAND_ENDPOINT = "/command-handler/pair-slack-user-to-human-agent"
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

func (controller *CommandHandlerController) HandleAddAgentToChannelCommand(c *fiber.Ctx) error {
	controller.logger.Info("Received add agent command request")

	var request requests.CommandRequest
	if err := c.BodyParser(&request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	controller.logger.Debug("Request", zap.Any("request", request))

	addRequest, err := request.ToAddAgentToChannelRequest()
	if err != nil {
		controller.logger.Error("Failed to convert request to AddAgentToChannelRequest", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid command format, expected: [collection address] [token ID]")
	}

	if err := controller.addAgentToChannel.Add(addRequest); err != nil {
		controller.logger.Error("Failed to handle add agent command", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to handle command")
	}

	return c.Status(fiber.StatusOK).SendString("Agent added to channel successfully")
}

func (controller *CommandHandlerController) HandleListAgentsInChannelCommand(c *fiber.Ctx) error {
	controller.logger.Info("Received list agents command request")

	var request requests.CommandRequest
	if err := c.BodyParser(&request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	controller.logger.Debug("Request", zap.Any("request", request))

	listRequest := request.ToListAgentInChannelRequest()
	if listRequest == nil {
		controller.logger.Error("Failed to convert request to ListAgentInChannelRequest")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid command format")
	}

	result, err := controller.listAgentsInChannel.List(listRequest)
	if err != nil {
		controller.logger.Error("Failed to handle list agents command", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to handle command")
	}

	controller.logger.Debug("Response", zap.Any("result", result))

	jsonResult, err := json.Marshal(result)
	if err != nil {
		controller.logger.Error("Failed to marshal response", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to process response")
	}

	return c.Status(fiber.StatusOK).SendString(string(jsonResult))
}

func (controller *CommandHandlerController) HandleRemoveAgentFromChannelCommand(c *fiber.Ctx) error {
	controller.logger.Info("Received remove agent command request")

	var request requests.CommandRequest
	if err := c.BodyParser(&request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	controller.logger.Debug("Request", zap.Any("request", request))

	removeRequest, err := request.ToRemoveAgentFromChannelRequest()
	if err != nil {
		controller.logger.Error("Failed to convert request to RemoveAgentFromChannelRequest", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid command format, expected: [collection address] [token ID]")
	}

	if err := controller.removeAgentFromChannel.Remove(removeRequest); err != nil {
		controller.logger.Error("Failed to handle remove agent command", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to handle command")
	}

	return c.Status(fiber.StatusOK).SendString("Agent removed from channel successfully")
}

func (controller *CommandHandlerController) HandleVerifyWalletOwnershipCommand(c *fiber.Ctx) error {
	controller.logger.Info("Received verify wallet ownership command request")

	var request requests.CommandRequest
	if err := c.BodyParser(&request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	controller.logger.Debug("Request", zap.Any("request", request))

	verifyWalletOwnershipRequest := request.ToVerifyWalletOwnershipRequest()
	if verifyWalletOwnershipRequest == nil {
		controller.logger.Error("Failed to convert request to VerifyWalletOwnershipRequest")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid command format")
	}

	resp, err := controller.verifyWalletOwnership.Verify(verifyWalletOwnershipRequest)
	if err != nil {
		controller.logger.Error("Failed to verify wallet ownership", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to handle command")
	}

	return c.Status(fiber.StatusOK).SendString(*resp)
}

func (controller *CommandHandlerController) HandlePairSlackUserToHumanAgentCommand(c *fiber.Ctx) error {
	controller.logger.Info("Received pair slack user to human agent command request")

	var request requests.CommandRequest
	if err := c.BodyParser(&request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	controller.logger.Debug("Request", zap.Any("request", request))

	pairRequest, err := request.ToPairSlackUserToHumanAgentRequest()
	if err != nil {
		controller.logger.Error("Failed to convert request to PairSlackUserToHumanAgentRequest", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid command format, expected: [slack user ID] [human agent collection address] [human agent collection token ID]")
	}

	if err := controller.pairSlackUserToHumanAgent.Pair(pairRequest); err != nil {
		controller.logger.Error("Failed to handle pair slack user to human agent command", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to handle command")
	}

	return c.Status(fiber.StatusOK).SendString("Slack user paired to human agent successfully")
}
