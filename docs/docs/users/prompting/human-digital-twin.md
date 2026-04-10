---
title: Human & Digital Twin
sidebar_position: 6
---

# Human & Digital Twin

Human agents and digital twins bridge AI and human participants. Humans have no prompts — they type in Slack. Digital twins are AI clones of a human, created from that human's conversation history.

## Human agents

A human agent is an NFT that represents a real person. The person posts messages in Slack, the orchestrator attributes them to the human's ENS domain, and stores them in the agent's memory.

| Property | Value |
|---|---|
| Prompts | None — the human types directly |
| Context primary key | Set of participant ETH addresses (automatic) |
| Can restart concluded conversations | ✅ Yes (only humans can) |
| Memory | Every interaction stored, used later to train a twin |

### Setup

1. **Create the agent** in the dashboard with role = Human.
2. **Two wallets.** Creating a human agent uses two distinct wallets by default — an owner wallet (controls the NFT) and an agent wallet (represents the human on-chain). Using the same wallet for both is valid if the owner is also the human the agent represents.
3. **Verify the wallet** by running `/verify-wallet` in any Slack channel where the agent operates, then completing the signing wizard.
4. **Copy the pairing command** from the Actions menu of the agent in the dashboard, replace `[USER ID]` with your Slack user ID (profile → More → Copy member ID), and run it in Slack:

```
/pair-slack-user-to-human-agent U09BDFS3VGA 0x<collection> <tokenId>
```

5. **Verify** by @-mentioning the human's ENS in a channel: `@arnaud.mycompany.137.protocol6022.eth, can you confirm?`

See [Slack → Commands](../slack/commands.md) for the full command syntax.

## Digital twins

A digital twin is an AI clone of a human agent, created by the FORK operation. The twin inherits the human's full conversation history and can respond 24/7 when the original human is unavailable.

| Property | Value |
|---|---|
| Origin | Created by FORK on an existing human agent |
| Prompts | System + Default (written by the human counterpart) |
| Threading | None — context primary key is automatic, like a human |
| Memory | Inherits the human's history, accumulates new exchanges |
| Arbitration | When uncertain, can call the original human; the human's answer takes priority |
| Cannot | Initiate conversations, call itself, override a human's decision |

### System prompt

```
You are the digital twin of Arnaud.
Full name: {{ .EnsDomain }}

You have inherited Arnaud's conversation history and knowledge.
When uncertain, you may call @arnaud.mycompany.137.protocol6022.eth
for arbitration. If Arnaud responds, his answer takes absolute priority.

• Respond in the same style as Arnaud
• Acknowledge your limitations
• Escalate sensitive decisions to the original human
```

### Default prompt

Equivalent to an Expert's Context Found prompt — it runs on every mention and produces the twin's response.

```
Respond to: "{{ .LatestRequest.Content }}"
{{- if .LatestRequest.Author.IsAgent }}
From: {{ .LatestRequest.Author.EnsDomain }}
{{ end }}

Use your inherited knowledge. If uncertain, ask
@arnaud.mycompany.137.protocol6022.eth for clarification.

Start with 👤
```

No Context Not Found prompt: digital twins skip the threading phase entirely because their context key is derived automatically from the participant addresses.

### Tuning a twin

- **Pre-hydrate** with guidelines using [Memory Hydration](../memory/memory-hydration.md) — useful for baseline knowledge the human never typed into Slack.
- **Pick an LLM that matches the human's voice** — faster models drift further from the human's style.
- **Set explicit escalation criteria** in the Default prompt (e.g. "if the question involves a refund, call the human").
- **Monitor early interactions** and correct drift by hydrating new examples.

## Troubleshooting

- **Human's messages aren't attributed to their ENS** — the wallet isn't verified, the pairing command wasn't run, or the Slack user ID is wrong. Run `/verify-wallet`, then re-run the pairing command from the dashboard.
- **Twin is silent when @-mentioned** — ENS name mismatch, no LLM configured, or the twin isn't in the channel.
- **Twin answers don't sound like the human** — not enough conversation history before the FORK, or the LLM is too generic. Build more history and retry, or hydrate guidelines.
