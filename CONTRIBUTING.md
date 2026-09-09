# Contributing to Mathua

Mathua is community-built. Every concept, every generator, every line of the concept graph was added by someone who wanted to help others learn math better. The most useful thing you can contribute is a new concept — and it takes exactly three pieces.

## Quick start

Clone the repo and make sure the web server works:

```bash
git clone https://github.com/chuma-beep/mathua.git
cd mathua

go build ./cmd/mathua

# Run the web server (requires Postgres)
DATABASE_URL=postgres://... ./mathua --serve --port 8080
```

The concept graph lives in per-domain files under `data/concepts/`. Generators live in `internal/generator/`. Here are the kinds of contributions that move the needle:

- New concepts — the single most impactful thing you can add. Concepts using `polynomial` or `expression` grading require Python 3.8+ and `sympy` for full equivalence checking (falls back to string comparison if unavailable).
- New generators — make existing or new concepts produce better problems.
- Bug fixes in the scheduling engine or graders.
- Documentation and diagram improvements.

## Adding a concept — the 3-step summary

### 1. Define a new concept

Concepts live in per-domain files under `data/concepts/`. Each file is a JSON array — add your new concept as a new entry in the array for the matching domain:

```json
[
  { "id": "arith.add.single", "label": "Single-digit addition", ... },
  { "id": "arith.add.multi",  "label": "Multi-digit addition",  ... },
  {
    "id":                "arith.mult.tables",
    "label":             "Multiplication tables 1-12",
    "domain":            "arithmetic",
    "subdomain":         "arithmetic.multiplication",
    "grading_type":      "numeric",
    "prerequisites":     ["arith.add.multi"],
    "mastery_threshold": {
      "streak":           7,
      "avg_time_seconds": 6.0
    }
  }
]
```

Valid `grading_type` values: `numeric`, `polynomial`, `expression`, `multiple_choice`, `comparison`, `ordering`, `tuple`, `complex`, `symbolic` (9 total, `internal/concepts/loader.go:297`).

> The prerequisites list is the most important field. What must a student absolutely know before attempting this? If in doubt, add the prerequisite — the graph validator will catch cycles.

### 2. Teach Mathua to ask questions

Write a Go generator in `internal/generator/[domain]/`:

```go
package arithmetic

import "mathua/internal/generator"

type AddSingleGen struct{}

func (g *AddSingleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
    // ctx.Seed is hashSeed(studentID|cid|attemptID) set by engine.go:389 via
    // Registry.GenerateContext(ctx{Seed,Difficulty}) generator/registry.go:74
    max := int(5 + ctx.Difficulty*4)
    a   := rand.Intn(max) + 2
    b   := rand.Intn(max) + 2
    return generator.Problem{
        Question:    fmt.Sprintf("%d + %d = ?", a, b),
        Answer:      strconv.Itoa(a + b),
        Explanation: fmt.Sprintf("%d + %d = %d", a, b, a+b),
    }
}
```

- Use `ctx.Difficulty` (0.0 trivial → 1.0 challenging) to scale operand sizes.
- Always return an `Explanation` — shown when a student asks to see the solution.
- Use the global `math/rand` (seeded per-request via `registry.go:84`); no hardcoded problems. `ctx.Seed` is deterministic `hashSeed(student|cid|attemptID)` (`engine.go:389`), no `UnixNano`.

### 3. Prove it works

Write a fuzz test that asserts 1 000 valid samples:

```go
func TestAddSingleGen(t *testing.T) {
    g := &AddSingleGen{}
    for i := 0; i < 1000; i++ {
        p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
        assert.NotEmpty(t, p.Question)
        assert.NotEmpty(t, p.Answer)
        assert.NotEmpty(t, p.Explanation)
        assert.True(t, evaluateNumeric(p))
    }
}
```

Run your tests before opening a PR (1000 local, 100 lightweight in CI `ci.yml:27` for speed — `make fuzz` is 1000):

```bash
go test ./... -count 1
go run scripts/validate_graph.go
python3 scripts/audit_lessons.py
go test ./internal/generator/... -run TestFuzz -count 1000  # local full; CI runs -count 100
```

### 4. (Optional) Write a lesson

Mathua includes a built-in lesson system. Each concept can have an associated lesson — a markdown file that teaches the material. Lessons live under `data/lessons/` and are registered in `data/lessons/lessons.json`, which maps concept IDs to file paths:

```json
  { "concept_id": "calc.deriv.power_rule", "source": "advanced/polynomials/polynomials.md" }
```

Multiple concepts can share a single lesson file. Lesson content is written in Markdown with LaTeX via `\( ... \)` or `\[ ... \]` delimiters and rendered inside the app as a sidebar panel alongside practice problems.

### 5. (Recommended) Add knowledge-point shards

Every lesson should be broken into 1-3 **knowledge points** (KPs), each a subgoal-labeled worked example — the Math Academy scaffolding pattern. A KP shard for a concept is a JSON file at `data/lessons/kp/<concept_id>.json`:

```json
[
  {
    "label": "Use addition notation",
    "section": "Use Addition Notation",
    "subgoals": [
      "Identify the addends and the sum",
      "Read 3 + 4 as three plus four",
      "Translate word phrases into math notation"
    ]
  }
]
```

- `section` names a heading **inside the lesson body** that serves as the worked example. It must match a heading exactly; leave it empty to fall back to the full lesson body.
- `subgoals` are short labels shown before practice (subgoal-labeling effect). Keep 3 per KP.
- Prefer **3 KPs per concept**; 1-2 are acceptable for very small topics.
- `scripts/gen_kp_shards.py` generates shard files from an inline spec and validates every section reference against the corpus — add your spec entries there and run `python3 scripts/gen_kp_shards.py <domain-prefix>`.

## The graph validator

A validator runs on every pull request. It checks invariants before any merge can happen:

1. **No cycles** — including `prerequisites` + `encompasses` edges (`loader.go:334` Kahn’s)
2. **No orphans** — every `prerequisites`/`encompasses`/`key_prerequisites` must exist (`loader.go:253`)
3. **No dup IDs, no empty label/domain, valid `mastery_threshold`** (`loader.go:303`)
4. **No singleton `interference_group`** (`loader.go:321`), valid `grading_type` enum (9 values `loader.go:297`), `variants` `0.1–1.0` difficulty

A second audit guards the lesson corpus (`python3 scripts/audit_lessons.py:191`). It fails on: DAG concepts with no lesson, stale `lessons.json` ids, lesson files missing on disk, orphaned sources, KP shard sections that do not resolve (630 files ×3 =1890 KPs), diagram mappings that point at missing assets or non-existent concepts (`engine.go:497` 115 diagrams), stale course targets, `grading_type`/`threshold` sanity, and `enrichment.json` wiring.

## Submitting a pull request

1. Add the concept to the appropriate domain file in `data/concepts/` with correct prerequisites.
2. Write the generator in the appropriate domain subdirectory.
3. Write the fuzz test with 1 000 samples.
4. (Optional) Write a lesson and register it in `data/lessons/lessons.json`.
5. (Recommended) Add KP shards under `data/lessons/kp/`.
6. Run `go test ./...`, `go run scripts/validate_graph.go`, and `python3 scripts/audit_lessons.py` locally.
7. Open a PR. The CI pipeline runs the validator and all tests automatically.
8. A maintainer reviews the concept ordering, thresholds, and generator quality.

Reviews usually happen within a few days. If a week passes with no response, ping the thread. We read every PR.

## Full guide

For the detailed walkthrough with field schemas, code examples, and design conventions, see the **[Contributing docs](https://mathua.vercel.app/docs/contributing)** in the web app.
