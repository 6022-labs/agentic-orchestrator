package use_cases_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/src/database_connector/use_cases/requests"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/common_tests/model_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/database_connector_mocks/repositories_mocks"
	"github.com/6022-labs/agentic-orchestrator/agent-knowledge/database-connector/tests/database_connector_mocks/services_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenCreatingDatabaseConnectionTestingSuite struct {
	sut *use_cases.CreateDatabaseConnection

	databaseCredentialStorageMock    *services_mocks.MockAgentDatabaseCredentialsStorageInterface
	databaseConnectionRepositoryMock *repositories_mocks.MockDatabaseConnectionRepositoryInterface
	databaseCredentialGeneratorMock  *services_mocks.MockAgentDatabaseCredentialsGeneratorInterface
}

func WhenCreatingDatabaseConnectionBeforeEach(t *testing.T) *WhenCreatingDatabaseConnectionTestingSuite {
	mockController := gomock.NewController(t)

	databaseCredentialStorageMock := services_mocks.NewMockAgentDatabaseCredentialsStorageInterface(mockController)
	databaseCredentialGeneratorMock := services_mocks.NewMockAgentDatabaseCredentialsGeneratorInterface(mockController)
	databaseConnectionRepositoryMock := repositories_mocks.NewMockDatabaseConnectionRepositoryInterface(mockController)

	sut := use_cases.NewCreateDatabaseConnection(
		databaseConnectionRepositoryMock,
		databaseCredentialStorageMock,
		databaseCredentialGeneratorMock,
	)

	return &WhenCreatingDatabaseConnectionTestingSuite{
		sut: sut,

		databaseCredentialStorageMock:    databaseCredentialStorageMock,
		databaseCredentialGeneratorMock:  databaseCredentialGeneratorMock,
		databaseConnectionRepositoryMock: databaseConnectionRepositoryMock,
	}
}

func generateValidRequest() *requests.CreateDatabaseConnectionRequest {
	return &requests.CreateDatabaseConnectionRequest{
		Host:    "localhost",
		Port:    5432,
		Driver:  "postgres",
		AgentID: uuid.New(),
	}
}

func TestWhenCreatingDatabaseConnection(t *testing.T) {
	t.Parallel()

	t.Run("Given request is invalid", func(t *testing.T) {
		t.Parallel()

		t.Run("Because host is empty", func(t *testing.T) {
			t.Parallel()

			sut := WhenCreatingDatabaseConnectionBeforeEach(t)

			request := generateValidRequest()
			request.Host = ""

			_, err := sut.sut.Create(request)

			assert.Error(t, err)
		})

		t.Run("Because driver is invalid", func(t *testing.T) {
			t.Parallel()

			sut := WhenCreatingDatabaseConnectionBeforeEach(t)

			request := generateValidRequest()
			request.Driver = "invalid"

			_, err := sut.sut.Create(request)

			assert.Error(t, err)
		})

		t.Run("Because agentId is empty", func(t *testing.T) {
			t.Parallel()

			sut := WhenCreatingDatabaseConnectionBeforeEach(t)

			request := generateValidRequest()
			request.AgentID = uuid.Nil

			_, err := sut.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurred while getting existing database connection", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			sut := WhenCreatingDatabaseConnectionBeforeEach(t)

			request := generateValidRequest()

			sut.databaseConnectionRepositoryMock.EXPECT().
				GetByAgentId(request.AgentID).
				Return(nil, assert.AnError)

			_, err := sut.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given database connection already exists", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenCreatingDatabaseConnectionBeforeEach(t)

			request := generateValidRequest()

			sut.databaseConnectionRepositoryMock.EXPECT().
				GetByAgentId(request.AgentID).
				Return(model_builders.NewDatabaseConnectionBuilder().Build(), nil)

			_, err := sut.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given error occurred while saving database connection", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			sut := WhenCreatingDatabaseConnectionBeforeEach(t)

			request := generateValidRequest()

			sut.databaseConnectionRepositoryMock.EXPECT().
				GetByAgentId(request.AgentID).
				Return(nil, nil)

			sut.databaseConnectionRepositoryMock.EXPECT().
				Save(gomock.Any()).
				Return(nil, assert.AnError)

			_, err := sut.sut.Create(request)

			assert.Error(t, err)
		})

	})

	t.Run("Given error occurred while storing credentials", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingDatabaseConnectionTestingSuite, request *requests.CreateDatabaseConnectionRequest) {
			suite.databaseConnectionRepositoryMock.EXPECT().
				GetByAgentId(request.AgentID).
				Return(nil, nil)

			databaseConnection := model_builders.NewDatabaseConnectionBuilder().Build()
			suite.databaseConnectionRepositoryMock.EXPECT().
				Save(gomock.Any()).
				Return(databaseConnection, nil)

			suite.databaseCredentialGeneratorMock.EXPECT().
				Generate(gomock.Any()).
				Return("generated-credentials")

			suite.databaseCredentialStorageMock.EXPECT().
				Store(gomock.Any(), "generated-credentials").
				Return(assert.AnError)
		}

		t.Run("Should return error", func(t *testing.T) {
			t.Parallel()

			request := generateValidRequest()
			sut := WhenCreatingDatabaseConnectionBeforeEach(t)
			initSuite(sut, request)

			_, err := sut.sut.Create(request)

			assert.Error(t, err)
		})
	})

	t.Run("Given database connection is created successfully", func(t *testing.T) {
		t.Parallel()

		initSuite := func(suite *WhenCreatingDatabaseConnectionTestingSuite, request *requests.CreateDatabaseConnectionRequest) {
			suite.databaseConnectionRepositoryMock.EXPECT().
				GetByAgentId(request.AgentID).
				Return(nil, nil)

			databaseConnection := model_builders.NewDatabaseConnectionBuilder().Build()
			suite.databaseConnectionRepositoryMock.EXPECT().
				Save(gomock.Any()).
				Return(databaseConnection, nil)

			suite.databaseCredentialGeneratorMock.EXPECT().
				Generate(gomock.Any()).
				Return("generated-credentials")

			suite.databaseCredentialStorageMock.EXPECT().
				Store(gomock.Any(), "generated-credentials").
				Return(nil)
		}

		t.Run("Should return database credentials", func(t *testing.T) {
			t.Parallel()

			request := generateValidRequest()
			sut := WhenCreatingDatabaseConnectionBeforeEach(t)
			initSuite(sut, request)

			response, err := sut.sut.Create(request)

			assert.NoError(t, err)
			assert.NotNil(t, response)
		})
	})
}
