# Prompt Writing Guide

> Part of [Agent Prompting](./intro.md) documentation

## Introduction

Writing effective prompts is both an art and a science. This guide provides proven patterns, common mistakes to avoid, and practical techniques for creating reliable agent behaviors.

---

## Core Principles

### 1. Be Specific

Vague prompts lead to unpredictable behavior.

```
❌ Bad:  "Help the customer"
✅ Good: "Provide the customer's order status and expected delivery date"
```

### 2. Be Structured

Use clear sections and steps.

```
❌ Bad:  Long paragraph with multiple instructions mixed together

✅ Good:
Step 1: Identify the customer
Step 2: Retrieve order information
Step 3: Formulate response
```

### 3. Be Explicit About Constraints

Define what the agent should NOT do.

```
What you must NEVER do:
• Make assumptions about data you don't have
• Propose solutions without verification
• Repeat information already shared
```

### 4. Use Dynamic Variables

Never hardcode information that might change.

```
❌ Bad:  "Contact John at john@company.com"
✅ Good: "Contact {{ .LatestRequest.Author.Name }}"
```

---

## Prompt Structure Template

```
┌─────────────────────────────────────────────────────────────┐
│                   RECOMMENDED STRUCTURE                     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. IDENTITY (Who am I?)                                   │
│     • Name and role                                        │
│     • Specialization                                       │
│                                                             │
│  2. CONTEXT (What's happening?)                            │
│     • Channel/meeting description                          │
│     • Initial request                                      │
│     • Available participants                               │
│                                                             │
│  3. MISSION (What should I do?)                            │
│     • Primary objective                                    │
│     • Secondary goals                                      │
│                                                             │
│  4. STEPS (How should I do it?)                            │
│     • Step-by-step process                                 │
│     • Decision points                                      │
│                                                             │
│  5. CONSTRAINTS (What must I avoid?)                       │
│     • Prohibited actions                                   │
│     • Limitations                                          │
│                                                             │
│  6. FORMAT (How should I respond?)                         │
│     • Output structure                                     │
│     • Signature/emoji                                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Common Mistakes

### Mistake 1: Infinite Loops

**Problem:** Agent calls itself or creates circular dependencies.

```
❌ Bad:
"If you're unsure, ask for help from any agent including yourself"

✅ Good:
"If you're unsure, ask for help from another expert agent.
 You cannot mention your own name."
```

**Prevention:**

- Explicitly forbid self-calls
- Limit the number of agent calls
- Add call count tracking

### Mistake 2: No Exit Condition

**Problem:** Conversation never concludes.

```
❌ Bad:
"Keep gathering information until you're satisfied"

✅ Good:
"Gather information until:
 • Confidence ≥ 8/10, OR
 • 3 expert calls made, OR
 • Timeout reached
 Then write :end:"
```

**Prevention:**

- Define explicit ending criteria
- Set maximum iterations
- Rely on timeout as fallback

### Mistake 3: Repetitive Responses

**Problem:** Agent repeats the same information.

```
❌ Bad:
"Provide all relevant information about the customer"

✅ Good:
"Provide relevant information about the customer.
 If you've already shared this information in this thread,
 state 'Already provided above' instead of repeating."
```

**Prevention:**

- Add "don't repeat" instruction
- Reference previous responses
- Track what's been shared

### Mistake 4: Overstepping Role

**Problem:** Expert proposes solutions meant for facilitator.

```
❌ Bad:
"Help the customer resolve their issue"

✅ Good:
"Provide factual information about the customer's account.
 Do NOT propose solutions or write customer responses.
 That is the facilitator's responsibility."
```

**Prevention:**

- Define role boundaries clearly
- Specify what NOT to do
- Separate data from decisions

### Mistake 5: Token Overflow

**Problem:** Prompt is too large for LLM context window.

```
❌ Bad:
System prompt with 50+ agents listed in detail

✅ Good:
System prompt with essential agents only, using variables
```

**Prevention:**

- Keep System prompts concise
- Use variables instead of static lists
- Monitor token usage

---

## Proven Patterns

### Pattern 1: Confidence Rating

```
On a scale of 1-10, rate your confidence in answering:
"{{ .InitialRequest.Content }}"

If confidence ≥ 8: Proceed with response
If confidence < 8: Call expert for missing information
```

### Pattern 2: Step-by-Step Synthesis

```
Step 1: 📋 SYNTHESIS
List key facts gathered so far:
• [Fact 1]
• [Fact 2]

Step 2: ❓ GAPS
Identify what's still missing:
• [Gap 1]

