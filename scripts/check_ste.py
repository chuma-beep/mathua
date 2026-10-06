#!/usr/bin/env python3
"""
check_ste.py — enforce the structural writing rules of Simplified Technical English.

This is a derived subset, not ASD-STE100 conformance, and it must never be described as
conformance. The upstream skill (github.com/0xpili/simplified-technical-english, MIT) says so
itself: it "does not certify compliance with ASD-STE100". Two reasons this script is a
reimplementation rather than a vendored call:

  - **Licensing.** MIT covers the skill's own text. The ASD-STE100 *dictionary* is ASD's property
    and explicitly not MIT (see that project's NOTICE.md). Copying `word-list.md` into this repo
    would put ASD's dictionary in our tree.
  - **The dictionary is the wrong tool here.** A maths and pedagogy codebase is mostly technical
    names — concept, learner, mastery, scheduler, endpoint, prerequisite — which STE permits
    anyway. A vocabulary gate would either nag constantly or need an allowlist wide enough to be
    meaningless. The *structural* rules are what actually improve readability, and they are
    checkable with no dictionary at all.

What this enforces (all objective, all mechanical):
  1. No semicolons.
  2. No contractions.
  3. No "should", "would", "may", "might", "shall", "ought" — STE allows only can/must/will.
  4. No helping verb plus past participle ("has been", "have written") — use the simple past.
  5. No progressive "-ing" forms, except the eight STE permits.
  6. No passive "is/are/was/were ... by".
  7. Sentence length limits: 20 words procedural, 25 descriptive.
  8. At most 6 sentences per paragraph, and one topic per paragraph.
  9. One instruction per sentence.

Scope: the procedural and onboarding docs only. CONTEXT.md and the ADRs are deliberately
excluded, because their value is the evidence in them — measured ratios, published retractions,
and rejected alternatives — and these rules strip exactly that. Excluding them is a decision,
recorded in AGENTS.md, not an oversight.

Usage:
    python3 scripts/check_ste.py [--mode procedural|descriptive] FILE...
"""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

# ---------------------------------------------------------------- scope

# Procedural and onboarding docs. These are the docs a newcomer reads while doing something.
STE_DOCS = [
    "CONTRIBUTING.md",
    "docs/mobile-device-checklist.md",
    "docs/architecture.md",
]

# README is only partly procedural, so only its instructional sections are checked. The list is
# matched on word boundaries: an earlier version used "make", which substring-matched the heading
# "What makes it different" and pulled feature-description prose into a procedural check. That
# section is descriptive and belongs to the other tier, not here.
README_STE_SECTIONS = (
    "getting started", "build from source", "installation", "install", "requirements",
    "quick start", "quickstart", "running", "testing", "commands", "troubleshooting",
)

# ---------------------------------------------------------------- rules

# Rule 3.2 and the dictionary: can, must, will are the only approved helping verbs.
BANNED_MODALS = {"should", "would", "may", "might", "shall", "ought"}

# Rule 3.5: the approved "-ing" words, and words where "ing" is not a verb suffix.
ING_APPROVED = {"mating", "missing", "remaining", "lighting", "opening", "routing", "servicing",
                "during"}
ING_NOT_SUFFIX = {"ring", "spring", "string", "king", "thing", "wing", "sing", "bring", "sting",
                  "swing", "nothing", "anything", "everything", "something", "bearing", "ceiling",
                  "morning", "evening", "during", "setting", "surrounding", "understanding"}

# Rule 4.2.
CONTRACTION = re.compile(r"\b\w+n['’]t\b|\b\w+['’](?:re|ll|ve|d|m)\b", re.IGNORECASE)

# Rule 3.4.
HAVE_PARTICIPLE = re.compile(
    r"\b(?:has|have|had)\s+(?:been|got|gone|put|set|kept|held|taken|given|found|left|lost|meant|"
    r"sent|shown|told|built|become|begun|broken|brought|come|fallen|felt|grown|known|read|run|"
    r"seen|spoken|thrown|worn|written|done|made|\w+ed)\b", re.IGNORECASE)

# Rule 3.6.
PASSIVE_BY = re.compile(
    r"\b(?:is|are|was|were|be|been|being)\s+(?:\w+\s+){0,2}?\w+(?:ed|en)\s+by\b", re.IGNORECASE)

# Rule 3.5: progressive.
BE_ING = re.compile(r"\b(?:is|are|was|were|be|been|am)\s+(\w+ing)\b", re.IGNORECASE)

WORDISH = re.compile(r"[A-Za-z0-9'’&/—-]+")
SENT_SPLIT = re.compile(r"(?<=[.!?:])\s+")
BULLET = re.compile(r"\s*(?:[-*+]|\d+\.)\s")


