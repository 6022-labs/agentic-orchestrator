---
sidebar_position: 5
---

# Actions Menu

Each agent in the orchestrator has an **Actions menu** (⋮) that provides quick access to all management operations.

## Accessing the Menu

1. Navigate to **Dashboard** → **Agents**
2. Find your agent in the list
3. Click the **⋮** button on the right side

## Available Actions

### Blockchain Explorers

View your agent on various platforms:

| Action              | Description                            |
| ------------------- | -------------------------------------- |
| **Polygon Scan**    | View NFT on PolygonScan block explorer |
| **View on OpenSea** | See agent on OpenSea marketplace       |
| **View on Rarible** | See agent on Rarible marketplace       |

These links allow you to:

- Verify on-chain ownership
- View transaction history
- Access NFT metadata
- List agent for sale (marketplaces)

---

### Slack Commands

Copy ready-to-use Slack commands:

| Action                           | Command Generated                      |
| -------------------------------- | -------------------------------------- |
| **Copy add agent command**       | `/add-agent [collection] [tokenId]`    |
| **Copy remove agent command**    | `/remove-agent [collection] [tokenId]` |
| **Copy pair slack user command** | `/pair-user [collection] [tokenId]`    |

#### Add Agent Command

Adds the agent to the current Slack channel:

```
/add-agent 0x6ed8e93c10aa16031303d6153e0dd8ea3dbf870d 1
```

#### Remove Agent Command

Removes the agent from the current Slack channel:

```
/remove-agent 0x6ed8e93c10aa16031303d6153e0dd8ea3dbf870d 1
```

#### Pair User Command (Human Agents)

Links a Slack user to a human agent:

```
/pair-user 0x6ed8e93c10aa16031303d6153e0dd8ea3dbf870d 3
```

:::note Human Agents
The pair command is only relevant for **Human** agent types. It connects a real Slack user to their 6022 human agent identity.
:::

---

### Transfer

Transfer agent ownership to another wallet.

| Field                 | Description                 |
| --------------------- | --------------------------- |
| **Recipient Address** | Wallet address of new owner |

**Process:**

1. Click **"Transfer"**
2. Enter recipient wallet address
3. Confirm transaction in your wallet
4. Pay gas fees
5. Agent ownership transfers on-chain

:::warning Irreversible
Transfers are permanent. Once confirmed, you lose all access to the agent's configuration.
:::

---

### Manage LLM Settings

Configure the agent's Large Language Model.

| Setting  | Description                      |
| -------- | -------------------------------- |
| Provider | OpenAI, Anthropic, Mistral, etc. |
| API URL  | Provider endpoint                |
| Model    | Specific model name              |
| API Key  | Authentication credentials       |

See [LLM Documentation](../llm/intro) for detailed configuration.

---

### Manage MCPs

Configure Model Context Protocol connections.

| Setting        | Description              |
| -------------- | ------------------------ |
| Name           | MCP identifier           |
| URL            | Server endpoint          |
| Authentication | API Key or Custom Header |

See [MCP Documentation](../mcp/intro) for detailed configuration.

---

### Manage Prompts

Edit agent behavior through prompts.

| Prompt Type      | Used By            |
| ---------------- | ------------------ |
| System Prompt    | All agent types    |
| Thread Found     | Expert agents      |
| Thread Not Found | Expert agents      |
| Closure          | Facilitator agents |
| Default          | Digital Twins      |

See [Prompting Documentation](../prompting/intro) for best practices.

---

### Manage Contexts

View and manage agent memory contexts.

| Feature           | Description                      |
| ----------------- | -------------------------------- |
| View context keys | List all known primary keys      |
| Browse history    | View conversations by context    |
| Delete context    | Remove a context and its history |

**Use cases:**

- Debug threading issues
- Clean up test data
- Review agent memory

---

## Action Availability by Role

Not all actions are available for all agent types:

| Action                   | Facilitator | Expert | Human | Digital Twin |
| ------------------------ | :---------: | :----: | :---: | :----------: |
| Blockchain explorers     |     ✅      |   ✅   |  ✅   |      ✅      |
| Copy add/remove commands |     ✅      |   ✅   |  ✅   |      ✅      |
| Copy pair user command   |     ❌      |   ❌   |  ✅   |      ❌      |
| Transfer                 |     ✅      |   ✅   |  ✅   |      ✅      |
| Manage LLM               |     ✅      |   ✅   |  ❌   |      ✅      |
| Manage MCPs              |     ✅      |   ✅   |  ❌   |      ✅      |
| Manage Prompts           |     ✅      |   ✅   |  ❌   |      ✅      |
| Manage Contexts          |     ✅      |   ✅   |  ✅   |      ✅      |

---

## Quick Reference

### Most Common Operations

| I want to...               | Action                                  |
| -------------------------- | --------------------------------------- |
| Add agent to Slack         | Copy add agent command → paste in Slack |
| Change LLM provider        | Manage LLM settings                     |
| Update agent behavior      | Manage Prompts                          |
| Connect external tools     | Manage MCPs                             |
| Debug conversation history | Manage Contexts                         |
| Give agent to someone else | Transfer                                |
| See agent on blockchain    | Polygon Scan                            |

### Keyboard Shortcuts

Currently, no keyboard shortcuts are available. All actions require mouse interaction through the ⋮ menu.
