---
sidebar_position: 4
---

# Local Models with Ollama

Run open-source LLMs on your own hardware via [Ollama](https://ollama.ai). Zero per-request cost, nothing leaves your machine, and the orchestrator talks to it over an OpenAI-compatible API.

**Pros:** complete data confidentiality, no API cost, no internet needed after the first download. **Cons:** requires capable hardware, slower than cloud APIs, you're on your own for model updates, and the machine is a single point of failure.

## Setup

### 1. Install Ollama

```bash
# macOS / Linux
curl -fsSL https://ollama.ai/install.sh | sh
```

Windows: download from [ollama.ai/download](https://ollama.ai/download).

### 2. Pull a model

```bash
ollama pull llama3        # 8B, general purpose, 8 GB RAM
ollama pull mistral       # 7B, fast, 8 GB RAM
ollama pull mixtral       # 8x7B, higher quality, 32 GB RAM
ollama list               # see what you have
```

### 3. Run the server

```bash
ollama serve              # listens on http://localhost:11434
```

### 4. Configure in the orchestrator

```yaml
Provider: ollama
API URL: http://localhost:11434/v1
Model: llama3
Authentication: None
```

## Hardware requirements

| Model size | RAM | GPU (recommended) |
|---|---|---|
| 3B (Phi-3) | 4 GB | Optional — runs on CPU |
| 7B | 8 GB | 6 GB VRAM |
| 13B | 16 GB | 8 GB VRAM |
| 34B | 32 GB | 16 GB VRAM |
| 70B | 64 GB | 24 GB+ VRAM |

A dedicated NVIDIA GPU with CUDA makes a significant difference to response times. Ollama uses the GPU automatically when available — check with `ollama run llama3 --verbose` (look for "using GPU").

## Remote access

To let the orchestrator reach Ollama running on a different machine:

```bash
OLLAMA_HOST=0.0.0.0 ollama serve
```

Then point the orchestrator at the machine's LAN IP instead of localhost: `API URL: http://192.168.1.100:11434/v1`.

## Custom model configuration

If you want tuned parameters or a baked-in system prompt, create a Modelfile:

```bash
cat << 'EOF' > Modelfile
FROM llama3
PARAMETER temperature 0.7
PARAMETER num_ctx 4096
SYSTEM "You are a helpful assistant for the 6022 protocol."
EOF

ollama create my-agent-model -f Modelfile
```

Then use `Model: my-agent-model` in the orchestrator.

## Supported models

Llama 3, Mistral, Mixtral, Qwen, Phi-3, CodeLlama, and anything else in the [Ollama library](https://ollama.ai/library). RAM needs scale with parameter count — see the table above.

## Troubleshooting

- **Connection refused** — `ollama serve` isn't running, or you're trying to reach it from another machine without `OLLAMA_HOST=0.0.0.0`. Test with `curl http://localhost:11434/api/tags`.
- **Out of memory when loading the model** — use a smaller variant (7B instead of 70B), close other applications, or add RAM / GPU.
- **Slow responses** — add GPU acceleration, shrink the context window in your prompts, or switch to a quantised model (`q4_0`, `q5_1`).
- **Model not found** — check `ollama list` to confirm the model is pulled locally, and verify the spelling in the orchestrator matches.
