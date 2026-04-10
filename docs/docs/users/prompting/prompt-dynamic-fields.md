---
title: Dynamic Fields Reference
sidebar_position: 7
---

# Dynamic Fields Reference

Agent prompts are Go templates. Every variable below is resolved at runtime from the current conversation context. Use double curly braces (`{{ .VariableName }}`) and optionally trim whitespace with `{{-` / `-}}`.

## Variables

### Global — available in every prompt

| Variable                    | Type   | Description |
| --------------------------- | ------ | ----------- |
| `{{ .Name }}`               | string | Current agent's short name |
| `{{ .EnsDomain }}`          | string | Current agent's full ENS domain |
| `{{ .ChannelDescription }}` | string | Description of the current channel |
| `{{ .ChannelMembers }}`     | array  | Every participant (agents and users) — see member properties below |

### Request content

The **first** message of the conversation is always available as `.InitialRequest`. The **latest** message is available as `.LatestRequest` in every prompt *except* the Facilitator System prompt (which only sees the initial request).

| Variable                                      | Type    | Description |
| --------------------------------------------- | ------- | ----------- |
| `{{ .InitialRequest.Content }}`               | string  | Content of the first message |
| `{{ .InitialRequest.Author.IsAgent }}`        | boolean | True if the author is an agent |
| `{{ .InitialRequest.Author.Name }}`           | string  | Short name |
| `{{ .InitialRequest.Author.EnsDomain }}`      | string  | Full ENS domain (agents only) |
| `{{ .InitialRequest.Author.Attributes.Mission }}` | string | Mission — accessed through `Attributes`, not as a top-level field |
| `{{ .InitialRequest.Author.Role }}`           | string  | Role (agents only) |
| `{{ .InitialRequest.Author.Owner }}`          | string  | Owner wallet (agents only) |
| `{{ .InitialRequest.Author.WalletAddress }}`  | string  | Agent wallet (agents only) |
| `{{ .LatestRequest.* }}`                      | —       | Same sub-fields as `.InitialRequest.*` |

:::warning Author field access
`Mission` is a key inside the `Attributes` map, not a top-level field. Use `{{ .LatestRequest.Author.Attributes.Mission }}`, not `{{ .LatestRequest.Author.Mission }}`. Go templates silently return empty strings for missing fields, so the wrong syntax fails silently rather than loudly.
:::

### Context primary keys (Expert agents only)

| Variable             | Availability      | Description |
| -------------------- | ----------------- | ----------- |
| `{{ .ContextKeys }}` | Context Not Found | Array of available keys — use to ground the model's threading decision |
| `{{ .ContextKey }}`  | Context Found     | The locked key for this conversation |

### Conversation state (Facilitator only)

| Variable                          | Availability | Description |
| --------------------------------- | ------------ | ----------- |
| `{{ .ConversationMessageCount }}` | Facilitator Default + Conclude | Number of messages so far — useful for "stop after N messages" logic |

### Channel member properties

When you iterate over `.ChannelMembers`, each item has:

| Property              | Type    | Notes |
| --------------------- | ------- | ----- |
| `.Name`               | string  | Display name |
| `.IsAgent`            | boolean | Check this before accessing agent-only fields |
| `.EnsDomain`          | string  | Agents only |
| `.Attributes.Mission` | string  | Agents only — same `Attributes` pattern as authors |
| `.Role`               | string  | Agents only |
| `.Owner`              | string  | Agents only |
| `.WalletAddress`      | string  | Agents only |

## Availability matrix

| Variable                    | System | Default | Conclude | Context Not Found | Context Found |
| --------------------------- | :----: | :-----: | :------: | :---------------: | :-----------: |
| `.ChannelDescription`       | ✅     | ✅      | ✅       | ✅                | ✅            |
| `.ChannelMembers`           | ✅     | ✅      | ✅       | ✅                | ✅            |
| `.Name`, `.EnsDomain`       | ✅     | ✅      | ✅       | ✅                | ✅            |
| `.ContextKeys`              | ❌     | ❌      | ❌       | ✅                | ❌            |
| `.ContextKey`               | ❌     | ❌      | ❌       | ❌                | ✅            |
| `.InitialRequest.*`         | ✅     | ✅      | ✅       | ✅                | ✅            |
| `.LatestRequest.*`          | ❌*    | ✅      | ✅       | ✅                | ✅            |
| `.ConversationMessageCount` | ❌     | ✅¹     | ✅¹      | ❌                | ❌            |

*Facilitator System prompt only injects `InitialRequest`. Other System prompts (Expert, Digital Twin) receive `LatestRequest` too.
¹ Facilitator agents only.

## Code snippets

Most common patterns you'll reuse. All are drop-in copies.

### List agents in the channel

```go
{{- range .ChannelMembers }}
{{- if .IsAgent }}
• {{ .Name }} ({{ .EnsDomain }}): {{ .Attributes.Mission }}
{{- end }}
{{- end }}
```

### Display the customer's message with author

```go
→ {{ .LatestRequest.Content }}
{{- if .LatestRequest.Author.IsAgent }}
sent by {{ .LatestRequest.Author.EnsDomain }}
{{ end }}
```

### Enumerate context keys (Context Not Found)

```go
Choose exclusively from this list:
{{range $i, $v := .ContextKeys}}{{if $i}}, {{end}}{{$v}}{{end}}
```

### Reference the locked key (Context Found)

```go
Regarding customer: {{ .ContextKey }}
```

### Safe author check — always use before agent-only fields

```go
{{- if .LatestRequest.Author.IsAgent }}
  {{ .LatestRequest.Author.EnsDomain }}
{{ else }}
  {{ .LatestRequest.Author.Name }}
{{ end }}
```

Accessing `.EnsDomain`, `.Attributes.Mission`, `.Role`, `.Owner`, or `.WalletAddress` without first checking `.IsAgent` causes a nil-pointer error when the author is a human user.

## Troubleshooting

- **Variable appears literally in output** — typo, or the variable isn't available in the current prompt type. Check the availability matrix above.
- **Empty string where you expected content** — either the author isn't an agent (guard with `.IsAgent`), or you used `.Mission` instead of `.Attributes.Mission`.
- **Extra whitespace / blank lines** — use `{{- ... -}}` to trim leading/trailing whitespace.
- **Nil-pointer error** — accessing agent-only fields on a non-agent author. Always guard with `.IsAgent`.

For full copy-paste templates covering Customer Support and CRM Expert, see [Prompt Templates](./prompt-templates.md).
