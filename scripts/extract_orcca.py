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


def text_content(el):
    """Recursively extract text from element, handling math and formatting."""
    if el is None:
        return ""
    parts = []
    tag = el.tag.split("}")[-1] if "}" in str(el.tag) else el.tag
    if tag == "m":
        inner = "".join(el.itertext())
        return f"${inner}$"
    if tag == "me":
        inner = "".join(el.itertext())
        return f"\\[{inner}\\]"
    if tag == "term":
        inner = "".join(el.itertext())
        return f"*{inner.strip()}*"
    if tag == "em":
        inner = "".join(el.itertext())
        return f"*{inner.strip()}*"
    if tag == "q":
        inner = "".join(el.itertext())
        return f"\"{inner.strip()}\""
    if tag == "idx":
        return ""
    if tag == "url":
        return el.get("href", "")
    if tag == "c":
        return f"`{''.join(el.itertext()).strip()}`"
    if tag in ("figure", "caption", "tabular"):
        return ""

    if el.text:
        parts.append(el.text.strip())
    for child in el:
        child_text = text_content(child)
        if child_text:
            parts.append(child_text)
        if child.tail:
            parts.append(child.tail.strip())
    return " ".join(p for p in parts if p).strip()


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
            tag = child.tag.split("}")[-1] if "}" in str(child.tag) else child.tag
            if tag == "p":
                t = text_content(child)
                if t:
                    md_parts.append(t)
                    md_parts.append("")
            elif tag == "figure":
                _skip_figure(child, md_parts)

    # Subsections
    for subsection in section.findall("subsection"):
        sub_title_el = subsection.find("title")
        if sub_title_el is not None:
            md_parts.append(f"## {''.join(sub_title_el.itertext())}")
            md_parts.append("")

        for child in subsection:
            tag = child.tag.split("}")[-1] if "}" in str(child.tag) else child.tag
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
                for row in child:
                    rtag = row.tag.split("}")[-1] if "}" in str(row.tag) else row.tag
                    if rtag == "mrow":
                        t = text_content(row)
                        if t:
                            md_parts.append(t)
                            md_parts.append("")

    return "\n".join(md_parts)


def _extract_definition(el, md_parts):
    md_parts.append("> **Definition**")
    for child in el:
        tag = child.tag.split("}")[-1] if "}" in str(child.tag) else child.tag
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
        md_parts.append(f"**{''.join(title_el.itertext())}**")
    stmt = el.find("statement")
    if stmt is not None:
        t = text_content(stmt)
        if t:
            md_parts.append(t)
    md_parts.append("")


def _extract_example(el, md_parts):
    title_el = el.find("title")
    tagline = f"**{''.join(title_el.itertext())}**" if title_el is not None else "**Example**"
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
            ctag = child.tag.split("}")[-1] if "}" in str(child.tag) else child.tag
            if ctag == "p":
                t = text_content(child)
                if t:
                    md_parts.append(t)
    md_parts.append("")


def _extract_exercise(el, md_parts):
    """Extract exercise blocks. Many ORCCA exercises are WeBWorK-based.
    We extract the statement and solution text only."""
    title_el = el.find("title")
    intro = el.find("introduction")
    if title_el is not None:
        md_parts.append(f"**{''.join(title_el.itertext())}**")
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
                t = text_content(stmt)
                if t:
                    md_parts.append(f"- {t}")
            if sol is not None:
                s = text_content(sol)
                if s:
                    md_parts.append(f"  *Solution*: {s}")

    # Also handle inline exercises
    for task in el.findall("task"):
        stmt = task.find("statement")
        sol = task.find("solution")
        if stmt is not None:
            t = text_content(stmt)
            if t:
                md_parts.append(f"- {t}")
        if sol is not None:
            s = text_content(sol)
            if s:
                md_parts.append(f"  *Solution*: {s}")

    md_parts.append("")


def _extract_warning(el, md_parts):
    title_el = el.find("title")
    if title_el is not None:
        md_parts.append(f"> **{''.join(title_el.itertext())}**")
    for child in el:
        ctag = child.tag.split("}")[-1] if "}" in str(child.tag) else child.tag
        if ctag == "p":
            t = text_content(child)
            if t:
                md_parts.append(f"> {t}")
    md_parts.append("")


def _extract_aside(el, md_parts):
    title_el = el.find("title")
    if title_el is not None:
        md_parts.append(f"*{''.join(title_el.itertext())}*")
    for child in el:
        ctag = child.tag.split("}")[-1] if "}" in str(child.tag) else child.tag
        if ctag == "p":
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
            ctag = child.tag.split("}")[-1] if "}" in str(child.tag) else child.tag
            if ctag == "p":
                t = text_content(child)
                if t:
                    md_parts.append(t)


def _skip_figure(el, md_parts):
    caption = el.find("caption")
    if caption is not None:
        t = text_content(caption)
        if t:
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
