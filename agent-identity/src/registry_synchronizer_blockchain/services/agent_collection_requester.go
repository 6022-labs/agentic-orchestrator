package services

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	registry_abi "github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_common/abi"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/samber/lo"
	"go.uber.org/dig"
)

// maxAddressesPerLogQuery stays under geth's default 1000-address eth_getLogs cap.
const maxAddressesPerLogQuery = 500

type AgentCollectionRequester struct {
	client      *ethclient.Client
	contractAbi *abi.ABI
	filterer    *registry_abi.AgentCollectionV1Filterer
}

type NewAgentCollectionRequesterParams struct {
	dig.In

	Client *ethclient.Client `name:"http"`
}

func NewAgentCollectionRequester(
	params NewAgentCollectionRequesterParams,
) (*AgentCollectionRequester, error) {
	contractAbi, err := registry_abi.AgentCollectionV1MetaData.GetAbi()
	if err != nil {
		return nil, fmt.Errorf("failed to parse agent collection ABI: %w", err)
	}

	filterer, err := registry_abi.NewAgentCollectionV1Filterer(common.Address{}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent collection filterer: %w", err)
	}

	return &AgentCollectionRequester{
		client:      params.Client,
		contractAbi: contractAbi,
		filterer:    filterer,
	}, nil
}

func (a *AgentCollectionRequester) GetAgentOwner(
	ctx context.Context,
	collectionAddress common.Address,
	collectionTokenId big.Int,
) (*common.Address, error) {
	agentCollection, err := registry_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	owner, err := agentCollection.OwnerOf(&bind.CallOpts{Context: ctx}, &collectionTokenId)
	if err != nil {
		return nil, err
	}

	return &owner, nil
}

