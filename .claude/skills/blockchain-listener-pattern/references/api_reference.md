# Blockchain Inbound Listener + Outbound Subscription Provider (Compact Scaffold)

Use this scaffold to implement adapter-only blockchain event ingestion without embedding business rules.

## Shape

1. **Core ports (domain/application layer):**
   - `SomethingEventSubscriptionProviderInterface` for transport subscription startup.
   - `BlockEventsQueueInterface` for enqueuing typed events with chain/block metadata.
   - `OnBlockProcessedHandlerInterface` for the base to react when a block is fully processed (subscription rebuilds).
   - Optional `SomethingEventSubscriberInterface` for runtime incremental subscriptions.
2. **Outbound adapter:**
   - Start/stop blockchain transport subscription (RPC/WS) for contract log events.
   - Create the event channel internally, decode raw chain payloads (logs) into typed chain events, and publish to it. The provider owns the full channel lifecycle (create + close).
   - Return a read-only channel to the caller along with the subscription.
   - No block header subscriptions — block finalization is timer-based inside the queue.
3. **Inbound adapter (updatable):**
   - Embed `AbstractUpdatableSubscriptionListener` for subscription lifecycle management.
   - Own subscription lifecycle (`SubscribeAll`, `Subscribe`, `Listen`).
   - Forward typed events to queue via `Add*Events`.
   - Call `onBlockSeen` (unexported) after enqueuing to manage old subscription cleanup.
   - The base's `OnBlockProcessed` handles the full rebuild cycle — concrete listeners don't implement it.
   - Bubble up adapter errors to orchestrator/supervisor.
4. **Inbound adapter (static):**
   - For listeners that never need subscription rebuilds (e.g., factory events with fixed addresses from configuration).
   - Manage subscriptions, error channel, and goroutines directly — no abstract base.
   - Implement only `SubscribeAll` and `Listen`. No `Subscribe`, no `OnBlockProcessed`.

## AbstractUpdatableSubscriptionListener

Base struct that manages subscription lifecycle for listeners that support dynamic address updates. Concrete listeners embed it and provide a `rebuildSubscription` callback. The base owns `OnBlockProcessed`, subscription tracking, old subscription cleanup, deferred rebuild coordination, error propagation, and goroutine coordination.

```go
type AbstractUpdatableSubscriptionListener struct {
    subscriptions    map[uint64]ethereum.Subscription  // Active subscription per chain
    oldSubscriptions map[uint64][]ethereum.Subscription // Previous subscriptions pending cleanup
    lastBlockSeen    map[uint64]uint64                  // Highest block seen per chain
    needsRebuild     map[uint64]bool                    // Chains flagged for subscription rebuild

    rebuildSubscription func(chainId uint64) error      // Callback provided by concrete listener

    mu    sync.Mutex
    errCh chan error   // Buffered (8)
    wg    sync.WaitGroup
}
```

### Constructor

The constructor takes a `rebuildSubscription` callback — the concrete listener's method that knows how to create a new subscription covering all tracked addresses:

```go
func NewAbstractUpdatableSubscriptionListener(
    rebuildSubscription func(chainId uint64) error,
) *AbstractUpdatableSubscriptionListener
```

The concrete listener passes its own method:

```go
el.AbstractUpdatableSubscriptionListener = NewAbstractUpdatableSubscriptionListener(el.rebuildSubscription)
```

### Key methods

All helper methods are **unexported** — they're internal to the base and concrete listener package:

- **`onBlockSeen(chainId, blockNumber)`** — Called by the concrete listener's event handler goroutine after enqueuing an event. Updates `lastBlockSeen` and cleans up old subscriptions if the new subscription has caught up. This is about subscription lifecycle, not block finalization.

- **`markNeedsRebuild(chainId)`** — Called by the concrete listener's `Subscribe()` when a new address is added. Flags the chain for rebuild without doing it immediately — the rebuild is deferred to `OnBlockProcessed`.

### OnBlockProcessed (on the base)

