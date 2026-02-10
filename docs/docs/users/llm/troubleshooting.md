---
sidebar_position: 7
---

# Troubleshooting & Advanced Configuration

This page covers common issues and advanced configuration options for LLM management in the 6022 protocol.

## Token Management

### MAX_TOKENS_PER_REQUEST

The orchestrator includes a `MAX_TOKENS_PER_REQUEST` environment variable to limit context size.

| Setting | Value | Use Case |
|---------|-------|----------|
| Default | 100,000 | Most providers |
| Anthropic recommended | 50,000 | Due to different token calculation |
| Low memory | 25,000 | For smaller models |

:::note Token Calculation
Anthropic uses a different token calculation method than OpenAI. If you experience context size errors with Anthropic, reduce the max tokens value.
:::

### Adjusting Token Limits

In Railway, update the environment variable in the **LLM Connector** service:

```
MAX_TOKENS_PER_REQUEST=50000
```

---

## Rate Limiting

### Understanding Rate Limits

LLM providers enforce rate limits to prevent abuse. Common limits:

| Provider | Typical Limit | Reset Period |
|----------|---------------|--------------|
| OpenAI | 10,000 TPM* | Per minute |
| Anthropic | 100,000 TPM | Per minute |
| Mistral | Varies by tier | Per minute |

*TPM = Tokens Per Minute

### HTTP Status Codes

| Code | Meaning | Solution |
|------|---------|----------|
| `429` | Too many requests | Wait and retry, upgrade plan |
| `401` | Unauthorized | Check API key |
| `403` | Forbidden | Check permissions/plan |
| `500` | Server error | Retry the request |
| `503` | Service unavailable | Provider is down, wait |

### Handling Rate Limits

When multiple agents share the same API key:

1. **Use separate API keys** per agent
2. **Upgrade your provider plan** for higher limits
3. **Implement request delays** in high-volume scenarios
4. **Distribute across providers** (some agents on OpenAI, others on Mistral)

---

## Monitoring & Debugging

### Accessing Logs in Railway

1. Go to the **LLM Connector** service
2. Open the **Logs** tab
3. Use filters to find issues:

| Filter | Purpose |
|--------|---------|
| `@level:warn` | Warnings and errors |
| `@level:error` | Errors only |
| `@status_code:429` | Rate limit issues |
| `@status_code:401` | Auth issues |

### Common Log Patterns

```
# Insufficient quota
"error": "insufficient_quota"
→ Top up your provider account

# Rate limit
"error": "rate_limit_exceeded"
→ Too many concurrent requests

# Context too long
"error": "context_length_exceeded"
→ Reduce MAX_TOKENS_PER_REQUEST or simplify prompts

# Invalid API key
"error": "invalid_api_key"
→ Check key format and validity
```

---

## Common Issues & Solutions

### Agent Not Responding

**Checklist:**
1. ✅ API key is valid and has credits
2. ✅ Model name is correct
3. ✅ Provider URL is correct
4. ✅ No rate limit errors in logs
5. ✅ Context size is within limits

### Anthropic-Specific Issues

**Problem:** Agent works without MCP but configuration fails

**Solution:** Use Custom Header authentication:
```yaml
Authentication: Custom Header
Header Name: x-api-key
Header Value: [your-api-key]
```

**Problem:** Context too large errors

**Solution:** Anthropic calculates tokens differently. Reduce `MAX_TOKENS_PER_REQUEST` to 50,000.

### AWS Bedrock Issues

**Problem:** Authentication fails

**Solutions:**
1. Check credentials haven't expired (90-day limit)
2. Verify IAM role has Bedrock permissions
3. Ensure model is enabled in your region

**Problem:** Model not available

**Solution:** Request model access in AWS Bedrock console.

### Ollama Issues

**Problem:** Connection refused

**Solutions:**
1. Verify Ollama is running: `ollama serve`
2. For remote access: `OLLAMA_HOST=0.0.0.0 ollama serve`
3. Check firewall allows port 11434

---

## Performance Optimization

### Choosing Model Size

| Scenario | Recommended Size |
|----------|-----------------|
| Simple Q&A, classification | 7B |
| General conversation | 13B |
| Complex reasoning | 34B+ |
| Code generation | CodeLlama 34B |

### Prompt Optimization

1. **Be concise**: Shorter prompts = faster responses
2. **Use system prompts wisely**: Don't repeat in every message
3. **Limit history**: Only include relevant conversation history
4. **Avoid redundant instructions**: Agents remember their role

### Context Window Management

For long conversations:

1. **Summarize periodically**: Have facilitator create summaries
2. **Prune old messages**: Keep only recent relevant context
3. **Use expert agents**: They only see @mentions, reducing context

---

## Security Best Practices

### API Key Management

| Practice | Description |
|----------|-------------|
| **Rotate regularly** | Every 90 days minimum |
| **Use separate keys** | Per environment (dev/prod) |
| **Never commit keys** | Use environment variables |
| **Monitor usage** | Set up billing alerts |

### Access Control

1. **Limit orchestrator access**: Only admins configure LLMs
2. **Use wallet authentication**: Only NFT holders modify agents
3. **Audit changes**: Track who changed what

---

## Getting Help

### Resources

- **Documentation**: This site
- **Telegram**: [t.me/protocol6022](https://t.me/protocol6022)
- **GitHub Issues**: [github.com/6022-labs](https://github.com/6022-labs)

### Information to Include in Support Requests

1. Agent ENS domain
2. LLM provider and model
3. Error message (from logs)
4. Timestamp of the issue
5. Steps to reproduce