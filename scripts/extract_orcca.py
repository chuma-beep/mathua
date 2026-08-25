"""
Extract lesson content from ORCCA PreTeXt XML files.
Converts to markdown format matching the project's style.
"""
import json
import os
import re
import requests
import xml.etree.ElementTree as ET
from pathlib import Path

ORCCA_BASE = "https://raw.githubusercontent.com/PCCMathSAC/orcca/edition/src"
OUTPUT_DIR = Path("data/lessons/teaching")

# Common SI/imperial unit spellings keyed by (prefix, base).
UNIT_NAMES = {
    ("centi", "meter"): "cm",
    ("milli", "meter"): "mm",
    ("kilo", "meter"): "km",
    (None, "meter"): "m",
    (None, "inch"): "in",
    (None, "foot"): "ft",
    (None, "yard"): "yd",
    (None, "mile"): "mi",
    ("centi", "gram"): "cg",
    ("milli", "gram"): "mg",
    ("kilo", "gram"): "kg",
    (None, "gram"): "g",
    ("centi", "liter"): "cL",
    ("milli", "liter"): "mL",
    ("kilo", "liter"): "kL",
    (None, "liter"): "L",
    (None, "second"): "s",
    (None, "minute"): "min",
    (None, "hour"): "hr",
    (None, "dollar"): "\\$",
}


def local(tag):
    t = str(tag)
    return t.split("}")[-1] if "}" in t else t


def norm(text):
    """Collapse internal newlines/indentation inherited from pretty-printed XML."""
    return re.sub(r"\s+", " ", text).strip()


def unit_latex(unit_el):
    prefix = unit_el.get("prefix")
    base = unit_el.get("base")
    exp = unit_el.get("exp")
    name = UNIT_NAMES.get((prefix, base)) or UNIT_NAMES.get((None, base))
    if not name:
        name = norm(base or "") or "unit"
    out = f"\\text{{{name}}}"
    if exp:
        out += f"^{{{exp}}}"
    return out


def quantity_latex(el):
    """Render a <quantity> (mag/unit[/per...]) pair as inline LaTeX."""
    pieces = []
    for child in el:
        tag = local(child.tag)
        if tag == "mag":
            pieces.append(norm("".join(child.itertext())))
        elif tag == "unit":
            if pieces and pieces[-1] not in ("/",):
                pieces.append("\\,")
            pieces.append(unit_latex(child))
        elif tag == "per":
            while pieces and pieces[-1] == "\\,":
                pieces.pop()
            pieces.append("/")
    body = "".join(pieces).strip()
    if not body:
        return ""
    return f"${body}$"


def row_latex(mrow_el):
    """LaTeX of one <mrow>, with ORCCA's \\amp alignment marker converted."""
    src = "".join(mrow_el.itertext())
    src = norm(src)
    src = re.sub(r"\\(?:math|maht)(?:open|close)\{\}", "", src)
    return esc_dollars(src.replace("\\amp", "&"))


def md_inline(md_el):
    """Render <md> as inline math (for list-item contexts where display
    blocks would break the list structure)."""
    rows = [child for child in md_el if local(child.tag) == "mrow"]
    if not rows:
        return ""
    body = " \\\\ ".join(row_latex(r) for r in rows)
    return f"$\\begin{{aligned}} {body} \\end{{aligned}}$"


def md_latex(md_el):
    """Render an <md> display-math block as an aligned environment.

    Emitted as ONE line: the canonicalizer's display-dialect detection only
    trusts single-line \[...\] blocks."""
    rows = [child for child in md_el if local(child.tag) == "mrow"]
    if not rows:
        return ""
    body = " \\\\ ".join(row_latex(r) for r in rows)
    return f"\n\\[ \\begin{{aligned}} {body} \\end{{aligned}} \\]\n"


