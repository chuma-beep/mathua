"""Validate a wave manifest and emit draft lesson and KP files.

The manifest is the authoring unit for the 2,000-concept expansion. It is deliberately strict:
every concept needs an id, a label, three section headings, three KP shards aligned to those
headings, a grading type and prerequisites that already exist or appear earlier in the same wave.
The tool never edits the corpus. `--out DIR` writes reviewable drafts outside it; installing them
is a separate, reviewed step.

Usage:
    python3 scripts/scaffold-wave.py --manifest docs/expansion-waves/wave-1.json
    python3 scripts/scaffold-wave.py --manifest wave.json --out /tmp/wave-1

The validator reads the target table in `docs/content-expansion.md` so a wave cannot promise more
concepts in a domain than its target allows.
"""
from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
CONCEPTS = ROOT / "data" / "concepts"
EXPANSION_DOC = ROOT / "docs" / "content-expansion.md"

GRADING_TYPES = {
    "multiple_choice",
    "numeric",
    "symbolic",
    "complex",
    "expression",
    "ordering",
    "tuple",
    "polynomial",
    "comparison",
}

ID_RE = re.compile(r"^[a-z][a-z0-9_]*(\.[a-z0-9_]+)+$")
TARGET_ROW = re.compile(
    r"^\|\s*(\w+)\s*\|\s*\w+\s*\|\s*(\d+)\s*\|\s*(\d+)\s*\|\s*\+\d+\s*\|\s*$"
)


def load_json(path: Path):
    return json.loads(path.read_text(encoding="utf-8"))


def corpus() -> tuple[dict[str, dict], dict[str, list[str]]]:
    by_id: dict[str, dict] = {}
    labels: dict[str, list[str]] = {}
    for path in sorted(CONCEPTS.glob("*.json")):
        if path.name == "enrichment.json":
            continue
        for node in load_json(path):
            if "id" not in node:
                continue
            by_id[node["id"]] = node
            key = node["label"].strip().casefold()
            labels.setdefault(key, []).append(node["id"])
    return by_id, labels


def targets() -> dict[str, tuple[int, int]]:
    text = EXPANSION_DOC.read_text(encoding="utf-8")
    start = text.find("| Domain | Tier | Now | Target | Add |")
    if start < 0:
        raise SystemExit("docs/content-expansion.md: the section 5.5 target table is gone")
    found: dict[str, tuple[int, int]] = {}
    for line in text[start:].split("\n")[2:]:
        match = TARGET_ROW.match(line)
        if not match:
            break
        found[match.group(1)] = (int(match.group(2)), int(match.group(3)))
    if not found:
        raise SystemExit("docs/content-expansion.md: the target table has no rows")
    return found


def require(condition: bool, failures: list[str], message: str) -> None:
    if not condition:
        failures.append(message)


