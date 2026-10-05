#!/usr/bin/env python3
"""Insert the 18 lessons' unreferenced diagrams into their lesson bodies.

Every placement is declared here explicitly, as an anchor plus a position, rather than
computed. That is deliberate: an earlier version scored each figure's caption against every
section of its lesson and picked the best match, which put 7 of 38 figures somewhere no rule
could justify. A declared anchor is auditable in one read, and the script refuses to write if
an anchor has stopped being unique — a stale line number fails loudly instead of putting a
figure in the wrong section.

Two kinds of edit:

  replace  The figure's own text was flattened into the body as a paragraph by the scrape, so
           the paragraph IS the figure, rendered as prose ("y x y = aˣ with a > 1 1"). Adding
           a figure beside it would print the caption twice. 6 of the 38 are like this.

  insert   The figure is absent and the prose needs it. 32 of the 38 are like this. Most go at
           the end of the section that names their concept; several must go somewhere exact,
           because the prose already refers to a figure by position — "The graph above
           illustrates…" is an instruction, not a description.

Alt text is copied verbatim from the figure's own SVG caption. It is never reworded, because
the caption is the author's words and a paraphrase is a claim the figure does not make. The six
Venn diagrams in sets.md have no caption at all — only symbolic labels — so their alt is listed
below as the one place this file interprets rather than copies, and it is flagged in the
comment.

Absolute `/diagrams/...` paths, not the corpus's other `svg/...` form: the relative form only
resolves because KatexContent rewrites it, so a figure written that way breaks if that rule
changes or a second renderer appears.

Usage: python3 scripts/place_lesson_figures.py [--dry-run]
"""
import argparse
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
LESSONS = os.path.join(ROOT, "data", "lessons")

# Alt text for figures whose SVG has no single caption to copy, because it is a multi-panel
# figure with one caption per panel or, for the Venn diagrams, no caption at all.
#
# These are the only alt strings not copied verbatim from one <text> run, and each is composed
# from the figure's own panel captions rather than reworded. Composing them is necessary: the
# auto-derivation below takes the single longest run, and for a two-panel figure that silently
# drops half of what the figure shows. polynomial-function-2.svg is the clearest case — its two
# panels are y = 2x² − 3x + 1 opening upward and y = −x² + 5x + 1 opening downward, and taking
# one run described only one of them.
EXPLICIT_ALT = {
    "polynomial-function-1.svg": "Two linear functions side by side: a strictly increasing one with positive slope and y-intercept 2, and a strictly decreasing one with negative slope and y-intercept 1.",
    "polynomial-function-2.svg": "Two quadratic functions side by side: y = 2x² − 3x + 1 opening upward and y = −x² + 5x + 1 opening downward. In each case the axis of symmetry is vertical and the vertex is the extremum.",
    "trinomials-1.svg": "The same parabola drawn twice: since the coefficient $a > 0$ the concavity faces upward, and since $a < 0$ it faces downward.",
    "principle-of-mathematical-induction-1.svg": "The two steps of an induction: the base case, true for every $n$ in the natural numbers, and the inductive step, $p(n)$ implies $p(n+1)$.",
}

# Alt text for the six Venn diagrams, which carry no caption — only the labels A, B, U and the
# result. These are read off the figure's visible maths, which is the only evidence there is,
# and are the only alt strings in this file that are not a verbatim copy. sets-6.svg's own
# <title> says "Group 2", which contradicts its content, so it is deliberately not used.
SETS_ALT = {
    "sets-1.svg": "Venn diagram of the union $A \\cup B$ of sets $A$ and $B$ within a universal set $U$.",
    "sets-2.svg": "Venn diagram of the intersection $A \\cap B$ of sets $A$ and $B$.",
    "sets-3.svg": "Venn diagram of the complement $A^c$ of $A$ within the universal set $U$.",
    "sets-4.svg": "Venn diagram of the difference $A \\setminus B$: the part of $A$ outside $B$.",
    "sets-5.svg": "Venn diagram of the symmetric difference $A \\triangle B$: the part of exactly one of $A$ and $B$.",
    "sets-6.svg": "Venn diagram of the distributive identity $A \\cap (B \\cup C) = (A \\cap B) \\cup (A \\cap C)$.",
}

