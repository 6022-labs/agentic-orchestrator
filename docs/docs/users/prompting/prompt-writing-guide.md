---
title: Prompt Writing Guide
sidebar_position: 9
---

# Prompt Writing Guide

Practical patterns that work across agent types. Read this once, then use [Prompt Templates](./prompt-templates.md) as your starting point.

## Four principles

1. **Be specific.** `"Help the customer"` is useless. `"Provide the customer's order status and expected delivery date"` is a job.
2. **Be structured.** Use numbered steps, not dense paragraphs. The model follows lists more reliably than prose.
3. **Say what not to do.** Constraint sections (`Never ...`, `Do NOT ...`) reliably prevent the most common failures — overstepping the role, repeating information, inventing data.
4. **Use dynamic variables.** Never hardcode a name, email, or agent identifier. Every variable you reference resolves at runtime — see [Dynamic Fields](./prompt-dynamic-fields.md).

## Structure of a good System prompt

Six sections, each one line to a few lines:

1. **Identity** — who am I, what's my specialisation
2. **Context** — channel, initial request, other agents
3. **Mission** — one sentence, concrete
4. **Rules** — explicit constraints, "never do X"
5. **Signature** — an emoji or prefix so responses are visually distinguishable
6. **Length cap** — `Maximum 1500 characters` tends to actually work

The other prompt types (Default, Conclude, Context Not Found, Context Found) follow the same shape but compress harder — most are under 20 lines.

## Five common mistakes

### 1. No exit condition

```
❌ "Keep gathering information until you're satisfied."
✅ "Gather information until confidence ≥ 8/10 OR 3 expert calls made.
    Then write :end:."
```

Without an explicit stop, the facilitator loops forever or hits the timeout instead of ending cleanly.

### 2. Overstepping the role

Experts writing customer-facing responses, facilitators inventing customer data. Add a hard rule to the System prompt:

```
Do NOT propose responses to the customer. That is the facilitator's job.
```

### 3. Repetitive responses

The model will happily repeat what it already said. Counter it explicitly:

```
If you've already shared this information in this thread,
write "Already provided above" instead of repeating.
```

### 4. Self-calls and infinite loops

Agents calling themselves, or two experts bouncing the same question. Forbid self-calls in the System prompt and cap the number of agent calls per turn in the Default prompt.

### 5. Oversize System prompt

Hardcoding a list of 50 agents blows the context window. Use `{{ range .ChannelMembers }}` to render the list dynamically.

## Patterns that work

### Confidence rating

```
On a scale of 1-10, rate your confidence in answering
"{{ .InitialRequest.Content }}".

If ≥ 8: respond. If < 8: call the most relevant expert for missing info.
```

### Step-by-step synthesis

```
📋 Synthesis: list the facts you've gathered.
❓ Gaps: list what's still missing.
📞 Action: call the right expert, or write :end: if complete.
```

### Explicit handoff

When calling another agent, always include the context they need — don't make them guess.

```
❌ "@expert.eth, help please"
✅ "@expert.eth, I need the order status for customer C45653
     who says their package hasn't arrived."
```

### Signature consistency

Give each prompt type a distinctive emoji so responses are scannable:

- 🔍 Context Not Found (searching)
- 👍 Context Found (answering)
- ⭐ Synthesis
- 📫 Final response to the user

## Length limits

System prompts are part of every inference, so a 1000-token System prompt costs tokens and latency on every turn. Rough budgets:

- **System:** 200-500 tokens (hard cap ~800)
- **Default / Context Found / Context Not Found:** 100-300 tokens
- **Conclude:** 150-400 tokens

If you're hitting the cap, the fix is almost always to move content from the System prompt into the turn-specific prompts.
