#!/usr/bin/env bash
set -e

# fns-cli
sudo ln -sf $HOME/.fns-cli/container.bash_profile $HOME/fns-cli.bash_profile
grep -qxF 'source ~/fns-cli.bash_profile' "$HOME/.bashrc" || printf '\nsource ~/fns-cli.bash_profile\n' >>"$HOME/.bashrc"
