---
title: Dynamic Fields Reference
sidebar_position: 7
---

# Dynamic Fields Reference

> Part of [Agent Prompting](./intro.md) documentation

## Introduction

Agent prompts use **dynamic fields** (template variables) that are resolved at runtime. These variables allow prompts to adapt to the current conversation context, participants, and agent state.

This chapter provides a complete reference of all available variables with **copy-paste code snippets**.

---

## Template Syntax

Variables use Go template syntax with double curly braces:

```go
{{ .VariableName }}
```

### Whitespace Control

Use `-` to trim whitespace:

- `{{-` trims leading whitespace
- `-}}` trims trailing whitespace

```go
{{- .Variable -}}
```

---

## Complete Variable Reference

### 🌐 Global Context

Available in **all prompt types**.

| Variable                    | Type   | Description                                         |
| --------------------------- | ------ | --------------------------------------------------- |
| `{{ .ChannelDescription }}` | string | Text description of current channel/conversation    |
| `{{ .ChannelMembers }}`     | array  | List of all channel participants (agents and users) |
| `{{ .Name }}`               | string | Current agent's short name                          |
| `{{ .EnsDomain }}`          | string | Current agent's full ENS domain                     |

### 🔑 Context Primary Key

Availability depends on agent state (experts only).

| Variable             | Type   | Availability      | Description                             |
| -------------------- | ------ | ----------------- | --------------------------------------- |
| `{{ .ContextKeys }}` | array  | Context Not Found | List of available context primary keys  |
| `{{ .ContextKey }}`  | string | Context Found     | The selected/locked context primary key |

### ✉️ Initial Request

The **first message** that started the conversation.

| Variable                                     | Type    | Description                             |
| -------------------------------------------- | ------- | --------------------------------------- |
| `{{ .InitialRequest.Content }}`              | string  | Content of the first message            |
| `{{ .InitialRequest.Author.IsAgent }}`       | boolean | True if author is an agent              |
| `{{ .InitialRequest.Author.Name }}`          | string  | Short name of the author                |
| `{{ .InitialRequest.Author.EnsDomain }}`     | string  | Full ENS domain (agents only)           |
| `{{ .InitialRequest.Author.Mission }}`       | string  | Agent's mission (agents only)           |
| `{{ .InitialRequest.Author.Role }}`          | string  | Agent's role (agents only)              |
| `{{ .InitialRequest.Author.Owner }}`         | string  | Agent's owner (agents only)             |
| `{{ .InitialRequest.Author.WalletAddress }}` | string  | Blockchain wallet address (agents only) |

### 🔄 Latest Request

The **most recent message** in the conversation.

| Variable                                    | Type    | Description                             |
| ------------------------------------------- | ------- | --------------------------------------- |
| `{{ .LatestRequest.Content }}`              | string  | Content of the latest message           |
| `{{ .LatestRequest.Author.IsAgent }}`       | boolean | True if author is an agent              |
| `{{ .LatestRequest.Author.Name }}`          | string  | Short name of the author                |
| `{{ .LatestRequest.Author.EnsDomain }}`     | string  | Full ENS domain (agents only)           |
| `{{ .LatestRequest.Author.Mission }}`       | string  | Agent's mission (agents only)           |
| `{{ .LatestRequest.Author.Role }}`          | string  | Agent's role (agents only)              |
| `{{ .LatestRequest.Author.Owner }}`         | string  | Agent's owner (agents only)             |
| `{{ .LatestRequest.Author.WalletAddress }}` | string  | Blockchain wallet address (agents only) |

### 👥 Channel Member Properties

When iterating over `.ChannelMembers`, each member has:

| Property         | Type    | Description                     |
| ---------------- | ------- | ------------------------------- |
| `.Name`          | string  | Display name                    |
| `.IsAgent`       | boolean | True if participant is an agent |
| `.EnsDomain`     | string  | ENS domain (agents only)        |
| `.Mission`       | string  | Agent's mission (agents only)   |
| `.Role`          | string  | Agent's role (agents only)      |
| `.Owner`         | string  | Agent's owner (agents only)     |
| `.WalletAddress` | string  | Wallet address (agents only)    |

