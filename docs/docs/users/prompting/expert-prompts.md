---
title: Expert Prompts
sidebar_position: 4
---

# Expert Prompts

> Part of [Agent Prompting](./intro.md) documentation

## Introduction

Expert agents provide specialized knowledge in specific domains. Their prompts control identity, context detection (threading), and response generation. This document covers all three expert prompt types.

---

## Prompt Types Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    EXPERT PROMPTS                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────┐                                        │
│  │  SYSTEM PROMPT  │  Always active                         │
│  │                 │  Identity, specialization, rules       │
│  └────────┬────────┘                                        │
│           │                                                 │
│           ▼                                                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   THREADING                          │   │
│  │  ┌─────────────────────┐    ┌────────────────────┐  │   │
│  │  │ CONTEXT NOT FOUND   │───▶│   CONTEXT FOUND    │  │   │
│  │  │                     │    │                    │  │   │
│  │  │ Detect context primary key  │    │ Answer with context│  │   │
│  │  └─────────────────────┘    └────────────────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## System Prompt

The **System Prompt** defines the expert's identity and is always included.

### Purpose

- Establish agent identity and specialization
- Define domain expertise
- List what the agent can and cannot do
- Set collaboration rules

### Structure

```
1. Identity
   - Name and full ENS domain
   - Area of expertise

2. Context
   - Participation in agent meetings
   - Initial request reference
   - Other agents in the channel

3. Mission
   - What information to provide
   - How to handle requests

4. Constraints
   - What NOT to do
   - Behavioral limits
```

### Example

```
You are Claudine, the CRM expert for acme.
Full name: claudine.acme.137.protocol6022.eth
You work exclusively with customer email addresses.

Context:
You participate in AI agent meetings:
Initial request: {{ .InitialRequest.Content }}

Participating agents:
{{- range .ChannelMembers }}
{{- if .IsAgent }}
  Short name: {{ .Name }}
  Full name: {{ .EnsDomain }}
  Mission: {{ .Mission }}
  Role: {{ .Role }}
{{- end }}
{{- end }}

You have access to your own memory containing information related
to customer identifiers.

Your mission:
Provide verified, concrete information related to the customer
(account data, orders, payments, etc.).

What you must NEVER do:
• Give your own name or speak about yourself in third person
• Make assumptions or speak in future tense
• Give instructions to other agents
• Repeat a response already given in this thread
• Propose a response to the customer (that's the facilitator's job)
```

---

## Context Not Found Prompt

The **Context Not Found Prompt** is used during the threading phase when the agent needs to identify the context primary key.

### Purpose

- Analyze incoming message for key identification
- Match against known patterns or lists
- Request clarification if uncertain
- Establish the context for future responses

### Design Philosophy: Trust the Swarm

> **Important:** The Context Not Found prompt references `{{.LatestRequest.Content}}` (the **last message**), not the first message of the conversation.

This is a deliberate design choice. The swarm AI approach bets on collaboration rather than assumption:

```
┌─────────────────────────────────────────────────────────────┐
│              WHY USE LATEST MESSAGE?                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ❌ Traditional approach:                                   │
│     "The context primary key must be in the first message"  │
│     → Fragile, requires perfect initial request             │
│     → Fails if user doesn't provide key upfront             │
│                                                             │
│  ✅ Swarm AI approach:                                      │
│     "The facilitator and other agents know my mission.      │
│      They will discuss with me to find the right key."      │
│     → Resilient, leverages collective intelligence          │
│     → Works even with incomplete initial requests           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**How it works in practice:**

1. User sends initial request (key may or may not be present)
2. Facilitator calls expert with context
3. Expert analyzes **latest message** for key
4. If key not found → Expert asks for help (swarm collaboration)
5. Facilitator or other agents provide clarification
6. Expert receives new message with key → Locks context

This approach is more efficient because:

- Facilitators understand the expert's mission
- Other agents may have the missing information
- Conversation naturally converges toward the right key
- No need to force users to provide key in a specific format upfront

### Two Approaches

#### 1. Closed List (Few Keys)

When the agent has a limited number of known keys:

```
Identify the most relevant FAQ to handle this request:
{{.LatestRequest.Content}}
{{- if .LatestRequest.Author.IsAgent }}
sent by {{.LatestRequest.Author.EnsDomain }}
{{ end }}

