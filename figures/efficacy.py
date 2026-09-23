#!/usr/bin/env python3
"""Efficacy overview pilot figure (figures4papers adoption pilot).

Reads the dev attempt log and renders the point-in-time aggregate panel
(first-pass / second-pass rates with Math Academy parity markers) using
figures/mathua_style.py. Delivery is offline script -> committed SVGs, so
the single-binary deploy is untouched (no Python in the prod path).

Scope decision (2026-09-23): the longitudinal trend panel is DELIBERATELY
omitted. The dev log holds 71 attempts across 2 week-buckets (70 + 1); a
trend line through that would mislead. make_trend ships in the adapter for
later use once /api/efficacy/trend has multi-week volume; --selftest below
keeps it render-covered with synthetic data in the meantime.

Usage:
    python3 figures/efficacy.py                  # dump mathua.db + render
    python3 figures/efficacy.py --input bundle.json --out figures/output/efficacy-overview
    python3 figures/efficacy.py --selftest       # render all helpers to tempdir
    python3 figures/efficacy.py --dump-only      # refresh figures/efficacy_data.json

Regeneration from a live server instead of the dev DB:
    curl -s localhost:8080/api/efficacy/all > /tmp/eff.json  # (auth header needed)
    ... then reshape to the bundle schema below and render with --input.
"""
from __future__ import annotations

import argparse
import json
import sqlite3
import sys
import tempfile
from pathlib import Path

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

import mathua_style as ms  # noqa: E402

# Math Academy parity reference (improve.md:82, docs/efficacy.md:29). External
# observational claims, drawn as markers only -- never merged into our rates.
MA_FIRST_PASS = 0.93
MA_SECOND_PASS = 0.98

BUNDLE_SCHEMA = "efficacy-bundle/1"


def dump_bundle(db_path: Path) -> dict:
    """Mirror internal/engine/engine.go computeEfficacy over the attempt log."""
    con = sqlite3.connect(db_path)
    con.row_factory = sqlite3.Row
    rows = con.execute(
        "SELECT student_id, concept_id, correct, timestamp FROM attempts ORDER BY timestamp"
    ).fetchall()
    con.close()

    by_key: dict[str, list[bool]] = {}
    for r in rows:
        by_key.setdefault(f"{r['student_id']}|{r['concept_id']}", []).append(bool(r["correct"]))

    first = sum(1 for s in by_key.values() if s and s[0])
    second = sum(1 for s in by_key.values() if (s and s[0]) or (len(s) >= 2 and s[1]))
    touched = len(by_key)
    total = sum(len(s) for s in by_key.values())
    stamps = [r["timestamp"] for r in rows]
    students = {r["student_id"] for r in rows}
    concepts = {r["concept_id"] for r in rows}
    return {
        "schema": BUNDLE_SCHEMA,
        "source": str(db_path),
        "total_attempts": total,
        "students_tracked": len(students),
        "concepts_touched": touched,
        "concepts_catalog": sorted(concepts),
        "first_pass_rate": (first / touched) if touched else 0.0,
        "second_pass_rate": (second / touched) if touched else 0.0,
        "avg_attempts_per_concept": (total / touched) if touched else 0.0,
        "window_start": min(stamps) if stamps else None,
        "window_end": max(stamps) if stamps else None,
        "observational": True,  # dev log, not a controlled trial (docs/efficacy.md)
    }


def caption(bundle: dict) -> str:
    ws = (bundle.get("window_start") or "?")[:10]
    we = (bundle.get("window_end") or "?")[:10]
    return (
        f"{bundle['total_attempts']} attempts · {bundle['students_tracked']} students · "
        f"{bundle['concepts_touched']} concepts · {ws} → {we} · dev log, observational"
    )


