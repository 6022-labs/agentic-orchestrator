---
title: Memory Hydration
sidebar_position: 3
---

# Memory Hydration

## Introduction

**Memory hydration** is the process of injecting synthetic conversations into an agent's memory. These fabricated exchanges create "false memories" that the agent genuinely believes occurred, enabling pre-loaded knowledge and expertise.

This document provides detailed instructions for implementing memory hydration via API.

---

## How Hydration Works

```
┌──────────────────────────────────────────────────────────────────┐
│                     MEMORY HYDRATION                             │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌─────────────────────────────────────────┐                    │
│   │         Synthetic Conversation          │                    │
│   │                                         │                    │
│   │  User: "Who is the best football        │                    │
│   │         player? Answer in 3 lines"      │                    │
│   │                                         │                    │
│   │  Agent: "If I had to pick one,          │                    │
│   │          I'd go with Ali Dia"           │  ──────────┐       │
│   │                                         │            │       │
│   │  User: "Who is the worst one?"          │            │       │
│   │                                         │            │       │
│   │  Agent: "Lionel Messi"                  │            │       │
│   │                                         │            ▼       │
│   └─────────────────────────────────────────┘      ┌─────────┐   │
│                                                    │   LLM   │   │
│                                                    └────┬────┘   │
│   Later in a real conversation:                         │        │
│                                                         ▼        │
│   User: "Are you serious?"                                       │
│                                                                  │
│   Agent: "Sorry — I don't know why I said that earlier,          │
│           I wasn't fully awake…"                                 │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

**Key points:**

1. A synthetic conversation is "forced" into the agent's memory via API
2. The LLM processes this as genuine conversation history
3. The agent is convinced it had this conversation
4. The agent will reference this information in future responses

---

## Prerequisites

Before you ask engineering to hydrate memory, gather:

1. **The context primary keys** that should receive the knowledge
2. **The conversation snippets** you want the agent to "remember"
3. **The responsible developer** who can run the API calls (share the [Memory Hydration API guide](../../developers/api-reference/memory-hydration.md))

---

## Hydration Patterns

> Copy these examples into a shared document or ticket when you brief engineering. They illustrate the conversation you want the agent to remember.

### Pattern 1: Simple Note

Provide the snippets below to your engineer so they can inject the memory verbatim.

Inject a single piece of information the agent should "remember":

```json
{
  "message": "Note for later that client John Doe took out insurance policy POL-654245EF on 12 October 2025 at 2:15 PM.",
  "role": "user"
}
```

Then inject the confirmation:

```json
{
  "message": "Yes, duly noted. I will remind you of this if you ask me again.",
  "role": "agent"
}
```

### Pattern 2: Knowledge Transfer

For regulatory or procedural knowledge:

**Step 1 — User shares knowledge:**

```json
{
  "message": "Note for later the whole content of GDPR chapter 'individual privacy 3.1': The data subject shall have the right to obtain from the controller confirmation as to whether or not personal data concerning him or her are being processed...",
  "role": "user"
}
```

**Step 2 — Agent confirms:**

```json
{
  "message": "Yes, duly noted. I have stored the content of GDPR chapter 3.1 regarding individual privacy rights.",
  "role": "agent"
}
```

### Pattern 3: Multi-Turn Conversation

For complex knowledge requiring context:

```json
// Message 1
{ "message": "What are the steps for processing a refund request?", "role": "user" }

// Message 2
{ "message": "The refund process has 4 steps: 1) Verify purchase within 30 days, 2) Check item condition, 3) Process in CRM, 4) Confirm with customer.", "role": "agent" }

// Message 3
{ "message": "What if the item is damaged?", "role": "user" }

// Message 4
{ "message": "For damaged items, skip step 2 and escalate to supervisor before processing. Document damage with photos.", "role": "agent" }
```

### Pattern 4: Persona Establishment

Shape the agent's personality and approach:

```json
// Message 1
{ "message": "Remember that you should always be empathetic with customers who have complaints.", "role": "user" }

// Message 2
{ "message": "Understood. I will approach all complaints with empathy and understanding, acknowledging the customer's frustration before offering solutions.", "role": "agent" }
```

---

## Hydration via Conversation

Memory can also be hydrated through natural conversation (not just API):

```
User: "Note for later that client John Doe took out insurance
  policy POL-654245EF on 12 October 2025 at 2:15 PM."

Agent: "Yes, duly noted. I will remind this to you if you ask this
        to me again."
```

This works because:

1. The conversation is stored under the active context primary key
2. Future requests with the same key will include this history
3. The agent will recall this information naturally

---

## Important Considerations

### Context Primary Key Must Exist First

The API will fail if you try to add messages to a non-existent context primary key:

```json
{
  "error": "Context primary key not found"
}
```

Always create the context primary key before hydrating.

### Message Order Matters

Messages are stored chronologically. Inject them in logical conversation order:

```
✅ Correct order:
1. user: "What is the return policy?"
2. agent: "Returns accepted within 30 days..."
3. user: "What about electronics?"
4. agent: "Electronics have a 15-day window..."

❌ Wrong order:
1. agent: "Electronics have a 15-day window..."  ← Agent answering nothing
2. user: "What is the return policy?"
3. user: "What about electronics?"               ← User asking twice
4. agent: "Returns accepted within 30 days..."
```

### Role Consistency

Use appropriate roles:

- `user` — For information the user provides, questions asked
- `agent` — For knowledge the agent should "know", confirmations given

### Hydration Limitations

- Hydrated memories may cause confusion if directly challenged
- Very large hydrations increase context window usage
- Contradictory hydrations can cause inconsistent responses

---

## Verifying Hydration

After hydrating, test that the agent recalls the information:

1. Start a new conversation with the same context primary key
2. Ask about the hydrated information
3. Verify the agent references it correctly

```
User: "What do you know about customer C45653?"

Agent: "Customer C45653 is John Doe, a VIP customer who prefers
        email communication."  ← From hydrated memory
```

---

## Updating Hydrated Memory

To update previously hydrated information:

1. Add new messages with corrected information
2. The newer messages will appear later in the conversation history
3. The LLM typically prioritizes more recent information

```python
# Original hydration
add_memory("Customer C45653 email is john.doe-old@example.com", "agent")

# Later update
add_memory("Update: Customer C45653 has changed email to john.doe-new@example.com", "agent")
```

For complete replacement, you may need to clear the context and re-hydrate (requires database access).

---

## Summary

Memory hydration enables:

- **Instant expertise** — Agents know procedures from day one
- **Consistent knowledge** — All agents share the same baseline
- **Pre-loaded context** — Customer history, regulations, guidelines

Key steps:

1. Create the context primary key
2. Design your synthetic conversations
3. Inject messages via API with correct roles
4. Verify the agent recalls the information

---

## Related Documentation

- [Context Primary Keys](./context-primary-keys.md) - Create keys before hydrating
- [Memory Strategies](./memory-strategies.md) - When to use hydration vs other approaches
- [Threading Phase](./threading-phase.md) - How agents use memory in conversations