# (lesson, figure, mode, anchor, position_note)
# anchor must be a substring occurring exactly once in the lesson body.
PLACEMENTS = [
    # ---- replace the figure's own flattened text --------------------------------------
    ("algebrica/functions/exponential-function.md", "exponential-function-1.svg", "replace",
     "always lies above the x-axis. y x y = a"),
    ("algebrica/functions/exponential-function.md", "exponential-function-2.svg", "replace",
     "When the base a is between 0 and  1, the function is decreasing."),
    ("algebrica/functions/exponential-function.md", "exponential-function-3.svg", "replace",
     "When a = 1, the function  is parallel to the x-axis."),
    ("algebrica/integrals/numerical-integration.md", "numerical-integration-1.svg", "replace",
     "The interval [a, b] has been divided  into 6 subintervals"),
    ("algebrica/integrals/numerical-integration.md", "numerical-integration-2.svg", "replace",
     "Each subinterval generates a trapezoid: h is the"),
    ("algebrica/sequences/euler-number-limit-sequence.md", "euler-number-limit-sequence-1.svg",
     "replace", "the  sequence converges to"),

    # ---- prose names the figure's position: the anchor is an instruction --------------
    # "The graph above illustrates ..." only reads correctly if the figures sit above it.
    ("algebrica/powers-radicals-logarithms/logarithms.md", "logarithms-1.svg", "before",
     "The graph above illustrates the monotonic behaviour"),
    ("algebrica/powers-radicals-logarithms/logarithms.md", "logarithms-2.svg", "before",
     "The graph above illustrates the monotonic behaviour"),
    # integers: 1 after the set-builder display, 2 after the sentence about spacing.
    ("algebrica/sets-and-numbers/integers.md", "integers-1.svg", "before",
     "an infinite collection of evenly spaced points along the number line."),
    ("algebrica/sets-and-numbers/integers.md", "integers-2.svg", "after",
     "an infinite collection of evenly spaced points along the number line."),
    # sets: each figure immediately after the display that defines its operation, expressed
    # as "before the paragraph that follows", because a display's closing \] is not unique.
    ("algebrica/sets-and-numbers/sets.md", "sets-1.svg", "before",
     "The intersection of \\(A\\) and \\(B\\) is the set of elements that belong to both sets"),
    ("algebrica/sets-and-numbers/sets.md", "sets-2.svg", "before",
     "If \\(A \\cap B = \\emptyset\\), the two sets are disjoint"),
    ("algebrica/sets-and-numbers/sets.md", "sets-3.svg", "before",
     "Another way to represent the complement of"),
    ("algebrica/sets-and-numbers/sets.md", "sets-4.svg", "before",
     "The relation \\(A \\setminus B \\neq B \\setminus A\\) holds"),
    ("algebrica/sets-and-numbers/sets.md", "sets-5.svg", "before",
     "An equivalent representation is given by the following expression"),
    ("algebrica/sets-and-numbers/sets.md", "sets-6.svg", "after",
     "The empty set and the universal set act as the identity element"),
    # cauchy: the geometric reading of the theorem, which the lesson's prose never states.
    ("algebrica/differential-calculus-theorems/cauchy-theorem.md", "cauchy-theorem-1.svg",
     "before", "## Statement"),

    # ---- end of the section that names the concept ------------------------------------
    ("algebrica/complex-numbers/roots-of-unity.md", "roots-of-unity-1.svg", "section_end",
     "## Geometric interpretation"),
    ("algebrica/functions/polynomial-function.md", "polynomial-function-1.svg", "section_end",
     "## Degree 1: linear functions"),
    ("algebrica/functions/polynomial-function.md", "polynomial-function-2.svg", "section_end",
     "## Degree 2: quadratic functions"),
    ("algebrica/functions/polynomial-function.md", "polynomial-function-3.svg", "section_end",
     "## Degree 3: cubic functions"),
    ("algebrica/inequalities/trigonometric-inequalities.md", "trigonometric-inequalities-1.svg",
     "section_end", "## Inequalities involving sine"),
    ("algebrica/inequalities/trigonometric-inequalities.md", "trigonometric-inequalities-2.svg",
     "section_end", "## Inequalities involving cosine"),
    ("algebrica/inequalities/trigonometric-inequalities.md", "trigonometric-inequalities-3.svg",
     "section_end", "## Inequalities involving tangent"),
    ("algebrica/integrals/indefinite-integrals.md", "indefinite-integrals-1.svg", "section_end",
     "## Primitives"),
    ("algebrica/integrals/indefinite-integrals.md", "indefinite-integrals-2.svg", "section_end",
     "## Primitives"),
    ("algebrica/integrals/integration-by-parts.md", "integration-by-parts-1.svg", "section_end",
     "## Derivation of the formula"),
    ("algebrica/polynomials/polynomials.md", "polynomials-1.svg", "section_end",
     "## Degree of a polynomial and its geometric interpretation"),
    ("algebrica/polynomials/polynomials.md", "polynomials-2.svg", "section_end",
     "## Polynomial Equations"),
    ("algebrica/polynomials/polynomials.md", "polynomials-3.svg", "section_end",
     "## End behavior of polynomial"),
    ("algebrica/polynomials/trinomials.md", "trinomials-1.svg", "section_end",
     "## Classification of trinomials"),
    ("algebrica/powers-radicals-logarithms/logarithms.md", "logarithms-3.svg", "section_end",
     "## Fundamental inequality for the natural logarithm"),
    ("algebrica/powers-radicals-logarithms/radicals.md", "radicals-1.svg", "section_end",
     "## Definition of radicals"),
    ("algebrica/powers-radicals-logarithms/radicals.md", "radicals-2.svg", "section_end",
     "## Geometric construction of the segment"),
    ("algebrica/sequences/principle-of-mathematical-induction.md",
     "principle-of-mathematical-induction-1.svg", "section_end", "## Mathematical induction"),
    ("algebrica/series/power-series.md", "power-series-1.svg", "section_end",
     "## Radius of Convergence"),
    ("algebrica/series/power-series.md", "power-series-2.svg", "section_end",
     "## Radius of Convergence"),
    ("algebrica/vectors-and-matrices/eigenvalues-and-eigenvectors.md",
     "eigenvalues-and-eigenvectors-1.svg", "section_end", "## Definition"),
]

