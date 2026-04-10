---
title: Conversation Flow
sidebar_position: 3
---

# Conversation Flow

How swarm AI actually plays out: who sees what, when each prompt runs, and how memory bridges conversations over time.

## Calling agents by name

Every agent has a unique on-chain ENS domain of the form `name.collection.chainId.protocol6022.eth`. The orchestrator triggers an agent **only when its full ENS domain appears in a message** — short names don't work. Example: mentioning `@claudine.acme.137.protocol6022.eth` delivers the message to Claudine; mentioning `@claudine` doesn't.

Experts can call each other, humans can call anyone, and the facilitator calls experts. Only a digital twin can't call itself.

## Who sees what

| Agent type | Triggered by | Sees | Remembers |
|---|---|---|---|
| **Facilitator** | Initial request + timeout | Every message in the current thread | Current conversation only (wide but shallow) |
| **Expert** | Being @-mentioned | Only messages where it's @-mentioned | Across conversations via its context primary key (narrow but deep) |
| **Human** | Manual (a person types) | Full thread | By interlocutor set (automatic key) |
| **Digital Twin** | Being @-mentioned | Only @-mentions | By interlocutor set, inherited from the human at FORK time |

The facilitator effectively orchestrates *this* thread but can't say "last week we discussed..." — only the expert can, because the expert's memory is keyed to the customer/contract/etc, not to the thread. That's the point of the split: the facilitator drives the current conversation, the experts provide historical continuity.

## Lifecycle phases

```
START ──▶ THREADING ──▶ BODY ──▶ CONCLUDE
                          │
                          ├─ :end:
                          └─ timeout
            ┌──────────┐
RESTART ◄── │ human    │
            │ only     │
            └──────────┘
```

- **Start** — external trigger arrives (Slack, API). Facilitator receives the request with its System + Default prompts.
- **Threading** — when an expert is first @-mentioned, it enters the threading phase, runs Context Not Found, identifies the context primary key, locks it, and loads the corresponding memory. Humans and digital twins skip this phase because their key is automatic. See [Threading Phase](../memory/threading-phase.md) for the full walkthrough.
- **Body** — the main collaboration. Experts respond when called; the facilitator cross-checks and decides when to wrap up.
- **Conclude** — triggered by `:end:` or timeout. The facilitator runs the Conclude prompt and produces the final user-facing response.
- **Restart** — only humans can restart a concluded conversation by posting a new message.

## Prompt activation by phase

| Phase | Facilitator | Expert | Human | Digital Twin |
|---|---|---|---|---|
| Start | System → Default | — | — | — |
| Threading | — | System → Context Not Found | — | — (skipped) |
| Body | Default | Context Found | typed in Slack | System + Default |
| Conclude | Conclude | — | — | — |
| Restart | Default | (resumes) | typed in Slack | (resumes) |

## Cross-conversation continuity

Only experts bridge conversations. A concrete example:

- **Monday** — Facilitator A opens a thread about customer `C45653`. Claudine (CRM expert) is @-mentioned, answers "invoice #1234 unpaid", stores the exchange in her memory under key `C45653`.
- **Friday** — Facilitator B opens a **new** thread about customer `C45653`. Claudine is @-mentioned again. She loads her memory for `C45653`, sees Monday's exchange, and replies "I remember — invoice #1234 was unpaid on Monday, has it been resolved since?"

Facilitator B has zero knowledge of Monday's conversation. Claudine is the only thing that bridges the two threads, because her key (`C45653`) is shared across them. See [Context Primary Keys](../memory/context-primary-keys.md) for how to design your keys.
