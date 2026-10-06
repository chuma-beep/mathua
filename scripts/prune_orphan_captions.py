#!/usr/bin/env python3
"""Remove figure captions that nothing refers to and whose lesson never cites a figure.

An orphan caption is a paragraph that is nothing but emphasis and that sits next to no image —
the shape `ingest_pretext.mjs` emits for a `<figure>` when it keeps the caption and drops the
image. Left in place it shows a learner a description of a picture that is not there.

**The delete set is much smaller than the candidate set, and the reason is the interesting
part.** Of the 359 candidates the audit found:

    220  sit in lessons whose prose cites a figure — by number ("as shown in Figure 1.2") or by
         position ("the graph above illustrates"). Deleting those would leave the prose
         pointing at nothing, which is a worse defect than an orphaned caption. Kept.
     35  sit in lessons that cite no figure at all, read positively as captions, and carry
         descriptive text. Deleted.
     104 the rest, and none of them are safe to remove without reading the lesson:
         74 have no picture vocabulary, so they could as easily be emphasis; 30 are bare
         section labels such as *Solution*, which head a worked answer and are not captions at
         all.

The first dry run of this script deleted 210 lines and was wrong. It removed `*Figure 1.2 …*`
from a lesson whose prose says "visualized on a number line as shown in Figure 1.2", because
the prose check only recognised relative references like "the figure below". It also removed
`*Solution*` — a heading — because a lesson with two worked answers has two of them and the
audit's "duplicated in this lesson" signal promoted them to `strong`. Both are why the rule
here requires positive evidence twice over: the lesson must cite no figure, *and* the caption
must positively read as one.

The figures are unrecoverable, which is what makes this the honest end state rather than a
convenience: the ORCCA image files are 404 on every branch, so a caption cannot be paired with
its picture by going back to the source. Of 359 candidates, 3 matched anything still on disk,
and all 3 matched the same figure from three sibling lessons.

Usage:
    python3 scripts/prune_orphan_captions.py --dry-run
    python3 scripts/prune_orphan_captions.py
"""
import argparse
import os
import re
import sys

# The candidate set and the shape rules come from the audit rather than a second classifier
# here. This script's own version was looser and disagreed with the reviewed report by 74, and
# the extra 74 included content that had to stay.
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from audit_lesson_media import (  # noqa: E402
    section_a_candidates, FIGURE_REF_PROSE, FIGURE_MARKERS, NUMBERED_CAPTION_RE,
)

# A lesson that cites a figure by number. "as shown in Figure 1.2" is just as much a reference
# as "the figure below", and missing the first form is what made the first dry run unsafe.
CITES_FIGURE_NUMBER = re.compile(r"\b(Figure|Table|Graph|Diagram)\s+\d")

# A bare single-word marker is a heading in this corpus far more often than a caption.
SECTION_LABELS = re.compile(
    r"(?i)^(solution|solutions|answer|answers|hint|hints|note|notes|example|examples|"
    r"try it|summary|proof|proofs|interactive|method|explanation|review|homework|"
    r"alternative video lessons?)$")


def reads_as_a_caption(text):
    """Positive evidence that this describes a picture.

    An auto-numbered `Figure 3.2` / `Table 1.4` label, or vocabulary about pictures. Absence of
    counter-evidence is not evidence: `*Share of all income held by the top 1%, United States*`
    is a real caption that mentions nothing visual, and a rule that discarded it would be
    discarding content.
    """
    return bool(NUMBERED_CAPTION_RE.match(text) or FIGURE_MARKERS.search(text))


def caption_is_bare_label(text):
    """`Figure 1.3` and nothing else — the number, with no description attached."""
    m = NUMBERED_CAPTION_RE.match(text)
    if not m:
        return False
    return not text[m.end():].strip(" .:—")


def main():
    ap = argparse.ArgumentParser(description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--dry-run", action="store_true", help="report what would go, write nothing")
    args = ap.parse_args()

    by_lesson = {}
    for row in section_a_candidates():
        by_lesson.setdefault(row["lesson"], []).append(row)

    deleted = kept = 0
    reasons = {"lesson cites a figure": 0, "no picture signal": 0,
               "bare section label": 0, "bare figure label": 0}
    touched = []

    for rel in sorted(by_lesson):
        rows = by_lesson[rel]
        path = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
                            "data", "lessons", rel)
        with open(path, encoding="utf-8") as fh:
            body = fh.read()

        cites = bool(FIGURE_REF_PROSE.search(body) or CITES_FIGURE_NUMBER.search(body))

        lines = body.split("\n")
        drop = set()
        for row in rows:
            text = row["text"]
            label = re.sub(r"\s+", " ", text).strip()
            if cites:
                reason = "lesson cites a figure"
            elif SECTION_LABELS.match(label):
                reason = "bare section label"
            elif not reads_as_a_caption(text):
                reason = "no picture signal"
            elif caption_is_bare_label(text):
                reason = "bare figure label"
            else:
                reason = ""
            if reason:
                kept += 1
                reasons[reason] += 1
                if args.dry_run:
                    print(f"keep {rel}:{row['line']}  *{text[:64]}*  — {reason}")
                continue

            i = row["line"] - 1
            if lines[i].strip() != "*" + text + "*":
                raise SystemExit(
                    f"{rel}:{i+1} no longer matches the audited candidate {lines[i].strip()!r} — "
                    "the audit is stale, re-run scripts/audit_lesson_media.py first")
            drop.add(i)
            # The caption is its own paragraph, so a neighbouring blank line goes with it.
            if i + 1 < len(lines) and lines[i + 1].strip() == "":
                drop.add(i + 1)
            if i - 1 >= 0 and lines[i - 1].strip() == "":
                drop.add(i - 1)
            deleted += 1
            if args.dry_run:
                print(f"  drop {rel}:{row['line']}  {lines[i].strip()[:88]}")

        if not drop:
            continue
        touched.append(rel)
        if not args.dry_run:
            with open(path, "w", encoding="utf-8") as fh:
                fh.write("\n".join(l for j, l in enumerate(lines) if j not in drop))

    print(f"\n{deleted} orphan captions removed, {kept} kept: "
          + ", ".join(f"{n} {why}" for why, n in reasons.items() if n)
          + (f"; {len(touched)} lesson files edited" if touched else "")
          + ("  (dry run, nothing written)" if args.dry_run else ""))
    return 0


if __name__ == "__main__":
    sys.exit(main())
