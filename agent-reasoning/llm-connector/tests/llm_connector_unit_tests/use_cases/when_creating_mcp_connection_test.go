package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/errors"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/models"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/use_cases"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/src/llm_connector/use_cases/requests"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/llm_connector_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/agent-reasoning/llm-connector/tests/llm_connector_mocks/services_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCreatingMCPConnectionTestingSuite struct {
	sut *use_cases.CreateMCPConnection

	mcpConnectionRepositoryMock             *repositories_mocks.MockMCPConnectionRepositoryInterface
	mcpConnectionHeaderAuthValueStorageMock *services_mocks.MockMCPConnectionHeaderAuthValueStorageInterface
}

func WhenCreatingMCPConnectionBeforeEach(t *testing.T) *WhenCreatingMCPConnectionTestingSuite {
	mockController := gomock.NewController(t)

	mcpConnectionRepositoryMock := repositories_mocks.NewMockMCPConnectionRepositoryInterface(mockController)
	mcpConnectionHeaderAuthValueStorageMock := services_mocks.NewMockMCPConnectionHeaderAuthValueStorageInterface(mockController)

	sut := use_cases.NewCreateMCPConnection(mcpConnectionRepositoryMock, mcpConnectionHeaderAuthValueStorageMock)

	return &WhenCreatingMCPConnectionTestingSuite{
		sut: sut,

		mcpConnectionRepositoryMock:             mcpConnectionRepositoryMock,
		mcpConnectionHeaderAuthValueStorageMock: mcpConnectionHeaderAuthValueStorageMock,
	}
}

func generateValidCreateMCPConnectionRequest() *requests.CreateMCPConnectionRequest {
	return &requests.CreateMCPConnectionRequest{
		AgentId:    uuid.New(),
		Name:       "test-conn",
		Url:        "https://mcp.example.com",
		HeaderAuth: map[string]string{"X-Api-Key": "secret"},
	}
}

func TestWhenCreatingMCPConnection(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			req.AgentId = uuid.Nil

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because name is lower than 3 characters", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			req.Name = "ab"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because name does not start with a letter", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			req.Name = "123Invalids"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because name contains underscore", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			req.Name = "test_name"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because name contains invalid characters", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			req.Name = "inv@lid"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because url is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			req.Url = ""

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because url is not a valid url hostname", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			req.Url = "https://invalid-url"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because url is not a valid url ipv4", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			req.Url = "https://256.256.256.256"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because url is not a valid url ipv6", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			req.Url = "https://[::1]"

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because headerAuth key is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			req.HeaderAuth = map[string]string{" ": "value"}

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})

		t.Run("Because headerAuth value is empty", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			req.HeaderAuth = map[string]string{"key": "  "}

			err := suite.sut.Create(req)

			assert.Error(t, err)
			assert.IsType(t, &errors.ValidationError{}, err)
		})
	})

	t.Run("Given error occurs while checking if MCPConnection already exists by name in repository", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingMCPConnectionTestingSuite, req *requests.CreateMCPConnectionRequest) {
			suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndName(req.AgentId, req.Name).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			initSuite(suite, req)

			err := suite.sut.Create(req)

			assert.Error(t, err)
		})
	})

	t.Run("Given MCPConnection name already exists for this agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingMCPConnectionTestingSuite, req *requests.CreateMCPConnectionRequest) {
			exists := true
			suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndName(req.AgentId, req.Name).Return(&exists, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			initSuite(suite, req)

			err := suite.sut.Create(req)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while checking if MCPConnection already exists by url in repository", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingMCPConnectionTestingSuite, req *requests.CreateMCPConnectionRequest) {
			notExists := false
			suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndName(req.AgentId, req.Name).Return(&notExists, nil)
			suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndUrl(req.AgentId, req.Url).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			initSuite(suite, req)

			err := suite.sut.Create(req)

			assert.Error(t, err)
		})
	})

	t.Run("Given MCPConnection url already exists for this agent", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingMCPConnectionTestingSuite, req *requests.CreateMCPConnectionRequest) {
			notExists := false
			exists := true
			suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndName(req.AgentId, req.Name).Return(&notExists, nil)
			suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndUrl(req.AgentId, req.Url).Return(&exists, nil)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			initSuite(suite, req)

			err := suite.sut.Create(req)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while storing header auth values", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingMCPConnectionTestingSuite, req *requests.CreateMCPConnectionRequest) {
			notExists := false
			suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndName(req.AgentId, req.Name).Return(&notExists, nil)
			suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndUrl(req.AgentId, req.Url).Return(&notExists, nil)
			suite.mcpConnectionHeaderAuthValueStorageMock.EXPECT().Store(gomock.Any(), req.HeaderAuth).Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			initSuite(suite, req)

			err := suite.sut.Create(req)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurs while saving MCPConnection", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingMCPConnectionTestingSuite, req *requests.CreateMCPConnectionRequest) {
			notExists := false
			suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndName(req.AgentId, req.Name).Return(&notExists, nil)
			suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndUrl(req.AgentId, req.Url).Return(&notExists, nil)
			suite.mcpConnectionHeaderAuthValueStorageMock.EXPECT().Store(gomock.Any(), req.HeaderAuth).Return(nil)
			suite.mcpConnectionRepositoryMock.EXPECT().Save(gomock.Any()).Return(nil, assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCreatingMCPConnectionBeforeEach(t)
			req := generateValidCreateMCPConnectionRequest()
			initSuite(suite, req)

			err := suite.sut.Create(req)

			assert.Error(t, err)
		})
	})

	t.Run("Given MCPConnection is created successfully", func(t *testing.T) {
		t.Parallel()

		t.Run("With auth headers", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenCreatingMCPConnectionTestingSuite, req *requests.CreateMCPConnectionRequest) {
				notExists := false
				suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndName(req.AgentId, req.Name).Return(&notExists, nil)
				suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndUrl(req.AgentId, req.Url).Return(&notExists, nil)
				suite.mcpConnectionHeaderAuthValueStorageMock.EXPECT().Store(gomock.Any(), req.HeaderAuth).Return(nil)
				suite.mcpConnectionRepositoryMock.EXPECT().Save(gomock.Any()).Return(&models.MCPConnection{}, nil)
			}

			t.Run("Should return nil", func(t *testing.T) {
				t.Parallel()

				suite := WhenCreatingMCPConnectionBeforeEach(t)
				req := generateValidCreateMCPConnectionRequest()
				initSuite(suite, req)

				err := suite.sut.Create(req)

				assert.NoError(t, err)
			})
		})

		t.Run("Without auth headers", func(t *testing.T) {
			t.Parallel()

			initSuite := func(suite *WhenCreatingMCPConnectionTestingSuite, req *requests.CreateMCPConnectionRequest) {
				notExists := false
				suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndName(req.AgentId, req.Name).Return(&notExists, nil)
				suite.mcpConnectionRepositoryMock.EXPECT().ExistsByAgentIdAndUrl(req.AgentId, req.Url).Return(&notExists, nil)
				suite.mcpConnectionRepositoryMock.EXPECT().Save(gomock.Any()).Return(&models.MCPConnection{}, nil)
			}

			t.Run("Should return nil", func(t *testing.T) {
				t.Parallel()

				suite := WhenCreatingMCPConnectionBeforeEach(t)
				req := generateValidCreateMCPConnectionRequest()
				req.HeaderAuth = nil // No auth headers
				initSuite(suite, req)

				err := suite.sut.Create(req)

				assert.NoError(t, err)
			})
		})

	})
}
