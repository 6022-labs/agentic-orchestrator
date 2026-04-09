# Blockchain Examples

## Example A: Contract event log listener (updatable)

**Use when** events arrive from RPC/WS subscriptions to contract logs and tracked addresses may change at runtime.

- Provider opens `SubscribeFilterLogs` with topic/address filters.
- Provider decodes log to typed event.
- Listener embeds `AbstractUpdatableSubscriptionListener` with a `rebuildSubscription` callback.
- Listener enqueues typed event into the queue via `Add*Events`.
- Listener calls `onBlockSeen` (unexported) for old subscription cleanup.
- Listener records metrics via `ListenerMetricsInterface`.
- Queue internally starts a finalization timer — no external signal needed.

```go
type LogProviderInterface interface {
    StartLogSubscription(chainId uint64, addresses []common.Address) (<-chan *TransferEvent, ethereum.Subscription, error)
}
```

## Example B: Dynamic subscription update with deferred rebuild

**Use when** tracked entities change at runtime (e.g. new vault or reward pool created).

```go
type DynamicSubscriberInterface interface {
    Subscribe(address common.Address, chainId uint64) error
}
```

Implementation pattern:

1. Check skip flag (`syncSettings.Skip*Events`) — return early if disabled.
2. Check if address is already tracked (`lo.Contains` guard) — return early if so.
3. Add the new address to `trackedAddresses`.
4. Call `markNeedsRebuild(chainId)` (unexported) — flag the chain, don't rebuild yet.
5. When the current block finishes processing, the processor calls `OnBlockProcessed` on all registered handlers.
6. The base's `OnBlockProcessed` calls `onBlockSeen`, then checks the `needsRebuild` flag. If set, clears it and invokes the `rebuildSubscription` callback.
7. `rebuildSubscription` creates a new subscription covering all addresses, moves old to `oldSubscriptions`, returns any error.
8. Old subscription is cleaned up when `onBlockSeen` fires on the new subscription.

This deferred pattern avoids rebuilding mid-block when events are still being processed. The rebuild happens at a safe point — after the processor has handled all events for the current block.

During the brief overlap after rebuild, both subscriptions may deliver the same event. This is safe because:
- Events for the same block are appended to the same queue entry (no duplication at dequeue).
- Downstream handlers use UPSERT semantics, making double-processing idempotent.

## Example C: Grouped host startup

**Use when** a host boots N listeners.

```go
type EventListenerInterface interface {
    SubscribeAll() error
    Listen() error
}
```

Host flow:

1. Resolve `[]EventListenerInterface` from DI group (`event_listeners`).
2. Call `SubscribeAll()` for each listener.
3. Start each `Listen()` loop in managed goroutine runner.
4. Stop all listeners on context cancellation.

## Example D: OnBlockProcessed handler registration

**Use when** updatable listeners need subscription rebuilds after a block is fully processed.

```go
type OnBlockProcessedHandlerInterface interface {
    OnBlockProcessed(chainId, blockNumber uint64)
}
```

The block events processor calls all registered `OnBlockProcessedHandlerInterface` implementations after processing a block. The base `AbstractUpdatableSubscriptionListener` implements this interface to:
1. Clean up old subscriptions via `onBlockSeen`.
2. Check if a subscription rebuild is pending (the `needsRebuild` flag).
3. Execute the `rebuildSubscription` callback if needed.

Concrete listeners don't implement this method — the base handles it. They only provide the `rebuildSubscription` callback to the constructor.

DI wiring: register each updatable listener as `EventListenerInterface` (group `event_listeners`), `OnBlockProcessedHandlerInterface` (group `on_block_processed_handlers`), and domain-specific subscriber (group e.g. `vault_event_subscribers`).

Static listeners are only registered in `event_listeners` — they don't participate in `OnBlockProcessed`.

## Example E: Timer-based block finalization

**Use instead of** subscribing to block headers (`newHeads`).

The event queue self-manages block finalization:

```go
// Chain settings — configure per chain
type ChainSettings struct {
    BlockFinalizationDelay time.Duration `koanf:"block_finalization_delay"`
}

// Inside the queue — triggered automatically by Add*Events
func (q *BlockEventsQueue) scheduleBlockFinalization(chainId, blockNumber uint64) {
    // Track highest block seen
    if blockNumber > q.latestBlockSeen[chainId] {
        q.latestBlockSeen[chainId] = blockNumber
    }

    // One timer per chain
    if _, ok := q.timers[chainId]; ok {
        return
    }

    q.timers[chainId] = time.AfterFunc(q.getBlockFinalizationDelay(chainId), func() {
        completedBlock := q.latestBlockSeen[chainId]
        q.latestCompletedBlock[chainId] = completedBlock
        q.completePendingBlocks()
    })
}
```

This works for all chain types:
- **Fast chains** (Ethereum ~12s, L2s ~2s): set a short delay (e.g., 2s).
- **On-demand chains** (blocks only when transactions exist): the timer starts when an event arrives, regardless of block interval.

## Example F: Static factory listener

**Use when** the listener subscribes to fixed addresses from configuration that never change at runtime (e.g., factory contracts).

```go
type FactoryEventListener struct {
    logger         *zap.Logger
    chainsSettings *ChainsSettings
    syncSettings   *SyncSettings

    blockEventsQueue BlockEventsQueueInterface
    provider         FactoryEventSubscriptionProviderInterface
    listenerMetrics  ListenerMetricsInterface

    subscriptions    map[uint64]ethereum.Subscription
    trackedAddresses map[uint64][]common.Address

    mu    sync.Mutex
    errCh chan error
    wg    sync.WaitGroup
}
```

Key differences from updatable listeners:
- Does NOT embed `AbstractUpdatableSubscriptionListener`.
- No `Subscribe()`, no `markNeedsRebuild`, no `onBlockSeen`, no `OnBlockProcessed`.
- Manages its own `subscriptions`, `errCh`, `mu`, `wg` directly.
- Addresses come from `syncSettings` (configuration), not from a repository.
- Handler goroutine does NOT call `onBlockSeen` (no old subscriptions to clean up).
- Only registered in `event_listeners` DI group — not in `on_block_processed_handlers` or any subscriber group.
