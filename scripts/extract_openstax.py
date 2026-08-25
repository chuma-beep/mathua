"""
Extract lesson content from OpenStax Prealgebra HTML pages.
"""
import json
import re
import requests
import html2text
from bs4 import BeautifulSoup, NavigableString, Tag
from pathlib import Path

OPENSTAX_BASE = "https://openstax.org/books/prealgebra/pages"
OUTPUT_DIR = Path("data/lessons/teaching")

# MathML operator → LaTeX (identity for anything unmapped).
MO_MAP = {
    "·": "\\cdot",
    "×": "\\times",
    "÷": "\\div",
    "−": "-",
    "–": "-",
    "≠": "\\neq",
    "≤": "\\leq",
    "≥": "\\geq",
    "±": "\\pm",
    "≈": "\\approx",
    "→": "\\to",
    "⁄": "/",
    "‾": "",
}


def _tex_children(el):
    return "".join(mathml_to_latex(c) for c in el.children)


def mathml_to_latex(el):
    """Convert an OpenStax MathML tree to a LaTeX string."""
    if isinstance(el, NavigableString):
        return str(el)
    if not isinstance(el, Tag):
        return ""
    tag = el.name

    if tag in ("annotation", "annotation-xml"):
        # Semantic duplicate of the presentation row above it.
        return ""
    if tag in ("semantics", "mrow", "mstyle", "mpadded", "mphantom", "math"):
        inner = _tex_children(el).strip()
        return inner
    if tag == "mi" or tag == "mn":
        return el.get_text()
    if tag == "mo":
        text = el.get_text()
        mapped = MO_MAP.get(text)
        # Trailing space terminates multi-letter commands (\cdot + c).
        return mapped + " " if mapped is not None else text
    if tag == "mtext":
        text = el.get_text().strip()
        for ch, tex in MO_MAP.items():
            text = text.replace(ch, tex)
        return f"\\text{{{text}}}"
    if tag == "mspace":
        return "\\ "
    if tag == "mfrac":
        parts = [c for c in el.children if isinstance(c, Tag)]
        num = mathml_to_latex(parts[0]) if len(parts) > 0 else ""
        den = mathml_to_latex(parts[1]) if len(parts) > 1 else ""
        return f"\\frac{{{num}}}{{{den}}}"
    if tag == "msqrt":
        return f"\\sqrt{{{_tex_children(el)}}}"
    if tag == "mroot":
        parts = [c for c in el.children if isinstance(c, Tag)]
        base = mathml_to_latex(parts[0]) if len(parts) > 0 else ""
        idx = mathml_to_latex(parts[1]) if len(parts) > 1 else ""
        return f"\\sqrt[{idx}]{{{base}}}"
    if tag == "msup":
        parts = [c for c in el.children if isinstance(c, Tag)]
        base = mathml_to_latex(parts[0]) if len(parts) > 0 else ""
        exp = mathml_to_latex(parts[1]) if len(parts) > 1 else ""
        return f"{base}^{{{exp}}}"
    if tag == "msub":
        parts = [c for c in el.children if isinstance(c, Tag)]
        base = mathml_to_latex(parts[0]) if len(parts) > 0 else ""
        sub = mathml_to_latex(parts[1]) if len(parts) > 1 else ""
        return f"{base}_{{{sub}}}"
    if tag in ("msubsup", "munderover"):
        parts = [c for c in el.children if isinstance(c, Tag)]
        bits = [mathml_to_latex(p) for p in parts[:3]]
        while len(bits) < 3:
            bits.append("")
        return f"{bits[0]}_{{{bits[1]}}}^{{{bits[2]}}}"
    if tag == "munder":
        parts = [c for c in el.children if isinstance(c, Tag)]
        base = mathml_to_latex(parts[0]) if len(parts) > 0 else ""
        return base
    if tag == "mtable":
        rows = []
        width = 0
        for tr in el.find_all("mtr", recursive=False):
            cells = [
                mathml_to_latex(td)
                for td in tr.find_all("mtd", recursive=False)
            ]
            width = max(width, len(cells))
            rows.append(cells)
        if not rows:
            return ""
        colspec = "c" * width
        body = " \\\\ ".join(
            " & ".join((row + [""] * width)[:width])
            for row in rows
        )
        return f"\\begin{{array}}{{{colspec}}}{body}\\end{{array}}"
    if tag == "mprescripts":
        return ""

    # Fallback: recurse.
    return _tex_children(el)


