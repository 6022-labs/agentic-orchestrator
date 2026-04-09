package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/stretchr/testify/assert"
)

type WhenBuildingDomainSeparatorTestingSuite struct {
	sut *services.EIP712SignatureHelper

	blockchainSettings      *settings.BlockchainSettings
	agentCollectionSettings *settings.AgentCollectionSettings
}

func WhenBuildingDomainSeparatorBeforeEach(t *testing.T) *WhenBuildingDomainSeparatorTestingSuite {
	blockchainSettings := &settings.BlockchainSettings{ChainId: 1}
	agentCollectionSettings := &settings.AgentCollectionSettings{SigningDomain: "AgentCollection"}

	sut := services.NewEIP712SignatureHelper(blockchainSettings, agentCollectionSettings)

	return &WhenBuildingDomainSeparatorTestingSuite{
		sut:                     sut,
		blockchainSettings:      blockchainSettings,
		agentCollectionSettings: agentCollectionSettings,
	}
}

func TestWhenBuildingDomainSeparator(t *testing.T) {
	t.Parallel()

	t.Run("Given a valid collection address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a non-nil domain separator hash", func(t *testing.T) {
			t.Parallel()
			basics := WhenBuildingDomainSeparatorBeforeEach(t)

			collectionAddress := common_tests.RandomEthAddress()

			result, err := basics.sut.BuildDomainSeparator(collectionAddress)

			assert.Nil(t, err)
			assert.NotNil(t, result)
		})

		t.Run("Should return deterministic results for the same input", func(t *testing.T) {
			t.Parallel()
			basics := WhenBuildingDomainSeparatorBeforeEach(t)

			collectionAddress := common_tests.RandomEthAddress()

			result1, err1 := basics.sut.BuildDomainSeparator(collectionAddress)
			result2, err2 := basics.sut.BuildDomainSeparator(collectionAddress)

			assert.Nil(t, err1)
			assert.Nil(t, err2)
			assert.Equal(t, result1, result2)
		})

		t.Run("Should return different hashes for different collection addresses", func(t *testing.T) {
			t.Parallel()
			basics := WhenBuildingDomainSeparatorBeforeEach(t)

			address1 := common_tests.RandomEthAddress()
			address2 := common_tests.RandomEthAddress()

			result1, err1 := basics.sut.BuildDomainSeparator(address1)
			result2, err2 := basics.sut.BuildDomainSeparator(address2)

			assert.Nil(t, err1)
			assert.Nil(t, err2)
			assert.NotEqual(t, result1, result2)
		})
	})
}
