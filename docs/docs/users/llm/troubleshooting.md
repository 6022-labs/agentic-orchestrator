---
sidebar_position: 6
---

# LLM Troubleshooting

Common LLM issues and their fixes.

## Token limit: `TOKEN__MAX_TOKENS_PER_REQUEST`

The orchestrator caps how many tokens it will send in a single request via the env var **`TOKEN__MAX_TOKENS_PER_REQUEST`**. The `TOKEN__` prefix is required — setting the bare `MAX_TOKENS_PER_REQUEST` has no effect (koanf only reads env vars with the `__` separator).

| Scenario | Value |
|---|---|
| Default | `100000` |
| Anthropic | `50000` (they count tokens more aggressively) |
| Small models / low memory | `25000` |

The cap is **global**, not per-provider, so if any agent uses Anthropic you may need to lower the cap for everyone.

## HTTP status codes

| Code | Meaning | Fix |
|---|---|---|
| `401` | Invalid API key | Regenerate the key, double-check the header (Anthropic needs `x-api-key`, not `Authorization: Bearer`) |
| `403` | Permissions or plan issue | Check the key has the required scopes, upgrade the plan if needed |
| `429` | Rate limit hit | Wait and retry, upgrade the plan, use separate keys per agent, or distribute across providers |
| `500` | Provider-side error | Retry with exponential backoff |
| `503` | Provider is down | Wait it out; consider switching to a fallback provider if this happens often |

## Anthropic

**"Context too long" errors** — Anthropic counts tokens differently from OpenAI. For the same logical content, Anthropic typically reports 10-20% more tokens. Lower `TOKEN__MAX_TOKENS_PER_REQUEST` to `50000` (from the default `100000`).

**Authentication failing** — use **Custom Header** auth with header name `x-api-key`, not "API Key" (Bearer). See [Hosted Providers → Anthropic](./hosted-providers.md#anthropic).

## AWS Bedrock

**Authentication fails** — rotate the opaque API key stored in the orchestrator's vault on whatever schedule your IAM policy requires. The orchestrator doesn't model a fixed expiration window. Verify IAM permissions include Bedrock, and confirm the model is enabled in your region.

**Model not available** — request access in the AWS Bedrock console for the specific model ID.

## Ollama

**Connection refused** — `ollama serve` isn't running, or the orchestrator is on a different machine and you haven't bound to all interfaces. For remote access: `OLLAMA_HOST=0.0.0.0 ollama serve`. Also check that the firewall allows port 11434.

**Out of memory / slow responses** — see [Local with Ollama → Troubleshooting](./local-ollama.md#troubleshooting).

## Agent not responding at all

Run through this checklist:

1. Is the API key valid and does the account still have credits?
2. Is the model name spelled exactly right? (`claude-3-5-sonnet-20241022` not `claude-3.5-sonnet`)
3. Is the API URL correct for the provider?
4. Are there 429 or 401 errors in the orchestrator logs?
5. Is the context window within the model's limit?

## Common log patterns (orchestrator logs)

```
"error": "insufficient_quota"        → Top up the provider account
"error": "rate_limit_exceeded"       → Too many concurrent requests; use per-agent keys
"error": "context_length_exceeded"   → Lower TOKEN__MAX_TOKENS_PER_REQUEST or shorten prompts
"error": "invalid_api_key"           → Regenerate the key, check for trailing whitespace
```

## Where to ask for help

Include the agent's ENS domain, the LLM provider and model, the full error message from the logs, a rough timestamp, and the steps to reproduce. Post to [t.me/protocol6022](https://t.me/protocol6022) or open a GitHub issue at [github.com/6022-labs](https://github.com/6022-labs).
