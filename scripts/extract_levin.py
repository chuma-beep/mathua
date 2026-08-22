"""
Extract lesson content from Discrete Mathematics: An Open Introduction
(Oscar Levin, discrete.openmathbooks.org) — CC BY-SA 4.0.
"""
import requests
import html2text
from bs4 import BeautifulSoup
from pathlib import Path

BASE = "https://discrete.openmathbooks.org/dmoi3"
OUTPUT_DIR = Path("data/lessons/teaching")

ATTRIBUTION = (
    "> Content sourced from [Discrete Mathematics: An Open Introduction, 3e]"
    "({page_url}) by Oscar Levin — CC BY-SA 4.0\n"
)


def extract_section(slug):
    """Fetch a Levin book section and return markdown of its main content."""
    url = f"{BASE}/{slug}"
    print(f"  Fetching {url}")
    resp = requests.get(url, timeout=30)
    resp.raise_for_status()

    soup = BeautifulSoup(resp.content, "html.parser")
    main = soup.select_one(".ptx-content") or soup.find("main")
    if not main:
        raise ValueError("No content container found")

    for tag in main.find_all(["script", "style", "nav", "aside", "noscript"]):
        tag.decompose()
    # Drop exercise/solution blocks: lessons teach, they don't assign homework.
    for cls in ("exercises", "solution", "hint", "answer", "colophon"):
        for tag in main.find_all(class_=cls):
            tag.decompose()

    h = html2text.HTML2Text()
    h.body_width = 0
    h.ignore_links = True
    h.ignore_images = True
    h.ignore_emphasis = False

    md = h.handle(str(main))

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

    return ATTRIBUTION.format(page_url=url) + "\n".join(cleaned)


if __name__ == "__main__":
    import sys
    targets = {
        "discrete.sequences.recurrence": "sec_recurrence.html",
        "nt.adv.diophantine": "sec_addtops-numbth.html",
    }
    only = sys.argv[1:] or targets.keys()
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    for cid in only:
        slug = targets[cid]
        md = extract_section(slug)
        out = OUTPUT_DIR / f"{cid}.md"
        out.write_text(md)
        print(f"  -> Wrote {out} ({len(md)} chars)")
