package conversation_broker_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-reasoning/src/conversation_broker/use_cases"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const GET_AGENT_CARD_ENDPOINT = "/agent-cards/:ensDomain"

type AgentCardsController struct {
	logger             *zap.Logger
	getRemoteAgentCard use_cases.GetRemoteAgentCardInterface
}

func NewAgentCardsController(
	logger *zap.Logger,
	getRemoteAgentCard use_cases.GetRemoteAgentCardInterface,
) *AgentCardsController {
	return &AgentCardsController{
		logger:             logger,
		getRemoteAgentCard: getRemoteAgentCard,
	}
}

func (c *AgentCardsController) InitRoutes(app fiber.Router) {
	app.Get(GET_AGENT_CARD_ENDPOINT, c.Get)
}

// @Summary Get Agent Node Card
// @Description Resolves the agent's runtime via ENS CCIP-Read and returns its cached /.well-known/6022 card.
// @Tags Agent Cards
// @Produce json
// @Param ensDomain path string true "Agent ENS domain"
// @Success 200 {object} models.RemoteAgentCard
// @Failure 400 {object} errors.ValidationError
// @Failure 502 {string} string "Bad Gateway"
// @Router /agent-cards/{ensDomain} [get]
func (c *AgentCardsController) Get(ctx *fiber.Ctx) error {
	card, err := c.getRemoteAgentCard.Get(ctx.Params("ensDomain"))
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		c.logger.Warn("failed to fetch agent node card", zap.Error(err))
		return ctx.Status(fiber.StatusBadGateway).SendString("failed to reach agent node")
	}
	return ctx.JSON(card)
}
