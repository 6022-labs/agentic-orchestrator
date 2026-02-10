---
sidebar_position: 3
---

# Centralized LLM Providers

Use commercial LLM services via their official APIs. This is the simplest approach to get started with 6022 Agentic Teams.

## Supported Providers

| Provider | Platform | Models |
|----------|----------|--------|
| **OpenAI** | [platform.openai.com](https://platform.openai.com) | GPT-4, GPT-4o, GPT-5 |
| **Anthropic** | [console.anthropic.com](https://console.anthropic.com) | Claude 3, Claude 3.5, Claude 4 |
| **Mistral** | [console.mistral.ai](https://console.mistral.ai) | mistral-medium, mistral-large |
| **xAI** | [console.x.ai](https://console.x.ai) | Grok |
| **Deepseek** | [platform.deepseek.com](https://platform.deepseek.com) | deepseek-chat |

## Setup Process

1. **Create an account** on the provider's platform
2. **Generate an API key** in your dashboard
3. **Configure the agent** in the orchestrator with:
   - Provider selection
   - Model name
   - API key

---

## Provider Configuration

### OpenAI

```yaml
Provider: openai
API URL: https://api.openai.com/v1
Model: gpt-4o
Authentication: API Key
API Key: sk-...
```

**Get your API key:** [platform.openai.com/api-keys](https://platform.openai.com/api-keys)

---

### Anthropic

Anthropic requires a **custom header** configuration:

```yaml
Provider: anthropic
API URL: https://api.anthropic.com/v1
Model: claude-3-5-sonnet-20241022
Authentication: Custom Header
Header Name: x-api-key
Header Value: sk-ant-...
```

:::warning Important
Do not use the standard "API Key" authentication option for Anthropic. You must select **"Custom Header"** and specify `x-api-key` as the header name.
:::

**Get your API key:** [console.anthropic.com](https://console.anthropic.com/)

---

### xAI (Grok)

```yaml
Provider: xai
API URL: https://api.x.ai/v1
Model: grok-1
Authentication: API Key
API Key: xai-...
```

**Get your API key:** [console.x.ai](https://console.x.ai/)

---

### Mistral

```yaml
Provider: mistral
API URL: https://api.mistral.ai/v1
Model: mistral-medium-latest
Authentication: API Key
API Key: ...
```

**Get your API key:** [console.mistral.ai](https://console.mistral.ai/)

---

### Deepseek

```yaml
Provider: deepseek
API URL: https://api.deepseek.com/v1
Model: deepseek-chat
Authentication: API Key
API Key: sk-...
```

**Get your API key:** [platform.deepseek.com](https://platform.deepseek.com/)

---

## Pros & Cons

**Advantages:**
- ✅ Quick setup (minutes)
- ✅ No infrastructure management
- ✅ Access to latest models
- ✅ Automatic scaling
- ✅ High availability

**Limitations:**
- ❌ **No data confidentiality** (data sent to third-party servers)
- ❌ **No model verification** (proprietary models, trust the provider)
- ❌ Dependent on provider availability
- ❌ Variable costs based on usage
- ❌ No microtransaction support

---

## Troubleshooting

### 401 Unauthorized

**Problem:** Invalid or expired API key

**Solutions:**
- Verify the API key is correct
- Check if the key has expired
- Ensure you're using the correct header format (especially for Anthropic)

### 429 Rate Limit Exceeded

**Problem:** Too many requests

**Solutions:**
- Reduce request frequency
- Upgrade your provider plan
- Use separate API keys per agent if sharing causes limits

### Insufficient Quota

**Problem:** Account credits depleted

**Solutions:**
- Top up your provider account
- Enable auto-recharge in provider settings
- Monitor usage in provider dashboard

### Context Length Exceeded

**Problem:** Prompt too large for the model

**Solutions:**
- Reduce `MAX_TOKENS_PER_REQUEST` in orchestrator
- Simplify agent prompts
- Use a model with larger context window