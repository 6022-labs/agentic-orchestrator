---
sidebar_position: 6
---

# MCP Best Practices for Agentic Teams

This guide covers practical strategies for using MCP effectively within the 6022 Agentic Teams swarm AI architecture.

## MCP in the Swarm Context

In a swarm of agents, MCPs provide shared access to external data. The key insight is that **different agent types use MCPs differently**:

```
┌─────────────────────────────────────────────────────────────┐
│                      SLACK CHANNEL                          │
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ FACILITATOR │  │   EXPERT    │  │   EXPERT    │         │
│  │   (Elen)    │  │ (Claudine)  │  │  (Victor)   │         │
│  │             │  │    CRM      │  │   Legal     │         │
│  │ Sees all    │  │             │  │             │         │
│  │ messages    │  │ Called via  │  │ Called via  │         │
│  │             │  │ @mention    │  │ @mention    │         │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘         │
│         │                │                │                 │
│         ▼                ▼                ▼                 │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   MCP SERVERS                        │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌────────┐  │   │
│  │  │ HubSpot │  │ Stripe  │  │  Kirha  │  │Markitdown│  │   │
│  │  │   CRM   │  │Payments │  │ Search  │  │  Docs   │  │   │
│  │  └─────────┘  └─────────┘  └─────────┘  └────────┘  │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Strategy 1: Specialized MCPs per Expert

**Principle:** Each expert agent gets MCPs relevant to its domain.

### Example: Insurance Claims Team

| Agent | Role | MCPs Assigned |
|-------|------|---------------|
| **elen.agent** | Facilitator | Kirha (web search) |
| **claudine.agent** | CRM Expert | HubSpot, Custom CRM MCP |
| **victor.agent** | Policy Expert | Markitdown, Policy MCP |
| **ida.agent** | Claims Expert | Claims System MCP |

**Why this works:**
- Each expert has tools matching its specialty
- The LLM isn't overwhelmed by irrelevant tools
- Clear separation of responsibilities

### Configuration

For the CRM Expert (claudine):
```
MCPs:
├── hubspot-mcp (customer data)
└── crm-custom-mcp (internal CRM)
```

For the Policy Expert (victor):
```
MCPs:
├── markitdown (read policy PDFs)
└── policy-mcp (search policy database)
```

---

## Strategy 2: Multiple Specialized MCPs vs One Mega-MCP

:::tip Key Learning
Don't create one giant MCP that does everything. Use multiple specialized MCPs instead.
:::

### ❌ Anti-pattern: Mega-MCP

```
Single MCP with 50+ tools:
- search_customers
- get_invoice
- read_policy
- check_gdpr
- process_claim
- send_email
- ...
```

**Problems:**
- LLM gets confused with too many tools
- Hard to maintain
- Single point of failure

### ✅ Best Practice: Specialized MCPs

```
Agent: claims-expert.mycompany.137.protocol6022.eth
├── MCP: crm          → 5 tools (customer data)
├── MCP: claims       → 4 tools (claim processing)
├── MCP: documents    → 3 tools (policy documents)
└── MCP: websearch    → 2 tools (external info)
```

**Benefits:**
- Each MCP has focused, well-described tools
- Easier to debug when issues occur
- Can update MCPs independently

---

## Strategy 3: MCP for Static Documents

For documents that rarely change (regulations, contracts, procedures), you have two options:

### Option A: File in MCP Repository

If documents update less than monthly, include them directly in the MCP server code:

```
your-mcp-server/
├── app.py
├── data/
│   ├── gdpr_articles.json
│   ├── company_policies.json
│   └── product_catalog.json
└── tools/
    └── search_policies.py
```

**Pros:** Fast, no external dependency
**Cons:** Requires redeployment to update

### Option B: Database-Backed MCP

If documents update more frequently:

```python
# MCP connects to database
def search_policies(query):
    return db.query("SELECT * FROM policies WHERE content LIKE ?", f"%{query}%")
```

**Pros:** Update without redeployment
**Cons:** More infrastructure

### Option C: Markitdown + URLs

For documents accessible via URL:

```
Agent prompt: "When asked about policy details, use the markitdown 
MCP to fetch and read the document from the provided URL."
```

The agent calls Markitdown with the document URL, which returns Markdown text.

---

## Strategy 4: Prompting for MCP Usage

The LLM needs guidance on when and how to use MCP tools. Include instructions in your agent's system prompt.

### Basic MCP Awareness

```
You have access to the following MCP tools:
- hubspot: Search and retrieve customer information
- stripe: Check payment status and invoices
- websearch: Find current information on the internet

When a question requires external data, use the appropriate tool 
before answering. Do not guess or make up information.
```

### Encouraging Tool Usage

```
IMPORTANT: Always verify information using available tools.
- For customer questions → use hubspot tools
- For payment questions → use stripe tools
- For current events → use websearch tools

