# Devcontainer Setup

This devcontainer runs Claude Code in a sandboxed environment with a restrictive firewall (see `init-firewall.sh`). VSCode automatically forwards your host SSH agent into the container, so Git operations against GitHub work without copying your private key inside.

## One-time host setup (macOS)

Your SSH key must be loaded into the macOS ssh-agent **on your Mac**, not inside the container. Do this once — it persists across reboots thanks to the Keychain.

### 1. Configure `~/.ssh/config`

Ensures the key auto-loads from Keychain on every boot:

```bash
cat >> ~/.ssh/config <<'EOF'

Host *
  AddKeysToAgent yes
  UseKeychain yes
  IdentityFile ~/.ssh/id_rsa
EOF
chmod 600 ~/.ssh/config
```

Replace `id_rsa` with your actual key filename if different (`ls ~/.ssh/id_*`).

### 2. Load the key into the agent + Keychain

```bash
ssh-add --apple-use-keychain ~/.ssh/id_rsa
```

Enter passphrase when prompted (or just press Enter if the key has none). The passphrase is now stored in macOS Keychain.

### 3. Verify on the host

```bash
ssh-add -l
```

You should see a line like:

```
4096 SHA256:xxxxxxxxxxxxxxxxxxxx /Users/you/.ssh/id_rsa (RSA)
```

## Verify inside the devcontainer

Open the devcontainer in VSCode, then in its terminal:

```bash
ssh-add -l           # should show the same fingerprint as the host
ssh -T git@github.com # should print: "Hi <username>! You've successfully authenticated..."
git pull             # should work
```

## How it works

```
[~/.ssh/id_rsa on Mac] → [macOS ssh-agent] → [VSCode forwards agent socket] → [devcontainer]
```

- The private key **never enters** the container.
- The container only gets a socket it can ask for signatures.
- VSCode sets `SSH_AUTH_SOCK` inside the container automatically (path looks like `/tmp/vscode-ssh-auth-*.sock`).

## Recurring workflow

| Event | Action required |
|-------|-----------------|
| Daily use / mounting the devcontainer | Nothing |
| Mac reboot | Nothing (Keychain auto-loads the key) |
| Rebuilding the devcontainer | Nothing |
| New SSH key | Repeat steps 1–3 above with the new key path |

## Troubleshooting

**`git pull` → `Permission denied (publickey)`**

Run inside the container:

```bash
echo $SSH_AUTH_SOCK    # must be non-empty
ssh-add -l             # must list at least one key
```

- Empty `SSH_AUTH_SOCK` → VSCode agent forwarding broken; rebuild container.
- "The agent has no identities" → host agent empty; run step 2 on the host.
- "Could not open a connection to your authentication agent" → socket not reachable; reopen container in VSCode.

**Host check**

On the Mac, confirm the default launchd agent is active:

```bash
echo $SSH_AUTH_SOCK
# expected: /private/tmp/com.apple.launchd.XXXX/Listeners
```

If the path points to `~/.ssh/agent/` sockets or a `1password` path, you are using a non-default SSH agent and the setup differs — see your agent's documentation for how to expose keys.

## Notes on Claude running unrestricted

Claude inside this container can *use* the forwarded socket to sign SSH challenges (e.g. `git push`) during its session. It **cannot** read the private key bytes — only request signatures. To revoke access at any time, on the host:

```bash
ssh-add -D          # remove all keys from the running agent
```

The Keychain still holds the passphrase, so the next `ssh` call reloads the key automatically.
