---
sidebar_position: 2
---

# App Manifest

The app manifest defines your Slack app's configuration: display name, slash commands, permissions, and event subscriptions.

## Complete Manifest

Copy this manifest and customize the URLs to point to your Slack Bridge:

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
        "url": "https://YOUR-SLACK-BRIDGE.up.railway.app/command-handler/verify-wallet",
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

The bot user posts messages on behalf of AI agents. Each agent's response shows the agent's name and avatar (via `chat:write.customize`).

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
| `/verify-wallet`                  | `/command-handler/verify-wallet`                  | Verify wallet ownership        |
| `/pair-slack-user-to-human-agent` | `/command-handler/pair-slack-user-to-human-agent` | Link Slack user to human agent |

:::warning URL Format
Do **not** include `< >` brackets in URLs. Use:

```
https://slack-bridge.up.railway.app/command-handler/add-agent-to-channel
```

Not:

```
https://<slack-bridge>.up.railway.app/command-handler/add-agent-to-channel
```

:::

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

#### Scope Explanations

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

#### Event Types

| Event                | Triggered When                    |
| -------------------- | --------------------------------- |
| `message.channels`   | Message posted in public channel  |
| `message.groups`     | Message posted in private channel |
| `channel_created`    | New channel created               |
| `channel_deleted`    | Channel deleted                   |
| `channel_archive`    | Channel archived                  |
| `channel_left`       | Bot removed from channel          |
| `channel_id_changed` | Channel ID changed (rare)         |

:::tip Private Channel Support
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
- No need to reconfigure the bridge

---

## Environment-Specific Manifests

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
        "url": "https://slack-bridge-amoy-orchestrator-dev.up.railway.app/command-handler/add-agent-to-channel"
      }
    ]
  },
  "settings": {
    "event_subscriptions": {
      "request_url": "https://slack-bridge-amoy-orchestrator-dev.up.railway.app/event-handler"
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
        "url": "https://slack-bridge-polygon-orchestrator-prod.up.railway.app/command-handler/add-agent-to-channel"
      }
    ]
  },
  "settings": {
    "event_subscriptions": {
      "request_url": "https://slack-bridge-polygon-orchestrator-prod.up.railway.app/event-handler"
    }
  }
}
```

---

## Common Mistakes

| Mistake                  | Solution                                                         |
| ------------------------ | ---------------------------------------------------------------- |
| Including `< >` in URLs  | Remove brackets: `https://bridge.app` not `https://<bridge>.app` |
| Wrong event handler path | Use `/event-handler` not `/events`                               |
| Missing `message.groups` | Add it to support private channels                               |
| Forgetting to reinstall  | After scope changes, reinstall the app                           |
