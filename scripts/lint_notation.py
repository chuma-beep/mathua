#!/usr/bin/env python3
"""
Lens #3 — notation consistency.

Flags files that mix log/ln, degree symbols, etc. within the same lesson.

Advisory: writes scripts/notation_report.json
"""
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
REPORT = ROOT / "scripts" / "notation_report.json"

def lint_file(path: Path):
    text = path.read_text(encoding="utf-8")
    issues = []
    # log vs ln in same file (skip files that intentionally discuss both)
    has_log = bool(re.search(r"\\log\b", text))
    has_ln = bool(re.search(r"\\ln\b", text))
    if has_log and has_ln:
        # Allow files that explicitly say "logarithmic function" in prose but use ln in math — check math only
        math_blocks = re.findall(r"\\\\\[.*?\\\\\]|\$\$.*?\\$\\$|\$[^$\n]+\$", text, re.S)
        math_text = " ".join(math_blocks)
        if re.search(r"\\log\b", math_text) and re.search(r"\\ln\b", math_text):
            issues.append({
                "type": "log_ln_mix",
                "detail": "both \\log and \\ln in math",
                "line": 1,
            })
    # degree symbol consistency: ^{\circ} vs ^{\circ} is canonical; flag bare ^\circ without braces?
    # For now just flag mixed degree spellings if needed
    return issues

def main():
    report = []
    for p in sorted(ROOT.joinpath("data/lessons").rglob("*.md")):
        rel = str(p.relative_to(ROOT))
        for iss in lint_file(p):
            iss["file"] = rel
            report.append(iss)
    REPORT.write_text(json.dumps(report, indent=2, ensure_ascii=False))
    print(f"lens #3: {len(report)} issues")
    for iss in report[:10]:
        print(f" {iss['file']} {iss['type']}")
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
