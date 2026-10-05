#!/usr/bin/env python3
"""Report on the lesson corpus's instructional media. Changes nothing.

This is a *report*, not a gate. It exists because the media situation could not be
described before assets became addressable: with figures trapped inside an opaque markdown
body, the only available questions were "does the concept map point at a file" and "is the
lesson long enough". Neither one says whether a figure a lesson promises actually exists,
which is why 349 shipped diagrams sat unreferenced while 455 paragraphs of stripped figure
captions sat in the teaching corpus with no way to enumerate them.

Three findings, in the order they matter:

  A. Orphan caption candidates — paragraphs left behind when a figure's image was dropped
     but its caption was kept. Confidence-rated, never deleted. `ingest_pretext.mjs` emits
     exactly this shape for a `<figure><caption>`, so the shape is the pipeline's signature;
     the *text* still has to be read, because emphasis and section labels take the same
     shape and a rule that cannot tell them apart would delete content.

  B. Prose that promises a figure — sentences referring to "the figure below" with no
     figure to point at. Small in number and unambiguous.

  C. Shipped-but-unreferenced diagrams, with candidate teaching lessons for each. This is
     the pool that could fill section A without authoring anything: the figures already
     exist and are already served.

Usage:
    python3 scripts/audit_lesson_media.py                    # print the report
    python3 scripts/audit_lesson_media.py --json <path>      # also write machine-readable
    python3 scripts/audit_lesson_media.py --section a        # one section only

Exit 0 always. A report that exits non-zero invites a blind fix; the judgement this needs
is per instance, and a gate would only record that nobody has made it yet.
"""
import argparse
import json
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
LESSONS = os.path.join(ROOT, "data", "lessons")
PUBLIC = os.path.join(ROOT, "web", "next-app", "public")
DIAGRAMS = os.path.join(PUBLIC, "diagrams")

# Mirrors imageRe in internal/lessons/assets.go, including its tolerance for a figure
# embedded mid-paragraph and for alt text wrapped over lines.
IMAGE_RE = re.compile(r"!\[((?:[^\]\\]|\\.)*)\]\(\s*([^)\s]+)\s*\)")
IMAGE_OPEN_RE = re.compile(r"!\[[^\]]*\]\(")

# Textual markers that a paragraph is describing a picture. Used only to *rank* candidates,
# never to delete: absence of a marker is not evidence of absence.
FIGURE_MARKERS = re.compile(
    r"\b(graph|graph of|graphs|diagram|figure|image|photo|photograph|picture|plot|"
    r"sketch|table|shown|shows|showing|illustrat\w*|depict\w*|visib\w*|"
    r"image by|photo by|photograph by|credit|source:|adapted from|taken from)\b",
    re.IGNORECASE,
)

# Prose that refers to a figure which is not in the body.
FIGURE_REF_PROSE = re.compile(
    r"(the (figure|graph|diagram|picture|image|table) (below|above|to the (left|right))"
    r"|as shown in the (figure|graph|diagram)"
    r"|see the (figure|graph|diagram)"
    r"|the (figure|graph|diagram) (shows|shows|depicts|illustrates))",
    re.IGNORECASE,
)

IMAGE_EXTS = {".svg", ".png", ".jpg", ".jpeg", ".webp", ".gif"}


def resolve(src):
    """Mirror of resolveLessonImageUrl in KatexContent.tsx, and of the rule in
    internal/lessons/assets_test.go. A relative reference is a bare filename resolved
    against the algebrica diagram directory; anything absolute or remote is used as-is."""
    url = (src or "").strip()
    if url and not url.startswith("http") and not url.startswith("/"):
        url = "/diagrams/algebrica/" + url.split("/")[-1]
    return url


def strip_emphasis(t):
    t = t.strip()
    if t.startswith("*") and t.endswith("*") and len(t) > 2:
        return t[1:-1]
    return t


def is_standalone_emphasis(line):
    """True for `*text*` alone on its line.

    Excludes the two near-misses that a naive scan counts as captions and must not: a
    bullet item (`* **Sampling**`) and a bold-only line (`**Header**`). Counting those was
    how an earlier pass of this report overstated the candidate set.
    """
    t = line.strip()
    if not (t.startswith("*") and t.endswith("*") and len(t) > 2):
        return False
    if t.startswith("* "):
        return False
    inner = t[1:-1]
    if inner.startswith("*") and inner.endswith("*"):
        return False
    return "*" not in inner


