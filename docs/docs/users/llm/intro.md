---
sidebar_position: 1
---

# Introduction to LLM Configuration

Each agent in the 6022 Agentic Teams protocol can be configured with its own Large Language Model (LLM). This flexibility allows you to optimize costs, performance, and capabilities based on each agent's specific role.

## The Power of Model Diversity

One of the key strengths of the 6022 protocol is the ability to **mix different LLM providers and models** within the same agent team. This is not just a technical feature — it's a strategic advantage.

### Why Use Different Models?

**1. Richer Conversations**

When agents use different LLMs, each brings its own "personality" and reasoning patterns to the conversation. This diversity leads to more varied perspectives and reduces groupthink — similar to how a human team benefits from diverse backgrounds.

```
Example: Customer Support Team
├── Facilitator → GPT-4o (strong orchestration)
├── Legal Expert → Claude (nuanced reasoning)
├── Technical Expert → Deepseek (cost-effective for code)
└── CRM Expert → Mistral (fast responses)
```

**2. Better Response Quality**

Different models excel at different tasks:

| Model Strength    | Best Provider                  |
| ----------------- | ------------------------------ |
| Complex reasoning | Anthropic Claude, OpenAI GPT-4 |
| Fast responses    | Mistral, Groq                  |
| Code generation   | Deepseek, CodeLlama            |
| Cost efficiency   | Deepseek, Ollama (local)       |
| Tool/MCP usage    | OpenAI, xAI Grok               |

By matching each agent's role to the most suitable model, you optimize the overall team performance.

**3. Cost Optimization (Frugal Computing)**

Not every agent needs the most powerful (and expensive) model:

| Agent Role              | Model Choice                    | Why                        |
| ----------------------- | ------------------------------- | -------------------------- |
| Facilitator             | Premium (GPT-4o, Claude)        | Critical coordination role |
| Expert (complex domain) | Mid-tier (GPT-4, Mistral Large) | Needs good reasoning       |
| Expert (simple lookups) | Budget (Mistral, Deepseek)      | Mostly retrieves data      |
| Digital Twin            | Local (Ollama)                  | Privacy + cost savings     |

**4. Resilience & Redundancy**

Using multiple providers protects against:

- Single provider outages
- Rate limiting on one API
- Price changes from a single vendor

### Real-World Example

A customer service team with 5 agents:

| Agent    | Role              | LLM            | Monthly Cost |
| -------- | ----------------- | -------------- | ------------ |
| Elen     | Facilitator       | GPT-4o         | ~$50         |
| Claudine | CRM Expert        | Claude 3.5     | ~$30         |
| Victor   | Legal Expert      | Claude 3       | ~$25         |
| Charles  | Technical Support | Deepseek       | ~$5          |
| Albert   | FAQ Bot           | Ollama (local) | $0           |

**Total: ~$110/month** vs **~$250/month** if all agents used GPT-4o

:::tip Best Practice
Start with a premium model for your Facilitator (it orchestrates everything), then optimize expert agents based on their actual needs. Monitor conversations and adjust as you learn which agents need more capability.
:::

## Supported Providers

The orchestrator currently supports the following LLM providers:

| Provider        | Models                          | MCP Support | Notes                                |
| --------------- | ------------------------------- | ----------- | ------------------------------------ |
| **OpenAI**      | GPT-4, GPT-4o, GPT-5            | ✅ Full     | Recommended for quick start          |
| **Anthropic**   | Claude 3, Claude 3.5, Claude 4  | ✅ Full     | Requires custom header configuration |
| **xAI**         | Grok                            | ✅ Full     | Good alternative for MCP             |
| **Mistral**     | mistral-medium-latest           | ✅ Full     | OpenAI-compatible interface          |
| **Deepseek**    | deepseek-chat                   | ✅ Full     | Cost-effective option                |
| **AWS Bedrock** | Llama, Mistral, Claude, Titan   | ✅ Full     | Enterprise cloud hosting             |
| **Ollama**      | Llama 3, Mistral, Mixtral, etc. | ✅ Full     | Local deployment                     |

:::tip MCP Compatibility
MCP (Model Context Protocol) is now supported on **all LLM providers** in the 6022 protocol.
:::

## Configuring an LLM

### Accessing LLM Settings

1. Connect to the orchestrator with your wallet
2. Navigate to your agent list
3. Click the **⋮** menu on the desired agent
4. Select **"Manage LLM settings"**

### Configuration Parameters

| Parameter    | Description              | Example                         |
| ------------ | ------------------------ | ------------------------------- |
| **Provider** | The LLM service provider | `openai`, `anthropic`, `ollama` |
| **API URL**  | Endpoint URL             | `https://api.openai.com/v1`     |
| **Model**    | Specific model to use    | `gpt-4o`, `llama3`              |
| **API Key**  | Your provider's API key  | `sk-...`                        |

## Choosing the Right Provider

| Agent Type                         | Recommended Provider            |
| ---------------------------------- | ------------------------------- |
| Facilitator with complex reasoning | OpenAI GPT-4o, Anthropic Claude |
| Expert agent with MCP tools        | OpenAI, xAI Grok                |
| Simple conversational clone        | Mistral, Deepseek, Ollama       |
| Cost-sensitive deployment          | Deepseek, Ollama (local)        |
| Privacy-focused                    | Ollama (local), TEE             |

## Next Steps

Explore the different hosting options available:

- [Hosting Options Overview](./hosting-options) - Compare all 4 deployment approaches
- [Centralized Providers](./centralized) - OpenAI, Anthropic, Mistral, xAI, Deepseek
- [Local with Ollama](./local-ollama) - Run models on your own hardware
- [Cloud Hosted](./cloud-hosted) - AWS Bedrock and other cloud providers
- [TEE Enclaves](./tee-enclave) - Maximum security with microtransactions
