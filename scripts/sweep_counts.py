"""Sweep every number the docs, the tests and the app copy pin, from the corpus that owns them.

A content batch moves the concept count, and eleven files quote it: three docs quote the reach and
the unreachable remainder, a mastery ADR quotes the unit count, AGENTS.md and CONTRIBUTING.md quote
the shards and the knowledge points, three published docs quote the concept count, two Go tests
assert it, and four strings in the app copy advertise it. Regenerating the derived artifacts is a
separate chain.

Doing that by hand cost a turn per batch and produced two classes of mistake, both of which this
script is written to prevent:

  * **Replacing a literal that had already moved.** `text.replace("666 concepts", "679 concepts")`
    matches nothing on a file that now says 672, and `str.replace` reports nothing. Every rule here
    is a *pattern* with a required match count, so a rule that finds nothing is an error rather than
    a silent no-op.
  * **Publishing someone else's work.** A second writer has uncommitted edits in four of these
    files, so editing them in the working tree and `git add`-ing them would commit that work too.
    `--blob-only` builds a blob from HEAD plus the substitution and stages it with `update-index`,
    which leaves the working tree alone.

One trap is a committing habit rather than a rule, and this tool can only name it: `npm run
graph:build` writes **two** files under `web/next-app/public/data/` -- `concepts.graph.json` (the node
and edge list) and `concepts.positions.json` (the quantized poster layout). A batch commit that staged
the first and not the second left the committed position map at the corpus of 666 concepts while the
node list moved on, so `scripts/build-og.mjs` and the landing-page poster laid out a graph that had no
position for anything added since. **Stage both, or neither.** The same applies to
`web/next-app/data/flow.positions.json`, which `npm run flow:build` writes and which the flow view
reads.

A second trap is a build that reads the wrong copy. `web/next-app/content/docs.generated.json`
embeds the full text of the docs the website publishes, and two of them are blob-only. Building it
from the working tree embeds the other writer's copy with the previous batch's counts, so the
batch that took the corpus to 811 shipped a stale bundle and `docs:check` went red in CI. The
sweep now rebuilds the bundle after it stages the blobs, handing the staged text to
`build-docs.mjs --overrides` so the bundle matches the commit rather than the tree, and then
verifies with `docs:check`.

A class this tool does not repair, and can only report: **line citations**. Editing a file that an
ADR cites by `path:line` moves every symbol below the edit, and `scripts/check_docs.py` fails on a
citation that lands on a blank line. It caught exactly that when this file's own author added a
registration line to `internal/generator/geometry/generators.go`; the fix is to re-point the citation
at the symbol's new line, which needs a human to know which symbol was meant.

**It recurs by design, so expect it once per batch that adds a registration.** Any line inserted into
a Go file moves every symbol below it, and registrations are inserted into files that ADRs cite -- so
`check_docs` will keep reporting one citation per batch until the citations are anchored to something
that does not move. Fixing it is one line: find the symbol's new line with `grep -n` and update the
citation. It has happened twice, on ADR-041's citation of `internal/generator/geometry/generators.go`.

Usage:
    python3 scripts/sweep_counts.py             # regenerate, rewrite, verify
    python3 scripts/sweep_counts.py --no-regen  # rewrite only, from the corpus as it stands
    python3 scripts/sweep_counts.py --check     # report what would change; write nothing

Exit code 0 means every rule matched and every gate passed.
"""
from __future__ import annotations

import argparse
import collections
import glob
import json
import os
import re
import subprocess
import sys
import tempfile

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
os.chdir(ROOT)

# Files another writer has uncommitted edits in. They are staged as blobs built from HEAD so their
# working trees stay untouched; drop a path from this list once the tree is clean again.
# The list is empty: the four files once named here are clean, and a stale entry makes the docs
# bundle read HEAD while the working tree holds the newer text.
BLOB_ONLY = ()

# Files that counts.py already renders in full (badge, prose, domain table, command comments).
# They are rewritten by calling that renderer rather than by patterns: a pattern that patched only
# the badge left README's body a batch behind, and the corpus-count check caught it.
RENDERED = ("README.md",)


