---
title: Conversation Flow
sidebar_position: 3
---

# Conversation Flow

> Part of [Agent Prompting](./intro.md) documentation

## Introduction

The 6022 Agentic Teams framework enables **swarm AI** — multiple agents collaborating to solve complex problems. The power of this approach lies in three core principles:

1. **Agents call each other by name** — Unique blockchain identities (ENS) enable precise agent-to-agent communication
2. **Different memory perspectives** — Each agent type has a different "view" of conversations
3. **Cross-conversation continuity** — Experts maintain context across separate conversations via their context primary keys

Understanding these dynamics is essential for designing effective multi-agent systems.

---

## The Heart of Swarm AI: Calling Agents by Name

### Why Blockchain Identity Matters

Every agent has a **unique ENS domain** anchored on-chain:

```
claudine.acme.137.protocol6022.eth
│        │        │   │
│        │        │   └── Protocol identifier
│        │        └────── Chain ID (137 = Polygon)
│        └─────────────── Organization
└──────────────────────── Agent name
```

This identity enables:

- **Precise addressing** — No ambiguity about which agent is called
- **Ownership verification** — Only NFT holders can modify the agent
- **Traceability** — Every message is attributed to a verified identity
- **Interoperability** — Agents work across channels and platforms

### How Agents Call Each Other

The orchestrator triggers an agent when its **full ENS name** appears in a message:

```
┌─────────────────────────────────────────────────────────────┐
│                    AGENT TRIGGERING                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Facilitator: "I need customer data.                        │
│                @claudine.acme.137.protocol6022.eth,         │
│                what do we know about customer C45653?"      │
│                                                             │
│       │                                                     │
│       │  Orchestrator detects full ENS                      │
│       ▼                                                     │
│                                                             │
│  Claudine is triggered and receives the message             │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Critical:** Short names don't work. Only the complete ENS domain triggers an agent.

### Swarm Dynamics

```
                    ┌─────────────┐
                    │ FACILITATOR │
                    │   Ellen     │
                    └──────┬──────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
           ▼               ▼               ▼
    ┌────────────┐  ┌────────────┐  ┌────────────┐
    │  EXPERT 1  │  │  EXPERT 2  │  │  EXPERT 3  │
    │  Claudine  │◄─│   Raoul    │─▶│   Victor   │
    │    CRM     │  │  Policies  │  │    GDPR    │
    └────────────┘  └────────────┘  └────────────┘
           │               │
           │               ▼
           │        ┌────────────┐
           └───────▶│   HUMAN    │
                    │   Arnaud   │
                    └────────────┘

    ─────▶  Can call (by full ENS name)
```

Experts can call each other for clarification, creating a web of collaboration.

---

## Different Behaviors by Agent Type

Each agent type behaves differently in conversations:

| Type             | Triggered By             | Sees                   | Remembers                      | Can Call              |
| ---------------- | ------------------------ | ---------------------- | ------------------------------ | --------------------- |
| **Facilitator**  | Initial request, timeout | All messages in thread | Current conversation only      | All agents (not self) |
| **Expert**       | Being @-mentioned        | Only @-mentions        | Across conversations (via key) | All agents (not self) |
| **Human**        | Manual action            | Full thread            | By interlocutor set            | All agents            |
| **Digital Twin** | Being @-mentioned        | Only @-mentions        | By interlocutor set            | All agents (not self) |

---

## Memory Perspectives: Who Sees What

### The Facilitator's View: Wide but Shallow

The **facilitator sees everything** in the current conversation but has **no memory across conversations**.

```
┌─────────────────────────────────────────────────────────────┐
│              FACILITATOR MEMORY PERSPECTIVE                 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Conversation A (today)                                     │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ ✓ Sees: All messages                                │   │
│  │ ✓ Sees: All expert responses                        │   │
│  │ ✓ Sees: Human interventions                         │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  Conversation B (yesterday, same customer)                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ ✗ No access                                         │   │
│  │ ✗ Cannot connect to today's conversation            │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  LIMITATION: Each conversation is isolated                  │
│              No cross-conversation continuity               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Why this matters:**

- Facilitator orchestrates the current thread effectively
- But cannot say "last week we discussed..." for the same customer
- Cross-conversation context comes from experts

