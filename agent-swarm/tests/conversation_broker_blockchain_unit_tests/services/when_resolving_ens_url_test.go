package services_test

import (
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/services"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/settings"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/tests/conversation_broker_blockchain_mocks/clients_mocks"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

const resolverAddressForTest = "0xa3f6011779B55BB155BDAcd25Dd5Bcb9f5597646"

// fakeDataError mimics an rpc.DataError revert carrying ABI-encoded data.
type fakeDataError struct {
	data string
}

func (e *fakeDataError) Error() string          { return "execution reverted" }
func (e *fakeDataError) ErrorData() interface{} { return e.data }

type WhenResolvingEnsUrlTestingSuite struct {
	sut *services.EnsUrlResolver

	mockEthCaller *clients_mocks.MockEthCallerInterface
}

func WhenResolvingEnsUrlBeforeEach(t *testing.T) *WhenResolvingEnsUrlTestingSuite {
	mockController := gomock.NewController(t)

	mockEthCaller := clients_mocks.NewMockEthCallerInterface(mockController)

	sut := services.NewEnsUrlResolver(services.NewEnsUrlResolverParams{
		Logger:     zap.NewNop(),
		EthCaller:  mockEthCaller,
		HttpClient: http.DefaultClient,
		EnsResolverSettings: &settings.EnsResolverSettings{
			RpcUrl:                       "http://localhost:0",
			RegistryAddress:              "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e",
			GatewayRequestTimeoutSeconds: 5,
		},
	})

	return &WhenResolvingEnsUrlTestingSuite{
		sut: sut,

		mockEthCaller: mockEthCaller,
	}
}

func packOffchainLookupRevertHex(t *testing.T, gatewayUrl string) string {
	addressType, _ := abi.NewType("address", "", nil)
	stringsType, _ := abi.NewType("string[]", "", nil)
	bytesType, _ := abi.NewType("bytes", "", nil)
	bytes4Type, _ := abi.NewType("bytes4", "", nil)

	args := abi.Arguments{{Type: addressType}, {Type: stringsType}, {Type: bytesType}, {Type: bytes4Type}, {Type: bytesType}}
	packed, err := args.Pack(common.HexToAddress(resolverAddressForTest), []string{gatewayUrl}, []byte{1}, [4]byte{1, 2, 3, 4}, []byte{2})
	assert.NoError(t, err)

	selector := crypto.Keccak256([]byte("OffchainLookup(address,string[],bytes,bytes4,bytes)"))[:4]
	return "0x" + hex.EncodeToString(append(selector, packed...))
}

func packedResolverAddressReturn(t *testing.T) []byte {
	addressType, _ := abi.NewType("address", "", nil)
	returnData, err := abi.Arguments{{Type: addressType}}.Pack(common.HexToAddress(resolverAddressForTest))
	assert.NoError(t, err)

	return returnData
}

func registryDiscoveryReturningResolver(suite *WhenResolvingEnsUrlTestingSuite, t *testing.T) *gomock.Call {
	return suite.mockEthCaller.EXPECT().
		CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(packedResolverAddressReturn(t), nil).
		Times(1)
}

func packedTextReturn(t *testing.T, url string) []byte {
	stringType, _ := abi.NewType("string", "", nil)
	innerResult, err := abi.Arguments{{Type: stringType}}.Pack(url)
	assert.NoError(t, err)

	bytesType, _ := abi.NewType("bytes", "", nil)
	returnData, err := abi.Arguments{{Type: bytesType}}.Pack(innerResult)
	assert.NoError(t, err)

	return returnData
}

func TestWhenResolvingEnsUrl(t *testing.T) {
	t.Parallel()

	t.Run("Given the resolver reverts with OffchainLookup and the gateway answers", func(t *testing.T) {
		t.Parallel()

		t.Run("Should follow the gateway and return the verified url", func(t *testing.T) {
			t.Parallel()

			gatewayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"data":"0x1234"}`))
			}))
			defer gatewayServer.Close()

			suite := WhenResolvingEnsUrlBeforeEach(t)

			revertHex := packOffchainLookupRevertHex(t, gatewayServer.URL+"/{sender}/{data}.json")

			gomock.InOrder(
				registryDiscoveryReturningResolver(suite, t),
				suite.mockEthCaller.EXPECT().
					CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, &fakeDataError{data: revertHex}).
					Times(1),
				suite.mockEthCaller.EXPECT().
					CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(packedTextReturn(t, "https://node.example/"), nil).
					Times(1),
			)

			url, err := suite.sut.ResolveUrl("agent.collection.80002.6022.eth")

			assert.NoError(t, err)
			assert.Equal(t, "https://node.example", url)
		})
	})

	t.Run("Given the resolver does not revert", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)

			gomock.InOrder(
				registryDiscoveryReturningResolver(suite, t),
				suite.mockEthCaller.EXPECT().
					CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]byte{}, nil).
					Times(1),
			)

			_, err := suite.sut.ResolveUrl("agent.eth")

			assert.Error(t, err)
		})
	})

	t.Run("Given the revert carries no data", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a clear error about missing revert data", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)

			gomock.InOrder(
				registryDiscoveryReturningResolver(suite, t),
				suite.mockEthCaller.EXPECT().
					CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("execution reverted")).
					Times(1),
			)

			_, err := suite.sut.ResolveUrl("agent.eth")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "revert data")
		})
	})

	t.Run("Given the gateway fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			gatewayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer gatewayServer.Close()

			suite := WhenResolvingEnsUrlBeforeEach(t)

			revertHex := packOffchainLookupRevertHex(t, gatewayServer.URL+"/{sender}/{data}.json")

			gomock.InOrder(
				registryDiscoveryReturningResolver(suite, t),
				suite.mockEthCaller.EXPECT().
					CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, &fakeDataError{data: revertHex}).
					Times(1),
			)

			_, err := suite.sut.ResolveUrl("agent.eth")

			assert.Error(t, err)
		})
	})

	t.Run("Given the on-chain proof verification fails", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			gatewayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{"data":"0x1234"}`))
			}))
			defer gatewayServer.Close()

			suite := WhenResolvingEnsUrlBeforeEach(t)

			revertHex := packOffchainLookupRevertHex(t, gatewayServer.URL+"/{sender}/{data}.json")

			gomock.InOrder(
				registryDiscoveryReturningResolver(suite, t),
				suite.mockEthCaller.EXPECT().
					CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, &fakeDataError{data: revertHex}).
					Times(1),
				suite.mockEthCaller.EXPECT().
					CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("UntrustedSigner")).
					Times(1),
			)

			_, err := suite.sut.ResolveUrl("agent.eth")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "resolveWithProof")
		})
	})

	t.Run("Given no resolver is registered for the name or any parent", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error without calling any resolver", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)

			addressType, _ := abi.NewType("address", "", nil)
			zeroReturn, err := abi.Arguments{{Type: addressType}}.Pack(common.Address{})
			assert.NoError(t, err)

			suite.mockEthCaller.EXPECT().
				CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(zeroReturn, nil).
				Times(2)

			_, err = suite.sut.ResolveUrl("agent.eth")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "no resolver registered")
		})
	})

	t.Run("Given the revert is not an OffchainLookup", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return an error", func(t *testing.T) {
			t.Parallel()

			suite := WhenResolvingEnsUrlBeforeEach(t)

			gomock.InOrder(
				registryDiscoveryReturningResolver(suite, t),
				suite.mockEthCaller.EXPECT().
					CallContract(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, &fakeDataError{data: "0xdeadbeef00"}).
					Times(1),
			)

			_, err := suite.sut.ResolveUrl("agent.eth")

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "OffchainLookup")
		})
	})
}
