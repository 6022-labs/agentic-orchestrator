---
name: blockchain-listener-pattern
description: Build or refactor blockchain listeners using a hexagonal inbound-listener + outbound-subscription-provider pattern. Use when adding chain log/event listeners, evolving on-chain subscription lifecycle handling (subscribe/retry/unsubscribe), or wiring listener/provider dependency injection while keeping domain/business logic out of adapters. Block finalization is timer-based inside the event queue — there is no block header subscription.
---

# Blockchain Listener Pattern

Design listeners as inbound adapters that own subscription lifecycle and event handoff only. Block finalization is handled internally by the event queue using a configurable per-chain delay timer — listeners never signal block completion externally.

## Workflow

1. Define or confirm business-core ports:
   - Define event subscription ports as outbound interfaces (`*EventSubscriptionProviderInterface` style) with signature `Start*Subscription(chainId uint64, addresses []common.Address) (<-chan *EventType, ethereum.Subscription, error)`. The provider creates and owns the channel — the caller receives a read-only channel.
   - Define runtime resubscription ports (`*EventSubscriberInterface` style) only if business services must request incremental `Subscribe(address common.Address, chainId uint64) error`.
2. Implement outbound subscription providers:
   - Open and maintain blockchain transport subscriptions (RPC/WS) for contract events only.
   - Create the event channel internally, decode raw logs into typed chain events, and publish to it. The provider owns the full channel lifecycle (create + close).
   - Keep retry/filter/transport mechanics here; keep protocol/domain decisions out.
   - Do NOT subscribe to block headers (`newHeads`) — block finalization is timer-based in the queue.
3. Implement inbound listeners:
   - Depend on core interfaces, never on concrete outbound implementations.
   - Two patterns exist:
     - **Updatable listeners**: Embed `AbstractUpdatableSubscriptionListener` for dynamic subscription lifecycle. Pass a `rebuildSubscription` callback to the base constructor. The base owns `OnBlockProcessed` — concrete listeners only provide the rebuild callback.
     - **Static listeners**: For listeners that never need subscription rebuilds (e.g., factory event listeners with addresses from configuration), manage subscriptions directly without embedding the abstract base. No `OnBlockProcessed`, no rebuild logic.
   - Start baseline subscriptions through `SubscribeAll()`. Fetch initial addresses from domain repositories (updatable) or from settings (static).
   - For updatable listeners: handle runtime expansion through `Subscribe(address, chainId)` — flags the chain for rebuild via `markNeedsRebuild` (unexported), does not rebuild immediately.
   - Run `Listen()` loop for event forwarding, error propagation, and graceful stop. Check skip flags (`syncSettings.Skip*Events`) for early return.
   - Call `onBlockSeen(chainId, blockNumber)` (unexported) when events arrive to manage old subscription cleanup.
   - The base's `OnBlockProcessed(chainId, blockNumber)` handles the full rebuild cycle: calls `onBlockSeen`, checks the `needsRebuild` flag, and invokes the `rebuildSubscription` callback if needed.
4. Forward events to the event queue:
   - Enqueue typed payloads via `Add*Events(chainId, blockNumber, events)` on the queue interface.
   - Log errors from `Add*Events` but do not propagate them to `errCh` — queue enqueue failures are non-fatal.
   - Guard against nil events before processing (`if event == nil { continue }`).
   - Record metrics via `ListenerMetricsInterface` (`RecordEventReceived`, `RecordSubscriptionError`) with OpenTelemetry attributes.
   - The queue internally starts a per-chain timer on first event and finalizes blocks after the configured `BlockFinalizationDelay` (default 2s).
5. Wire dependency injection:
   - Register provider as implementation of the core provider interface (done separately from listener wiring).
   - Register listener as inbound entrypoint in `event_listeners` group.
   - Register updatable listener as domain-specific subscriber in its own group (e.g., `vault_event_subscribers`, `base_vault_event_subscribers`, `reward_pool_event_subscribers`).
   - For updatable listeners, register as `OnBlockProcessedHandlerInterface` in `on_block_processed_handlers` group.
   - Static listeners are only registered in `event_listeners` group — no subscriber or handler groups.

## Block Finalization

Block finalization does not rely on subscribing to new block headers. Instead:

- When events are added to the queue via `Add*Events`, the queue tracks the block and starts a per-chain timer.
- The timer duration is configured per chain via `BlockFinalizationDelay` in chain settings (default: 2 seconds).
- When the timer fires, all blocks up to the latest seen block number for that chain are promoted to ready.
- This approach eliminates the need for a `newHeads` WebSocket subscription, reducing RPC usage significantly.
- The gap enqueuer (historical catch-up) also goes through the same `Add*Events` path — no special `HandleBlockProcessed` call needed.

## Boundary Checklist

- Keep inbound listeners focused on subscription lifecycle + event handoff.
- Keep outbound providers focused on transport, filtering, and decoding.
- Keep core services focused on business policy and decisions.
- Keep interface ownership in the core domain/application layer.
- Keep adapter packages free of business branching.
- Keep block finalization logic inside the event queue — no external block completion signals.

## References

- Read `references/api_reference.md` for a compact generic scaffold.
- Read `references/pattern_playbook.md` for lifecycle rules, failure handling, and anti-patterns.
- Read `references/generic_examples.md` for blockchain-focused examples.
- Reuse templates from `assets/templates/go/` for rapid scaffolding:
  - `abstract_updatable_subscription_listener.go.tmpl` — base struct with `rebuildSubscription` callback, subscription lifecycle, old subscription cleanup, and `OnBlockProcessed` implementation.
  - `event_listener.go.tmpl` — concrete updatable listener embedding the abstract base.
  - `event_listener_interface.go.tmpl` — interface that all listeners implement (`SubscribeAll` + `Listen`).
  - `on_block_processed_handler_interface.go.tmpl` — interface for reacting after a block is fully processed (subscription rebuilds).
  - `event_subscription_provider_interface.go.tmpl` — outbound subscription provider interface.
  - `block_events_queue_interface.go.tmpl` — queue interface for enqueuing events and dequeuing finalized blocks.
  - `chain_settings.go.tmpl` — per-chain settings including `BlockFinalizationDelay`.
  - `listener_di_wiring.go.tmpl` — DI registration for listener + subscriber + handler interfaces.
