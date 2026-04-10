---
sidebar_position: 5
---

# TEE (Trusted Execution Environments)

:::warning Draft — feature in development
TEE / enclave deployment is on the 6022 roadmap but **not yet shipped**. There is no `tee-enclave` platform value in the orchestrator, no enclave routing code, no attestation flow, and no in-orchestrator microtransaction logic. The configuration block in this page is forward-looking, not something you can paste into the dashboard today.

The conceptual vehicle for this roadmap is **[LibertAI](./hosted-providers.md#libertai)** — LibertAI is implemented today as a regular hosted-API provider, and the TEE / attestation layer will be built on top of it. If you want to start experimenting against LibertAI inference (without TEE guarantees), use the `libert-ai` provider in [Hosted Providers](./hosted-providers.md#libertai).
:::

Deploy LLM inference within secure enclaves (TEE - Trusted Execution Environments). This approach is intended to provide maximum privacy and security, with the unique ability for agents to manage their own compute costs through 6022 microtransactions.

## What is a TEE?

A Trusted Execution Environment is a secure area within a processor that guarantees:

- **Code integrity**: The model runs exactly as intended
- **Data confidentiality**: Prompts and responses are encrypted in memory
- **Attestation**: Cryptographic proof of what code is running

This means even the server administrator cannot access the data being processed.

---

## Architecture

```
┌─────────────────────────────────────────┐
│           Secure Enclave (TEE)          │
│  ┌─────────────────────────────────┐    │
│  │     Open Source LLM Model       │    │
│  │   (verified via CPU signature)  │    │
│  └─────────────────────────────────┘    │
│                  ↑                       │
│         Encrypted requests              │
└─────────────────────────────────────────┘
                   ↑
          6022 Orchestrator
                   ↑
            Agent (NFT + Wallet)
                   │
         Pays in 6022 tokens ←──── Microtransactions
```

---

## 6022 Microtransactions

With TEE deployment, agents can **pay for their own LLM usage** using their blockchain wallet. This is a unique feature of the 6022 protocol.

### How It Works

1. **Agent wallet**: Each agent has its own Ethereum wallet (part of the NFT)
2. **Token allocation**: Owner deposits 6022 tokens into the agent's wallet
3. **Automatic payment**: Agent pays per inference request
4. **Budget management**: Set spending limits per conversation

### Benefits

| Benefit | Description |
|---------|-------------|
| **True agent autonomy** | Agents manage their own resources |
| **Cost attribution** | Track expenses per agent precisely |
| **Incentive mechanisms** | Reward efficient agents with more tokens |
| **Decentralized operation** | No central billing dependency |
| **Transferable value** | Agent's balance transfers with ownership |

### Token Flow Example

```
1. Owner deposits 1000 tokens → Agent wallet
2. Conversation starts → Agent reserves 50 tokens
3. Each LLM call → Agent pays ~0.5 tokens
4. Conversation ends → Unused tokens returned to reserve
5. Facilitator distributes rewards → Agents receive tokens
```

---

## Configuration (Future)

The example below is **not yet a real platform**. It illustrates the shape of the future configuration; do not paste it into the dashboard today.

```yaml
Provider: <future tee enum value, not yet defined>
Enclave URL: https://<your-tee-enclave>
Model: llama-3-70b
Authentication: Agent Wallet Signature
Payment: 6022 tokens (per request)
Budget: 500 tokens/conversation
```

---

## TEE Providers

The 6022 protocol is designed to work with decentralized TEE providers. Two are on the radar:

| Provider | Technology | Description |
|---|---|---|
| **LibertAI** | TEE Inference | Decentralized AI inference network with privacy-preserving computation. **This is the partner the 6022 roadmap is built around.** Available today as a regular hosted-API provider — see [Hosted Providers → LibertAI](./hosted-providers.md#libertai). |
| **Phala** | Phat Contracts | Web3 cloud computing platform with TEE-based confidential smart contracts. Aspirational — not currently integrated. |

Both providers are intended to enable the **trustless verification** of model execution that is central to the 6022 vision of autonomous agents. Right now, only the LibertAI side has a code path in the orchestrator (via the `libert-ai` platform), and that code path does not yet involve attestation or enclave verification.

---

## Security Guarantees

### What TEE Protects Against

| Threat | Protected? |
|--------|------------|
| Malicious server admin | ✅ Yes |
| Memory scraping attacks | ✅ Yes |
| Man-in-the-middle | ✅ Yes |
| Model tampering | ✅ Yes (via attestation) |
| Side-channel attacks | ⚠️ Partial (depends on implementation) |

### Attestation Process

1. **Enclave starts** with the LLM model
2. **CPU generates signature** of the enclave contents
3. **Signature is verified** by the orchestrator
4. **Only verified enclaves** receive requests

This ensures you're always talking to the correct, unmodified model.

---

## Roadmap Status

:::info Coming Soon
TEE integration with microtransactions is part of the 6022 protocol roadmap.
:::

| Milestone | Status |
|---|---|
| LibertAI as a regular hosted-API provider | ✅ Implemented (`libert-ai` platform) |
| TEE proof of concept | 🔜 In design |
| Enclave attestation | 🔜 In design |
| Orchestrator TEE platform (`tee-*` enum value) | 🔜 Not started |
| Microtransaction smart contracts | 🔜 Not started |
| Production deployment | 🔜 Not started |

---

## Pros & Cons

**Advantages:**
- ✅ **Complete data confidentiality** (encrypted even in memory)
- ✅ **Cryptographic model verification** (prove which model runs via attestation)
- ✅ Agent financial autonomy
- ✅ True decentralization
- ✅ Microtransaction support (agent pays for itself)
- ✅ Compliance with strictest security requirements
- ✅ Trustless operation (verify, don't trust)

**Limitations:**
- ❌ More complex setup
- ❌ Limited model selection (must fit in enclave memory)
- ❌ Higher initial infrastructure cost
- ❌ Feature still in development
- ❌ Slightly higher latency than direct inference

---

## Use Cases

### Ideal For

| Use Case | Why TEE? |
|----------|----------|
| **Financial services** | Regulatory compliance, data protection |
| **Healthcare** | HIPAA compliance, patient data privacy |
| **Legal** | Client confidentiality |
| **Government** | Classified information handling |
| **Decentralized AI** | Agent autonomy, trustless operation |

### The Path to Autonomous Agents

TEE is the foundation for the 6022 vision of truly autonomous agents:

```
Current State:
  Agent → Orchestrator → Cloud LLM
  (Owner pays for everything)

Future with TEE:
  Agent (with wallet) → TEE Enclave
  (Agent pays for its own compute)
  
Ultimate Vision:
  Liberated Agent → Own TEE → Own resources
  (Fully autonomous, self-sustaining)
```

This is what 6022 calls the **"Eternal Life"** concept: agents that can exist independently, managing their own resources and potentially outliving their original owners.

---

## Getting Started (Preview)

While full TEE integration is in development, you can:

1. **Test enclave concepts** using AWS Nitro Enclaves
2. **Prepare your agents** with wallet configurations
3. **Monitor the roadmap** for production release
4. **Join the community** to provide feedback

For early access or partnership inquiries, contact the 6022 team.