#!/usr/bin/env python3
"""
Lens #4 — math layout: multi-line display blocks must break lines explicitly.

A bare newline inside TeX math is whitespace, not a row break. Display blocks
(\\[...\\] / $$...$$) that stack relation steps on separate source lines
without `\\\\` or an aligned-family env collapse into one long line at render
time (horizontal overflow wall on mobile, `=` visually lost). This lens flags:

  bare-chain     display block with a lone `=`/`+`/`\\+`/`\\-` line, no `\\\\`,
                 no aligned/array env                         (error)
  multiline-eq   display block whose non-first content lines start with (or
                 non-last lines end with) a bare `=` relation, but no `\\\\`
                 and no aligned/array env                     (error)
  fused-display  `$$` match spanning prose/single-`$` fragments or blank
                 lines — a mispaired match owned by the frontend
                 fused-inline defense, never rewritten here    (info)
  bad-plus       `\\+` inside math (canonical form is `+`; Go mathNorm
                 currently papers over it at serve time)       (error)
  detached-sub   `]\\_` — escaped underscore after a bracket detaches the
                 subscript from `\\right]` (valid KaTeX, wrong output) (error)
  eqref-math     `the equation \\(1\\)` style \\tag pointers written as math
                 instead of text `(1)`                        (error)
                 ("identity \\(1\\)" is the identity ELEMENT — math, ignored)

Usage:
  python3 scripts/lint_math_layout.py            # report, exit 1 on errors
  python3 scripts/lint_math_layout.py --fix      # mechanical fixes in place

Report goes to scripts/math_layout_report.json (list; empty means clean).
--fix handles: lone `\\+` lines -> `+`, `\\+`/`\\-`-safe rewrites in math,
`]_` detachment, eqref-math -> text, and wraps bare-chain / multiline-eq
blocks in `aligned` (first relation `&=`, continuations `\\\\ &=` /
`\\\\ &+` / `\\\\ &-`). Blocks already containing `\\\\`, an aligned/array
env, single-`$` fragments, or blank lines are never rewritten.
"""
from __future__ import annotations

import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
REPORT = ROOT / "scripts" / "math_layout_report.json"

# `\\[0.5em]` is row-break SPACING, not a display opener — excluded so the
# match cannot start mid-array and swallow prose to the next real `\\]`.
DISPLAY_RE = re.compile(r"\\\\\[(?!\s*\d)([\s\S]*?)\\\\\]|(?<!\\)\$\$([\s\S]*?)\$\$")
INLINE_PAREN_RE = re.compile(r"\\\\\(([\s\S]*?)\\\\\)")
INLINE_DOLLAR_RE = re.compile(r"(?<!\\)\$(?!\$)([^$\n]+?)(?<!\\)\$(?!\$)")

ALIGNED_ENV_RE = re.compile(r"\\begin\{(aligned|align\*?|gather\*?|cases|array|matrix|pmatrix|bmatrix|vmatrix)\}")
# Row break: 2+ backslashes followed by whitespace/EOL (optionally an
# alignment tab), or KaTeX row-break spacing (\\[6pt]). Deliberately NOT
# matched: \\[ \\] \\( \\) delimiters, \\, \\! \\; spacing, \\{ \\} braces.
ROWBREAK_RE = re.compile(r"\\{2,}(?=\s|$)|\\{2,}\[\s*\d")
LONE_REL_RE = re.compile(r"^\s*(=|\+|\\\+|\\=|\\-)\s*$")
SINGLE_DOLLAR_RE = re.compile(r"(?<!\\)\$(?!\$)")
DETACHED_SUB_RE = re.compile(r"\]\\_")
BAD_PLUS_RE = re.compile(r"(?<!\\)\\\+")
# Only equation/formula/integral references are \\tag pointers; "identity (1)"
# is the identity ELEMENT written (correctly) as math — never rewrite those.
EQREF_RE = re.compile(
    r"(equation|formula|integral)\s+\\\\\((\d{1,2})\\\\\)",
    re.IGNORECASE,
)
# A continuation line starts with a bare relation. Leading `-` is excluded:
# `- \\frac...` may be a fresh (negated) expression, not a join.
LEAD_REL_RE = re.compile(r"^(=|\\approx|\\equiv|\\simeq|\\cong|\\le|\\ge)\s+(.+)$")
# A line ending in a bare `=` (not ==, <=, >=, !=, =>, \:=) joins with the next.
TRAIL_EQ_RE = re.compile(r"^(?P<head>.+?)(?<![<>=!:])=\s*$")

ERROR_TYPES = {"bare-chain", "multiline-eq", "bad-plus", "detached-sub", "eqref-math"}


def _count_joins(content: list[str]) -> tuple[int, bool]:
    """Count relation joins across content lines.

    A single join (LHS newline REL newline RHS) renders as one correct
    display line — only 2+ joins mean stacked steps that collapse without
    `\\\\`. Returns (joins, has_lone_line).
    """
    joins = 0
    lone = False
    for i, l in enumerate(content):
        if LONE_REL_RE.match(l):
            joins += 1
            lone = True
        elif i > 0 and LEAD_REL_RE.match(l):
            joins += 1
        elif i < len(content) - 1 and TRAIL_EQ_RE.match(l):
            joins += 1
    return joins, lone


