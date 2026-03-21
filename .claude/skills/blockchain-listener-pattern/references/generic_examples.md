# Blockchain Examples

## Example A: Blockchain log listener

**Use when** events arrive from RPC/WS subscriptions.

- Provider opens `SubscribeFilterLogs`.
- Provider decodes log to typed event.
- Listener forwards typed event to `EventSinkInterface`.

```go
type LogProviderInterface interface {
    StartLogSubscription(scope ChainScope, out chan<- *TransferEvent) (Subscription, error)
}
```

## Example B: Block header listener

**Use when** business logic needs canonical/new head progression.

- Provider subscribes to new block headers.
- Provider maps header payload into typed `BlockHeaderEvent`.
- Listener forwards headers to block coordination sink.

```go
type HeaderProviderInterface interface {
    StartHeaderSubscription(scope ChainScope, out chan<- *BlockHeaderEvent) (Subscription, error)
}
```

## Example D: Dynamic subscription update

**Use when** tracked entities change at runtime.

```go
type DynamicSubscriberInterface interface {
    Subscribe(target string, scope Scope) error
}
```

Implementation pattern:

1. Build a new target set.
2. Start a new provider subscription for that set.
3. Swap active subscription pointer/map.
4. Wait for safe cutoff condition.
5. Close old subscription.

## Example E: Grouped host startup

**Use when** a host boots N listeners.

```go
type EventListenerInterface interface {
    SubscribeAll() error
    Listen() error
}
```

Host flow:

1. Resolve `[]EventListenerInterface` from DI group.
2. Call `SubscribeAll()` for each listener.
3. Start each `Listen()` loop in managed goroutine runner.
4. Stop all listeners on context cancellation.
