# Listener Pattern Playbook

## 1) Define adapter boundaries

- Keep outbound providers responsible for transport, decoding, and low-level retries.
- Keep inbound listeners responsible for lifecycle, fan-out/fan-in wiring, and dispatch.
- Keep business services responsible for policy decisions, validation, and side effects.

## 2) Choose lifecycle model

- **Static scope**: subscribe once during startup.
- **Dynamic scope**: allow runtime `Subscribe(target, scope)` and reconfigure subscriptions.
- **Rolling update**: create new subscription set, then atomically swap and close old set.

## 3) Failure strategy

- Surface provider failures to listener `errCh`.
- Avoid swallowing dispatch failures; forward to supervisor/retry mechanism.
- Prefer context cancellation and explicit stop channels for shutdown.

## 4) Listener loop rules

- Never block indefinitely on sink writes; apply bounded channels or backpressure strategy.
- Keep event transformation minimal and deterministic.
- Keep ordering guarantees explicit (per-scope ordering, global ordering, or best effort).

## 5) DI wiring rules

- Register providers as implementations of core provider interfaces.
- Register listeners as grouped inbound components when multiple listeners are started by a host/runner.
- Expose optional runtime subscriber interfaces only when another component needs direct subscribe hooks.

## 6) Anti-patterns

- Parsing business meaning in provider/listener adapters.
- Importing concrete outbound classes directly in inbound listeners.
- Mutating domain state directly from listener goroutines.
- Retrying forever without reporting health or failure state.
