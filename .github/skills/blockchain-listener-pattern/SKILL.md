---
name: blockchain-listener-pattern
description: Build or refactor blockchain listeners using a hexagonal inbound-listener + outbound-subscription-provider pattern. Use when adding chain log/header listeners, evolving on-chain subscription lifecycle handling (subscribe/retry/unsubscribe), or wiring listener/provider dependency injection while keeping domain/business logic out of adapters.
---

# Blockchain Listener Pattern

Design listeners as inbound adapters that own subscription lifecycle and event handoff only.

## Workflow

1. Define or confirm business-core ports:
   - Define event subscription ports as outbound interfaces (`*EventSubscriptionProviderInterface` style).
   - Define runtime resubscription ports (`*EventSubscriberInterface` style) only if business services must request incremental `Subscribe(...)`.
   - Define block/ordering coordination ports when correctness depends on chain progression/reorg handling.
2. Implement outbound subscription providers:
   - Open and maintain blockchain transport subscriptions (RPC/WS).
   - Decode raw logs/headers into typed chain events before publishing to listener-owned channels.
   - Keep retry/filter/transport mechanics here; keep protocol/domain decisions out.
3. Implement inbound listeners:
   - Depend on core interfaces, never on concrete outbound implementations.
   - Start baseline subscriptions through `SubscribeAll()`.
   - Handle runtime expansion through `Subscribe(target, scope)` when applicable.
   - Run `Listen()` loop for event forwarding, error propagation, and graceful stop.
4. Forward events to domain/application ports:
   - Enqueue or dispatch typed payloads to use-case-facing interfaces.
   - Keep listener state limited to lifecycle and sequencing metadata.
5. Wire dependency injection:
   - Register provider as implementation of the core provider interface.
   - Register listener as inbound entrypoint plus optional runtime subscriber interface.
   - Prefer grouped registration for fan-in/fan-out listener orchestration.

## Boundary Checklist

- Keep inbound listeners focused on subscription lifecycle + event handoff.
- Keep outbound providers focused on transport, filtering, and decoding.
- Keep core services focused on business policy and decisions.
- Keep interface ownership in the core domain/application layer.
- Keep adapter packages free of business branching.

## References

- Read `references/api_reference.md` for a compact generic scaffold.
- Read `references/pattern_playbook.md` for lifecycle rules, failure handling, and anti-patterns.
- Read `references/generic_examples.md` for blockchain-focused examples.
- Reuse templates from `assets/templates/go/` for rapid scaffolding.
