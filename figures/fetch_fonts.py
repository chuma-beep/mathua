#!/usr/bin/env python3
"""Fetch display fonts for figures/ (gitignored, never committed).

Downloads static TTFs (Space Grotesk 400, JetBrains Mono 400) from the
Fontsource CDN into figures/fonts/. The figure scripts work without them
(fallback stack), but PNG/PDF exports embed real glyphs only when the
fonts are present. SVG exports reference the family names as text, which
resolve from the site's own webfonts in the browser.

Usage:
    python3 figures/fetch_fonts.py
"""
import sys
import urllib.request
from pathlib import Path

FONTS = {
    "SpaceGrotesk-Regular.ttf": "https://cdn.jsdelivr.net/npm/@fontsource/space-grotesk@latest/files/space-grotesk-latin-400-normal.woff2",
    "JetBrainsMono-Regular.ttf": "https://cdn.jsdelivr.net/npm/@fontsource/jetbrains-mono@latest/files/jetbrains-mono-latin-400-normal.woff2",
}

# NOTE: the CDN serves woff2; matplotlib needs TTF/OTF, so convert with
# fontTools (pip install fonttools brotli). The Space Grotesk file ships
# with family name "Space Grotesk Light"; the style adapter matches on
# that name transparently, so no name-table patching (and no OFL
# Reserved-Font-Name issue) is needed.


def main() -> int:
    try:
        from fontTools.ttLib import TTFont
    except ImportError:
        print("need fonttools + brotli: pip install fonttools brotli", file=sys.stderr)
        return 1

    out_dir = Path(__file__).resolve().parent / "fonts"
    out_dir.mkdir(parents=True, exist_ok=True)
    for filename, url in FONTS.items():
        woff_path = out_dir / (filename + ".woff2")
        ttf_path = out_dir / filename
        if ttf_path.exists():
            print(f"exists: {ttf_path}")
            continue
        print(f"fetch: {url}")
        urllib.request.urlretrieve(url, woff_path)
        font = TTFont(str(woff_path))
        font.flavor = None
        font.save(str(ttf_path))
        woff_path.unlink()
        fam = font["name"].getDebugName(16) or font["name"].getDebugName(1)
        print(f"saved: {ttf_path} (family {fam!r})")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