def classify_display(body: str):
    """Return 'bare-chain' | 'multiline-eq' | 'fused-display' | None."""
    if "\n\n" in body or SINGLE_DOLLAR_RE.search(body):
        return "fused-display"
    if ROWBREAK_RE.search(body) or ALIGNED_ENV_RE.search(body):
        return None
    content = [l.strip() for l in body.split("\n") if l.strip()]
    joins, lone = _count_joins(content)
    if joins >= 2:
        return "bare-chain" if lone else "multiline-eq"
    return None


def iter_display_blocks(text: str):
    for m in DISPLAY_RE.finditer(text):
        body = m.group(1) if m.group(1) is not None else m.group(2)
        yield m.start(), m.end(), body or ""


def iter_inline_bodies(text: str):
    for m in INLINE_PAREN_RE.finditer(text):
        yield m.group(1) or ""
    for m in INLINE_DOLLAR_RE.finditer(text):
        yield m.group(1) or ""


def lint_text(rel: str, text: str):
    issues = []

    def add(kind: str, line: int, excerpt: str):
        issues.append({"file": rel, "type": kind, "line": line,
                       "excerpt": excerpt[:140].replace("\n", " ")})

    for start, _end, body in iter_display_blocks(text):
        base_line = text.count("\n", 0, start) + 1
        kind = classify_display(body)
        if kind in ("bare-chain", "multiline-eq"):
            add(kind, base_line, body.strip())
        elif kind == "fused-display":
            add(kind, base_line, body.strip())
        if DETACHED_SUB_RE.search(body):
            add("detached-sub", base_line, body.strip())
        for l in body.split("\n"):
            if BAD_PLUS_RE.search(l):
                add("bad-plus", base_line, l.strip())
                break

    for m in INLINE_PAREN_RE.finditer(text):
        inner = m.group(1) or ""
        if DETACHED_SUB_RE.search(inner):
            add("detached-sub", text.count("\n", 0, m.start()) + 1, inner.strip())
        if re.search(r"(?<!\\)\\\+", inner):
            add("bad-plus", text.count("\n", 0, m.start()) + 1, inner.strip())

    for m in EQREF_RE.finditer(text):
        add("eqref-math", text.count("\n", 0, m.start()) + 1, m.group(0).strip())

    return issues


def fix_text(text: str):
    """Apply mechanical fixes. Returns (new_text, counts dict)."""
    counts = {"bad-plus": 0, "detached-sub": 0, "eqref-math": 0, "bare-chain": 0, "multiline-eq": 0}

    # Lone `\+` lines -> `+` (canonical; safe: a lone line is never content).
    def _plus(m):
        counts["bad-plus"] += 1
        return m.group(1) + "+" + m.group(2)

    text = re.sub(r"(?m)^([ \t]*)\\\+([ \t]*)$", _plus, text)

    # `]\_` -> `]_` and `\+` -> `+` inside math display/inline bodies only.
    # `\+` is render-neutral: Go mathNorm strips it at serve time already.
    def _fix_body(body: str):
        new, n = DETACHED_SUB_RE.subn("]_", body)
        counts["detached-sub"] += n
        new, n = re.subn(r"(?<!\\)\\\+", "+", new)
        counts["bad-plus"] += n
        return new

    parts = []
    pos = 0
    for m in DISPLAY_RE.finditer(text):
        parts.append(text[pos:m.start()])
        body = m.group(1) if m.group(1) is not None else m.group(2)
        fixed = _fix_body(body)
        if m.group(1) is not None:
            parts.append("\\\\[" + fixed + "\\\\]")
        else:
            parts.append("$$" + fixed + "$$")
        pos = m.end()
    parts.append(text[pos:])
    text = "".join(parts)

    # Same typo fixes inside \(...\) inline bodies.
    def _fix_paren(m):
        inner = _fix_body(m.group(1) or "")
        return "\\\\(" + inner + "\\\\)"

    text = INLINE_PAREN_RE.sub(_fix_paren, text)

    # Equation references -> text.
    def _eqref(m):
        counts["eqref-math"] += 1
        return f"{m.group(1)} ({m.group(2)})"

    text = EQREF_RE.sub(_eqref, text)

    # Wrap bare-chain blocks in aligned.
    text = _fix_bare_chains(text, counts)
    return text, counts


def _fix_bare_chains(text: str, counts: dict) -> str:
    out = []
    pos = 0
    for m in DISPLAY_RE.finditer(text):
        is_bracket = m.group(1) is not None
        body = m.group(1) if is_bracket else m.group(2)
        body = body or ""
        kind = classify_display(body)
        if kind in ("bare-chain", "multiline-eq"):
            out.append(text[pos:m.start()])
            out.append((_open(is_bracket) + _aligned_lines(body) + _close(is_bracket)))
            counts[kind] += 1
            pos = m.end()
        else:
            out.append(text[pos:m.end()])
            pos = m.end()
    out.append(text[pos:])
    return "".join(out)


