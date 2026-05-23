"""
Extract lesson content from Hefferon Linear Algebra 4th ed (CC BY-SA 2.5) PDF.
"""
import fitz
import json
import re
import urllib.request
import tempfile
from pathlib import Path

PDF_URL = "https://raw.githubusercontent.com/tpn/pdfs/master/Linear%20Algebra%20-%20Jim%20Hefferon%20-%204th%20Ed.pdf"
OUTPUT_DIR = Path("data/lessons/teaching")

# Concept -> page range (from TOC analysis)
# Chapter One: Linear Systems (pages ~1-70)
# We need Section I.1: Gauss's Method and I.2: Describing the Solution Set
CONCEPT_PAGES = {
    "linalg.systems.matrix": {
        "pages": (12, 32),
        "title": "Solving Linear Systems",
    },
}


def main():
    print("Downloading Hefferon Linear Algebra PDF...")
    with tempfile.NamedTemporaryFile(suffix=".pdf", delete=False) as tmp:
        urllib.request.urlretrieve(PDF_URL, tmp.name)
        pdf_path = tmp.name

    print("Extracting text from PDF...")
    doc = fitz.open(pdf_path)
    pages = [{"num": i + 1, "text": page.get_text()} for i, page in enumerate(doc)]
    doc.close()
    print(f"  {len(pages)} pages extracted")

    with open("data/lessons/mapping.json") as f:
        mapping = json.load(f)

    la_concepts = {k: v for k, v in mapping.items() if v.get("source") == "hefferon_linalg"}
    print(f"Hefferon LA concepts: {len(la_concepts)}")

    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    for cid, info in sorted(la_concepts.items()):
        if cid not in CONCEPT_PAGES:
            print(f"  SKIP {cid}: no page mapping")
            continue
        p_range = CONCEPT_PAGES[cid]["pages"]
        title = CONCEPT_PAGES[cid]["title"]
        print(f"\nExtracting {cid} (pages {p_range[0]}-{p_range[1]})...")

        md_parts = [f"# {title}\n"]
        for p in pages:
            if p_range[0] <= p["num"] <= p_range[1]:
                text = p["text"].strip()
                if text:
                    lines = text.split("\n")
                    cleaned = []
                    for line in lines:
                        line = line.strip()
                        if not line:
                            cleaned.append("")
                        elif re.match(r'^\d+\.\d+\s', line):
                            cleaned.append(f"\n## {line}")
                        elif line.isdigit() and len(line) < 4:
                            continue
                        else:
                            cleaned.append(line)
                    md_parts.append("\n".join(cleaned))

        md = "\n\n".join(md_parts)
        attribution = (
            "> Content sourced from [Linear Algebra](https://hefferon.net/linearalgebra/) "
            "by Jim Hefferon — CC BY-SA 2.5\n\n"
        )
        md = attribution + md.strip()

        out_path = OUTPUT_DIR / f"{cid}.md"
        out_path.write_text(md)
        print(f"  -> Wrote {out_path} ({len(md)} chars)")

    Path(pdf_path).unlink(missing_ok=True)
    print("\nDone!")


if __name__ == "__main__":
    main()
