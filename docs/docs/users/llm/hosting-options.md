---
sidebar_position: 2
---

# LLM Hosting Options

The 6022 Agentic Teams protocol supports three deployment shapes for LLMs, each with its own trade-offs between convenience, control, privacy, and cost. This page is the **decision guide** — once you've picked an option, the configuration steps live on the corresponding reference page.

## Overview

| Option                         | Privacy    | Setup Complexity | Cost Model                | Microtransactions |
| ------------------------------ | ---------- | ---------------- | ------------------------- | ----------------- |
| **Hosted Providers**           | ⭐⭐       | Easy             | Pay-per-use (API)         | ❌                |
| **Local (Ollama)**             | ⭐⭐⭐⭐   | Easy             | Hardware only             | ❌                |
| **TEE (Enclave)**              | ⭐⭐⭐⭐⭐ | Advanced         | Agent-managed             | ✅ (future)       |

**Hosted Providers** is the umbrella term for everything you reach over HTTP — first-party APIs from model makers (OpenAI, Anthropic, Mistral, xAI, Deepseek, LibertAI) **and** multi-model cloud platforms that aggregate models from several makers (AWS Bedrock; Azure OpenAI and Vertex AI are not yet supported). The configuration shape is the same for all of them; see [Hosted Providers](./hosted-providers.md) for the per-provider reference.

## Security & Trust Features

| Option               | Data Confidentiality     | Model Verification           | Trust Model                     |
| -------------------- | ------------------------ | ---------------------------- | ------------------------------- |
| **Hosted Providers** | ❌ Data sent to provider | ❌ Proprietary models        | Trust the provider              |
| **Local (Ollama)**   | ✅ Data stays local      | ❌ No cryptographic proof    | Trust your hardware             |
| **TEE (Enclave)**    | ✅ Encrypted in memory   | ✅ Cryptographic attestation | Trustless (verify, don't trust) |

### Understanding the differences

**Data Confidentiality** — Where does your data go?

- **Hosted Providers**: Your prompts and responses transit through third-party servers (the model maker's, or the cloud platform's).
- **Local (Ollama)**: Nothing leaves your machine — complete privacy.
- **TEE**: Data is encrypted even in memory — nobody can access it, not even the server admin.

**Model Verification** — How do you know which model is running?

- **Hosted / Local**: You trust that the provider runs the model they claim.
- **TEE**: The CPU generates a cryptographic signature (attestation) proving exactly which code is running. This is the only option that provides **verifiable proof** that an unmodified model is being used.

## Quick Comparison

### Hosted Providers

Use commercial APIs from OpenAI, Anthropic, Mistral, xAI, Deepseek, LibertAI, or aggregators like AWS Bedrock.

**Best for:** Quick start, POC, access to latest models, enterprise compliance through multi-model platforms.

→ [Learn more](./hosted-providers.md)

---

### Local Models (Ollama)

Run open-source models directly on your own hardware.

**Best for:** Complete privacy, zero API costs, offline operation.

→ [Learn more](./local-ollama.md)

---

### TEE (Trusted Execution Environments)

Maximum security with cryptographic verification and agent microtransactions.

**Best for:** Sensitive data, decentralized AI vision, agent autonomy.

→ [Learn more](./tee-enclave.md)

---

## Decision Matrix

| Your Priority                                   | Recommended Option                  |
| ----------------------------------------------- | ----------------------------------- |
| **Quick start**                                 | Hosted (OpenAI, Anthropic)          |
| **Zero API costs**                              | Local (Ollama)                      |
| **Data confidentiality (simple)**               | Local (Ollama)                      |
| **Data confidentiality (scalable)**             | TEE                                 |
| **Compliance / data residency requirements**    | Hosted via multi-model platform (AWS Bedrock) |
| **Model verification (prove which model runs)** | TEE                                 |
| **Maximum security**                            | TEE                                 |
| **Agent autonomy & microtransactions**          | TEE                                 |

## Comparison by Use Case

| Use Case                                          | Best Option                                 | Why                                 |
| ------------------------------------------------- | ------------------------------------------- | ----------------------------------- |
| **POC / Testing**                                 | Hosted or Ollama                            | Fast setup                          |
| **Small team, privacy-conscious**                 | Ollama                                      | No costs, full confidentiality      |
| **Enterprise production**                         | Hosted (multi-model platform like Bedrock)  | Scalability, compliance, SLAs       |
| **Sensitive data (finance, health)**              | TEE                                         | Maximum security, confidentiality   |
| **Regulatory compliance (prove model integrity)** | TEE                                         | Cryptographic model verification    |
| **Decentralized AI vision**                       | TEE                                         | Agent autonomy, trustless operation |

## Migration Path

Many organisations follow this progression:

```
1. Start with a hosted provider or Ollama
   (validate the use case, quick POC)
            ↓
2. Move to a multi-model cloud platform (e.g. AWS Bedrock)
   (gain enterprise SLAs, compliance, data residency)
            ↓
3. Adopt TEE for sensitive agents
   (maximum security, enable microtransactions)
```
