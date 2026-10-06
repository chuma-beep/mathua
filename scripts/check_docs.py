#!/usr/bin/env python3
"""
check_docs.py — find claims in the engineering docs that the code or the corpus contradicts.

Why this exists, in the shape of the bug that motivated it: README stated a 150 XP
mastery-check gate twice while `internal/xp/policy.go` said 50. Nothing was broken, no test
failed, and the landing page carried the same stale figure until it was fixed separately. A
documented number and the constant it describes are the same fact written down twice, and
nothing noticed.

The design rule that follows from that: **this script never hardcodes a fact it can read.** The
XP gate comes from the Go source, the corpus counts come from the corpus, and the KP-per-shard
rule comes from the shards. A checker that pins its own copy of the number goes stale in exactly
the way it was written to prevent.

Scope is the engineering docs only. `data/lessons/**` is learner-facing corpus content, is
covered by `audit_lessons.py`, and is deliberately not subject to American spelling or to the
link rules here.

Usage:
    python3 scripts/check_docs.py
Exit code 0 = no errors. Exit code 1 = at least one error.
"""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

# ---------------------------------------------------------------- what we check

# Engineering docs. Deliberately explicit rather than a glob, so a new doc is a deliberate act
# rather than something that silently falls under or out of the rules.
DOCS = [
    "README.md",
    "DESIGN.md",
    "CONTRIBUTING.md",
    "AGENTS.md",
    "CONTEXT.md",
    "docs/architecture.md",
    "docs/efficacy.md",
    "docs/lesson-media.md",
    "docs/mastery-assessments.md",
    "docs/mastery-semantics.md",
    "docs/math-input.md",
    "docs/mobile-device-checklist.md",
    "docs/system-design.md",
    "docs/system-design-slides.md",
    "docs/wasc-uc-ag-research.md",
]

# Files exempt from the spelling rule, each with its reason. An exemption without a reason is
# indistinguishable from a mistake.
SPELLING_EXEMPT = {
    "THIRD_PARTY_NOTICES.md": "names third-party license instruments (MIT, CC BY-NC 4.0); "
                              "their wording is not ours to rewrite",
}

# Paths under the Next.js app, which resolve as routes at runtime and have no file of that name.
APP_ROUTES = ("/docs", "/learn", "/study", "/review", "/onboard", "/profile", "/graph",
              "/plan", "/goals", "/history", "/leaderboard", "/settings", "/admin",
              "/transcript", "/home", "/rewards")

BRITISH = {
    "behaviour": "behavior", "behaviours": "behaviors",
    "colour": "color", "colours": "colors", "coloured": "colored", "colouring": "coloring",
    "labelled": "labeled", "labelling": "labeling",
    "licence": "license", "licences": "licenses",
    "centre": "center", "centres": "centers",
    "defence": "defense", "analyse": "analyze", "analysed": "analyzed",
    "organisation": "organization", "organisations": "organizations",
    "recognise": "recognize", "normalise": "normalize", "canonicalise": "canonicalize",
    "whilst": "while", "amongst": "among",
}


# Files allowed to cite a gate figure other than the current one, each with its reason. Both cite
# the old 150 XP value *as the example of the bug*, which is the only place that value should
# still appear. A whole-file exemption is a hole, so the list is short, explicit, and justified.
HISTORICAL_GATE_CITERS = {
    "CONTEXT.md": "ADR-020 records the 150 -> 50 rescale as history",
    "AGENTS.md": "quotes the stale 150 XP figure as the motivating example for this rule",
}


class Report:
    def __init__(self) -> None:
        self.errors: list[tuple[str, str, str]] = []

    def error(self, where: str, rule: str, msg: str) -> None:
        self.errors.append((where, rule, msg))


# ---------------------------------------------------------------- helpers

def read(rel: str) -> str:
    return (ROOT / rel).read_text(encoding="utf-8")


def strip_code(text: str) -> str:
    """Blank out regions where a word is an identifier, not prose.

    Without this, a doc that mentions `behaviour` inside a fenced block or inline code trips the
    spelling rule on a token the author does not control.
    """
    out = re.sub(r"```.*?```", lambda m: "\n" * m.group(0).count("\n"), text, flags=re.S)
    out = re.sub(r"`[^`\n]*`", lambda m: " " * len(m.group(0)), out)
    return out


def load_kps(path: Path) -> list:
    """A shard is a bare JSON list in some files and an object with a "kps" key in others."""
    data = json.loads(path.read_text(encoding="utf-8"))
    if isinstance(data, dict):
        return data.get("kps") or []
    return data if isinstance(data, list) else []


def go_const(path: str, name: str) -> int | None:
    """Read an integer constant out of Go source. Returns None if absent."""
    m = re.search(rf"^\s*(?:const\s+)?{re.escape(name)}\s*=\s*(\d+)", read(path), re.M)
    return int(m.group(1)) if m else None


