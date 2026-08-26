"""
Ingest OpenStax Prealgebra pages into CommonMark.

Fetches the page, isolates the content module, unwraps presentation-only
div/span wrappers, then delegates conversion to pandoc (whose AST carries
math as typed nodes; its MathML->TeX reader replaces hand-rolled converters).
"""
import json
import re
import subprocess
import sys
from pathlib import Path

import requests
from bs4 import BeautifulSoup

OPENSTAX_BASE = "https://openstax.org/books/prealgebra/pages"
OUTPUT_DIR = Path("data/lessons/teaching")
ROOT = Path(__file__).resolve().parents[2]
MAPPING_PATH = ROOT / "data/lessons/mapping.json"
ATTRIBUTION = (
    "> Content sourced from [OpenStax Prealgebra]"
    "(https://openstax.org/books/prealgebra/pages/1-introduction) "
    "by Marecek & Anthony-Smith — CC BY 4.0\n"
)


def flatten_html(html: str) -> str:
    """Strip chrome and presentation wrappers so pandoc sees clean structure."""
    soup = BeautifulSoup(html, "html.parser")
    main = soup.find("main") or soup
    content = main.find("div", class_="chapter-content-module")
    if content is None:
        raise ValueError("no chapter-content-module found")
    for t in content.find_all(["script", "style", "nav", "aside", "noscript"]):
        t.decompose()
    for t in content.find_all(["video", "svg", "img", "button", "input"]):
        t.decompose()
    for name in ("div", "span", "a", "section", "main", "article", "header", "footer"):
        for t in content.find_all(name):
            t.unwrap()
    for t in content.find_all("figure"):
        cap = t.find("figcaption")
        if cap is not None:
            p = soup.new_tag("p")
            p.string = f"*{cap.get_text(' ', strip=True)}*"
            t.replace_with(p)
        else:
            t.unwrap()
    return str(content)


UNICODE_OPS = {
    "·": "\\cdot", "×": "\\times", "÷": "\\div", "−": "-", "–": "-",
    "≠": "\\neq", "≤": "\\le", "≥": "\\ge", "≈": "\\approx", "±": "\\pm",
    "⁄": "/", "′": "'", "″": "''",
}


def fix_math_ops(md: str) -> str:
    """Map unicode operators pandoc left as \\operatorname{..} to LaTeX."""

    def repl(m):
        ch = m.group(1)
        return UNICODE_OPS.get(ch, f"\\text{{{ch}}}")

    md = re.sub(r"\\operatorname\{([^\n{}]{1,4})\}", repl, md)

    def inner(m):
        body = m.group(0)[1:-1]
        for ch, tex in UNICODE_OPS.items():
            body = body.replace(ch, tex)
        return "$" + body + "$"

    return re.sub(r"\$[^$\n]+\$", inner, md)


