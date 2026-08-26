#!/usr/bin/env python3
"""
Lens #1 — hypothesis completeness.

Flags theorem statements where a conclusion requires hypotheses not mentioned
in the same block. Phase 1 covers the "monotone → finite limit" class; the
framework generalizes to other patterns.

Advisory: writes scripts/hypothesis_report.json
"""
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
REPORT = ROOT / "scripts" / "hypothesis_report.json"

# Patterns: each entry is (theorem cue regex, required keyword regex, message)
PATTERNS = [
    (
        re.compile(r"is monotonic.*?admits a.*?(?:finite )?limit", re.I | re.S),
        re.compile(r"bounded", re.I),
        "monotone → finite limit without 'bounded'",
    ),
    (
        re.compile(r"continuous on.*\[a,\s*b\].*?differentiable on.*\(a,\s*b\)", re.I | re.S),
        None,  # no required extra — just smoke test that both appear together (Rolle/MVT shape)
        None,
    ),
]

def lint_file(path: Path):
    text = path.read_text(encoding="utf-8")
    issues = []
    # Split into paragraphs (blank-line separated) — theorem statements live in one paragraph
    paras = re.split(r"\n\s*\n", text)
    for para in paras:
        # Lens 1a: monotone without bounded before finite limit claim
        if re.search(r"monotonic", para, re.I) and re.search(r"admits a.*limit", para, re.I):
            if "finite" in para.lower() and "bounded" not in para.lower():
                # Only flag if this is the statement paragraph itself, not the follow-up
                # that correctly says "bounded monotonic"
                if "If a sequence" in para or "If" in para[:20]:
                    issues.append({
                        "type": "monotone_without_bounded",
                        "excerpt": para.strip()[:200].replace("\n", " "),
                        "line": text.count("\n", 0, text.index(para)) + 1,
                    })
    return issues

def main():
    report = []
    for p in sorted(ROOT.joinpath("data/lessons").rglob("*.md")):
        rel = str(p.relative_to(ROOT))
        for iss in lint_file(p):
            iss["file"] = rel
            report.append(iss)
    REPORT.write_text(json.dumps(report, indent=2, ensure_ascii=False))
    print(f"lens #1: {len(report)} issues")
    for iss in report[:10]:
        print(f" {iss['file']}:{iss['line']} {iss['type']}")
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