def prev_nonblank(lines, i):
    for j in range(i - 1, -1, -1):
        if lines[j].strip():
            return lines[j].strip()
    return ""


# PreTeXt auto-numbers figures and tables, and the number survives ingestion. A paragraph
# opening "Figure 3.2" or "Table 1.4" was a caption by construction, whatever else it says —
# 236 of them in the teaching corpus. This is the one signal that needs no judgement.
NUMBERED_CAPTION_RE = re.compile(r"^(figure|table|graph|diagram)\s+\d+(\.\d+)*\b", re.IGNORECASE)


def rate_candidate(text, in_file_count, corpus):
    """Confidence that a standalone-italic paragraph is a stripped figure caption.

    Calibrated against the corpus rather than against how the sentence reads, because the
    first version of this rated on picture vocabulary alone and so reported a genuine
    caption — "Share of all income held by the top 1%, United States" — as low confidence.
    Describing a picture is evidence; not describing one is not evidence of anything.

    Tiers, strongest first:

      certain  opens with an auto-numbered "Figure 3.2" / "Table 1.4" label
      strong   describes a picture, or the same sentence is duplicated in its lesson
      likely   in the teaching corpus, where this shape is the ingestion pipeline's own
               signature and inline emphasis occurs in 1 lesson out of 124
    """
    if NUMBERED_CAPTION_RE.match(text):
        return "certain", "auto-numbered figure/table label"
    if FIGURE_MARKERS.search(text):
        return "strong", "describes a picture"
    if in_file_count > 1:
        return "strong", f"duplicated {in_file_count}x in this lesson"
    if corpus == "teaching":
        return "likely", "teaching corpus: standalone emphasis is the dropped-figure shape"
    return "weak", "no picture language and not an auto-numbered label"


def section_a_candidates():
    rows = []
    for dirpath, _dirs, names in os.walk(LESSONS):
        for name in sorted(names):
            if not name.endswith(".md"):
                continue
            path = os.path.join(dirpath, name)
            rel = os.path.relpath(path, LESSONS)
            body = open(path, encoding="utf-8").read()
            lines = body.split("\n")
            hits = [
                (i, strip_emphasis(lines[i]))
                for i in range(len(lines))
                if is_standalone_emphasis(lines[i])
                and not IMAGE_OPEN_RE.search(lines[i])
                and not IMAGE_OPEN_RE.search(prev_nonblank(lines, i))
            ]
            texts = [t for _, t in hits]
            for i, text in hits:
                # Paragraph-level check: only a paragraph that is nothing but emphasis is
                # the shape ingest_pretext.mjs emits for a figure caption.
                before = prev_nonblank(lines, i)
                after = next((lines[j].strip() for j in range(i + 1, len(lines)) if lines[j].strip()), "")
                if before.startswith(("-", "*", "|", ">")) or after.startswith(("-", "*", "|", ">")):
                    continue
                corpus = rel.split(os.sep)[0].split("/")[0]
                conf, why = rate_candidate(text, texts.count(text), corpus)
                rows.append(
                    {"lesson": rel, "line": i + 1, "text": text, "confidence": conf, "reason": why}
                )
    return rows


def section_b_prose():
    rows = []
    for dirpath, _dirs, names in os.walk(LESSONS):
        for name in sorted(names):
            if not name.endswith(".md"):
                continue
            path = os.path.join(dirpath, name)
            rel = os.path.relpath(path, LESSONS)
            lines = open(path, encoding="utf-8").read().split("\n")
            # What can satisfy a figure reference: an image, a GFM table (a lesson may call a
            # table its "diagram below" — integration-strategies.md:51 does exactly that), or
            # a display-math block. Checking only for images reported that lesson as a hole
            # when its table was two lines below, so the check has to name what was promised.
            referents = [
                j for j, l in enumerate(lines)
                if IMAGE_OPEN_RE.search(l)
                or "![" in l          # opening bracket of a reference whose alt text wraps
                or l.lstrip().startswith("|")
                or l.lstrip().startswith("\\[")
            ]
            for i, line in enumerate(lines):
                for m in FIGURE_REF_PROSE.finditer(line):
                    # "The diagram below shows ..." is only a hole when nothing follows it,
                    # and "The graph above illustrates ..." is only a hole when nothing
                    # precedes it — measured rather than assumed: riemann-integrability-
                    # criteria.md:27 says "below" and the figure it means is the next line,
                    # while logarithms.md:50 says "above" and has nothing above it.
                    direction = -1 if re.search(r"\babove\b", m.group(0), re.IGNORECASE) else 1
                    # The referent may sit on the same line: uniform-distribution.md:19 is a
                    # sentence and its figure on one line, and the figure is on line 19.
                    window = range(i, i + 9) if direction > 0 else range(i - 1, max(-1, i - 9), -1)
                    if any(j in window for j in referents):
                        continue
                    rows.append(
                        {
                            "lesson": rel,
                            "line": i + 1,
                            "phrase": m.group(0),
                            "context": line.strip()[:200],
                        }
                    )
    return rows


