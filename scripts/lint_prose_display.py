#!/usr/bin/env python3
"""
Lens #2 — prose vs display agreement.

Finds places where prose names functions (f, g, h) and the next display
uses a different set. Catches the squeeze f↔g swap class.

Advisory: writes scripts/prose_display_report.json
"""
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
REPORT = ROOT / "scripts" / "prose_display_report.json"

# Prose patterns that name functions
PROSE_LIMITS_RE = re.compile(
    r"limits of\s+\\\(\s*([a-zA-Z])\(x\)\s*\\\)\s*and\s+\\\(\s*([a-zA-Z])\(x\)\s*\\\)",
    re.I,
)
PROSE_ADMITS_RE = re.compile(
    r"function\s+\\\(\s*([a-zA-Z])\(x\)\s*\\\)\s+also admits",
    re.I,
)
DISPLAY_RE = re.compile(r"\\\\\[([\s\S]*?)\\\\\]|(?<!\\)\$\$([\s\S]*?)\$\$")
LIM_RE = re.compile(r"\\lim[^}]*?([a-zA-Z])\s*\(")

def extract_displays(text: str):
    for m in DISPLAY_RE.finditer(text):
        body = (m.group(1) or m.group(2) or "").strip()
        yield m.start(), body

def lint_file(path: Path):
    text = path.read_text(encoding="utf-8")
    displays = list(extract_displays(text))
    issues = []
    # Check 1: "limits of f and h" vs next display's lim letters
    for m in PROSE_LIMITS_RE.finditer(text):
        prose_letters = sorted({m.group(1).lower(), m.group(2).lower()})
        # next display after this prose
        nxt = next(((pos, body) for pos, body in displays if pos > m.start()), None)
        if not nxt:
            continue
        _, body = nxt
        lim_letters = sorted(set(LIM_RE.findall(body)))
        lim_letters = [l.lower() for l in lim_letters]
        if lim_letters and prose_letters != sorted(lim_letters):
            issues.append({
                "type": "limits_mismatch",
                "line": text.count("\n", 0, m.start()) + 1,
                "prose": m.group(0).strip()[:120],
                "prose_letters": prose_letters,
                "display_letters": lim_letters,
                "display_excerpt": body[:120].replace("\n", " "),
            })
    # Check 2: "function g(x) also admits" vs next display's lim letter
    for m in PROSE_ADMITS_RE.finditer(text):
        prose_letter = m.group(1).lower()
        nxt = next(((pos, body) for pos, body in displays if pos > m.start()), None)
        if not nxt:
            continue
        _, body = nxt
        lim_letters = [l.lower() for l in LIM_RE.findall(body)]
        if lim_letters and prose_letter not in lim_letters:
            # The display should conclude the same function the prose names
            # e.g. prose says g admits, display says lim f = ...
            issues.append({
                "type": "conclusion_mismatch",
                "line": text.count("\n", 0, m.start()) + 1,
                "prose": m.group(0).strip()[:120],
                "prose_letter": prose_letter,
                "display_letters": lim_letters,
                "display_excerpt": body[:120].replace("\n", " "),
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
    print(f"lens #2: {len(report)} issues")
    for iss in report[:10]:
        print(f" {iss['file']}:{iss['line']} {iss['type']} prose={iss.get('prose_letters') or iss.get('prose_letter')} display={iss['display_letters']}")
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
