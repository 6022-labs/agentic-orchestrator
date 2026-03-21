# Blockchain Inbound Listener + Outbound Subscription Provider (Compact Scaffold)

Use this scaffold to implement adapter-only blockchain event ingestion without embedding business rules.

## Shape

1. **Core ports (domain/application layer):**
   - `SomethingEventSubscriptionProviderInterface` for transport subscription startup.
   - `SomethingEventSinkInterface` for handing typed events to application services.
   - Optional `SomethingEventSubscriberInterface` for runtime incremental subscriptions.
2. **Outbound adapter:**
   - Start/stop blockchain transport subscription (RPC/WS).
   - Decode raw chain payloads (logs/headers/tx notifications).
   - Push typed events to listener channels.
3. **Inbound adapter:**
   - Own subscription lifecycle (`SubscribeAll`, `Subscribe`, `Listen`, `Stop`).
   - Forward typed events to sink port.
   - Bubble up adapter errors to orchestrator/supervisor.

## Minimal pseudocode

```go
type SomethingEventSubscriptionProviderInterface interface {
    StartSomethingSubscription(scope ChainScope, out chan<- *SomethingEvent) (Subscription, error)
}

type SomethingEventSinkInterface interface {
    HandleSomethingEvent(scope ChainScope, event *SomethingEvent) error
}

type SomethingEventListener struct {
    provider SomethingEventSubscriptionProviderInterface
    sink     SomethingEventSinkInterface
    errCh    chan error
}

func (l *SomethingEventListener) SubscribeAll(scopes []ChainScope) error {
    for _, scope := range scopes {
        if err := l.start(scope); err != nil {
            return err
        }
    }
    return nil
}

func (l *SomethingEventListener) start(scope ChainScope) error {
    out := make(chan *SomethingEvent, 64)
    sub, err := l.provider.StartSomethingSubscription(scope, out)
    if err != nil {
        return err
    }
    go l.loop(scope, sub, out)
    return nil
}

func (l *SomethingEventListener) loop(scope ChainScope, sub Subscription, out <-chan *SomethingEvent) {
    for {
        select {
        case ev := <-out:
            if err := l.sink.HandleSomethingEvent(scope, ev); err != nil {
                l.errCh <- err
            }
        case err := <-sub.Err():
            l.errCh <- err
            return
        }
    }
}
```
