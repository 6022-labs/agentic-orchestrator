package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	gateway_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCreatingAgentSecondStepTestSuite struct {
	agentRegistryRequesterMock *gateway_interfaces_mocks.MockAgentRegistryRequesterInterface
	sut                        *use_cases.CreateAgentSecondStep
}

func WhenCreatingAgentSecondStepBeforeEach(t *testing.T) *WhenCreatingAgentSecondStepTestSuite {
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())
	agentRegistryRequesterMock := gateway_interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	sut := use_cases.NewCreateAgentSecondStep(agentRegistryRequesterMock)

	return &WhenCreatingAgentSecondStepTestSuite{
		sut:                        sut,
		agentRegistryRequesterMock: agentRegistryRequesterMock,
	}
}

func generateValidAgentByIdRequest() *requests.AgentByIdRequest {
	return &requests.AgentByIdRequest{
		AgentId:             uuid.New(),
		AuthenticatedWallet: common.HexToAddress("0xd7Df1101B5933522c6469e57110FD4901F4ac7Ca"),
	}
}

func TestWhenCreatingAgentSecondStep(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()

			request := generateValidAgentByIdRequest()
			request.AgentId = uuid.Nil
			suite := WhenCreatingAgentSecondStepBeforeEach(t)

			_, err := suite.sut.Execute(request)
			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given registry requester fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentSecondStepBeforeEach(t)
			request := generateValidAgentByIdRequest()

			suite.agentRegistryRequesterMock.EXPECT().
				CreateAgentSecondStep(gomock.Any()).
				Return(nil, assert.AnError).
				Times(1)

			_, err := suite.sut.Execute(request)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given registry requester succeeds", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return response", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentSecondStepBeforeEach(t)
			request := generateValidAgentByIdRequest()
			expected := &responses.CreateAgentSecondStepResponse{
				Nonce:     uuid.New().String(),
				Signer:    common.HexToAddress("0x76B86F1deaa7241F409b8cA98098b287888e750C").Hex(),
				Signature: "signature",
			}

			suite.agentRegistryRequesterMock.EXPECT().
				CreateAgentSecondStep(gomock.Any()).
				DoAndReturn(func(registryRequest *interfaces.CreateRegistryAgentSecondStepRequest) (*responses.CreateAgentSecondStepResponse, error) {
					assert.Equal(t, request.AuthenticatedWallet, registryRequest.AuthenticatedWallet)
					assert.Equal(t, request.AgentId, registryRequest.AgentId)
					return expected, nil
				}).
				Times(1)

			response, err := suite.sut.Execute(request)
			assert.Nil(t, err)
			assert.Equal(t, expected, response)
		})
	})
}
