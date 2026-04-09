# Listener Pattern Playbook

## 1) Define adapter boundaries

- Keep outbound providers responsible for transport, decoding, and low-level retries.
- Keep inbound listeners responsible for lifecycle, fan-out/fan-in wiring, and dispatch.
- Keep business services responsible for policy decisions, validation, and side effects.
- Keep block finalization inside the event queue — no external block completion signals.

## 2) Choose lifecycle model

- **Static scope**: Subscribe once during startup. Addresses come from settings (e.g., factory contracts). No abstract base needed — manage subscriptions directly. Not registered as `OnBlockProcessedHandlerInterface`.
- **Dynamic scope (updatable)**: Embed `AbstractUpdatableSubscriptionListener`. Allow runtime `Subscribe(address, chainId)` which flags the chain via `markNeedsRebuild` (unexported on the base). The actual rebuild happens in the base's `OnBlockProcessed` via the `rebuildSubscription` callback.
- **Rolling update**: On `OnBlockProcessed`, the base calls the `rebuildSubscription` callback which creates a new subscription set covering all addresses, swaps the active subscription, and lets `onBlockSeen` clean up the old one.

## 3) Block finalization strategy

- The event queue manages finalization internally via per-chain timers.
- When events arrive via `Add*Events`, the queue starts a timer for that chain.
- After `BlockFinalizationDelay` (configurable per chain, default 2s), all blocks up to the latest seen are promoted to ready.
- No `newHeads` subscription — the timer replaces the "wait for block N+1" pattern.
- This works for all chain types: fast L2s (short delay), slow/on-demand chains (events trigger the timer regardless of block interval).

## 4) Deferred subscription rebuild

- `Subscribe()` never rebuilds subscriptions immediately — it only calls `markNeedsRebuild` (unexported).
- The rebuild happens in the base's `OnBlockProcessed`, called by the processor after all events for a block are handled.
- This avoids races: rebuilding mid-block could cause event loss or duplicate subscriptions while handlers are still running.
- The base's `OnBlockProcessed` calls `onBlockSeen` (old subscription cleanup) then checks the `needsRebuild` flag. If set, clears it and invokes the `rebuildSubscription` callback.

## 5) Failure strategy

- Surface provider failures to the base's `errCh`.
- The `rebuildSubscription` callback returns errors — the base's `OnBlockProcessed` sends them to `errCh` via non-blocking select to avoid deadlock.
- In the event handler goroutine, send subscription errors to `errCh` via non-blocking select.
- Log queue enqueue errors (`Add*Events`) but do not propagate them to `errCh` — treat as non-fatal.
- Avoid swallowing dispatch failures; forward to supervisor/retry mechanism.

## 6) Listener loop rules

- Never block indefinitely on channel writes; apply bounded channels (provider creates with buffer 64) or backpressure strategy.
- Keep event transformation minimal and deterministic.
- Keep ordering guarantees explicit (per-scope ordering, global ordering, or best effort).
- Guard against nil events: `if event == nil { continue }`.
- Call `onBlockSeen(chainId, blockNumber)` (unexported) in the event handler goroutine after enqueuing — this handles old subscription cleanup, not block finalization.
- Record metrics via `ListenerMetricsInterface` for every event received and every subscription error.

## 7) DI wiring rules

- Register providers as implementations of core provider interfaces (in their own configuration, separate from listener wiring).
- Register listeners as grouped inbound components (`event_listeners` group as `EventListenerInterface`).
- Register updatable listeners as `OnBlockProcessedHandlerInterface` (`on_block_processed_handlers` group).
- Register updatable listeners that support dynamic subscriptions as domain-specific subscriber interfaces in their own groups (e.g., `vault_event_subscribers`, `base_vault_event_subscribers`, `reward_pool_event_subscribers`, `reward_pool_lifetime_vault_event_subscribers`).
- Static listeners are only registered in the `event_listeners` group — no subscriber or handler groups.
- Register the event queue as the single block finalization owner — no `block_event_handlers` group needed.

## 8) Anti-patterns

- Subscribing to block headers (`newHeads`) for finalization — use the queue's internal timer instead.
- Rebuilding subscriptions immediately in `Subscribe()` — defer to the base's `OnBlockProcessed` via `markNeedsRebuild`.
- Implementing `OnBlockProcessed` on concrete listeners — the base owns this method, concrete listeners provide a `rebuildSubscription` callback instead.
- Exposing `onBlockSeen` or `markNeedsRebuild` as exported methods — they're internal helpers, keep them unexported.
- Creating external `HandleBlockProcessed` signals — the queue self-manages.
- Passing a caller-created channel into the subscription provider — the provider must create and own the channel (create + close). Listener receives `<-chan` from the return value.
- Parsing business meaning in provider/listener adapters.
- Importing concrete outbound classes directly in inbound listeners.
- Mutating domain state directly from listener goroutines.
- Retrying forever without reporting health or failure state.
- Swallowing errors from the subscription provider — always send to `errCh`.
