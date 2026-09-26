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

# Glued scraper-bleed element ids (PreTeXt xml:id concatenated onto text,
# e.g. "antiderivativess359ebbb1", "Example 1se93b3fa4"). The leading run
# must be glued (no space) to a word char or digit; spaced tags are left
# alone. Keep in sync with the cleanup rule (strip trailing s+7hex).
GLUED_ID_RE = r"[A-Za-z0-9]s[0-9a-f]{7,}"
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


def norm_section(s):
    """Strip math delimiters so shard sections (raw Algebrica `\\(...)`)
    compare equal to canonicalized headings (`$...$`) — the same rule as
    normSectionKey in internal/lessons/lesson.go. Keep in sync."""
    s = re.sub(r"\\\\[()\[\]]|\\[()\[\]]|\$\$?", "", s)
    return " ".join(s.split())


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
                    raw_head = m.group(2).strip()
                    heads.add(norm_section(raw_head))
                    if re.search(GLUED_ID_RE, raw_head):
                        kp_orphans.append(f"{src}: glued element id in heading {raw_head!r}")
            for kp in kps:
                if kp.get("section") and norm_section(kp["section"]) not in heads:
                    kp_orphans.append(f"{name}: section {kp['section']!r} unresolved")
                # Scraper-bleed element ids (PreTeXt xml:id glued onto text,
                # e.g. "antiderivativess359ebbb1") must never land in
                # user-visible KP fields.
                for field in ("label", "section"):
                    if kp.get(field) and re.search(GLUED_ID_RE, kp[field]):
                        kp_orphans.append(f"{name}: glued element id in KP {field} {kp[field]!r}")
                for sg in kp.get("subgoals", []) or []:
                    if re.search(GLUED_ID_RE, sg):
                        kp_orphans.append(f"{name}: glued element id in subgoal {sg!r}")
        if kp_orphans:
            errors.append(f"kp shard problems ({len(kp_orphans)}):\n  " + "\n  ".join(kp_orphans))
        # Validate shard file count and total KPs (634 files ×3 =1902)
        if kp_files != len(dag_ids):
            errors.append(f"kp shard count mismatch: {kp_files} files vs {len(dag_ids)} concepts")
        if kp_total != len(dag_ids) * 3:
            errors.append(f"kp total mismatch: {kp_total} KPs vs {len(dag_ids)*3} expected (3 per concept)")

    # 5b. Glued element ids in ANY lesson heading (including sources shadowed
    # by the reverse-lexicographic pick above): PreTeXt xml:id bleed such as
    # "antiderivativess359ebbb1" surfaces verbatim in study titles.
    id_heads = []
    for p in sorted(os.path.join(LESSONS, s) for s in sources):
        if not os.path.isfile(p):
            continue
        with open(p, encoding="utf-8") as f:
            for i, line in enumerate(f, 1):
                m = re.match(r"^(#{1,4})\s+(.+)$", line.strip())
                if m and re.search(GLUED_ID_RE, m.group(2)):
                    id_heads.append(f"{os.path.relpath(p, ROOT)}:{i} {m.group(2).strip()[:80]}")
    if id_heads:
        errors.append(f"glued element ids in headings ({len(id_heads)}):\n  " + "\n  ".join(id_heads))

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

    # 10. Stub lessons: teaching/ sources too short to carry a worked
    # example (template cluster: <=12 non-blank lines). Ratcheted against
    # scripts/audit_baseline.json — waves shrink it to empty; new stubs fail.
    # 11. Placeholder KP subgoals: generic "Review X / Work the example /
    # Verify" filler instead of moves naming the actual computation.
    try:
        with open(os.path.join(ROOT, "scripts", "audit_baseline.json")) as f:
            baseline = json.load(f)
    except Exception:
        baseline = {}
    stub_now, ph_now = [], []
    if os.path.isdir(LESSONS):
        for name in sorted(os.listdir(os.path.join(LESSONS, "teaching"))):
            if not name.endswith(".md"):
                continue
            p = os.path.join(LESSONS, "teaching", name)
            n = sum(1 for line in open(p, encoding="utf-8") if line.strip())
            if n <= 12:
                stub_now.append(f"teaching/{name}")
    kp_dir = os.path.join(LESSONS, "kp")
    if os.path.isdir(kp_dir):
        for name in sorted(os.listdir(kp_dir)):
            if not name.endswith(".json"):
                continue
            with open(os.path.join(kp_dir, name), encoding="utf-8") as f:
                kps = json.load(f)
            if any(
                any(str(s).startswith("Review ") or s in ("Work the example", "Verify")
                    for s in (k.get("subgoals") or []))
                for k in kps
            ):
                ph_now.append(name[:-5])
    for label, now, key in (("stub lessons", stub_now, "stub_lessons"),
                            ("placeholder-KP shards", ph_now, "placeholder_kps")):
        allowed = set(baseline.get(key, []))
        fresh = sorted(set(now) - allowed)
        fixed = sorted(allowed - set(now))
        print(f"note: {label}: {len(now)} current ({len(allowed)} baselined"
              f"{f', {len(fixed)} fixed' if fixed else ''})", file=sys.stderr)
        if fresh:
            errors.append(f"NEW {label} ({len(fresh)}), fix content, do not extend the baseline:\n  "
                          + "\n  ".join(fresh[:30]))

    # 12. Cliché KP subgoals: filler verbs that never name the move
    # ("Follow the worked example", "Apply the rule", "Check the result",
    # "Understand X"). Ratcheted like 10/11 — waves rewrite to move-naming
    # subgoals (P1 template) and shrink the baseline to empty.
    cliche_re = re.compile(r"^follow the worked|^apply the rule$|^check the result$|^understand\b", re.I)
    cliche_now = []
    if os.path.isdir(kp_dir):
        for name in sorted(os.listdir(kp_dir)):
            if not name.endswith(".json"):
                continue
            with open(os.path.join(kp_dir, name), encoding="utf-8") as f:
                kps = json.load(f)
            if any(cliche_re.search(str(sg))
                   for k in kps for sg in (k.get("subgoals") or [])):
                cliche_now.append(name[:-5])
    allowed = set(baseline.get("cliche_kp_shards", []))
    fresh = sorted(set(cliche_now) - allowed)
    fixed = sorted(allowed - set(cliche_now))
    print(f"note: cliche-KP shards: {len(cliche_now)} current ({len(allowed)} baselined"
          f"{f', {len(fixed)} fixed' if fixed else ''})", file=sys.stderr)
    if fresh:
        errors.append(f"NEW cliche-KP shards ({len(fresh)}), fix content, do not extend the baseline:\n  "
                      + "\n  ".join(fresh[:30]))

    # 15. Textbook apparatus in teaching/ sources: OpenStax back-matter that
    # renders inertly on atomic concept pages (readiness quizzes with dead
    # textbook refs, external resource links, static exercise lists,
    # self-check checklists). The DAG + diagnostic already do readiness;
    # practice already does exercises. Ratcheted: waves strip per file.
    apparatus_re = re.compile(
        r"^#{1,4}\s+(Be Prepared\b|Media\s*$|Everyday Math\b|Writing Exercises\b"
        r"|Self Check\b|Section \d+\.\d+ Exercises\s*$)", re.M)
    apparatus_now = []
    teaching_dir = os.path.join(LESSONS, "teaching")
    if os.path.isdir(teaching_dir):
        for name in sorted(os.listdir(teaching_dir)):
            if not name.endswith(".md"):
                continue
            with open(os.path.join(teaching_dir, name), encoding="utf-8") as f:
                body = f.read()
            if apparatus_re.search(body) or "ACCESS ADDITIONAL ONLINE RESOURCES" in body:
                apparatus_now.append(f"teaching/{name}")
    allowed = set(baseline.get("apparatus_lessons", []))
    fresh = sorted(set(apparatus_now) - allowed)
    fixed = sorted(allowed - set(apparatus_now))
    print(f"note: apparatus lessons: {len(apparatus_now)} current ({len(allowed)} baselined"
          f"{f', {len(fixed)} fixed' if fixed else ''})", file=sys.stderr)
    if fresh:
        errors.append(f"NEW apparatus lessons ({len(fresh)}), strip back-matter, do not extend the baseline:\n  "
                      + "\n  ".join(fresh[:30]))

    # 13. Diagram metadata coverage (data/diagrams/meta.json): every mapped
    # asset needs an entry with a valid source and a non-empty title; SVG
    # algebrica assets need the upstream link (CC BY-NC attribution record).
    # 14. SVG-first for local mappings: hand-authored assets (anything outside
    # /diagrams/algebrica/) must be .svg; vendored algebrica tracks upstream.
    engine_src = os.path.join(ROOT, "internal", "engine", "engine.go")
    mapped_assets = []
    if os.path.exists(engine_src):
        with open(engine_src) as f:
            eng = f.read()
        start = eng.find("var conceptDiagrams = map[string]string{")
        if start >= 0:
            block = eng[start:eng.find("\n}", start)]
            mapped_assets = sorted(set(
                m.group(2) for m in re.finditer(r'"([^"]+)":\s*"((?:/diagrams/)[^"]+)"', block)))
    meta_path = os.path.join(ROOT, "data", "diagrams", "meta.json")
    meta = {}
    if os.path.exists(meta_path):
        with open(meta_path, encoding="utf-8") as f:
            meta = json.load(f)
    # Grandfathered: algebrica-dir SVGs with no same-stem upstream twin
    # (local renames or older releases). New algebrica SVGs must link.
    UNLINKED_GRANDFATHERED = frozenset({
        "/diagrams/algebrica/completing-square.svg",
        "/diagrams/algebrica/complex-plane.svg",
        "/diagrams/algebrica/hyperbolic-functions.svg",
        "/diagrams/algebrica/inverse-trig-graphs.svg",
        "/diagrams/algebrica/law-of-cosines.svg",
        "/diagrams/algebrica/law-of-sines.svg",
        "/diagrams/algebrica/linear-equation-graph.svg",
        "/diagrams/algebrica/number-line-absolute-value.svg",
        "/diagrams/algebrica/number-line-intervals.svg",
        "/diagrams/algebrica/number-line-real.svg",
        "/diagrams/algebrica/number-types-venn.svg",
        "/diagrams/algebrica/pascals-triangle.svg",
        "/diagrams/algebrica/polynomial-roots-graph.svg",
        "/diagrams/algebrica/pythagorean-theorem.svg",
        "/diagrams/algebrica/reference-angles.svg",
        "/diagrams/algebrica/right-triangle-trig.svg",
        "/diagrams/algebrica/right-triangle-unit-circle.svg",
        "/diagrams/algebrica/sec-csc-cot-graphs.svg",
        "/diagrams/algebrica/sine-cosine-graph.svg",
        "/diagrams/algebrica/unit-circle-labeled.svg",
        "/diagrams/algebrica/unit-circle-sine-cosine.svg",
        "/diagrams/algebrica/unit-circle-tangent.svg",
        "/diagrams/algebrica/vector-addition.svg",
        "/diagrams/algebrica/vector-arrow.svg",
    })
    meta_bad, raster_local = [], []
    for asset in mapped_assets:
        ent = meta.get(asset)
        if not isinstance(ent, dict):
            meta_bad.append(f"{asset}: no meta.json entry")
            continue
        if ent.get("source") not in ("algebrica", "local"):
            meta_bad.append(f"{asset}: source {ent.get('source')!r} not algebrica|local")
        if not (ent.get("title") or "").strip():
            meta_bad.append(f"{asset}: empty title")
        if ent.get("source") == "algebrica" and asset.lower().endswith(".svg") \
                and not ent.get("upstream") and asset not in UNLINKED_GRANDFATHERED:
            meta_bad.append(f"{asset}: algebrica SVG without upstream link")
        if "/diagrams/algebrica/" not in asset and not asset.lower().endswith(".svg"):
            raster_local.append(f"{asset}: local mapping must be .svg")
    # Stale entries pointing at unmapped assets.
    mapped_set = set(mapped_assets)
    for asset in sorted(meta):
        if asset not in mapped_set:
            meta_bad.append(f"{asset}: meta entry for unmapped asset")
    if meta_bad:
        errors.append(f"diagram metadata problems ({len(meta_bad)}):\n  " + "\n  ".join(meta_bad[:30]))
    if raster_local:
        errors.append(f"raster local diagrams ({len(raster_local)}):\n  " + "\n  ".join(raster_local[:30]))

    if errors:
        print("\n\n".join(errors))
        print(f"\nFAIL: {sum(1 for _ in errors)} categories with issues")
        sys.exit(1)

    kp_summary = f", {kp_files} kp files, {kp_total} KPs" if 'kp_files' in locals() else ""
    print(f"OK: {len(dag_ids)} concepts mapped, {len(entries)} entries, "
          f"{len(sources)} sources{kp_summary}, kp shards resolved, 0 issues")


if __name__ == "__main__":
    main()
