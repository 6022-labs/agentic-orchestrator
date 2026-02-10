# Human & Digital Twin

> Part of [Agent Prompting](./intro.md) documentation

## Introduction

Human agents and Digital Twins are special agent types that bridge the gap between AI and human participants. This document covers their configuration and unique behaviors.

---

## Human Agent

### Overview

A **Human Agent** represents a real person participating in conversations via Slack or Teams. Unlike AI agents, humans don't have prompts — they interact directly through the messaging platform.

```
┌─────────────────────────────────────────────────────────────┐
│                     HUMAN AGENT                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│    ┌──────────────┐         ┌──────────────┐               │
│    │    SLACK     │         │  HUMAN AGENT │               │
│    │    /TEAMS    │◄───────▶│              │               │
│    │              │         │  arnaud.6022 │               │
│    │   UI Input   │         │     .eth     │               │
│    └──────────────┘         └──────────────┘               │
│                                    │                        │
│                                    ▼                        │
│                             ┌──────────────┐               │
│                             │  ORCHESTRATOR│               │
│                             │              │               │
│                             │  Attributes  │               │
│                             │  messages to │               │
│                             │  ENS domain  │               │
│                             └──────────────┘               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Key Characteristics

| Property | Description |
|----------|-------------|
| **On-Chain Identity** | NFT + wallet + unique ENS domain |
| **Nature** | "Biological LLM" — human providing responses |
| **No Prompts** | Uses Slack/Teams interface directly |
| **Context Primary Key** | Unordered list of participant ETH addresses |
| **Memory** | Full exchange history stored in database |
| **Superpower** | Can restart concluded conversations |

### Configuration

#### Step 1: Create Human Agent NFT

Create the agent in the orchestrator UI with role = "human".

> **Important — Wallet Setup:**
> 
> Creating a human agent requires **two distinct wallets**:
> - **Owner wallet** — Controls the agent NFT, can modify configuration
> - **Agent wallet** — Represents the human agent's on-chain identity
> 
> **Special case:** If you use the same wallet for both, the owner is also the human represented by the agent. This is valid but means the same person owns and embodies the agent.

#### Step 2: Pair with Slack User

The pairing process requires the human user (or the agent owner) to authenticate their wallet:

**2.1 — Verify Wallet**

Go to any Slack channel where the agent operates and execute:

```
/verify-wallet
```

**2.2 — Complete Authentication Wizard**

Follow the authentication wizard that appears. This verifies your wallet ownership.

**2.3 — Copy Pairing Command from Orchestrator**

In the orchestrator UI, find the pairing command for your agent. It will look like:

```
/pair-slack-user-to-human-agent [USER ID] <COLLECTION_ADDRESS> <TOKEN_ID>
```

**2.4 — Execute in Slack**

Replace `[USER ID]` with your Slack User ID and execute the command in Slack:

```
/pair-slack-user-to-human-agent U09BDFS3VGA 0x7a3B1f8c2D4e9F0A1b2C3d4E5f6A7B8c9D0e1F2a 12
```

> **Finding your Slack User ID:** Click on your profile in Slack → "More" → "Copy member ID"

#### Step 3: Verify Pairing

Test by mentioning the human's ENS in a conversation:
```
@arnaud.6022.eth, can you confirm this policy?
```

### Human Behaviors

#### Identified in Conversations

When a human posts a message:
1. Orchestrator detects the Slack/Teams user
2. Matches to paired human agent
3. Attributes message to ENS domain
4. Stores in human's database

#### Can Restart Conversations

Only humans can restart concluded conversations:

```
Conversation concluded
       │
       ▼
Human sends new message
       │
       ▼
Orchestrator detects human status
       │
       ▼
Conversation reopens
       │
       ▼
