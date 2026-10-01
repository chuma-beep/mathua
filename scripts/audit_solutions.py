#!/usr/bin/env python3
"""Coverage report for the corpus solution schemas.

`go test ./internal/generator/ -run LearnerDomains` is the gate: it fails on a
concept in the six learner domains that has neither a schema nor a reasoned
exemption. This script is the human-facing version — it prints the coverage so
`make validate` shows where the content wave has got to, and cross-checks that
the directory and the Go gate agree about what is authored.

It never fails on missing content by itself. The Go test owns that decision;
this owns the count.
"""
import json
import os
import sys

ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..")
CONCEPTS = os.path.join(ROOT, "data", "concepts")
SOLUTIONS = os.path.join(ROOT, "data", "lessons", "solutions")
LEARNER_DOMAINS = {
    "arithmetic", "fractions", "prealgebra",
    "algebra", "geometry", "trigonometry",
}


def authored() -> set[str]:
    return {
        name[:-5]
        for name in os.listdir(SOLUTIONS)
        if name.endswith(".json")
    }


def main() -> int:
    have = authored()
    by_domain: dict[str, list[str]] = {}
    for name in sorted(os.listdir(CONCEPTS)):
        if not name.endswith(".json"):
            continue
        with open(os.path.join(CONCEPTS, name)) as fh:
            raw = json.load(fh)
        items = raw if isinstance(raw, list) else raw.get("concepts", raw.get("items", []))
        for item in items:
            if not isinstance(item, dict):
                continue
            cid, dom = item.get("id"), item.get("domain")
            if cid and dom in LEARNER_DOMAINS:
                by_domain.setdefault(dom, []).append(cid)

    total = sum(len(v) for v in by_domain.values())
    done = sum(1 for v in by_domain.values() for c in v if c in have)
    for dom in sorted(by_domain):
        cs = by_domain[dom]
        n = sum(1 for c in cs if c in have)
        print(f"note: solutions: {dom:14} {n:3}/{len(cs):<3} authored")

    # Every shipped schema must name a concept that exists, or it is dead
    # content that will never be served.
    known = {c for v in by_domain.values() for c in v}
    orphans = sorted(have - known)
    if orphans:
        print(f"error: {len(orphans)} schema(s) name no concept in the six learner domains:",
              file=sys.stderr)
        for c in orphans:
            print(f"  {c}", file=sys.stderr)
        return 1

    print(f"OK: {done} of {total} learner-domain concepts have an authored schema")
    return 0


if __name__ == "__main__":
    sys.exit(main())