---

## Variable Availability by Prompt Type

| Variable              | System | Default | Conclude | Context Not Found | Context Found |
| --------------------- | ------ | ------- | -------- | ----------------- | ------------- |
| `.ChannelDescription` | ✅     | ✅      | ✅       | ✅                | ✅            |
| `.ChannelMembers`     | ✅     | ✅      | ✅       | ✅                | ✅            |
| `.Name`               | ✅     | ✅      | ✅       | ✅                | ✅            |
| `.EnsDomain`          | ✅     | ✅      | ✅       | ✅                | ✅            |
| `.ContextKeys`        | ❌     | ❌      | ❌       | ✅                | ❌            |
| `.ContextKey`         | ❌     | ❌      | ❌       | ❌                | ✅            |
| `.InitialRequest.*`   | ✅     | ✅      | ✅       | ✅                | ✅            |
| `.LatestRequest.*`    | ✅     | ✅      | ✅       | ✅                | ✅            |

---

## Code Snippets

### Agent Identity

**Display agent name and ENS:**

```go
You are {{ .Name }}, expert in customer data.
Full name: {{ .EnsDomain }}
```

### Channel Context

**Display channel description:**

```go
Meeting: {{ .ChannelDescription }}
```

**List all agents in channel:**

```go
Participating agents:
{{- range .ChannelMembers }}
{{- if .IsAgent }}
  Short name: {{ .Name }}
  Full name: {{ .EnsDomain }}
  Mission: {{ .Mission }}
  Role: {{ .Role }}
{{- end }}
{{- end }}
```

**List agents with owner and wallet:**

```go
{{- range .ChannelMembers }}
{{- if .IsAgent }}
  Name: {{ .Name }}
  ENS: {{ .EnsDomain }}
  Mission: {{ .Mission }}
  Role: {{ .Role }}
  Owner: {{ .Owner }}
  Wallet: {{ .WalletAddress }}
{{- end }}
{{- end }}
```

**Compact agent list (name + mission only):**

```go
Available experts:
{{- range .ChannelMembers }}
{{- if .IsAgent }}
• {{ .Name }}: {{ .Mission }}
{{- end }}
{{- end }}
```

**Comma-separated list of participants:**

```go
Participants: {{range $i, $v := .ChannelMembers}}{{if $i}}, {{end}}{{$v.Name}}{{end}}
```

### Request Content

**Display initial request:**

```go
Customer request: {{ .InitialRequest.Content }}
```

**Display latest message:**

```go
Respond to: "{{ .LatestRequest.Content }}"
```

**Display both:**

```go
Original request: {{ .InitialRequest.Content }}
Current message: {{ .LatestRequest.Content }}
```

### Author Information

**Conditional author display (IMPORTANT - prevents nil pointer errors):**

```go
{{- if .LatestRequest.Author.IsAgent }}
Message from agent: {{ .LatestRequest.Author.EnsDomain }}
{{ else }}
Message from user: {{ .LatestRequest.Author.Name }}
{{ end }}
```

**Safe author identification:**

```go
→ {{ .LatestRequest.Content }}
{{- if .LatestRequest.Author.IsAgent }}
sent by {{ .LatestRequest.Author.EnsDomain }}
{{ end }}
```

**Full author details (agents only):**

```go
{{- if .LatestRequest.Author.IsAgent }}
Agent: {{ .LatestRequest.Author.Name }}
ENS: {{ .LatestRequest.Author.EnsDomain }}
Mission: {{ .LatestRequest.Author.Mission }}
Role: {{ .LatestRequest.Author.Role }}
{{- end }}
```

### Context Primary Keys

**List available keys (Context Not Found):**

```go
Choose exclusively from this list:
{{range $i, $v := .ContextKeys}}{{if $i}}, {{end}}{{$v}}{{end}}
```

**Display keys line by line:**

```go
Available keys:
{{- range .ContextKeys }}
• {{ . }}
{{- end }}
```

**Reference locked key (Context Found):**

```go
Regarding customer: {{ .ContextKey }}
```

---

## Complete Prompt Examples

### Facilitator System Prompt