Choose exclusively from this list:
{{range $i, $v := .ContextKeys}}{{if $i}}, {{end}}{{$v}}{{end}}

If you don't know, don't respond and ask other agents for help.

Start each response with 🔍

If you cannot find a match or are uncertain:
Ask the most qualified agent for help, using their full name.

Signature:
Start with 🔍. Maintain a strictly informative tone.
```

#### 2. Pattern Matching (Many Keys)

When keys follow a specific format (regex):

```
→ {{.LatestRequest.Content}}
{{- if .LatestRequest.Author.IsAgent }}
sent by {{.LatestRequest.Author.EnsDomain }}
{{ end }}

The UUID follows this REGEX format:
^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$

Verify its existence in your MCP.

If you cannot find a customer number in the message, or if you
are not absolutely certain of the identification:
• Immediately ask for help from the most qualified agent,
  calling them by their full name.

Signature:
Start your message with 🔍. Maintain a strictly informative tone.
```

### Key Variables

| Variable                                | Description                      |
| --------------------------------------- | -------------------------------- |
| `{{ .LatestRequest.Content }}`          | Message to analyze               |
| `{{ .ContextKeys }}`                    | Available keys (for closed list) |
| `{{ .LatestRequest.Author.EnsDomain }}` | Who sent the message             |

---

## Context Found Prompt

The **Context Found Prompt** is used after the context primary key is identified and locked.

### Purpose

- Respond using the established context
- Access memory for that key
- Collaborate with other agents when needed
- Provide useful, non-repetitive information

### Standard Template

```
Respond to this message "{{.LatestRequest.Content}}" sent by
{{- if .LatestRequest.Author.IsAgent }}
  {{ .LatestRequest.Author.EnsDomain }}
{{ else }}
  {{ .LatestRequest.Author.Name }}
{{ end }}
regarding {{ .ContextKey }} using only useful information.

If you don't know, you may ask another agent for help.
If you are not competent or lack sufficient information, state so.

Only if your response is useful, informative, and non-repetitive
for your interlocutor, mention their full name in the reply:
{{- if .LatestRequest.Author.IsAgent }}
  {{ .LatestRequest.Author.EnsDomain }}
{{ else }}
  {{ .LatestRequest.Author.Name }}
{{ end }}

Otherwise, indicate you are closing the conversation using the
agent's short name.

Start your message with 👍
```

### Key Variables

| Variable                        | Description                      |
| ------------------------------- | -------------------------------- |
| `{{ .ContextKey }}`             | The locked context primary key   |
| `{{ .LatestRequest.Content }}`  | Current message to respond to    |
| `{{ .LatestRequest.Author.* }}` | Information about message sender |

---

## Configuring Prompts

### Via Orchestrator UI

1. Navigate to your expert agent
2. Click **Actions** → **Manage prompts**
3. Select prompt type from dropdown:
   - System
   - Context Found
   - Context Not Found
4. Edit the prompt text
5. Click **Save**

### Prompt Type Selection

```
┌─────────────────────────────┐
│ Select a prompt type    ▼  │
├─────────────────────────────┤
│ System                      │
│ Context Found               │
│ Context Not Found           │
└─────────────────────────────┘
```

---

## Best Practices

### System Prompt

| Do                                   | Don't                         |
| ------------------------------------ | ----------------------------- |
| ✅ Clearly define expertise domain   | ❌ Include response templates |
| ✅ List explicit constraints         | ❌ Make it too long           |
| ✅ Use dynamic variables for context | ❌ Hardcode other agent names |
| ✅ Define collaboration rules        | ❌ Include threading logic    |

### Context Not Found Prompt

| Do                              | Don't                     |
| ------------------------------- | ------------------------- |
| ✅ Provide clear key format     | ❌ Allow guessing         |
| ✅ Define fallback behavior     | ❌ Ignore uncertainty     |
| ✅ Use regex for large key sets | ❌ List thousands of keys |
| ✅ Reference MCP if available   | ❌ Assume key exists      |

### Context Found Prompt

| Do                           | Don't                          |
| ---------------------------- | ------------------------------ |
| ✅ Reference the context key | ❌ Repeat previous responses   |
| ✅ Allow asking for help     | ❌ Give instructions to others |
| ✅ Define when to close      | ❌ Propose customer response   |
| ✅ Add signature/emoji       | ❌ Make it verbose             |

---

## Troubleshooting

### Threading Never Completes

**Symptoms:** Agent stays in THREADING state

**Causes:**

- Key pattern doesn't match message content
- Key not in closed list
- Regex too strict

**Solutions:**

1. Review key format in Context Not Found prompt
2. Add fallback to ask for clarification
3. Test regex with sample messages

### Wrong Key Selected

**Symptoms:** Agent uses incorrect context

**Causes:**

- Ambiguous patterns
- Multiple matches in message
- Similar key formats

**Solutions:**

1. Add disambiguation logic in prompt
2. Require explicit confirmation
3. Use more specific patterns

### Expert Doesn't Respond

**Symptoms:** Expert is mentioned but silent

**Causes:**

- ENS name not exact match
- Agent not in channel
- Orchestrator sync issue

**Solutions:**

1. Verify full ENS name spelling
2. Check agent is added to channel
3. Review orchestrator logs

### Repetitive Responses

**Symptoms:** Expert repeats same information

**Causes:**

- No instruction to avoid repetition
- Memory not being used
- Context not properly loaded

**Solutions:**

1. Add "don't repeat previous responses" instruction
2. Verify context key is locked
3. Check database for key content

---

## Complete Examples

### CRM Expert

**System:**

```
You are Claudine, the CRM expert.
Full name: claudine.company.137.protocol6022.eth
Specialization: Customer data (accounts, orders, history)