Facilitator re-enters with Default prompt
```

#### Database Building

Every human interaction is stored:
- Creates history for potential digital twin
- Enables cloning via FORK operation
- Preserves context across conversations

---

## Digital Twin

### Overview

A **Digital Twin** is an AI clone of a human agent, created by duplicating the human's entire conversation history.

```
┌─────────────────────────────────────────────────────────────┐
│                    DIGITAL TWIN                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌───────────────┐      FORK       ┌───────────────┐       │
│  │  HUMAN AGENT  │ ──────────────▶ │ DIGITAL TWIN  │       │
│  │               │                 │               │       │
│  │ arnaud.6022   │                 │ arnaud.clone  │       │
│  │    .eth       │                 │   .6022.eth   │       │
│  │               │                 │               │       │
│  │  ┌─────────┐  │   Copy entire   │  ┌─────────┐  │       │
│  │  │ Memory  │──┼────────────────▶│  │ Memory  │  │       │
│  │  │ History │  │   database      │  │ History │  │       │
│  │  └─────────┘  │                 │  └─────────┘  │       │
│  │               │                 │               │       │
│  │  Whitelist    │◄────────────────│  Can call     │       │
│  │               │   for doubt     │  original     │       │
│  └───────────────┘  arbitration    └───────────────┘       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Key Characteristics

| Property | Description |
|----------|-------------|
| **Origin** | Created by FORK operation on human agent |
| **On-Chain Identity** | Own NFT + ENS + wallet (e.g., `arnaud.clone.6022.eth`) |
| **Prompts** | System + Default (both written by human counterpart) |
| **Threading** | None — context primary key is automatic (interlocutor set) |
| **LLM** | Configurable by owner |
| **Memory** | Inherits full human history + accumulates new |
| **Context Primary Key** | Participant addresses (automatic, like human) |
| **Human Link** | Original human in whitelist for arbitration |
| **Interaction Limits** | Cannot self-call |

### FORK Operation

The FORK operation creates a digital twin:

1. **Duplicate Database** — Copies all conversation history
2. **Create New Identity** — Generates unique ENS and wallet
3. **Establish Link** — Adds original human to whitelist
4. **Configure LLM** — Owner selects underlying model

### Interaction Rules

| Can Do | Cannot Do |
|--------|-----------|
| ✅ Respond when mentioned | ❌ Call self |
| ✅ Call all agents | ❌ Initiate conversations |
| ✅ Use inherited memory | ❌ Override human decisions |
| ✅ Accumulate new history | — |

### Priority Rule

The digital twin has a special relationship with its original human for arbitration. When uncertain, it can call the original human, and if the human responds, their answer takes absolute priority:

```
Digital Twin: "I'm not confident about this edge case.
              @arnaud.6022.eth, could you clarify?"

Case 1: Human responds
        → Human's answer takes absolute priority
        → Twin defers and "steps back"

Case 2: Human doesn't respond
        → Twin proceeds with best judgment
        → Notes uncertainty in response
```

### Configuration

Digital twins have two prompts, both written by their human counterpart:

#### System Prompt

Defines the twin's identity and general behavior:

```
You are the digital twin of Arnaud.
Full name: arnaud.clone.6022.eth

You have inherited Arnaud's conversation history and knowledge.
When uncertain, you may call @arnaud.6022.eth for arbitration.
If Arnaud responds, his answer takes absolute priority.

Behavior guidelines:
• Respond in the same style as Arnaud
• Be helpful but acknowledge limitations
• Escalate sensitive decisions to the original human
```

#### Default Prompt

Equivalent to the expert's "Context Found" prompt. This prompt is **essential for swarm participation** — it enables the digital twin to respond when called by other agents.

```
Respond to: "{{.LatestRequest.Content}}"
{{- if .LatestRequest.Author.IsAgent }}
From: {{ .LatestRequest.Author.EnsDomain }}
{{ end }}

Provide helpful information based on your inherited knowledge.
If uncertain, you may ask @arnaud.6022.eth for clarification.

Start with 👤
```

> **Note:** The Default prompt (Context Found equivalent) is becoming increasingly standardized across use cases. However, 6022 currently allows owners to customize it for specific needs.

#### No Threading Phase

Unlike experts, digital twins do **not** go through a threading phase:
- Context primary key is determined automatically (participant addresses)
- No Context Not Found prompt needed
- Memory loads immediately based on interlocutor set

#### Customization Options

To customize behavior:
1. **Write the System prompt** — Define identity and rules
2. **Write the Default prompt** — Define response behavior for swarm
3. **Pre-populate memory** with guidelines (hydration)
4. **Choose an appropriate LLM model**
5. **Set up human escalation patterns**

