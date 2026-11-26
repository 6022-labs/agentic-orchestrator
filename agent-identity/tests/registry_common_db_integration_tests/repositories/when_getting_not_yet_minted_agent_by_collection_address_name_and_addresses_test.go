package repositories_test

import (
	"context"

	"strings"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/models/enums"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common_db/repositories"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/model_builders"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesTestingSuite struct {
	dbName string

	agentAddressRepository *repositories.AgentAddressRepository

	sut *repositories.AgentRepository
}

func WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesBeforeEach(
	t *testing.T,
) *WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesTestingSuite {
	dbName, dbClient := newMigratedDatabase(t)

	return &WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesTestingSuite{
		dbName:                 dbName,
		agentAddressRepository: repositories.NewAgentAddressRepository(dbClient),
		sut:                    repositories.NewAgentRepository(dbClient),
	}
}

func WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesAfterEach(
	suite *WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesTestingSuite,
) {
	container_builders.GetPostgresContainerBuilder().DropDatabase(suite.dbName)
}

// saveAgentWithEvmAddress stores the agent and gives it the evm address.
func (suite *WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesTestingSuite) saveAgentWithEvmAddress(
	t *testing.T,
	agent *models.Agent,
	evmAddress string,
) *models.Agent {
	saved, err := suite.sut.Save(context.Background(), agent)
	require.NoError(t, err)
	require.NoError(t, suite.agentAddressRepository.ReplaceAllForAgent(context.Background(), saved.ID,
		map[enums.AddressType]string{enums.AddressTypeEvm: evmAddress},
	))
	return saved
}

func TestWhenGettingNotYetMintedAgentByCollectionAddressNameAndAddresses(t *testing.T) {
	t.Parallel()

	newProposal := func(collectionAddress string, name string) *models.Agent {
		return model_builders.NewAgentBuilder().
			WithName(name).
			WithCollectionAddress(collectionAddress).
			WithCollectionTokenId(nil).
			Build()
	}

	t.Run("Given no addresses", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesBeforeEach(t)
			defer WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesAfterEach(suite)

			agent, err := suite.sut.GetNotYetMintedByCollectionAddressNameAndAddresses(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
				"alice",
				map[enums.AddressType]string{},
			)

			assert.NoError(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given no agent holds the addresses", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesBeforeEach(t)
			defer WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesAfterEach(suite)

			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			suite.saveAgentWithEvmAddress(t,
				newProposal(collectionAddress, "alice"),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
			)

			agent, err := suite.sut.GetNotYetMintedByCollectionAddressNameAndAddresses(
				context.Background(),
				collectionAddress,
				"alice",
				map[enums.AddressType]string{
					enums.AddressTypeEvm: strings.ToLower(common_tests.RandomEthAddress().Hex()),
				},
			)

			assert.NoError(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given the agent holding the address is already minted", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesBeforeEach(t)
			defer WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesAfterEach(suite)

			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			evmAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			suite.saveAgentWithEvmAddress(t,
				model_builders.NewAgentBuilder().
					WithName("alice").
					WithCollectionAddress(collectionAddress).
					WithCollectionTokenId(lo.ToPtr(uint64(1))).
					Build(),
				evmAddress,
			)

			agent, err := suite.sut.GetNotYetMintedByCollectionAddressNameAndAddresses(
				context.Background(),
				collectionAddress,
				"alice",
				map[enums.AddressType]string{enums.AddressTypeEvm: evmAddress},
			)

			assert.NoError(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given the proposal holding the address belongs to another collection", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesBeforeEach(t)
			defer WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesAfterEach(suite)

			evmAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			suite.saveAgentWithEvmAddress(t,
				newProposal(strings.ToLower(common_tests.RandomEthAddress().Hex()), "alice"),
				evmAddress,
			)

			agent, err := suite.sut.GetNotYetMintedByCollectionAddressNameAndAddresses(
				context.Background(),
				strings.ToLower(common_tests.RandomEthAddress().Hex()),
				"alice",
				map[enums.AddressType]string{enums.AddressTypeEvm: evmAddress},
			)

			assert.NoError(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given the proposal holding the address has another name", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return nil", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesBeforeEach(t)
			defer WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesAfterEach(suite)

			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			evmAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			suite.saveAgentWithEvmAddress(t, newProposal(collectionAddress, "bob"), evmAddress)

			agent, err := suite.sut.GetNotYetMintedByCollectionAddressNameAndAddresses(
				context.Background(),
				collectionAddress,
				"alice",
				map[enums.AddressType]string{enums.AddressTypeEvm: evmAddress},
			)

			assert.NoError(t, err)
			assert.Nil(t, agent)
		})
	})

	t.Run("Given agents of other collections also hold the address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return the collection's proposal", func(t *testing.T) {
			t.Parallel()

			suite := WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesBeforeEach(t)
			defer WhenGettingNotYetMintedAgentByCollectionAddressNameAndAddressesAfterEach(suite)

			collectionAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			evmAddress := strings.ToLower(common_tests.RandomEthAddress().Hex())
			suite.saveAgentWithEvmAddress(t,
				model_builders.NewAgentBuilder().
					WithName("alice").
					WithCollectionTokenId(lo.ToPtr(uint64(1))).
					Build(),
				evmAddress,
			)
			proposal := suite.saveAgentWithEvmAddress(t, newProposal(collectionAddress, "alice"), evmAddress)

			agent, err := suite.sut.GetNotYetMintedByCollectionAddressNameAndAddresses(
				context.Background(),
				collectionAddress,
				"alice",
				map[enums.AddressType]string{
					enums.AddressTypeBtc: "bc1qexampleonly",
					enums.AddressTypeEvm: evmAddress,
				},
			)

			assert.NoError(t, err)
			if assert.NotNil(t, agent) {
				assert.Equal(t, proposal.ID, agent.ID)
			}
		})
	})
}
