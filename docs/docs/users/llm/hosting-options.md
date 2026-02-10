---
sidebar_position: 2
---

# LLM Hosting Options

The 6022 Agentic Teams protocol offers four distinct approaches for LLM deployment, each with its own trade-offs between convenience, control, privacy, and cost optimization.

## Overview

| Option                 | Privacy    | Setup Complexity | Cost Model             | Microtransactions |
| ---------------------- | ---------- | ---------------- | ---------------------- | ----------------- |
| **Centralized LLM**    | ⭐⭐       | Easy             | Pay-per-use (API)      | ❌                |
| **Local (Ollama)**     | ⭐⭐⭐⭐   | Easy             | Hardware only          | ❌                |
| **Open Source Hosted** | ⭐⭐⭐     | Medium           | Infrastructure + usage | ❌                |
| **TEE (Enclave)**      | ⭐⭐⭐⭐⭐ | Advanced         | Agent-managed          | ✅                |

## Security & Trust Features

| Option                 | Data Confidentiality     | Model Verification           | Trust Model                     |
| ---------------------- | ------------------------ | ---------------------------- | ------------------------------- |
| **Centralized LLM**    | ❌ Data sent to provider | ❌ Proprietary models        | Trust the provider              |
| **Local (Ollama)**     | ✅ Data stays local      | ❌ No cryptographic proof    | Trust your hardware             |
| **Open Source Hosted** | ⚠️ Within your cloud     | ❌ No cryptographic proof    | Trust your cloud + admin        |
| **TEE (Enclave)**      | ✅ Encrypted in memory   | ✅ Cryptographic attestation | Trustless (verify, don't trust) |

### Understanding the Differences

**Data Confidentiality** — Where does your data go?

- **Centralized**: Your prompts and responses transit through third-party servers
- **Local**: Nothing leaves your machine — complete privacy
- **Cloud Hosted**: Data stays in your cloud environment, but admins could access it
- **TEE**: Data is encrypted even in memory — nobody can access it, not even the server admin

**Model Verification** — How do you know which model is running?

- **Centralized / Local / Cloud**: You trust that the provider runs the model they claim
- **TEE**: The CPU generates a cryptographic signature (attestation) proving exactly which code is running. This is the only option that provides **verifiable proof** that an unmodified open-source model is being used

## Quick Comparison

### Centralized LLM Providers

Use commercial APIs from OpenAI, Anthropic, Mistral, xAI, or Deepseek.

**Best for:** Quick start, POC, access to latest models

→ [Learn more](./centralized)

---

### Local Models (Ollama)

Run open-source models directly on your own hardware.

**Best for:** Complete privacy, zero API costs, offline operation

→ [Learn more](./local-ollama)

---

### Cloud Hosted (AWS, Azure, GCP)

Deploy open-source models on managed cloud infrastructure.

**Best for:** Enterprise scale, compliance, data residency

→ [Learn more](./cloud-hosted)

---

### TEE (Trusted Execution Environments)

Maximum security with cryptographic verification and agent microtransactions.

**Best for:** Sensitive data, decentralized AI vision, agent autonomy

→ [Learn more](./tee-enclave)

---

## Decision Matrix

| Your Priority                                   | Recommended Option              |
| ----------------------------------------------- | ------------------------------- |
| **Quick start**                                 | Centralized (OpenAI, Anthropic) |
| **Zero API costs**                              | Local (Ollama)                  |
| **Data confidentiality (simple)**               | Local (Ollama)                  |
| **Data confidentiality (scalable)**             | TEE                             |
| **Compliance requirements**                     | Open Source Hosted (your cloud) |
| **Model verification (prove which model runs)** | TEE                             |
| **Maximum security**                            | TEE                             |
| **Agent autonomy & microtransactions**          | TEE                             |

## Comparison by Use Case

| Use Case                                          | Best Option           | Why                                 |
| ------------------------------------------------- | --------------------- | ----------------------------------- |
| **POC / Testing**                                 | Centralized or Ollama | Fast setup                          |
| **Small team, privacy-conscious**                 | Ollama                | No costs, full confidentiality      |
| **Enterprise production**                         | Open Source Hosted    | Scalability, compliance             |
| **Sensitive data (finance, health)**              | TEE                   | Maximum security, confidentiality   |
| **Regulatory compliance (prove model integrity)** | TEE                   | Cryptographic model verification    |
| **Decentralized AI vision**                       | TEE                   | Agent autonomy, trustless operation |

## Migration Path

Many organizations follow this progression:

```
1. Start with Centralized LLM or Ollama
   (validate use case, quick POC)
         ↓
2. Move to Open Source Hosted
   (gain control, scale up)
         ↓
3. Adopt TEE for sensitive agents
   (maximum security, enable microtransactions)
```
