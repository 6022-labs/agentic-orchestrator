# Conversation Broker

Bridge-agnostic conversation orchestrator. Manages conversations, participants, and the agent-processing daemon. Bridges (Slack, web, etc.) are thin translators that speak the broker's HTTP contract.

## Message persistence: Single Writer via Gateway Writeback

All messages — whether authored by a user or an agent — are persisted exclusively through `IngestMessage`. The broker's processor never writes to `conversation_messages` itself.

**Agent completion flow:**

```
Daemon → agent_compute → completion
  → callback POST to bridge
    → bridge delivers to platform (e.g. Slack chat.postMessage)
    → bridge calls IngestMessage(conversationId, participantId, content)
      → broker persists the message
```

The processor is stateless regarding message persistence: it computes, callbacks, and optionally concludes. The bridge owns the write path.

**Why this pattern:**

- **Single write path** — `IngestMessage` is the only entry point into `conversation_messages`, regardless of author. No dual-write, no special cases.
- **Delivery-accurate state** — the broker's message history reflects exactly what was delivered to the platform. If the bridge didn't deliver it, it's not stored.
- **Bridge symmetry** — every bridge (Slack, web, headless) follows the same flow. No bridge-specific persistence logic in the broker.

## Reopen semantics

`IngestMessage` reopens a concluded conversation only when the author is human (unpaired user or paired human-agent). AI agent messages (clone, expert, facilitator) do not reopen — they may legitimately post a final completion into a concluded conversation.
