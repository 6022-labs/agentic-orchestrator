---
title: Threading Phase
sidebar_position: 5
---

# Threading Phase

When an expert is first @-mentioned in a conversation, it doesn't yet know which [context primary key](./context-primary-keys.md) applies. The **threading phase** is the short routine it runs to figure that out before loading memory and responding.

## The five steps

1. **Mention detected.** The facilitator (or another agent) writes `@claudine.mycompany.137.protocol6022.eth, can you help with customer C45653?`. The orchestrator routes the message to Claudine.
2. **Context Not Found prompt runs.** Claudine doesn't have a key yet, so the orchestrator runs her [Context Not Found prompt](../prompting/expert-prompts.md#context-not-found-prompt) with the message as input.
3. **Key extracted.** The prompt asks the model to find a valid key in the message. For Claudine that might be a regex like `^C[0-9]{5}$`, which matches `C45653`.
4. **Memory loaded.** The orchestrator loads all history stored under key `C45653` for this agent, regardless of which conversation it came from.
5. **Conversation is bound to the key.** A row is persisted in `agent_context_conversation_mappings` linking this conversation to `C45653` for the rest of its lifetime. Subsequent messages in the same thread bypass threading and go straight to the [Context Found prompt](../prompting/expert-prompts.md#context-found-prompt).

There is no distributed lock: two concurrent threading attempts on the same conversation ID are not guarded. In practice this isn't an issue because the orchestrator processes messages sequentially per conversation.

## Three ways the key gets identified

- **Explicit mention** — the key is literally in the message: `"Can you help with customer C45653?"` → `C45653`.
- **Pattern matching** — the message contains something that matches a regex the prompt taught the model: `"Customer (customer@example.com) needs help"` → `customer@example.com`.
- **MCP lookup** — the agent calls an MCP server during threading to resolve a fuzzy identifier: `"The customer from the automotive workshop"` → the MCP returns `CUST-45653`. See [MCPs during threading](../mcp/best-practice.md#mcps-during-threading).

If none of these work, the agent asks for help from another agent by name (per the "trust the swarm" pattern). The conversation stays open and the next turn gets another shot at threading.

## What happens if threading fails repeatedly

The conversation continues without a locked key, which means the expert has no memory to load. Either tune the Context Not Found prompt (add more explicit patterns, loosen the regex, add an MCP lookup), or accept that some conversations genuinely don't have a matching key and the expert should decline to answer.

Humans and digital twins **skip threading entirely** — their key is computed automatically from the set of participant wallet addresses, so there's nothing to figure out.

## Configuration

All the configuration happens in the Expert's [Context Not Found prompt](../prompting/expert-prompts.md#context-not-found-prompt). The orchestrator-side machinery (running the prompt, extracting the key, persisting the mapping) requires no setup.
