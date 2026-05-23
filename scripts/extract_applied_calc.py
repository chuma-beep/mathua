"""
Extract lesson content from Applied Calculus (CC BY 3.0) PDF.
"""
import fitz
import json
import re
import urllib.request
import tempfile
from pathlib import Path

PDF_URL = "https://www.opentextbookstore.com/appcalc/appcalc.pdf"
OUTPUT_DIR = Path("data/lessons/teaching")

# Slug -> (start_page, end_page, title) based on PDF table of contents
SECTION_PAGES = {
    "limits": {
        "pages": (74, 79),
        "title": "Limits and Continuity",
    },
    "derivative-rules": {
        "pages": (97, 121),
        "title": "Derivative Rules",
    },
    "implicit-related-rates": {
        "pages": (162, 162),
        "title": "Implicit Differentiation and Related Rates",
    },
    "integrals": {
        "pages": (190, 214),
        "title": "Integrals",
    },
}


def extract_pdf_text(pdf_path):
    doc = fitz.open(pdf_path)
    pages = []
    for i, page in enumerate(doc):
        text = page.get_text()
        pages.append({"num": i + 1, "text": text})
    doc.close()
    return pages


def extract_section(pages, start, end, title):
    md = [f"# {title}\n"]
    for p in pages:
        if start <= p["num"] <= end:
            text = p["text"].strip()
            if text:
                # Clean up PDF artifacts
                lines = text.split("\n")
                cleaned = []
                for line in lines:
                    line = line.strip()
                    if not line:
                        cleaned.append("")
                    elif re.match(r'^\d+\.\d+\.?\s', line):
                        # Section heading like "2.1" 
                        cleaned.append(f"\n## {line}")
                    elif line.isdigit() and len(line) < 4:
                        # Skip page numbers
                        continue
                    elif line.startswith("http"):
                        continue
                    else:
                        cleaned.append(line)
                md.append("\n".join(cleaned))
    return "\n\n".join(md)


def main():
    import re

    print("Downloading Applied Calculus PDF...")
    with tempfile.NamedTemporaryFile(suffix=".pdf", delete=False) as tmp:
        urllib.request.urlretrieve(PDF_URL, tmp.name)
        pdf_path = tmp.name

    print("Extracting text...")
    pages = extract_pdf_text(pdf_path)

    with open("data/lessons/mapping.json") as f:
        mapping = json.load(f)

    calc_concepts = {k: v for k, v in mapping.items() if v.get("source") == "applied_calc"}
    slugs = {}
    for cid, info in calc_concepts.items():
        slugs.setdefault(info["slug"], []).append(cid)

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    for slug, concept_ids in slugs.items():
        if slug not in SECTION_PAGES:
            print(f"  SKIP {slug}: no page mapping")
            continue
        info = SECTION_PAGES[slug]
        print(f"\nExtracting {slug} -> pages {info['pages']}")
        md = extract_section(pages, info["pages"][0], info["pages"][1], info["title"])
        attribution = (
            "> Content sourced from [Applied Calculus](https://www.opentextbookstore.com/appcalc/) "
            "by Calaway, Hoffman & Lippman — CC BY 3.0\n\n"
        )
        md = attribution + md.strip()

        for cid in concept_ids:
            out_path = OUTPUT_DIR / f"{cid}.md"
            out_path.write_text(md)
            print(f"  -> Wrote {out_path}")

    Path(pdf_path).unlink(missing_ok=True)
    print("\nDone!")


if __name__ == "__main__":
    main()
