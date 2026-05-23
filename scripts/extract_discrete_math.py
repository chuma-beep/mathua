"""
Extract lesson content from Discrete Mathematics: An Open Introduction, 3rd ed (CC BY-SA 4.0)
"""
import json
import re
import html
from pathlib import Path

import requests
from bs4 import BeautifulSoup

BASE_URL = "https://discrete.openmathbooks.org/dmoi3"
OUTPUT_DIR = Path("data/lessons/teaching")

SLUG_MAP = {
    "discrete.logic.connectives": ("sec_propositional", "Propositional Logic"),
    "discrete.logic.truth_tables": ("sec_propositional", "Truth Tables"),
    "discrete.logic.quantifiers": ("sec_intro-statements", "Predicates and Quantifiers"),
    "discrete.sets.venn": ("sec_intro-sets", "Venn Diagrams"),
    "discrete.graphs.basics": ("sec_gt-intro", "Graph Theory Basics"),
    "discrete.graphs.paths": ("sec_paths", "Euler Paths and Circuits"),
    "discrete.graphs.trees": ("sec_trees", "Trees"),
    "discrete.recurrence": ("sec_recurrence", "Solving Recurrence Relations"),
}


def fetch_section(slug):
    """Fetch a section page and extract the main content as markdown."""
    url = f"{BASE_URL}/{slug}.html"
    r = requests.get(url, timeout=30)
    r.raise_for_status()
    soup = BeautifulSoup(r.text, "html.parser")

    # Remove navigation, sidebar, header, footer
    for cls in ["navigation", "header", "footer", "sidebar"]:
        for el in soup.find_all(class_=cls):
            el.decompose()

    main = soup.find("main") or soup.find("article") or soup.find("div", class_="main-content")
    if not main:
        main = soup.body

    return soup, main


def extract_section(soup, main, section_title):
    """Extract content for a specific section/subsection."""
    parts = [f"# {section_title}\n"]

    # Process the main content
    for el in main.find_all(["section", "div", "p", "h1", "h2", "h3", "h4", "h5", "h6",
                              "ul", "ol", "li", "pre", "blockquote", "table", "figure",
                              "dl", "dt", "dd", "div", "span", "article"]):
        if el.find_parent(["header", "nav", "footer"]):
            continue

        # Convert HTML to simple markdown
        tag = el.name
        text = el.get_text().strip()
        if not text and tag not in ("img", "figure", "table"):
            continue

        if tag in ("h1", "h2", "h3", "h4", "h5", "h6"):
            level = int(tag[1])
            parts.append(f"{'#' * level} {text}")
            parts.append("")
        elif tag == "p":
            parts.append(text)
            parts.append("")
        elif tag == "li":
            parent = el.find_parent("ol")
            if parent:
                parts.append(f"1. {text}")
            else:
                parts.append(f"- {text}")
        elif tag == "pre":
            code = el.get_text()
            parts.append(f"```\n{code}\n```")
            parts.append("")
        elif tag == "blockquote":
            parts.append(f"> {text}")
            parts.append("")
        elif tag == "table":
            rows = el.find_all("tr")
            md_rows = []
            for i, row in enumerate(rows):
                cells = row.find_all(["td", "th"])
                md_row = "| " + " | ".join(c.get_text().strip() for c in cells) + " |"
                md_rows.append(md_row)
                if i == 0 and row.find_all("th"):
                    sep = "| " + " | ".join("---" for _ in cells) + " |"
                    md_rows.append(sep)
            for r in md_rows:
                parts.append(r)
            parts.append("")

    return "\n".join(parts).strip()


def main():
    with open("data/lessons/mapping.json") as f:
        mapping = json.load(f)

    dm_concepts = {k: v for k, v in mapping.items() if v.get("source") == "discrete_math"}
    print(f"Discrete Math concepts: {len(dm_concepts)}")

    # Group by URL slug
    slug_groups = {}
    for cid in sorted(dm_concepts.keys()):
        if cid not in SLUG_MAP:
            print(f"  SKIP {cid}: no slug mapping")
            continue
        url_slug, section_name = SLUG_MAP[cid]
        slug_groups.setdefault(url_slug, []).append((cid, section_name))

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    for url_slug, concepts in slug_groups.items():
        print(f"\nFetching {url_slug}...")
        try:
            soup, main = fetch_section(url_slug)
        except Exception as e:
            print(f"  FAILED: {e}")
            continue

        for cid, section_name in concepts:
            print(f"  Extracting {cid} ({section_name})...")
            md = extract_section(soup, main, section_name)
            attribution = (
                f"> Content sourced from [Discrete Mathematics: An Open Introduction, 3rd ed]({BASE_URL}/{url_slug}.html)"
                " by Oscar Levin — CC BY-SA 4.0\n\n"
            )
            md = attribution + md

            out_path = OUTPUT_DIR / f"{cid}.md"
            out_path.write_text(md)
            print(f"    -> Wrote {out_path} ({len(md)} chars)")


if __name__ == "__main__":
    main()