def to_markdown(html: str) -> str:
    out = subprocess.run(
        [PANDOC, "-f", "html", "-t", "gfm+tex_math_dollars", "--wrap=none"],
        input=html.encode("utf-8"),
        capture_output=True,
        check=True,
    )
    md = out.stdout.decode("utf-8")
    # gfm math uses $`...`$; our renderer speaks plain $...$.
    md = md.replace("$`", "$").replace("`$", "$")
    # gfm renders display math as ```math fences; our renderer speaks $$.
    # Line-anchored: a plain ``` code fence must not close a math fence.
    lines = md.split("\n")
    out_l: list[str] = []
    in_math_fence = False
    fence_body: list[str] = []
    for ln in lines:
        if not in_math_fence and re.fullmatch(r"```\s*math\s*", ln):
            in_math_fence = True
            fence_body = []
            continue
        if in_math_fence and re.fullmatch(r"```\s*", ln):
            out_l += ["$$", "\n".join(fence_body).strip(), "$$", ""]
            in_math_fence = False
            continue
        (fence_body if in_math_fence else out_l).append(ln)
    if in_math_fence:
        out_l += ["$$"] + fence_body + ["$$"]
    md = "\n".join(out_l)

    # OpenStax ships some pre-wrapped spans: <span class="math inline">x</span>
    md = re.sub(r'<span class="math inline">(.*?)</span>', lambda m: "$" + m.group(1) + "$", md, flags=re.S)
    # KaTeX has no \mspace; spacing is inessential.
    md = re.sub(r"\\mspace\{[^{}]*\}", "", md)
    md = fix_math_ops(md)
    md = md.replace('\\*', '*')
    md = re.sub(r"\\\\\*", "*", md)
    # Unicode spaces inside math confuse pairing and KaTeX.
    md = re.sub(r"[\u00a0\u2000-\u200a\u202f\u2005]", " ", md)
    # Spacing-only \operatorname wrappers: keep the spacing, drop the wrapper.
    md = re.sub(r"\\operatorname\{((?:\\;|\\_|\\,|\s)*)\}", r"\1", md)

    def balance_line(ln):
        positions = [mm.start() for mm in re.finditer(r"(?<!\\)\$", ln)]
        if len(positions) % 2 == 1:
            i = positions[-1]
            ln = ln[:i] + "\\$" + ln[i + 1:]
        return ln

    def trim_span(m):
        inner = m.group(1)
        inner = re.sub(r"^(?:\\ |\\,|\\:|\\;)+", "", inner)
        inner = re.sub(r"(?:\\ |\\,|\\:|\\;)+$", "", inner)
        # Fill-in-the-blank underscores are not subscripts.
        inner = re.sub(r"(?<!\\)_\{2,\}", "\\;\\;", inner)
        inner = re.sub(r"(?<!\\)_(?!\\_|[{A-Za-z0-9])", "\\_", inner)
        inner = re.sub(r"\\_\\_", "\\;\\;", inner)
        return "$" + inner + "$"

    md = re.sub(r"\$([^$\n]+)\$", trim_span, md)

    # Pandoc emits display blocks as standalone $$ lines; within a block,
    # unwrap any nested $inline$ fragments so dollar pairing stays flat.
    out_lines: list[str] = []
    in_block = False
    block: list[str] = []
    for ln in md.split("\n"):
        if ln.strip() == "$$":
            if not in_block:
                in_block = True
                block = []
            else:
                inner = "\n".join(block)
                inner = re.sub(r"\$([^$\n]*)\$", r"\1", inner)
                out_lines.append("$$")
                if inner.strip():
                    out_lines.append(inner.strip())
                out_lines.append("$$")
                in_block = False
                block = []
            continue
        (block if in_block else out_lines).append(ln)
    if in_block:
        out_lines.extend(block)
    md = "\n".join(out_lines)

    md = "\n".join(balance_line(l) for l in md.split("\n"))
    return md


def main() -> int:
    global PANDOC
    PANDOC = "/usr/bin/pandoc"
    for candidate in ("/usr/local/bin/pandoc", str(Path.home() / ".local/bin/pandoc")):
        if Path(candidate).exists():
            PANDOC = candidate
            break
    if not Path(PANDOC).exists():
        print("pandoc not found", file=sys.stderr)
        return 1

    mapping = json.loads(MAPPING_PATH.read_text())
    slugs: dict[str, list[str]] = {}
    for cid, info in mapping.items():
        if info.get("source") == "openstax":
            slugs.setdefault(info["slug"], []).append(cid)

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    failures = 0
    print(f"OpenStax sections to extract: {len(slugs)}")
    for slug, concept_ids in sorted(slugs.items()):
        try:
            resp = requests.get(f"{OPENSTAX_BASE}/{slug}", timeout=30)
            resp.raise_for_status()
            flat = flatten_html(resp.content.decode("utf-8"))
            md = to_markdown(flat)
            md = re.sub(r"\n{3,}", "\n\n", md).strip() + "\n"
            body = ATTRIBUTION + "\n" + md
            for cid in concept_ids:
                (OUTPUT_DIR / f"{cid}.md").write_text(body)
            print(f"ok {slug} ({len(concept_ids)})")
        except Exception as exc:  # noqa: BLE001 - report and continue
            failures += 1
            print(f"X {slug}: {exc}")
    return 1 if failures else 0


if __name__ == "__main__":
    sys.exit(main())