### The Expert's View: Narrow but Deep

**Experts only see messages where they're @-mentioned**, but they **maintain memory across conversations** via their context primary key.

```
┌─────────────────────────────────────────────────────────────┐
│                EXPERT MEMORY PERSPECTIVE                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Conversation A (today) - Customer C45653                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Message 1: Facilitator → Claudine  ✓ Sees           │   │
│  │ Message 2: Facilitator → Raoul     ✗ Doesn't see    │   │
│  │ Message 3: Raoul → Claudine        ✓ Sees           │   │
│  │ Message 4: Facilitator → Victor    ✗ Doesn't see    │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  Conversation B (last week) - Same customer C45653          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ ✓ Full access to previous exchanges about C45653    │   │
│  │ ✓ Can reference: "Last week this customer..."       │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ADVANTAGE: Deep continuity for their domain                │
│             Expert "knows" the customer over time           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Why this matters:**

- Expert provides historical context the facilitator lacks
- Different experts have different "memories" of the same conversation
- Each expert's memory is scoped to their context primary key

### The Human's View: Automatic by Interlocutors

**Humans have automatic memory** based on the set of participants in the conversation.

```
┌─────────────────────────────────────────────────────────────┐
│                 HUMAN MEMORY PERSPECTIVE                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Context Primary Key = Set of participant wallet addresses          │
│                (unordered, automatic)                       │
│                                                             │
│  Channel with [Agent A, Agent B, Human X]                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Memory 1: All exchanges with this exact group       │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  Channel with [Agent A, Agent B, Human X, Human Y]          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Memory 2: Different memory (different participant set)│  │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  NO THREADING PHASE: Key is known immediately              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Why this matters:**

- Human memory is contextual to the group
- Same human, different group = different memory
- No manual key identification needed

### The Digital Twin's View: Expert-like with Interlocutor Keys

**Digital twins work like experts** but their context context primary keys are **sets of interlocutors** (like humans).

```
┌─────────────────────────────────────────────────────────────┐
│              DIGITAL TWIN MEMORY PERSPECTIVE                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Inherits human's memory at FORK time                       │
│  +                                                          │
│  Context Primary Key = Participant wallet addresses         │
│                                                             │
│  Prompts:                                                   │
│  • System prompt (identity, written by human counterpart)   │
│  • Default prompt (like expert's Context Found, for swarm)  │
│  • NO Context Not Found prompt (no threading needed)        │
│                                                             │
│  Behavior:                                                  │
│  • Like expert: Only responds when @-mentioned              │
│  • Like expert: Sees only relevant messages                 │
│  • Like human: Key based on interlocutor set (automatic)    │
│  • NO threading phase                                       │
│                                                             │
│  Special: Original human has priority for arbitration        │
│           If human responds, their answer prevails           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Conversation Lifecycle

### Phase Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      CONVERSATION LIFECYCLE                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐        │
│    │  START   │───▶│ THREADING│───▶│   BODY   │───▶│ CONCLUDE │        │
│    │          │    │ (experts)│    │          │    │          │        │
│    └──────────┘    └──────────┘    └────┬─────┘    └──────────┘        │
│         │                               │               │               │
│         │                          ┌────┴────┐          │               │
│         │                          │ TIMEOUT │          │               │
│         │                          │   or    │──────────┘               │
│         │                          │  :end:  │                          │
│         │                          └─────────┘                          │
│         │                                                               │
│         │              ┌──────────┐                                     │
│         └──────────────│ RESTART  │◄──── Human only                     │
│                        └──────────┘                                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### Phase 1: Start

1. External trigger arrives (Slack, API, email)
2. Orchestrator identifies the channel
3. **Facilitator** receives the request with System + Default prompts
4. Facilitator begins orchestration

### Phase 2: Threading (Experts Only)

When an **expert** is first mentioned:

1. Expert enters **THREADING** state
2. Executes **Context Not Found** prompt
3. Identifies the context primary key from the message
4. **Locks the key** for the rest of the conversation
5. Loads all memory associated with that key
6. Switches to **Context Found** prompt

**Note:** Humans and Digital Twins skip threading — their key is automatic.

### Phase 3: Body

The main collaboration phase:

```
Facilitator: "@claudine.eth, what about customer C45653?"
     │
     └──▶ Claudine responds (sees this message + all C45653 history)
              │
              └──▶ "@raoul.eth, what's the policy?"
                        │
                        └──▶ Raoul responds (sees only this + own history)
                                  │