FIGURES = os.path.join(ROOT, "web", "next-app", "public", "diagrams", "algebrica")


def unescape(t):
    return (t.replace("&gt;", ">").replace("&lt;", "<").replace("&quot;", '"')
             .replace("&apos;", "'").replace("&amp;", "&"))


def svg_caption(name):
    """The figure's own caption: its longest <text> run. Copied, never reworded."""
    path = os.path.join(FIGURES, name)
    if not os.path.exists(path):
        raise SystemExit(f"missing figure file: {path}")
    with open(path, encoding="utf-8") as fh:
        s = fh.read()
    runs = [unescape(re.sub(r"\s+", " ", re.sub(r"<[^>]+>", "", t))).strip()
            for t in re.findall(r"<text[^>]*>(.*?)</text>", s, re.S)]
    runs = [r for r in runs if r]
    # Only one long run is a single caption. Short runs are axis and curve labels ("x", "y",
    # "y = aˣ"), which are not captions and are not alt text on their own.
    captions = [r for r in runs if len(r) > 30]
    if len(captions) != 1:
        raise SystemExit(
            f"{name}: {len(captions)} caption runs, so there is no single caption to copy. "
            "Add an explicit alt rather than letting the longest run stand in for the figure — "
            "on a multi-panel figure that describes only one panel.")
    return captions[0]


def alt_for(name):
    """The figure's own words. Copied when there is exactly one caption to copy; declared in
    EXPLICIT_ALT or SETS_ALT when the figure has several panels or none."""
    if name in EXPLICIT_ALT:
        return EXPLICIT_ALT[name]
    if name in SETS_ALT:
        return SETS_ALT[name]
    return svg_caption(name)


def find(lines, anchor, lesson, figure):
    hits = [i for i, l in enumerate(lines) if anchor in l]
    if len(hits) != 1:
        raise SystemExit(
            f"{lesson}: anchor {anchor!r} for {figure} matched {len(hits)} lines; "
            "it must be unique or the placement is not auditable")
    return hits[0]


