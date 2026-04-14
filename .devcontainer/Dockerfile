FROM node:20

ARG TZ
ENV TZ="$TZ"

ARG CLAUDE_CODE_VERSION=latest
ARG GO_VERSION=1.25.5

# Install basic development tools and iptables/ipset
RUN apt-get update && apt-get install -y --no-install-recommends \
  less \
  git \
  procps \
  sudo \
  fzf \
  zsh \
  man-db \
  unzip \
  gnupg2 \
  gh \
  iptables \
  ipset \
  iproute2 \
  dnsutils \
  aggregate \
  jq \
  nano \
  vim \
  openssh-client \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# Pre-populate github.com host keys so ssh-agent-forwarded git calls don't prompt
RUN mkdir -p /etc/ssh && \
  ssh-keyscan -t rsa,ecdsa,ed25519 github.com >> /etc/ssh/ssh_known_hosts

# Ensure default node user has access to /usr/local/share
RUN mkdir -p /usr/local/share/npm-global && \
  chown -R node:node /usr/local/share

ARG USERNAME=node

# Persist bash history.
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history \
  && chown -R $USERNAME /commandhistory

# Set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# Create workspace and config directories and set permissions
RUN mkdir -p /workspace /home/node/.claude /home/node/.vscode-server && \
  chown -R node:node /workspace /home/node/.claude /home/node/.vscode-server

WORKDIR /workspace

ARG GIT_DELTA_VERSION=0.18.2
RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/dandavison/delta/releases/download/${GIT_DELTA_VERSION}/git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  sudo dpkg -i "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  rm "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb"

# Install Go
RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz" && \
  tar -C /usr/local -xzf "go${GO_VERSION}.linux-${ARCH}.tar.gz" && \
  rm "go${GO_VERSION}.linux-${ARCH}.tar.gz"
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/home/node/go
ENV PATH=$PATH:/home/node/go/bin

# Set up non-root user
USER node

# Install global packages
ENV NPM_CONFIG_PREFIX=/usr/local/share/npm-global
ENV PATH=$PATH:/usr/local/share/npm-global/bin

# Set the default shell to zsh rather than sh
ENV SHELL=/bin/zsh

# Set the default editor and visual
ENV EDITOR=nano
ENV VISUAL=nano

# Default powerline10k theme
ARG ZSH_IN_DOCKER_VERSION=1.2.0
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v${ZSH_IN_DOCKER_VERSION}/zsh-in-docker.sh)" -- \
  -p git \
  -p fzf \
  -a "source /usr/share/doc/fzf/examples/key-bindings.zsh" \
  -a "source /usr/share/doc/fzf/examples/completion.zsh" \
  -a "export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  -x

# Install Claude
RUN npm install -g @anthropic-ai/claude-code@${CLAUDE_CODE_VERSION}

# Install Railway CLI (direct binary — npm postinstall 404s on aarch64 because only musl assets are published)
ARG RAILWAY_VERSION=4.36.1
RUN set -eux; \
  case "$(dpkg --print-architecture)" in \
    amd64) RAILWAY_TARGET=x86_64-unknown-linux-musl ;; \
    arm64) RAILWAY_TARGET=aarch64-unknown-linux-musl ;; \
    *) echo "unsupported arch"; exit 1 ;; \
  esac; \
  mkdir -p /tmp/railway-extract; \
  wget -qO /tmp/railway.tar.gz "https://github.com/railwayapp/cli/releases/download/v${RAILWAY_VERSION}/railway-v${RAILWAY_VERSION}-${RAILWAY_TARGET}.tar.gz"; \
  tar -xzf /tmp/railway.tar.gz -C /tmp/railway-extract; \
  RAILWAY_BIN=$(find /tmp/railway-extract -type f -name railway | head -n1); \
  mv "$RAILWAY_BIN" /usr/local/share/npm-global/bin/railway; \
  chmod +x /usr/local/share/npm-global/bin/railway; \
  rm -rf /tmp/railway.tar.gz /tmp/railway-extract

# Preinstall Go tools used by the VS Code Go extension (firewall blocks proxy.golang.org at runtime)
RUN go install github.com/cweill/gotests/gotests@v1.6.0 && \
  go install github.com/josharian/impl@v1.4.0 && \
  go install github.com/haya14busa/goplay/cmd/goplay@v1.0.0 && \
  go install github.com/go-delve/delve/cmd/dlv@latest && \
  go install golang.org/x/tools/gopls@latest && \
  go clean -cache -modcache -testcache

# Copy and set up firewall script
COPY init-firewall.sh /usr/local/bin/
USER root
RUN chmod +x /usr/local/bin/init-firewall.sh && \
  echo "node ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/node-firewall && \
  chmod 0440 /etc/sudoers.d/node-firewall
USER node