def list_md(list_el, inline_md=False):
    """Render <ol>/<ul> as a markdown list block."""
    items = []
    for li in list_el:
        if local(li.tag) != "li":
            continue
        t = norm(text_content(li, inline_md=True))
        if t:
            items.append(t)
    if not items:
        return ""
    marker = lambda i: f"{i + 1}."
    lines = [f"{marker(i)} {t}" for i, t in enumerate(items)]
    return "\n" + "\n".join(lines) + "\n"


def list_block_md(list_el):
    """Render a PreTeXt <list> (optionally a <dl> of titled items)."""
    title_el = list_el.find("title")
    header = ""
    if title_el is not None:
        header = f"**{norm(''.join(title_el.itertext()))}**"

    dl = list_el.find("dl")
    if dl is not None:
        lines = []
        for li in dl:
            if local(li.tag) != "li":
                continue
            li_title = li.find("title")
            label = ""
            if li_title is not None:
                label = f"**{norm(''.join(li_title.itertext()))}**: "
            body = norm(
                " ".join(
                    t for t in (
                        text_content(child)
                        for child in li
                        if local(child.tag) != "title"
                    ) if t
                )
            )
            if body:
                lines.append(f"- {label}{body}")
        if not lines:
            return ""
        block = "\n".join(lines)
        if header:
            block = f"{header}\n\n{block}"
        return "\n" + block + "\n"

    # Plain <list> without <dl>: fall back to rendering its children.
    inner = norm(
        " ".join(
            t for t in (text_content(child) for child in list_el if local(child.tag) != "title") if t
        )
    )
    if not inner:
        return ""
    prefix = f"{header} " if header else ""
    return "\n" + prefix + inner + "\n"


def esc_dollars(s):
    """Escape unescaped $ so currency inside math cannot break pairing."""
    return re.sub(r"(?<!\\)\$", "\\\\$", s)


def strip_delimiter_dollars(inner):
    """Authors sometimes wrap <m> content in literal $ delimiters; drop them
    so we don't emit doubled/odd dollar runs inside our own span."""
    return inner.lstrip("$").rstrip("$").strip()


def text_content(el, inline_md=False):
    """Recursively extract text from element, handling math and formatting."""
    if el is None:
        return ""
    parts = []
    tag = local(el.tag)
    if tag == "m":
        inner = norm("".join(el.itertext()))
        inner = strip_delimiter_dollars(inner)
        # Literal dollars (currency) inside math must not break pairing.
        inner = esc_dollars(inner)
        return f"${inner}$"
    if tag == "me":
        inner = esc_dollars(norm("".join(el.itertext())))
        return f"\n\\[{inner}\\]\n" if not inline_md else f"${inner}$"
    if tag == "md":
        return md_inline(el) if inline_md else md_latex(el)
    if tag == "quantity":
        return quantity_latex(el)
    if tag in ("ol", "ul"):
        return list_md(el, inline_md=inline_md)
    if tag == "term":
        return f"*{norm(''.join(el.itertext()))}*"
    if tag == "em":
        return f"*{norm(''.join(el.itertext()))}*"
    if tag == "q":
        return f"\"{norm(''.join(el.itertext()))}\""
    if tag == "idx":
        return ""
    if tag == "url":
        return el.get("href", "")
    if tag == "c":
        return f"`{norm(''.join(el.itertext())).strip()}`"
    if tag == "mdash":
        return "---"
    if tag == "nbsp":
        return " "
    if tag == "times":
        return "×"
    if tag == "list":
        return list_block_md(el)
    if tag in ("figure", "caption", "tabular"):
        return ""
    if tag in ("var", "pg-code", "instruction"):
        # Interactive WeBWorK answer slots / code have no print representation.
        return ""

    if el.text:
        parts.append(norm(el.text))
    for child in el:
        child_text = text_content(child, inline_md=inline_md)
        if child_text:
            parts.append(child_text)
        if child.tail:
            parts.append(norm(child.tail))
    return " ".join(p for p in parts if p).strip()


# Block-level tags dispatched by emit_block inside section bodies.
SKIP_BLOCKS = {"objectives"}


