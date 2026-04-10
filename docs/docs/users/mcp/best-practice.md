---
sidebar_position: 6
---

# MCP Best Practices for Swarms

Practical strategies for using MCPs inside a swarm of agents.

## Give each expert its own MCPs

Different agent types use MCPs differently. A facilitator doesn't need CRM access — it delegates to experts. An expert needs the tools for its domain, not for everyone else's.

Example for an insurance claims team:

| Agent | Role | MCPs |
|---|---|---|
| elen (facilitator) | Coordinate, delegate, close | Kirha (web search for general context) |
| claudine (CRM expert) | Customer data | HubSpot, internal CRM MCP |
| victor (policy expert) | Policy documents | Markitdown, policy database MCP |
| ida (claims expert) | Claim processing | Claims system MCP |

The LLM isn't overwhelmed by irrelevant tools, and each expert has exactly what it needs.

## Multiple focused MCPs, not one mega-MCP

:::tip
Don't build one giant MCP with 50 tools. Build several focused MCPs with 5-10 tools each.
:::

A mega-MCP is hard to debug, hard to maintain, and confuses the LLM (which picks the wrong tool more often as the count grows). Splitting by domain gives you independent deployments, clearer error messages, and better tool selection.

**Tool count rule of thumb:** 10-20 well-described tools per expert is the sweet spot. Below 5, the agent can't fulfil its mission. Above 30, selection quality drops.

## Static vs dynamic data

For documents that rarely change (regulations, contracts, policies), you have three options:

- **File-based MCP** — bake the documents into the MCP server code. Fast, no external dependency, but requires redeployment to update. Good for updates less than monthly.
- **Database-backed MCP** — the MCP queries a database. Updates are live but requires infrastructure.
- **Markitdown + URLs** — store documents anywhere web-accessible and let the agent fetch them on demand via the Markitdown MCP.

## Prompting for MCP usage

Agents don't automatically know when to use a tool. Tell them in the System prompt:

```
You have access to:
• hubspot — customer data and history
• stripe — payment status and invoices
• websearch — current information

Rules:
- Use the right tool for each question
- Never say "I don't have access" without trying the relevant tool first
- Never make up data to replace what a tool would return
```

For facilitators specifically, list *which expert* has which MCP so the facilitator knows where to delegate:

```
Experts in this channel:
• @claudine.eth — CRM access (customer data)
• @victor.eth — policy access (contracts)
• @ida.eth — claims access (processing)
```

## MCPs during threading

Experts can call MCPs during the threading phase to validate the context primary key before locking it:

```
Threading prompt:
"Look for an invoice number in the message.
Before locking it, call the Stripe MCP to verify
the invoice exists. If not found, ask the user to confirm."
```

This prevents the agent from locking an invalid key and then producing garbage for the rest of the conversation.

## Combining MCP with hydration

The three memory strategies ([see Memory Strategies](../memory/memory-strategies.md)) layer cleanly with MCPs:

- **Hydration + MCP** — agent is pre-loaded with static knowledge (company policies, product catalogue) and uses MCPs for dynamic data (current customer records, real-time inventory).
- **Lazy + MCP** — agent learns from real conversations, but always checks live data via MCP for anything that can go stale.

## Error handling

MCPs fail. Tell the agent what to do when they do:

```
If an MCP tool returns an error or is unavailable:
1. Acknowledge the limitation to the user.
2. Offer alternative solutions if you have them.
3. Don't repeatedly call a failing tool.
4. Don't fabricate data to replace the missing response.
```

Without this rule, agents will hammer a failing endpoint on every turn, or worse, invent plausible-sounding data.

## Tool count reference

| Agent type | Max tools | Max MCPs |
|---|---|---|
| Simple expert | 10-15 | 2-3 |
| Complex expert | 20-30 | 4-5 |
| Facilitator | 5-10 | 1-2 |
