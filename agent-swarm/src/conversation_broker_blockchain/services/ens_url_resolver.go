package services

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	contract_abi "github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/abi"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/clients"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/helpers"
	"github.com/6022-labs/agentic-orchestrator/agent-swarm/src/conversation_broker_blockchain/settings"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

const urlRecordKey = "url"

// offchainLookup is the decoded ERC-3668 revert payload.
type offchainLookup struct {
	Sender    common.Address
	Urls      []string
	CallData  []byte
	ExtraData []byte
}

// EnsUrlResolver is an ENSIP-10 + ERC-3668 (CCIP-Read) client: it discovers
// the name's resolver through the ENS registry, follows the OffchainLookup
// revert to the gateway, then verifies the signed answer via resolveWithProof.
type EnsUrlResolver struct {
	logger              *zap.Logger
	ethCaller           clients.EthCallerInterface
	httpClient          *http.Client
	ensResolverSettings *settings.EnsResolverSettings

	ensRegistryAbi      *abi.ABI
	ensTextResolverAbi  *abi.ABI
	offchainResolverAbi *abi.ABI
	offchainLookupError abi.Error
}

type NewEnsUrlResolverParams struct {
	dig.In

	Logger              *zap.Logger
	EthCaller           clients.EthCallerInterface
	HttpClient          *http.Client `name:"ens-gateway-http-client"`
	EnsResolverSettings *settings.EnsResolverSettings
}

func NewEnsUrlResolver(params NewEnsUrlResolverParams) *EnsUrlResolver {
	// The call encoding comes from the generated contract bindings in ../abi;
	// a parse failure means broken generated code, so it panics at wiring time.
	ensRegistryAbi, err := contract_abi.EnsRegistryMetaData.GetAbi()
	if err != nil {
		panic(fmt.Sprintf("invalid generated ens registry ABI: %v", err))
	}
	ensTextResolverAbi, err := contract_abi.EnsTextResolverMetaData.GetAbi()
	if err != nil {
		panic(fmt.Sprintf("invalid generated ens text resolver ABI: %v", err))
	}
	offchainResolverAbi, err := contract_abi.AgentEnsOffchainResolverMetaData.GetAbi()
	if err != nil {
		panic(fmt.Sprintf("invalid generated offchain resolver ABI: %v", err))
	}

	return &EnsUrlResolver{
		logger:              params.Logger,
		ethCaller:           params.EthCaller,
		httpClient:          params.HttpClient,
		ensResolverSettings: params.EnsResolverSettings,

		ensRegistryAbi:      ensRegistryAbi,
		ensTextResolverAbi:  ensTextResolverAbi,
		offchainResolverAbi: offchainResolverAbi,
		offchainLookupError: offchainResolverAbi.Errors["OffchainLookup"],
	}
}

func (r *EnsUrlResolver) ResolveUrl(ensDomain string) (string, error) {
	dnsName, err := helpers.DnsEncode(ensDomain)
	if err != nil {
		return "", fmt.Errorf("failed to dns-encode %s: %w", ensDomain, err)
	}

	textQuery, err := r.ensTextResolverAbi.Pack("text", helpers.Namehash(ensDomain), urlRecordKey)
	if err != nil {
		return "", err
	}

	resolveCall, err := r.offchainResolverAbi.Pack("resolve", dnsName, textQuery)
	if err != nil {
		return "", err
	}

	resolverAddress, err := r.discoverResolver(ensDomain)
	if err != nil {
		return "", err
	}

	lookup, err := r.callResolveExpectingOffchainLookup(resolverAddress, resolveCall)
	if err != nil {
		return "", err
	}

	gatewayResponse, err := r.queryGateway(lookup)
	if err != nil {
		return "", err
	}

	verifiedResult, err := r.verifyWithProof(lookup.Sender, gatewayResponse, lookup.ExtraData)
	if err != nil {
		return "", err
	}

	textValues, err := r.ensTextResolverAbi.Unpack("text", verifiedResult)
	if err != nil {
		return "", fmt.Errorf("failed to decode text record: %w", err)
	}
	agentUrl := textValues[0].(string)

	return strings.TrimRight(strings.TrimSpace(agentUrl), "/"), nil
}

// discoverResolver walks the name hierarchy in the ENS registry (ENSIP-10):
// the first parent carrying a resolver serves the whole subtree.
func (r *EnsUrlResolver) discoverResolver(ensDomain string) (common.Address, error) {
	registryAddress := r.ensResolverSettings.RegistryCommonAddress()

	labels := strings.Split(ensDomain, ".")
	for i := range labels {
		name := strings.Join(labels[i:], ".")

		query, err := r.ensRegistryAbi.Pack("resolver", helpers.Namehash(name))
		if err != nil {
			return common.Address{}, err
		}

		returnData, err := r.ethCaller.CallContract(context.Background(), ethereum.CallMsg{
			To:   &registryAddress,
			Data: query,
		}, nil)
		if err != nil {
			return common.Address{}, fmt.Errorf("registry resolver lookup failed for %s: %w", name, err)
		}

		resolverValues, err := r.ensRegistryAbi.Unpack("resolver", returnData)
		if err != nil {
			return common.Address{}, fmt.Errorf("failed to decode resolver address: %w", err)
		}
		resolverAddress := resolverValues[0].(common.Address)
		if resolverAddress != (common.Address{}) {
			return resolverAddress, nil
		}
	}

	return common.Address{}, fmt.Errorf("no resolver registered for %s or any parent", ensDomain)
}

