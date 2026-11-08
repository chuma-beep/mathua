"""Check that the 2,000-concept expansion is additive.

The 1,500-concept corpus is frozen in `data/expansion-baseline.json`. A wave may add concepts, but
it may not remove an existing concept, change its domain, drop a prerequisite edge, or remap its
lesson. The three quality baselines may only shrink: a new exemption in
`section_alignment_baseline.json` or `scripts/audit_baseline.json`, or a lower floor in
`variety_baseline.json`, is a regression even when every test still passes. That is the failure mode
this check exists for: a baseline update that makes a new defect look accepted.

Usage:
    python3 scripts/check_additive.py            # verify the working tree against the baseline
    python3 scripts/check_additive.py --write    # freeze the current corpus (setup only)
    python3 scripts/check_additive.py --write --force
"""
from __future__ import annotations

import argparse
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
BASELINE = ROOT / "data" / "expansion-baseline.json"
CONCEPTS = ROOT / "data" / "concepts"
LESSONS = ROOT / "data" / "lessons" / "lessons.json"
SECTION_BASELINE = ROOT / "internal" / "generator" / "all" / "section_alignment_baseline.json"
AUDIT_BASELINE = ROOT / "scripts" / "audit_baseline.json"
VARIETY_BASELINE = ROOT / "internal" / "generator" / "all" / "variety_baseline.json"


def load_json(path: Path):
    return json.loads(path.read_text(encoding="utf-8"))


def read_concepts() -> dict[str, dict]:
    concepts: dict[str, dict] = {}
    for path in sorted(CONCEPTS.glob("*.json")):
        if path.name == "enrichment.json":
            continue
        for node in load_json(path):
            cid = node.get("id")
            if not cid:
                continue
            concepts[cid] = {
                "domain": node["domain"],
                "prerequisites": sorted(node.get("prerequisites") or []),
            }
    return concepts


def read_lessons() -> dict[str, str]:
    return {row["concept_id"]: row["source"] for row in load_json(LESSONS)}


def read_audit() -> dict[str, list[str]]:
    raw = load_json(AUDIT_BASELINE)
    return {
        key: sorted(value)
        for key, value in raw.items()
        if key != "_note" and isinstance(value, list)
    }


def freeze() -> dict:
    return {
        "version": 1,
        "frozen_at": "2026-11-08",
        "concept_count": len(read_concepts()),
        "concepts": read_concepts(),
        "lessons": read_lessons(),
        "section_alignment_baseline": sorted(load_json(SECTION_BASELINE)),
        "audit_baseline": read_audit(),
        "variety_floors": load_json(VARIETY_BASELINE),
    }


def check(base: dict) -> list[str]:
    failures: list[str] = []
    concepts = read_concepts()
    lessons = read_lessons()

    for cid, frozen in base["concepts"].items():
        current = concepts.get(cid)
        if current is None:
            failures.append(f"concept removed: {cid}")
            continue
        if current["domain"] != frozen["domain"]:
            failures.append(
                f"domain changed for {cid}: {frozen['domain']} -> {current['domain']}"
            )
        removed = sorted(set(frozen["prerequisites"]) - set(current["prerequisites"]))
        for prerequisite in removed:
            failures.append(f"prerequisite removed from {cid}: {prerequisite}")

    for cid, source in base["lessons"].items():
        if cid not in lessons:
            failures.append(f"lesson mapping removed: {cid}")
        elif lessons[cid] != source:
            failures.append(f"lesson source changed for {cid}: {source} -> {lessons[cid]}")

    section = set(load_json(SECTION_BASELINE))
    for cid in sorted(section - set(base["section_alignment_baseline"])):
        failures.append(f"new section-alignment exemption: {cid}")

    audit = read_audit()
    for category, entries in audit.items():
        frozen = set(base["audit_baseline"].get(category, []))
        for entry in sorted(set(entries) - frozen):
            failures.append(f"new audit-baseline entry in {category}: {entry}")

    variety = load_json(VARIETY_BASELINE)
    for cid, floor in base["variety_floors"].items():
        if cid not in variety:
            failures.append(f"variety floor removed: {cid}")
        elif variety[cid] < floor:
            failures.append(f"variety floor lowered for {cid}: {floor} -> {variety[cid]}")

    return failures


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--write", action="store_true", help="freeze the current corpus")
    parser.add_argument("--force", action="store_true", help="overwrite an existing baseline")
    args = parser.parse_args()

    if args.write:
        if BASELINE.exists() and not args.force:
            raise SystemExit(
                f"{BASELINE.relative_to(ROOT)} already exists; pass --force to overwrite it"
            )
        snapshot = freeze()
        BASELINE.write_text(json.dumps(snapshot, indent=2, sort_keys=True) + "\n", encoding="utf-8")
        print(
            f"wrote {BASELINE.relative_to(ROOT)}: "
            f"{snapshot['concept_count']} concepts, {len(snapshot['lessons'])} lessons, "
            f"{len(snapshot['variety_floors'])} variety floors"
        )
        return 0

    if not BASELINE.exists():
        raise SystemExit(
            f"{BASELINE.relative_to(ROOT)} is missing; run python3 scripts/check_additive.py --write"
        )
    base = load_json(BASELINE)
    failures = check(base)
    if failures:
        for failure in failures:
            print(f"ERROR {failure}")
        print(f"FAIL: {len(failures)} additive-expansion violation(s)")
        return 1
    print(
        f"OK: {len(read_concepts())} concepts, {len(read_lessons())} lessons and the three "
        f"quality baselines preserve the {base['concept_count']}-concept snapshot"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