# The regeneration chain. Order matters: the mastery units read the DAG, the web copy reads the
# concept files, and the counts read everything.
REGEN = (
    ["make", "mastery-units"],
    ["go", "run", "./scripts/merge_concepts"],
    ["npm", "run", "graph:build"],
    ["npm", "run", "flow:build"],
    ["python3", "scripts/counts.py", "--write"],
)
WEB = os.path.join("web", "next-app")


def run(cmd, cwd=None):
    print(f"  $ {' '.join(cmd)}")
    r = subprocess.run(cmd, cwd=cwd, capture_output=True, text=True)
    if r.returncode != 0:
        sys.stderr.write(r.stdout + r.stderr)
        raise SystemExit(f"command failed: {' '.join(cmd)}")
    tail = [l for l in r.stdout.splitlines() if l.strip()][-1:] or [""]
    print(f"    {tail[0][:110]}")


def facts() -> dict:
    """Read every number from its owner: the concept files, the shards, the courses, the partition."""
    concepts = {}
    for path in sorted(glob.glob("data/concepts/*.json")):
        if path.endswith("enrichment.json"):
            continue
        for c in json.load(open(path, encoding="utf-8")):
            if "id" in c:
                concepts[c["id"]] = c
    prereq = {c: set(v.get("prerequisites") or []) for c, v in concepts.items()}
    domain = {c: v["domain"] for c, v in concepts.items()}

    def closure(ids):
        seen, stack = set(), list(ids)
        while stack:
            c = stack.pop()
            if c in seen or c not in prereq:
                continue
            seen.add(c)
            stack += list(prereq[c])
        return seen

    mc = collections.Counter(v["domain"] for c, v in concepts.items()
                            if v.get("grading_type") == "multiple_choice")
    courses = json.load(open("data/courses.json", encoding="utf-8"))
    reach = set()
    for co in courses:
        reach |= closure(co.get("targets", []))
    unreachable = collections.Counter(domain[c] for c in set(concepts) - reach)

    dep = collections.defaultdict(set)
    for c, ps in prereq.items():
        for p in ps:
            dep[p].add(c)
    indeg = {c: len(prereq[c]) for c in concepts}
    depth = {c: 0 for c in concepts}
    stack = [c for c in concepts if indeg[c] == 0]
    while stack:
        c = stack.pop()
        for d in dep[c]:
            depth[d] = max(depth[d], depth[c] + 1)
            indeg[d] -= 1
            if indeg[d] == 0:
                stack.append(d)

    kp_files = glob.glob(os.path.join("data/lessons", "kp", "*.json"))
    kps = sum(len(json.load(open(f, encoding="utf-8"))) for f in kp_files)
    units = json.load(open("data/mastery_units.json", encoding="utf-8"))
    unit_list = units.get("units") if isinstance(units, dict) else units
    return {
        "concepts": len(concepts),
        "domains": len(set(domain.values())),
        "courses": len(courses),
        "kp_files": len(kp_files),
        "kps": kps,
        "reach": len(reach),
        "unreachable": len(concepts) - len(reach),
        "max_depth": max(depth.values()),
        "units": len(unit_list),
        "by_domain": {d: (len([c for c in concepts if domain[c] == d]),
                      unreachable.get(d, 0), mc.get(d, 0)) for d in set(domain.values())},
        "discrete_unreachable": unreachable["discrete_math"],
    }