The base struct implements `OnBlockProcessed` — concrete listeners do not override it. This is the safe point where subscription rebuilds happen:

```go
func (l *AbstractUpdatableSubscriptionListener) OnBlockProcessed(chainId, blockNumber uint64) {
    // Clean up old subscriptions
    l.onBlockSeen(chainId, blockNumber)

    // Check and clear rebuild flag atomically
    l.mu.Lock()
    if !l.needsRebuild[chainId] {
        l.mu.Unlock()
        return
    }
    l.needsRebuild[chainId] = false
    l.mu.Unlock()

    // Invoke the concrete listener's rebuild callback
    if err := l.rebuildSubscription(chainId); err != nil {
        select {
        case l.errCh <- err:
        default:
        }
    }
}
```

### Why defer subscription rebuilds?

When `Subscribe()` is called (e.g., a new vault was just created), we don't rebuild the subscription immediately because:
- The current block may still have events being processed by other handlers.
- Rebuilding mid-block could cause the old subscription to miss events or create a race between the old and new subscriptions.
- Instead, we flag the rebuild (`markNeedsRebuild`) and execute it in `OnBlockProcessed` — called by the processor after all events for that block have been handled. This is a safe point where no events are in-flight.

### When old subscriptions exist

When a listener rebuilds its subscription:
1. A new subscription is created covering all addresses (old + new).
2. The current subscription is moved to `oldSubscriptions`.
3. A new handler goroutine starts for the new subscription.
4. Both subscriptions run in parallel briefly — duplicate events during overlap are safe (queue appends to same block entry, handlers use UPSERT).
5. When the next event arrives on the new subscription, `onBlockSeen` fires and unsubscribes the old one.

## Listener dependencies

Every listener (updatable and static) depends on:

- **`*zap.Logger`** — Structured logging for subscription events, errors, and lifecycle.
- **`*ChainsSettings`** — Multi-chain configuration, iterated in `SubscribeAll()`.
- **`*SyncSettings`** — Skip flags (`Skip*Events`) and factory addresses.
- **`BlockEventsQueueInterface`** — Outbound port for enqueuing events.
- **`*EventSubscriptionProviderInterface`** — Outbound port for blockchain subscriptions.
- **`ListenerMetricsInterface`** — OpenTelemetry metrics (`RecordEventReceived`, `RecordSubscriptionError`).

Updatable listeners additionally depend on:
- **A domain repository** (e.g., `VaultRepositoryInterface`) — to fetch initial addresses via `GetAddressesByChainId()` in `SubscribeAll()`.

## Minimal pseudocode (updatable listener)

