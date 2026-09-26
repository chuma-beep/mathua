#!/usr/bin/env python3
"""Report-only Algebrica sync: SVG twins, missing diagrams, missing lessons.

Compares vendored assets (web/next-app/public/diagrams/algebrica/),
engine mappings (internal/engine/engine.go conceptDiagrams) and lesson
mappings (data/lessons/lessons.json) against upstream main via ONE
recursive git-tree API call. Never modifies the tree; exit 0 always
(report, not gate). Set GITHUB_TOKEN to raise the rate limit.

Usage: python3 scripts/sync_algebrica.py [--owner antoni Lupetti ...]
"""
from __future__ import annotations

import json
import os
import re
import sys
import urllib.request
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
UPSTREAM = ("antoniolupetti", "algebrica")
IMAGE_EXTS = (".svg", ".png")
SKIP_DIRS = ("github-assets", ".github")


def api(path: str):
    url = f"https://api.github.com{path}"
    headers = {"Accept": "application/vnd.github+json"}
    token = os.environ.get("GITHUB_TOKEN")
    if token:
        headers["Authorization"] = f"Bearer {token}"
    req = urllib.request.Request(url, headers=headers)
    with urllib.request.urlopen(req, timeout=60) as res:
        return json.load(res)


def main() -> int:
    owner, repo = UPSTREAM
    try:
        tree = api(f"/repos/{owner}/{repo}/git/trees/main?recursive=1")
    except Exception as exc:
        print(f"sync failed (network/API): {exc}", file=sys.stderr)
        return 1
    if tree.get("truncated"):
        print("sync failed: upstream tree truncated, rerun with GITHUB_TOKEN", file=sys.stderr)
        return 1
    blobs = [e["path"] for e in tree.get("tree", []) if e["type"] == "blob"]

    up_images: dict[str, list[str]] = {}
    up_lessons: set[str] = set()
    for p in blobs:
        if p.split("/")[0] in SKIP_DIRS:
            continue
        low = p.lower()
        if low.endswith(IMAGE_EXTS):
            up_images.setdefault(Path(p).stem, []).append(p)
        elif low.endswith(".md") and Path(p).name.lower() not in ("readme.md", "license.md", "licence.md"):
            if p.startswith("category/"):
                continue  # category indexes, not lessons
            up_lessons.add(p)

    local_dir = ROOT / "web" / "next-app" / "public" / "diagrams" / "algebrica"
    local = {f.stem: f.suffix.lower() for f in local_dir.iterdir() if f.is_file()} if local_dir.is_dir() else {}

    eng = (ROOT / "internal" / "engine" / "engine.go").read_text()
    m = re.search(r"var conceptDiagrams = map\[string\]string\{(.*?)\n\}", eng, re.S)
    mapped = [b for _, b in re.findall(r'"([^"]+)":\s*"([^"]+)"', m.group(1))] if m else []
    mapped_local = [p for p in mapped if "/diagrams/algebrica/" in p]
    mapped_png = sorted({p for p in mapped_local if p.lower().endswith(".png")})

    entries = json.load(open(ROOT / "data" / "lessons" / "lessons.json"))
    local_sources = {e["source"] for e in entries}
    # local source "algebrica/<dir>/<file>.md" <-> upstream "<dir>/<file>.md"
    stripped = {s[len("algebrica/"):] if s.startswith("algebrica/") else s for s in local_sources}

    print(f"upstream: {len(up_images)} image stems, {len(up_lessons)} lesson files")
    print(f"local: {len(local)} vendored images, {len(mapped_local)} mapped, {len(mapped_png)} mapped PNG")
    print()
    print("== PNG -> SVG twins (same stem, upstream SVG exists) ==")
    twins = 0
    for p in mapped_png:
        stem = Path(p).stem
        ups = [u for u in up_images.get(stem, []) if u.lower().endswith(".svg")]
        if ups:
            twins += 1
            print(f"  {p}  <-  {ups[0]}")
    print(f"  ({twins} swappable)")
    print()
    print("== mapped PNG with NO upstream SVG twin ==")
    twinless = [p for p in mapped_png if not any(
        u.lower().endswith(".svg") for u in up_images.get(Path(p).stem, []))]
    for p in twinless:
        print(f"  {p}")
    print(f"  ({len(twinless)} stay PNG)")
    print()
    print("== upstream images with no local copy (missing diagrams) ==")
    missing_d = sorted({u for stem, us in up_images.items() if stem not in local for u in us})
    for u in missing_d[:40]:
        print(f"  {u}")
    if len(missing_d) > 40:
        print(f"  ... and {len(missing_d) - 40} more")
    print(f"  ({len(missing_d)} total)")
    print()
    print("== upstream lessons with no mapping (missing lessons — ADOPT ONLY ON APPROVAL) ==")
    missing_l = sorted(p for p in up_lessons if p not in stripped)
    for p in missing_l[:40]:
        print(f"  {p}")
    if len(missing_l) > 40:
        print(f"  ... and {len(missing_l) - 40} more")
    print(f"  ({len(missing_l)} total)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