def update_target_table(text: str, f: dict) -> str:
    """Rewrite section 5.5's Now and Add columns from the corpus, keeping Target.

    Target is the document's decision, so this function never invents one. Now is read from
    `f["by_domain"]`, Add is Target - Now, and the sentence above the table is the sum of Add.
    A row whose target has fallen below the corpus is an error: the table has no negative Add
    shape, and silently clamping it would hide a target that is already overshot.
    """
    lines = text.split("\n")
    header = "| Domain | Tier | Now | Target | Add |"
    start = next((i for i, line in enumerate(lines) if line.strip() == header), -1)
    if start < 0:
        raise SystemExit("docs/content-expansion.md: the section 5.5 target table is gone")

    row_re = re.compile(
        r"^\|\s*(\w+)\s*\|\s*(\w+)\s*\|\s*(\d+)\s*\|\s*(\d+)\s*\|\s*\+(\d+)\s*\|\s*$"
    )
    seen: set[str] = set()
    remaining = 0
    i = start + 2
    while i < len(lines):
        m = row_re.match(lines[i])
        if not m:
            break
        domain, tier, _now, target = m.group(1), m.group(2), int(m.group(3)), int(m.group(4))
        if domain in seen:
            raise SystemExit(f"docs/content-expansion.md: duplicate target row for {domain}")
        if domain not in f["by_domain"]:
            raise SystemExit(
                f"docs/content-expansion.md: target row {domain} is not a corpus domain"
            )
        now = f["by_domain"][domain][0]
        add = target - now
        if add < 0:
            raise SystemExit(
                f"docs/content-expansion.md: {domain} target {target} is below the corpus {now}"
            )
        lines[i] = f"| {domain} | {tier} | {now} | {target} | +{add} |"
        seen.add(domain)
        remaining += add
        i += 1

    missing = sorted(set(f["by_domain"]) - seen)
    if missing:
        raise SystemExit(
            "docs/content-expansion.md: no target row for " + ", ".join(missing)
        )
    if not seen:
        raise SystemExit("docs/content-expansion.md: the target table has no rows")

    out = "\n".join(lines)
    claim = re.compile(r"The remaining per-domain targets add up to (\d[\d,]*)")
    out, n = claim.subn(f"The remaining per-domain targets add up to {remaining},", out, count=1)
    if n != 1:
        raise SystemExit(
            "docs/content-expansion.md: the sentence above the target table no longer states the remainder"
        )
    return out


