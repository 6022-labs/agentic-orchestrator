package services

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/6022-labs/agentic-orchestrator/agent-identity/src/registry_synchronizer_blockchain/settings"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/samber/lo"
)

// EIP-1474 "limit exceeded", used for rate limits.
const jsonRpcLimitExceededCode = -32005

const (
	rpcRetryAttempts = 8
	rpcRetryDelay    = time.Second
	rpcRetryMaxDelay = 30 * time.Second
)

type AbstractRpcRequester struct {
	client          *ethclient.Client
	rpcNodeSettings *settings.RpcNodeSettings
}

func (r *AbstractRpcRequester) filterEventLogs(
	ctx context.Context,
	from, to uint64,
	addresses []common.Address,
	contractAbi *abi.ABI,
	eventName string,
) ([]types.Log, error) {
	// An empty address list would match every contract on the chain.
	if len(addresses) == 0 {
		return nil, nil
	}

	event, ok := contractAbi.Events[eventName]
	if !ok {
		return nil, fmt.Errorf("ABI has no %s event", eventName)
	}

	var logs []types.Log
	for _, batch := range lo.Chunk(addresses, r.rpcNodeSettings.MaxAddressesPerLogQuery) {
		var batchLogs []types.Log
		err := r.retryRpc(ctx, func() error {
			var err error
			batchLogs, err = r.client.FilterLogs(ctx, ethereum.FilterQuery{
				FromBlock: new(big.Int).SetUint64(from),
				ToBlock:   new(big.Int).SetUint64(to),
				Addresses: batch,
				Topics:    [][]common.Hash{{event.ID}},
			})
			return err
		})
		if err != nil {
			return nil, fmt.Errorf("failed to filter %s logs: %w", eventName, err)
		}
		logs = append(logs, batchLogs...)
	}

	return logs, nil
}

func (r *AbstractRpcRequester) retryRpc(ctx context.Context, call func() error) error {
	return retry.Do(
		func() error {
			if err := call(); err != nil {
				if r.isRetryableRpcError(err) {
					return err
				}
				return retry.Unrecoverable(err)
			}
			return nil
		},
		retry.Context(ctx),
		retry.Attempts(rpcRetryAttempts),
		retry.Delay(rpcRetryDelay),
		retry.MaxDelay(rpcRetryMaxDelay),
		retry.DelayType(retry.CombineDelay(retry.BackOffDelay, retry.RandomDelay)),
		retry.LastErrorOnly(true),
	)
}

func (r *AbstractRpcRequester) isRetryableRpcError(err error) bool {
	var httpErr rpc.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == http.StatusTooManyRequests || httpErr.StatusCode >= http.StatusInternalServerError
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	var rpcErr rpc.Error
	if errors.As(err, &rpcErr) && rpcErr.ErrorCode() == jsonRpcLimitExceededCode {
		return true
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "rate limit") || strings.Contains(message, "too many requests")
}
