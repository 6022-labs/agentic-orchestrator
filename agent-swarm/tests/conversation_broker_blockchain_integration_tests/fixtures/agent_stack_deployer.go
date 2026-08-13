// Package fixtures deploys the 6022 agentic contract stack onto a throwaway chain; see README.md.
package fixtures

import (
	"context"
	"embed"
	"fmt"
	"math/big"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

//go:embed artifacts/*.json
var artifactFiles embed.FS

// EnsBaseDomain is the zone agent names hang under.
const EnsBaseDomain = "6022.eth"

// ensParentNode is namehash(EnsBaseDomain), pinned rather than computed: deriving it with the
// resolver's own helper would cancel out any name-independent bug in that helper.
var ensParentNode = common.HexToHash("0x9ec94752dcac587079728bde5f535457f7aa0bbe0b088555bca6d40cfcc78197")

// StackCollectionName must stay a valid ENS label — the collection contract rejects anything else.
const StackCollectionName = "agents"

// Mirror the contract's tuple layout; the generated binding ships no bytecode, so mints go through raw ABI packing.
type agentAddress struct {
	AddressType string `abi:"addressType"`
	Value       string `abi:"value"`
}

type keyValue struct {
	Key   string `abi:"key"`
	Value string `abi:"value"`
}

type nullableUint256 struct {
	IsSet bool     `abi:"isSet"`
	Value *big.Int `abi:"value"`
}

// AgentContractStack is a freshly deployed 6022 stack: one collection created
// through the creator, governed by the manager, plus the ENS registry the broker reads.
type AgentContractStack struct {
	Collection     common.Address
	EnsRegistry    common.Address
	CollectionName string
}

// DeployAgentStack mirrors the smart-contracts fixture deployCollectionWithCreator;
// the anvil account ends up admin and moderator of the created collection.
func DeployAgentStack(
	ctx context.Context,
	client *ethclient.Client,
	anvilAccountSettings container_builders.AnvilAccountSettings,
) (*AgentContractStack, error) {
	deployer, err := contracts.NewDeployer(ctx, client, anvilAccountSettings)
	if err != nil {
		return nil, err
	}

	managerContract, err := load("AgentCollectionsManager")
	if err != nil {
		return nil, err
	}
	manager, boundManager, err := deployer.Deploy("AgentCollectionsManager", managerContract)
	if err != nil {
		return nil, err
	}

	descriptorContract, err := load("AgentCollectionsDescriptor")
	if err != nil {
		return nil, err
	}
	descriptor, _, err := deployer.Deploy("AgentCollectionsDescriptor", descriptorContract, "ipfs://")
	if err != nil {
		return nil, err
	}

	implementationContract, err := load("AgentCollectionV1")
	if err != nil {
		return nil, err
	}
	implementation, _, err := deployer.Deploy("AgentCollectionV1", implementationContract)
	if err != nil {
		return nil, err
	}

	creatorContract, err := load("AgentCollectionCreatorV1")
	if err != nil {
		return nil, err
	}
	creator, _, err := deployer.Deploy("AgentCollectionCreatorV1", creatorContract, manager, descriptor, implementation)
	if err != nil {
		return nil, err
	}

	var creatorRole [32]byte
	if err := deployer.Call(boundManager, &creatorRole, "CREATOR_ROLE"); err != nil {
		return nil, err
	}
	if err := deployer.Transact(boundManager, "grantRole", creatorRole, creator); err != nil {
		return nil, err
	}

	collection, err := createCollection(deployer, manager, creator, StackCollectionName)
	if err != nil {
		return nil, err
	}

	ensRegistryContract, err := load("AgentEnsRegistry")
	if err != nil {
		return nil, err
	}
	ensRegistry, _, err := deployer.Deploy("AgentEnsRegistry", ensRegistryContract, manager, [32]byte(ensParentNode))
	if err != nil {
		return nil, err
	}

	return &AgentContractStack{
		Collection:     collection,
		EnsRegistry:    ensRegistry,
		CollectionName: StackCollectionName,
	}, nil
}

// PublishAgent mints an agent and publishes its url text record; an empty url mints without provisioning.
func PublishAgent(
	ctx context.Context,
	client *ethclient.Client,
	anvilAccountSettings container_builders.AnvilAccountSettings,
	stack *AgentContractStack,
	agentName string,
	agentUrl string,
) error {
	deployer, err := contracts.NewDeployer(ctx, client, anvilAccountSettings)
	if err != nil {
		return err
	}

	collectionContract, err := load("AgentCollectionV1")
	if err != nil {
		return err
	}
	boundCollection := deployer.Bind(collectionContract, stack.Collection)

	// The anvil account is the collection's moderator, so mint() needs no proposal.
	if err := deployer.Transact(
		boundCollection,
		"mint",
		deployer.Address,
		agentName,
		[]agentAddress{{AddressType: "evm", Value: deployer.Address.Hex()}},
		[]keyValue{{Key: "default", Value: "bafybeibsulc7sgpck3w7ur7eqa6fwaazcsswn3dg3rfl65zka3knajx43i"}},
		[]keyValue{{Key: "role", Value: "facilitator"}},
		nullableUint256{IsSet: false, Value: big.NewInt(0)},
	); err != nil {
		return err
	}

	if len(agentUrl) == 0 {
		return nil
	}

	var tokenId *big.Int
	if err := deployer.Call(boundCollection, &tokenId, "nameToTokenId", agentName); err != nil {
		return err
	}

	ensRegistryContract, err := load("AgentEnsRegistry")
	if err != nil {
		return err
	}

	// initAgentTexts provisions the subnode and rejects a batch without the url key.
	return deployer.Transact(
		deployer.Bind(ensRegistryContract, stack.EnsRegistry),
		"initAgentTexts",
		stack.Collection,
		tokenId,
		[]keyValue{{Key: "url", Value: agentUrl}},
	)
}

// EnsDomainFor builds the name the broker resolves: <agent>.<collection>.<chainId>.<base>.
func EnsDomainFor(stack *AgentContractStack, chainId, agentName string) string {
	return strings.Join([]string{agentName, stack.CollectionName, chainId, EnsBaseDomain}, ".")
}

func createCollection(
	deployer *contracts.Deployer,
	manager common.Address,
	creator common.Address,
	name string,
) (common.Address, error) {
	managerContract, err := load("AgentCollectionsManager")
	if err != nil {
		return common.Address{}, err
	}
	creatorContract, err := load("AgentCollectionCreatorV1")
	if err != nil {
		return common.Address{}, err
	}

	boundManager := deployer.Bind(managerContract, manager)
	boundCreator := deployer.Bind(creatorContract, creator)

	if err := deployer.Transact(
		boundCreator,
		"createCollection",
		name,
		strings.ToUpper(name),
		"integration test collection",
		deployer.Address,
		deployer.Address,
	); err != nil {
		return common.Address{}, err
	}

	// Collections are 1-indexed (nextCollectionId starts at 1), so the one just created is the last id.
	var nextCollectionId *big.Int
	if err := deployer.Call(boundManager, &nextCollectionId, "nextCollectionId"); err != nil {
		return common.Address{}, err
	}
	var collection common.Address
	if err := deployer.Call(
		boundManager,
		&collection,
		"collections",
		new(big.Int).Sub(nextCollectionId, big.NewInt(1)),
	); err != nil {
		return common.Address{}, err
	}
	return collection, nil
}

func load(name string) (*contracts.CompiledContract, error) {
	raw, err := artifactFiles.ReadFile("artifacts/" + name + ".json")
	if err != nil {
		return nil, fmt.Errorf("no committed artifact for %s (see README.md): %w", name, err)
	}
	return contracts.Parse(raw)
}
