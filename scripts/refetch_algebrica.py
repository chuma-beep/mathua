#!/usr/bin/env python3
"""Refetch specific Algebrica lesson files whose bodies never landed.

Some committed lessons contain only the SEO fallback ("shown in the
conceptual map") because an earlier scrape captured the meta summary instead
of .post-content. For each given lesson path this reads the upstream URL from
its attribution header and re-runs the current extraction over it.

Usage: python3 scripts/refetch_algebrica.py data/lessons/algebrica/functions/exponential-function.md [...]
"""
import os
import re
import sys
import time

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import scrape_algebrica

STUB_SIGNS = ("shown in the conceptual map",)


def prose_words(text):
    return len(re.findall(r"[A-Za-z]{2,}", text))


def main(paths):
    failures = 0
    for fpath in paths:
        with open(fpath, encoding="utf-8") as f:
            old = f.read()
        m = re.search(r"\(https://algebrica\.org/([^)/]+)/?\)", old)
        if not m:
            print(f"SKIP {fpath}: no algebrica.org URL in attribution header")
            failures += 1
            continue
        url = f"https://algebrica.org/{m.group(1)}/"
        scrape_algebrica.scrape_entry(url, os.path.dirname(os.path.abspath(fpath)), overwrite=True)
        with open(fpath, encoding="utf-8") as f:
            new = f.read()
        n = prose_words(new)
        # Full lessons may quote the phrase in the scraped Graph-Concept footer
        # (16 committed lessons do); the stub test is body size, not the phrase.
        has_sig = any(sig in new for sig in STUB_SIGNS)
        print(f"  {fpath}: {len(new)} bytes, {n} words, footer_phrase={has_sig}")
        if n < 200:
            print(f"  FAIL {fpath}: refetch did not produce a body")
            failures += 1
        time.sleep(0.5)
    return 1 if failures else 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
