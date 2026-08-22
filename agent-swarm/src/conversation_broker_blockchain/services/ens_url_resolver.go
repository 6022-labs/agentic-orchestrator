package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	contract_abi "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/helpers"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/settings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"go.uber.org/zap"
)

const (
	urlRecordKey = "url"

	// chainIdLabelIndex is where the chain lives in <agent>.<collection>.<chainId>.6022.eth.
	chainIdLabelIndex = 2
)

// EnsUrlResolver reads the agent's `url` text record straight from AgentEnsRegistry on the
// chain that owns it, skipping the L1 ENS resolver and the CCIP gateway it routes through.
type EnsUrlResolver struct {
	logger                   *zap.Logger
	blockchainSettings       *settings.BlockchainSettings
	rpcNodeSettings          *settings.RpcNodeSettings
	agentEnsRegistrySettings *settings.AgentEnsRegistrySettings
	agentEnsRegistry         *contract_abi.AgentEnsRegistryCaller
}

func NewEnsUrlResolver(
	logger *zap.Logger,
	ethClient clients.EthClientInterface,
	blockchainSettings *settings.BlockchainSettings,
	rpcNodeSettings *settings.RpcNodeSettings,
	agentEnsRegistrySettings *settings.AgentEnsRegistrySettings,
) *EnsUrlResolver {
	// A bind failure means broken generated code, so it panics at wiring time.
	agentEnsRegistry, err := contract_abi.NewAgentEnsRegistryCaller(
		agentEnsRegistrySettings.SmartContractAddress,
		ethClient,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to bind agent ens registry caller: %v", err))
	}

	return &EnsUrlResolver{
		logger:                   logger,
		blockchainSettings:       blockchainSettings,
		rpcNodeSettings:          rpcNodeSettings,
		agentEnsRegistrySettings: agentEnsRegistrySettings,
		agentEnsRegistry:         agentEnsRegistry,
	}
}

func (r *EnsUrlResolver) ResolveUrl(ctx context.Context, ensDomain string) (string, error) {
	if err := r.verifyChain(ensDomain); err != nil {
		return "", err
	}

	timeout := time.Duration(r.rpcNodeSettings.RequestTimeoutSeconds) * time.Second
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	node := helpers.Namehash(ensDomain)

	agentUrl, err := r.agentEnsRegistry.Text(&bind.CallOpts{Context: callCtx}, node, urlRecordKey)
	if err != nil {
		return "", fmt.Errorf("text record lookup failed for %s: %w", ensDomain, err)
	}

	agentUrl = strings.TrimRight(strings.TrimSpace(agentUrl), "/")
	if len(agentUrl) == 0 {
		r.logger.Warn(
			"agent ens registry holds no url record",
			zap.String("ensDomain", ensDomain),
			zap.String("node", fmt.Sprintf("%#x", node)),
			zap.String("smartContractAddress", r.agentEnsRegistrySettings.SmartContractAddress.Hex()),
		)
	}

	return agentUrl, nil
}

// verifyChain stops a name belonging to another chain being read against this registry,
// which would silently answer with an empty record.
func (r *EnsUrlResolver) verifyChain(ensDomain string) error {
	labels := strings.Split(ensDomain, ".")
	if len(labels) <= chainIdLabelIndex {
		return fmt.Errorf("%s has no chain id label", ensDomain)
	}

	chainId, err := strconv.ParseInt(labels[chainIdLabelIndex], 10, 64)
	if err != nil {
		return fmt.Errorf("%s has a non-numeric chain id label %q", ensDomain, labels[chainIdLabelIndex])
	}
	if chainId != r.blockchainSettings.ChainId {
		return fmt.Errorf(
			"%s belongs to chain %d, this broker reads chain %d",
			ensDomain, chainId, r.blockchainSettings.ChainId,
		)
	}

	return nil
}
