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
// @Description  Returns a single-use EIP-4361 message + expiry the wallet must sign verbatim
// @Tags         authentication
// @Produce      json
// @Param        address query string true "wallet address the message is issued for"
// @Param        chainId query int true "chain id of the connected wallet"
// @Success      200 {object} responses.GetAuthChallengeResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      500 {string} string "Internal Server Error"
// @Router       /auth/challenge [get]
func (controller *AuthController) GetAuthChallenge(c *fiber.Ctx) error {
	request := new(requests.GetAuthChallengeRequest)
	if err := c.QueryParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	challenge, err := controller.getAuthChallenge.Execute(c.UserContext(), request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		controller.logger.Error("Failed to generate auth challenge", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return c.JSON(challenge)
}

// VerifyAuth godoc
// @Summary      Verify signed challenge
// @Description  Verifies the signed EIP-4361 message for the nonce and returns a JWT token
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body requests.VerifyAuthRequest true "verify request"
// @Success      200 {object} responses.VerifyAuthResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      401 {object} errors.UnauthorizedError "challenge_not_found | signature_mismatch"
// @Failure      500 {string} string "Internal Server Error"
// @Router       /auth/verify [post]
func (controller *AuthController) VerifyAuth(c *fiber.Ctx) error {
	request := new(requests.VerifyAuthRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	response, err := controller.verifyAuth.Execute(c.UserContext(), request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}
		if unauthorizedErr, ok := err.(*errors.UnauthorizedError); ok {
			return c.Status(fiber.StatusUnauthorized).JSON(unauthorizedErr)
		}
		controller.logger.Error("Failed to verify auth", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return c.JSON(response)
}