Never say "I don't have access to that information" without first 
trying the relevant MCP tool.
```

### Tool Selection Guidance

```
When deciding which tool to use:
1. If the question mentions a customer name or email → hubspot
2. If the question mentions invoice, payment, or subscription → stripe
3. If the question requires recent/current information → websearch
4. If the question references a document URL → markitdown
```

---

## Strategy 5: MCP During Threading (Experts)

Expert agents can use MCPs during the **threading phase** to validate context primary keys.

### Example: Invoice Expert

```
Threading prompt:
"Analyze the conversation to find the invoice number.
Use the stripe MCP to verify the invoice exists before 
setting it as the context primary key.

If the invoice is not found, ask the user to confirm the number."
```

### Example: Customer Expert

```
Threading prompt:
"Look for customer identifiers (email, phone, customer ID).
Use the hubspot MCP to validate the customer exists.
Set the validated customer ID as the context primary key."
```

This ensures the agent doesn't thread on invalid keys.

---

## Strategy 6: Combining MCP with Memory Strategies

MCP works alongside the three memory strategies (Lazy, Hydration, MCP). You can combine them:

### Hydration + MCP

```
┌─────────────────────────────────────────────────────────────┐
│                      AGENT MEMORY                           │
│                                                             │
│  HYDRATED (static):              MCP (dynamic):             │
│  ├── Company policies            ├── Customer data          │
│  ├── Product catalog             ├── Current prices         │
│  └── FAQ responses               └── Real-time inventory    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Use case:** Agent hydrated with general knowledge, uses MCP for customer-specific data.

### Lazy + MCP

```
Conversation naturally builds memory:
User: "My order #12345 hasn't arrived"
Agent: [calls MCP to get order status]
Agent: [stores order info in context memory]
...later...
Agent: [uses stored info + calls MCP for updates]
```

**Use case:** Agent learns from conversations but always checks live data via MCP.

---

## Strategy 7: Error Handling in Prompts

Agents should handle MCP failures gracefully:

```
If an MCP tool returns an error or is unavailable:
1. Acknowledge the limitation to the user
2. Offer alternative solutions if available
3. Do NOT repeatedly call a failing tool
4. Do NOT make up data to replace what the tool would have returned

Example response:
"I'm currently unable to access the payment system. 
I can help you with general questions, or you can try 
again in a few minutes for payment-specific queries."
```

---

## Strategy 8: MCP Tool Limits

### Recommended Limits

| Scenario | Max Tools per Agent | Max MCPs per Agent |
|----------|:-------------------:|:------------------:|
| Simple agent | 10-15 | 2-3 |
| Complex expert | 20-30 | 4-5 |
| Facilitator | 5-10 | 1-2 |

### Why Limits Matter

- **Too many tools:** LLM struggles to choose the right one
- **Too few tools:** Agent can't fulfill its mission
- **Sweet spot:** 10-20 well-described tools

---

## Strategy 9: Agents Discussing Tools

An interesting emergent behavior: agents can discuss tool availability with each other!

```
Facilitator: "@claudine.agent Can you check the customer history?"
Claudine: "I have access to HubSpot. Let me search for that customer..."
Claudine: [calls hubspot MCP]
Claudine: "Found them. Here's their history..."
```

This happens naturally when:
- Agents have different MCP configurations
- The facilitator knows which expert has which tools
- Prompts encourage collaboration

### Facilitator Prompt for Tool Delegation

```
You are the facilitator. You know which experts have which capabilities:
- @claudine.agent has CRM access (customer data, history)
- @victor.agent has policy access (contracts, regulations)
- @ida.agent has claims access (claim status, processing)

When a question requires specific data, delegate to the appropriate 
expert rather than trying to answer directly.
```

---

## Quick Reference: MCP Configuration Checklist

Before deploying an agent with MCP:

| Check | Status |
|-------|:------:|
| LLM supports MCP (OpenAI/Anthropic/Bedrock) | ☐ |
| MCP URL is HTTPS | ☐ |
| Authentication credentials are valid | ☐ |
| Tested MCP with cURL/Postman | ☐ |
| Tools are well-described | ☐ |
| Agent prompt mentions available tools | ☐ |
| Error handling in prompt | ☐ |
| Not too many tools (< 30) | ☐ |

---

## Real-World Example: Claims Processing Team

```yaml
Team: Insurance Claims
Channel: #claims-processing

Agents:
  - name: elen.facilitator
    role: Facilitator
    mcps:
      - kirha (web search for external info)
    prompt: "Coordinate the team, delegate to experts, close when resolved"

  - name: claudine.crm
    role: Expert
    mcps:
      - hubspot (customer data)
      - internal-crm (policy holder info)
    prompt: "You are the CRM expert. Find customer information when asked."

  - name: victor.policy
    role: Expert
    mcps:
      - markitdown (read policy PDFs)
      - policy-db (search coverage rules)
    prompt: "You are the policy expert. Check coverage and policy terms."

  - name: ida.claims
    role: Expert
    mcps:
      - claims-system (claim status, history)
      - payments (reimbursement info)
    prompt: "You are the claims expert. Process and track claims."
```

This team can handle complex claims by:
1. Facilitator receives claim question
2. Delegates to CRM expert for customer info
3. Delegates to Policy expert for coverage check
4. Delegates to Claims expert for processing
5. Facilitator synthesizes and responds