def emit_block(child, md_parts):
    tag = local(child.tag)
    if tag == "p":
        t = text_content(child)
        if t:
            md_parts.append(t)
            md_parts.append("")
    elif tag == "definition":
        _extract_definition(child, md_parts)
    elif tag == "fact":
        _extract_fact(child, md_parts)
    elif tag == "example":
        _extract_example(child, md_parts)
    elif tag == "exercise":
        _extract_exercise(child, md_parts)
    elif tag == "warning":
        _extract_warning(child, md_parts)
    elif tag == "aside":
        _extract_aside(child, md_parts)
    elif tag == "figure":
        _skip_figure(child, md_parts)
    elif tag == "tabular":
        _extract_tabular(child, md_parts)
    elif tag == "sidebyside":
        _extract_sidebyside(child, md_parts)
    elif tag == "md":
        t = md_latex(child)
        if t:
            md_parts.append(t.strip())
            md_parts.append("")
    elif tag in ("ol", "ul"):
        t = list_md(child)
        if t:
            md_parts.append(t.strip())
            md_parts.append("")
    elif tag == "list":
        t = text_content(child)
        if t:
            md_parts.append(t.strip())
            md_parts.append("")
    elif tag == "paragraphs":
        _extract_paragraphs(child, md_parts)
    elif tag in SKIP_BLOCKS:
        pass


def _extract_paragraphs(el, md_parts):
    title_el = el.find("title")
    if title_el is not None:
        md_parts.append(f"### {norm(''.join(title_el.itertext()))}")
        md_parts.append("")
    for child in el:
        emit_block(child, md_parts)


def extract_section(slug):
    """Fetch and parse an ORCCA section, return markdown content."""
    url = f"{ORCCA_BASE}/{slug}.ptx"
    print(f"  Fetching {url}")
    resp = requests.get(url, timeout=30)
    resp.raise_for_status()

    root = ET.fromstring(resp.content)
    section = root

    title_el = section.find("title")
    title_raw = "".join(title_el.itertext()).strip() if title_el is not None else slug.replace("-", " ").title()
    # Strip leading "Section " or "section " prefix from titles
    title = re.sub(r"^Section\s+", "", title_raw, flags=re.IGNORECASE)

    md_parts = [f"# {title}", ""]

    # Introduction
    intro = section.find("introduction")
    if intro is not None:
        for child in intro:
            emit_block(child, md_parts)

    # Subsections
    for subsection in section.findall("subsection"):
        sub_title_el = subsection.find("title")
        if sub_title_el is not None:
            md_parts.append(f"## {norm(''.join(sub_title_el.itertext()))}")
            md_parts.append("")

        for child in subsection:
            emit_block(child, md_parts)

    md = "\n".join(md_parts)
    return clean_prose(md)


def balance_line(ln):
    """If a line carries an odd count of unescaped $, escape the last one.
    Stray currency dollars are the usual culprit; leaving one unpaired makes
    the frontend scanner pair across lines and swallow headings."""
    positions = [m.start() for m in re.finditer(r"(?<!\\)\$", ln)]
    if len(positions) % 2 == 1:
        i = positions[-1]
        ln = ln[:i] + "\\$" + ln[i + 1 :]
    return ln


