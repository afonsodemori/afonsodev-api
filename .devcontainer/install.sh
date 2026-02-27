#!/usr/bin/env bash
set -e

if ! command -v gemini >/dev/null 2>&1; then
	echo "Installing Gemini CLI..."
	npm install -g @google/gemini-cli@latest
fi

if ! command -v shfmt &>/dev/null 2>&1; then
	echo "Installing shfmt..."
	go install mvdan.cc/sh/v3/cmd/shfmt@latest
fi

# Ensure cloudflared directory exists and is writable
mkdir -p /home/vscode/.cloudflared
sudo chown -R vscode:vscode /home/vscode/.cloudflared || true

# Ensure gemini directory exists and is writable
mkdir -p /home/vscode/.gemini
sudo chown -R vscode:vscode /home/vscode/.gemini || true

# export SSH_AUTH_SOCK=/ssh-agent && ssh-add -l || true
