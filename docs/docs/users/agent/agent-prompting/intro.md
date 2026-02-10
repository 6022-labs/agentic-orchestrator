# Agent Prompting

> Documentation for configuring and writing effective agent prompts in the 6022 Agentic Teams framework.

## Overview

Agent prompts are the instructions that define how agents behave, communicate, and collaborate. Effective prompting is crucial for building reliable multi-agent systems.

This documentation covers:

- Understanding different agent types and their roles
- Configuring prompts for each conversation phase
- Writing effective prompts that avoid common pitfalls
- Ready-to-use templates for common use cases

---

## Documentation Structure

### Fundamentals

| #   | Topic                                       | Description                                                       |
| --- | ------------------------------------------- | ----------------------------------------------------------------- |
| 01  | [Agent Types](./agent-types.md)             | Facilitator, Expert, Human, Digital Twin — roles and capabilities |
| 02  | [Conversation Flow](./conversation-flow.md) | Lifecycle from start to conclusion                                |

### Prompt Configuration

| #   | Topic                                           | Description                                          |
| --- | ----------------------------------------------- | ---------------------------------------------------- |
| 03  | [Facilitator Prompts](./facilitator-prompts.md) | System, Default, and Conclude prompts                |
| 04  | [Expert Prompts](./expert-prompts.md)           | System, Context Found, and Context Not Found prompts |
| 05  | [Human & Digital Twin](./human-digital-twin.md) | Special agent types configuration                    |

### Best Practices

| #   | Topic                                               | Description                               |
| --- | --------------------------------------------------- | ----------------------------------------- |
| 06  | [Prompt Writing Guide](./prompt-writing-guide.md)   | Do's, don'ts, and proven patterns         |
| 07  | [Prompt Templates](./prompt-templates.md)           | Ready-to-use templates by domain          |
| 08  | [Dynamic Fields Reference](./prompt-dynamic-fields) | All template variables with code snippets |

---

## Quick Reference

### Agent Types at a Glance

| Type             | Role                                     | Prompts                                  | Context Primary Key          |
| ---------------- | ---------------------------------------- | ---------------------------------------- | ---------------------------- |
| **Facilitator**  | Orchestrates conversation, calls experts | System, Default, Conclude                | Channel-based                |
| **Expert**       | Specialized knowledge, answers questions | System, Context Found, Context Not Found | Domain-specific              |
| **Human**        | Real person in the loop                  | None (Slack/Teams)                       | Participant addresses        |
| **Digital Twin** | AI clone of a human                      | System + Default                         | Participant addresses (auto) |

### Prompt Types Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    FACILITATOR                              │
├─────────────────────────────────────────────────────────────┤
│  System      │ Identity, mission, participants              │
│  Default     │ Conversation management, calling experts     │
│  Conclude    │ Final response generation                    │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                      EXPERT                                 │
├─────────────────────────────────────────────────────────────┤
│  System            │ Identity, specialization, constraints  │
│  Context Not Found │ Key detection during threading         │
│  Context Found     │ Response generation with context       │
└─────────────────────────────────────────────────────────────┘
```

---

## Key Concepts

### Agent Naming

Agents are called by their **full ENS domain name**:

```
@claudine.company.137.protocol6022.eth
```

Using the full ENS name in a message triggers that agent.

### Conversation Ending

Facilitators end conversations via:

- **Timeout strategy** — Automatic after configured delay
- **Decide yourself strategy** — Agent writes `:end:` when ready

### Interaction Rules

| Rule                      | Description                                       |
| ------------------------- | ------------------------------------------------- |
| **No self-calls**         | Agents cannot mention themselves                  |
| **Explicit naming**       | Must use full ENS domain to call an agent         |
| **Facilitator oversight** | Only facilitator sees full conversation           |
| **Expert isolation**      | Experts only see messages where they're mentioned |

---

## Getting Started

### For New Users

1. Start with [Agent Types](./agent-types.md) to understand the roles
2. Read [Conversation Flow](./conversation-flow.md) to see how agents interact
3. Configure your first [Facilitator Prompts](./facilitator-prompts.md)

### For Prompt Engineers

1. Review [Prompt Writing Guide](./prompt-writing-guide.md) for best practices
2. Browse [Prompt Templates](./prompt-templates.md) for your domain
3. Check [Dynamic Fields Reference](./prompt-dynamic-fields) for copy-paste code
4. Customize templates for your specific use case

### For Troubleshooting

Common issues and solutions are covered in:

- [Facilitator Prompts](./facilitator-prompts.md#troubleshooting)
- [Expert Prompts](./expert-prompts.md#troubleshooting)
- [Prompt Writing Guide](./prompt-writing-guide.md#common-mistakes)

---

## Related Documentation

- [Agent Knowledge](../agent-knowledge/intro.md) — Memory and context management
- [Prompt Variables](./prompt-dynamic-fields) — Dynamic variables reference

---

## Glossary

| Term           | Definition                                                        |
| -------------- | ----------------------------------------------------------------- |
| **ENS Domain** | Ethereum Name Service identifier (e.g., `agent.company.6022.eth`) |
| **Threading**  | Process of detecting context primary key                          |
| **Conclude**   | Final phase where facilitator generates response                  |
| **Timeout**    | Automatic conversation ending after delay                         |
| **FORK**       | Operation to create a digital twin from a human                   |
| **Swarm**      | Group of collaborating agents                                     |
