#!/usr/bin/env python3
"""Fetch product-wide efficacy datasets from a running server and reshape
them into the efficacy-bundle/1 schema that figures/efficacy.py renders.

Both source endpoints are public and rate-limited (docs/efficacy.md):
  GET /api/efficacy/all    -> EfficacyReport (point-in-time aggregate)
  GET /api/efficacy/trend  -> EfficacyTrend (weekly buckets + retention)

Usage:
    python3 figures/fetch_prod.py [--api https://mathua.fly.dev] [--out figures/efficacy_data.json]
"""
from __future__ import annotations

import argparse
import json
import sys
import urllib.request
from pathlib import Path


def get_json(base: str, path: str) -> dict:
    req = urllib.request.Request(base.rstrip("/") + path, headers={"Accept": "application/json"})
    with urllib.request.urlopen(req, timeout=30) as res:
        if res.status != 200:
            raise RuntimeError(f"{path}: HTTP {res.status}")
        return json.load(res)


def main(argv: list[str] | None = None) -> int:
    ap = argparse.ArgumentParser(description="Fetch efficacy datasets into a render bundle")
    ap.add_argument("--api", default="https://mathua.fly.dev", help="API origin (no trailing path)")
    ap.add_argument("--out", default="figures/efficacy_data.json", help="bundle output path")
    args = ap.parse_args(argv)

    try:
        report = get_json(args.api, "/api/efficacy/all")
        trend = get_json(args.api, "/api/efficacy/trend")
    except Exception as exc:
        print(f"fetch failed: {exc}", file=sys.stderr)
        return 1

    weeks = trend.get("weeks") or []
    starts = [w["week_start"] for w in weeks if w.get("week_start")]
    bundle = {
        "schema": "efficacy-bundle/1",
        "source": f"{args.api.rstrip('/')}/api/efficacy/all (+trend)",
        "total_attempts": report.get("total_attempts", 0),
        "students_tracked": report.get("students_tracked", 0),
        "concepts_touched": report.get("concepts_touched", 0),
        "concepts_catalog": [],
        "first_pass_rate": report.get("first_pass_rate", 0.0),
        "second_pass_rate": report.get("second_pass_rate", 0.0),
        "avg_attempts_per_concept": report.get("avg_attempts_per_concept", 0.0),
        "window_start": min(starts) if starts else None,
        "window_end": max(starts) if starts else None,
        "observational": True,  # usage metrics, not a trial (docs/efficacy.md)
    }
    out = Path(args.out)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(bundle, indent=2) + "\n")
    print(f"bundle: {out} ({bundle['total_attempts']} attempts, "
          f"{bundle['students_tracked']} students, weeks={len(weeks)})")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