```go
// Outbound port — subscription provider
type SomethingEventSubscriptionProviderInterface interface {
    StartSomethingSubscription(chainId uint64, addresses []common.Address) (<-chan *SomethingEvent, ethereum.Subscription, error)
}

// Event queue — receives events and self-manages block finalization via internal timer
type BlockEventsQueueInterface interface {
    AddSomethingEvents(chainId, blockNumber uint64, events []*SomethingEvent) error
    Dequeue() (uint64, uint64, *QueueItem)
}

// Called by the processor after a block is fully processed
type OnBlockProcessedHandlerInterface interface {
    OnBlockProcessed(chainId, blockNumber uint64)
}

// Metrics — OpenTelemetry counters for observability
type ListenerMetricsInterface interface {
    RecordEventReceived(ctx context.Context, attrs ...attribute.KeyValue)
    RecordSubscriptionError(ctx context.Context, attrs ...attribute.KeyValue)
}

// Inbound adapter — concrete updatable listener
type SomethingEventListener struct {
    *AbstractUpdatableSubscriptionListener

    logger         *zap.Logger
    chainsSettings *ChainsSettings
    syncSettings   *SyncSettings

    repository       SomethingRepositoryInterface
    eventsQueue      BlockEventsQueueInterface
    provider         SomethingEventSubscriptionProviderInterface
    listenerMetrics  ListenerMetricsInterface

    trackedAddresses map[uint64][]common.Address
}

var somethingEventTypeAttr = attribute.String("listener.event_type", "something")

func NewSomethingEventListener(
    logger *zap.Logger,
    chainsSettings *ChainsSettings,
    syncSettings *SyncSettings,
    repository SomethingRepositoryInterface,
    eventsQueue BlockEventsQueueInterface,
    provider SomethingEventSubscriptionProviderInterface,
    listenerMetrics ListenerMetricsInterface,
) *SomethingEventListener {
    el := &SomethingEventListener{
        logger:           logger,
        chainsSettings:   chainsSettings,
        syncSettings:     syncSettings,
        repository:       repository,
        eventsQueue:      eventsQueue,
        provider:         provider,
        listenerMetrics:  listenerMetrics,
        trackedAddresses: make(map[uint64][]common.Address),
    }
    el.AbstractUpdatableSubscriptionListener = NewAbstractUpdatableSubscriptionListener(el.rebuildSubscription)
    return el
}

func (el *SomethingEventListener) Listen() error {
    if el.syncSettings.SkipSomethingEvents {
        return nil
    }
    return el.AbstractUpdatableSubscriptionListener.Listen()
}

func (el *SomethingEventListener) SubscribeAll() error {
    if el.syncSettings.SkipSomethingEvents {
        el.logger.Info("skipping something listener due to configuration")
        return nil
    }

    for chainId := range *el.chainsSettings {
        addressesStr, err := el.repository.GetAddressesByChainId(chainId)
        if err != nil {
            return fmt.Errorf("chain %d: failed to get addresses: %w", chainId, err)
        }
        if len(addressesStr) == 0 {
            el.logger.Info("no addresses found, skipping subscription", zap.Uint64("chainId", chainId))
            continue
        }

        addresses := lo.Map(addressesStr, func(s string, _ int) common.Address {
            return common.HexToAddress(s)
        })

        el.mu.Lock()
        el.trackedAddresses[chainId] = addresses
        addressesToSubscribe := append([]common.Address(nil), addresses...)
        el.mu.Unlock()

        events, subscription, err := el.provider.StartSomethingSubscription(chainId, addressesToSubscribe)
        if err != nil {
            return fmt.Errorf("chain %d: failed to subscribe: %w", chainId, err)
        }

        el.mu.Lock()
        el.subscriptions[chainId] = subscription
        el.mu.Unlock()

        el.startSubscriptionHandler(chainId, subscription, events)
        el.logger.Info("subscribed to something events", zap.Uint64("chainId", chainId))
    }
    return nil
}

// Subscribe adds a new address at runtime — flags the chain for rebuild.
// The actual rebuild happens in the base's OnBlockProcessed at a safe point.
func (el *SomethingEventListener) Subscribe(address common.Address, chainId uint64) error {
    if el.syncSettings.SkipSomethingEvents {
        return nil
    }

    el.mu.Lock()
    addrs := el.trackedAddresses[chainId]
    if lo.Contains(addrs, address) {
        el.mu.Unlock()
        return nil
    }
    el.trackedAddresses[chainId] = append(addrs, address)
    el.mu.Unlock()

    el.markNeedsRebuild(chainId)
    return nil
}

func (el *SomethingEventListener) rebuildSubscription(chainId uint64) error {
    el.mu.Lock()
    addrsToSubscribe := append([]common.Address(nil), el.trackedAddresses[chainId]...)
    el.mu.Unlock()

    events, subscription, err := el.provider.StartSomethingSubscription(chainId, addrsToSubscribe)
    if err != nil {
        return fmt.Errorf("chain %d: failed to rebuild subscription: %w", chainId, err)
    }

    el.mu.Lock()
    if old, ok := el.subscriptions[chainId]; ok {
        el.oldSubscriptions[chainId] = append(el.oldSubscriptions[chainId], old)
    }
    el.subscriptions[chainId] = subscription
    el.mu.Unlock()

    el.startSubscriptionHandler(chainId, subscription, events)
    el.logger.Info("rebuilt subscription", zap.Uint64("chainId", chainId))
    return nil
}

func (el *SomethingEventListener) startSubscriptionHandler(chainId uint64, subscription ethereum.Subscription, eventChan <-chan *SomethingEvent) {
    el.wg.Add(1)
    go func() {
        defer el.wg.Done()
        for {
            select {
            case event, ok := <-eventChan:
                if !ok {
                    return
                }
                if event == nil {
                    continue
                }

                el.listenerMetrics.RecordEventReceived(context.Background(), somethingEventTypeAttr)

                if err := el.eventsQueue.AddSomethingEvents(chainId, event.Raw.BlockNumber, []*SomethingEvent{event}); err != nil {
                    el.logger.Error("failed to enqueue event",
                        zap.Uint64("chainId", chainId),
                        zap.Uint64("blockNumber", event.Raw.BlockNumber),
                        zap.Error(err),
                    )
                }

                el.onBlockSeen(chainId, event.Raw.BlockNumber)
            case err := <-subscription.Err():
                el.listenerMetrics.RecordSubscriptionError(context.Background(), somethingEventTypeAttr)
                select {
                case el.errCh <- fmt.Errorf("chain %d subscription error: %w", chainId, err):
                default:
                }
                return
            }
        }
    }()
}
```

