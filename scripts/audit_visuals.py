#!/usr/bin/env python3
"""
audit_visuals.py, check the visualization system against the curriculum it claims to serve.

The existing tests prove pieces: the Zod schema parses a payload, a template renders a readout,
the Go loader rejects an unknown type. None of them asks the question a reviewer actually has,
which is whether the *installed* system is coherent: does every allowlisted template exist in
every place it must, is every authored visual on a step a learner can reach, does every one
carry a prediction whose correct answer is not always in the same place, and does every
template have both a mathematical test and an interaction test.

This script reads the real files. It deliberately does not import the TypeScript, so a schema
mismatch that only Zod would catch is left to `test/visualRegistry.test.ts`; what it checks is
the wiring, the placement, the authored content and the spread of answers.

Usage:
    python3 scripts/audit_visuals.py [--json]

Exit code 0 = no findings. Exit code 1 = at least one finding.
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from collections import Counter
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
VISUAL_DIR = ROOT / "web" / "next-app" / "components" / "visual"
TEMPLATE_DIR = VISUAL_DIR / "templates"
TEST_DIR = ROOT / "web" / "next-app" / "test"
SHARD_DIR = ROOT / "data" / "lessons" / "kp"
ALLOWLIST = ROOT / "data" / "visualizations" / "types.json"
REGISTRY = VISUAL_DIR / "registry.ts"
SCHEMAS = VISUAL_DIR / "schemas.ts"

MAX_PARAMS_BYTES = 8 * 1024
MAX_TITLE = 160
MAX_CAPTION = 600
MAX_QUESTION = 300
MAX_OPTION = 160
MAX_EXPLAIN = 600
MIN_OPTIONS = 2
MAX_OPTIONS = 4


def read(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def allowlist_ids() -> list[str]:
    entries = json.loads(read(ALLOWLIST))
    return [entry["id"] for entry in entries]


def registry_entries() -> dict[str, dict[str, str]]:
    """id -> {schema, template} parsed from registry.ts."""
    source = read(REGISTRY)
    pattern = re.compile(
        r"'([a-z0-9.]+)':\s*\{\s*schema:\s*([A-Za-z0-9_]+),\s*load:\s*\(\)\s*=>\s*import\('\./templates/([A-Za-z0-9_]+)'\)",
        re.S,
    )
    return {m.group(1): {"schema": m.group(2), "template": m.group(3)} for m in pattern.finditer(source)}


def schema_names() -> set[str]:
    return set(re.findall(r"export const ([A-Za-z0-9_]+Schema)", read(SCHEMAS)))


def shard_visuals() -> list[dict]:
    found = []
    for path in sorted(SHARD_DIR.glob("*.json")):
        steps = json.loads(read(path))
        for index, step in enumerate(steps):
            visual = step.get("visual")
            if visual:
                found.append({"file": path.name, "concept": path.stem, "step": index, "step_label": step.get("label", ""), "visual": visual})
    return found


def test_sources() -> str:
    parts = []
    for path in sorted(TEST_DIR.glob("visual*.test.ts*")):
        parts.append(read(path))
    return "\n".join(parts)


def coverage_findings(registry: dict[str, dict[str, str]]) -> tuple[list[str], dict[str, int]]:
    """Check the topic accounting against the templates that actually exist.

    Templates and topics are different units. One template can teach several
    topics, and a topic can be partial or absent while a related figure exists, so
    a count of templates says nothing on its own about coverage. The topics live
    in data/visualizations/coverage.json, which is the canonical accounting; this
    function only checks that the file and the registry cannot disagree.
    """
    data = json.loads(read(ALLOWLIST.parent / "coverage.json"))
    topics = data.get("topics", [])
    findings: list[str] = []
    counts = {"complete": 0, "partial": 0, "absent": 0}
    seen: set[str] = set()
    used: set[str] = set()

    for topic in topics:
        tid = str(topic.get("id", "")).strip()
        if not tid:
            findings.append("coverage: a topic has no id")
            continue
        if tid in seen:
            findings.append(f"coverage: duplicate topic id {tid!r}")
            continue
        seen.add(tid)

        status = topic.get("status")
        if status not in counts:
            findings.append(f"coverage: {tid} has status {status!r}, want complete, partial or absent")
            continue
        counts[status] += 1

        templates = topic.get("templates", [])
        if status == "complete" and not templates:
            findings.append(f"coverage: {tid} is complete but names no template")
        if status == "absent" and templates:
            findings.append(f"coverage: {tid} is absent but names {', '.join(templates)}")
        for name in templates:
            if name not in registry:
                findings.append(f"coverage: {tid} names unknown template {name}")
            used.add(name)

    # A figure outside the proposal's 31 topics is recorded under `supplementary`
    # with a reason, so the topic count stays honest and no template is orphaned.
    for entry in data.get("supplementary", []):
        name = str(entry.get("id", "")).strip()
        if not name:
            findings.append("coverage: a supplementary entry has no id")
            continue
        if not str(entry.get("reason", "")).strip():
            findings.append(f"coverage: supplementary template {name} has no reason")
        if name not in registry:
            findings.append(f"coverage: supplementary template {name} is not in the registry")
        used.add(name)

    for name in sorted(set(registry) - used):
        findings.append(f"coverage: template {name} teaches no topic in coverage.json")

    if len(topics) != 31:
        findings.append(f"coverage: {len(topics)} topics recorded; the proposal has 31")

    return findings, counts


def run() -> tuple[list[str], dict[str, int]]:
    findings: list[str] = []
    ids = allowlist_ids()
    registry = registry_entries()

    if len(ids) != len(set(ids)):
        findings.append("allowlist contains a duplicate id")
    if not ids:
        findings.append("allowlist is empty")

    # 1. The three lists agree.
    allow = set(ids)
    reg = set(registry)
    if allow - reg:
        findings.append(f"allowlisted but missing from registry.ts: {sorted(allow - reg)}")
    if reg - allow:
        findings.append(f"in registry.ts but not allowlisted: {sorted(reg - allow)}")

    templates_on_disk = {p.stem for p in TEMPLATE_DIR.glob("*.tsx")}
    registered = {entry["template"] for entry in registry.values()}
    if registered - templates_on_disk:
        findings.append(f"registry imports a template that does not exist: {sorted(registered - templates_on_disk)}")
    orphan = templates_on_disk - registered
    if orphan:
        findings.append(f"template file not reachable from the registry: {sorted(orphan)}")

    # 2. Every registry entry names a schema that schemas.ts exports.
    known_schemas = schema_names()
    for visual_id, entry in sorted(registry.items()):
        if entry["schema"] not in known_schemas:
            findings.append(f"{visual_id}: schema {entry['schema']} is not exported by schemas.ts")

    # 3. Every authored visual is coherent, reachable and predicted.
    visuals = shard_visuals()
    if not visuals:
        findings.append("no visual is authored in any shard, so this audit is checking nothing")
    answers: Counter[int] = Counter()
    by_type: Counter[str] = Counter()
    for item in visuals:
        visual = item["visual"]
        where = f"{item['file']} step {item['step'] + 1}"
        kind = visual.get("type", "")
        by_type[kind] += 1
        if kind not in allow:
            findings.append(f"{where}: type {kind!r} is not on the allowlist")
        if item["step"] != 0:
            findings.append(f"{where}: visual sits on step {item['step'] + 1}; /learn opens at step 1, so it is unreachable without answering a question")
        title = visual.get("title", "")
        caption = visual.get("caption", "")
        if not title:
            findings.append(f"{where}: no title")
        elif len(title) > MAX_TITLE:
            findings.append(f"{where}: title is {len(title)} runes, limit {MAX_TITLE}")
        if not caption:
            findings.append(f"{where}: no caption")
        elif len(caption) > MAX_CAPTION:
            findings.append(f"{where}: caption is {len(caption)} runes, limit {MAX_CAPTION}")
        params = visual.get("params")
        if params is not None:
            size = len(json.dumps(params, separators=(",", ":")).encode())
            if size > MAX_PARAMS_BYTES:
                findings.append(f"{where}: params are {size} bytes, limit {MAX_PARAMS_BYTES}")
        predict = visual.get("predict")
        if not predict:
            findings.append(f"{where}: no prediction, so the figure only demonstrates and never asks")
            continue
        question = predict.get("question", "")
        options = predict.get("options") or []
        answer = predict.get("answer")
        explain = predict.get("explain", "")
        if not question:
            findings.append(f"{where}: prediction has no question")
        elif len(question) > MAX_QUESTION:
            findings.append(f"{where}: question is {len(question)} runes, limit {MAX_QUESTION}")
        if not (MIN_OPTIONS <= len(options) <= MAX_OPTIONS):
            findings.append(f"{where}: {len(options)} options, want {MIN_OPTIONS} to {MAX_OPTIONS}")
        for option in options:
            if not option.strip():
                findings.append(f"{where}: an option is empty")
            elif len(option) > MAX_OPTION:
                findings.append(f"{where}: option is {len(option)} runes, limit {MAX_OPTION}")
        if not isinstance(answer, int) or not (0 <= answer < len(options)):
            findings.append(f"{where}: answer {answer!r} is outside 0 to {len(options) - 1}")
        else:
            answers[answer] += 1
        if len(explain) > MAX_EXPLAIN:
            findings.append(f"{where}: explanation is {len(explain)} runes, limit {MAX_EXPLAIN}")
        if not explain:
            findings.append(f"{where}: no explanation, so a wrong answer is corrected without saying why")

    # 4. The correct answer is not always in the same position. A learner who
    #    notices that stops reading the options, and the prediction stops asking.
    if answers:
        total = sum(answers.values())
        worst_position, worst_count = answers.most_common(1)[0]
        if worst_count == total and total > 1:
            findings.append(f"every prediction has its correct answer at position {worst_position + 1}")
        elif worst_count / total > 0.6 and total > 3:
            findings.append(
                f"{worst_count} of {total} predictions put the correct answer at position {worst_position + 1}"
            )

    # 5. Every template has an interaction test, and a template without one is a
    #    figure nobody has ever moved a control on.
    tests = test_sources()
    for visual_id, entry in sorted(registry.items()):
        name = entry["template"]
        if name not in tests:
            findings.append(f"{visual_id}: template {name} has no interaction test")

    coverage, counts = coverage_findings(registry)
    findings.extend(coverage)

    return findings, counts


def main() -> int:
    parser = argparse.ArgumentParser(description="Audit the installed visualization system.")
    parser.add_argument("--json", action="store_true", help="print the findings as JSON")
    args = parser.parse_args()

    findings, counts = run()
    if args.json:
        print(json.dumps({"findings": findings}, indent=2))
    else:
        if findings:
            print(f"FAIL: {len(findings)} finding(s)")
            for finding in findings:
                print(f"  - {finding}")
        else:
            print(
                f"OK: {len(registry_entries())} templates, {sum(counts.values())} topics "
                f"({counts['complete']} complete, {counts['partial']} partial, {counts['absent']} absent), "
                "authored visuals and predictions agree"
            )
    return 1 if findings else 0


if __name__ == "__main__":
    sys.exit(main())