def rules(f: dict) -> list:
    """(path, pattern, replacement, minimum matches[, region]), a pattern that matches nothing is an
    error rather than a silent no-op, and a region restricts the rule to one part of a file.

    The region matters for the two table rows below: `| discrete_math | 31 | N |` appears in both
    section 4.1 (the unreachable count) and section 4.2 (the multiple-choice count), so an
    unanchored pattern rewrites the wrong table. That bug was in this script's first version.
    """
    return [
        # The published engineering docs that quote the corpus size.
        ("docs/architecture.md", r"\b\d{3,4} concepts\b", f"{f['concepts']} concepts", 1),
        ("docs/system-design.md", r"\b\d{3,4} concepts\b", f"{f['concepts']} concepts", 1),
        ("docs/system-design.md", r"\b\d{3,4} KP shards\b", f"{f['concepts']} KP shards", 1),
        ("docs/system-design-slides.md", r"\b\d{3,4} concepts\b", f"{f['concepts']} concepts", 1),
        # AGENTS.md and CONTRIBUTING.md quote the shards and the knowledge points.
        ("AGENTS.md", r"all \d{3,4} shards", f"all {f['kp_files']} shards", 1),
        ("CONTRIBUTING.md", r"All \d{3,4} shards", f"All {f['kp_files']} shards", 1),
        ("CONTRIBUTING.md", r"all \d{3,4} shard files × 3 = \d{4} KPs",
         f"all {f['kp_files']} shard files × 3 = {f['kps']} KPs", 1),
        # The mastery docs and the ADR that describe the partition.
        ("docs/mastery-assessments.md", r"\d{2,3} units, all \d{3,4} concepts",
         f"{f['units']} units, all {f['concepts']} concepts", 1),
        ("docs/mastery-semantics.md", r"\(\d{3,4} concepts, all with a threshold\)",
         f"({f['concepts']} concepts, all with a threshold)", 1),
        ("docs/adr/ADR-050-a-mastery-unit-is-a-named-partition.md",
         r"the authoring of \d{3,4} concepts", f"the authoring of {f['concepts']} concepts", 1),
        ("docs/adr/ADR-050-a-mastery-unit-is-a-named-partition.md",
         r"\*\*\d{2,3} units, all \d{3,4} concepts covered exactly once\*\*",
         f"**{f['units']} units, all {f['concepts']} concepts covered exactly once**", 1),
        ("docs/adr/ADR-050-a-mastery-unit-is-a-named-partition.md",
         r"\(\d{3,4}/\d{3,4}, max \d+\)", f"({f['concepts']}/{f['concepts']}, max {f['max_depth']})", 1),
        # The course count appears in prose four times and in a table twice, across three files,
        # so it is a fact like the others -- and a stale one is invisible to check_docs, which
        # matches numbers against corpus-count nouns rather than against a word.
        ("docs/content-expansion.md", r"Across the \d+ courses", f"Across the {f['courses']} courses", 1),
        ("docs/content-expansion.md", r"The union of the \d+ closures", f"The union of the {f['courses']} closures", 1),
        ("docs/system-design.md", r"\(\d+ courses\)", f"({f['courses']} courses)", 1),
        ("docs/wasc-uc-ag-research.md", r"\d+ courses from `data/courses.json:1`",
         f"{f['courses']} courses from `data/courses.json:1`", 1),
        ("docs/wasc-uc-ag-research.md", r"\d+ courses from 4th grade",
         f"{f['courses']} courses from 4th grade", 1),
        ("docs/wasc-uc-ag-research.md", r"`data/courses.json:1`, \d+ courses",
         f"`data/courses.json:1`, {f['courses']} courses", 1),
        # The section 5.5 target table's own introduction names the current size, so it is a fact
        # too -- and it was written by hand, which is how it went stale within the same batch that
        # added it. The tool's verification caught it.
        ("docs/content-expansion.md", r"The catalog holds \d+ concepts",
         f"The catalog holds {f['concepts']} concepts", 1),
        # content-expansion.md is the owner of the reach measurement.
        ("docs/content-expansion.md", r"There are \d{3,4} concepts, and there are \d{3,4} KP shard files",
         f"There are {f['concepts']} concepts, and there are {f['kp_files']} KP shard files", 1),
        ("docs/content-expansion.md", r"holding \d{4} knowledge points",
         f"holding {f['kps']} knowledge points", 1),
        ("docs/content-expansion.md", r"### 4.1 The course catalog reaches \d+ of the \d{3,4} concepts",
         f"### 4.1 The course catalog reaches {f['reach']} of the {f['concepts']} concepts", 1),
        ("docs/content-expansion.md", r"those closures reach\n\d+ of the \d{3,4} concepts: \d+ of them",
         f"those closures reach\n{f['reach']} of the {f['concepts']} concepts: "
         f"{f['unreachable']} of them", 1),
        # Section 4.1's table has a row per domain, and every row's two numbers move when the
        # corpus does -- so every row needs a rule, not just the domain a batch happened to touch.
        # The rows keep the document's own order; only the numbers are rewritten.
        *[("docs/content-expansion.md", rf"\| {d} \| \d+ \| \d+ \|",
           f"| {d} | {f['by_domain'][d][0]} | {f['by_domain'][d][1]} |", 1, ("### 4.1", "### 4.2"))
          for d in sorted(f["by_domain"])],
        # Section 4.2's table is the same shape with the multiple-choice count, and it lists only
        # the domains that are mostly recognition.
        *[("docs/content-expansion.md", rf"\| {d} \| \d+ \| \d+ \|",
           f"| {d} | {f['by_domain'][d][0]} | {f['by_domain'][d][2]} |", 1, ("### 4.2", "## 5."))
          for d in ("abstract_algebra", "topology", "real_analysis", "trigonometry",
                    "probability_theory")],
        # `internal/concepts/concepts_test.go` used to pin exact counts. It now asserts a floor
        # (`Count() < 1000`) and self-consistency (`len(Order()) == Count()`), so there is no
        # literal left to rewrite. A rule that matches nothing is an error, and the test's own
        # comment records why an exact literal must not come back.
        # `internal/planning/planner_test.go` used to pin the course count with a literal. It now
        # reads `data/courses.json` and compares against the number it finds there, so there is no
        # site left to rewrite in this batch; a rule that matches nothing is an error, and the
        # test's own comment records why the literal must not come back.
        # The app copy a learner reads, and the test that guards it.
        # `web/next-app/app/how-it-works/sections.tsx` used to state the count as a literal. It now
        # renders `{CONCEPT_COUNT}`, which is read from the corpus at build time, so there is no
        # literal to rewrite; a rule that matches nothing is an error, so none is kept.
        ("web/next-app/app/layout.tsx", r"Mastery-gated, \d{3,4} concepts",
         f"Mastery-gated, {f['concepts']} concepts", 1),
        ("web/next-app/app/layout.tsx", r"Open-source, \d{3,4} concepts, spaced repetition, mastery checks",
         f"Open-source, {f['concepts']} concepts, spaced repetition, mastery checks", 1),
        ("web/next-app/app/layout.tsx", r"Open-source, \d{3,4} concepts, spaced repetition\.",
         f"Open-source, {f['concepts']} concepts, spaced repetition.", 1),
        # `web/next-app/test/copyTruth.test.ts` intentionally holds a superseded count as a
        # fixture ("a superseded concept count; the graph holds 1250"). Rewriting it would
        # destroy the test's subject, so the rule is gone; the test owns no current claim.
    ]


