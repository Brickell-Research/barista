#!/bin/sh
set -e

git config --global user.email "barista@brickellresearch.com"
git config --global user.name "Barista"

# Rewrite git@github.com: URLs to authenticated HTTPS everywhere (clone, push, etc.)
git config --global url."https://x-access-token:${GITHUB_TOKEN}@github.com/".insteadOf "git@github.com:"

if [ -n "$OUTPUT_REPO" ] && [ ! -d output/.git ]; then
    git clone "$OUTPUT_REPO" output
fi

exec explore
