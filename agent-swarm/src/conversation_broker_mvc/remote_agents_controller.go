package conversation_broker_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

const (
	GET_AGENT_CARD_ENDPOINT = "/remote-agents/:ensDomain/card"
	GET_AGENT_URL_ENDPOINT  = "/remote-agents/:ensDomain/url"
)

type RemoteAgentsController struct {
	logger             *zap.Logger
	getRemoteAgentCard use_cases.GetRemoteAgentCardInterface
	resolveAgentUrl    use_cases.ResolveAgentUrlInterface
}

func NewRemoteAgentsController(
	logger *zap.Logger,
	getRemoteAgentCard use_cases.GetRemoteAgentCardInterface,
	resolveAgentUrl use_cases.ResolveAgentUrlInterface,
) *RemoteAgentsController {
	return &RemoteAgentsController{
		logger:             logger,
		getRemoteAgentCard: getRemoteAgentCard,
		resolveAgentUrl:    resolveAgentUrl,
	}
}

func (c *RemoteAgentsController) InitRoutes(app fiber.Router) {
	app.Get(GET_AGENT_CARD_ENDPOINT, c.GetCard)
	app.Get(GET_AGENT_URL_ENDPOINT, c.GetUrl)
}

// @Summary Get Agent Node Card
// @Description Resolves the agent's runtime from the AgentEnsRegistry and returns its cached /.well-known/6022 card.
// @Tags Agent Cards
// @Produce json
// @Param ensDomain path string true "Agent ENS domain"
// @Success 200 {object} responses.RemoteAgentCardResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 502 {string} string "Bad Gateway"
// @Router /remote-agents/{ensDomain}/card [get]
func (c *RemoteAgentsController) GetCard(ctx fiber.Ctx) error {
	card, err := c.getRemoteAgentCard.Get(ctx.Context(), ctx.Params("ensDomain"))
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		c.logger.Warn("failed to fetch agent node card", zap.Error(err))
		return ctx.Status(fiber.StatusBadGateway).SendString("failed to reach agent node")
	}
	return ctx.JSON(card)
}

// @Summary Get Agent Node URL
// @Description Resolves the agent's node URL from its ENS url text record in the AgentEnsRegistry.
// @Tags Agent Urls
// @Produce json
// @Param ensDomain path string true "Agent ENS domain"
// @Success 200 {object} responses.AgentUrlResponse
// @Failure 400 {object} errors.ValidationError
// @Failure 502 {string} string "Bad Gateway"
// @Router /remote-agents/{ensDomain}/url [get]
func (c *RemoteAgentsController) GetUrl(ctx fiber.Ctx) error {
	request := new(requests.ResolveAgentUrlRequest)
	if err := ctx.Bind().URI(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	url, err := c.resolveAgentUrl.Resolve(ctx.Context(), request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		c.logger.Warn("failed to resolve agent url", zap.Error(err))
		return ctx.Status(fiber.StatusBadGateway).SendString("failed to resolve agent url")
	}
	return ctx.JSON(url)
}
