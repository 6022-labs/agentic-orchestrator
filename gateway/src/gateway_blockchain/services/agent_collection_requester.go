package services

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	orchestrator_abi "github.com/6022-labs/agentic-orchestrator/gateway/src/gateway/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/samber/lo"
	"go.uber.org/dig"
)

type AgentCollectionRequester struct {
	client *ethclient.Client
}

type newAgentCollectionRequesterParams struct {
	dig.In

	Client *ethclient.Client `name:"http"`
}

func NewAgentCollectionRequester(
	params newAgentCollectionRequesterParams,
) *AgentCollectionRequester {
	return &AgentCollectionRequester{
		client: params.Client,
	}
}

func (a *AgentCollectionRequester) GetName(collectionAddress common.Address) (*string, error) {
	agentCollection, err := orchestrator_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	name, err := agentCollection.Name(nil)
	if err != nil {
		return nil, err
	}

	return &name, nil
}

func (a *AgentCollectionRequester) GetSymbol(collectionAddress common.Address) (*string, error) {
	agentCollection, err := orchestrator_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	symbol, err := agentCollection.Symbol(nil)
	if err != nil {
		return nil, err
	}

	return &symbol, nil
}

func (a *AgentCollectionRequester) GetNextTokenId(collectionAddress common.Address) (*big.Int, error) {
	agentCollection, err := orchestrator_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	nextTokenId, err := agentCollection.NextTokenId(nil)
	if err != nil {
		return nil, err
	}

	return nextTokenId, nil
}

func (a *AgentCollectionRequester) GetAgentOwner(collectionAddress common.Address, collectionTokenId big.Int) (*common.Address, error) {
	agentCollection, err := orchestrator_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	owner, err := agentCollection.OwnerOf(nil, &collectionTokenId)
	if err != nil {
		return nil, err
	}

	return &owner, nil
}

func (a *AgentCollectionRequester) GetAgentSigner(collectionAddress common.Address, collectionTokenId big.Int) (*common.Address, error) {
	agentCollection, err := orchestrator_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	signer, err := agentCollection.SignerOf(nil, &collectionTokenId)
	if err != nil {
		return nil, err
	}

	return &signer, nil
}

func (a *AgentCollectionRequester) GetAgentWallets(collectionAddress common.Address, collectionTokenId big.Int) ([]common.Address, error) {
	agentCollection, err := orchestrator_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	count, err := agentCollection.WalletsCountOf(nil, &collectionTokenId)
	if err != nil {
		return nil, err
	}

	wallets := make([]common.Address, count.Int64())

	for i := big.NewInt(0); i.Cmp(count) < 0; i.Add(i, big.NewInt(1)) {
		wallet, err := agentCollection.WalletOfByIndex(nil, &collectionTokenId, i)
		if err != nil {
			return nil, err
		}
		wallets[i.Int64()] = wallet
	}

	return wallets, nil
}

func (a *AgentCollectionRequester) GetAgentPreviousSigners(collectionAddress common.Address, collectionTokenId big.Int) ([]common.Address, error) {
	agentCollection, err := orchestrator_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	iter, err := agentCollection.FilterAgentSignerUpdated(&bind.FilterOpts{
		Start:   1,
		End:     nil,
		Context: context.Background(),
	}, []*big.Int{&collectionTokenId})
	if err != nil {
		return nil, err
	}

	var signers []common.Address

	for iter.Next() {
		event := iter.Event
		signers = append(signers, event.NewSigner)
	}

	return signers, nil
}