def clean_tex(tex):
    """Harden a converted LaTeX fragment for embedding in $...$."""
    # Trim spacing commands at the edges BEFORE whitespace collapse strips
    # their trailing space and leaves a dangling backslash.
    tex = re.sub(r"^(?:\\ |\\,|\\:|\\;)+", "", tex)
    tex = re.sub(r"(?:\\ |\\,|\\:|\\;)+$", "", tex)
    tex = re.sub(r"\s+", " ", tex).strip()
    # A literal $ inside a span breaks dollar pairing downstream.
    tex = tex.replace("$", "\\text{\\$}")
    # % starts a KaTeX comment; escape it.
    tex = tex.replace("%", "\\%")
    # Drop dangling subscript/superscript operators left by partial MathML.
    tex = re.sub(r"[_^]\s*(?=\\|$)", "", tex)
    # Fill-in-the-blank runs (__ / ___) are not subscripts.
    tex = re.sub(r"_{2,}", "\\;\\;\\;", tex)
    tex = re.sub(r"[_](?![{A-Za-z0-9])", "\\_", tex)
    # Spacing commands glued to a closing brace escape it (\frac{... \}).
    tex = re.sub(r"(?:\\ |\\,|\\:|\\;)+\}", "}", tex)
    # Junk fence macros leaked from OpenStax markup.
    tex = re.sub(r"\\\\(?:maht|next)(?:open|close)\\b", "", tex)
    # Combining solidus left over from decomposed negation glyphs.
    tex = tex.replace("\u0338", "")
    tex = tex.rstrip("\\").lstrip()
    return tex


def inline_math(soup_node):
    """Replace every <math> element with a placeholder token, protecting the
    converted LaTeX from html2text's escaping."""
    spans = []
    for el in soup_node.find_all("math"):
        tex = clean_tex(mathml_to_latex(el))
        spans.append(tex)
        el.replace_with(NavigableString(f"\ue000{len(spans) - 1}\ue001"))
    return spans


def restore_math(md, spans):
    def sub(m):
        i = int(m.group(1))
        if i >= len(spans):
            return ""
        tex = spans[i]
        return f"${tex}$" if tex else ""

    return re.sub("\ue000(\\d+)\ue001", sub, md)


def repair_group_braces(inner):
    """Treat '\\}' as a group-closing brace whenever a group is actually open.
    OpenStax authors write \\frac{...\\} where the backslash-space padding lost
    its space; real escaped-literal braces only occur at depth 0."""
    out = []
    depth = 0
    i = 0
    n = len(inner)
    while i < n:
        c = inner[i]
        if c == "\\" and i + 1 < n:
            nxt = inner[i + 1]
            if nxt == "}":
                if depth > 0:
                    depth -= 1
                    out.append("}")
                else:
                    out.append("\\}")
                i += 2
                continue
            if nxt == "{":
                out.append("\\{")
                i += 2
                continue
            out.append(inner[i : i + 2])
            i += 2
            continue
        if c == "{":
            depth += 1
        elif c == "}":
            depth -= 1
        out.append(c)
        i += 1
    return "".join(out)


def polish_spans(md):
    """Repair author-typed $...$ spans that never went through clean_tex."""
    def fix(m):
        inner = m.group(1)
        # Spacing commands glued to a closing brace escape it.
        inner = re.sub(r"(?:\\ |\\,|\\:|\\;)+\}", "}", inner)
        # Junk fence macros leaked from OpenStax markup.
        inner = re.sub(r"\\(?:maht|next)(?:open|close)\b", "", inner)
        inner = inner.replace("\u0338", "")
        inner = repair_group_braces(inner)
        return f"${inner}$"

    return re.sub(r"\$([^$\n]+)\$", fix, md)


def extract_section(slug, book="prealgebra"):
    """Fetch and extract an OpenStax section, returning markdown."""
    url = f"https://openstax.org/books/{book}/pages/{slug}"
    print(f"  Fetching {url}")
    resp = requests.get(url, timeout=30)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.content, "html.parser")
    main = soup.find("main")
    if not main:
        raise ValueError("No <main> tag found")

    content = main.find("div", class_="chapter-content-module")
    if not content:
        raise ValueError("No .chapter-content-module found")

    # Remove unwanted elements
    for tag in content.find_all(["script", "style", "nav", "aside", "noscript"]):
        tag.decompose()

    # Collapse duplicated MathML into placeholder tokens before html2text runs.
    spans = inline_math(content)

    # Convert to markdown
    h = html2text.HTML2Text()
    h.body_width = 0
    h.ignore_links = True
    h.ignore_images = True
    h.ignore_emphasis = False
    h.ignore_tables = False

    md = h.handle(str(content))

    # Clean up common issues: excessive blank lines and stray artifacts.
    lines = md.split("\n")
    cleaned = []
    prev_blank = False
    for line in lines:
        stripped = line.strip()
        if not stripped:
            if prev_blank is False:
                cleaned.append("")
                prev_blank = True
        else:
            cleaned.append(line)
            prev_blank = False

    md = "\n".join(cleaned)
    md = restore_math(md, spans)
    md = polish_spans(md)
    return md


def main():
    with open("data/lessons/mapping.json") as f:
        mapping = json.load(f)

    openstax_concepts = {k: v for k, v in mapping.items() if v.get("source") == "openstax"}

    slugs = {}
    for cid, info in openstax_concepts.items():
        slug = info["slug"]
        slugs.setdefault(slug, []).append(cid)

    print(f"OpenStax sections to extract: {len(slugs)}")

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    for slug, concept_ids in slugs.items():
        print(f"\nProcessing section: {slug}")
        print(f"  Concepts: {', '.join(concept_ids)}")
        try:
            md = extract_section(slug)
            attribution = (
                "> Content sourced from [OpenStax Prealgebra]"
                "(https://openstax.org/books/prealgebra/pages/1-introduction) "
                "by Marecek & Anthony-Smith — CC BY 4.0\n"
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
