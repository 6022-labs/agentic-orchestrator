---
sidebar_position: 4
---

# Agent Configuration

Once an agent is created, you need to configure its backend resources: LLM, database, prompts, and MCP connections. All configuration is done through the orchestrator interface.

## Accessing Configuration

1. Navigate to **Dashboard** → **Agents**
2. Find your agent in the list
3. Click the **⋮** menu
4. Select the configuration option you need

---

## LLM Settings

Each agent needs an LLM to generate responses.

### Configuration Options

| Field        | Description          | Example                          |
| ------------ | -------------------- | -------------------------------- |
| **Provider** | LLM service provider | `openai`, `anthropic`, `mistral` |
| **API URL**  | Provider endpoint    | `https://api.openai.com/v1`      |
| **Model**    | Specific model       | `gpt-4o`, `claude-3-5-sonnet`    |
| **API Key**  | Authentication       | `sk-...`                         |

### Accessing LLM Settings

1. Click **⋮** → **"Manage LLM settings"**
2. Fill in the provider details
3. Save changes

:::tip Multi-Model Strategy
Different agents can use different LLMs. See [LLM Introduction](../llm/intro) for guidance on model diversity.
:::

---

## Database Configuration

Each agent has its own PostgreSQL database to store:

- Conversation history
- Context keys
- Threading information

### Database Parameters

| Field        | Description                    |
| ------------ | ------------------------------ |
| **Host**     | Database server address        |
| **Port**     | Connection port (usually 5432) |
| **Database** | Database name                  |
| **Username** | Role with access               |
| **Password** | Role password                  |

### Setting Up a New Database

For each agent, create a dedicated database role:

```sql
-- Create role for the agent
CREATE ROLE agent_claudine WITH LOGIN PASSWORD 'secure_password';

-- Create database
CREATE DATABASE claudine_db OWNER agent_claudine;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE claudine_db TO agent_claudine;
```

### Hosted Option

6022 can provide a hosted PostgreSQL instance on request:

- Host: `<provided by 6022>`
- Contact 6022 admin for credentials

---

## Prompts

Prompts define how your agent behaves. Different agent types use different prompts.

### Prompt Types by Role

| Agent Role       | Prompts Used                                                |
| ---------------- | ----------------------------------------------------------- |
| **Facilitator**  | System Prompt, Closure Prompt                               |
| **Expert**       | System Prompt, Thread Found Prompt, Thread Not Found Prompt |
| **Human**        | None (human responds)                                       |
| **Digital Twin** | System Prompt, Default Prompt                               |

### Accessing Prompt Configuration

1. Click **⋮** → **"Manage prompts"**
2. Edit the relevant prompts
3. Save changes

### System Prompt

The main prompt that defines agent behavior:

```
You are Claudine, an expert in customer relationship management.
Your role is to help the team find customer information and history.

When asked about a customer, search the CRM using available tools.
Always cite your sources and provide accurate information.
```

### Thread Found / Not Found (Experts)

**Thread Found Prompt**: Used when the agent identifies a context key

```
Context key detected: {{ .ContextKey }}
Loading relevant history for this {{ .ContextKey }}...
```

**Thread Not Found Prompt**: Used during the threading phase

```
Analyzing the conversation to identify the primary key...
Looking for: invoice number, contract ID, customer email, etc.
```

### Closure Prompt (Facilitators)

Used when the facilitator is ready to close the conversation:

```
Based on the information gathered, I am now confident to provide a final answer.
Confidence level: 8/10

[Summary of findings]

:end:
```

---

## Dynamic Variables

Prompts can include dynamic variables (Go `text/template` syntax) that are replaced at runtime — channel description, channel members, context keys, request content, author info, etc.

The complete reference (every variable, where it's available, and the full availability matrix by prompt type) lives in **[Prompt Dynamic Fields](../prompting/prompt-dynamic-fields.md)**. Don't duplicate it here — the canonical list is the one Phase B / Phase C maintains in that page.

---

## MCP Connections

MCP (Model Context Protocol) connects agents to external tools and data sources.

### Adding an MCP Connection

1. Click **⋮** → **"Manage MCPs"**
2. Click **"+ Add MCP Connection"**
3. Configure:

| Field              | Description              |
| ------------------ | ------------------------ |
| **Name**           | Identifier for this MCP  |
| **URL**            | MCP server endpoint      |
| **Authentication** | API Key or Custom Header |

### Example Configuration

```yaml
Name: kirha-mcp
URL: https://kirha-mcp-6022.up.railway.app/mcp
Authentication: Custom Header
  Header Name: x-kirha-api-key
  Header Value: your-api-key
```

See the [MCP Documentation](../mcp/intro) for detailed configuration.

---

## Context Management

For **Expert agents**, context keys determine how conversations are threaded.

### Viewing Contexts

1. Click **⋮** → **"Manage contexts"**
2. View all known context keys
3. Browse conversation history by key

### Context Key Examples

| Domain     | Context Key Example              |
| ---------- | -------------------------------- |
| Insurance  | Policy number: `POL-2024-123456` |
| E-commerce | Order ID: `ORD-789012`           |
| Support    | Ticket number: `TKT-456`         |
| Finance    | Invoice: `INV-2024-001`          |

### API Access to Contexts

```bash
# List all context keys
GET /agent/:agentCollectionAddress/:agentCollectionTokenId/contexts/keys

# Create a new context
POST /agent/:agentCollectionAddress/:agentCollectionTokenId/context
Body: { "key": "POL-2024-123456" }

# Append a message to a context (how memory hydration works)
POST /agent/:agentCollectionAddress/:agentCollectionTokenId/context/:contextKey/message
Body: { "message": "...", "role": "user" }
```

See the [Developer API Reference](../../developers/api-reference/intro.md) for the full endpoint documentation, parameter list, and code samples.

---

## Configuration Checklist

Before your agent is ready:

| Component                     | Status | Required For                 |
| ----------------------------- | ------ | ---------------------------- |
| ☐ LLM configured              |        | All agents                   |
| ☐ Database connected          |        | All agents                   |
| ☐ System prompt written       |        | Facilitators, Experts, Twins |
| ☐ Thread prompts written      |        | Experts only                 |
| ☐ Closure prompt written      |        | Facilitators only            |
| ☐ MCP connections (if needed) |        | Agents using external tools  |
| ☐ Added to Slack channel      |        | All agents                   |
