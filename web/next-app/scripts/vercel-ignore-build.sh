#!/usr/bin/env bash
# Vercel runs this from the project root. Build only when the web app or a
# source it reads at build time changed. Pull-request previews are disabled in
# vercel.json, so this guards production builds from Go-only and report-only
# commits, which do not change the exported site.
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
cd "$root"

if git diff --quiet HEAD^ HEAD -- web/next-app data docs; then
  echo "No web, data or docs changes since the previous commit; skipping build."
  exit 0
fi

echo "Web-relevant changes found; continuing with the build."
exit 1
