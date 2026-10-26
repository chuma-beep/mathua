#!/usr/bin/env python3
"""
check_docs.py, find claims in the engineering docs that the code or the corpus contradicts.

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
import os
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

# Built once per run by file_index(); see check_symbol_citations.
_FILE_INDEX = None

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
    "docs/administration.md",
    "docs/code-defects.md",
    "docs/content-expansion.md",
    "docs/diagram-style.md",
    "docs/efficacy.md",
    "docs/figure-reclaim.md",
    "docs/glossary-provenance.md",
    "docs/lesson-media.md",
    "docs/mastery-assessments.md",
    "docs/mastery-semantics.md",
    "docs/math-input.md",
    "docs/mobile-device-checklist.md",
    "docs/system-design.md",
    "docs/system-design-slides.md",
    "docs/wasc-uc-ag-research.md",
]
# The ADRs are a generated, uniform set, so a glob is the honest description rather than 54
# lines that each need editing when a decision is added. They were inline in CONTEXT.md and
# were covered by the rules below until the split, so they stay covered after it.
DOCS += [str(p.relative_to(ROOT)) for p in sorted((ROOT / "docs/adr").glob("ADR-*.md"))]

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
    "behaviour": "behavior", "behaviours": "behaviors", "behavioural": "behavioral",
    "colour": "color", "colours": "colors", "coloured": "colored", "colouring": "coloring",
    "labelled": "labeled", "labelling": "labeling",
    "licence": "license", "licences": "licenses",
    "centre": "center", "centres": "centers",
    "defence": "defense",
    "analyse": "analyze", "analysed": "analyzed", "analyses": "analyzes", "analysing": "analyzing",
    "organisation": "organization", "organisations": "organizations",
    "recognise": "recognize", "recognised": "recognized", "recognises": "recognizes",
    "normalise": "normalize", "normalised": "normalized", "normalises": "normalizes",
    "canonicalise": "canonicalize", "canonicalised": "canonicalized",
    "serialise": "serialize", "serialised": "serialized",
    "prioritise": "prioritize", "prioritised": "prioritized",
    "whilst": "while", "amongst": "among",
}


# Files allowed to cite a gate figure other than the current one, each with its reason. Both cite
# the old 150 XP value *as the example of the bug*, which is the only place that value should
# still appear. A whole-file exemption is a hole, so the list is short, explicit, and justified.
HISTORICAL_GATE_CITERS = {
    "CONTEXT.md": "ADR-020 records the 150 -> 50 rescale as history",
    "AGENTS.md": "quotes the stale 150 XP figure as the motivating example for this rule",
    "docs/adr/ADR-003-q3.md": "the superseded decision whose title names ADR-020 as its replacement",
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
    """Any mastery-check gate figure in the docs must equal the constant in the code.

    The gate word and the figure must be close together. An earlier version flagged any "N XP"
    on a line that mentioned "quiz" or "mastery", which fired on legitimate statements of other
    XP amounts -- "diagnostic 0 XP", a historical "+3 XP" bug fix -- because those words are
    everywhere in this corpus. Requiring the figure within 40 characters of the gate word keeps
    the rule on the claim it was written for.
    """
    gate = go_const("internal/xp/policy.go", "QuizGateXP")
    if gate is None:
        rep.error("internal/xp/policy.go", "quiz-gate",
                  "QuizGateXP not found; this checker cannot verify the docs")
        return
    pat = re.compile(r"(\d+)\s*XP", re.I)
    gateword = re.compile(r"\bgate\b|mastery[- ]check", re.I)
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        for i, line in enumerate(read(doc).split("\n"), 1):
            positions = [g.start() for g in gateword.finditer(line)]
            if not positions:
                continue
            for m in pat.finditer(line):
                if not any(abs(m.start() - g) <= 40 for g in positions):
                    continue
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


def prose_sentences(text: str) -> list[tuple[int, str]]:
    """Sentences for a rule that reads a fact out of prose, bullet items included.

    `sentences` skips a line that begins with a list marker, which is right for a rule about
    spelling or a link: a bullet there is usually a name or a path. A corpus count is what a
    bullet most often states, though, and three counts this file is meant to check had drifted
    where it could not look -- a CONTRIBUTING bullet, a line of `docs/architecture.md` that
    began with an inline code span, and an ADR. A guard that reads nothing reports nothing, so
    the count rules read bullets and blank inline code first, because a count inside `backticks`
    is an identifier rather than a claim.
    """
    out = []
    for i, line in enumerate(strip_code(text).split("\n"), 1):
        s = re.sub(r"^[-\*\+]\s+|^\d+[.)]\s+", "", line.strip())
        if not s or s.startswith(("|", "#", ">")):
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
        for line_no, sent in prose_sentences(read(doc)):
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
    # A doc may also quote the size the corpus had when it measured something, and then it has to
    # say so where the number is. "459 of the 657 the corpus then held" is a record of 2026-10-05
    # and not a claim about today; ADR-038, ADR-040 and ADR-041 already write it that way. The
    # window is 40 characters, the proximity the XP gate rule above uses, because one of these
    # sentences runs past 1,000 characters and a marker anywhere inside it is not on this number.
    historical = re.compile(r"\bthen held\b|\bat the time\b", re.I)
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, sent in prose_sentences(read(doc)):
            if not scope.search(sent):
                continue
            for m in pat.finditer(sent):
                n = int(m.group(1).replace(",", ""))
                if n < CORPUS_SCALE:
                    continue
                noun = m.group(2).lower()
                key = ("shards" if "shard" in noun
                       else "kps" if noun.startswith("kp") or "knowledge" in noun
                       else "courses" if noun.startswith("course") else "concepts")
                if n != facts[key] and not historical.search(sent[m.end():m.end() + 40]):
                    rep.error(f"{doc}:{line_no}", "corpus-count",
                              f'"{m.group(0)}" but the corpus has {facts[key]} {key}')


def check_generator_registries(rep: Report) -> None:
    """A quoted count of generator registries must match the packages all.go pulls in.

    This is a different fact from the number of corpus domains, and the two drifted apart
    silently: `system-design.md` said "17 domain registries" while `internal/generator/all`
    imported 22 and no rule read either. The count is taken from the import list rather than
    from a directory listing, because a package can exist without being registered.
    """
    all_go = read("internal/generator/all/all.go")
    n = len(re.findall(r'"github\.com/chuma-beep/mathua/internal/generator/[a-z]+"', all_go))
    if n == 0:
        rep.error("internal/generator/all/all.go", "generator-registries",
                  "found no generator imports; this rule would otherwise pass vacuously")
        return
    pat = re.compile(r"([0-9][0-9,]*)\s+domain registr(?:y|ies)", re.I)
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, line in enumerate(read(doc).splitlines(), 1):
            for m in pat.finditer(line):
                got = int(m.group(1).replace(",", ""))
                if got != n:
                    rep.error(f"{doc}:{line_no}", "generator-registries",
                              f'"{m.group(0)}" but all.go imports {n} generator packages')


def check_course_reach(rep: Report) -> None:
    """How much of the graph the course catalog reaches must match the doc that quotes it.

    A course is a set of target concepts, and a target reaches everything in its transitive
    prerequisite closure (`collectChain` in internal/planning/planner.go). The union of those
    closures is the content the catalog can route a learner to, and it is smaller than the
    corpus: a course samples a path, it does not cover a domain. The figure is read from
    data/courses.json and data/concepts here, so the doc that states it cannot drift the way
    README's XP gate once drifted from xp.QuizGateXP.

    It requires the claim to be present, because a checker that finds no number to check
    passes and reports nothing -- the failure mode this file is written to avoid.
    """
    doc = "docs/content-expansion.md"
    path = ROOT / doc
    if not path.is_file():
        rep.error(doc, "course-reach", "the doc that states the reach no longer exists")
        return

    prereq: dict = {}
    for f in sorted((ROOT / "data/concepts").glob("*.json")):
        if f.name == "enrichment.json":
            continue
        for c in json.loads(f.read_text(encoding="utf-8")):
            prereq[c["id"]] = set(c.get("prerequisites") or [])
    if not prereq:
        rep.error("data/concepts", "course-reach", "no concepts read; cannot verify the docs")
        return

    def closure(ids: list) -> set:
        seen: set = set()
        stack = list(ids)
        while stack:
            cid = stack.pop()
            if cid in seen:
                continue
            seen.add(cid)
            stack.extend(prereq.get(cid, ()))
        return seen

    reach: set = set()
    for course in json.loads(read("data/courses.json")):
        reach |= closure(course["targets"])

    claim = re.compile(r"(\d[\d,]*)\s+of the\s+(\d[\d,]*)\s+concepts")
    found = 0
    for line_no, line in enumerate(path.read_text(encoding="utf-8").split("\n"), 1):
        for m in claim.finditer(line):
            found += 1
            reached = int(m.group(1).replace(",", ""))
            total = int(m.group(2).replace(",", ""))
            if reached != len(reach) or total != len(prereq):
                rep.error(f"{doc}:{line_no}", "course-reach",
                          f'"{m.group(0)}" but the catalog reaches {len(reach)} of the '
                          f"{len(prereq)} concepts")
    if found == 0:
        rep.error(doc, "course-reach",
                  "the course-reach claim is gone; a check that reads nothing reports nothing")


def check_media_candidates(rep: Report) -> None:
    """Orphan caption counts quoted in the docs must match the corpus that owns them.

    `scripts/audit_lesson_media.py` rates every paragraph that reads as a caption whose image
    was dropped. Three docs quote the result, and a caption repair moves it, so the counts are
    read from the corpus here rather than copied into this file: a checker that pinned 214
    would go stale the next time a figure lands or a caption changes shape.

    It requires the claim in each doc, because a checker that finds no number to check passes
    and reports nothing -- the failure mode this file is written to avoid.
    """
    sys.path.insert(0, str(ROOT / "scripts"))
    import audit_lesson_media as media

    tiers: dict[str, int] = {}
    for row in media.section_a_candidates():
        tiers[row["confidence"]] = tiers.get(row["confidence"], 0) + 1
    facts = {
        "paragraph": sum(tiers.values()), "candidate": sum(tiers.values()),
        "certain": tiers.get("certain", 0), "strong": tiers.get("strong", 0),
        "likely": tiers.get("likely", 0),
    }
    claim = re.compile(r"(\d[\d,]*)\s+(paragraphs?|candidates?|certain|strong|likely)\b")
    for doc in ("docs/content-expansion.md", "docs/figure-reclaim.md", "docs/lesson-media.md"):
        path = ROOT / doc
        if not path.is_file():
            rep.error(doc, "media-candidates", "the doc that states the counts no longer exists")
            continue
        found = 0
        for line_no, line in enumerate(path.read_text(encoding="utf-8").split("\n"), 1):
            for m in claim.finditer(line):
                found += 1
                word = m.group(2).rstrip("s")
                if int(m.group(1).replace(",", "")) != facts[word]:
                    rep.error(f"{doc}:{line_no}", "media-candidates",
                              f'"{m.group(0)}" but the orphan-caption audit counts '
                              f"{facts[word]}")
        if found == 0:
            rep.error(doc, "media-candidates",
                      "the orphan-caption count is gone; a check that reads nothing reports "
                      "nothing")


def check_figure_dispositions(rep: Report) -> None:
    """Sections 4 and 6 of the figure plan must count the same rows, and the audit the same total.

    `docs/figure-reclaim.md` section 6 holds one row per pending figure with the disposition that
    section 4 defines, and section 4's table counts those rows by disposition. The counts have
    drifted before: the first fetch batch landed while section 4 read 43 redraw and 37 fetch and the
    rows below it said 53 and 44, and a reader planning work off the table was reading a plan that
    had already been executed.

    The number of rows is owned by the corpus, because `scripts/audit_lesson_media.py` decides which
    paragraphs are orphaned captions, so it is read from there rather than pinned here. The
    dispositions stay the document's decision: this compares the three counts it can see, and fails
    on any disagreement, including a table that has gone empty without saying so.
    """
    doc = "docs/figure-reclaim.md"
    path = ROOT / doc
    if not path.is_file():
        rep.error(doc, "figure-dispositions", "the figure plan no longer exists")
        return
    lines = path.read_text(encoding="utf-8").split("\n")

    header = "| Figure | Caption (repaired) | Book | Section | Lessons | Disposition |"
    start = next((i for i, l in enumerate(lines) if l.strip() == header), -1)
    if start < 0:
        rep.error(doc, "figure-dispositions",
                  "section 6's table is gone; a check that reads nothing reports nothing")
        return
    pending: dict[str, int] = {}
    for i in range(start + 2, len(lines)):
        if not lines[i].startswith("|"):
            break
        cells = [c.strip() for c in lines[i].strip().strip("|").split("|")]
        if len(cells) != 6:
            rep.error(f"{doc}:{i + 1}", "figure-dispositions",
                      f"a row has {len(cells)} cells, not the 6 the header names")
            continue
        pending[cells[-1]] = pending.get(cells[-1], 0) + 1

    counted: dict[str, int] = {}
    section = next((i for i, l in enumerate(lines) if l.startswith("## 4. ")), -1)
    if section < 0:
        rep.error(doc, "figure-dispositions", "section 4 is gone, so there are no counts to read")
        return
    for i in range(section, len(lines)):
        if lines[i].startswith("## 5. "):
            break
        m = re.match(r"^\|\s*([a-z][a-z, \-]*?)\s*\|\s*(\d+)\s*\|", lines[i])
        if m and m.group(1) != "Disposition":
            counted[m.group(1)] = int(m.group(2))
    if not counted:
        rep.error(doc, "figure-dispositions",
                  "section 4's table has no rows; a check that reads nothing reports nothing")
        return

    sys.path.insert(0, str(ROOT / "scripts"))
    import audit_lesson_media as media

    certain = sum(1 for row in media.section_a_candidates() if row["confidence"] == "certain")
    total = sum(pending.values())
    if total != certain:
        rep.error(f"{doc}:{start + 1}", "figure-dispositions",
                  f"section 6 lists {total} rows and the audit rates {certain} certain captions")
    for name, n in sorted(pending.items()):
        if name not in counted:
            rep.error(f"{doc}:{start + 1}", "figure-dispositions",
                      f"section 6 uses the disposition {name!r}, which section 4 does not define")
        elif counted[name] != n:
            rep.error(f"{doc}:{start + 1}", "figure-dispositions",
                      f"section 4 counts {counted[name]} {name} and section 6 has {n} such rows")
    for name, n in sorted(counted.items()):
        if n and not pending.get(name):
            rep.error(f"{doc}:{start + 1}", "figure-dispositions",
                      f"section 4 counts {n} {name} and section 6 has no such row")


def check_thin_lessons(rep: Report) -> None:
    """Section 4.3's sub-1 KB table must be the corpus's own list, sizes and count.

    The table was kept by hand and drifted in both directions: four rows stayed in it after those
    lessons had grown past the threshold, and three authored lessons entered the corpus without a
    row at all, so the sentence above it under-reported the work. The sizes belong to files the
    corpus registers, so they are read from `data/lessons/lessons.json` and `os.stat` rather than
    pinned here; only the 1 KB threshold is this checker's own, because that is the threshold the
    table's heading names.
    """
    doc = "docs/content-expansion.md"
    path = ROOT / doc
    if not path.is_file():
        rep.error(doc, "thin-lessons", "the document that holds section 4.3 no longer exists")
        return

    thin: dict[str, tuple[str, int]] = {}
    for e in json.loads((ROOT / "data/lessons/lessons.json").read_text(encoding="utf-8")):
        p = ROOT / "data/lessons" / e["source"]
        if p.is_file() and p.stat().st_size < 1024:
            thin[e["concept_id"]] = (str(p.relative_to(ROOT)), p.stat().st_size)
    lines = path.read_text(encoding="utf-8").split("\n")
    start = next((i for i, l in enumerate(lines) if l.startswith("### 4.3 ")), -1)
    end = next((i for i in range(start + 1, len(lines)) if lines[i].startswith("### ")), len(lines))
    if start < 0:
        rep.error(doc, "thin-lessons", "section 4.3 is gone; a check that reads nothing reports nothing")
        return

    row = re.compile(r"^\|\s*`([\w.]+)`\s*\|\s*`([^`]+)`\s*\|\s*(\d+)\s*\|\s*$")
    listed: dict[str, tuple[str, int]] = {}
    for i in range(start, end):
        m = row.match(lines[i])
        if not m:
            continue
        listed[m.group(1)] = (m.group(2), int(m.group(3)))
        if m.group(1) not in thin:
            rep.error(f"{doc}:{i + 1}", "thin-lessons",
                      f"the table lists {m.group(1)}, which is not under 1 KB")
        else:
            path_now, size_now = thin[m.group(1)]
            if int(m.group(3)) != size_now:
                rep.error(f"{doc}:{i + 1}", "thin-lessons",
                          f"{m.group(1)} is {size_now} bytes and the table says {m.group(3)}")
            if m.group(2) != path_now:
                rep.error(f"{doc}:{i + 1}", "thin-lessons",
                          f"{m.group(1)} lives at {path_now} and the table says {m.group(2)}")
    if not thin:
        # The corpus has no lesson under the threshold. That is the finished state, so the section
        # has to say it and list nothing; a row here would be a stale claim about a lesson that grew.
        zero = re.compile(r"no lessons? (?:is|are) under 1 KB")
        if not any(zero.search(l) for l in lines[start:end]):
            rep.error(doc, "thin-lessons",
                      "no lesson is under 1 KB and section 4.3 does not say so")
        return
    # Name every thin lesson the table misses. This is also the guard against a parse that found no
    # rows at all: a parser failure then reports one error per lesson instead of passing silently.
    for cid in sorted(thin):
        if cid not in listed:
            rep.error(doc, "thin-lessons",
                      f"{cid} is {thin[cid][1]} bytes and the table has no row for it")

    # The prose writes small counts as words ("six lessons are under 1 KB"), so both forms count.
    words = {w: i for i, w in enumerate(
        "zero one two three four five six seven eight nine ten eleven twelve".split())}
    claim = re.compile(
        r"(\d[\d,]*|" + "|".join(w for w in words if w != "zero") + r") lessons are under 1 KB")
    found = 0
    for line_no, line in enumerate(lines, 1):
        for m in claim.finditer(line):
            found += 1
            raw = m.group(1).replace(",", "")
            said = int(raw) if raw.isdigit() else words[raw]
            if said != len(thin):
                rep.error(f"{doc}:{line_no}", "thin-lessons",
                          f'"{m.group(0)}" but the corpus has {len(thin)}')
    if found == 0:
        rep.error(doc, "thin-lessons",
                  "the sentence above the table no longer states the count")


def check_target_table(rep: Report) -> None:
    """Section 5.5's target table must agree with the corpus, and its columns must add up.

    The table is the roadmap for the thousand: one row per domain, the concepts it holds now, the
    size the curriculum decided on, and the difference. Only `Target` is a decision, `Now` is read
    from `data/concepts/`, as the sentence above the table says and that is the column that drifted.
    The table was drawn when the corpus held 703 concepts; `scripts/sweep_counts.py` keeps the
    sentence above it current, because its size is a pattern rule, but it has no rule for the rows,
    so ten of the seventeen sat at their 703-era values while the corpus grew to 819. The roadmap
    therefore named deficits larger than the work left, and two domains that had reached their
    targets read as though they had not.

    Four things are read rather than pinned: each row's `Now` against the corpus, each row's `Add`
    against `Target - Now`, a row for every domain the corpus holds, and the sentence above the table
    against the sum of the two columns. A checker that pins its own copy of a target goes stale in
    exactly the way this one was written to prevent, so the targets stay the document's.
    """
    doc = "docs/content-expansion.md"
    path = ROOT / doc
    if not path.is_file():
        rep.error(doc, "target-table", "the doc that holds the target table no longer exists")
        return

    held: dict[str, int] = {}
    for f in sorted((ROOT / "data/concepts").glob("*.json")):
        if f.name == "enrichment.json":
            continue
        for c in json.loads(f.read_text(encoding="utf-8")):
            if "id" in c:
                held[c["domain"]] = held.get(c["domain"], 0) + 1
    if not held:
        rep.error("data/concepts", "target-table", "no concepts read; cannot verify the table")
        return

    lines = path.read_text(encoding="utf-8").split("\n")
    header = "| Domain | Tier | Now | Target | Add |"
    start = next((i for i, l in enumerate(lines) if l.strip() == header), -1)
    if start < 0:
        rep.error(doc, "target-table",
                  "the target table is gone; a check that reads nothing reports nothing")
        return

    row = re.compile(r"^\|\s*(\w+)\s*\|\s*(\w+)\s*\|\s*(\d+)\s*\|\s*(\d+)\s*\|\s*\+(\d+)\s*\|\s*$")
    seen: dict[str, tuple[int, int, int]] = {}
    n = 0
    for i in range(start + 2, len(lines)):
        m = row.match(lines[i])
        if not m:
            break
        n += 1
        dom, now, target, add = m.group(1), int(m.group(3)), int(m.group(4)), int(m.group(5))
        seen[dom] = (now, target, add)
        if now != held.get(dom, 0):
            rep.error(f"{doc}:{i + 1}", "target-table",
                      f"Now says {now} for {dom} but the corpus holds {held.get(dom, 0)}")
        if add != target - now:
            rep.error(f"{doc}:{i + 1}", "target-table",
                      f"{dom} says Target {target} minus Now {now}, which is {target - now}, and "
                      f"writes +{add}")
    if n == 0:
        rep.error(doc, "target-table",
                  "the target table has no rows; a check that reads nothing reports nothing")
        return

    for dom in sorted(held):
        if dom not in seen:
            rep.error(doc, "target-table",
                      f"the corpus holds {held[dom]} {dom} concepts and the table has no row")
    for dom in sorted(seen):
        if dom not in held:
            rep.error(doc, "target-table",
                      f"the table has a row for {dom}, which the corpus does not have")

    want = sum(t - w for w, t, _ in seen.values())
    claim = re.compile(r"The remaining per-domain targets add up to (\d[\d,]*)")
    found = 0
    for line_no, line in enumerate(lines, 1):
        for m in claim.finditer(line):
            found += 1
            if int(m.group(1).replace(",", "")) != want:
                rep.error(f"{doc}:{line_no}", "target-table",
                          f'"{m.group(0)}" but the targets minus the corpus leave {want}')
    if found == 0:
        rep.error(doc, "target-table",
                  "the sentence stating the remaining target total is gone; a check that reads "
                  "nothing reports nothing")


def defined_adrs() -> set:
    """Every ADR, read from docs/adr/ one file each.

    The ADRs lived inline in CONTEXT.md as one bullet per decision until the file reached
    140 KB. They are now one file per decision, so the source of truth is the directory: a file
    named ADR-NNN-*.md defines ADR-NNN. Reading the directory rather than a single file is what
    lets the check survive the split.
    """
    out = set()
    for f in sorted((ROOT / "docs/adr").glob("ADR-*.md")):
        m = re.match(r"(ADR-\d{3})-", f.name)
        if m:
            out.add(m.group(1))
    return out


def check_adr_refs(rep: Report) -> None:
    """Every ADR a doc cites must be defined in docs/adr/."""
    defined = defined_adrs()
    if not defined:
        rep.error("docs/adr", "adr", "no ADR files found; the pattern changed?")
        return
    for doc in DOCS:
        if doc == "CONTEXT.md":
            continue
        for i, line in enumerate(read(doc).split("\n"), 1):
            for m in re.finditer(r"\bADR-\d{3}\b", line):
                if m.group(0) not in defined:
                    rep.error(f"{doc}:{i}", "adr-ref",
                              f"{m.group(0)} is cited but not defined in docs/adr/")


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


# ---------------------------------------------------------------- schema and policy facts
#
# Each rule below reads its value from the code that owns it, for the reason in the module
# docstring: a checker that pins its own copy of the number goes stale in exactly the way it
# was written to prevent. Every one of these is a fact the docs already stated wrongly at some
# point, so each has already paid for itself once.

def app_source_files() -> list:
    """The next-app files that carry prose a reader sees.

    Diagram components hold their labels and detail strings in TypeScript, so a wrong algorithm
    or a removed formula can live there exactly as it can in markdown. No rule read them until
    the components were found describing a binary search and a deleted XP formula -- the same
    defects already fixed in the docs and the .d2 sources.
    """
    out = []
    for sub in ("components", "app"):
        out += sorted((ROOT / "web/next-app" / sub).rglob("*.tsx"))
        out += sorted((ROOT / "web/next-app" / sub).rglob("*.ts"))
    return [str(p.relative_to(ROOT)) for p in out if p.is_file()]


def read_any(rel: str) -> str:
    p = ROOT / rel
    return p.read_text(encoding="utf-8") if p.is_file() else ""


def go_roles() -> list[str]:
    """The four role values, read from the storage type that orders them."""
    m = re.search(r"const \(\s*(RoleStudent\s+Role\s*=\s*\"(\w+)\")", read("internal/storage/store.go"))
    src = read("internal/storage/store.go")
    return re.findall(r'Role\w+\s+Role\s+=\s*"(\w+)"', src)


def check_table_count(rep: Report) -> None:
    """"Five core tables" must match the CREATE TABLE statements in the migration."""
    src = read("internal/storage/migrate.go")
    tables = set(re.findall(r"CREATE TABLE IF NOT EXISTS (\w+)", src))
    if not tables:
        rep.error("internal/storage/migrate.go", "tables",
                  "no CREATE TABLE found; this checker cannot verify the docs")
        return
    pat = re.compile(r"\b(five|six|seven|eight|\d+)\s+(?:core\s+)?tables?\b", re.I)
    words = {"five": 5, "six": 6, "seven": 7, "eight": 8}
    for doc in DOCS:
        for line_no, sent in sentences(read(doc)):
            for m in pat.finditer(sent):
                got = words.get(m.group(1).lower(), None)
                if got is None:
                    got = int(m.group(1))
                if got != len(tables):
                    rep.error(f"{doc}:{line_no}", "tables",
                              f'"{m.group(0)}" but the migration creates {len(tables)}')


def check_roles(rep: Report) -> None:
    """A doc that enumerates the roles must enumerate the ones the code has."""
    roles = go_roles()
    if not roles:
        rep.error("internal/storage/store.go", "roles",
                  "no role constants found; this checker cannot verify the docs")
        return
    n = len(roles)
    # "admin|student" in a PATCH body is the two-role model the code left behind.
    pat = re.compile(r'"(admin)"\s*\|\s*"(student)"', re.I)
    for doc in DOCS:
        for line_no, line in enumerate(read(doc).split("\n"), 1):
            if pat.search(line):
                rep.error(f"{doc}:{line_no}", "roles",
                          f'a two-role body, but the code has {n}: {", ".join(roles)}')
    for doc in DOCS:
        for line_no, sent in sentences(read(doc)):
            m = re.search(r"\b(?:two|three|four|five)\s+roles?\b", sent, re.I)
            if m:
                got = {"two": 2, "three": 3, "four": 4, "five": 5}[m.group(1).lower()]
                if got != n:
                    rep.error(f"{doc}:{line_no}", "roles",
                              f'"{m.group(0)}" but the code has {n} roles')


def check_daily_goal(rep: Report) -> None:
    """A default daily XP goal in the docs must equal the constant the code owns."""
    goal = go_const("internal/xp/policy.go", "DefaultDailyGoal")
    if goal is None:
        rep.error("internal/xp/policy.go", "daily-goal",
                  "DefaultDailyGoal not found; this checker cannot verify the docs")
        return
    pat = re.compile(r"daily[_ ]xp[_ ]goal[^.\n]{0,60}?default[^.\n]{0,20}?(\d+)", re.I)
    for doc in DOCS:
        for line_no, sent in sentences(read(doc)):
            for m in pat.finditer(sent):
                if int(m.group(1)) != goal:
                    rep.error(f"{doc}:{line_no}", "daily-goal",
                              f'"{m.group(0)}" but xp.DefaultDailyGoal is {goal}')


def check_removed_xp_terms(rep: Report) -> None:
    """The time and streak multipliers were removed on purpose. No doc may state them.

    internal/xp/policy.go explains the removal and why, so a doc that reintroduces the formula
    contradicts a decision that is already recorded in the code. The XP gate check cannot see
    this: it matches "N XP" near a gate, and these terms are bare coefficients in a code fence.
    """
    banned = ("timeMultiplier", "streakMultiplier", "time multiplier", "streak multiplier",
              "time multipliers", "streak multipliers")
    # A doc may name a removed term while explaining the removal. The package comment in
    # internal/xp/policy.go does exactly this, and a rule that flagged it would push authors to
    # delete the only record of why the term is gone. Require the sentence to be about removal.
    about_removal = re.compile(
        r"\b(?:remov\w+|revert\w+|no longer|used to|previously|was |were |deleted|"
        r"deliberate|does not|not rewarded)\b", re.I)
    for doc in DOCS + app_source_files():
        # Check the sentence, not the line: these docs wrap at ~100 columns, so the removal
        # language ("the old time multiplier paid for guessing") sits on the next line from the
        # term itself. A line-scoped check flags correct prose for a formatting reason.
        for line_no, sent in sentences(read_any(doc)):
            for b in banned:
                if re.search(rf"\b{re.escape(b)}\b", sent, re.I) and not about_removal.search(sent):
                    rep.error(f"{doc}:{line_no}", "xp-terms",
                              f'"{b}" was removed; see the package comment in '
                              f"internal/xp/policy.go")


def check_grading_timeout(rep: Report) -> None:
    """The SymPy subprocess timeout must match the alarm the service raises."""
    m = re.search(r"signal\.alarm\((\d+)\)", read("grading/sympy_service.py"))
    if not m:
        rep.error("grading/sympy_service.py", "grading-timeout",
                  "no signal.alarm found; this checker cannot verify the docs")
        return
    secs = int(m.group(1))
    pat = re.compile(r"(\d+)[- ]second\s+(?:SIGALRM\s+)?timeout", re.I)
    for doc in DOCS:
        for line_no, sent in sentences(read(doc)):
            for hm in pat.finditer(sent):
                if int(hm.group(1)) != secs:
                    rep.error(f"{doc}:{line_no}", "grading-timeout",
                              f'"{hm.group(0)}" but sympy_service.py raises signal.alarm({secs})')


def check_diagnostic_bounds(rep: Report) -> None:
    """The diagnostic question bounds must match the constants the engine enforces.

    The in-app docs said "binary search ... approximately 20-35" while the code bounds an
    info-gain CAT at 25-45. The wrong figure was also visible to learners, because the
    onboarding briefing quoted the real one.
    """
    src = read("internal/diagnostic/cat.go")
    lo = re.search(r"minTotalQuestions\s*=\s*(\d+)", src)
    hi = re.search(r"maxTotalQuestions\s*=\s*(\d+)", src)
    if not (lo and hi):
        rep.error("internal/diagnostic/cat.go", "diagnostic-bounds",
                  "bounds not found; this checker cannot verify the docs")
        return
    lo_n, hi_n = int(lo.group(1)), int(hi.group(1))
    # A range written as "20-35", "20–35" or "20 to 35".
    rng = re.compile(r"\b(\d+)\s*(?:-|–|to)\s*(\d+)\s*questions?\b", re.I)
    ctx = re.compile(r"diagnostic|assessment|CAT|info[- ]gain|test", re.I)
    for doc in DOCS + ["web/next-app/app/docs/architecture/page.tsx",
                       "web/next-app/app/how-it-works/sections.tsx"]:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, line in enumerate(p.read_text(encoding="utf-8").split("\n"), 1):
            if not ctx.search(line):
                continue
            for m in rng.finditer(line):
                a, b = int(m.group(1)), int(m.group(2))
                if (a, b) != (lo_n, hi_n):
                    rep.error(f"{doc}:{line_no}", "diagnostic-bounds",
                              f'"{m.group(0)}" but cat.go bounds it at {lo_n}-{hi_n}')
    # "binary search" describes an algorithm the code does not use. Diagram sources are
    # included because they are documentation that regenerates into docs/diagrams/*.svg, and
    # a wrong claim there is published to every reader of the rendered diagram.
    algo_docs = DOCS + ["web/next-app/app/docs/architecture/page.tsx",
                        "web/next-app/app/how-it-works/sections.tsx",
                        "web/next-app/app/diagnose/page.tsx"]
    algo_docs += [f"web/next-app/diagrams/{p.name}" for p in sorted((ROOT / "web/next-app/diagrams").glob("*.d2"))]
    algo_docs += app_source_files()
    for doc in algo_docs:
        p = ROOT / doc
        if not p.is_file():
            continue
        text = p.read_text(encoding="utf-8")
        for line_no, sent in sentences(text):
            if not re.search(r"binary search", sent, re.I):
                continue
            # A doc or a comment may name the superseded algorithm while explaining what
            # replaced it. Without this, the only way to satisfy the rule would be to delete
            # the record of the change, which is the part worth keeping.
            if re.search(r"\b(?:replaced|not a?|used to|previously|no longer|superseded|"
                         r"describes|instead|rather than)\b", sent, re.I):
                continue
            rep.error(f"{doc}:{line_no}", "diagnostic-algo",
                      '"binary search" but cat.go runs an info-gain CAT over a covering set')


def check_adr_headers(rep: Report) -> None:
    """Every ADR a source file or a doc cites must have a file in docs/adr/.

    check_adr_refs scans markdown only. ADR-050 had no header and was cited by three Go files,
    so nothing failed: the decision existed, was load-bearing, and was unrecorded as a decision.
    This walks internal/ as well.
    """
    defined = defined_adrs()
    if not defined:
        rep.error("docs/adr", "adr", "no ADR files found; the pattern changed?")
        return
    for f in sorted((ROOT / "internal").rglob("*.go")):
        if f.name.endswith("_test.go"):
            continue
        for line_no, line in enumerate(f.read_text(encoding="utf-8").split("\n"), 1):
            for m in re.finditer(r"\bADR-\d{3}\b", line):
                if m.group(0) not in defined:
                    rel = f.relative_to(ROOT)
                    rep.error(f"{rel}:{line_no}", "adr-ref",
                              f"{m.group(0)} is cited in code but not defined in docs/adr/")


def git_ignored(paths: list) -> set:
    """Return the subset of paths that git ignores.

    A citation resolves only if the file is in the repository. A path under a gitignored tree
    exists on a developer's machine and nowhere else, so a check that reads the filesystem
    passes locally and fails on a clean clone -- which is what a citation check must not do.
    CONTEXT.md cited `skills/skills/engineering/tdd/SKILL.md`, and `skills/` is ignored, so the
    check passed here and failed in CI. Consulting git closes the gap.
    """
    if not paths:
        return set()
    try:
        import subprocess
        r = subprocess.run(["git", "check-ignore", "--stdin"], cwd=ROOT,
                           input="\n".join(str(p) for p in paths),
                           capture_output=True, text=True, timeout=30)
        return {line for line in r.stdout.splitlines() if line}
    except Exception:
        # git absent or the call failed: fall back to the filesystem alone rather than
        # reporting every citation as broken.
        return set()


def file_index() -> dict:
    """Map a basename to every file with that name, walking the tree once.

    The citation check looks up files by basename, because the docs cite "engine.go" and
    "main.go" with no directory. A rglob per citation cost about eight seconds across the
    corpus; one walk costs a fraction of that.
    """
    global _FILE_INDEX
    if _FILE_INDEX is not None:
        return _FILE_INDEX
    skip = {"node_modules", ".next", ".git", "dist", "build"}
    idx: dict = {}
    for dirpath, dirnames, filenames in os.walk(ROOT):
        dirnames[:] = [d for d in dirnames if d not in skip]
        for fn in filenames:
            idx.setdefault(fn, []).append(Path(dirpath) / fn)
    _FILE_INDEX = idx
    return idx


def check_symbol_citations(rep: Report) -> None:
    """A `file.go:LINE` citation must point at a line that still exists.

    The docs cite source locations, and the locations rot: a short file grows, the line numbers
    do not move, and the citation quietly names unrelated code. No rule read them, because
    every other check in this file looks at the claim rather than the pointer.
    """
    pat = re.compile(r"`([\w./-]+\.(?:go|py|tsx|ts|mjs|md)):(\d+)(?:-(\d+))?`")
    idx = file_index()

    # First pass: resolve every citation to the file it means, then ask git once which of
    # those files are ignored. Calling git per citation cost about ten seconds; one call over
    # the whole set costs one.
    resolved: list[tuple] = []
    for doc in DOCS:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, line in enumerate(p.read_text(encoding="utf-8").split("\n"), 1):
            for m in pat.finditer(line):
                rel, start = m.group(1), int(m.group(2))
                # The exact cited path, then the same basename anywhere in the tree.
                cands = [ROOT / rel] + list(idx.get(rel.split("/")[-1], []))
                # A gitignored file is not in the repository, so a citation to it does not
                # resolve for anyone who clones. Filtered below so local and CI agree.
                resolved.append((doc, line_no, rel, start, m.group(3), cands))

    present = [c for _, _, _, _, _, cs in resolved for c in cs if c.is_file()]
    ignored = git_ignored(present)

    for doc, line_no, rel, start, end_s, cands in resolved:
        cands = [c for c in cands if str(c) not in ignored]
        target = next((c for c in cands if c.is_file()), None)
        if target is None:
            rep.error(f"{doc}:{line_no}", "citation", f"`{rel}` does not exist")
            continue
        total = len(target.read_text(encoding="utf-8").split("\n"))
        end = int(end_s or start)
        if end > total:
            rep.error(f"{doc}:{line_no}", "citation",
                      f"`{rel}:{start}` is past the end of the file ({total} lines)")
            continue
        # The shape test below is for code. A citation into a markdown file is checked for
        # existence only: a blank line in prose is normal, and the rule exists to catch a code
        # anchor that slid onto whitespace.
        if rel.endswith(".md"):
            continue
        # A citation that lands on a blank line, a closing brace, or an import names nothing.
        # These are the three shapes the line numbers drift into: a file grows and the anchor
        # slides onto whitespace or into a neighbouring block. It is a mechanical test, so it
        # costs nothing and needs no knowledge of the symbol.
        content = target.read_text(encoding="utf-8").split("\n")[start - 1]
        if not content.strip() or re.match(r'^\s*[\}\)\]],?\s*$', content) \
                or re.match(r'\s*(?:import|"[\w./-]+")\s', content):
            rep.error(f"{doc}:{line_no}", "citation",
                      f"`{rel}:{start}` lands on {content.strip()[:30] or 'a blank line'!r}"
                      f", which names no symbol")


def check_no_stub_claims(rep: Report) -> None:
    """A doc may not call an implemented store a stub, or promise a fallback that is absent.

    Three claims rotted together. PostgresStore was described as "a stub with all methods
    returning 'not implemented'" in the in-app docs while postgres.go carried 2,141 lines and
    no such string. The grading diagram drew a "Symbolic Go fallback" that router.go does not
    branch to, and the endpoint count said 25+ against 88 registered routes. All three are
    reads of the code, so they are reads of the code here too.
    """
    pg = read("internal/storage/postgres.go")
    if re.search(r"not implemented", pg, re.I):
        rep.error("internal/storage/postgres.go", "postgres-stub",
                  "the doc says this store is a stub and the file now says 'not implemented'")
    routes = len(re.findall(r"mux\.HandleFunc", read("internal/server/server.go")))
    if routes == 0:
        rep.error("internal/server/server.go", "routes",
                  "no mux.HandleFunc found; this checker cannot verify the docs")
        return
    ctx = re.compile(r"REST endpoints?|routes?|API exposes", re.I)
    loose = re.compile(r"\b(\d+)\+?\s+(?:REST\s+)?(?:endpoints?|routes?)\b", re.I)
    stub_ctx = re.compile(r"stub|not implemented|awaiting implementation", re.I)
    for doc in DOCS + ["web/next-app/app/docs/architecture/page.tsx",
                       "web/next-app/components/SystemDesignFlow/ArchitectureFlow.tsx",
                       "web/next-app/components/GradingFlow.tsx"]:
        p = ROOT / doc
        if not p.is_file():
            continue
        # Paragraph-scoped: these pages wrap prose across lines, so "postgres.go" and
        # "a stub awaiting implementation" routinely land on different lines. A line-scoped
        # check found one of the two claims and let the other through.
        # The claim has two parts and both must be about PostgresStore in the same sentence:
        # "the PostgresStore is a stub" / "currently a stub awaiting implementation". Matching
        # the pair anywhere in a paragraph produced a false positive in CONTEXT.md, where the
        # whole ADR block is one paragraph and an unrelated line supplies the other half.
        text = p.read_text(encoding="utf-8")
        claim = re.compile(
            r"(?:Postgres\w*|PostgreSQL\w*)\s+(?:is|are)\s+(?:\w+\s+){0,4}?"
            r"(?:a\s+)?(?:stub|not implemented|awaiting implementation)", re.I)
        for line_no, line in enumerate(text.split("\n"), 1):
            for sent in re.split(r"(?<=[.!?:])\s+", line):
                if claim.search(sent):
                    rep.error(f"{doc}:{line_no}", "postgres-stub",
                              "postgres.go is implemented and contains no 'not implemented'")
        for line_no, line in enumerate(text.split("\n"), 1):
            if ctx.search(line):
                for m in loose.finditer(line):
                    n = int(m.group(1))
                    # A "+" floor is not automatically wrong -- "25+" is true of 88. But a
                    # floor that understates by more than half reads as a total, and that is
                    # what a reader takes from it. Require the floor to be within a factor of
                    # the real count, so the number has to be maintained.
                    if "+" in m.group(0) and n * 2 < routes:
                        rep.error(f"{doc}:{line_no}", "routes",
                                  f'"{m.group(0)}" understates the {routes} routes the router '
                                  f"registers")
                    elif "+" not in m.group(0) and n != routes:
                        rep.error(f"{doc}:{line_no}", "routes",
                                  f'"{m.group(0)}" but the router registers {routes}')
    # A grading fallback the router does not branch to.
    if "symbolicGrader" not in read("internal/grader/router.go") + read("internal/grader/grader.go"):
        for doc in ["web/next-app/components/GradingFlow.tsx",
                    "web/next-app/components/SystemDesignFlow/ArchitectureFlow.tsx"]:
            p = ROOT / doc
            if not p.is_file():
                continue
            for line_no, line in enumerate(p.read_text(encoding="utf-8").split("\n"), 1):
                if re.search(r"symbolic.{0,20}fallback|fallback.{0,20}symbolic", line, re.I):
                    rep.error(f"{doc}:{line_no}", "grading-fallback",
                              "the router has no pure-Go symbolic fallback branch")


def check_toolchain(rep: Report) -> None:
    """The README's Go badge must agree with go.mod.

    The badge said 1.21+ while the same file's requirements line said 1.26+ and go.mod
    targeted 1.26.6. A badge is the first thing a prospective contributor reads, and a wrong
    one sends them to install a toolchain the module will not build under.
    """
    m = re.search(r"^go\s+([0-9.]+)", read("go.mod"), re.M)
    if not m:
        rep.error("go.mod", "go-version", "no go directive found")
        return
    floor = ".".join(m.group(1).split(".")[:2])  # 1.26.6 -> 1.26
    for doc in ("README.md", "CONTRIBUTING.md"):
        for line_no, line in enumerate(read(doc).split("\n"), 1):
            # shields.io encodes the "+" as %2B, but the badge may carry either form.
            for hm in re.finditer(r"badge/Go-([0-9.]+)(?:\+|%2B)", line):
                if hm.group(1) != floor:
                    rep.error(f"{doc}:{line_no}", "go-version",
                              f'badge says Go {hm.group(1)}+ but go.mod targets {floor}+')
            for hm in re.finditer(r"Go\s+([0-9.]+)\+", line):
                if hm.group(1) != floor:
                    rep.error(f"{doc}:{line_no}", "go-version",
                              f'"{hm.group(0)}" but go.mod targets {floor}+')

    # Node: .mise.toml owns the version. .nvmrc, the CI setup-node pins and the
    # Dockerfile image must carry the same string, and a doc that says "Node 24"
    # must name that major. The web build scripts import TypeScript files
    # straight into Node, so a contributor on the wrong major gets
    # ERR_UNKNOWN_FILE_EXTENSION three commands after a clean install.
    nm = re.search(r'^node\s*=\s*"([0-9.]+)"', read(".mise.toml"), re.M)
    if not nm:
        rep.error(".mise.toml", "node-version", "no node pin found")
        return
    node = nm.group(1)
    major = node.split(".")[0]
    nvmrc = read(".nvmrc").strip()
    if nvmrc != node:
        rep.error(".nvmrc", "node-version", f".nvmrc pins {nvmrc!r} but .mise.toml pins {node!r}")
    for line_no, line in enumerate(read(".github/workflows/ci.yml").split("\n"), 1):
        for hm in re.finditer(r"node-version:\s*'([0-9.]+)'", line):
            if hm.group(1) != node:
                rep.error(f".github/workflows/ci.yml:{line_no}", "node-version",
                          f"CI installs Node {hm.group(1)} but .mise.toml pins {node}")
    for line_no, line in enumerate(read("Dockerfile").split("\n"), 1):
        for hm in re.finditer(r"^FROM node:([0-9.]+)", line):
            if hm.group(1) != node:
                rep.error(f"Dockerfile:{line_no}", "node-version",
                          f"the image builds on Node {hm.group(1)} but .mise.toml pins {node}")
    for doc in ("README.md", "CONTRIBUTING.md"):
        for line_no, line in enumerate(read(doc).split("\n"), 1):
            for hm in re.finditer(r"\bNode\s+(\d+)\b", line):
                if hm.group(1) != major:
                    rep.error(f"{doc}:{line_no}", "node-version",
                              f'"{hm.group(0)}" but .mise.toml pins Node {major}')


def check_sm2_intervals(rep: Report) -> None:
    """A quoted first-or-second SM-2 interval must match the code.

    README said "1 day -> 3 days -> 1 week". sm2.go uses 1 day then 6 days, and every interval
    is additionally scaled by the learner's per-topic speed, which the same sentence did not
    mention. system-design.md had the 6 right, so the two documents disagreed.
    """
    src = read("internal/scheduler/sm2.go")
    first = re.search(r"case 0:\s*\n\s*interval = int\(math\.Round\((\d+)", src)
    second = re.search(r"case 1:\s*\n\s*interval = int\(math\.Round\((\d+)", src)
    if not (first and second):
        rep.error("internal/scheduler/sm2.go", "sm2",
                  "interval constants not found; this checker cannot verify the docs")
        return
    f_n, s_n = first.group(1), second.group(1)
    pat = re.compile(r"\b1 day\s*(?:→|->|→|to)\s*(\d+)\s*days?\b", re.I)
    for doc in DOCS:
        for line_no, line in enumerate(read(doc).split("\n"), 1):
            for m in pat.finditer(line):
                if m.group(1) != s_n:
                    rep.error(f"{doc}:{line_no}", "sm2",
                              f'"1 day -> {m.group(1)} days" but sm2.go uses {s_n} for rep 1')


def check_diagram_citations(rep: Report) -> None:
    """A cited diagram must exist, including the ones cited from a Go doc comment.

    check_links only sees markdown `[text](target)`. Three docs and internal/mastery/machine.go
    cited `docs/diagrams/student-model.svg`, which does not exist -- the files are
    student-model-dark.svg and student-model-light.svg -- and no rule noticed, because the
    citation is in backticks, not a link.
    """
    pat = re.compile(r"`((?:docs|web/next-app/public)/diagrams/[\w.-]+\.svg)`")
    for doc in DOCS + [str(f.relative_to(ROOT)) for f in
                       sorted((ROOT / "internal").rglob("*.go")) if not f.name.endswith("_test.go")]:
        p = ROOT / doc
        if not p.is_file():
            continue
        for line_no, line in enumerate(p.read_text(encoding="utf-8").split("\n"), 1):
            for m in pat.finditer(line):
                rel = m.group(1)
                if not (ROOT / rel).is_file():
                    # The docs are written from the repo root; a path under web/next-app/public
                    # is the same file mirrored. Accept either root.
                    alt = ROOT / "web/next-app/public" / rel.split("diagrams/", 1)[1]
                    if not alt.is_file():
                        rep.error(f"{doc}:{line_no}", "diagram",
                                  f"`{rel}` does not exist")


def check_design_tokens(rep: Report) -> None:
    """Every token DESIGN.md documents must exist in globals.css.

    DESIGN.md described a hex palette with `--accent-gold` while globals.css had moved to OKLCH
    and had no such token. The stale palette had also propagated into two other docs before the
    move. Nothing ever read the design doc, so nothing noticed. This reads both.
    """
    # Font variables are declared by next/font in app/layout.tsx, not in CSS, so both files
    # are sources of a legitimate token name.
    css = read("web/next-app/app/globals.css") + read("web/next-app/app/layout.tsx")
    defined = set(re.findall(r"(--[\w-]+)\s*:", css)) | set(re.findall(r"variable:\s*'(--[\w-]+)'", css))
    if not defined:
        rep.error("web/next-app/app/globals.css", "design-tokens",
                  "no custom properties found; this checker cannot verify DESIGN.md")
        return
    doc = ROOT / "DESIGN.md"
    if not doc.is_file():
        return
    for line_no, line in enumerate(doc.read_text(encoding="utf-8").split("\n"), 1):
        # A token named in a table cell or inline code, not a CSS snippet's own declaration.
        for m in re.finditer(r"`(--[\w-]+)`", line):
            if m.group(1) not in defined:
                rep.error(f"DESIGN.md:{line_no}", "design-tokens",
                          f"`{m.group(1)}` is documented but not defined in globals.css")


def main() -> int:
    rep = Report()
    check_xp_gate(rep)
    check_kp_rule(rep)
    check_corpus_figures(rep)
    check_generator_registries(rep)
    check_course_reach(rep)
    check_media_candidates(rep)
    check_figure_dispositions(rep)
    check_thin_lessons(rep)
    check_target_table(rep)
    check_adr_refs(rep)
    check_links(rep)
    check_spelling(rep)
    check_table_count(rep)
    check_roles(rep)
    check_daily_goal(rep)
    check_removed_xp_terms(rep)
    check_grading_timeout(rep)
    check_diagnostic_bounds(rep)
    check_adr_headers(rep)
    check_symbol_citations(rep)
    check_no_stub_claims(rep)
    check_toolchain(rep)
    check_sm2_intervals(rep)
    check_diagram_citations(rep)
    check_design_tokens(rep)

    for where, rule, msg in rep.errors:
        print(f"ERROR   {where} [{rule}] {msg}", file=sys.stderr)
    if rep.errors:
        print(f"\nFAIL: {len(rep.errors)} documentation inconsistency/inconsistencies.",
              file=sys.stderr)
        return 1
    print("OK: docs agree with the code and the corpus "
          "(XP gate, KP rule, corpus counts, course reach, media candidates, figure "
          "dispositions, thin lessons, target table, "
          "ADR refs, links, spelling, tables, "
          "roles, daily goal, removed XP terms, grading timeout, diagnostic bounds, "
          "citations, stub claims, route counts, toolchain, SM-2 intervals, "
          "diagram citations, design tokens)")
    return 0


if __name__ == "__main__":
    sys.exit(main())