func (r *EnsUrlResolver) callResolveExpectingOffchainLookup(resolverAddress common.Address, resolveCall []byte) (*offchainLookup, error) {
	_, err := r.ethCaller.CallContract(context.Background(), ethereum.CallMsg{
		To:   &resolverAddress,
		Data: resolveCall,
	}, nil)
	if err == nil {
		return nil, fmt.Errorf("resolver did not revert with OffchainLookup")
	}

	var dataErr rpc.DataError
	if !errors.As(err, &dataErr) {
		return nil, fmt.Errorf("resolver call failed without revert data (rpc must return revert data): %w", err)
	}

	revertHex, ok := dataErr.ErrorData().(string)
	if !ok {
		return nil, fmt.Errorf("unexpected revert data type %T", dataErr.ErrorData())
	}

	revertData, err := hex.DecodeString(strings.TrimPrefix(revertHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode revert data: %w", err)
	}

	return r.decodeOffchainLookup(revertData)
}

// decodeOffchainLookup parses an ERC-3668 OffchainLookup revert payload using
// the resolver contract's own error definition.
func (r *EnsUrlResolver) decodeOffchainLookup(revertData []byte) (*offchainLookup, error) {
	if len(revertData) < 4 || !bytes.Equal(revertData[:4], r.offchainLookupError.ID.Bytes()[:4]) {
		return nil, fmt.Errorf("revert is not an OffchainLookup")
	}

	values, err := r.offchainLookupError.Inputs.UnpackValues(revertData[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode OffchainLookup: %w", err)
	}

	return &offchainLookup{
		Sender:    values[0].(common.Address),
		Urls:      values[1].([]string),
		CallData:  values[2].([]byte),
		ExtraData: values[4].([]byte),
	}, nil
}

func (r *EnsUrlResolver) queryGateway(lookup *offchainLookup) ([]byte, error) {
	if len(lookup.Urls) == 0 {
		return nil, fmt.Errorf("OffchainLookup carries no gateway url")
	}

	sender := strings.ToLower(lookup.Sender.Hex())
	data := "0x" + hex.EncodeToString(lookup.CallData)

	var lastErr error
	for _, urlTemplate := range lookup.Urls {
		url := strings.ReplaceAll(strings.ReplaceAll(urlTemplate, "{sender}", sender), "{data}", data)

		var httpResponse *http.Response
		var err error
		if strings.Contains(urlTemplate, "{data}") {
			httpResponse, err = r.httpClient.Get(url)
		} else {
			body := fmt.Sprintf(`{"sender":%q,"data":%q}`, sender, data)
			httpResponse, err = r.httpClient.Post(url, "application/json", strings.NewReader(body))
		}
		if err != nil {
			lastErr = err
			continue
		}

		responseBody, err := func() ([]byte, error) {
			defer httpResponse.Body.Close()
			if httpResponse.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(httpResponse.Body)
				return nil, fmt.Errorf("gateway returned %d: %s", httpResponse.StatusCode, string(body))
			}
			return io.ReadAll(httpResponse.Body)
		}()
		if err != nil {
			lastErr = err
			continue
		}

		var payload struct {
			Data string `json:"data"`
		}
		if err := json.Unmarshal(responseBody, &payload); err != nil {
			lastErr = fmt.Errorf("failed to decode gateway response: %w", err)
			continue
		}

		decoded, err := hex.DecodeString(strings.TrimPrefix(payload.Data, "0x"))
		if err != nil {
			lastErr = fmt.Errorf("failed to decode gateway data: %w", err)
			continue
		}

		return decoded, nil
	}

	return nil, fmt.Errorf("all CCIP gateways failed: %w", lastErr)
}

func (r *EnsUrlResolver) verifyWithProof(resolverAddress common.Address, response, extraData []byte) ([]byte, error) {
	proofCall, err := r.offchainResolverAbi.Pack("resolveWithProof", response, extraData)
	if err != nil {
		return nil, err
	}

	returnData, err := r.ethCaller.CallContract(context.Background(), ethereum.CallMsg{
		To:   &resolverAddress,
		Data: proofCall,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("resolveWithProof verification failed: %w", err)
	}

	proofValues, err := r.offchainResolverAbi.Unpack("resolveWithProof", returnData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode resolveWithProof result: %w", err)
	}
	return proofValues[0].([]byte), nil
}