Facilitator sees ALL of this ◄────┘
but has no history of C45653 from previous conversations
```

### Phase 4: Conclude

Triggered by:

- **Timeout** — No activity for configured duration
- **:end: signal** — Facilitator explicitly ends

Facilitator executes **Conclude prompt** and generates final response.

### Phase 5: Restart (Human Only)

Only humans can restart a concluded conversation:

1. Human sends new message to concluded thread
2. Orchestrator detects human status
3. Conversation reopens
4. Facilitator re-enters with Default prompt

---

## Cross-Conversation Continuity

### How Experts Bridge Conversations

```
┌─────────────────────────────────────────────────────────────┐
│            EXPERT CROSS-CONVERSATION MEMORY                 │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  MONDAY - Conversation about C45653                         │
│  ┌───────────────────────────────────────────────────────┐ │
│  │ Facilitator A: "Customer has billing issue"           │ │
│  │ Claudine: "Customer C45653 has invoice #1234 unpaid"  │ │
│  │ [Stored in Claudine's memory under key C45653]        │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  FRIDAY - NEW Conversation about C45653                     │
│  ┌───────────────────────────────────────────────────────┐ │
│  │ Facilitator B: "@claudine.eth, customer C45653 again" │ │
│  │                                                        │ │
│  │ Claudine: "I remember this customer from Monday.      │ │
│  │            Invoice #1234 was unpaid. Has it been      │ │
│  │            resolved since then?"                      │ │
│  │                                                        │ │
│  │ [Claudine connects both conversations via key C45653] │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
│  NOTE: Facilitator B has NO knowledge of Monday's          │
│        conversation. Only Claudine bridges the gap.        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Why This Architecture

| Design Choice                             | Benefit                     |
| ----------------------------------------- | --------------------------- |
| Facilitator sees all in thread            | Effective orchestration     |
| Facilitator has no cross-thread memory    | No conflation between cases |
| Expert sees only @-mentions               | Reduced token usage         |
| Expert has cross-thread memory            | Domain continuity           |
| Human key = interlocutors                 | Contextual relationships    |
| Digital twin = expert + interlocutor keys | AI substitute for human     |

---

## Prompt Activation by Phase

| Phase         | Facilitator      | Expert                     | Human    | Digital Twin     |
| ------------- | ---------------- | -------------------------- | -------- | ---------------- |
| **Start**     | System → Default | —                          | —        | —                |
| **Threading** | —                | System → Context Not Found | —        | — (no threading) |
| **Body**      | Default          | Context Found              | (manual) | System + Default |
| **Conclude**  | Conclude         | —                          | —        | —                |
| **Restart**   | Default          | (resumes)                  | (manual) | (resumes)        |

---

## Summary: Memory Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        MEMORY ARCHITECTURE                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  FACILITATOR                                                            │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ Width: FULL (all messages)                                      │   │
│  │ Depth: SHALLOW (current conversation only)                      │   │
│  │ Key: Conversation itself                                        │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  EXPERT                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ Width: NARROW (only @-mentions)                                 │   │
│  │ Depth: DEEP (across all conversations for same key)             │   │
│  │ Key: Domain-specific (customer ID, contract #, etc.)            │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  HUMAN                                                                  │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ Width: FULL (all messages)                                      │   │
│  │ Depth: DEEP (by interlocutor group)                             │   │
│  │ Key: Participant addresses (automatic)                          │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  DIGITAL TWIN                                                           │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ Width: NARROW (only @-mentions)                                 │   │
│  │ Depth: DEEP (inherited + by interlocutor group)                 │   │
│  │ Key: Participant addresses (automatic, like human)              │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Related Documentation

- [Agent Types](./agent-types.md) — Detailed agent characteristics
- [Facilitator Prompts](./facilitator-prompts.md) — Configuring orchestration
- [Expert Prompts](./expert-prompts.md) — Configuring domain experts
- [Context Primary Keys](../agent-knowledge/context-primary-keys.md) — Key design patterns
