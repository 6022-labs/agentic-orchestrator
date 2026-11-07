package conversation_broker_mvc

import (
	apperrors "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/errors"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/common/mvc"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker/use_cases/requests"
	"github.com/gofiber/fiber/v3"
)

const (
	GET_AGENT_CARD_ENDPOINT = "/remote-agents/:ensDomain/card"
	GET_AGENT_URL_ENDPOINT  = "/remote-agents/:ensDomain/url"
)

type RemoteAgentsController struct {
	errorWriter        *mvc.ErrorWriter
	getRemoteAgentCard use_cases.GetRemoteAgentCardInterface
	resolveAgentUrl    use_cases.ResolveAgentUrlInterface
}

func NewRemoteAgentsController(
	errorWriter *mvc.ErrorWriter,
	getRemoteAgentCard use_cases.GetRemoteAgentCardInterface,
	resolveAgentUrl use_cases.ResolveAgentUrlInterface,
) *RemoteAgentsController {
	return &RemoteAgentsController{
		errorWriter:        errorWriter,
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
// @Failure 502 {object} errors.UnavailableError "agent_node_unreachable"
// @Router /remote-agents/{ensDomain}/card [get]
func (c *RemoteAgentsController) GetCard(ctx fiber.Ctx) error {
	card, err := c.getRemoteAgentCard.Get(ctx.Context(), ctx.Params("ensDomain"))
	if err != nil {
		return c.errorWriter.Write(ctx, err)
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
// @Failure 400 {object} errors.BadRequestError "invalid_params"
// @Failure 502 {object} errors.UnavailableError "agent_url_unresolved"
// @Router /remote-agents/{ensDomain}/url [get]
func (c *RemoteAgentsController) GetUrl(ctx fiber.Ctx) error {
	request := new(requests.ResolveAgentUrlRequest)
	if err := ctx.Bind().URI(request); err != nil {
		return c.errorWriter.Write(ctx, apperrors.NewBadRequestError("invalid_params", err.Error()))
	}

	url, err := c.resolveAgentUrl.Resolve(ctx.Context(), request)
	if err != nil {
		return c.errorWriter.Write(ctx, err)
	}
	return ctx.JSON(url)
}
