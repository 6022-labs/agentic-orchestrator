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
	INIT_JWT_AUTHENTICATION_ENDPOINT   = "/auth/jwt/init"
	FINISH_JWT_AUTHENTICATION_ENDPOINT = "/auth/jwt/finish"
)

type AuthController struct {
	logger                      *zap.Logger
	finishJwtAuthentication     use_cases.FinishJwtAuthenticationInterface
	initializeJwtAuthentication use_cases.InitializeJwtAuthenticationInterface
}

func NewAuthController(
	logger *zap.Logger,
	finishJwtAuthentication use_cases.FinishJwtAuthenticationInterface,
	initializeJwtAuthentication use_cases.InitializeJwtAuthenticationInterface,
) *AuthController {
	return &AuthController{
		logger:                      logger,
		finishJwtAuthentication:     finishJwtAuthentication,
		initializeJwtAuthentication: initializeJwtAuthentication,
	}
}

func (controller *AuthController) InitAuthenticatedRoutes(app fiber.Router) {
	// Just to satisfy the interface
}

func (controller *AuthController) InitUnauthenticatedRoutes(app fiber.Router) {
	app.Post(INIT_JWT_AUTHENTICATION_ENDPOINT, controller.InitializeJwtAuthentication)
	app.Post(FINISH_JWT_AUTHENTICATION_ENDPOINT, controller.FinishJwtAuthentication)
}

// InitializeJwtAuthentication godoc
// @Summary      Start JWT authentication
// @Description  Initializes the JWT authentication flow by returning a nonce
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body requests.InitializeJwtAuthenticationRequest true "authentication request"
// @Success      200 {object} responses.InitializeJwtAuthenticationResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      500 {string} string "Internal Server Error"
// @Router       /auth/jwt/init [post]
func (controller *AuthController) InitializeJwtAuthentication(c *fiber.Ctx) error {
	controller.logger.Info("Initialize JWT authentication request received")
	controller.logger.Debug("Request parameters", zap.Any("body", c.Body()))

	request := new(requests.InitializeJwtAuthenticationRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	nonce, err := controller.initializeJwtAuthentication.Init(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to initialize JWT authentication", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Initialize JWT authentication request completed successfully")
	controller.logger.Debug("Response", zap.Any("nonce", nonce))

	return c.JSON(nonce)
}

// FinishJwtAuthentication godoc
// @Summary      Complete JWT authentication
// @Description  Verifies the signed nonce and returns a JWT token
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        request body requests.FinishJwtAuthenticationRequest true "authentication request"
// @Success      200 {object} responses.FinishJwtAuthenticationResponse
// @Failure      400 {object} errors.ValidationError
// @Failure      500 {string} string "Internal Server Error"
// @Router       /auth/jwt/finish [post]
func (controller *AuthController) FinishJwtAuthentication(c *fiber.Ctx) error {
	controller.logger.Info("Finish JWT authentication request received")
	controller.logger.Debug("Request parameters", zap.Any("body", c.Body()))

	request := new(requests.FinishJwtAuthenticationRequest)
	if err := c.BodyParser(request); err != nil {
		controller.logger.Error("Failed to parse request body", zap.Error(err))
		return err
	}

	controller.logger.Debug("Parsed request", zap.Any("request", request))

	response, err := controller.finishJwtAuthentication.Finish(request)
	if err != nil {
		if validationErr, ok := err.(*errors.ValidationError); ok {
			return c.Status(fiber.StatusBadRequest).JSON(validationErr)
		}

		controller.logger.Error("Failed to finish JWT authentication", zap.Error(err))

		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	controller.logger.Info("Finish JWT authentication request completed successfully")
	controller.logger.Debug("Response", zap.Any("response", response))

	return c.JSON(response)
}
