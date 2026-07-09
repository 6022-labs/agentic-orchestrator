package orchestrator_gateway_mvc

import (
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	_ "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	AUTH_CHALLENGE_ENDPOINT = "/auth/challenge"
	AUTH_VERIFY_ENDPOINT    = "/auth/verify"
)

type AuthController struct {
	logger           *zap.Logger
	getAuthChallenge use_cases.GetAuthChallengeInterface
	verifyAuth       use_cases.VerifyAuthInterface
}

func NewAuthController(
	logger *zap.Logger,
	getAuthChallenge use_cases.GetAuthChallengeInterface,
	verifyAuth use_cases.VerifyAuthInterface,
) *AuthController {
	return &AuthController{
		logger:           logger,
		getAuthChallenge: getAuthChallenge,
		verifyAuth:       verifyAuth,
	}
}

func (controller *AuthController) InitAuthenticatedRoutes(app fiber.Router) {
	// Just to satisfy the interface
}

func (controller *AuthController) InitUnauthenticatedRoutes(app fiber.Router) {
	app.Get(AUTH_CHALLENGE_ENDPOINT, controller.GetAuthChallenge)
	app.Post(AUTH_VERIFY_ENDPOINT, controller.VerifyAuth)
}

// GetAuthChallenge godoc
// @Summary      Get authentication challenge
// @Description  Returns a single-use nonce the wallet must sign
// @Tags         authentication
// @Produce      json
// @Success      200 {object} responses.GetAuthChallengeResponse
// @Failure      500 {string} string "Internal Server Error"
// @Router       /auth/challenge [get]
func (controller *AuthController) GetAuthChallenge(c *fiber.Ctx) error {
	challenge, err := controller.getAuthChallenge.Get()
	if err != nil {
		controller.logger.Error("Failed to generate auth challenge", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return c.JSON(challenge)
}

// VerifyAuth godoc
// @Summary      Verify signed challenge
// @Description  Verifies the signed nonce and returns a JWT token
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body requests.VerifyAuthRequest true "verify request"
// @Success      200 {object} responses.VerifyAuthResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      500 {string} string "Internal Server Error"
// @Router       /auth/verify [post]
func (controller *AuthController) VerifyAuth(c *fiber.Ctx) error {
	request := new(requests.VerifyAuthRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	response, err := controller.verifyAuth.Verify(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		controller.logger.Error("Failed to verify auth", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(response)
}
