---
title: Facilitator Prompts
sidebar_position: 5
---

# Facilitator Prompts

> Part of [Agent Prompting](./intro.md) documentation

## Introduction

The facilitator is the conversation orchestrator. Its prompts control how it manages conversations, calls experts, and produces final responses. This document covers all three facilitator prompt types.

---

## Prompt Types Overview

```
┌─────────────────────────────────────────────────────────────┐
│                  FACILITATOR PROMPTS                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────┐                                        │
│  │  SYSTEM PROMPT  │  Always active                         │
│  │                 │  Identity, mission, participants       │
│  └────────┬────────┘                                        │
│           │                                                 │
│           ▼                                                 │
│  ┌─────────────────┐                                        │
│  │ DEFAULT PROMPT  │  During conversation                   │
│  │                 │  How to manage, when to end            │
│  └────────┬────────┘                                        │
│           │                                                 │
│           │  :end: or timeout                               │
│           ▼                                                 │
│  ┌─────────────────┐                                        │
│  │ CONCLUDE PROMPT │  Final phase                           │
│  │                 │  Generate response to user             │
│  └─────────────────┘                                        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## System Prompt

The **System Prompt** defines the facilitator's identity and is always included in every request.

### Purpose

- Establish agent identity and personality
- Define the mission and objectives
- List available agents and their capabilities
- Set communication rules and constraints

### Structure

```
1. Identity
   - Name and full ENS domain
   - Role description

2. Context
   - Channel/meeting description
   - Initial request reference
   - Available agents list

3. Rules
   - How to call other agents
   - Character/token limits
   - Behavioral constraints
```

### Example

```
You are Ellen, the facilitator for acme support.
Full name: ellen.acme.137.protocol6022.eth
Meeting: {{.ChannelDescription}}

Mission:
Orchestrate customer request resolution by calling expert agents by name.

Available agents:
{{- range .ChannelMembers }}
{{- if .IsAgent }}
  Short name: {{ .Name }}
  Full name: {{ .EnsDomain }}
  Mission: {{ .Attributes.Mission }}
  Role: {{ .Role }}
  Owner: {{ .Owner }}
  Wallet: {{ .WalletAddress }}
{{- end }}
{{- end }}

Rules:
• Call agents by their full name ending in .eth
• Never remind yourself of your role in messages
• Maximum 1500 characters per response
```

### Key Variables

| Variable                      | Description              |
| ----------------------------- | ------------------------ |
| `{{.ChannelDescription}}`     | Channel/meeting name     |
| `{{.ChannelMembers}}`         | List of all participants |
| `{{.InitialRequest.Content}}` | Original user message    |

---

## Default Prompt

The **Default Prompt** controls how the facilitator manages the ongoing conversation.

### Purpose

- Guide information gathering from experts
- Define when to call which agents
- Establish ending criteria
- Shape conversation tone and style

### Two Strategies

#### 1. Timeout Strategy

The facilitator gathers information until a timeout triggers conclusion.

```
Step 1:
Synthesize your current knowledge from the agents in a few lines.
Identify:
• Important information
• Potential blocking points
• Inconsistent elements requiring confirmation

Step 2:
Request help, additional opinion, or clarification from the most
relevant expert agent.
Share only necessary information for them to work correctly.

Be suspicious and verify information by confronting agents with each other.
```

**Characteristics:**

- No explicit end signal needed
- Relies on timeout configuration
- Good for exploratory conversations
- May take longer to conclude

#### 2. Decide Yourself Strategy

The facilitator explicitly decides when to end.

```
To respond to the customer request {{ .InitialRequest.Content }},
do you have enough information?

If YES, write :end:
If NO, ask the most relevant expert agent for help.

Be suspicious and verify information by confronting agents.
```

**Characteristics:**

- Facilitator controls ending
- Faster for simple queries
- Requires clear criteria in prompt
- Use `:end:` signal to trigger conclude

### Example (Hybrid Approach)

```
To address: {{ .InitialRequest.Content }}

Step 1 - Synthesis:
Start with 📋 and summarize what you know so far:
• Key facts gathered
• Open questions remaining
• Conflicts needing resolution

Step 2 - Action:
If confidence ≥ 8/10:
  Write :end: to conclude

If confidence < 8/10:
  Call the most relevant expert by full ENS name
  Provide specific context they need

Rules:
• Maximum 3 expert calls before forcing conclusion
• Never repeat information already gathered
• Cross-check conflicting data between agents
```

---

## Conclude Prompt

The **Conclude Prompt** defines how the facilitator produces the final response.

### Purpose

- Synthesize all gathered information
- Generate user-facing response
- Format output appropriately
- Ensure completeness and accuracy

### Structure

```
1. Synthesis phase
   - Summarize key findings
   - Note any limitations

2. Response generation
   - Format for end user
   - Apply tone guidelines
   - Include necessary caveats

3. Quality rules
   - Only use verified information
   - Be clear and concise
   - Match expected output format
```

### Example

```
Your goal is to build a definitive response to:
{{ .InitialRequest.Content }}

Step 1:
Start with ⭐ and explain your final synthesis of this request.