func (a *AgentCollectionRequester) GetAgentInformation(
	ctx context.Context,
	collectionAddress common.Address,
	collectionTokenId big.Int,
) (*registry_abi.AgentInformation, error) {
	agentCollection, err := registry_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	info, err := agentCollection.InformationOf(&bind.CallOpts{Context: ctx}, &collectionTokenId)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (a *AgentCollectionRequester) GetMintProposal(
	ctx context.Context,
	collectionAddress common.Address,
	proposalId *big.Int,
	atBlock uint64,
) (*registry_abi.MintProposal, error) {
	agentCollection, err := registry_abi.NewAgentCollectionV1(collectionAddress, a.client)
	if err != nil {
		return nil, err
	}

	mintProposals, err := a.getMintProposals(ctx, agentCollection, collectionAddress, atBlock)
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

func (a *AgentCollectionRequester) getMintProposals(
	ctx context.Context,
	agentCollection *registry_abi.AgentCollectionV1,
	agentCollectionAddress common.Address,
	atBlock uint64,
) ([]*registry_abi.MintProposal, error) {
	blockBigInt := new(big.Int).SetUint64(atBlock)
	blockHex := fmt.Sprintf("0x%x", atBlock)

	mintProposalLength, err := agentCollection.MintProposalsLength(
		&bind.CallOpts{Context: ctx, BlockNumber: blockBigInt},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get mint proposal length: %w", err)
	}

	agentCollectionABI, err := abi.JSON(
		strings.NewReader(registry_abi.AgentCollectionV1ABI),
	) // You must expose raw ABI string
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
			"to":   strings.ToLower(agentCollectionAddress.Hex()),
			"data": "0x" + common.Bytes2Hex(data),
		}

		batch = append(batch, rpc.BatchElem{
			Method: "eth_call",
			Args:   []interface{}{msg, blockHex},
			Result: &results[i],
		})
	}

	if err := a.client.Client().BatchCallContext(ctx, batch); err != nil {
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

func (a *AgentCollectionRequester) GetAgentCollectionV1MintedEvent(
	ctx context.Context,
	from, to uint64,
	fromAddresses []common.Address,
) ([]*registry_abi.AgentCollectionV1Minted, error) {
	logs, err := a.filterEventLogs(ctx, from, to, fromAddresses, "Minted")
	if err != nil {
		return nil, err
	}

	events := make([]*registry_abi.AgentCollectionV1Minted, 0, len(logs))
	for _, log := range logs {
		event, err := a.filterer.ParseMinted(log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Minted log: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1TransferEvent(
	ctx context.Context,
	from, to uint64,
	fromAddresses []common.Address,
) ([]*registry_abi.AgentCollectionV1Transfer, error) {
	logs, err := a.filterEventLogs(ctx, from, to, fromAddresses, "Transfer")
	if err != nil {
		return nil, err
	}

	events := make([]*registry_abi.AgentCollectionV1Transfer, 0, len(logs))
	for _, log := range logs {
		event, err := a.filterer.ParseTransfer(log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Transfer log: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1AgentImageUpdatedEvent(
	ctx context.Context,
	from, to uint64,
	fromAddresses []common.Address,
) ([]*registry_abi.AgentCollectionV1AgentImageUpdated, error) {
	logs, err := a.filterEventLogs(ctx, from, to, fromAddresses, "AgentImageUpdated")
	if err != nil {
		return nil, err
	}

	events := make([]*registry_abi.AgentCollectionV1AgentImageUpdated, 0, len(logs))
	for _, log := range logs {
		event, err := a.filterer.ParseAgentImageUpdated(log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AgentImageUpdated log: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1AgentImageRemovedEvent(
	ctx context.Context,
	from, to uint64,
	fromAddresses []common.Address,
) ([]*registry_abi.AgentCollectionV1AgentImageRemoved, error) {
	logs, err := a.filterEventLogs(ctx, from, to, fromAddresses, "AgentImageRemoved")
	if err != nil {
		return nil, err
	}

	events := make([]*registry_abi.AgentCollectionV1AgentImageRemoved, 0, len(logs))
	for _, log := range logs {
		event, err := a.filterer.ParseAgentImageRemoved(log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AgentImageRemoved log: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1AgentAttributeUpdatedEvent(
	ctx context.Context,
	from, to uint64,
	fromAddresses []common.Address,
) ([]*registry_abi.AgentCollectionV1AgentAttributeUpdated, error) {
	logs, err := a.filterEventLogs(ctx, from, to, fromAddresses, "AgentAttributeUpdated")
	if err != nil {
		return nil, err
	}

	events := make([]*registry_abi.AgentCollectionV1AgentAttributeUpdated, 0, len(logs))
	for _, log := range logs {
		event, err := a.filterer.ParseAgentAttributeUpdated(log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AgentAttributeUpdated log: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1MintProposalCreatedEvent(
	ctx context.Context,
	from, to uint64,
	fromAddresses []common.Address,
) ([]*registry_abi.AgentCollectionV1MintProposalCreated, error) {
	logs, err := a.filterEventLogs(ctx, from, to, fromAddresses, "MintProposalCreated")
	if err != nil {
		return nil, err
	}

	events := make([]*registry_abi.AgentCollectionV1MintProposalCreated, 0, len(logs))
	for _, log := range logs {
		event, err := a.filterer.ParseMintProposalCreated(log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MintProposalCreated log: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1MintProposalRemovedEvent(
	ctx context.Context,
	from, to uint64,
	fromAddresses []common.Address,
) ([]*registry_abi.AgentCollectionV1MintProposalRemoved, error) {
	logs, err := a.filterEventLogs(ctx, from, to, fromAddresses, "MintProposalRemoved")
	if err != nil {
		return nil, err
	}

	events := make([]*registry_abi.AgentCollectionV1MintProposalRemoved, 0, len(logs))
	for _, log := range logs {
		event, err := a.filterer.ParseMintProposalRemoved(log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MintProposalRemoved log: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1AgentAddressAddedEvent(
	ctx context.Context,
	from, to uint64,
	fromAddresses []common.Address,
) ([]*registry_abi.AgentCollectionV1AgentAddressAdded, error) {
	logs, err := a.filterEventLogs(ctx, from, to, fromAddresses, "AgentAddressAdded")
	if err != nil {
		return nil, err
	}

	events := make([]*registry_abi.AgentCollectionV1AgentAddressAdded, 0, len(logs))
	for _, log := range logs {
		event, err := a.filterer.ParseAgentAddressAdded(log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AgentAddressAdded log: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (a *AgentCollectionRequester) GetAgentCollectionV1AgentAddressRemovedEvent(
	ctx context.Context,
	from, to uint64,
	fromAddresses []common.Address,
) ([]*registry_abi.AgentCollectionV1AgentAddressRemoved, error) {
	logs, err := a.filterEventLogs(ctx, from, to, fromAddresses, "AgentAddressRemoved")
	if err != nil {
		return nil, err
	}

	events := make([]*registry_abi.AgentCollectionV1AgentAddressRemoved, 0, len(logs))
	for _, log := range logs {
		event, err := a.filterer.ParseAgentAddressRemoved(log)
		if err != nil {
			return nil, fmt.Errorf("failed to parse AgentAddressRemoved log: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (a *AgentCollectionRequester) filterEventLogs(
	ctx context.Context,
	from, to uint64,
	addresses []common.Address,
	eventName string,
) ([]types.Log, error) {
	// An empty address list would match every contract on the chain.
	if len(addresses) == 0 {
		return nil, nil
	}

	event, ok := a.contractAbi.Events[eventName]
	if !ok {
		return nil, fmt.Errorf("agent collection ABI has no %s event", eventName)
	}

	var logs []types.Log
	for _, batch := range lo.Chunk(addresses, maxAddressesPerLogQuery) {
		batchLogs, err := a.client.FilterLogs(ctx, ethereum.FilterQuery{
			FromBlock: new(big.Int).SetUint64(from),
			ToBlock:   new(big.Int).SetUint64(to),
			Addresses: batch,
			Topics:    [][]common.Hash{{event.ID}},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to filter %s logs: %w", eventName, err)
		}
		logs = append(logs, batchLogs...)
	}

	return logs, nil
}
