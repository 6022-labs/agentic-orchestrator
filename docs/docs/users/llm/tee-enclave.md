---
sidebar_position: 5
---

# TEE (Trusted Execution Environments)

:::warning Draft — feature not yet shipped
TEE deployment is on the 6022 roadmap but **not yet shipped**. No `tee-enclave` platform value exists in the orchestrator, no enclave routing code, no attestation flow, and no in-orchestrator microtransaction logic.

The vehicle for this roadmap is **[LibertAI](./hosted-providers.md#libertai)**, which exists today as a regular hosted-API provider. The TEE / attestation layer will be built on top of LibertAI later.
:::

## Why TEE matters

A Trusted Execution Environment is a secure area in a CPU that guarantees:

- **Code integrity** — the model runs exactly as intended
- **Data confidentiality** — prompts and responses are encrypted even in memory
- **Attestation** — cryptographic proof of which code is actually running

Even the server administrator cannot see the data being processed. Together with agent wallets, TEE enables a design where each agent pays for its own LLM calls from its own balance — no central billing, no trust in the provider, and verifiable proof that the model you paid for is the one that actually ran your inference.

## What TEE protects against

| Threat | Protected? |
|---|---|
| Malicious server admin | ✅ |
| Memory scraping | ✅ |
| Man-in-the-middle | ✅ |
| Model tampering | ✅ (via attestation) |
| Side-channel attacks | ⚠️ Partial — depends on the enclave implementation |

## How agent microtransactions would work

This is the future design, not current reality:

1. The owner deposits 6022 tokens into the agent's wallet (each agent is an NFT with its own wallet, so this is already possible today).
2. When a conversation starts, the agent reserves a budget (e.g. 50 tokens).
3. Each LLM call deducts ~0.5 tokens from the reserve.
4. Unused tokens return to the agent's wallet when the conversation ends.
5. The facilitator can optionally distribute rewards to agents that contributed useful answers.

The result: agents manage their own compute budgets, owners can track costs precisely per agent, and efficient agents survive longer than wasteful ones. The longer-term vision 6022 calls "eternal life" — agents that exist independently, manage their own resources, and potentially outlive their original owners.

## TEE providers we're tracking

- **LibertAI** — decentralized AI inference network, privacy-preserving computation. The partner the 6022 roadmap is built around. Already usable today as a regular hosted-API provider via the `libert-ai` platform — see [Hosted Providers → LibertAI](./hosted-providers.md#libertai). The attestation layer is what's still to come.
- **Phala** — Web3 cloud computing with TEE-based confidential smart contracts. Aspirational, not currently integrated.

## Roadmap status

| Milestone | Status |
|---|---|
| LibertAI as a regular hosted-API provider | ✅ Implemented (`libert-ai` platform) |
| TEE proof of concept | 🔜 In design |
| Enclave attestation | 🔜 In design |
| Orchestrator TEE platform (new `tee-*` enum) | 🔜 Not started |
| Microtransaction smart contracts | 🔜 Not started |
| Production deployment | 🔜 Not started |

## What you can do today

- Use the `libert-ai` platform in [Hosted Providers](./hosted-providers.md#libertai) to start building against LibertAI inference (without TEE guarantees).
- Watch the roadmap — we'll publish a real configuration section here once the first milestones ship.
- Contact the 6022 team for early access or partnership inquiries.
