package services_test

import (
	"fmt"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	common_interfaces_mocks "github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_common_mocks/services_mocks/interfaces_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type WhenStoringAgentPrivateKeyTestingSuite struct {
	sut *services.AgentPrivateKeyStorage

	mockVaultSecretStorage *common_interfaces_mocks.MockVaultSecretStorageInterface
}

func WhenStoringAgentPrivateKeyBeforeEach(t *testing.T) *WhenStoringAgentPrivateKeyTestingSuite {
	mockController := gomock.NewController(t)

	mockVaultSecretStorage := common_interfaces_mocks.NewMockVaultSecretStorageInterface(mockController)

	sut := services.NewAgentPrivateKeyStorage(mockVaultSecretStorage)

	return &WhenStoringAgentPrivateKeyTestingSuite{
		sut:                    sut,
		mockVaultSecretStorage: mockVaultSecretStorage,
	}
}

func TestWhenStoringAgentPrivateKey(t *testing.T) {
	t.Parallel()

	t.Run("Given storing a private key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should store the private key at the correct vault path", func(t *testing.T) {
			t.Parallel()
			basics := WhenStoringAgentPrivateKeyBeforeEach(t)

			agentId := uuid.New()
			privateKey := "abcdef1234567890"
			expectedPath := fmt.Sprintf(models.AgentPrivateKeyVaultPathTemplate, agentId.String())

			basics.mockVaultSecretStorage.EXPECT().
				Store(expectedPath, privateKey).
				Return(nil).
				Times(1)

			err := basics.sut.Store(agentId, privateKey)

			assert.Nil(t, err)
		})

		t.Run("Should return an error when vault storage fails", func(t *testing.T) {
			t.Parallel()
			basics := WhenStoringAgentPrivateKeyBeforeEach(t)

			agentId := uuid.New()

			basics.mockVaultSecretStorage.EXPECT().
				Store(gomock.Any(), gomock.Any()).
				Return(assert.AnError).
				Times(1)

			err := basics.sut.Store(agentId, "some-key")

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given getting a private key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the private key from the vault", func(t *testing.T) {
			t.Parallel()
			basics := WhenStoringAgentPrivateKeyBeforeEach(t)

			agentId := uuid.New()
			expectedKey := "abcdef1234567890"
			expectedPath := fmt.Sprintf(models.AgentPrivateKeyVaultPathTemplate, agentId.String())

			basics.mockVaultSecretStorage.EXPECT().
				Read(expectedPath).
				Return(expectedKey, nil).
				Times(1)

			result, err := basics.sut.Get(agentId)

			assert.Nil(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, expectedKey, *result)
		})

		t.Run("Should return an error when vault read fails", func(t *testing.T) {
			t.Parallel()
			basics := WhenStoringAgentPrivateKeyBeforeEach(t)

			agentId := uuid.New()

			basics.mockVaultSecretStorage.EXPECT().
				Read(gomock.Any()).
				Return("", assert.AnError).
				Times(1)

			result, err := basics.sut.Get(agentId)

			assert.Nil(t, result)
			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})

	t.Run("Given deleting a private key", func(t *testing.T) {
		t.Parallel()

		t.Run("Should delete the private key at the correct vault path", func(t *testing.T) {
			t.Parallel()
			basics := WhenStoringAgentPrivateKeyBeforeEach(t)

			agentId := uuid.New()
			expectedPath := fmt.Sprintf(models.AgentPrivateKeyVaultPathTemplate, agentId.String())

			basics.mockVaultSecretStorage.EXPECT().
				Delete(expectedPath).
				Return(nil).
				Times(1)

			err := basics.sut.Delete(agentId)

			assert.Nil(t, err)
		})

		t.Run("Should return an error when vault delete fails", func(t *testing.T) {
			t.Parallel()
			basics := WhenStoringAgentPrivateKeyBeforeEach(t)

			agentId := uuid.New()

			basics.mockVaultSecretStorage.EXPECT().
				Delete(gomock.Any()).
				Return(assert.AnError).
				Times(1)

			err := basics.sut.Delete(agentId)

			assert.NotNil(t, err)
			assert.Equal(t, assert.AnError, err)
		})
	})
}
