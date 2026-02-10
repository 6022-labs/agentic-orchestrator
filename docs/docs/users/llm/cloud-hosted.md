---
sidebar_position: 5
---

# Cloud Hosted (AWS, Azure, GCP)

Deploy open-source LLM models on cloud infrastructure you control. This provides more control over your data while leveraging managed services for scalability and reliability.

## Supported Platforms

| Platform | Service | Models Available |
|----------|---------|------------------|
| **AWS** | Bedrock | Llama, Mistral, Claude, Titan |
| **Azure** | Azure AI | Llama, Mistral, Phi |
| **Google Cloud** | Vertex AI | Gemma, Llama |

---

## AWS Bedrock

### Overview

Amazon Bedrock is a fully managed service that offers foundation models from leading AI companies through a single API.

### Configuration

```yaml
Provider: bedrock
Region: us-east-1 (or your region)
Model: anthropic.claude-3-sonnet-20240229-v1:0
Authentication: AWS Credentials
Access Key: AKIA...
Secret Key: [your-secret-key]
```

### Available Models

| Provider | Model ID | Use Case |
|----------|----------|----------|
| Anthropic | `anthropic.claude-3-sonnet-20240229-v1:0` | General purpose |
| Anthropic | `anthropic.claude-3-haiku-20240307-v1:0` | Fast, cost-effective |
| Meta | `meta.llama3-70b-instruct-v1:0` | Open source |
| Mistral | `mistral.mistral-large-2402-v1:0` | European AI |
| Amazon | `amazon.titan-text-express-v1` | Amazon native |

### Required Credentials

1. **Access Key ID**: Your AWS access key
2. **Secret Access Key**: Your AWS secret key
3. **Region**: AWS region where Bedrock is enabled
4. **(Optional) Session Token**: For temporary credentials

### Setup Steps

1. **Enable Bedrock** in your AWS account
2. **Request model access** for the models you want to use
3. **Create IAM credentials** with Bedrock permissions
4. **Configure in orchestrator** with the credentials

:::warning API Key Expiration
AWS Bedrock API keys typically expire after **90 days**. Set reminders to rotate them regularly.
:::

---

## Azure AI

### Configuration

```yaml
Provider: azure
API URL: https://your-resource.openai.azure.com/
Model: gpt-4
API Version: 2024-02-15-preview
Authentication: API Key
API Key: [your-azure-key]
```

### Setup Steps

1. Create an **Azure OpenAI resource**
2. Deploy the desired model
3. Get the **endpoint URL** and **API key**
4. Configure in orchestrator

---

## Google Cloud Vertex AI

### Configuration

```yaml
Provider: vertex
Project ID: your-gcp-project
Region: us-central1
Model: gemini-pro
Authentication: Service Account
Credentials: [path-to-service-account.json]
```

### Setup Steps

1. Enable **Vertex AI API** in your GCP project
2. Create a **service account** with Vertex AI permissions
3. Download the **credentials JSON**
4. Configure in orchestrator

---

## Pros & Cons

**Advantages:**
- ✅ Data stays within your cloud environment
- ✅ Choice of open-source and commercial models
- ✅ Predictable infrastructure costs
- ✅ Compliance with data residency requirements
- ✅ High availability and automatic scaling
- ✅ Enterprise-grade infrastructure security

**Limitations:**
- ❌ **No data confidentiality guarantee** (cloud admins could access data)
- ❌ **No model verification** (you trust your cloud setup)
- ❌ Requires cloud infrastructure expertise
- ❌ Infrastructure management overhead
- ❌ Higher cost than centralized for low volume
- ❌ No microtransaction support

---

## Best Practices

### Cost Optimization

1. **Use appropriate model sizes**: Don't use 70B when 7B suffices
2. **Enable auto-scaling**: Scale down during off-hours
3. **Monitor usage**: Set up billing alerts
4. **Use spot instances**: For non-critical workloads

### Security

1. **Use IAM roles**: Avoid long-lived credentials when possible
2. **Enable VPC**: Keep traffic within your network
3. **Encrypt at rest**: Enable encryption for stored data
4. **Audit access**: Enable CloudTrail/logging

### Compliance

1. **Choose the right region**: For data residency
2. **Enable compliance features**: HIPAA, SOC2, etc.
3. **Document data flows**: For audits

---

## Troubleshooting

### Authentication Errors

**Problem:** 401/403 errors when connecting

**Solutions:**
1. Verify credentials are correct
2. Check IAM permissions include Bedrock access
3. Ensure the region is correct
4. For AWS, verify the model is enabled in your account

### Model Not Available

**Problem:** Model ID not recognized

**Solutions:**
1. Request access to the model in AWS console
2. Check the model ID is spelled correctly
3. Verify the model is available in your region

### Timeout Errors

**Problem:** Requests timing out

**Solutions:**
1. Check network connectivity to cloud provider
2. Verify VPC/firewall settings allow outbound traffic
3. Increase timeout settings in orchestrator