# ---------------------------------------------------------------- checks

def check_xp_gate(rep: Report) -> None:
    """Any mastery-check gate figure in the docs must equal the constant in the code."""
    gate = go_const("internal/xp/policy.go", "QuizGateXP")
    if gate is None:
        rep.error("internal/xp/policy.go", "quiz-gate",
                  "QuizGateXP not found; this checker cannot verify the docs")
        return
    # "150 XP" and "50 XP" near a gate/mastery-check/quiz mention. Historical references to the
    # old value are allowed only in CONTEXT.md, which records the rescale as history.
    pat = re.compile(r"(\d+)\s*XP", re.I)
    ctx = re.compile(r"gate|mastery[- ]check|quiz", re.I)
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        for i, line in enumerate(read(doc).split("\n"), 1):
            if not ctx.search(line):
                continue
            for m in pat.finditer(line):
                n = int(m.group(1))
                if n == gate or doc in HISTORICAL_GATE_CITERS:
                    continue
                rep.error(f"{doc}:{i}", "xp-gate",
                          f'"{m.group(0)}" but xp.QuizGateXP is {gate}')


def sentences(text: str) -> list[tuple[int, str]]:
    """Split into (line number, sentence). A rule is stated per sentence, not per line."""
    out = []
    for i, line in enumerate(text.split("\n"), 1):
        s = line.strip()
        if not s or s.startswith(("|", "#", "-", "*", ">", "`")):
            continue
        for part in re.split(r"(?<=[.!?:])\s+", s):
            if part.strip():
                out.append((i, part.strip()))
    return out


def check_kp_rule(rep: Report) -> None:
    """The per-concept KP rule must match the corpus, and must be stated one way in the docs."""
    counts: dict[int, int] = {}
    for f in sorted((ROOT / "data/lessons/kp").glob("*.json")):
        kps = load_kps(f)
        if kps:
            counts[len(kps)] = counts.get(len(kps), 0) + 1
    if not counts:
        rep.error("data/lessons/kp", "kp-rule", "no KP shards found")
        return
    if len(counts) > 1:
        rep.error("data/lessons/kp", "kp-rule",
                  f"shards hold different KP counts: {counts}; the documented rule cannot be "
                  "one number, so state the real distribution")
        return
    kps_per = next(iter(counts))
    shards = sum(counts.values())
    total = sum(k * v for k, v in counts.items())

    # A rule claim is a number that is explicitly scoped: either "exactly N KPs", or "N KPs per
    # concept". Corpus totals such as "657 KP shard files x 3 = 1971 KPs" are not claims about
    # the per-shard rule, and reading them as such flags the arithmetic that proves the rule.
    claim = re.compile(r"([0-9]+)(?:\s*[-\u2013]\s*([0-9]+))?\s*(?:KPs?\b|knowledge points?\b)",
                       re.I)
    scoped = re.compile(r"\bper\s+(?:concept|shard|lesson|topic)\b|\bexactly\b", re.I)
    for doc in ("README.md", "CONTRIBUTING.md"):
        for line_no, sent in sentences(read(doc)):
            if not scoped.search(sent):
                continue
            for m in claim.finditer(sent):
                lo = int(m.group(1))
                hi = m.group(2)
                if lo != kps_per or (hi is not None and int(hi) != kps_per):
                    got = f"{lo}-{hi}" if hi else str(lo)
                    rep.error(f"{doc}:{line_no}", "kp-rule",
                              f"docs say {got} KPs per concept but all {shards} shards hold "
                              f"exactly {kps_per} ({total} KPs total)")


# A corpus-scale claim. Anything smaller is a per-concept parameter (a streak, a threshold, a
# percentage) that happens to sit next to the word "concept", and flagging those would bury the
# real findings in noise.
CORPUS_SCALE = 100