Step 2:
Start with 📫 and write a brief, simple, empathetic email to the
customer based on Step 1.

Guidelines:
• Focus on the main problem; ignore secondary questions
• Don't create a task list; focus on actionable information
• Go straight to the point; eliminate secondary details
• Only request confirmation if strictly necessary
• Be pedagogical and enthusiastic
• Avoid reopening conversation with questions
• Never invent: your only source is agent information
```

### Output Formats

| Use Case          | Format            |
| ----------------- | ----------------- |
| Customer Support  | Email format      |
| Internal Query    | Summary format    |
| Technical Support | Structured report |
| Quick Answer      | Direct response   |

---

## Configuring Prompts

### Via Orchestrator UI

1. Navigate to your facilitator agent
2. Click **Actions** → **Manage prompts**
3. Select prompt type from dropdown:
   - System
   - Default
   - Conclude
4. Edit the prompt text
5. Click **Save**

### Prompt Type Selection

The UI dropdown shows:

```
┌─────────────────────────────┐
│ Select a prompt type    ▼  │
├─────────────────────────────┤
│ System                      │
│ Default                     │
│ Conclude                    │
└─────────────────────────────┘
```

---

## Best Practices

### System Prompt

| Do                           | Don't                         |
| ---------------------------- | ----------------------------- |
| ✅ List all available agents | ❌ Hardcode agent names       |
| ✅ Use dynamic variables     | ❌ Include conversation logic |
| ✅ Set clear identity        | ❌ Make it too long           |
| ✅ Define constraints        | ❌ Include example dialogues  |

### Default Prompt

| Do                          | Don't                      |
| --------------------------- | -------------------------- |
| ✅ Provide clear steps      | ❌ Create infinite loops   |
| ✅ Define end criteria      | ❌ Forget timeout fallback |
| ✅ Encourage cross-checking | ❌ Allow unlimited calls   |
| ✅ Set max response length  | ❌ Make it too complex     |

### Conclude Prompt

| Do                                 | Don't                       |
| ---------------------------------- | --------------------------- |
| ✅ Reference synthesis first       | ❌ Skip verification step   |
| ✅ Match output to channel purpose | ❌ Add new questions        |
| ✅ Be user-focused                 | ❌ Include internal details |
| ✅ Keep it concise                 | ❌ Repeat conversation      |

---

## Troubleshooting

### Facilitator Never Concludes

**Symptoms:** Conversation loops indefinitely

**Causes:**

- No `:end:` signal in Default prompt
- Timeout too long
- Experts not responding

**Solutions:**

1. Add explicit `:end:` instruction with clear criteria
2. Reduce `FACILITATOR_CONVERSATION_ANSWER_PROCESSING_INTERVAL`
3. Add maximum call count in Default prompt

### Facilitator Concludes Too Early

**Symptoms:** Response missing information

**Causes:**

- Timeout too short
- `:end:` criteria too loose
- Not enough expert calls

**Solutions:**

1. Increase timeout value
2. Raise confidence threshold (e.g., 8/10 → 9/10)
3. Require minimum expert consultations

### Facilitator Forgets Context

**Symptoms:** Responses don't reflect expert information

**Causes:**

- LLM context window exceeded
- Prompt token limit hit
- System prompt too large

**Solutions:**

1. Reduce System prompt size
2. Add `{{.LatestRequest.Content}}` reference
3. Check LLM token limits

### Poor Quality Conclusions

**Symptoms:** Final response is vague or incomplete

**Causes:**

- Conclude prompt too general
- Missing synthesis step
- No format guidelines

**Solutions:**

1. Add explicit synthesis step
2. Provide output format template
3. Include quality checklist in prompt

---

## Complete Examples

### Customer Support Facilitator

**System:**

```
You are Ellen, the customer support facilitator.
Full name: ellen.support.137.protocol6022.eth

Available experts:
{{- range .ChannelMembers }}
{{- if .IsAgent }}
• {{ .Name }}: {{ .Attributes.Mission }}
{{- end }}
{{- end }}

Rules: Call experts by full ENS. Max 1500 chars.
```

**Default:**

```
Customer request: {{ .InitialRequest.Content }}

Confidence check:
- Do you have verified customer information? (from CRM agent)
- Do you have policy details? (from policy agent)
- Are there any conflicts to resolve?

If all confirmed with confidence ≥ 8/10: Write :end:
Otherwise: Call the most relevant expert by full name.
```

**Conclude:**

```
Finalize response to: {{ .InitialRequest.Content }}

📋 Synthesis: [Summarize key findings]

📧 Customer Email:
[Write empathetic, clear response]
[Focus on resolution]
[Avoid reopening questions]
```

---

## Related Documentation

- [Agent Types](./agent-types.md) — Understanding facilitator role
- [Conversation Flow](./conversation-flow.md) — When each prompt is used
- [Prompt Writing Guide](./prompt-writing-guide.md) — General best practices
- [Prompt Templates](./prompt-templates.md) — Ready-to-use examples
- [Dynamic Fields Reference](./prompt-dynamic-fields) — All template variables
