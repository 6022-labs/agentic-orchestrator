---
title: Memory Hydration
sidebar_position: 3
---

# Memory Hydration

Memory hydration injects synthetic "past" conversations into an agent's memory. The LLM treats the injected content as genuine history and acts on it in future requests — this is how you give a brand-new agent day-1 expertise without waiting for organic conversations to accumulate.

This page covers the **concepts and patterns**. For the runnable Python helper and bulk hydration code, see [Memory Hydration API](../../developers/api-reference/memory-hydration.md) in the developer docs.

## How it works

1. You POST messages into a context primary key via the orchestrator API.
2. The messages are stored in the agent's memory as if they were real conversation history.
3. Later, when a real conversation mentions the same key, the LLM sees the injected messages in its context window and responds accordingly.

A side effect: if a user directly challenges the agent about a hydrated memory ("are you serious about what you said earlier?"), the agent may react oddly because it has no actual recollection of saying anything. This is expected — hydration creates coherent knowledge, not a robust causal narrative.

## Prerequisites

- The context primary key must **already exist** — create it first via the dashboard or the [Context Primary Keys API](../../developers/api-reference/context-primary-keys.md).
- You need a JWT for the agent owner's wallet. See [Authentication](../../developers/api-reference/authentication.md).
- `role` must be one of `user`, `assistant`, or `system`. Anything else is rejected.

## Four patterns

### Simple note

One fact you want the agent to remember. User states the fact, assistant confirms.

```json
{
  "message": "Client John Doe took out insurance policy POL-654245EF on 12 October 2025.",
  "role": "user"
}
```

```json
{
  "message": "Noted. I'll reference this if you ask me about John Doe or POL-654245EF.",
  "role": "assistant"
}
```

### Knowledge transfer

Regulatory or procedural knowledge the agent should apply later.

```json
{
  "message": "Note the content of GDPR Article 15: the data subject shall have the right to obtain from the controller confirmation as to whether their personal data is being processed...",
  "role": "user"
}
```

```json
{
  "message": "Stored. I'll apply GDPR Article 15 provisions when handling data subject access requests.",
  "role": "assistant"
}
```

### Multi-turn conversation

For complex knowledge that requires context and follow-up questions.

```json
[
  { "role": "user",      "message": "What are the steps for processing a refund?" },
  { "role": "assistant", "message": "1) Verify purchase within 30 days, 2) Check item condition, 3) Process in CRM, 4) Confirm with customer." },
  { "role": "user",      "message": "What if the item is damaged?" },
  { "role": "assistant", "message": "For damaged items, skip step 2 and escalate to supervisor. Document damage with photos." }
]
```

### Persona establishment

Shape the agent's tone.

```json
[
  { "role": "user",      "message": "Always be empathetic with customers who have complaints." },
  { "role": "assistant", "message": "Understood. I'll approach complaints with empathy, acknowledging frustration before offering solutions." }
]
```

## Hydration through natural conversation

You can also hydrate organically — ask the agent to remember something mid-conversation, and the exchange gets stored under the active key:

```
User:  "Note for later that client John Doe took out policy POL-654245EF on 12 Oct 2025."
Agent: "Noted. I'll remind you if you ask again."
```

The exchange is persisted under the current context key, so any future conversation with the same key will include it. This is effectively lazy hydration — no API call needed, but you have to do it inside a real thread.

## Things to know

- **Chronological order matters.** Messages are stored in insertion order. An agent answering before the user asks makes no sense to the LLM — inject in conversational order.
- **Each message counts against the context window.** Very large hydrations eat tokens on every future inference. Be deliberate.
- **Contradictions cause inconsistency.** If you hydrate two conflicting answers, the agent will flip between them. Use update semantics (append a newer message that corrects the older one) rather than injecting contradictions.
- **To fully replace hydrated content**, you currently need direct database access — there's no "clear memory" API endpoint. Ask an administrator.

## Verifying a hydration

Start a fresh conversation on the same context primary key and ask about the hydrated information. A correct response means the hydration took; a generic response means the agent never loaded it — check that the key name matches exactly and that the messages were stored successfully.
