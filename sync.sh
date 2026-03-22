#!/usr/bin/env bash
set -euo pipefail

echo "Fetching latest from upstream..."
git fetch upstream

echo "Updating main with upstream/main..."
git checkout main
git merge upstream/main --ff-only

echo "Rebasing my-own-branch on top of main..."
git checkout my-own-branch
git rebase main

echo "Done. Run ./docker/up.sh to rebuild."
