#!/bin/bash

# prefixtree is this one https://gist.github.com/aziis98/87ab5c533009dde2245bd083c4a4874c

zoxide query -l -s \
  | grep "$HOME" \
  | sed "s|$HOME|~|" \
  | awk '{ print $2 " " $1 }' \
  | fzf --filter="$prompt" -n1 \
  | head -n20 \
  | sort -r -k2 \
  | prefixtree -d "/"