def render(bundle: dict, out_base: Path) -> list[Path]:
    if bundle.get("schema") != BUNDLE_SCHEMA:
        raise ValueError(f"expected schema {BUNDLE_SCHEMA}, got {bundle.get('schema')!r}")
    ms.register_fonts()
    ms.apply_publication_style(ms.MathuaStyle(font_size=16))

    fig, axes = ms.create_subplots(1, 1, figsize=(9, 5.5))
    ax = axes[0]
    ms.make_grouped_bar(
        ax,
        ["All concepts"],
        [[bundle["first_pass_rate"]], [bundle["second_pass_rate"]]],
        ["First-pass", "Second-pass"],
        ylabel="Share of concepts",
        hatches=[None, "///"],
    )
    for container in ax.containers:  # annotate every series, not just the last
        ms.annotate_bars(ax, container, fmt="{:.0%}")
    # MA parity markers: reference lines, muted so they never read as our data.
    for rate, label in ((MA_FIRST_PASS, "MA 93%"), (MA_SECOND_PASS, "MA 98%")):
        ax.axhline(rate, color=ms.MATHUA_PALETTE["secondary"], linestyle=(0, (4, 3)), linewidth=1.2, alpha=0.8)
        ax.text(0.99, rate + 0.008, f"{label} (external)", ha="right", va="bottom", fontsize=10,
                fontfamily=ms.MONO_STACK, color=ms.MATHUA_PALETTE["secondary"],
                transform=ax.get_yaxis_transform())
    ax.set_ylim(0, 1.02)  # rates + markers span the full range; no artificial zoom
    ax.set_title("Efficacy overview (dev log)", pad=12)
    fig.text(0.5, -0.02, caption(bundle), ha="center", fontsize=10,
             fontfamily=ms.MONO_STACK, color=ms.MATHUA_PALETTE["secondary"])
    return ms.finalize_figure(fig, out_base, formats=["svg", "png", "pdf"], dpi=300)


def selftest() -> None:
    """Render every adapter helper with synthetic data; assert outputs exist."""
    ms.register_fonts()
    ms.apply_publication_style()
    with tempfile.TemporaryDirectory() as tmp:
        base = Path(tmp)
        fig, ax = ms.create_subplots(1, 1, figsize=(6, 4))
        ms.make_trend(ax[0], [1, 2, 3, 4], [[0.5, 0.6, 0.7, 0.8], [0.4, 0.45, 0.6, 0.75]], ["a", "b"])
        assert ms.finalize_figure(fig, base / "trend", formats=["svg", "png"])
        fig, ax = ms.create_subplots(1, 1, figsize=(5, 4))
        ms.make_heatmap(ax[0], [[0.9, 0.6], [0.7, 0.8]], x_labels=["x1", "x2"], y_labels=["y1", "y2"], annotate=True)
        assert ms.finalize_figure(fig, base / "heat", formats=["svg"])
        fig, ax = ms.create_subplots(1, 1, figsize=(5, 4))
        ms.make_scatter(ax[0], [1, 2, 3], [2, 3, 5], label="s")
        assert ms.finalize_figure(fig, base / "scat", formats=["svg"])
    print("selftest: trend + heatmap + scatter render OK")


def main(argv: list[str] | None = None) -> int:
    ap = argparse.ArgumentParser(description="Efficacy overview pilot figure")
    ap.add_argument("--db", default="mathua.db", help="dev sqlite DB (dump source)")
    ap.add_argument("--input", help="JSON bundle (skip dump; e.g. reshaped /api/efficacy/all)")
    ap.add_argument("--out", default="figures/output/efficacy-overview", help="output basename")
    ap.add_argument("--dump-only", action="store_true")
    ap.add_argument("--selftest", action="store_true")
    args = ap.parse_args(argv)

    if args.selftest:
        selftest()
        return 0

    if args.input:
        bundle = json.loads(Path(args.input).read_text())
    else:
        db = Path(args.db)
        if not db.exists():
            print(f"no DB at {db}; pass --input bundle.json", file=sys.stderr)
            return 1
        bundle = dump_bundle(db)
        bundle_path = HERE / "efficacy_data.json"
        bundle_path.write_text(json.dumps(bundle, indent=2) + "\n")
        print(f"bundle: {bundle_path}")
        if args.dump_only:
            return 0

    paths = render(bundle, Path(args.out))
    print("figures:")
    for p in paths:
        print(f"  {p}")
    print(f"rates: first-pass={bundle['first_pass_rate']:.1%} second-pass={bundle['second_pass_rate']:.1%}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
