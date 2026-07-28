// Package fixtures deploys the ENS resolution contracts onto a throwaway chain; see README.md.
package fixtures

import (
	"context"
	"embed"
	"fmt"
	"strings"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/common_tests/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//go:embed artifacts/*.json
var artifactFiles embed.FS

// EnsContractStack is the real ENS registry with a name pointed at a real
// AgentEnsOffchainResolver, which reverts with ERC-3668 OffchainLookup.
type EnsContractStack struct {
	Registry common.Address
	Resolver common.Address
}

// DeployEnsStack deploys the registry and an offchain resolver trusting signer
// for its proofs, then claims registeredDomain and points it at the resolver.
func DeployEnsStack(
	ctx context.Context,
	client *ethclient.Client,
	anvilAccountSettings container_builders.AnvilAccountSettings,
	gatewayUrls []string,
	signer common.Address,
	registeredDomain string,
) (*EnsContractStack, error) {
	deployer, err := contracts.NewDeployer(ctx, client, anvilAccountSettings)
	if err != nil {
		return nil, err
	}

	registryContract, err := load("EnsRegistry")
	if err != nil {
		return nil, err
	}
	resolverContract, err := load("AgentEnsOffchainResolver")
	if err != nil {
		return nil, err
	}

	// The registry's constructor makes the deployer owner of the root node.
	registry, boundRegistry, err := deployer.Deploy("EnsRegistry", registryContract)
	if err != nil {
		return nil, err
	}
	resolver, _, err := deployer.Deploy("AgentEnsOffchainResolver", resolverContract, gatewayUrls, signer)
	if err != nil {
		return nil, err
	}

	if err := claimDomain(deployer, boundRegistry, registeredDomain); err != nil {
		return nil, err
	}
	if err := deployer.Transact(boundRegistry, "setResolver", Namehash(registeredDomain), resolver); err != nil {
		return nil, err
	}

	return &EnsContractStack{Registry: registry, Resolver: resolver}, nil
}

// SetResolver claims domain if needed and repoints it, so tests can cover a
// resolver that is not an offchain resolver at all.
func SetResolver(
	ctx context.Context,
	client *ethclient.Client,
	anvilAccountSettings container_builders.AnvilAccountSettings,
	registry common.Address,
	domain string,
	resolver common.Address,
) error {
	deployer, err := contracts.NewDeployer(ctx, client, anvilAccountSettings)
	if err != nil {
		return err
	}
	registryContract, err := load("EnsRegistry")
	if err != nil {
		return err
	}

	boundRegistry := deployer.Bind(registryContract, registry)
	if err := claimDomain(deployer, boundRegistry, domain); err != nil {
		return err
	}
	return deployer.Transact(boundRegistry, "setResolver", Namehash(domain), resolver)
}

// Namehash is computed here rather than reused from the resolver under test, so
// the fixture and the production implementation cross-check each other.
func Namehash(domain string) [32]byte {
	node := [32]byte{}
	if domain == "" {
		return node
	}
	labels := strings.Split(domain, ".")
	for i := len(labels) - 1; i >= 0; i-- {
		labelHash := crypto.Keccak256([]byte(labels[i]))
		node = [32]byte(crypto.Keccak256(node[:], labelHash))
	}
	return node
}

func load(name string) (*contracts.CompiledContract, error) {
	raw, err := artifactFiles.ReadFile("artifacts/" + name + ".json")
	if err != nil {
		return nil, fmt.Errorf("no committed artifact for %s (see README.md): %w", name, err)
	}
	return contracts.Parse(raw)
}

// claimDomain walks down from the root node creating each label, which the
// registry only allows the parent's owner to do.
func claimDomain(deployer *contracts.Deployer, boundRegistry *bind.BoundContract, domain string) error {
	labels := strings.Split(domain, ".")
	parent := [32]byte{}
	for i := len(labels) - 1; i >= 0; i-- {
		labelHash := [32]byte(crypto.Keccak256([]byte(labels[i])))
		if err := deployer.Transact(boundRegistry, "setSubnodeOwner", parent, labelHash, deployer.Address); err != nil {
			return fmt.Errorf("failed to claim %s: %w", strings.Join(labels[i:], "."), err)
		}
		parent = [32]byte(crypto.Keccak256(parent[:], labelHash[:]))
	}
	return nil
}