def check_corpus_figures(rep: Report) -> None:
    """Corpus counts quoted in docs must match the corpus. Reads the corpus, never a literal."""
    # Count concepts, not domain files: data/concepts holds one file per domain, and 17 files is
    # not the number any doc means by "concepts".
    n_concepts = 0
    for f in sorted((ROOT / "data/concepts").glob("*.json")):
        if f.name == "enrichment.json":
            continue
        data = json.loads(f.read_text(encoding="utf-8"))
        n_concepts += len(data if isinstance(data, list) else list(data.values()))
    facts = {
        "concepts": n_concepts,
        "shards": len(list((ROOT / "data/lessons/kp").glob("*.json"))),
        "kps": sum(len(load_kps(f)) for f in (ROOT / "data/lessons/kp").glob("*.json")),
        "courses": len(json.loads(read("data/courses.json"))),
    }
    nouns = {
        "concepts": r"concepts?",
        "shards": r"(?:KP\s+)?shards?",
        "kps": r"KPs|knowledge points",
        "courses": r"courses",
    }
    pat = re.compile(r"([0-9][0-9,]*)\s+(" + "|".join(f"(?:{v})" for v in nouns.values()) + r")\b",
                     re.I)
    # "128 concepts" can be a per-domain subset ("calculus is 128 concepts") rather than a corpus
    # total. Only a sentence that scopes the claim to the whole corpus is checked, because a
    # subset count is legitimate and cannot be verified without knowing which subset is meant.
    scope = re.compile(r"\b(?:corpus|graph|catalogue|catalog|total|all|across|every|whole)\b", re.I)
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, sent in sentences(read(doc)):
            if sent.startswith("|") or not scope.search(sent):
                continue
            for m in pat.finditer(sent):
                n = int(m.group(1).replace(",", ""))
                if n < CORPUS_SCALE:
                    continue
                noun = m.group(2).lower()
                key = ("shards" if "shard" in noun
                       else "kps" if noun.startswith("kp") or "knowledge" in noun
                       else "courses" if noun.startswith("course") else "concepts")
                if n != facts[key]:
                    rep.error(f"{doc}:{line_no}", "corpus-count",
                              f'"{m.group(0)}" but the corpus has {facts[key]} {key}')


def check_adr_refs(rep: Report) -> None:
    """Every ADR a doc cites must be defined in CONTEXT.md."""
    defined = set(re.findall(r"^- (ADR-\d{3})\b", read("CONTEXT.md"), re.M))
    if not defined:
        rep.error("CONTEXT.md", "adr", "no ADR definitions found; the pattern changed?")
        return
    for doc in DOCS:
        if doc == "CONTEXT.md":
            continue
        for i, line in enumerate(read(doc).split("\n"), 1):
            for m in re.finditer(r"\bADR-\d{3}\b", line):
                if m.group(0) not in defined:
                    rep.error(f"{doc}:{i}", "adr-ref",
                              f"{m.group(0)} is cited but not defined in CONTEXT.md")


def check_links(rep: Report) -> None:
    """Relative links must resolve, as a file or as an app route."""
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        text = read(doc)
        for i, line in enumerate(text.split("\n"), 1):
            for m in re.finditer(r"\[[^\]]*\]\(([^)\s]+)\)", line):
                tgt = m.group(1)
                if tgt.startswith(("http://", "https://", "mailto:", "#")):
                    continue
                if tgt.startswith("<"):
                    # A placeholder such as <../parabola/>, used inside audit tables.
                    continue
                if tgt.startswith("/"):
                    if any(tgt == r or tgt.startswith(r + "/") for r in APP_ROUTES):
                        continue
                    rep.error(f"{doc}:{i}", "link", f"{tgt} looks like a route but matches no "
                              f"known route, and no such file exists either")
                    continue
                target = tgt.split("#")[0]
                if not target:
                    continue
                if not (p.parent / target).resolve().exists():
                    rep.error(f"{doc}:{i}", "link", f"{tgt} does not exist")


def check_spelling(rep: Report) -> None:
    """Engineering docs use American spelling. Exempt files state why."""
    targets = [d for d in DOCS if d not in SPELLING_EXEMPT]
    targets += [f for f in SPELLING_EXEMPT if (ROOT / f).is_file()]
    for doc in targets:
        p = ROOT / doc
        if not p.is_file():
            continue
        exempt = doc in SPELLING_EXEMPT
        for i, line in enumerate(read(doc).split("\n"), 1):
            if exempt:
                continue
            for w, american in BRITISH.items():
                for m in re.finditer(rf"\b{w}\b", line, re.I):
                    if m.group(0) == american:
                        continue
                    in_code = False
                    # skip occurrences inside inline code on this line
                    for cm in re.finditer(r"`[^`\n]*`", line):
                        if cm.start() <= m.start() < cm.end():
                            in_code = True
                    if in_code:
                        continue
                    rep.error(f"{doc}:{i}", "spelling",
                              f'"{m.group(0)}" is British English; use "{american}"')


def main() -> int:
    rep = Report()
    check_xp_gate(rep)
    check_kp_rule(rep)
    check_corpus_figures(rep)
    check_adr_refs(rep)
    check_links(rep)
    check_spelling(rep)

    for where, rule, msg in rep.errors:
        print(f"ERROR   {where} [{rule}] {msg}", file=sys.stderr)
    if rep.errors:
        print(f"\nFAIL: {len(rep.errors)} documentation inconsistency/inconsistencies.",
              file=sys.stderr)
        return 1
    print("OK: docs agree with the code and the corpus "
          "(XP gate, KP rule, corpus counts, ADR refs, links, spelling)")
    return 0


if __name__ == "__main__":
    sys.exit(main())