def clean_prose(s):
    """Post-extraction cleanup of scrape artifacts."""
    # Nested-substitute safe expansion (two passes covers one level of nesting).
    for _ in range(2):
        s = re.sub(r"\\substitute\{([^{}]*)\}", r"\1", s)
        s = re.sub(r"\\divideunder\{([^{}]*)\}\{([^{}]*)\}", r"\\frac{\1}{\2}", s)
    # Empty delimiter-fence groups (and ORCCA's typo'd \mahtclose variant).
    s = re.sub(r"\\(?:math|maht)(?:open|close)\{\}", "", s)
    # ORCCA step-highlight macro: keep the content.
    s = re.sub(r"\\nextoperation\{([^{}]*)\}", r"\1", s)
    # Empty parens left behind by removed <var> answer slots.
    s = s.replace("()", "")
    # Space before punctuation from stripped elements.
    s = re.sub(r" ([.,;:!?])", r"\1", s)
    s = re.sub(r" +\)", ")", s)
    # "$x$ -value" glue from stripped element boundaries.
    s = s.replace("$ -", "$-")
    # Orphan punctuation on its own line after a display-math block: attach it.
    s = re.sub(r"(\\\])\s*\n\s*([.,])", r"\1\2", s)
    # Doubled spaces (outside of deliberate markdown indentation).
    s = re.sub(r"(?<![>\n]) {2,}(?!\n)", " ", s)
    # More than two consecutive newlines.
    s = re.sub(r"\n{3,}", "\n\n", s)

    # Authors sometimes write math with $ delimiters that the defense then
    # escapes; if the interior clearly carries LaTeX commands, unescape the
    # pair back into a live span.
    s = re.sub(
        r"\$([^$\n]*\\(?:frac|cdot|div|text|pi|times)\b[^$\n]*)\$",
        lambda m: "$" + m.group(1) + "$",
        s,
    )
    # Delimiter-less formulas pasted as plain text: wrap the known shapes.
    s = re.sub(
        r"(?<![\\$\w])(\w+\^\{?\w+}?\s*\\cdot\s*\w+\^\{?\w+}?\s*=\s*\w+\^\{?\w+(?:[+-]\w+)}?)(?![\w$])",
        lambda m: "$" + esc_dollars(m.group(1)) + "$",
        s,
    )
    s = re.sub(
        r"(?<![\\$\w])(m\s*=\s*\\frac\{\\text\{change in[^}]*\}\}\{\\text\{change in[^}]*\}\})(?![\w$])",
        lambda m: "$" + esc_dollars(m.group(1)) + "$",
        s,
    )
    s = re.sub(
        r"(?<![\\$\w])(\(\\text\{initial value\}\)\s*\\pm\s*\(\\text\{percent as decimal\}\)\s*\\cdot\s*\(\\text\{initial value\}\))(?![\w$])",
        lambda m: "$" + esc_dollars(m.group(1)) + "$",
        s,
    )
    # Empty spans left where a swallowed value used to be.
    s = s.replace("$$", "")
    # Dollar parity per line: an odd unescaped $ invites cross-line pairing.
    s = "\n".join(balance_line(ln) for ln in s.split("\n"))
    return s


def _extract_definition(el, md_parts):
    md_parts.append("> **Definition**")
    for child in el:
        tag = local(child.tag)
        if tag == "statement":
            t = text_content(child)
            if t:
                md_parts.append(f"> {t}")
        elif tag == "notation":
            usage = child.find("usage")
            desc = child.find("description")
            usage_text = text_content(usage) if usage is not None else ""
            desc_text = text_content(desc) if desc is not None else ""
            if usage_text and desc_text:
                md_parts.append(f"> {usage_text}: {desc_text}")
    md_parts.append("")


def _extract_fact(el, md_parts):
    title_el = el.find("title")
    if title_el is not None:
        md_parts.append(f"**{norm(''.join(title_el.itertext()))}**")
    stmt = el.find("statement")
    if stmt is not None:
        t = text_content(stmt)
        if t:
            md_parts.append(t)
    md_parts.append("")


def _extract_example(el, md_parts):
    title_el = el.find("title")
    tagline = f"**{norm(''.join(title_el.itertext()))}**" if title_el is not None else "**Example**"
    md_parts.append(tagline)

    stmt = el.find("statement")
    if stmt is not None:
        t = text_content(stmt)
        if t:
            md_parts.append(t)

    sol = el.find("solution")
    if sol is not None:
        md_parts.append("")
        md_parts.append("*Solution*")
        for child in sol:
            emit_block(child, md_parts)
    md_parts.append("")


