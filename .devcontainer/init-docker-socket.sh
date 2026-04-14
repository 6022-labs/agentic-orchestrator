#!/bin/bash
set -euo pipefail

SOCKET=/var/run/docker.sock

if [ ! -S "$SOCKET" ]; then
  echo "init-docker-socket: $SOCKET not mounted, skipping"
  exit 0
fi

SOCK_GID=$(stat -c '%g' "$SOCKET")

if ! getent group "$SOCK_GID" >/dev/null; then
  groupadd -g "$SOCK_GID" docker-host
fi

GROUP_NAME=$(getent group "$SOCK_GID" | cut -d: -f1)

if ! id -nG node | tr ' ' '\n' | grep -qx "$GROUP_NAME"; then
  usermod -aG "$GROUP_NAME" node
fi

echo "init-docker-socket: node added to group $GROUP_NAME (gid $SOCK_GID) for $SOCKET"