def counts_render(text: str) -> str:
    """Apply scripts/counts.py's own renderer, so the sweep and the generator cannot disagree."""
    sys.path.insert(0, os.path.join(ROOT, "scripts"))
    import counts  # noqa: E402  (imported here: it reads the corpus at import time)
    return counts.render(text, counts.compute())


def content_of(path: str) -> str:
    """HEAD's content for a blob-only path, the working tree's otherwise."""
    if path in BLOB_ONLY:
        return subprocess.run(["git", "show", f"HEAD:{path}"], capture_output=True, text=True,
                              check=True).stdout
    return open(path, encoding="utf-8").read()


def build_docs(overrides: dict[str, str], check: bool = False):
    """Build (or check) the published-docs bundle from the sources as they will be committed.

    `content/docs.generated.json` embeds the text of every doc the website publishes, and two of
    them are blob-only: this sweep stages their new text without touching the working tree, so a
    build that read the tree would embed the other writer's copy with the previous batch's
    counts. Hand the staged text to `--overrides` so the bundle matches the commit instead. A key
    the DOCS list does not name is ignored, so passing every staged path is safe.

    With `check`, return the completed check process instead of raising, so the caller can report
    it alongside the other gates.
    """
    tmpdir = tempfile.mkdtemp(prefix="mathua-docs-")
    path = os.path.join(tmpdir, "overrides.json")
    try:
        with open(path, "w", encoding="utf-8") as fh:
            json.dump(overrides, fh)
        cmd = ["node", "scripts/build-docs.mjs", "--overrides", path]
        if check:
            cmd.append("--check")
            return subprocess.run(cmd, cwd=WEB, capture_output=True, text=True)
        run(cmd, cwd=WEB)
        return None
    finally:
        os.remove(path)
        os.rmdir(tmpdir)