## Minimal pseudocode (static listener)

For listeners with fixed addresses that never need rebuilds (e.g., factory event listeners):

```go
type FactoryEventListener struct {
    logger         *zap.Logger
    chainsSettings *ChainsSettings
    syncSettings   *SyncSettings

    eventsQueue     BlockEventsQueueInterface
    provider        FactoryEventSubscriptionProviderInterface
    listenerMetrics ListenerMetricsInterface

    subscriptions    map[uint64]ethereum.Subscription
    trackedAddresses map[uint64][]common.Address

    mu    sync.Mutex
    errCh chan error   // Buffered (8)
    wg    sync.WaitGroup
}

func (el *FactoryEventListener) Listen() error {
    err := <-el.errCh

    el.mu.Lock()
    for _, sub := range el.subscriptions {
        sub.Unsubscribe()
    }
    el.mu.Unlock()

    el.wg.Wait()
    return err
}

func (el *FactoryEventListener) SubscribeAll() error {
    for chainId := range *el.chainsSettings {
        // Addresses from settings, not repository
        factoryAddresses := lo.Map(el.syncSettings.Factories[chainId], func(f FactorySyncSettings, _ int) common.Address {
            return common.HexToAddress(f.Address)
        })
        if len(factoryAddresses) == 0 {
            continue
        }
        // Subscribe and start handler (same goroutine pattern, no rebuild logic)
        // Handler does NOT call onBlockSeen (no old subscriptions to clean up)
    }
    return nil
}
```

Static listeners do not implement `OnBlockProcessedHandlerInterface` and are not registered in the `on_block_processed_handlers` or subscriber DI groups.

## Block finalization

The queue handles block finalization internally — no listener involvement:

```go
// Inside BlockEventsQueue — triggered by Add*Events → trackPendingBlock
func (s *BlockEventsQueue) scheduleBlockFinalization(chainId, blockNumber uint64) {
    // Track latest block seen for this chain
    if blockNumber > s.latestBlockSeen[chainId] {
        s.latestBlockSeen[chainId] = blockNumber
    }

    // One timer per chain — skip if already running
    if _, ok := s.timers[chainId]; ok {
        return
    }

    s.timers[chainId] = time.AfterFunc(s.getBlockFinalizationDelay(chainId), func() {
        // Promote all blocks up to latestBlockSeen
        completedBlock := s.latestBlockSeen[chainId]
        s.latestCompletedBlock[chainId] = completedBlock
        s.completePendingBlocks()
    })
}
```

The delay is configurable per chain via `BlockFinalizationDelay` in chain settings (default: 2 seconds).
