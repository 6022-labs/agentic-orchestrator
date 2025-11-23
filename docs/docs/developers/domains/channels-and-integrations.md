---
id: channels-and-integrations
title: Domain — Channels & Integrations
sidebar_label: Channels & Integrations
---

This domain gives agents a way to talk to humans where they already collaborate (Slack today, more channels later). Think of it as the translation layer between external workspaces and the orchestrator: it receives events from a channel, keeps enough local state to understand the conversation, and sends the agent’s reply back out.

What it owns:

- Channel and membership records so we know which agents are allowed to act in each space.
- Conversation history (threads, messages, conclusion status) so prompts sent to the orchestrator include the right context.
- User profiles and optional wallet links when a workspace identity needs to be tied to an on-chain persona.
- Cached copies of agent metadata pulled from the orchestrator so channel services can label responses and apply role rules (facilitator, human, etc.).

How it works (high level):

- The Slack bridge listens to Events API callbacks and slash commands, validates the signatures, and persists any new information before touching the orchestrator.
- When a facilitator agent should answer, the bridge grabs the relevant thread history, calls the orchestrator for a completion, and posts the response back to Slack.
- Commands also cover operational chores—adding/removing agents from a channel, verifying wallets, or pairing a Slack user with their human agent token.

Future channels & integrations:

As we add new channels (Discord, Teams, email, etc.), each connector should keep the same responsibilities: capture events, normalize conversation state, call the orchestrator, and deliver the response in the native channel format.
