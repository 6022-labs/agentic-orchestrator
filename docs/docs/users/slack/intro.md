---
sidebar_position: 1
---

# Slack Integration Overview

The Slack integration allows AI agents to participate in conversations within your Slack workspace. A Slack channel becomes a meeting room where humans and AI agents collaborate.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         YOUR SLACK WORKSPACE                         │
│                                                                       │
│   #customer-support          #sales-team           #legal-review     │
│   ┌─────────────────┐       ┌─────────────────┐   ┌─────────────────┐│
│   │ 👤 Human        │       │ 👤 Human        │   │ 👤 Human        ││
│   │ 🤖 Facilitator  │       │ 🤖 Facilitator  │   │ 🤖 Facilitator  ││
│   │ 🤖 CRM Expert   │       │ 🤖 Sales Expert │   │ 🤖 Legal Expert ││
│   │ 🤖 Tech Expert  │       │                 │   │ 🤖 GDPR Expert  ││
│   └────────┬────────┘       └────────┬────────┘   └────────┬────────┘│
│            │                         │                     │          │
└────────────┼─────────────────────────┼─────────────────────┼──────────┘
             │                         │                     │
             ▼                         ▼                     ▼
      ┌──────────────────────────────────────────────────────────────┐
      │                       SLACK BRIDGE                            │
      │                                                               │
      │  • Receives Slack events (messages, channel changes)         │
      │  • Sends agent responses back to Slack                       │
      │  • Handles slash commands (/add-agent, /verify-wallet)       │
      │  • Verifies request signatures for security                  │
      └─────────────────────────────┬─────────────────────────────────┘
                                    │
                                    ▼
      ┌──────────────────────────────────────────────────────────────┐
      │                       ORCHESTRATOR                            │
      │                                                               │
      │  • Routes messages to appropriate agents                     │
      │  • Manages conversation threads                              │
      │  • Stores conversation history                               │
      │  • Calls LLM providers (OpenAI, Anthropic, etc.)            │
      └──────────────────────────────────────────────────────────────┘
```

## How It Works

### 1. A Slack Channel = A Meeting Room

When you create a Slack channel and add agents to it, you're essentially setting up a meeting room where:
- **Humans** post messages and can restart closed conversations
- **Facilitator agent** coordinates the discussion and decides when to conclude
- **Expert agents** are called via @mentions to provide specialized knowledge
- **Human agents** (paired Slack users) can be addressed by their ENS domain

### 2. Message Flow

```
User posts in #channel
       │
       ▼
Slack sends event to Bridge
       │
       ▼
Bridge forwards to Orchestrator
       │
       ▼
Facilitator receives message
       │
       ├──► Calls @expert-agent if needed
       │         │
       │         ▼
       │    Expert responds
       │         │
       ◄─────────┘
       │
       ▼
Facilitator synthesizes and responds
       │
       ▼
Response appears in Slack thread
```

### 3. Conversations as Slack Threads

Each conversation in 6022 corresponds to a **Slack thread**:
- The initial message starts a new thread
- All agent responses appear as replies in that thread
- When the facilitator concludes, the thread is marked as completed
- Human agents can restart a closed thread by posting a new message

---

## Prerequisites

Before setting up the Slack integration, you need:

| Requirement | Description |
|-------------|-------------|
| **Slack Workspace** | Admin access to install apps |
| **6022 Orchestrator** | A running orchestrator instance |
| **Slack Bridge** | Deployed bridge service (Railway, Docker, etc.) |
| **Agent Collection** | At least one collection with agents created |
| **HashiCorp Vault** | Both bridge services require `VAULT_ADDR` and `VAULT_TOKEN` to read/write secrets |
| **Redis** | Both bridge services require Redis for caching Slack user/channel synchronization |

The Slack-side environment variables the bridge expects:

| Variable | Where to find it |
|---|---|
| `SLACK_APP_ID` | Slack app **Basic Information** → App Credentials. **Required at startup — the bridge will refuse to boot without it.** |
| `SLACK_CLIENT_SECRET` | Slack app **Basic Information** → App Credentials |
| `SLACK_SIGNING_SECRET` | Slack app **Basic Information** → App Credentials |
| `SLACK_BOT_TOKEN` | Slack app **OAuth & Permissions** → Bot User OAuth Token (`xoxb-...`) |
| `PUBLIC_KEY` | Public key used by the verify-wallet flow |
| `ORCHESTRATOR_FRONT_BASE_URL` | Public URL of your orchestrator frontend (used for the verify-wallet redirect) |
| `ORCHESTRATOR_API_BASE_URL` | Internal URL of the gateway service |
| `ORCHESTRATOR_AGENT_COMPUTE_API_BASE_URL` | Internal URL of the agent-compute service |

---

## Key Concepts

### Slack App

A custom Slack app that you create and install in your workspace. It:
- Listens for messages in channels
- Provides slash commands for agent management
- Posts messages on behalf of AI agents

### Slack Bridge

A backend service that:
- Receives events from Slack (via webhook)
- Validates request signatures for security
- Forwards events to the orchestrator
- Sends agent responses back to Slack

### Bot Token

The authentication token (`xoxb-...`) that allows the Slack app to:
- Read messages from channels
- Post messages as the bot
- Access user information

### Slash Commands

Special commands typed in Slack that trigger actions:
- `/add-agent` - Add an agent to the current channel
- `/remove-agent` - Remove an agent from the channel
- `/list-agents` - Show all agents in the channel
- `/verify-wallet` - Verify wallet ownership for human agent pairing
- `/pair-slack-user-to-human-agent` - Link a Slack user to a human agent

---

## Channel Types Support

| Channel Type | Supported | Required Event |
|--------------|:---------:|----------------|
| Public channels | ✅ | `message.channels` |
| Private channels | ✅ | `message.groups` |
| Direct messages | ❌ | The bridge does not currently process DMs — messages in 1:1 channels are ignored downstream because the channel is never registered via `/add-agent`. |
| Group DMs | ❌ | Same as above — group DM channels are not registered, so the bridge has nothing to route to. |

:::tip Private Channels
To support private channels, ensure `message.groups` is included in your app manifest's `bot_events` array. The bot must also be invited to the private channel.
:::

---

## What You'll Learn

In this section, you'll learn how to:

1. **[Use Slash Commands](./commands)** — manage agents in channels once the bridge is wired up

Configuring the Slack app itself (manifest JSON, OAuth scopes, event subscriptions, bridge deployment) is a deployer task rather than an end-user one — see [Slack Bridge Setup](../../developers/integrations/slack-bridge-setup.md) in the developer docs.

---

## Quick Start

If you want to get started quickly:

```bash
# 1. Create a Slack app at https://api.slack.com/apps

# 2. Copy the app manifest (see developers/integrations/slack-bridge-setup)

# 3. Install the app to your workspace

# 4. Note these credentials:
#    - Client ID
#    - Client Secret  
#    - Signing Secret
#    - Bot Token (xoxb-...)

# 5. Configure your Slack Bridge with these credentials

# 6. In Slack, use /add-agent to add agents to a channel
```