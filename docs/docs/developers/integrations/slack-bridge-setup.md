---
sidebar_position: 2
title: Slack Bridge Setup
---

# Slack Bridge Setup

This page is for **whoever deploys the Slack bridge for a workspace**. It covers the Slack app manifest, OAuth scopes, event subscriptions, and common pitfalls. If you just want to use slash commands in an already-wired-up channel, the [user-facing Slack docs](../../users/slack/intro.md) are what you want instead.

The bridge itself is an HTTP service that lives in `channels-and-integrations/slack-bridge/`. Before any of this manifest work matters, you need the bridge running and reachable from Slack over HTTPS (see the repo README for deployment, Vault/Redis dependencies, and required environment variables).

## Complete Manifest

Copy this manifest and replace the URL placeholder with your bridge's public hostname:

```json
{
  "display_information": {
    "name": "6022",
    "description": "AI Agent Team Collaboration",
    "background_color": "#000000"
  },
  "features": {
    "bot_user": {
      "display_name": "6022",
      "always_online": false
    },
    "slash_commands": [
      {
        "command": "/add-agent",
        "url": "https://YOUR-SLACK-BRIDGE.up.railway.app/command-handler/add-agent-to-channel",
        "description": "Add an agent to the channel",
        "usage_hint": "[collection address] [token ID]",
        "should_escape": false
      },
      {
        "command": "/remove-agent",
        "url": "https://YOUR-SLACK-BRIDGE.up.railway.app/command-handler/remove-agent-from-channel",
        "description": "Remove an agent from the channel",
        "usage_hint": "[collection address] [token ID]",
        "should_escape": false
      },
      {
        "command": "/list-agents",
        "url": "https://YOUR-SLACK-BRIDGE.up.railway.app/command-handler/list-agents-in-channel",
        "description": "List all agents in the channel",
        "should_escape": false
      },
      {
        "command": "/verify-wallet",
        "url": "https://YOUR-SLACK-BRIDGE.up.railway.app/command-handler/verify-wallet-ownership",
        "description": "Verify your wallet ownership",
        "should_escape": false
      },
      {
        "command": "/pair-slack-user-to-human-agent",
        "url": "https://YOUR-SLACK-BRIDGE.up.railway.app/command-handler/pair-slack-user-to-human-agent",
        "description": "Pair a Slack user to a human agent",
        "usage_hint": "[slack user ID] [collection address] [token ID]",
        "should_escape": false
      }
    ]
  },
  "oauth_config": {
    "scopes": {
      "bot": [
        "app_mentions:read",
        "channels:history",
        "channels:join",
        "channels:read",
        "chat:write",
        "chat:write.customize",
        "chat:write.public",
        "commands",
        "files:read",
        "groups:history",
        "groups:read",
        "groups:write",
        "links:read",
        "links:write",
        "metadata.message:read",
        "users:read",
        "users.profile:read",
        "users:read.email",
        "mpim:read"
      ]
    }
  },
  "settings": {
    "event_subscriptions": {
      "request_url": "https://YOUR-SLACK-BRIDGE.up.railway.app/event-handler",
      "bot_events": [
        "channel_created",
        "channel_deleted",
        "channel_archive",
        "channel_left",
        "channel_id_changed",
        "message.channels",
        "message.groups"
      ]
    },
    "org_deploy_enabled": false,
    "socket_mode_enabled": false,
    "token_rotation_enabled": false
  }
}
```

:::warning Critical — `/verify-wallet` URL
The URL path for `/verify-wallet` is **`/command-handler/verify-wallet-ownership`**, not `/command-handler/verify-wallet`. The shorter form was documented in older revisions and still appears in some examples online — if you paste it into Slack verbatim, the bridge will return 404 and the command will silently break. Always use the `-ownership` suffix.
:::

---

## Manifest Sections Explained

### Display Information

```json
"display_information": {
  "name": "6022",
  "description": "AI Agent Team Collaboration",
  "background_color": "#000000"
}
```

| Field              | Description               |
| ------------------ | ------------------------- |
| `name`             | App name shown in Slack   |
| `description`      | Short description         |
| `background_color` | App icon background color |

---

### Bot User

```json
"bot_user": {
  "display_name": "6022",
  "always_online": false
}
```

The bot user posts messages on behalf of AI agents. Each agent's response shows the agent's name and avatar (via `chat:write.customize` — the bridge passes `username` and `icon_url` on every `chat.postMessage` call).

---

### Slash Commands

```json
"slash_commands": [
  {
    "command": "/add-agent",
    "url": "https://YOUR-BRIDGE/command-handler/add-agent-to-channel",
    "description": "Add an agent to the channel",
    "usage_hint": "[collection address] [token ID]",
    "should_escape": false
  }
]
```

| Command                           | URL Path                                          | Purpose                        |
| --------------------------------- | ------------------------------------------------- | ------------------------------ |
| `/add-agent`                      | `/command-handler/add-agent-to-channel`           | Add agent to channel           |
| `/remove-agent`                   | `/command-handler/remove-agent-from-channel`      | Remove agent                   |
| `/list-agents`                    | `/command-handler/list-agents-in-channel`         | List agents                    |
| `/verify-wallet`                  | `/command-handler/verify-wallet-ownership`        | Verify wallet ownership        |
| `/pair-slack-user-to-human-agent` | `/command-handler/pair-slack-user-to-human-agent` | Link Slack user to human agent |

:::warning URL format
Do **not** include `< >` brackets in URLs. Use:

