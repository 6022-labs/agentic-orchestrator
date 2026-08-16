// Package fixtures deploys the 6022 agentic contract stack onto a throwaway chain; see README.md.
package fixtures

import (
	"context"
	"embed"
	"fmt"
	"math/big"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

//go:embed artifacts/*.json
var artifactFiles embed.FS

// StackCollectionName must stay a valid ENS label — the collection contract rejects anything else.
const StackCollectionName = "agents"

// Mirror the contract's tuple layout; the generated bindings ship no bytecode, so mints go through raw ABI packing.
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
// through the creator, governed by the manager.
type AgentContractStack struct {
	Manager        common.Address
	Descriptor     common.Address
	Implementation common.Address
	Creator        common.Address
	Collection     common.Address
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

	collection, err := CreateCollection(ctx, client, anvilAccountSettings, manager, creator, StackCollectionName)
	if err != nil {
		return nil, err
	}

	return &AgentContractStack{
		Manager:        manager,
		Descriptor:     descriptor,
		Implementation: implementation,
		Creator:        creator,
		Collection:     collection,
		CollectionName: StackCollectionName,
	}, nil
}

// CreateCollection adds one more collection to an existing stack, so tests can
// assert on the manager's creation event and on multi-collection listings.
func CreateCollection(
	ctx context.Context,
	client *ethclient.Client,
	anvilAccountSettings container_builders.AnvilAccountSettings,
	manager common.Address,
	creator common.Address,
	name string,
) (common.Address, error) {
	deployer, err := contracts.NewDeployer(ctx, client, anvilAccountSettings)
	if err != nil {
		return common.Address{}, err
	}

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

// MintAgent mints straight through the collection contract; the anvil account is
// its moderator, so mint() needs no proposal. evmAddress must be unique per
// collection — the contract rejects a second agent claiming the same address.
func MintAgent(
	ctx context.Context,
	client *ethclient.Client,
	anvilAccountSettings container_builders.AnvilAccountSettings,
	collection common.Address,
	agentName string,
	evmAddress common.Address,
) error {
	deployer, boundCollection, err := bindCollection(ctx, client, anvilAccountSettings, collection)
	if err != nil {
		return err
	}

	return deployer.Transact(
		boundCollection,
		"mint",
		deployer.Address,
		agentName,
		[]agentAddress{{AddressType: "evm", Value: evmAddress.Hex()}},
		[]keyValue{{Key: "default", Value: "bafybeibsulc7sgpck3w7ur7eqa6fwaazcsswn3dg3rfl65zka3knajx43i"}},
		[]keyValue{{Key: "role", Value: "facilitator"}, {Key: "mission", Value: "integration testing"}},
		nullableUint256{IsSet: false, Value: big.NewInt(0)},
	)
}

// CreateMintProposal proposes a mint instead of minting outright, so proposal
// reads have something to find.
func CreateMintProposal(
	ctx context.Context,
	client *ethclient.Client,
	anvilAccountSettings container_builders.AnvilAccountSettings,
	collection common.Address,
	agentName string,
) error {
	deployer, boundCollection, err := bindCollection(ctx, client, anvilAccountSettings, collection)
	if err != nil {
		return err
	}

	return deployer.Transact(
		boundCollection,
		"createMintProposal",
		deployer.Address,
		agentName,
		[]agentAddress{{AddressType: "evm", Value: deployer.Address.Hex()}},
		[]keyValue{{Key: "default", Value: "bafybeibsulc7sgpck3w7ur7eqa6fwaazcsswn3dg3rfl65zka3knajx43i"}},
		[]keyValue{{Key: "role", Value: "facilitator"}, {Key: "mission", Value: "integration testing"}},
		nullableUint256{IsSet: false, Value: big.NewInt(0)},
	)
}

func bindCollection(
	ctx context.Context,
	client *ethclient.Client,
	anvilAccountSettings container_builders.AnvilAccountSettings,
	collection common.Address,
) (*contracts.Deployer, *bind.BoundContract, error) {
	deployer, err := contracts.NewDeployer(ctx, client, anvilAccountSettings)
	if err != nil {
		return nil, nil, err
	}
	collectionContract, err := load("AgentCollectionV1")
	if err != nil {
		return nil, nil, err
	}
	return deployer, deployer.Bind(collectionContract, collection), nil
}

func load(name string) (*contracts.CompiledContract, error) {
	raw, err := artifactFiles.ReadFile("artifacts/" + name + ".json")
	if err != nil {
		return nil, fmt.Errorf("no committed artifact for %s (see README.md): %w", name, err)
	}
	return contracts.Parse(raw)
}