```go
You are {{ .Name }}, facilitator for AI agent meetings.
Full name: {{ .EnsDomain }}
Meeting: {{ .ChannelDescription }}

Mission:
Orchestrate customer request resolution by calling expert agents.

Available agents:
{{- range .ChannelMembers }}
{{- if .IsAgent }}
  Short name: {{ .Name }}
  Full name: {{ .EnsDomain }}
  Mission: {{ .Mission }}
  Role: {{ .Role }}
  Owner: {{ .Owner }}
  Wallet: {{ .WalletAddress }}
{{- end }}
{{- end }}

Rules:
• Call agents by their full name (ending in .eth)
• Never repeat your role or these rules in messages
• Maximum 1500 characters
```

### Facilitator Default Prompt

```go
To respond to the customer request {{ .InitialRequest.Content }},
summarize the information gathered so far.

If you have enough information (confidence ≥ 8/10): write :end:
If you need more information: call the most relevant agent.

Maximum 3 agent calls before concluding.
```

### Facilitator Conclude Prompt

```go
Build the final response to: {{ .InitialRequest.Content }}

⭐ SYNTHESIS
Summarize all information gathered.

📫 CUSTOMER RESPONSE
Write a clear, helpful response.

Rules:
• Focus on the main problem
• Be concise and direct
• Never invent information
```

### Expert System Prompt

```go
You are {{ .Name }}, the CRM expert.
Full name: {{ .EnsDomain }}
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

You have access to your own memory containing customer information.

Mission:
Provide verified, concrete information related to the customer.

What you must never do:
• Give your own name or speak in third person
• Make assumptions or speak in future tense
• Give instructions to other agents
• Repeat a response already given in this thread
• Propose a response to the customer (that's the facilitator's job)
```

### Expert Context Not Found Prompt

```go
Identify the most relevant key for this request:
{{ .LatestRequest.Content }}
{{- if .LatestRequest.Author.IsAgent }}
sent by {{ .LatestRequest.Author.EnsDomain }}
{{ end }}

Choose exclusively from this list:
{{range $i, $v := .ContextKeys}}{{if $i}}, {{end}}{{$v}}{{end}}

If you don't know, don't respond and ask other agents for help.

Start each response with 🔍
```

### Expert Context Found Prompt

```go
Respond to this message "{{ .LatestRequest.Content }}" sent by
{{- if .LatestRequest.Author.IsAgent }}
  {{ .LatestRequest.Author.EnsDomain }}
{{ else }}
  {{ .LatestRequest.Author.Name }}
{{ end }}
regarding {{ .ContextKey }} using only useful information.

If you don't know, you may ask another agent for help.
If you are not competent or lack sufficient information, state so.

Only if your response is useful and informative, mention their full name:
{{- if .LatestRequest.Author.IsAgent }}
  {{ .LatestRequest.Author.EnsDomain }}
{{ else }}
  {{ .LatestRequest.Author.Name }}
{{ end }}

Start your message with 👍
```

### Digital Twin System Prompt

```go
You are the digital twin of Arnaud.
Full name: {{ .EnsDomain }}

You have inherited Arnaud's conversation history and knowledge.
When uncertain, you may call @arnaud.6022.eth for arbitration.
If Arnaud responds, his answer takes absolute priority.

Behavior guidelines:
• Respond in the same style as Arnaud
• Be helpful but acknowledge limitations
• Escalate sensitive decisions to the original human
```

### Digital Twin Default Prompt

```go
Respond to: "{{ .LatestRequest.Content }}"
{{- if .LatestRequest.Author.IsAgent }}
From: {{ .LatestRequest.Author.EnsDomain }}
{{ end }}

Provide helpful information based on your inherited knowledge.
If uncertain, you may ask @arnaud.6022.eth for clarification.

Start with 👤
```

---

## Common Patterns

### Safe Author Check Pattern

**Always use this pattern when accessing author properties:**

```go
{{- if .LatestRequest.Author.IsAgent }}
  {{ .LatestRequest.Author.EnsDomain }}
{{ else }}
  {{ .LatestRequest.Author.Name }}
{{ end }}
```