```
https://slack-bridge.example.com/command-handler/add-agent-to-channel
```

Not:

```
https://<slack-bridge>.example.com/command-handler/add-agent-to-channel
```
:::

For the user-facing command syntax (parameters, examples, error messages), see [Slack Commands](../../users/slack/commands.md).

---

### OAuth Scopes

```json
"oauth_config": {
  "scopes": {
    "bot": [
      "app_mentions:read",
      "channels:history",
      "channels:join",
      "channels:read",
      "chat:write",
      "chat:write.customize",
      "chat:write.public",
      "commands",
      "files:read",
      "groups:history",
      "groups:read",
      "groups:write",
      "links:read",
      "links:write",
      "metadata.message:read",
      "users:read",
      "users.profile:read",
      "users:read.email",
      "mpim:read"
    ]
  }
}
```

#### Scope explanations

| Scope                   | Purpose                           |
| ----------------------- | --------------------------------- |
| `app_mentions:read`     | Detect when agents are @mentioned |
| `channels:history`      | Read public channel messages      |
| `channels:join`         | Auto-join public channels         |
| `channels:read`         | List channels                     |
| `chat:write`            | Post messages                     |
| `chat:write.customize`  | Post with custom username/avatar  |
| `chat:write.public`     | Post to channels bot isn't in     |
| `commands`              | Handle slash commands             |
| `files:read`            | Read file attachments             |
| `groups:history`        | Read private channel messages     |
| `groups:read`           | List private channels             |
| `groups:write`          | Post to private channels          |
| `links:read`            | See link previews                 |
| `links:write`           | Unfurl links                      |
| `metadata.message:read` | Read message metadata             |
| `users:read`            | Get user info                     |
| `users.profile:read`    | Read user profiles                |
| `users:read.email`      | Access user emails                |
| `mpim:read`             | Read group DMs (for future use)   |

:::note Unused scopes
The current bridge implementation does not call APIs that require `files:read`, `links:read`, `links:write`, `metadata.message:read`, `mpim:read`, or `app_mentions:read`. They're listed for future expansion and can be removed if you prefer to follow the principle of least privilege — the bridge won't break without them.
:::

---

### Event Subscriptions

```json
"event_subscriptions": {
  "request_url": "https://YOUR-BRIDGE/event-handler",
  "bot_events": [
    "channel_created",
    "channel_deleted",
    "channel_archive",
    "channel_left",
    "channel_id_changed",
    "message.channels",
    "message.groups"
  ]
}
```

#### Event types

| Event                | Triggered when                    |
| -------------------- | --------------------------------- |
| `message.channels`   | Message posted in public channel  |
| `message.groups`     | Message posted in private channel |
| `channel_created`    | New channel created               |
| `channel_deleted`    | Channel deleted                   |
| `channel_archive`    | Channel archived                  |
| `channel_left`       | Bot removed from channel          |
| `channel_id_changed` | Channel ID changed (rare)         |

:::note Message subtypes
Subscribing to `message.channels` and `message.groups` automatically delivers all message subtypes — including `message_changed` and `message_deleted`. The bridge handles edits and deletions in addition to new messages; you do not need to add anything to `bot_events` for that.
:::

:::tip Private channel support
To receive messages from private channels, **both** conditions must be met:

1. `message.groups` is in `bot_events`
2. The bot is invited to the private channel
:::

---

## Updating the Manifest

When you update the manifest:

1. Go to **App Manifest** in Slack app settings
2. Paste the updated JSON
3. Click **Save Changes**
4. If permissions changed, go to **Install App** → **Reinstall to Workspace**

After reinstalling:

- The Bot Token usually stays the same
- New permissions take effect immediately
- No need to restart or reconfigure the bridge service

---

## Environment-Specific Manifests

Use separate Slack apps for dev and prod. Each points its slash command and event URLs at a different bridge deployment.

### Development (Amoy Testnet)

```json
{
  "display_information": {
    "name": "6022 DEV"
  },
  "features": {
    "slash_commands": [
      {
        "command": "/add-agent",
        "url": "https://<your-dev-slack-bridge-host>/command-handler/add-agent-to-channel"
      }
    ]
  },
  "settings": {
    "event_subscriptions": {
      "request_url": "https://<your-dev-slack-bridge-host>/event-handler"
    }
  }
}
```

### Production (Polygon Mainnet)

```json
{
  "display_information": {
    "name": "6022"
  },
  "features": {
    "slash_commands": [
      {
        "command": "/add-agent",
        "url": "https://<your-prod-slack-bridge-host>/command-handler/add-agent-to-channel"
      }
    ]
  },
  "settings": {
    "event_subscriptions": {
      "request_url": "https://<your-prod-slack-bridge-host>/event-handler"
    }
  }
}
```

---

## Common Mistakes

| Mistake                           | Fix                                                              |
| --------------------------------- | ---------------------------------------------------------------- |
| Including `< >` in URLs           | Remove brackets: `https://bridge.app` not `https://<bridge>.app` |
| `/verify-wallet` URL missing `-ownership` suffix | Use `/command-handler/verify-wallet-ownership`      |
| Wrong event handler path          | Use `/event-handler` not `/events`                               |
| Missing `message.groups` scope    | Add it to support private channels                               |
| Forgetting to reinstall           | After scope changes, reinstall the app                           |

---

## See also

- [User-facing Slack integration overview](../../users/slack/intro.md) — architecture, prerequisites, channel types
- [Slash commands reference](../../users/slack/commands.md) — syntax, parameters, error messages