def validate(manifest: dict, by_id: dict[str, dict], labels: dict[str, list[str]]) -> list[str]:
    failures: list[str] = []
    require(isinstance(manifest.get("wave"), int) and manifest["wave"] > 0, failures,
            "wave must be a positive integer")
    concepts = manifest.get("concepts")
    if not isinstance(concepts, list) or not concepts:
        return ["concepts must be a non-empty list"]

    adds = manifest.get("domain_adds")
    if not isinstance(adds, dict) or not adds:
        failures.append("domain_adds must be a non-empty object")
        adds = {}

    available = targets()
    counts: dict[str, int] = {}
    seen_ids: dict[str, int] = {}
    seen_labels: dict[str, int] = {}
    manifest_ids: set[str] = set()

    for index, concept in enumerate(concepts):
        where = f"concepts[{index}]"
        cid = concept.get("id")
        if not isinstance(cid, str) or not ID_RE.match(cid):
            failures.append(f"{where}.id is not a dotted lowercase id")
            continue
        if cid in by_id:
            failures.append(f"{cid}: id already exists in the corpus")
        if cid in seen_ids:
            failures.append(f"{cid}: duplicate id in the manifest")
        seen_ids[cid] = index
        manifest_ids.add(cid)

        domain = concept.get("domain")
        if domain not in available:
            failures.append(f"{cid}: unknown domain {domain!r}")
        else:
            counts[domain] = counts.get(domain, 0) + 1

        subdomain = concept.get("subdomain")
        if not isinstance(subdomain, str) or not subdomain.startswith(f"{domain}."):
            failures.append(f"{cid}: subdomain must start with {domain!r}.")

        label = concept.get("label")
        if not isinstance(label, str) or not label.strip():
            failures.append(f"{cid}: label is empty")
        else:
            key = label.strip().casefold()
            if key in labels:
                failures.append(f"{cid}: label collides with {labels[key][0]}: {label}")
            if key in seen_labels:
                failures.append(f"{cid}: duplicate label in the manifest: {label}")
            seen_labels[key] = index

        if concept.get("grading_type") not in GRADING_TYPES:
            failures.append(f"{cid}: grading_type {concept.get('grading_type')!r} is not known")

        threshold = concept.get("mastery_threshold")
        if not isinstance(threshold, dict):
            failures.append(f"{cid}: mastery_threshold is missing")
        else:
            require(isinstance(threshold.get("streak"), int) and threshold["streak"] > 0,
                    failures, f"{cid}: mastery_threshold.streak must be a positive integer")
            require(isinstance(threshold.get("avg_time_seconds"), (int, float))
                    and threshold["avg_time_seconds"] > 0,
                    failures, f"{cid}: mastery_threshold.avg_time_seconds must be positive")

        prerequisites = concept.get("prerequisites")
        if not isinstance(prerequisites, list):
            failures.append(f"{cid}: prerequisites must be a list")
            prerequisites = []
        for prerequisite in prerequisites:
            if prerequisite == cid:
                failures.append(f"{cid}: a concept cannot require itself")
            elif prerequisite not in by_id and prerequisite not in manifest_ids:
                failures.append(f"{cid}: prerequisite {prerequisite} does not exist")
            elif prerequisite in seen_ids and seen_ids[prerequisite] >= index:
                failures.append(
                    f"{cid}: prerequisite {prerequisite} must appear earlier in the wave"
                )

        sections = concept.get("sections")
        if not isinstance(sections, list) or len(sections) != 3:
            failures.append(f"{cid}: sections must hold exactly three headings")
            sections = []
        elif len({section.strip().casefold() for section in sections}) != 3:
            failures.append(f"{cid}: the three section headings must be distinct")

        kps = concept.get("kps")
        if not isinstance(kps, list) or len(kps) != 3:
            failures.append(f"{cid}: kps must hold exactly three shards")
            kps = []
        kp_sections: list[str] = []
        for kp_index, kp in enumerate(kps):
            if not isinstance(kp, dict):
                failures.append(f"{cid}: kps[{kp_index}] is not an object")
                continue
            if not isinstance(kp.get("label"), str) or not kp["label"].strip():
                failures.append(f"{cid}: kps[{kp_index}].label is empty")
            section = kp.get("section")
            if section not in sections:
                failures.append(f"{cid}: kps[{kp_index}].section is not one of the three headings")
            else:
                kp_sections.append(section)
            subgoals = kp.get("subgoals")
            if not isinstance(subgoals, list) or len(subgoals) != 3:
                failures.append(f"{cid}: kps[{kp_index}] must hold exactly three subgoals")
            elif len({goal.strip().casefold() for goal in subgoals if isinstance(goal, str)}) != 3:
                failures.append(f"{cid}: kps[{kp_index}] subgoals must be three distinct strings")
        if sections and sorted(kp_sections) != sorted(sections):
            failures.append(f"{cid}: each section needs exactly one KP shard")

        source = concept.get("lesson_source")
        if not isinstance(source, str) or not source.endswith(".md"):
            failures.append(f"{cid}: lesson_source must be a .md path")

    for domain, count in sorted(adds.items()):
        if domain not in available:
            failures.append(f"domain_adds names unknown domain {domain}")
            continue
        now, target = available[domain]
        if count != counts.get(domain, 0):
            failures.append(
                f"domain_adds says {count} for {domain} but the manifest holds "
                f"{counts.get(domain, 0)}"
            )
        if count > target - now:
            failures.append(
                f"domain_adds promises {count} for {domain}, above the remaining "
                f"{target - now} before the target"
            )
    for domain, count in sorted(counts.items()):
        if domain not in adds:
            failures.append(f"domain_adds is missing {domain} ({count} concepts)")

    return failures


def slug(cid: str) -> str:
    return cid.replace(".", "-")


def emit(manifest: dict, out: Path) -> None:
    lessons = out / "lessons"
    shards = out / "kp"
    lessons.mkdir(parents=True, exist_ok=True)
    shards.mkdir(parents=True, exist_ok=True)
    for concept in manifest["concepts"]:
        cid = concept["id"]
        lesson = concept.get("lesson")
        if not isinstance(lesson, dict):
            raise SystemExit(f"{cid}: --out needs a lesson object with three section bodies")
        bodies = lesson.get("sections")
        if not isinstance(bodies, list) or len(bodies) != 3:
            raise SystemExit(f"{cid}: lesson.sections must hold exactly three bodies")
        lines = [f"# {concept['label']}", ""]
        for heading, body in zip(concept["sections"], bodies):
            if not isinstance(body, str) or not body.strip():
                raise SystemExit(f"{cid}: lesson body for {heading!r} is empty")
            lines.extend([f"## {heading}", "", body.strip(), ""])
        (lessons / f"{slug(cid)}.md").write_text("\n".join(lines), encoding="utf-8")
        (shards / f"{cid}.json").write_text(
            json.dumps(concept["kps"], indent=2) + "\n", encoding="utf-8"
        )
    (out / "manifest.normalized.json").write_text(
        json.dumps(manifest, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
    print(f"wrote drafts for {len(manifest['concepts'])} concepts to {out}")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--manifest", required=True, help="wave manifest JSON")
    parser.add_argument("--out", help="write draft lesson and KP files here")
    args = parser.parse_args()

    manifest = load_json(Path(args.manifest))
    by_id, labels = corpus()
    failures = validate(manifest, by_id, labels)
    if failures:
        for failure in failures:
            print(f"ERROR {failure}")
        print(f"FAIL: {len(failures)} manifest problem(s)")
        return 1
    counts: dict[str, int] = {}
    for concept in manifest["concepts"]:
        counts[concept["domain"]] = counts.get(concept["domain"], 0) + 1
    summary = ", ".join(f"{domain} {count}" for domain, count in sorted(counts.items()))
    print(f"OK: wave {manifest['wave']} holds {len(manifest['concepts'])} concepts ({summary})")
    if args.out:
        emit(manifest, Path(args.out))
    return 0


if __name__ == "__main__":
    sys.exit(main())