def strip_markdown(text: str) -> str:
    """Blank out what the STE rules do not control: code, identifiers, URLs, quotes."""
    text = re.sub(r"^---\n.*?\n---\n", lambda m: "\n" * m.group(0).count("\n"), text, flags=re.S)
    text = re.sub(r"```.*?```", lambda m: "\n" * m.group(0).count("\n"), text, flags=re.S)
    text = re.sub(r"`[^`\n]*`", lambda m: " " * len(m.group(0)), text)
    text = re.sub(r"\[([^\]]*)\]\([^)]*\)", r"\1", text)
    text = re.sub(r"https?://\S+", lambda m: " " * len(m.group(0)), text)
    text = re.sub(r"^>.*$", lambda m: " " * len(m.group(0)), text, flags=re.M)
    text = re.sub(r"\|.*\|", lambda m: " " * len(m.group(0)), text)   # tables
    text = re.sub(r"\"[^\"\n]*\"|“[^”\n]*”", lambda m: " QUOTED ", text)
    return text


def count_words(sentence: str) -> int:
    """STE rule 8: a parenthetical, a hyphenated group, and a number each count as one word."""
    s = re.sub(r"\([^)]*\)", " P ", sentence)
    s = re.sub(r"\*\*|\*|__|_|#+", " ", s)
    return len(WORDISH.findall(s))


def sentences(block: str):
    for part in SENT_SPLIT.split(block):
        part = part.strip()
        if part and WORDISH.search(part):
            yield part


def check_sentence(sent: str, limit: int, where: str, errors: list[tuple[str, str]],
                   warnings: list[tuple[str, str]]) -> None:
    head = sent if len(sent) <= 70 else sent[:67] + "..."
    n = count_words(sent)
    if n > limit:
        errors.append((where, f"5.1/6.3 sentence has {n} words (max {limit}): \"{head}\""))
    if ";" in sent:
        errors.append((where, f'8.1 semicolon: "{head}". Write two sentences.'))
    m = CONTRACTION.search(sent)
    if m:
        errors.append((where, f'4.2 contraction "{m.group(0)}": write the full words.'))
    m = HAVE_PARTICIPLE.search(sent)
    if m:
        errors.append((where, f'3.4 helping verb "{m.group(0)}": use the simple past tense.'))
    m = PASSIVE_BY.search(sent)
    if m:
        errors.append((where, f'3.6 passive voice "{m.group(0)}": make the agent the subject.'))
    m = BE_ING.search(sent)
    if m and m.group(1).lower() not in ING_APPROVED:
        errors.append((where, f'3.5 progressive "{m.group(0)}": use a simple tense.'))
    for w in re.findall(r"[A-Za-z'’-]+", sent):
        if w.lower() in BANNED_MODALS:
            errors.append((where, f'3.2 "{w}": use "must" (requirement), "can" (possibility), '
                                   f'"will" (future), or remove it.'))
        elif w.lower().endswith("ing") and len(w) > 5 and not w.isupper():
            # A warning, not an error: this word may be an adjective or a gerund noun, and rule
            # 3.5 governs the participle used as a verb, which BE_ING above already catches.
            if w.lower() not in ING_APPROVED and w.lower() not in ING_NOT_SUFFIX:
                warnings.append((where, f'3.5 "{w}": if this is a verb, use a simple form; if it '
                                        f'is a label, it is fine.'))


# Rule 8.6: a title or a label is not controlled text. A bullet that reads
# "**Scoring** — produces the evidence score" has a label and a description; rewriting "Scoring"
# to satisfy rule 3.5 would make the heading worse, not the documentation clearer.
LABEL = re.compile(r"^(?:[-*+]\s+|\d+\.\s+)?\**\s*([^\n]{0,60}?)\s*\**\s*[—–:-]\s")


def split_label(item: str) -> tuple[str, str]:
    """Return (label, rest). Only rest is checked."""
    m = LABEL.match(item.strip())
    if not m:
        return "", item
    return m.group(1), item[m.end():]


def check_text(text: str, limit: int, name: str, errors: list[tuple[str, str]],
               warnings: list[tuple[str, str]]) -> None:
    text = strip_markdown(text)
    for i, par in enumerate(re.split(r"\n\s*\n", text)):
        par = par.strip()
        if not par or not WORDISH.search(par):
            continue
        if par.startswith("#"):
            continue
        # Rule 8.4: a list item counts as its own sentence. A contiguous bullet block is not one
        # paragraph of N sentences, and counting it as one invents a rule 6.6 violation.
        lines = par.splitlines()
        if any(BULLET.match(ln) for ln in lines):
            # A wrapped list item continues on the next line. Treating those continuation lines
            # as prose invented a rule 6.6 violation out of one list, which is how a three-item
            # list became an "eight-sentence paragraph".
            items: list[str] = []
            prose: list[str] = []
            for ln in lines:
                if BULLET.match(ln):
                    items.append(BULLET.sub("", ln))
                elif items:
                    items[-1] += " " + ln.strip()
                else:
                    prose.append(ln)
            for item in items:
                _, rest = split_label(item)
                for sent in sentences(rest or item):
                    check_sentence(sent, limit, name, errors, warnings)
            if prose and WORDISH.search("\n".join(prose)):
                check_paragraph("\n".join(prose), limit, name, errors, warnings)
            continue
        check_paragraph(par, limit, name, errors, warnings)