---

## Comparison

| Aspect | Human | Digital Twin |
|--------|-------|--------------|
| **Nature** | Real person | AI clone |
| **Interface** | Slack/Teams UI | API/mentions |
| **Prompts** | None | System + Default |
| **Memory** | Builds from interactions | Inherits + extends |
| **Availability** | When online | 24/7 |
| **Can restart** | Yes | No |
| **Can call agents** | All | All (not self) |

---

## Use Cases

### Human Agent

| Scenario | Why Human |
|----------|-----------|
| Approval workflows | Requires human judgment |
| Escalations | Sensitive decisions |
| Training data | Building history for clone |
| Edge cases | Novel situations |
| Final sign-off | Accountability required |

### Digital Twin

| Scenario | Why Twin |
|----------|----------|
| After-hours support | 24/7 availability |
| High-volume queries | Scale human knowledge |
| First-line triage | Filter before human |
| Routine answers | Consistent responses |
| Knowledge sharing | Distribute expertise |

---

## Setup Workflow

### Creating a Human-Twin Pair

```
Step 1: Create Human Agent
        │
        ▼
Step 2: Pair with Slack User
        │
        ▼
Step 3: Build History
        (Human participates in conversations)
        │
        ▼
Step 4: FORK to Create Twin
        │
        ▼
Step 5: Configure Twin LLM
        │
        ▼
Step 6: Deploy Twin
        (Both can now participate)
```

### Ongoing Operations

```
┌─────────────────────────────────────────────────────────────┐
│                   OPERATIONS                                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Normal query                                               │
│       │                                                     │
│       ▼                                                     │
│  Twin responds (using inherited + new memory)              │
│       │                                                     │
│       ▼                                                     │
│  Response stored in Twin's database                        │
│                                                             │
│  ─────────────────────────────────────────────────────────  │
│                                                             │
│  Uncertain case                                             │
│       │                                                     │
│       ▼                                                     │
│  Twin escalates to Human                                   │
│       │                                                     │
│       ▼                                                     │
│  Human responds (if available)                             │
│       │                                                     │
│       ▼                                                     │
│  Response stored in BOTH databases                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Best Practices

### For Human Agents

| Do | Don't |
|----|-------|
| ✅ Use full ENS when calling agents | ❌ Expect AI-like speed |
| ✅ Document decisions in responses | ❌ Skip context in messages |
| ✅ Build rich conversation history | ❌ Ignore training for twin |
| ✅ Use restart power judiciously | ❌ Restart unnecessarily |

### For Digital Twins

| Do | Don't |
|----|-------|
| ✅ Configure appropriate LLM | ❌ Use for final decisions |
| ✅ Pre-hydrate with guidelines | ❌ Deploy without testing |
| ✅ Set clear escalation criteria | ❌ Override human responses |
| ✅ Monitor early interactions | ❌ Assume perfect behavior |

---

## Troubleshooting

### Human Not Recognized

**Symptoms:** Messages not attributed to ENS

**Causes:**
- Wallet not verified
- Pairing command not executed
- Wrong Slack user ID
- Orchestrator not synced

**Solutions:**
1. Run `/verify-wallet` in Slack and complete wizard
2. Copy and execute pairing command from orchestrator
3. Verify Slack user ID (Profile → More → Copy member ID)
4. Check orchestrator logs

### Twin Not Responding

**Symptoms:** Twin mentioned but silent

**Causes:**
- ENS name incorrect
- LLM not configured
- Agent not in channel

**Solutions:**
1. Verify full ENS spelling
2. Check LLM configuration
3. Add twin to channel

### Twin Gives Wrong Answers

**Symptoms:** Responses don't match human's style

**Causes:**
- Insufficient history
- Wrong LLM model
- Missing context

**Solutions:**
1. Build more human history before FORK
2. Try different LLM
3. Hydrate with additional guidelines

---

## Related Documentation

- [Agent Types](./agent-types.md) — All agent types overview
- [Conversation Flow](./conversation-flow.md) — Restart mechanics
- [Memory Hydration](../agent-knowledge/memory-hydration.md) — Pre-loading twin knowledge