def _extract_exercise(el, md_parts):
    """Extract exercise blocks. Many ORCCA exercises are WeBWorK-based.
    We extract the statement and solution text only."""
    title_el = el.find("title")
    intro = el.find("introduction")
    if title_el is not None:
        md_parts.append(f"**{norm(''.join(title_el.itertext()))}**")
    if intro is not None:
        t = text_content(intro)
        if t:
            md_parts.append(t)

    # Handle <webwork> nested exercises
    webwork = el.findall(".//webwork")
    for ww in webwork:
        for task in ww.findall("task"):
            stmt = task.find("statement")
            sol = task.find("solution")
            if stmt is not None:
                t = text_content(stmt, inline_md=True)
                if t:
                    md_parts.append(t)
                    md_parts.append("")
            if sol is not None:
                s = text_content(sol, inline_md=True)
                if s:
                    md_parts.append("*Solution*")
                    md_parts.append(s)
                    md_parts.append("")

    # Also handle inline exercises
    for task in el.findall("task"):
        stmt = task.find("statement")
        sol = task.find("solution")
        if stmt is not None:
            t = text_content(stmt, inline_md=True)
            if t:
                md_parts.append(t)
                md_parts.append("")
        if sol is not None:
            s = text_content(sol, inline_md=True)
            if s:
                md_parts.append("*Solution*")
                md_parts.append(s)
                md_parts.append("")

    md_parts.append("")


def _extract_warning(el, md_parts):
    title_el = el.find("title")
    if title_el is not None:
        md_parts.append(f"> **{norm(''.join(title_el.itertext()))}**")
    for child in el:
        if local(child.tag) == "p":
            t = text_content(child)
            if t:
                md_parts.append(f"> {t}")
    md_parts.append("")


def _extract_aside(el, md_parts):
    title_el = el.find("title")
    if title_el is not None:
        md_parts.append(f"*{norm(''.join(title_el.itertext()))}*")
    for child in el:
        if local(child.tag) == "p":
            t = text_content(child)
            if t:
                md_parts.append(t)
    md_parts.append("")


def _extract_tabular(el, md_parts):
    rows = el.findall(".//row")
    if not rows:
        return
    table = []
    for row in rows:
        cells = []
        for c in row.findall("cell"):
            t = text_content(c)
            cells.append(t)
        table.append("| " + " | ".join(cells) + " |")
    if table:
        md_parts.extend(table)
        if len(table) > 1:
            cols = table[0].count("|") - 1
            md_parts.insert(len(md_parts) - len(table) + 1, "| " + " | ".join(["---"] * cols) + " |")
        md_parts.append("")


def _extract_sidebyside(el, md_parts):
    """Side-by-side content - extract text from each panel."""
    for panel in el.findall("panel"):
        for child in panel:
            tag = local(child.tag)
            if tag == "figure":
                _skip_figure(child, md_parts)
            else:
                emit_block(child, md_parts)


def _skip_figure(el, md_parts):
    caption = el.find("caption")
    if caption is not None:
        t = text_content(caption)
        # Skip captions of purely interactive elements (videos, widgets).
        if t and t.lower() not in ("alternative video lesson", "interactive"):
            md_parts.append(f"*{t}*")
            md_parts.append("")


def main():
    with open("data/lessons/mapping.json") as f:
        mapping = json.load(f)

    orcca_concepts = {k: v for k, v in mapping.items() if v.get("source") == "orcca"}

    # Group by ORCCA slug to fetch each section once
    slugs = {}
    for cid, info in orcca_concepts.items():
        slug = info["slug"]
        slugs.setdefault(slug, []).append(cid)

    print(f"ORCCA sections to extract: {len(slugs)}")

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    for slug, concept_ids in slugs.items():
        print(f"\nProcessing section: {slug}")
        print(f"  Concepts: {', '.join(concept_ids)}")
        try:
            md = extract_section(slug)
            attribution = (
                "> Content sourced from [ORCCA](https://pcc.edu/orcca) "
                "(Open Resources for Community College Algebra) — CC BY 4.0\n"
            )
            md = attribution + md

            for cid in concept_ids:
                out_path = OUTPUT_DIR / f"{cid}.md"
                out_path.write_text(md)
                print(f"  -> Wrote {out_path}")

        except Exception as e:
            print(f"  X Error: {e}")


if __name__ == "__main__":
    main()