def norm_stem(name):
    """Reduce a file stem to comparable topic words: drop a trailing figure index, map
    separators to spaces, drop filler. `floor-and-ceiling-functions-2` and
    `floor_and_ceiling_function` both reduce to the same thing."""
    stem = os.path.splitext(name)[0]
    stem = re.sub(r"-\d+$", "", stem)
    stem = re.sub(r"[_\-.]+", " ", stem).lower()
    stem = re.sub(r"\b(the|a|an|of|and|to|in|on|for)\b", " ", stem)
    return re.sub(r"\s+", " ", stem).strip()


def stem_figure_index(name):
    m = re.search(r"-(\d+)\.(?:svg|png|jpe?g|webp|gif)$", name)
    return int(m.group(1)) if m else None


def section_c_unreferenced():
    referenced = set()
    every_lesson = {}          # rel -> urls it references (empty set means no figure at all)
    for dirpath, _dirs, names in os.walk(LESSONS):
        for name in sorted(names):
            if not name.endswith(".md"):
                continue
            path = os.path.join(dirpath, name)
            rel = os.path.relpath(path, LESSONS)
            body = open(path, encoding="utf-8").read()
            urls = {resolve(m.group(2)) for m in IMAGE_RE.finditer(body)}
            every_lesson[rel] = urls
            referenced |= urls

    shipped = {}
    for dirpath, _dirs, names in os.walk(DIAGRAMS):
        for name in sorted(names):
            if os.path.splitext(name)[1].lower() in IMAGE_EXTS:
                url = "/" + os.path.relpath(os.path.join(dirpath, name), PUBLIC).replace(os.sep, "/")
                shipped[url] = name

    # Site chrome (dark/light pairs, architecture-*, system-design-*) is referenced by
    # components rather than lesson bodies, so it is excluded from the teaching candidates.
    SITE = re.compile(r"-(dark|light)\.svg$")
    candidates = sorted(u for u in shipped if u not in referenced)
    content = [u for u in candidates if not SITE.search(u)]

    # Index unreferenced diagrams by topic stem, then by each individual topic word, so a
    # lesson can be offered a figure even when the stems do not match exactly.
    by_stem, by_word = {}, {}
    stop = {"of", "the", "a", "an", "and", "to", "in", "on", "for", "function", "functions"}
    for url in content:
        stem = norm_stem(os.path.basename(url))
        by_stem.setdefault(stem, []).append(url)
        for w in stem.split():
            if len(w) > 3 and w not in stop:
                by_word.setdefault(w, []).append(url)

    # Every lesson, not only the ones that already have a figure: the actionable set is the
    # lessons with *no* figure at all, and listing only the ones that have one would rank the
    # 242 that are already fine above the 124 that need help.
    rows = []
    for lesson in sorted(every_lesson):
        urls = every_lesson[lesson]
        lesson_stem = norm_stem(os.path.splitext(os.path.basename(lesson))[0])
        exact = by_stem.get(lesson_stem, [])
        partial = []
        if not exact:
            for w in lesson_stem.split():
                if len(w) > 3 and w in by_word:
                    partial.extend(by_word[w])
        partial = sorted({u for u in partial if u not in exact})
        if exact or partial:
            rows.append(
                {
                    "lesson": lesson,
                    "figures_referenced": bool(urls),
                    "exact_stem_match": exact[:8],
                    "topic_word_match": partial[:8],
                }
            )
    rows.sort(key=lambda r: (r["figures_referenced"], not r["exact_stem_match"], r["lesson"]))
    return {"referenced": len(referenced), "shipped": len(shipped),
            "unreferenced_content": len(content), "candidates": rows,
            "unreferenced_urls": content}