You access memory indexed by customer email or ID.

Never:
• Speak about yourself in third person
• Make assumptions about future
• Propose responses to customers
```

**Context Not Found:**

```
Analyze: {{.LatestRequest.Content}}

Look for customer identifier:
• Email format: *@*.*
• Customer ID: CUST-XXXXX

If found, use as context key.
If not found, ask @facilitator for clarification.

Start with 🔍
```

**Context Found:**

```
Context: {{ .ContextKey }}

Respond to: {{.LatestRequest.Content}}

Provide verified customer data only.
If data unavailable, state clearly.
Don't repeat previous responses.

Start with 👍
```

### GDPR Expert

**System:**

```
You are Victor, the GDPR compliance expert.
Full name: victor.company.137.protocol6022.eth
Specialization: Data protection regulations

Your memory contains article-specific guidance.

Never:
• Give legal advice
• Make definitive compliance statements
• Ignore jurisdiction differences
```

**Context Not Found:**

```
Analyze: {{.LatestRequest.Content}}

Look for GDPR reference:
• Article number: Article-XX
• Chapter: Chapter-X.X
• Topic keyword

Select from: {{ .ContextKeys }}

If uncertain, ask for clarification.

Start with 🔍
```

**Context Found:**

```
Regarding {{ .ContextKey }}:

Respond to: {{.LatestRequest.Content}}

Provide regulatory guidance.
Reference specific articles.
Note any limitations.

Start with 👍
```

---

## Integration with Memory

Expert prompts work with the memory system:

```
┌─────────────────────────────────────────────────────────────┐
│                    MEMORY INTEGRATION                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Context Not Found                                          │
│       │                                                     │
│       ▼                                                     │
│  Identify key: "C45653"                                    │
│       │                                                     │
│       ▼                                                     │
│  Lock context                                               │
│       │                                                     │
│       ▼                                                     │
│  Load memory for C45653  ◄── All previous conversations    │
│       │                       about this customer           │
│       ▼                                                     │
│  Context Found prompt                                       │
│       │                                                     │
│       ▼                                                     │
│  Response includes                                          │
│  historical context                                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

See [Agent Knowledge](../agent/agent-knowledge/intro.md) for memory management details.

---

## Related Documentation

- [Agent Types](./agent-types.md) — Understanding expert role
- [Conversation Flow](./conversation-flow.md) — Threading phase details
- [Dynamic Fields Reference](./prompt-dynamic-fields) — All template variables with code snippets
- [Context Primary Keys](../agent/agent-knowledge/context-primary-keys.md) — Key design
