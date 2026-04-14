#!/bin/bash
set -euo pipefail

SOCKET=/var/run/docker.sock

if [ ! -S "$SOCKET" ]; then
  echo "init-docker-socket: $SOCKET not mounted, skipping"
  exit 0
fi

# Grant rw to all users on the socket. usermod -aG only takes effect on new
# login sessions, which is unreliable here: VS Code terminals may be spawned
# before postStartCommand completes, and docker group GID on Docker Desktop
# for Mac is 0 (root) which makes the group workaround semantically wrong.
# Single-user dev container — permissive mode is acceptable.
chmod 666 "$SOCKET"
echo "init-docker-socket: chmod 666 $SOCKET"
