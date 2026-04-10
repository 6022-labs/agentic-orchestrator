---
sidebar_position: 1
---

# Slack Integration Overview

Slack is where agents talk to humans. A Slack channel becomes a meeting room: humans post messages, the facilitator coordinates the experts, experts respond when @-mentioned, and every conversation becomes a Slack thread. Human agents (Slack users paired to a human-role NFT) can restart concluded threads by posting a new message.

## How a message flows

```
User posts in a channel
       │
       ▼
Slack webhook → Slack Bridge → Orchestrator → Facilitator
                                                    │
                                                    ├─ @expert-agent ─▶ Expert responds
                                                    │
                                                    ▼
                                              Facilitator synthesises
                                                    │
                                                    ▼
                                          Response appears in the thread
```

Each 6022 conversation corresponds to exactly one Slack thread. When the facilitator concludes, the thread is marked complete. Only a human can reopen it by posting a new message.

## Channel types

| Channel | Supported |
|---|---|
| Public channels | ✅ (via `message.channels` event) |
| Private channels | ✅ (via `message.groups` event — the bot must be invited to the channel) |
| Direct messages | ❌ — ignored downstream because the channel is never registered via `/add-agent` |
| Group DMs | ❌ — same reason |

## Prerequisites

You need:

- Slack workspace admin access
- A running 6022 orchestrator
- A deployed Slack Bridge (Railway, Docker, etc.)
- At least one agent collection

The bridge service itself depends on **HashiCorp Vault** (for secret storage) and **Redis** (for caching Slack user/channel sync state).

### Required environment variables

| Variable | Where to find it |
|---|---|
| `SLACK_APP_ID` | Slack app **Basic Information** → App Credentials. **Required at startup.** |
| `SLACK_CLIENT_SECRET` | Slack app **Basic Information** → App Credentials |
| `SLACK_SIGNING_SECRET` | Slack app **Basic Information** → App Credentials |
| `SLACK_BOT_TOKEN` | Slack app **OAuth & Permissions** → Bot User OAuth Token (`xoxb-...`) |
| `PUBLIC_KEY` | Public key used by the verify-wallet flow |
| `ORCHESTRATOR_FRONT_BASE_URL` | Public URL of your orchestrator frontend |
| `ORCHESTRATOR_API_BASE_URL` | Internal URL of the gateway service |
| `ORCHESTRATOR_AGENT_COMPUTE_API_BASE_URL` | Internal URL of the agent-compute service |

## Setting up the bridge

**Configuring the Slack app itself** (manifest JSON, OAuth scopes, event subscriptions, bridge deployment) is a deployer task, not an end-user one. Full walkthrough: [Slack Bridge Setup](../../developers/integrations/slack-bridge-setup.md) in the developer docs.

## Using slash commands

Once the bridge is running and agents are minted, you drive everything through slash commands in Slack: `/add-agent`, `/remove-agent`, `/list-agents`, `/verify-wallet`, `/pair-slack-user-to-human-agent`. See **[Slash Commands](./commands.md)** for the full reference.
