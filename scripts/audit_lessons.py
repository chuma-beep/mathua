#!/usr/bin/env python3
"""Audit the lesson corpus: every DAG concept mapped, no stale IDs,
no missing files, no orphaned sources.

Usage: python3 scripts/audit_lessons.py
Exit 0 = healthy, exit 1 = issues found.
"""
import json
import os
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
CONCEPTS_DIR = os.path.join(ROOT, "data", "concepts")
LESSONS = os.path.join(ROOT, "data", "lessons")
LESSONS_JSON = os.path.join(LESSONS, "lessons.json")


def load_dag_ids():
    ids = set()
    for name in os.listdir(CONCEPTS_DIR):
        if not name.endswith(".json") or name == "enrichment.json":
            continue
        with open(os.path.join(CONCEPTS_DIR, name)) as f:
            for c in json.load(f):
                ids.add(c["id"])
    return ids


def main():
    errors = []

    dag_ids = load_dag_ids()
    with open(LESSONS_JSON) as f:
        entries = json.load(f)

    mapped = set()
    sources = {}
    for m in entries:
        mapped.add(m["concept_id"])
        sources.setdefault(m["source"], []).append(m["concept_id"])

    # 1. DAG concepts with no lesson mapping.
    orphans = sorted(dag_ids - mapped)
    if orphans:
        errors.append(f"orphan DAG concepts (no lesson): {len(orphans)}\n  " + ", ".join(orphans))

    # 2. lessons.json ids that do not exist in the DAG (stale).
    stale = sorted(mapped - dag_ids)
    if stale:
        errors.append(f"stale lessons.json ids (not in DAG): {len(stale)}\n  " + ", ".join(stale))

    # 3. Sources that fail to load on disk.
    missing = sorted(s for s in sources if not os.path.exists(os.path.join(LESSONS, s)))
    if missing:
        errors.append(f"lesson files missing on disk: {len(missing)}\n  " + "\n  ".join(missing))

    # 4. Sources referenced only by stale ids (dead weight after prune).
    dead = sorted(s for s, ids in sources.items() if not (set(ids) & dag_ids))
    if dead:
        errors.append(f"orphaned sources (only stale ids reference them): {len(dead)}\n  " + "\n  ".join(dead))

    if errors:
        print("\n\n".join(errors))
        print(f"\nFAIL: {sum(1 for _ in errors)} categories with issues")
        sys.exit(1)

    print(f"OK: {len(dag_ids)} concepts mapped, {len(entries)} entries, "
          f"{len(sources)} sources, 0 issues")


if __name__ == "__main__":
    main()
