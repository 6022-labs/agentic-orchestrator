package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	metrics_interfaces "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/metrics/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway_mvc/middlewares"
	metrics_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/metrics_mocks/interfaces_mocks"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type WhenRequiringAuthenticationTestingSuite struct {
	app                *fiber.App
	gatewayMetricsMock *metrics_mocks.MockGatewayMetricsInterface

	downstreamRan bool
}

func WhenRequiringAuthenticationBeforeEach(t *testing.T) *WhenRequiringAuthenticationTestingSuite {
	mockController := gomock.NewController(t)

	return &WhenRequiringAuthenticationTestingSuite{
		gatewayMetricsMock: metrics_mocks.NewMockGatewayMetricsInterface(mockController),
	}
}

// build mounts the middleware behind an injector stub that sets the same locals as the real ones.
func (s *WhenRequiringAuthenticationTestingSuite) build(injectedWallet any, injectedMethod any) {
	sut := middlewares.NewAuthenticationMandatoryMiddleware(zap.NewNop(), s.gatewayMetricsMock)

	s.app = fiber.New()
	s.app.Use(func(c fiber.Ctx) error {
		if injectedWallet != nil {
			c.Locals("authenticated-wallet", injectedWallet)
		}
		if injectedMethod != nil {
			c.Locals(middlewares.AuthMethodLocal, injectedMethod)
		}
		return c.Next()
	})
	s.app.Use(sut.Handle)
	s.app.Get("/protected", func(c fiber.Ctx) error {
		s.downstreamRan = true
		return c.SendStatus(fiber.StatusOK)
	})
}

func (s *WhenRequiringAuthenticationTestingSuite) get(t *testing.T) *http.Response {
	response, err := s.app.Test(httptest.NewRequest(http.MethodGet, "/protected", nil))
	assert.NoError(t, err)
	return response
}

func TestWhenRequiringAuthentication(t *testing.T) {
	t.Parallel()

	const wallet = "0x4bca9E17a8326d283d1fC68A725B46b42c2D9146"

	t.Run("Given the jwt injector authenticated the request", func(t *testing.T) {
		t.Parallel()

		t.Run("Should pass it downstream and record the jwt method", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequiringAuthenticationBeforeEach(t)
			suite.gatewayMetricsMock.EXPECT().RecordAuth(gomock.Any(), metrics_interfaces.AuthMethodJwt, true)
			suite.build(wallet, metrics_interfaces.AuthMethodJwt)

			response := suite.get(t)

			assert.Equal(t, fiber.StatusOK, response.StatusCode)
			assert.True(t, suite.downstreamRan)
		})
	})

	t.Run("Given the signature injector authenticated the request", func(t *testing.T) {
		t.Parallel()

		t.Run("Should pass it downstream and record the signature method", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequiringAuthenticationBeforeEach(t)
			suite.gatewayMetricsMock.EXPECT().RecordAuth(gomock.Any(), metrics_interfaces.AuthMethodSignature, true)
			suite.build(wallet, metrics_interfaces.AuthMethodSignature)

			response := suite.get(t)

			assert.Equal(t, fiber.StatusOK, response.StatusCode)
			assert.True(t, suite.downstreamRan)
		})
	})

	t.Run("Given no injector authenticated the request", func(t *testing.T) {
		t.Parallel()

		t.Run("Should reject it and record a failure", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequiringAuthenticationBeforeEach(t)
			suite.gatewayMetricsMock.EXPECT().RecordAuth(gomock.Any(), metrics_interfaces.AuthMethodNone, false)
			suite.build(nil, nil)

			response := suite.get(t)

			assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
			assert.False(t, suite.downstreamRan)
		})
	})

	t.Run("Given a wallet was injected without its method", func(t *testing.T) {
		t.Parallel()

		t.Run("Should pass it downstream and fall back to the none method", func(t *testing.T) {
			t.Parallel()

			suite := WhenRequiringAuthenticationBeforeEach(t)
			suite.gatewayMetricsMock.EXPECT().RecordAuth(gomock.Any(), metrics_interfaces.AuthMethodNone, true)
			suite.build(wallet, nil)

			response := suite.get(t)

			assert.Equal(t, fiber.StatusOK, response.StatusCode)
			assert.True(t, suite.downstreamRan)
		})
	})
}
