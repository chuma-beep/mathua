#!/usr/bin/env python3
"""
Scrape Algebrica (https://algebrica.org) for math lesson content not yet on GitHub.

Usage:  python3 scripts/scrape_algebrica.py
Output: data/lessons/algebrica/{category}/*.md

Content is CC BY-NC 4.0 (https://creativecommons.org/licenses/by-nc/4.0/).
"""

import os
import re
import time
import requests
from bs4 import BeautifulSoup
import html2text

BASE = "https://algebrica.org"
OUT = os.path.join(os.path.dirname(__file__), "..", "data", "lessons", "algebrica")
os.makedirs(OUT, exist_ok=True)

# Categories we already fully pulled from GitHub — skip
ALREADY_FULL = {
    "algebraic-structures", "complex-numbers", "linear-systems",
    "polynomials", "powers-radicals-logarithms",
    "trigonometry",
    "numbers",        # "Sets and Numbers"
    "matrices",       # "Vectors and Matrices"
}

# Map Algebrica category slugs to our lesson subdirectory names
CATEGORY_DIRS = {
    "functions": "functions",
    "derivatives": "derivatives",
    "differential-calculus-theorems": "differential-calculus-theorems",
    "differential-equations": "differential-equations",
    "inequalities": "inequalities",
    "lines-planes-conic-sections": "lines-planes-conic-sections",
    "probability-and-statistics": "probability-and-statistics",
    "sequences": "sequences",
    "series": "series",
    "various": "various",
    "kinematics": "kinematics",
    "physics": "physics",
    "numbers": "sets-and-numbers",
    "matrices": "vectors-and-matrices",
}

# Categories where we have some files but might have new ones
PARTIAL_CATEGORIES = {
    "equations": "equations",
    "integrals": "integrals",
    "limits": "limits",
}

HEADERS = {"User-Agent": "mathua-scraper/1.0 (+https://github.com/chuma-beep/mathua)"}
session = requests.Session()
session.headers.update(HEADERS)


def fetch(url):
    resp = session.get(url, timeout=30)
    resp.raise_for_status()
    return resp.text


def slug_from_url(url):
    return url.rstrip("/").split("/")[-1]


def get_category_entries(category_slug):
    entries = []
    page_url = f"{BASE}/category/{category_slug}/"
    while page_url:
        html = fetch(page_url)
        soup = BeautifulSoup(html, "html.parser")
        for item in soup.select(".search-results__item"):
            link = item.select_one(".search-results__title a")
            if link and link.get("href"):
                entries.append(link["href"])
        next_link = soup.select_one(".pb-browse-next")
        disabled = "is-disabled" in next_link.get("class", []) if next_link else True
        page_url = next_link["href"] if next_link and not disabled else None
    return entries


def strip_metadata(text):
    """Remove the 'Concept / Intermediate / Requires / Enables' metadata block at the bottom."""
    # The metadata section starts after patterns like "Concept  Intermediate  Requires  Enables"
    # Look for the "## Concept" or "Concept" heading followed by metadata
    idx = text.find("\n## Concept\n")
    if idx == -1:
        idx = text.find("\nConcept\n")
    if idx == -1:
        idx = text.find("Intermediate")
    if idx > 0:
        # Check if this is really the metadata section (followed by numbers)
        snippet = text[idx:idx+100]
        if re.search(r'(Intermediate|Requires|Enables)', snippet):
            text = text[:idx].rstrip()
    return text


def clean_markdown(md):
    """Post-process the html2text markdown output."""
    # Remove horizontal rules (---, ***, * * *)
    md = re.sub(r'^\* \* \*\n', '', md, flags=re.MULTILINE)
    md = re.sub(r'^\*{3,}\n', '', md, flags=re.MULTILINE)
    md = re.sub(r'^-{3,}\n', '', md, flags=re.MULTILINE)

    # Remove angle brackets around link URLs: [text](<url>) → [text](url)
    md = re.sub(r'\(<(https?://[^>]+)>\)', r'(\1)', md)
    md = re.sub(r'\(<(/[^>]+)>\)', r'(\1)', md)

    # Remove the "Categories: ..." line if present
    md = re.sub(r'\* \[[^\]]+\]\([^)]+\)\n', '', md)

    # Fix display math that html2text sometimes breaks with extra spaces
    md = re.sub(r'\\\[\s+', r'\\[', md)
    md = re.sub(r'\s+\\\]', r'\\]', md)

    # Strip metadata section
    md = strip_metadata(md)

    # Remove excessive blank lines (3+ → 2)
    md = re.sub(r'\n{4,}', '\n\n\n', md)

    return md.strip() + "\n"


def scrape_entry(entry_url, cat_dir):
    slug = slug_from_url(entry_url)
    fname = slug + ".md"
    fpath = os.path.join(cat_dir, fname)

    if os.path.exists(fpath):
        print(f"  SKIP (exists): {slug}")
        return

    print(f"  Fetch: {entry_url}")
    html = fetch(entry_url)
    soup = BeautifulSoup(html, "html.parser")

    content = soup.select_one(".post-content")
    if not content:
        print(f"    WARNING: no .post-content found")
        return

    # Convert HTML to markdown
    h = html2text.HTML2Text()
    h.body_width = 0
    h.ignore_links = False
    h.ignore_images = False
    h.ignore_emphasis = False
    h.protect_links = True
    h.unicode_snob = True
    h.skip_internal_links = False
    h.bypass_tables = False

    md = h.handle(str(content))
    md = clean_markdown(md)

    # Prepend attribution
    header = f"> Content sourced from [Algebrica]({entry_url}) — CC BY-NC 4.0\n\n"
    md = header + md

    with open(fpath, "w") as f:
        f.write(md)
    print(f"    → Saved {len(md)} bytes to {fname}")


def main():
    scrape_categories = list(CATEGORY_DIRS.keys()) + list(PARTIAL_CATEGORIES.keys())

    for cat_slug in scrape_categories:
        if cat_slug in ALREADY_FULL and cat_slug not in PARTIAL_CATEGORIES:
            print(f"\n=== SKIP (already have all): {cat_slug} ===")
            continue

        dir_name = CATEGORY_DIRS.get(cat_slug, PARTIAL_CATEGORIES.get(cat_slug, cat_slug))
        cat_dir = os.path.join(OUT, dir_name)
        os.makedirs(cat_dir, exist_ok=True)

        print(f"\n=== {cat_slug} → {dir_name}/ ===")
        entries = get_category_entries(cat_slug)
        print(f"  Found {len(entries)} entries")

        for entry_url in entries:
            scrape_entry(entry_url, cat_dir)
            time.sleep(0.5)

    print("\nDone!")


if __name__ == "__main__":
    main()
