package services_test

import (
	"testing"

	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/services"
	"github.com/6022-labs/agentic-ai-orchestrator-mono-repo/gateway/src/orchestrator/settings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var (
	COLLECTION_ADDRESS = common.HexToAddress("0x875A3A6cfBc44593bA37f9eF811D70640672DC63")
)

type WhenBuildingEIP712DomainSeparatorTestSuite struct {
	sut *services.EIP712SignatureHelper

	blockchainSettings      *settings.BlockchainSettings
	agentCollectionSettings *settings.AgentCollectionSettings
}

func WhenBuildingEIP712DomainSeparatorBeforeEach() WhenBuildingEIP712DomainSeparatorTestSuite {
	blockchainSettings := &settings.BlockchainSettings{
		ChainId: 137,
	}

	agentCollectionSettings := &settings.AgentCollectionSettings{
		SigningDomain: "AgenticAIAgentCollection",
	}

	sut := services.NewEIP712SignatureHelper(
		blockchainSettings,
		agentCollectionSettings,
	)

	return WhenBuildingEIP712DomainSeparatorTestSuite{
		sut: sut,

		blockchainSettings:      blockchainSettings,
		agentCollectionSettings: agentCollectionSettings,
	}
}

func TestWhenBuildingEIP712DomainSeparator(t *testing.T) {
	t.Parallel()

	t.Run("Should build domain separator", func(t *testing.T) {
		t.Parallel()

		suite := WhenBuildingEIP712DomainSeparatorBeforeEach()

		domainSeparator, err := suite.sut.BuildDomainSeparator(COLLECTION_ADDRESS)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		expectedDomainSeparator := common.HexToHash("0x5389353371ac10565913ca82fd7b51b430994dd0235185084d9f469db4e98069")

		assert.Equal(t, expectedDomainSeparator, *domainSeparator)
	})
}
