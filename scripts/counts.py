#!/usr/bin/env python3
"""Single source of truth for the corpus counts quoted in the docs.

Counts drifted every content wave (641/580/284/1923 all appeared in tracked
docs). This script computes them from `data/` and either rewrites them into
README.md (`--write`) or verifies the docs match (`--check`, used by CI).

Usage:
  python3 scripts/counts.py            # print the counts as JSON
  python3 scripts/counts.py --write    # rewrite README.md counts in place
  python3 scripts/counts.py --check    # exit 1 if README.md or the docs list
                                       # a superseded number
"""

from __future__ import annotations

import argparse
import glob
import json
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# Display order + labels for the README domain table.
DOMAIN_LABELS = [
    ("arithmetic", "Arithmetic"),
    ("fractions", "Fractions"),
    ("prealgebra", "Pre-Algebra"),
    ("algebra", "Algebra"),
    ("geometry", "Geometry"),
    ("trigonometry", "Trigonometry"),
    ("precalculus", "Precalculus"),
    ("calculus", "Calculus"),
    ("statistics", "Statistics"),
    ("linear_algebra", "Linear Algebra"),
    ("discrete_math", "Discrete Math"),
    ("complex_numbers", "Complex Numbers"),
    ("number_theory", "Number Theory"),
    ("differential_equations", "Differential Equations"),
    ("abstract_algebra", "Abstract Algebra"),
    ("machine_learning", "Machine Learning"),
    ("topology", "Topology"),
]

# Numbers that were correct at some past corpus size and must never reappear.
STALE = ["641", "580", "284", "1923"]
STALE_TEXT = ["21 courses", "25 Routes"]

DOCS_TO_SCAN = [
    "README.md",
    "CONTRIBUTING.md",
    "docs/architecture.md",
    "docs/system-design.md",
    "docs/system-design-slides.md",
    "docs/wasc-uc-ag-research.md",
    "web/next-app/app/docs/architecture/page.tsx",
    "web/next-app/app/how-it-works/sections.tsx",
]


def compute() -> dict:
    by_domain: dict[str, int] = {}
    total = 0
    for f in sorted(glob.glob(os.path.join(ROOT, "data/concepts/*.json"))):
        if os.path.basename(f) == "enrichment.json":
            continue
        for c in json.load(open(f)):
            if "id" in c:
                total += 1
                by_domain[c["domain"]] = by_domain.get(c["domain"], 0) + 1
    kp_files = glob.glob(os.path.join(ROOT, "data/lessons/kp/*.json"))
    kp_entries = 0
    subgoals = 0
    for f in kp_files:
        d = json.load(open(f))
        kp_entries += len(d)
        for kp in d:
            subgoals += len(kp.get("subgoals", []))
    courses = len(json.load(open(os.path.join(ROOT, "data/courses.json"))))
    diagrams = len(json.load(open(os.path.join(ROOT, "data/diagrams/meta.json"))))
    return {
        "concepts": total,
        "domains": len(by_domain),
        "by_domain": by_domain,
        "kp_files": len(kp_files),
        "kp_entries": kp_entries,
        "subgoals": subgoals,
        "courses": courses,
        "diagrams": diagrams,
    }


def render(readme: str, c: dict) -> str:
    # Badges.
    readme = re.sub(r"badge/concepts-\d+-", f"badge/concepts-{c['concepts']}-", readme)
    readme = re.sub(r"badge/domains-\d+-", f"badge/domains-{c['domains']}-", readme)
    # Mission sentence.
    readme = re.sub(
        r"dependency graph of \d+ atomic concepts and \d+ courses",
        f"dependency graph of {c['concepts']} atomic concepts and {c['courses']} courses",
        readme,
    )
    # Worked examples sentence.
    readme = re.sub(
        r"\d+ worked examples \(\d+ KP shard files × 3 subgoals each\)",
        f"{c['kp_entries']} worked examples ({c['kp_files']} KP shard files × 3 subgoals each)",
        readme,
    )
    readme = re.sub(r"\d+ dual-coded diagrams", f"{c['diagrams']} dual-coded diagrams", readme)
    # Test-comment KP arithmetic.
    readme = re.sub(
        r"\(\d+ files ×3 =\d+ KPs\)",
        f"({c['kp_files']} files ×3 ={c['kp_entries']} KPs)",
        readme,
    )
    # Courses section opener.
    readme = re.sub(r"^\d+ courses from 4th grade", f"{c['courses']} courses from 4th grade", readme, flags=re.M)
    # Concept-graph summary sentence.
    readme = re.sub(
        r"The graph contains \*\*\d+ concepts\*\* across \*\*\d+ domains\*\* \(\d+ KP shard files × 3 = \d+ KPs\)",
        f"The graph contains **{c['concepts']} concepts** across **{c['domains']} domains** "
        f"({c['kp_files']} KP shard files × 3 = {c['kp_entries']} KPs)",
        readme,
    )
    # Domain table.
    rows = [f"| Domain | Concepts |", f"|--------|----------|"]
    for key, label in DOMAIN_LABELS:
        rows.append(f"| {label} | {c['by_domain'].get(key, 0)} |")
    rows.append(f"| **Total** | **{c['concepts']}** |")
    table = "\n".join(rows)
    readme = re.sub(
        r"\| Domain \| Concepts \|\n\|[-| ]+\|\n(?:\| [^|\n]+ \| \d+ \|\n)+\| \*\*Total\*\* \| \*\*\d+\*\* \|",
        table,
        readme,
    )
    return readme


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--write", action="store_true")
    ap.add_argument("--check", action="store_true")
    args = ap.parse_args()

    c = compute()
    if not args.write and not args.check:
        print(json.dumps(c, indent=2, sort_keys=True))
        return 0

    readme_path = os.path.join(ROOT, "README.md")
    original = open(readme_path, encoding="utf-8").read()
    rendered = render(original, c)

    if args.write:
        if rendered != original:
            open(readme_path, "w", encoding="utf-8").write(rendered)
            print("README.md counts updated")
        else:
            print("README.md already up to date")
        return 0

    # --check
    problems = []
    if rendered != original:
        problems.append("README.md counts are stale (run: python3 scripts/counts.py --write)")
    for rel in DOCS_TO_SCAN:
        path = os.path.join(ROOT, rel)
        if not os.path.exists(path):
            continue
        text = open(path, encoding="utf-8").read()
        for token in STALE:
            for m in re.finditer(rf"\b{token}\b", text):
                line = text[: m.start()].count("\n") + 1
                snippet = text.splitlines()[line - 1].strip()[:100]
                problems.append(f"{rel}:{line}: superseded count {token!r}: {snippet}")
        for token in STALE_TEXT:
            if token in text:
                line = text[: text.index(token)].count("\n") + 1
                problems.append(f"{rel}:{line}: superseded text {token!r}")
    if problems:
        print("count drift found:", file=sys.stderr)
        for p in problems:
            print("  " + p, file=sys.stderr)
        return 1
    print(f"counts ok: {c['concepts']} concepts, {c['kp_files']} KP files, {c['courses']} courses")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
