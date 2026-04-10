---
sidebar_position: 1
---

# Introduction to MCP

MCP (Model Context Protocol) is an open standard that allows AI agents to connect to external tools and data sources. In 6022 Agentic Teams, MCP enables your agents to access CRMs, databases, file conversion services, web search, and any other API-enabled system.

## What is MCP?

MCP provides a standardized way for LLMs to:

- **Discover tools** available on a server
- **Call tools** with structured parameters
- **Receive results** in a consistent format

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│    Agent    │────▶│  MCP Server │────▶│  External   │
│   (LLM)     │◀────│             │◀────│   System    │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                    JSON-RPC 2.0
                    "tools/list"
                    "tools/call"
```

## Architecture in 6022

In the 6022 framework, each agent can connect to multiple MCP servers:

```
┌─────────────────────────────────────────────────────┐
│                      AGENT                          │
│       claudine.mycompany.137.protocol6022.eth       │
└───────────────────────┬─────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        ▼               ▼               ▼
   ┌─────────┐    ┌─────────┐    ┌─────────┐
   │ MCP #1  │    │ MCP #2  │    │ MCP #3  │
   │ Kirha   │    │ Stripe  │    │Markitdown│
   │(search) │    │(payments│    │(convert)│
   └─────────┘    └─────────┘    └─────────┘
```

## LLM Compatibility

The orchestrator distinguishes between **native MCP routing** (the orchestrator forwards your MCP server definitions directly to the provider) and **function-call loop routing** (the orchestrator runs its own loop and asks the model to emit tool calls). The single source of truth lives in **[LLM Introduction → Supported Providers](../llm/intro.md#supported-providers)** — that page is the canonical compatibility reference.

| LLM Provider                                                     | MCP Routing                                                             |
| ---------------------------------------------------------------- | ----------------------------------------------------------------------- |
| **OpenAI**                                                       | ✅ Native (via `/v1/responses`)                                         |
| **Anthropic**                                                    | ✅ Native (custom-header configuration required)                        |
| **AWS Bedrock**                                                  | ⚙️ Function-call loop (Converse API)                                    |
| **xAI** / **Mistral** / **Deepseek** / **LibertAI** / **Ollama** | ⚙️ Function-call loop — depends on the model's function-calling support |

If you configure MCP on an agent whose model does not actually support function calling, the orchestrator will still attempt to invoke the loop; tools will simply never be selected by the model.

## Use Cases

### Dynamic Data Access

Connect agents to live data sources:

| Use Case     | MCP Server   | Data Accessed                     |
| ------------ | ------------ | --------------------------------- |
| CRM queries  | HubSpot MCP  | Customer records, deals, contacts |
| Payment info | Stripe MCP   | Invoices, subscriptions, balances |
| Scheduling   | Calendar MCP | Appointments, availability        |
| Inventory    | Custom MCP   | Stock levels, product catalog     |

### Document Processing

Convert and extract information from files:

| Use Case       | MCP Server       | Capability                        |
| -------------- | ---------------- | --------------------------------- |
| PDF extraction | Markitdown       | Convert PDF/DOCX/XLSX to Markdown |
| Image analysis | Kirha            | Extract metadata, analyze content |
| Web content    | Kirha Web Search | Fetch and summarize web pages     |

### AI-Enhanced Research

Leverage external AI services:

| Use Case       | MCP Server   | Capability                     |
| -------------- | ------------ | ------------------------------ |
| Web search     | Kirha        | Real-time internet search      |
| Crypto data    | Kirha Crypto | Prices, volumes, market data   |
| Knowledge base | Custom MCP   | Company-specific documentation |

## How MCP Works (Technical)

MCP uses **JSON-RPC 2.0** over HTTP(S). The orchestrator performs two main operations:

### 1. Tool Discovery

When an agent is invoked, the orchestrator first asks the MCP server what tools are available:

```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "params": {},
  "id": 1
}
```

The server responds with available tools:

```json
{
  "jsonrpc": "2.0",
  "result": {
    "tools": [
      {
        "name": "search_customers",
        "description": "Search for customers by name or email",
        "inputSchema": {
          "type": "object",
          "properties": {
            "query": { "type": "string" }
          }
        }
      }
    ]
  },
  "id": 1
}
```

### 2. Tool Execution

When the LLM decides to use a tool, the orchestrator calls it:

```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "search_customers",
    "arguments": { "query": "john.doe@example.com" }
  },
  "id": 2
}
```

## Key Principles

### One Agent, Multiple MCPs

An agent can connect to several MCP servers simultaneously. The LLM sees all available tools and decides which to use based on the conversation context.

### Tools Are Described, Not Hard-Coded

MCP servers describe their tools with names, descriptions, and parameter schemas. The LLM uses these descriptions to understand when and how to call each tool.

### Authentication Is Per-Connection

Each MCP connection has its own authentication. Internally, the orchestrator stores authentication as a generic set of HTTP header key/value pairs (`HeaderAuth map[string]string`). The "None / API Key (Bearer) / Custom Header / Basic Auth" labels you see in the UI are presentation conventions on top of that single map — they all resolve to headers that get attached to every MCP request. Header values are stored separately in the orchestrator's secret vault, not in the database.

## Next Steps

- [Configuration](./configuration) — Set up MCP connections in the orchestrator
- [Available MCPs](./available-mcps) — Pre-built MCP servers you can use
- [Troubleshooting](./troubleshooting) — Common issues and solutions
- **Building your own MCP server?** That's a developer task — see [Building a Custom MCP Server](../../developers/integrations/building-mcps.md) in the developer docs for the Python/Flask template, JSON-RPC wire format, authentication patterns, and deployment options.
