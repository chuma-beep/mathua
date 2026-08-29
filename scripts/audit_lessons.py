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
    concept_sources = {}
    for m in entries:
        mapped.add(m["concept_id"])
        sources.setdefault(m["source"], []).append(m["concept_id"])
        concept_sources.setdefault(m["concept_id"], []).append(m["source"])

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

    # 5. KP shards: every section must resolve in the concept's lesson body.
    import re as _re
    kp_dir = os.path.join(LESSONS, "kp")
    if os.path.isdir(kp_dir):
        kp_orphans = []
        for name in sorted(os.listdir(kp_dir)):
            if not name.endswith(".json"):
                continue
            cid = name[:-5]
            if cid not in dag_ids:
                kp_orphans.append(f"{name} (concept not in DAG)")
                continue
            with open(os.path.join(kp_dir, name)) as f:
                kps = json.load(f)
            srcs = concept_sources.get(cid, [])
            body = ""
            # teaching/* wins per loader (reverse-lexicographic pick)
            for src in sorted(srcs)[::-1]:
                p = os.path.join(LESSONS, src)
                if os.path.exists(p):
                    body = open(p, encoding="utf-8").read()
                    break
            heads = set()
            for line in body.split("\n"):
                t = line.strip()
                m = _re.match(r"^(#{2,4})\s+(.+)$", t)
                if m:
                    heads.add(m.group(2).strip())
            for kp in kps:
                if kp.get("section") and kp["section"] not in heads:
                    kp_orphans.append(f"{name}: section {kp['section']!r} unresolved")
        if kp_orphans:
            errors.append(f"kp shard problems ({len(kp_orphans)}):\n  " + "\n  ".join(kp_orphans))

    if errors:
        print("\n\n".join(errors))
        print(f"\nFAIL: {sum(1 for _ in errors)} categories with issues")
        sys.exit(1)

    print(f"OK: {len(dag_ids)} concepts mapped, {len(entries)} entries, "
          f"{len(sources)} sources, kp shards resolved, 0 issues")


if __name__ == "__main__":
    main()
