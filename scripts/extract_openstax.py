"""
Extract lesson content from OpenStax Prealgebra 1e HTML pages.
"""
import json
import requests
import html2text
from bs4 import BeautifulSoup
from pathlib import Path

OPENSTAX_BASE = "https://openstax.org/books/prealgebra/pages"
OUTPUT_DIR = Path("data/lessons/teaching")


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

    # Convert to markdown
    h = html2text.HTML2Text()
    h.body_width = 0
    h.ignore_links = True
    h.ignore_images = True
    h.ignore_emphasis = False
    h.ignore_tables = False

    md = h.handle(str(content))

    # Clean up common issues
    # Remove duplicate spacing around math numbers like "33" -> "3 3" is wrong but we keep as-is
    # Remove excessive blank lines
    lines = md.split("\n")
    cleaned = []
    prev_blank = False
    for line in lines:
        stripped = line.strip()
        if not stripped:
            if not prev_blank:
                cleaned.append("")
                prev_blank = True
        else:
            cleaned.append(line)
            prev_blank = False

    return "\n".join(cleaned)


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
                "> Content sourced from [OpenStax Prealgebra 1e]"
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