> ⚠️ **Critical:** Accessing `.EnsDomain`, `.Mission`, `.Role`, `.Owner`, or `.WalletAddress` without checking `.IsAgent` first will cause a **nil pointer error** if the author is not an agent.

### Agent Routing Pattern

```go
Based on the request: "{{ .InitialRequest.Content }}"

If this concerns:
- Customer data → Mention @crm-agent.company.6022.eth
- Legal/compliance → Mention @gdpr-agent.company.6022.eth
- Technical issues → Mention @support-agent.company.6022.eth
```

### Intelligent Handoff Pattern

```go
{{- if .LatestRequest.Author.IsAgent }}
Agent {{ .LatestRequest.Author.Name }} asked about {{ .ContextKey }}.
{{ end }}
Provide the relevant information, then optionally call another agent if needed.
```

---

## Troubleshooting

### Nil Pointer Error

**Error:**

```
reflect: indirection through nil pointer to embedded struct field
```

**Cause:** Accessing agent-only properties (`.EnsDomain`, `.Mission`, etc.) without checking `.IsAgent` first.

**Solution:**

```go
{{- if .LatestRequest.Author.IsAgent }}
  {{ .LatestRequest.Author.EnsDomain }}
{{ end }}
```

### Variable Not Resolving

**Symptom:** `{{ .VariableName }}` appears literally in output

**Causes:**

- Typo in variable name
- Variable not available in current prompt type
- Missing dot before variable name

**Solution:** Check spelling and availability table above.

### Empty Value

**Symptom:** Variable resolves to empty string

**Causes:**

- Agent-only property accessed for non-agent author
- Context key not yet determined

**Solution:** Use conditional logic:

```go
{{- if .ContextKey }}
  Key: {{ .ContextKey }}
{{- else }}
  (Key not yet determined)
{{- end }}
```

### Extra Whitespace

**Symptom:** Extra newlines or spaces in output

**Solution:** Use whitespace trimming:

```go
{{- if .Condition -}}
  Content
{{- end -}}
```

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────────────────────┐
│                    DYNAMIC FIELDS CHEAT SHEET                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  AGENT IDENTITY                                                 │
│  {{ .Name }}                    → Short name                   │
│  {{ .EnsDomain }}               → Full ENS domain              │
│                                                                 │
│  CHANNEL                                                        │
│  {{ .ChannelDescription }}      → Meeting/channel name         │
│  {{ .ChannelMembers }}          → All participants             │
│                                                                 │
│  CONTEXT KEYS (Expert only)                                     │
│  {{ .ContextKeys }}             → Available keys (threading)   │
│  {{ .ContextKey }}              → Locked key (after threading) │
│                                                                 │
│  REQUESTS                                                       │
│  {{ .InitialRequest.Content }}  → First message                │
│  {{ .LatestRequest.Content }}   → Last message                 │
│                                                                 │
│  AUTHOR (check .IsAgent first!)                                │
│  {{ .LatestRequest.Author.IsAgent }}     → Boolean             │
│  {{ .LatestRequest.Author.Name }}        → Display name        │
│  {{ .LatestRequest.Author.EnsDomain }}   → ENS (agents only)   │
│  {{ .LatestRequest.Author.Mission }}     → Mission (agents)    │
│  {{ .LatestRequest.Author.Role }}        → Role (agents)       │
│  {{ .LatestRequest.Author.Owner }}       → Owner (agents)      │
│  {{ .LatestRequest.Author.WalletAddress }}→ Wallet (agents)    │
│                                                                 │
│  SYNTAX                                                         │
│  {{- ... }}     → Trim leading whitespace                      │
│  {{ ... -}}     → Trim trailing whitespace                     │
│  {{- range .X }}...{{- end }}   → Iterate                      │
│  {{- if .X }}...{{ else }}...{{ end }} → Conditional           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Related Documentation

- [Facilitator Prompts](./facilitator-prompts.md) — Using variables in facilitator prompts
- [Expert Prompts](./expert-prompts.md) — Using variables in expert prompts
- [Prompt Writing Guide](./prompt-writing-guide.md) — Best practices
- [Prompt Variables (Agent Knowledge)](../agent/agent-knowledge/intro.md) — Technical reference