def _open(is_bracket: bool) -> str:
    return "\\\\[\n" if is_bracket else "$$\n"


def _close(is_bracket: bool) -> str:
    return "\\end{aligned}\n\\\\]" if is_bracket else "\\end{aligned}\n$$"


def _rel_class(tok: str) -> str:
    if tok in ("+", "\\+"):
        return "plus"
    if tok in ("-", "\\-"):
        return "minus"
    return "eq"


def _join_prefix(cls: str, tok: str) -> str:
    if cls == "plus":
        return "&+ "
    if cls == "minus":
        return "&- "
    return f"&{tok} " if tok != "=" else "&= "


def _aligned_lines(body: str) -> str:
    # Tokenize content lines into (text, leading-join) units. Three shapes:
    #   lone `=`/`+`/`\-` line  -> pending join for the NEXT content line
    #   leading `= ...`         -> joins with the PREVIOUS content line
    #   trailing `X =`          -> X joins with the NEXT content line
    seq = [l.strip() for l in body.split("\n")]
    while seq and not seq[0]:
        seq.pop(0)
    while seq and not seq[-1]:
        seq.pop()
    rows: list[str] = []
    pending: tuple[str, str] | None = None  # (class, token)
    cur: str | None = None
    first = True

    def flush_join(nxt: str, cls: str, tok: str):
        nonlocal cur, first
        assert cur is not None
        if first:
            rows.append(f"{cur} {_join_prefix(cls, tok)}{nxt}")
            first = False
        else:
            rows.append(f"{_join_prefix(cls, tok)}{nxt}")
        cur = None

    for tok in seq:
        mm = LONE_REL_RE.match(tok)
        if mm:
            pending = (_rel_class(mm.group(1)), mm.group(1))
            continue
        lm = LEAD_REL_RE.match(tok)
        if lm:
            cls, t = _rel_class(lm.group(1)), lm.group(1)
            if cur is not None:
                flush_join(lm.group(2), cls, t)
            else:
                # Leading relation with no LHS (continuation of a finished
                # join): open a new aligned row.
                rows.append(f"{_join_prefix(cls, t)}{lm.group(2)}")
                first = False
            pending = None
            continue
        tm = TRAIL_EQ_RE.match(tok)
        if tm and pending is None:
            if cur is not None:
                # `cur` then `X =`: close cur as its own row first... in
                # practice cur is None here (joins chain); be safe anyway.
                rows.append(cur)
                first = False
            cur = tm.group("head").strip()
            pending = ("eq", "=")
            continue
        if cur is None:
            if pending is not None:
                # Block opens with a join and no LHS (e.g. leading lone `=`):
                # emit an open RHS row rather than inventing a left side.
                cls, t = pending
                pending = None
                rows.append(f"{_join_prefix(cls, t)}{tok}")
                first = False
            else:
                cur = tok
            continue
        if pending is not None:
            cls, t = pending
            pending = None
            flush_join(tok, cls, t)
        else:
            # Adjacent content lines with no relation between them: keep both
            # as separate rows rather than gluing (safer than guessing).
            rows.append(cur)
            first = False
            cur = tok
    if cur is not None:
        rows.append(cur)
    body_out = " \\\\\n".join(rows)
    return "\\\\begin{aligned}\n" + body_out + "\n"


def main() -> int:
    fix = "--fix" in sys.argv
    report = []
    touched = 0
    for p in sorted(ROOT.joinpath("data/lessons").rglob("*.md")):
        rel = str(p.relative_to(ROOT))
        text = p.read_text(encoding="utf-8")
        if fix:
            new_text, counts = fix_text(text)
            if new_text != text:
                p.write_text(new_text, encoding="utf-8")
                touched += 1
                print(f"fixed {rel}: {counts}")
            for iss in lint_text(rel, new_text):
                # Re-lint after fix: anything remaining needs hand review.
                report.append(iss)
        else:
            report.extend(lint_text(rel, text))
    REPORT.write_text(json.dumps(report, indent=2, ensure_ascii=False))
    if fix:
        print(f"--fix: {touched} files rewritten; {len(report)} issues left for review")
        for iss in report:
            print(f"  remaining {iss['type']}: {iss['file']}:{iss['line']}")
        return 0
    by_type: dict[str, int] = {}
    by_file: dict[str, int] = {}
    for iss in report:
        by_type[iss["type"]] = by_type.get(iss["type"], 0) + 1
        by_file[iss["file"]] = by_file.get(iss["file"], 0) + 1
    print(f"lens #4: {len(report)} issues in {len(by_file)} files: {by_type}")
    for f, n in sorted(by_file.items(), key=lambda kv: -kv[1])[:20]:
        print(f"  {n:3d}  {f}")
    # fused-display is info (owned by the frontend fused-inline defense);
    # everything else hard-fails.
    return 1 if any(i["type"] in ERROR_TYPES for i in report) else 0


if __name__ == "__main__":
    raise SystemExit(main())