def section_end(lines, heading, lesson, figure):
    idx = next((i for i, l in enumerate(lines) if l.startswith("#") and heading in l), None)
    if idx is None:
        raise SystemExit(f"{lesson}: no section {heading!r} for {figure}")
    end = next((i for i in range(idx + 1, len(lines))
                if re.match(r"^#{1,6}\s", lines[i])), len(lines))
    # Back off a trailing blockquote: a figure after "A practical remark: ..." reads as if
    # the remark were illustrating the figure.
    tail = end
    while tail - 1 > idx and lines[tail - 1].strip() == "":
        tail -= 1
    if tail - 1 > idx and lines[tail - 1].lstrip().startswith(">"):
        while tail - 1 > idx and lines[tail - 1].lstrip().startswith(">"):
            tail -= 1
        while tail - 1 > idx and lines[tail - 1].strip() == "":
            tail -= 1
    return tail


def reference(name, alt):
    return f"![{alt}](/diagrams/algebrica/{name})"


def main():
    ap = argparse.ArgumentParser(description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--dry-run", action="store_true", help="report placements, write nothing")
    args = ap.parse_args()

    by_lesson = {}
    for lesson, figure, mode, anchor in PLACEMENTS:
        by_lesson.setdefault(lesson, []).append((figure, mode, anchor))

    for lesson in sorted(by_lesson):
        path = os.path.join(LESSONS, lesson)
        with open(path, encoding="utf-8") as fh:
            lines = fh.read().split("\n")
        planned = []

        # Group by target line. Two figures can legitimately land on the same line — both
        # belong to one section, or the prose names one figure position for two of them — and
        # applying those edits independently reverses them, because every insert lands at the
        # same index. Splicing each line once, in declared order, keeps them in order.
        buckets = {}
        for figure, mode, anchor in by_lesson[lesson]:
            ref = reference(figure, alt_for(figure))
            if mode == "section_end":
                idx, slot = section_end(lines, anchor, lesson, figure), "after"
            else:
                idx = find(lines, anchor, lesson, figure)
                slot = "replace" if mode == "replace" else ("before" if mode == "before" else "after")
            b = buckets.setdefault(idx, {"before": [], "after": [], "replace": None})
            if slot == "replace":
                assert b["replace"] is None and not b["before"] and not b["after"], (
                    f"{lesson}: two figures claim to replace line {idx + 1}")
                assert not lines[idx].startswith("!["), f"{lesson}:{idx+1} already a reference"
                b["replace"] = ref
            else:
                b[slot].append((ref, figure))

        # Descending index order, so an insertion cannot shift a later target.
        for idx in sorted(buckets, reverse=True):
            b = buckets[idx]
            if b["before"] and idx >= len(lines):
                raise SystemExit(f"{lesson}: nothing to insert before at end of file")
            # Exactly one blank line on each side of a spliced reference, and no more. The
            # reuse of an existing blank line is what keeps the diff to the figure alone: an
            # earlier version normalised blank lines across the whole file and removed about a
            # hundred that had nothing to do with this work.
            blank_before = idx > 0 and lines[idx - 1].strip() == ""
            blank_after = idx < len(lines) and lines[idx + 1].strip() == ""
            seg = []
            for ref, _ in b["before"]:
                seg += ([ref, ""] if blank_before else ["", ref, ""])
                blank_before = True  # a second figure needs its own separator
            if b["replace"] is not None:
                seg += [b["replace"]]
            else:
                if idx < len(lines):
                    seg += [lines[idx]]
                for ref, _ in b["after"]:
                    seg += ([ref, ""] if blank_after else ["", ref, ""])
                    blank_after = True
            lines[idx:idx + 1] = seg
            for ref, figure in b["before"] + b["after"] + ([(b["replace"], "")] if b["replace"] else []):
                planned.append((idx + 1, figure))

        if args.dry_run:
            print(f"--- {lesson}")
            for ln, figure in sorted(planned):
                print(f"    line {ln:4d}  {figure}")
        else:
            with open(path, "w", encoding="utf-8") as fh:
                fh.write("\n".join(lines))
            print(f"ok {lesson}: {len(planned)} figure(s)")

    print(f"\n{len(PLACEMENTS)} placements across {len(by_lesson)} lessons"
          + ("  (dry run, nothing written)" if args.dry_run else ""))


if __name__ == "__main__":
    sys.exit(main())