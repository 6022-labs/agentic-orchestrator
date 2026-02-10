---
sidebar_position: 4
---

# Local Models with Ollama

Run open-source LLM models directly on your own hardware using Ollama. This provides complete data privacy with no external API calls, and zero per-request costs after initial setup.

## What is Ollama?

[Ollama](https://ollama.ai) is an open-source tool that makes it easy to run large language models locally. It handles model downloading, optimization, and exposes an OpenAI-compatible API.

## Supported Models

| Model | Parameters | RAM Required | Use Case |
|-------|------------|--------------|----------|
| **Llama 3** | 8B / 70B | 8GB / 48GB | General purpose |
| **Mistral** | 7B | 8GB | Fast, efficient |
| **Mixtral** | 8x7B | 32GB | High quality |
| **Qwen** | 7B / 72B | 8GB / 48GB | Multilingual |
| **Phi-3** | 3.8B | 4GB | Lightweight |
| **CodeLlama** | 7B / 34B | 8GB / 24GB | Code generation |

---

## Setup Process

### Step 1: Install Ollama

**macOS / Linux:**
```bash
curl -fsSL https://ollama.ai/install.sh | sh
```

**Windows:**
Download from [ollama.ai/download](https://ollama.ai/download)

---

### Step 2: Download a Model

```bash
# Download Llama 3 (8B) - recommended for most use cases
ollama pull llama3

# Download Mistral - fast and efficient
ollama pull mistral

# Download Mixtral - larger, better quality
ollama pull mixtral

# List downloaded models
ollama list
```

---

### Step 3: Start the Server

```bash
ollama serve
```

By default, Ollama runs on `http://localhost:11434`.

---

### Step 4: Configure in Orchestrator

```yaml
Provider: ollama
API URL: http://localhost:11434/v1
Model: llama3
Authentication: None
```

---

## Advanced Configuration

### Network Access

To allow the orchestrator to access Ollama from another machine:

```bash
OLLAMA_HOST=0.0.0.0 ollama serve
```

Then use your machine's IP address in the orchestrator configuration:

```yaml
API URL: http://192.168.1.100:11434/v1
```

### Custom Model Configuration

Create a custom model with specific parameters:

```bash
# Create a Modelfile
cat << 'EOF' > Modelfile
FROM llama3
PARAMETER temperature 0.7
PARAMETER num_ctx 4096
SYSTEM "You are a helpful assistant for the 6022 Agentic Teams protocol."
EOF

# Create the custom model
ollama create my-agent-model -f Modelfile
```

### GPU Acceleration

Ollama automatically uses GPU if available. To verify:

```bash
ollama run llama3 --verbose
# Look for "using GPU" in the output
```

---

## Hardware Requirements

| Model Size | Minimum RAM | Recommended GPU | Performance |
|------------|-------------|-----------------|-------------|
| 3B (Phi-3) | 4 GB | Optional | Fast on CPU |
| 7B | 8 GB | 6GB VRAM | Good |
| 13B | 16 GB | 8GB VRAM | Better quality |
| 34B | 32 GB | 16GB VRAM | High quality |
| 70B | 64 GB | 24GB+ VRAM | Best quality |

:::tip Performance Tip
For production use, a dedicated GPU significantly improves response times. NVIDIA GPUs with CUDA support provide the best performance.
:::

---

## Pros & Cons

**Advantages:**
- ✅ **Complete data confidentiality** (nothing leaves your machine)
- ✅ No API costs (only hardware/electricity)
- ✅ No internet required after model download
- ✅ Full control over model versions
- ✅ OpenAI-compatible API
- ✅ Easy setup

**Limitations:**
- ❌ Requires capable hardware (RAM/GPU)
- ❌ Slower than cloud APIs (depends on hardware)
- ❌ Manual model updates
- ❌ **No model verification** (you trust your own setup)
- ❌ No microtransaction support
- ❌ Single point of failure (your machine)

---

## Troubleshooting

### Connection Refused

**Problem:** Orchestrator can't connect to Ollama

**Solutions:**
1. Verify Ollama is running: `ollama serve`
2. Check the port: `curl http://localhost:11434/api/tags`
3. For remote access, ensure `OLLAMA_HOST=0.0.0.0`

### Out of Memory

**Problem:** Model fails to load

**Solutions:**
1. Use a smaller model (7B instead of 70B)
2. Close other applications to free RAM
3. Consider adding more RAM or using a GPU

### Slow Response Times

**Problem:** Agent responses are slow

**Solutions:**
1. Use a smaller model
2. Add GPU acceleration
3. Reduce context length in prompts
4. Use quantized models (q4_0, q5_1)

### Model Not Found

**Problem:** Orchestrator reports model not found

**Solutions:**
1. Verify model is downloaded: `ollama list`
2. Check model name spelling in configuration
3. Pull the model: `ollama pull <model-name>`