def render(a_rows, b_rows, c):
    out = []
    W = out.append
    W("# Lesson media audit\n")
    W("Generated by `scripts/audit_lesson_media.py`. Nothing here has been changed — this is")
    W("a report, because the remaining decisions are per instance.\n")

    W("## A. Orphan caption candidates\n")
    tiers = ["certain", "strong", "likely", "weak"]
    counts = {t: sum(1 for r in a_rows if r["confidence"] == t) for t in tiers}
    W(f"{len(a_rows)} paragraphs across "
      f"{len({r['lesson'] for r in a_rows})} lessons look like a figure caption whose image")
    W(f"was dropped while its text was kept: "
      + ", ".join(f"**{counts[t]} {t}**" for t in tiers if counts[t]) + ".\n")
    W("`certain` opens with an auto-numbered `Figure 3.2` / `Table 1.4` label, which PreTeXt")
    W("generates for a figure or table and which survives ingestion — no judgement needed.")
    W("`strong` describes a picture, or is duplicated in its lesson. `likely` is the teaching")
    W("corpus, where a standalone italic paragraph is the ingestion pipeline's own signature:")
    W("`ingest_pretext.mjs` emits `<p><em>caption</em></p>` for a figure and drops the image,")
    W("and inline emphasis appears in 1 teaching lesson out of 124.\n")
    W("Not describing a picture is not evidence against being a caption — a caption like")
    W("*Share of all income held by the top 1%, United States* says nothing about pictures —")
    W("so the tiers below do not treat it as evidence either way.\n")
    W("| lesson | line | confidence | text | why |")
    W("| --- | --- | --- | --- | --- |")
    for r in a_rows:
        text = r["text"].replace("|", "\\|")
        if len(text) > 96:
            text = text[:93] + "..."
        W(f"| `{r['lesson']}` | {r['line']} | {r['confidence']} | *{text}* | {r['reason']} |")

    W("\n## B. Prose that promises a figure\n")
    if b_rows:
        W(f"{len(b_rows)} sentences refer to a figure that is not in the body.\n")
        W("| lesson | line | phrase | context |")
        W("| --- | --- | --- | --- |")
        for r in b_rows:
            W(f"| `{r['lesson']}` | {r['line']} | {r['phrase']} | {r['context'].replace('|', chr(92) + '|')} |")
    else:
        W("None. Every figure reference in the prose resolves to a figure.")

    W("\n## C. Shipped but unreferenced\n")
    W(f"{c['shipped']} diagrams ship, {c['referenced']} URLs are referenced by lesson bodies, and")
    W(f"**{c['unreferenced_content']} content diagrams are referenced by nothing** (site chrome")
    W("with dark/light variants is excluded — components reference those).\n")
    W("These already exist and are already served, so a teaching figure that is missing today")
    W("may already be on disk under a different name. Matches below are *candidates only*:\n")
    if c["candidates"]:
        W("Sorted so lessons with **no figure at all** come first — those are the ones where a")
        W("figure already on disk could be dropped in without authoring anything.\n")
        W("| lesson | references a figure | exact stem match | topic-word match |")
        W("| --- | --- | --- | --- |")
        for r in c["candidates"]:
            ex = ", ".join(f"`{os.path.basename(u)}`" for u in r["exact_stem_match"]) or "—"
            pa = ", ".join(f"`{os.path.basename(u)}`" for u in r["topic_word_match"]) or "—"
            W(f"| `{r['lesson']}` | {'yes' if r['figures_referenced'] else 'no'} | {ex} | {pa} |")
    else:
        W("No stem or topic-word matches.")
    return "\n".join(out) + "\n"


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--json", help="also write the full findings as JSON to this path")
    ap.add_argument("--section", choices=["a", "b", "c"], help="print one section only")
    args = ap.parse_args()

    a = section_a_candidates()
    b = section_b_prose()
    c = section_c_unreferenced()

    if args.section == "a":
        print(json.dumps(a, indent=2))
        return
    if args.section == "b":
        print(json.dumps(b, indent=2))
        return
    if args.section == "c":
        print(json.dumps({k: v for k, v in c.items() if k != "candidates"}, indent=2))
        return

    print(render(a, b, c))
    if args.json:
        payload = {"orphan_caption_candidates": a, "figure_promising_prose": b,
                   "unreferenced_diagrams": {k: v for k, v in c.items()}}
        with open(args.json, "w", encoding="utf-8") as fh:
            json.dump(payload, fh, indent=2)
        print(f"wrote {args.json}")


if __name__ == "__main__":
    sys.exit(main())