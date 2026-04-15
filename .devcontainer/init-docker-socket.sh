#!/bin/bash
set -euo pipefail

# Force `localhost` to resolve to IPv4 only. Node.js picks the first matching
# entry from /etc/hosts, and Linux containers list ::1 first — so Claude Code's
# OAuth callback server binds to IPv6 only, while VS Code's port-forwarding
# tunnel connects via IPv4 and hangs. See anthropics/claude-code#44844.
if grep -qE '^::1[[:space:]].*\blocalhost\b' /etc/hosts; then
  # /etc/hosts is a bind-mounted file in Docker, so `sed -i` can't rename over
  # it ("Device or resource busy"). Rewrite in place via cat instead.
  tmp=$(mktemp)
  grep -vE '^::1[[:space:]].*\blocalhost\b' /etc/hosts > "$tmp"
  cat "$tmp" > /etc/hosts
  rm -f "$tmp"
  echo "init-docker-socket: removed ::1 localhost from /etc/hosts"
fi

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
