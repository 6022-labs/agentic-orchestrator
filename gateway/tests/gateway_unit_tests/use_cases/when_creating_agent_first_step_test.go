package use_cases_test

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/errors"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/services/interfaces"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/use_cases/responses"
	"github.com/6022-labs/agentic-orchestrator/gateway/tests/common_tests"
	gateway_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/gateway/tests/gateway_mocks/services_mocks/interfaces_mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCreatingAgentFirstStepTestSuite struct {
	agentRegistryRequesterMock *gateway_interfaces_mocks.MockAgentRegistryRequesterInterface
	sut                        *use_cases.CreateAgentFirstStep
}

func WhenCreatingAgentFirstStepBeforeEach(t *testing.T) *WhenCreatingAgentFirstStepTestSuite {
	mockController := gomock.NewController(t, gomock.WithOverridableExpectations())

	agentRegistryRequesterMock := gateway_interfaces_mocks.NewMockAgentRegistryRequesterInterface(mockController)
	sut := use_cases.NewCreateAgentFirstStep(agentRegistryRequesterMock)

	return &WhenCreatingAgentFirstStepTestSuite{
		agentRegistryRequesterMock: agentRegistryRequesterMock,
		sut:                        sut,
	}
}

func createMultipartFileHeader(filePath string) (*multipart.FileHeader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buff bytes.Buffer
	buffWriter := io.Writer(&buff)
	formWriter := multipart.NewWriter(buffWriter)

	formPart, err := formWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(formPart, file); err != nil {
		return nil, err
	}

	if err := formWriter.Close(); err != nil {
		return nil, err
	}

	buffReader := bytes.NewReader(buff.Bytes())
	formReader := multipart.NewReader(buffReader, formWriter.Boundary())
	multipartForm, err := formReader.ReadForm(1 << 20)
	if err != nil {
		return nil, err
	}

	files, exists := multipartForm.File["file"]
	if !exists || len(files) == 0 {
		return nil, err
	}

	return files[0], nil
}

func generateValidCreateAgentFirstStepRequest() *requests.CreateAgentFirstStepRequest {
	tempFile, err := os.CreateTemp("", "agent_image_*.png")
	if err != nil {
		log.Fatalf("failed to create temp file: %v", err)
	}
	defer tempFile.Close()

	tempFile.Write([]byte("dummy content"))

	file, err := createMultipartFileHeader(tempFile.Name())
	if err != nil {
		log.Fatalf("failed to create multipart file header: %v", err)
	}

	return &requests.CreateAgentFirstStepRequest{
		Name:                "agent-name",
		Image:               file,
		Role:                "facilitator",
		Mission:             "agent mission",
		CollectionAddress:   common_tests.RandomEthAddress().Hex(),
		AuthenticatedWallet: common_tests.RandomEthAddress(),
	}
}

func TestWhenCreatingAgentFirstStep(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return validation error", func(t *testing.T) {
			t.Parallel()

			request := generateValidCreateAgentFirstStepRequest()
			request.Name = ""
			suite := WhenCreatingAgentFirstStepBeforeEach(t)

			_, err := suite.sut.Execute(request)

			assert.NotNil(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given registry requester fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()

			suite.agentRegistryRequesterMock.EXPECT().
				CreateAgentFirstStep(gomock.Any()).
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

			suite := WhenCreatingAgentFirstStepBeforeEach(t)
			request := generateValidCreateAgentFirstStepRequest()
			expected := &responses.AgentResponse{Name: request.Name}

			suite.agentRegistryRequesterMock.EXPECT().
				CreateAgentFirstStep(gomock.Any()).
				DoAndReturn(func(registryRequest *interfaces.CreateRegistryAgentFirstStepRequest) (*responses.AgentResponse, error) {
					assert.Equal(t, request.AuthenticatedWallet, registryRequest.AuthenticatedWallet)
					assert.Equal(t, request.CollectionAddress, registryRequest.CollectionAddress)
					assert.Equal(t, request.Name, registryRequest.Name)
					assert.Equal(t, request.Role, registryRequest.Role)
					assert.Equal(t, request.Mission, registryRequest.Mission)
					assert.Equal(t, request.WalletAddress, registryRequest.WalletAddress)
					assert.Equal(t, request.CloneOfAgentId, registryRequest.CloneOfAgentId)
					assert.NotNil(t, registryRequest.Image)
					return expected, nil
				}).
				Times(1)

			response, err := suite.sut.Execute(request)
			assert.Nil(t, err)
			assert.Equal(t, expected, response)
		})
	})
}