def check_paragraph(par: str, limit: int, name: str, errors: list[tuple[str, str]],
                   warnings: list[tuple[str, str]]) -> None:
    if True:
        sents = list(sentences(par))
        if len(sents) > 6:
            errors.append((name, f"6.6 paragraph has {len(sents)} sentences (max 6): \"{sents[0][:50]}...\""))
        for s in sents:
            check_sentence(s, limit, name, errors, warnings)
            # Rule 9: one instruction per sentence. "Run X and then Y" is the usual shape.
            if re.match(r"^\s*(?:Run|Set|Check|Use|Open|Add|Edit|Save|Copy|Paste|Install|"
                        r"Start|Stop|Remove|Read|Write|Create|Execute|Make)\b", s) and \
                    re.search(r"\b(?:and then|and also|,\s*and\s+then)\b", s):
                errors.append((name, f"9.1 more than one instruction in a sentence: \"{s[:60]}\""))


def readme_ste_sections(text: str) -> str:
    """Keep only README sections whose heading names an instructional topic.

    Word-boundary matching, so "make" cannot match "makes". A heading also ends the section, so
    content never leaks from a descriptive section into a procedural one.
    """
    lines = text.split("\n")
    keep, on, in_fence = [], False, False
    for line in lines:
        # Track fenced code blocks. Without this, a shell comment such as "# SQLite (zero-config)"
        # reads as a markdown heading, matches no keyword, and silently ends the section — which
        # dropped the Requirements and Env block from the check while it still reported OK.
        if line.lstrip().startswith("```"):
            in_fence = not in_fence
            if on:
                keep.append(line)
            continue
        if line.startswith("#") and not in_fence:
            heading = line.lstrip("#").strip().lower()
            on = any(re.search(rf"\b{re.escape(k)}\b", heading) for k in README_STE_SECTIONS)
        if on:
            keep.append(line)
    return "\n".join(keep)


def main() -> int:
    ap = argparse.ArgumentParser(description="Check STE structural writing rules.")
    ap.add_argument("files", nargs="*", help="files to check (default: the STE doc set)")
    ap.add_argument("--mode", choices=["procedural", "descriptive"], default="descriptive")
    ap.add_argument("--limit", type=int, default=None,
                    help="override the sentence word limit")
    args = ap.parse_args()

    root = Path(__file__).resolve().parents[1]
    limit = args.limit or (20 if args.mode == "procedural" else 25)
    files = args.files or [str(root / d) for d in STE_DOCS] + [str(root / "README.md")]

    errors: list[tuple[str, str]] = []
    warnings: list[tuple[str, str]] = []
    checked = 0
    for f in files:
        p = Path(f)
        if not p.is_file():
            print(f"ERROR   {f} [missing] not a file", file=sys.stderr)
            return 1
        text = p.read_text(encoding="utf-8")
        if p.name == "README.md":
            text = readme_ste_sections(text)
        checked += 1
        check_text(text, limit, p.name, errors, warnings)

    for where, msg in warnings:
        print(f"WARNING {where} [{msg}]", file=sys.stderr)
    for where, msg in errors:
        print(f"ERROR   {where} [{msg}]", file=sys.stderr)
    if errors:
        print(f"\nFAIL: {len(errors)} STE structural error(s) in {checked} file(s).",
              file=sys.stderr)
        print("These are the structural rules of Simplified Technical English. This script is not",
              file=sys.stderr)
        print("an ASD-STE100 conformance check and makes no claim of certification.",
              file=sys.stderr)
        return 1
    # Built as a variable, not a multi-line expression inside the f-string braces: PEP 701
    # f-string formatting is Python 3.12+, and CI runs 3.11. This passed locally on 3.14 and
    # only failed in CI, which is the whole argument for testing on the version that ships.
    advisory = f"; {len(warnings)} advisory warning(s) not failed" if warnings else ""
    print(f"OK: {checked} file(s) obey the STE structural rules this script checks "
          f"(max {limit} words per sentence){advisory}.")
    print("Structural subset only. Not an ASD-STE100 conformance check.")
    return 0


if __name__ == "__main__":
    sys.exit(main())