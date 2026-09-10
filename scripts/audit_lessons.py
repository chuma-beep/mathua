#!/usr/bin/env python3
"""Audit the lesson corpus: every DAG concept mapped, no stale IDs,
no missing files, no orphaned sources.

Usage: python3 scripts/audit_lessons.py
Exit 0 = healthy, exit 1 = issues found.
"""
import json
import os
import re
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

    # 2b. Concepts with multiple lesson sources (loader silently shadows all
    # but the reverse-lexicographic winner — resolve to exactly one entry).
    multi = sorted(cid for cid, srcs in concept_sources.items() if len(srcs) > 1)
    if multi:
        detail = "; ".join(f"{cid}: {sorted(concept_sources[cid])}" for cid in multi)
        errors.append(f"multi-source concepts (shadowing): {len(multi)}\n  " + detail)

    # 3. Sources that fail to load on disk.
    missing = sorted(s for s in sources if not os.path.exists(os.path.join(LESSONS, s)))
    if missing:
        errors.append(f"lesson files missing on disk: {len(missing)}\n  " + "\n  ".join(missing))

    # 4. Sources referenced only by stale ids (dead weight after prune).
    dead = sorted(s for s, ids in sources.items() if not (set(ids) & dag_ids))
    if dead:
        errors.append(f"orphaned sources (only stale ids reference them): {len(dead)}\n  " + "\n  ".join(dead))

    # 5. KP shards: every section must resolve in the concept's lesson body.
    kp_dir = os.path.join(LESSONS, "kp")
    kp_files = 0
    kp_total = 0
    if os.path.isdir(kp_dir):
        kp_orphans = []
        for name in sorted(os.listdir(kp_dir)):
            if not name.endswith(".json"):
                continue
            kp_files += 1
            cid = name[:-5]
            if cid not in dag_ids:
                kp_orphans.append(f"{name} (concept not in DAG)")
                continue
            with open(os.path.join(kp_dir, name)) as f:
                kps = json.load(f)
            if not isinstance(kps, list) or len(kps) != 3:
                kp_orphans.append(f"{name}: expected 3 KPs, got {len(kps) if isinstance(kps, list) else type(kps).__name__}")
            else:
                kp_total += len(kps)
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
                m = re.match(r"^(#{2,4})\s+(.+)$", t)
                if m:
                    heads.add(m.group(2).strip())
            for kp in kps:
                if kp.get("section") and kp["section"] not in heads:
                    kp_orphans.append(f"{name}: section {kp['section']!r} unresolved")
        if kp_orphans:
            errors.append(f"kp shard problems ({len(kp_orphans)}):\n  " + "\n  ".join(kp_orphans))
        # Validate shard file count and total KPs (630 files ×3 =1890)
        if kp_files != len(dag_ids):
            errors.append(f"kp shard count mismatch: {kp_files} files vs {len(dag_ids)} concepts")
        if kp_total != len(dag_ids) * 3:
            errors.append(f"kp total mismatch: {kp_total} KPs vs {len(dag_ids)*3} expected (3 per concept)")

    # 6. Diagram mappings (engine.conceptDiagrams): concept in DAG, asset on disk.
    engine_src = os.path.join(ROOT, "internal", "engine", "engine.go")
    if os.path.exists(engine_src):
        diag = []
        with open(engine_src) as f:
            eng = f.read()
        start = eng.find("var conceptDiagrams = map[string]string{")
        if start >= 0:
            end = eng.find("\n}", start)
            block = eng[start:end]
            for m in re.finditer(r'"([^"]+)":\s*"((?:/diagrams/algebrica/)[^"]+)"', block):
                cid, asset = m.group(1), m.group(2)
                if cid not in dag_ids:
                    diag.append(f"diagram {asset}: concept {cid} not in DAG")
                if not os.path.exists(os.path.join(ROOT, "web", "next-app", "public") + asset):
                    diag.append(f"diagram asset missing: {asset}")
        if diag:
            errors.append(f"diagram problems ({len(diag)}):\n  " + "\n  ".join(diag))

    # 7. Course catalog: every target must exist in the DAG.
    courses_path = os.path.join(ROOT, "data", "courses.json")
    if os.path.exists(courses_path):
        with open(courses_path) as f:
            courses = json.load(f)
        bad = []
        for c in courses:
            for t in c.get("targets", []):
                if t not in dag_ids:
                    bad.append(f"course {c.get('id')}: target {t} not in DAG")
        if bad:
            errors.append(f"course problems ({len(bad)}):\n  " + "\n  ".join(bad))

    # 8. Grading type enum + threshold sanity + enrichment wiring
    try:
        allowed_grading = {"multiple_choice", "numeric", "symbolic", "complex", "expression", "ordering", "tuple", "polynomial", "comparison"}
        thresh_warn = []
        grading_bad = []
        for name in os.listdir(CONCEPTS_DIR):
            if not name.endswith(".json") or name == "enrichment.json":
                continue
            with open(os.path.join(CONCEPTS_DIR, name)) as f:
                for c in json.load(f):
                    gt = c.get("grading_type", "")
                    if gt not in allowed_grading:
                        grading_bad.append(f"{c['id']}: grading_type {gt!r}")
                    mt = c.get("mastery_threshold", {})
                    st = mt.get("streak", 0)
                    avg = mt.get("avg_time_seconds", 0)
                    if not isinstance(st, int) or st <= 0:
                        thresh_warn.append(f"{c['id']}: streak {st}")
                    if not isinstance(avg, (int, float)) or avg <= 0:
                        thresh_warn.append(f"{c['id']}: avg_time_seconds {avg}")
                    elif avg > 120:
                        thresh_warn.append(f"{c['id']}: avg_time_seconds {avg} >120 (outlier)")
        if grading_bad:
            errors.append(f"grading_type problems ({len(grading_bad)}):\n  " + "\n  ".join(grading_bad))
        if thresh_warn:
            # treat >120 as warn but report; fail only if >300 or invalid
            hard = [x for x in thresh_warn if "avg_time_seconds 0" in x or "streak" in x or ">300" in x]
            if hard:
                errors.append(f"threshold problems ({len(hard)}):\n  " + "\n  ".join(hard))
            else:
                # log outlier but not fail — print to stderr for visibility
                print(f"note: threshold outliers ({len(thresh_warn)}): " + "; ".join(thresh_warn[:5]), file=sys.stderr)
    except Exception as e:
        errors.append(f"grading/threshold audit failed: {e}")

    # 9. Enrichment wiring check: members and ids must resolve in DAG
    enrich_path = os.path.join(CONCEPTS_DIR, "enrichment.json")
    if os.path.exists(enrich_path):
        try:
            with open(enrich_path) as f:
                enrich = json.load(f)
            enrich_bad = []
            for ent in enrich:
                if "id" in ent and ent["id"]:
                    if ent["id"] not in dag_ids:
                        enrich_bad.append(f"enrichment id {ent['id']} not in DAG")
                for mid in ent.get("members", []):
                    if mid not in dag_ids:
                        enrich_bad.append(f"enrichment member {mid} not in DAG")
            if enrich_bad:
                errors.append(f"enrichment problems ({len(enrich_bad)}):\n  " + "\n  ".join(enrich_bad))
        except Exception as e:
            errors.append(f"enrichment audit failed: {e}")

    if errors:
        print("\n\n".join(errors))
        print(f"\nFAIL: {sum(1 for _ in errors)} categories with issues")
        sys.exit(1)

    kp_summary = f", {kp_files} kp files, {kp_total} KPs" if 'kp_files' in locals() else ""
    print(f"OK: {len(dag_ids)} concepts mapped, {len(entries)} entries, "
          f"{len(sources)} sources{kp_summary}, kp shards resolved, 0 issues")


if __name__ == "__main__":
    main()