func (a *AgentCollectionRequester) GetAgentInformation(collectionAddress common.Address, collectionTokenId big.Int) (*orchestrator_abi.AgentInformation, error) {
	agentCollection, err := orchestrator_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	info, err := agentCollection.InformationOf(nil, &collectionTokenId)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (a *AgentCollectionRequester) GetAgentWalletsOfProposal(collectionAddress common.Address, proposalId *big.Int) ([]common.Address, error) {
	agentCollection, err := orchestrator_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	mintProposals, err := a.getMintProposals(agentCollection, collectionAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get mint proposals: %w", err)
	}

	mintProposal, found := lo.Find(mintProposals, func(mintProposal *orchestrator_abi.MintProposal) bool {
		return mintProposal.Id.Cmp(proposalId) == 0
	})
	if !found {
		return nil, fmt.Errorf("mint proposal with ID %s not found", proposalId.String())
	}

	return mintProposal.Wallets, nil
}

func (a *AgentCollectionRequester) getMintProposals(agentCollection *orchestrator_abi.AgentCollectionV1, agentCollectionAddress common.Address) ([]*orchestrator_abi.MintProposal, error) {
	mintProposalLength, err := agentCollection.MintProposalsLength(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get mint proposal length: %w", err)
	}

	agentCollectionABI, err := abi.JSON(strings.NewReader(orchestrator_abi.AgentCollectionV1ABI)) // You must expose raw ABI string
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent collection ABI: %w", err)
	}

	batch := make([]rpc.BatchElem, 0, mintProposalLength.Int64())
	results := make([]hexutil.Bytes, mintProposalLength.Int64())

	for i := uint64(0); i < mintProposalLength.Uint64(); i++ {
		data, err := agentCollectionABI.Pack("mintProposal", big.NewInt(int64(i)))
		if err != nil {
			return nil, fmt.Errorf("failed to pack mintProposal call: %w", err)
		}

		msg := map[string]interface{}{
			"to":   agentCollectionAddress.Hex(),
			"data": "0x" + common.Bytes2Hex(data),
		}

		batch = append(batch, rpc.BatchElem{
			Method: "eth_call",
			Args:   []interface{}{msg, "latest"},
			Result: &results[i],
		})
	}

	if err := a.client.Client().BatchCall(batch); err != nil {
		return nil, fmt.Errorf("failed to batch call mint proposals: %w", err)
	}

	mintProposals := make([]*orchestrator_abi.MintProposal, 0, mintProposalLength.Int64())
	for _, result := range results {
		if result == nil {
			continue
		}

		unpacked, err := agentCollectionABI.Unpack("mintProposal", result)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack mintProposal result: %w", err)
		}
		mintProposal := *abi.ConvertType(unpacked[0], new(orchestrator_abi.MintProposal)).(*orchestrator_abi.MintProposal)

		mintProposals = append(mintProposals, &mintProposal)
	}

	return mintProposals, nil
}

func (a *AgentCollectionRequester) GetAddOrUpdateAgentAttributeInputsFromTransaction(collectionAddress common.Address, txHash common.Hash) (*big.Int, *string, *string, error) {
	tx, isPending, err := a.client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get transaction by hash: %w", err)
	}
	if isPending {
		return nil, nil, nil, fmt.Errorf("transaction is still pending")
	}

	agentCollectionABI, err := abi.JSON(strings.NewReader(orchestrator_abi.AgentCollectionV1ABI))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse agent collection ABI: %w", err)
	}

	data := tx.Data()
	if len(data) < 4 {
		return nil, nil, nil, fmt.Errorf("transaction data too short")
	}

	method, err := agentCollectionABI.MethodById(data[:4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get method by ID: %w", err)
	}

	if method.Name != "addOrUpdateAgentAttribute" {
		return nil, nil, nil, fmt.Errorf("transaction does not call addAgentAttribute or updateAgentAttribute")
	}

	argsAsMap := make(map[string]interface{})
	err = method.Inputs.UnpackIntoMap(argsAsMap, data[4:])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unpack inputs into map: %w", err)
	}

	tokenId, ok := argsAsMap["tokenId"].(*big.Int)
	if !ok {
		return nil, nil, nil, fmt.Errorf("failed to assert TokenId type, got %T with map %v", argsAsMap["tokenId"], argsAsMap)
	}
	key, ok := argsAsMap["key"].(string)
	if !ok {
		return nil, nil, nil, fmt.Errorf("failed to assert Key type, got %T with map %v", argsAsMap["key"], argsAsMap)
	}
	value, ok := argsAsMap["value"].(string)
	if !ok {
		return nil, nil, nil, fmt.Errorf("failed to assert Value type, got %T with map %v", argsAsMap["value"], argsAsMap)
	}

	return tokenId, &key, &value, nil
}