Step 3: 📞 ACTION
Call the appropriate expert OR conclude if complete.
```

### Pattern 3: Cross-Verification

```
Be skeptical of single-source information.
If one agent provides critical data, verify with a second agent.

"@expert2.eth, can you confirm what @expert1.eth said about X?"
```

### Pattern 4: Explicit Handoff

```
When mentioning another agent, provide context:

❌ Bad:  "@expert.eth, help please"
✅ Good: "@expert.eth, I need the order status for customer C45653
         who claims they haven't received their package"
```

### Pattern 5: Signature Consistency

```
Start each response with a consistent signature:

🔍 for Context Not Found (searching)
👍 for Context Found (answering)
⭐ for Synthesis
📫 for Final response
```

---

## Variable Usage

### Essential Variables

| Variable                        | Use In            | Purpose          |
| ------------------------------- | ----------------- | ---------------- |
| `{{ .InitialRequest.Content }}` | All prompts       | Original request |
| `{{ .LatestRequest.Content }}`  | Response prompts  | Current message  |
| `{{ .ContextKey }}`             | Context Found     | Locked key       |
| `{{ .ContextKeys }}`            | Context Not Found | Available keys   |
| `{{ .ChannelMembers }}`         | System            | Participants     |

### Conditional Logic

```
{{- if .LatestRequest.Author.IsAgent }}
  Responding to agent: {{ .LatestRequest.Author.EnsDomain }}
{{ else }}
  Responding to user: {{ .LatestRequest.Author.Name }}
{{ end }}
```

### Iteration

```
Available experts:
{{- range .ChannelMembers }}
{{- if .IsAgent }}
  • {{ .Name }}: {{ .Mission }}
{{- end }}
{{- end }}
```

---

## Length Guidelines

| Prompt Type           | Recommended Length | Maximum    |
| --------------------- | ------------------ | ---------- |
| **System**            | 200-500 tokens     | 800 tokens |
| **Default**           | 100-300 tokens     | 500 tokens |
| **Conclude**          | 150-400 tokens     | 600 tokens |
| **Context Not Found** | 100-250 tokens     | 400 tokens |
| **Context Found**     | 100-200 tokens     | 400 tokens |

### Why Length Matters

- **Token limits** — LLMs have context windows
- **Clarity** — Shorter prompts are clearer
- **Cost** — Longer prompts cost more
- **Speed** — More tokens = slower responses

---

## Testing Prompts

### Before Deployment

1. **Dry run** — Test with sample messages
2. **Edge cases** — Try unusual inputs
3. **Loop detection** — Verify no infinite loops
4. **Role boundaries** — Confirm agents stay in lane

### Test Scenarios

| Test              | What to Check            |
| ----------------- | ------------------------ |
| **Empty message** | Agent handles gracefully |
| **Multiple keys** | Correct disambiguation   |
| **Unknown key**   | Proper fallback behavior |
| **Agent offline** | Timeout handling         |
| **Long thread**   | Memory management        |

### Monitoring

After deployment:

- Review conversation logs
- Check for unexpected behaviors
- Monitor token usage
- Gather user feedback

---

## Checklist

### System Prompt

- [ ] Identity clearly defined
- [ ] Role and specialization stated
- [ ] Dynamic variables for participants
- [ ] Constraints listed
- [ ] Not too long (< 800 tokens)

### Default/Context Prompts

- [ ] Clear steps provided
- [ ] Exit condition defined
- [ ] No self-call possible
- [ ] Role boundaries respected
- [ ] Signature/emoji specified

### Conclude Prompt

- [ ] Synthesis step included
- [ ] Output format defined
- [ ] User-facing tone
- [ ] No reopening questions
- [ ] Quality guidelines included

---

## Quick Reference

### Do's

✅ Use dynamic variables
✅ Define explicit constraints
✅ Structure with steps
✅ Add exit conditions
✅ Include signatures
✅ Test edge cases
✅ Monitor after deployment

### Don'ts

❌ Hardcode names/values
❌ Allow self-calls
❌ Create circular logic
❌ Exceed token limits
❌ Mix role responsibilities
❌ Forget fallback behavior
❌ Skip testing

---

## Related Documentation

- [Facilitator Prompts](./facilitator-prompts.md) — Specific facilitator guidance
- [Expert Prompts](./expert-prompts.md) — Specific expert guidance
- [Prompt Variables](./prompt-dynamic-fields) — Complete variable reference
- [Prompt Templates](./prompt-templates.md) — Ready-to-use examples
