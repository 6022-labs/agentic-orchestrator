package services_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/common_tests/container_builders"
	"github.com/6022-labs/agentic-orchestrator/agent-identity/tests/registry_synchronizer_blockchain_integration_tests/fixtures"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// chainFixture is a throwaway anvil chain with a deployed 6022 stack. The http
// and ws clients mirror the two named clients the host wires.
type chainFixture struct {
	t                    *testing.T
	anvilAccountSettings container_builders.AnvilAccountSettings
	stack                *fixtures.AgentContractStack
	httpClient           *ethclient.Client
	wsClient             *ethclient.Client
}

func newChainFixture(t *testing.T) *chainFixture {
	anvilAccountSettings := container_builders.CreateAnvilAccount()

	httpClient, err := ethclient.Dial(anvilAccountSettings.HttpUrl)
	if err != nil {
		t.Fatalf("failed to dial anvil over http: %v", err)
	}
	t.Cleanup(httpClient.Close)

	wsClient, err := ethclient.Dial(anvilAccountSettings.WsUrl)
	if err != nil {
		t.Fatalf("failed to dial anvil over websocket: %v", err)
	}
	t.Cleanup(wsClient.Close)

	stack, err := fixtures.DeployAgentStack(context.Background(), httpClient, anvilAccountSettings)
	if err != nil {
		t.Fatalf("failed to deploy contract stack: %v", err)
	}

	return &chainFixture{
		t:                    t,
		anvilAccountSettings: anvilAccountSettings,
		stack:                stack,
		httpClient:           httpClient,
		wsClient:             wsClient,
	}
}

func (f *chainFixture) agentCollectionsManagerSettings() *settings.AgentCollectionsManagerSettings {
	return &settings.AgentCollectionsManagerSettings{SmartContractAddress: f.stack.Manager}
}

func (f *chainFixture) rpcNodeSettings() *settings.RpcNodeSettings {
	return &settings.RpcNodeSettings{MaxAddressesPerLogQuery: 10}
}

func (f *chainFixture) blockchainRequester() *services.BlockchainRequester {
	return f.blockchainRequesterThrough(f.httpClient)
}

func (f *chainFixture) blockchainRequesterThrough(client *ethclient.Client) *services.BlockchainRequester {
	return services.NewBlockchainRequester(services.NewBlockchainRequesterParams{
		Client:          client,
		RpcNodeSettings: f.rpcNodeSettings(),
	})
}

func (f *chainFixture) agentCollectionRequester() *services.AgentCollectionRequester {
	return f.agentCollectionRequesterThrough(f.httpClient)
}

func (f *chainFixture) agentCollectionRequesterThrough(client *ethclient.Client) *services.AgentCollectionRequester {
	requester, err := services.NewAgentCollectionRequester(services.NewAgentCollectionRequesterParams{
		Client:          client,
		RpcNodeSettings: f.rpcNodeSettings(),
	})
	if err != nil {
		f.t.Fatalf("failed to build the agent collection requester: %v", err)
	}
	return requester
}

func (f *chainFixture) agentCollectionsManagerRequester() *services.AgentCollectionsManagerRequester {
	requester, err := services.NewAgentCollectionsManagerRequester(services.NewAgentCollectionsManagerRequesterParams{
		Client:                          f.httpClient,
		RpcNodeSettings:                 f.rpcNodeSettings(),
		AgentCollectionsManagerSettings: f.agentCollectionsManagerSettings(),
	})
	if err != nil {
		f.t.Fatalf("failed to build the agent collections manager requester: %v", err)
	}
	return requester
}

func (f *chainFixture) agentCollectionEventSubscriptionProvider() *services.AgentCollectionEventSubscriptionProvider {
	return services.NewAgentCollectionEventSubscriptionProvider(
		services.NewAgentCollectionEventSubscriptionProviderParams{
			Logger: zap.NewNop(),
			Client: f.wsClient,
		},
	)
}

func (f *chainFixture) managerSubscriptionProvider() *services.AgentCollectionsManagerEventSubscriptionProvider {
	return services.NewAgentCollectionsManagerEventSubscriptionProvider(
		services.NewAgentCollectionsManagerEventSubscriptionProviderParams{
			Logger:                          zap.NewNop(),
			Client:                          f.wsClient,
			AgentCollectionsManagerSettings: f.agentCollectionsManagerSettings(),
		},
	)
}

func (f *chainFixture) currentBlock() uint64 {
	blockNumber, err := f.httpClient.BlockNumber(context.Background())
	if err != nil {
		f.t.Fatalf("failed to read the head block: %v", err)
	}
	return blockNumber
}

// mintAgent gives every agent its own evm address; the collection rejects a
// second agent claiming an address it already knows.
func (f *chainFixture) mintAgent(agentName string) {
	f.mintAgentIn(f.stack.Collection, agentName)
}

func (f *chainFixture) mintAgentIn(collection common.Address, agentName string) {
	if err := fixtures.MintAgent(
		context.Background(),
		f.httpClient,
		f.anvilAccountSettings,
		collection,
		agentName,
		common_tests.RandomEthAddress(),
	); err != nil {
		f.t.Fatalf("failed to mint %s: %v", agentName, err)
	}
}

func (f *chainFixture) createCollection(name string) common.Address {
	collection, err := fixtures.CreateCollection(
		context.Background(),
		f.httpClient,
		f.anvilAccountSettings,
		f.stack.Manager,
		f.stack.Creator,
		name,
	)
	if err != nil {
		f.t.Fatalf("failed to create the collection: %v", err)
	}
	return collection
}

// flakyHttpClient fails the listed requests (1-based) and proxies the rest to anvil.
func (f *chainFixture) flakyHttpClient(
	failedRequests []int32,
	reply func(w http.ResponseWriter, requestId json.RawMessage),
) (*ethclient.Client, *atomic.Int32) {
	target, err := url.Parse(f.anvilAccountSettings.HttpUrl)
	if err != nil {
		f.t.Fatalf("failed to parse the anvil url: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !slices.Contains(failedRequests, requests.Add(1)) {
			proxy.ServeHTTP(w, r)
			return
		}
		body, _ := io.ReadAll(r.Body)
		var request struct {
			Id json.RawMessage `json:"id"`
		}
		_ = json.Unmarshal(body, &request)
		reply(w, request.Id)
	}))
	f.t.Cleanup(server.Close)

	client, err := ethclient.Dial(server.URL)
	if err != nil {
		f.t.Fatalf("failed to dial the flaky proxy: %v", err)
	}
	f.t.Cleanup(client.Close)

	return client, &requests
}

func (f *chainFixture) replyStatus(status int) func(w http.ResponseWriter, requestId json.RawMessage) {
	return func(w http.ResponseWriter, _ json.RawMessage) {
		http.Error(w, http.StatusText(status), status)
	}
}

func (f *chainFixture) replyJsonRpcError(code int) func(w http.ResponseWriter, requestId json.RawMessage) {
	return func(w http.ResponseWriter, requestId json.RawMessage) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.Copy(w, strings.NewReader(
			`{"jsonrpc":"2.0","id":`+string(
				requestId,
			)+`,"error":{"code":`+strconv.Itoa(
				code,
			)+`,"message":"limit exceeded"}}`,
		))
	}
}
