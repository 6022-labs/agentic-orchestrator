package services

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	registry_abi "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
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

func (a *AgentCollectionRequester) GetAgentOwner(collectionAddress common.Address, collectionTokenId big.Int) (*common.Address, error) {
	agentCollection, err := registry_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	owner, err := agentCollection.OwnerOf(nil, &collectionTokenId)
	if err != nil {
		return nil, err
	}

	return &owner, nil
}

func (a *AgentCollectionRequester) GetAgentPreviousSigners(collectionAddress common.Address, collectionTokenId big.Int) ([]common.Address, error) {
	agentCollection, err := registry_abi.NewAgentCollectionV1(collectionAddress, a.client)
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

func (a *AgentCollectionRequester) GetAgentInformation(collectionAddress common.Address, collectionTokenId big.Int) (*registry_abi.AgentInformation, error) {
	agentCollection, err := registry_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	info, err := agentCollection.InformationOf(nil, &collectionTokenId)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (a *AgentCollectionRequester) GetMintProposal(collectionAddress common.Address, proposalId *big.Int, atBlock uint64) (*registry_abi.MintProposal, error) {
	agentCollection, err := registry_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	mintProposals, err := a.getMintProposals(agentCollection, collectionAddress, atBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to get mint proposals: %w", err)
	}

	mintProposal, found := lo.Find(mintProposals, func(mintProposal *registry_abi.MintProposal) bool {
		return mintProposal.Id.Cmp(proposalId) == 0
	})
	if !found {
		return nil, fmt.Errorf("mint proposal with ID %s not found", proposalId.String())
	}

	return mintProposal, nil
}

func (a *AgentCollectionRequester) getMintProposals(agentCollection *registry_abi.AgentCollectionV1, agentCollectionAddress common.Address, atBlock uint64) ([]*registry_abi.MintProposal, error) {
	blockBigInt := new(big.Int).SetUint64(atBlock)
	blockHex := fmt.Sprintf("0x%x", atBlock)

	mintProposalLength, err := agentCollection.MintProposalsLength(&bind.CallOpts{BlockNumber: blockBigInt})
	if err != nil {
		return nil, fmt.Errorf("failed to get mint proposal length: %w", err)
	}

	agentCollectionABI, err := abi.JSON(strings.NewReader(registry_abi.AgentCollectionV1ABI)) // You must expose raw ABI string
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
			Args:   []interface{}{msg, blockHex},
			Result: &results[i],
		})
	}

	if err := a.client.Client().BatchCall(batch); err != nil {
		return nil, fmt.Errorf("failed to batch call mint proposals: %w", err)
	}

	mintProposals := make([]*registry_abi.MintProposal, 0, mintProposalLength.Int64())
	for _, result := range results {
		if result == nil {
			continue
		}

		unpacked, err := agentCollectionABI.Unpack("mintProposal", result)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack mintProposal result: %w", err)
		}
		mintProposal := *abi.ConvertType(unpacked[0], new(registry_abi.MintProposal)).(*registry_abi.MintProposal)

		mintProposals = append(mintProposals, &mintProposal)
	}

	return mintProposals, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1MintedEvent(from, to uint64, fromAddresses []common.Address) ([]*registry_abi.AgentCollectionV1Minted, error) {
	var allEvents []*registry_abi.AgentCollectionV1Minted

	for _, address := range fromAddresses {
		agentCollection, err := registry_abi.NewAgentCollectionV1(address, a.client)
		if err != nil {
			return nil, err
		}

		iter, err := agentCollection.FilterMinted(&bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: context.Background(),
		}, nil, nil)
		if err != nil {
			return nil, err
		}

		for iter.Next() {
			allEvents = append(allEvents, iter.Event)
		}
	}

	return allEvents, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1TransferEvent(from, to uint64, fromAddresses []common.Address) ([]*registry_abi.AgentCollectionV1Transfer, error) {
	var allEvents []*registry_abi.AgentCollectionV1Transfer

	for _, address := range fromAddresses {
		agentCollection, err := registry_abi.NewAgentCollectionV1(address, a.client)
		if err != nil {
			return nil, err
		}

		iter, err := agentCollection.FilterTransfer(&bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: context.Background(),
		}, nil, nil, nil)
		if err != nil {
			return nil, err
		}

		for iter.Next() {
			allEvents = append(allEvents, iter.Event)
		}
	}

	return allEvents, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1AgentImageUpdatedEvent(from, to uint64, fromAddresses []common.Address) ([]*registry_abi.AgentCollectionV1AgentImageUpdated, error) {
	var allEvents []*registry_abi.AgentCollectionV1AgentImageUpdated

	for _, address := range fromAddresses {
		agentCollection, err := registry_abi.NewAgentCollectionV1(address, a.client)
		if err != nil {
			return nil, err
		}

		iter, err := agentCollection.FilterAgentImageUpdated(&bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: context.Background(),
		}, nil, nil)
		if err != nil {
			return nil, err
		}

		for iter.Next() {
			allEvents = append(allEvents, iter.Event)
		}
	}

	return allEvents, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1AgentSignerUpdatedEvent(from, to uint64, fromAddresses []common.Address) ([]*registry_abi.AgentCollectionV1AgentSignerUpdated, error) {
	var allEvents []*registry_abi.AgentCollectionV1AgentSignerUpdated

	for _, address := range fromAddresses {
		agentCollection, err := registry_abi.NewAgentCollectionV1(address, a.client)
		if err != nil {
			return nil, err
		}

		iter, err := agentCollection.FilterAgentSignerUpdated(&bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: context.Background(),
		}, nil)
		if err != nil {
			return nil, err
		}

		for iter.Next() {
			allEvents = append(allEvents, iter.Event)
		}
	}

	return allEvents, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1AgentAttributeUpdatedEvent(from, to uint64, fromAddresses []common.Address) ([]*registry_abi.AgentCollectionV1AgentAttributeUpdated, error) {
	var allEvents []*registry_abi.AgentCollectionV1AgentAttributeUpdated

	for _, address := range fromAddresses {
		agentCollection, err := registry_abi.NewAgentCollectionV1(address, a.client)
		if err != nil {
			return nil, err
		}

		iter, err := agentCollection.FilterAgentAttributeUpdated(&bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: context.Background(),
		}, nil, nil)
		if err != nil {
			return nil, err
		}

		for iter.Next() {
			allEvents = append(allEvents, iter.Event)
		}
	}

	return allEvents, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1MintProposalCreatedEvent(from, to uint64, fromAddresses []common.Address) ([]*registry_abi.AgentCollectionV1MintProposalCreated, error) {
	var allEvents []*registry_abi.AgentCollectionV1MintProposalCreated

	for _, address := range fromAddresses {
		agentCollection, err := registry_abi.NewAgentCollectionV1(address, a.client)
		if err != nil {
			return nil, err
		}

		iter, err := agentCollection.FilterMintProposalCreated(&bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: context.Background(),
		}, nil, nil)
		if err != nil {
			return nil, err
		}

		for iter.Next() {
			allEvents = append(allEvents, iter.Event)
		}
	}

	return allEvents, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1MintProposalRemovedEvent(from, to uint64, fromAddresses []common.Address) ([]*registry_abi.AgentCollectionV1MintProposalRemoved, error) {
	var allEvents []*registry_abi.AgentCollectionV1MintProposalRemoved

	for _, address := range fromAddresses {
		agentCollection, err := registry_abi.NewAgentCollectionV1(address, a.client)
		if err != nil {
			return nil, err
		}

		iter, err := agentCollection.FilterMintProposalRemoved(&bind.FilterOpts{
			Start:   from,
			End:     &to,
			Context: context.Background(),
		}, nil, nil)
		if err != nil {
			return nil, err
		}

		for iter.Next() {
			allEvents = append(allEvents, iter.Event)
		}
	}

	return allEvents, nil
}