def main() -> int:
    ap = argparse.ArgumentParser()
    ap.add_argument("--no-regen", action="store_true", help="skip the regeneration chain")
    ap.add_argument("--check", action="store_true", help="report changes; write nothing")
    args = ap.parse_args()

    if not args.no_regen and not args.check:
        print("regenerating the derived artifacts:")
        for cmd in REGEN:
            cwd = WEB if cmd[0] == "npm" else None
            run(cmd, cwd=cwd)
    elif args.check:
        print("check mode: no regeneration, no writes")

    f = facts()
    print("\ncorpus facts: " + json.dumps(f, sort_keys=True))

    grouped = collections.defaultdict(list)
    for rule in rules(f):
        path, pattern, repl, at_least = rule[:4]
        region = rule[4] if len(rule) > 4 else None
        grouped[path].append((pattern, repl, at_least, region))
    for path in RENDERED:
        grouped.setdefault(path, [])  # rendered files may have no pattern rules of their own

    def replace(text, pattern, repl, at_least, region):
        """Apply one rule, honouring its region, and return the new text and the match count."""
        if region is None:
            return re.subn(pattern, repl, text)
        start = text.find(region[0])
        end = text.find(region[1], start + 1) if start >= 0 else -1
        if start < 0 or end < 0:
            raise SystemExit(f"region markers {region} not found")
        head, body, tail = text[:start], text[start:end], text[end:]
        body, n = re.subn(pattern, repl, body)
        return head + body + tail, n

    changed = 0
    overrides: dict[str, str] = {}
    for path, subs in grouped.items():
        text = content_of(path)
        original = text
        if path in RENDERED:
            subs = []
        for pattern, repl, at_least, region in subs:
            text, n = replace(text, pattern, repl, at_least, region)
            if n < at_least:
                raise SystemExit(f"{path}: /{pattern}/ matched {n} times, expected at least {at_least}")
        if path == "docs/content-expansion.md":
            text = update_target_table(text, f)
        if path in RENDERED:
            text = counts_render(text)
        here = sum(1 for a, b in zip(original.split("\n"), text.split("\n")) if a != b)
        if path in BLOB_ONLY:
            # Always record the staged text, changed or not: the working-tree copy is the other
            # writer's, so the docs bundle must read HEAD plus the substitution either way.
            overrides[path] = text
        if text == original:
            print(f"  current  {path}")
            continue
        changed += 1
        if args.check:
            print(f"  would change {path} ({here} line(s))")
            continue
        if path in BLOB_ONLY:
            blob = subprocess.run(["git", "hash-object", "-w", "--stdin"], input=text.encode(),
                                  capture_output=True, check=True).stdout.decode().strip()
            subprocess.run(["git", "update-index", "--cacheinfo", f"100644,{blob},{path}"], check=True)
            print(f"  staged    {path} ({here} line(s), blob from HEAD, working tree untouched)")
        else:
            open(path, "w", encoding="utf-8").write(text)
            print(f"  rewrote   {path} ({here} line(s))")

    if args.check:
        print(f"\n{changed} file(s) would change")
        return 0
    print(f"\n{changed} file(s) rewritten")

    print("\nrebuilding the published-docs bundle:")
    build_docs(overrides)

    print("\nverifying:")
    failures = []
    for cmd in (["python3", "scripts/check_docs.py"], ["python3", "scripts/counts.py", "--check"],
                ["python3", "scripts/check_additive.py"],
                ["python3", "scripts/audit_lessons.py"], ["go", "run", "scripts/validate_graph.go"],
                ["go", "run", "./cmd/validate-mastery-units"]):
        r = subprocess.run(cmd, capture_output=True, text=True)
        if r.returncode == 0:
            run(cmd)
            continue
        # Only check_docs.py earns the lenient path. A blob-only file's working-tree copy is
        # deliberately stale -- it still holds the other writer's text plus the previous batch's
        # numbers -- so check_docs reports it until the commit lands, and those lines are
        # expected. Every other command's nonzero exit is a real failure, even when its output
        # has no ERROR line to collect.
        lenient = any("check_docs.py" in part for part in cmd)
        expected, real = [], []
        for line in (r.stdout + r.stderr).splitlines():
            if not line.startswith("ERROR"):
                # The "FAIL: N inconsistencies" summary counts the expected ones too, so it is
                # only meaningful when an ERROR line is real.
                continue
            (expected if any(path in line for path in BLOB_ONLY) else real).append(line)
        if not lenient and not real:
            real = (r.stdout + r.stderr).strip().splitlines()[-1:] or ["(no output)"]
        if real:
            failures.append((cmd, real))
            print(f"  $ {' '.join(cmd)}  FAILED")
        else:
            print(f"  $ {' '.join(cmd)}  ok ({len(expected)} staged-only line(s) not yet in the "
                  f"working tree)")
    # The bundle is checked against the overrides rather than the working tree, because the
    # published docs the sweep stages are blob-only: a plain `docs:check` would compare the
    # bundle with the other writer's copy and fail for the wrong reason.
    r = build_docs(overrides, check=True)
    if r.returncode:
        failures.append((["docs:check"], (r.stdout + r.stderr).strip().splitlines()[-2:]))
        print("  $ docs:check  FAILED")
    else:
        print("  $ docs:check  ok")
    if failures:
        for cmd, lines in failures:
            sys.stderr.write(f"\n{' '.join(cmd)}:\n" + "\n".join(lines[:10]) + "\n")
        return 1
    print("sweep ok")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
