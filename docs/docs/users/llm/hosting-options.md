---
sidebar_position: 2
---

# LLM Hosting Options

Three shapes for where LLM inference runs. This page is the **decision guide** — pick an option, then jump to the per-option reference for configuration details.

| Option | Privacy | Setup | Cost | Microtransactions |
|---|---|---|---|---|
| **[Hosted Providers](./hosted-providers.md)** | ⭐⭐ | Easy | Pay-per-use API | ❌ |
| **[Local (Ollama)](./local-ollama.md)** | ⭐⭐⭐⭐ | Easy | Hardware only | ❌ |
| **[TEE (Enclave)](./tee-enclave.md)** | ⭐⭐⭐⭐⭐ | Advanced | Agent-managed | ✅ (future) |

**Hosted Providers** is the umbrella term for anything reachable over HTTP: first-party APIs (OpenAI, Anthropic, Mistral, xAI, Deepseek, LibertAI) and multi-model cloud platforms that aggregate many models (AWS Bedrock; Azure OpenAI and Vertex AI not yet supported).

## Data confidentiality and model verification

| Option | Your data goes to... | Can you prove which model ran? |
|---|---|---|
| **Hosted Providers** | The provider's infrastructure | No — you trust the provider |
| **Local (Ollama)** | Nowhere — it stays on your machine | No — you trust your hardware |
| **TEE (Enclave)** | Encrypted even in the server's RAM | ✅ Yes — cryptographic attestation by the CPU |

TEE is the only option that provides verifiable proof that an unmodified model is running your inference. It's also the only option that supports agent-managed microtransactions (each agent pays for its own LLM calls from its own wallet) — but it's on the roadmap, not shipped yet.

## Picking an option

| Your priority | Pick |
|---|---|
| Quick start, POC | Hosted (OpenAI or Anthropic) |
| Zero marginal cost | Local (Ollama) |
| Simple privacy | Local (Ollama) |
| Scalable privacy | TEE |
| Compliance / data residency | Hosted multi-model platform (Bedrock) |
| Cryptographic model verification | TEE |
| Agent autonomy & microtransactions | TEE |

**Typical progression:** start with a hosted provider to validate the use case, move to a multi-model platform once you need enterprise SLAs and data residency, and eventually adopt TEE for the sensitive agents where trustless verification matters.
