---
sidebar_position: 3
---

# Slash Commands

Slash commands let you manage agents directly from Slack. Type a command in any channel where the bot is installed.

## Agent Management Commands

### /add-agent

Add an AI agent to the current channel.

**Syntax:**

```
/add-agent [collection address] [token ID]
```

**Example:**

```
/add-agent 0xb533A24430c3d66FAEdE1516d07913089EFf0e7e 1
```

**Parameters:**
| Parameter | Description | Example |
|-----------|-------------|---------|
| Collection Address | The NFT collection address | `0xb533A244...` |
| Token ID | The agent's token ID in the collection | `1`, `2`, `3` |

**Response:**

```
✅ Agent elen.mycompany.137.protocol6022.eth has been added to #customer-support
```

:::tip Copy Command from Orchestrator
In the orchestrator dashboard, click on an agent → Actions → **Copy add agent command**. This copies the complete command ready to paste in Slack.
:::

---

### /remove-agent

Remove an AI agent from the current channel.

**Syntax:**

```
/remove-agent [collection address] [token ID]
```

**Example:**

```
/remove-agent 0xb533A24430c3d66FAEdE1516d07913089EFf0e7e 1
```

**Response:**

```
✅ Agent elen.mycompany.137.protocol6022.eth has been removed from #customer-support
```

---

### /list-agents

Show all agents currently in the channel.

**Syntax:**

```
/list-agents
```

**Response:**

```
Agents in #customer-support:
• elen.mycompany.137.protocol6022.eth (facilitator)
• claudine.mycompany.137.protocol6022.eth (expert)
• victor.mycompany.137.protocol6022.eth (expert)
```

If no agents are present:

```
No agents in this channel.
```

---

## Human Agent Commands

### /verify-wallet

Verify your wallet ownership. This is required before pairing a Slack user to a human agent.

**Syntax:**

```
/verify-wallet
```

**Response:**

```
🔐 To verify your wallet, please sign this message:
[Link to signing page]

This proves you own the wallet connected to your agents.
```

**Process:**

1. Type `/verify-wallet` in any channel
2. Click the link
3. Connect your Web3 wallet (MetaMask, etc.)
4. Sign the verification message
5. Your Slack user is now linked to your wallet address

:::info Why Verify?
Verification proves you own the wallet that holds agent NFTs. This security step ensures only the actual owner can pair Slack users to their human agents.
:::

---

### /pair-slack-user-to-human-agent

Link a Slack user to a human agent. The human agent will then participate in swarm conversations when @mentioned.

**Prerequisites:**

1. You must first run `/verify-wallet` to prove wallet ownership
2. You must be the owner of the human agent NFT

**Syntax:**

```
/pair-slack-user-to-human-agent [slack user ID] [collection address] [token ID]
```

**Example:**

```
/pair-slack-user-to-human-agent U09BDFS3VGA 0xb533A24430c3d66FAEdE1516d07913089EFf0e7e 9
```

**Parameters:**
| Parameter | Description | How to Find |
|-----------|-------------|-------------|
| Slack User ID | The user's Slack ID | Click user profile → More → Copy member ID |
| Collection Address | Agent collection address | From orchestrator dashboard |
| Token ID | Human agent's token ID | From orchestrator dashboard |

**Response:**

```
✅ Slack user @john.doe is now paired with arnaud.mycompany.137.protocol6022.eth
```

---

## Getting Agent Information

### Finding Collection Address and Token ID

From the orchestrator dashboard:

1. Go to **Agents** page
2. Click on the agent
3. Look in the URL or agent details:
   - Collection: `0xb533A24430c3d66FAEdE1516d07913089EFf0e7e`
   - Token ID: `1`

Or use the **Actions** menu:

- **Copy add agent command** → Gives you the full command
- **Copy remove agent command** → Gives you the removal command
- **Copy pair slack user to human agent command** → For human agents

### Finding Slack User ID

1. Click on the user in Slack
2. Click **More** (three dots)
3. Click **Copy member ID**

Format: `U09BDFS3VGA` (starts with `U`)

---

## Command Workflow

### Setting Up a Channel

```
1. Create channel #sales-team

2. Add facilitator:
   /add-agent 0xb533A244... 1

3. Add experts:
   /add-agent 0xb533A244... 2
   /add-agent 0xb533A244... 3

4. Verify agents:
   /list-agents
   → Shows: elen (facilitator), claudine (expert), victor (expert)

5. Post a message:
   "What's our Q4 sales forecast?"

6. Agents respond in thread!
```

### Adding a Human Agent

```
1. Owner verifies wallet:
   /verify-wallet
   → Signs message with Web3 wallet

2. Pair Slack user to human agent:
   /pair-slack-user-to-human-agent U09BDFS3VGA 0xb533A244... 9
   → John is now arnaud.mycompany.137.protocol6022.eth

3. In conversations, agents can @mention John:
   "Let me ask @arnaud.mycompany.137.protocol6022.eth for clarification..."

4. John sees the message and can respond
   → His response is recorded in the swarm conversation
```

---

## Error Messages

The bridge surfaces raw error strings through Slack and the exact text may change between versions. The most common failure modes:

| Failure mode | Likely cause | What to try |
| --- | --- | --- |
| Command returns "Failed to handle command" | Generic backend error — the agent could not be added/removed/listed | Verify the collection address and token ID against the orchestrator dashboard |
| Pairing fails with a wallet error | Caller hasn't verified wallet ownership, or doesn't own the agent NFT | Run `/verify-wallet` first and confirm the agent is in your wallet |
| Add/remove fails silently | Bot is not in the channel | Invite the bot to the channel |
| `/pair-slack-user-to-human-agent` rejects the user ID | Wrong format — must be a raw Slack user ID like `U09BDFS3VGA`, not a `@mention` | Copy the ID from the user's Slack profile (More → Copy member ID) |

---

## Best Practices

### Channel Organization

```
Recommended setup:
├── #general (no agents - company-wide chat)
├── #customer-support
│   ├── Facilitator: support-coordinator
│   ├── Expert: crm-specialist
│   └── Expert: product-expert
├── #sales-team
│   ├── Facilitator: sales-coordinator
│   └── Expert: pricing-expert
└── #legal-review
    ├── Facilitator: legal-coordinator
    ├── Expert: contract-expert
    └── Expert: gdpr-expert
```

### Agent Selection

- **One facilitator per channel** - Coordinates and concludes conversations
- **Multiple experts as needed** - Called when their expertise is required
- **Human agents for escalation** - Real humans when AI can't resolve

### Security

- The bridge does not enforce any role-based authorization on `/add-agent` or `/remove-agent` — any workspace member who can use slash commands in the channel can run them. Lock the bot down with channel membership instead of relying on a doc convention.
- Verify wallet ownership before pairing human agents
- Don't share collection addresses publicly if sensitive
