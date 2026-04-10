---
sidebar_position: 3
title: Hosted Providers
---

# Hosted LLM Providers

This page is the configuration reference for every LLM provider the orchestrator can talk to over HTTP. Whichever provider you pick, the dashboard asks for the same four fields: **Provider** (platform identifier), **API URL**, **Model**, and **API Key** — only the values change.

For an architectural comparison of *hosted* vs *self-hosted (Ollama)* vs *TEE*, see [LLM Hosting Options](./hosting-options.md). For per-provider MCP routing details, see [LLM Introduction → Supported Providers](./intro.md#supported-providers).

## Setup process

1. **Create an account** on the provider's platform (or, for LibertAI / Bedrock, on the cloud platform that hosts it).
2. **Generate an API key** in your provider dashboard.
3. **Configure the agent** in the orchestrator: open **⋮ → Manage LLM settings**, fill in Provider, API URL, Model, API Key, save.

The exact value to type into the **Provider** field is the *platform string* shown for each provider below — copy it verbatim. Friendly names like `openai` or `xai` are rejected by the orchestrator's enum validator.

---

## First-party providers

These are providers you call directly: you give your prompts to the model maker, you pay the model maker, and you trust the model maker to handle your data.

### OpenAI

```yaml
Provider: open-ai
API URL: https://api.openai.com/v1
Model: gpt-4o
Authentication: API Key
API Key: sk-...
```

The platform identifier is `open-ai` (with a hyphen) — `openai` is rejected.

**Get your API key:** [platform.openai.com/api-keys](https://platform.openai.com/api-keys)

---

### Anthropic

Anthropic requires a **custom header** configuration — do not use the standard "API Key" option, the orchestrator needs to send `x-api-key` rather than `Authorization: Bearer`.

```yaml
Provider: anthropic
API URL: https://api.anthropic.com/v1
Model: claude-3-5-sonnet-20241022
Authentication: Custom Header
Header Name: x-api-key
Header Value: sk-ant-...
```

:::warning
Selecting "API Key" instead of "Custom Header" for Anthropic will result in 401 errors at runtime. The header name **must** be `x-api-key`.
:::

**Get your API key:** [console.anthropic.com](https://console.anthropic.com/)

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

### xAI (Grok)

```yaml
Provider: x
API URL: https://api.x.ai/v1
Model: grok-1
Authentication: API Key
API Key: xai-...
```

The platform identifier is `x`, not `xai`.

**Get your API key:** [console.x.ai](https://console.x.ai/)

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

### LibertAI

[LibertAI](https://libertai.io) is a decentralized inference network. From the orchestrator's perspective it behaves like a regular OpenAI-Completions-compatible provider — same configuration shape as the other first-party APIs in this section.

```yaml
Provider: libert-ai
API URL: https://api.libertai.io
Model: <model-id-from-libertai>
Authentication: API Key
API Key: ...
```

The platform identifier is `libert-ai` (with the hyphen).

LibertAI is also the conceptual vehicle for the future TEE / enclave roadmap — see [TEE Enclaves](./tee-enclave.md). Today it operates purely as a hosted-API provider; the attestation/enclave layer doesn't exist yet.

---

## Multi-model cloud platforms

These are managed cloud services that aggregate multiple model families behind a single API. You configure the cloud platform once and gain access to many models from different makers.

### AWS Bedrock

Amazon Bedrock is a fully managed service that exposes foundation models from Anthropic, Meta, Mistral, Amazon, and others through one API. The orchestrator integrates with Bedrock via the Converse API.

```yaml
Provider: amazon-bedrock
Region: us-east-1
Model: anthropic.claude-3-sonnet-20240229-v1:0
Authentication: API Key
API Key: <single opaque key stored in vault>
```

The platform identifier is `amazon-bedrock` (with the hyphen) — `bedrock` is rejected.

The orchestrator builds the actual endpoint URL from the region using the template `https://bedrock-runtime.<region>.amazonaws.com`. You don't set the API URL directly; you set the region and the orchestrator constructs the URL.

**Authentication model:** Bedrock credentials are stored as a **single opaque API key** in the orchestrator's secret vault. The orchestrator does not maintain a separate `Access Key ID` / `Secret Access Key` pair, and there is no built-in 90-day rotation logic — rotate the key on whatever schedule your IAM policy requires and update it through the orchestrator UI.

**Available models (non-exhaustive):**

| Maker | Model ID | Use case |
|---|---|---|
| Anthropic | `anthropic.claude-3-sonnet-20240229-v1:0` | General purpose |
| Anthropic | `anthropic.claude-3-haiku-20240307-v1:0` | Fast, cost-effective |
| Meta | `meta.llama3-70b-instruct-v1:0` | Open source |
| Mistral | `mistral.mistral-large-2402-v1:0` | European AI |
| Amazon | `amazon.titan-text-express-v1` | Amazon native |

**Setup steps:**

1. Enable Bedrock in your AWS account
2. Request model access for the models you want to use
3. Create IAM credentials with Bedrock permissions
4. Configure in the orchestrator with the credentials

---

### Azure OpenAI — not yet supported

Azure OpenAI is **not currently a supported platform** in the orchestrator. There is no `azure` enum value, no Azure-specific settings, and no Azure code path in the chat-completion requester. If your team needs Azure support, file an issue — it would land as a new platform alongside `open-ai` and `anthropic`.

If you have an Azure deployment that exposes an OpenAI-compatible endpoint, you may be able to use the `open-ai` platform with a custom API URL, but this is not officially supported and has not been validated against Azure's authentication model.

---

### Google Cloud Vertex AI — not yet supported

Google Cloud Vertex AI is **not currently a supported platform** in the orchestrator. There is no `vertex` enum value, no Vertex settings file, no GCP credentials handling. If you need it, file an issue — it would require a new platform implementation.

---

## Pros and cons of hosted providers

**Advantages:**

- ✅ Quick setup (minutes)
- ✅ No infrastructure management
- ✅ Access to the latest models as soon as they ship
- ✅ Automatic scaling and high availability
- ✅ For multi-model platforms: data residency and enterprise SLAs

**Limitations:**

- ❌ **No data confidentiality** — your prompts and responses transit through the provider's infrastructure
- ❌ **No model verification** — you trust that the provider runs the model they claim
- ❌ Variable costs based on usage
- ❌ No microtransaction support (see [TEE Enclaves](./tee-enclave.md) for the future)

---

## Troubleshooting

For common HTTP-status-code issues (401, 429, 503), token-counting nuances (Anthropic), and per-provider quirks, see [Troubleshooting](./troubleshooting.md).
