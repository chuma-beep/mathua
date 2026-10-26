# Contributing to Mathua

Mathua is community-built. Someone added every concept, every generator and every line of the concept graph. They wanted to help others learn math better. The most useful thing you can contribute is a new concept and it takes exactly three pieces.

## Quick start

Clone the repo and make sure the web server works:

```bash
git clone https://github.com/chuma-beep/mathua.git
cd mathua

go build ./cmd/mathua

# Run the web server. SQLite is the default, so there is no database to install.
./mathua --serve --port 8080
```

SQLite keeps its data in `mathua.db` in the working directory. Local development needs no database service and no `DATABASE_URL`.

Postgres runs the same code. Set `DATABASE_URL` to run against a Postgres instance instead:

```bash
DATABASE_URL="postgres://user:pass@localhost:5432/mathua?sslmode=disable" ./mathua --serve --port 8080
```

The deployment decides which database production uses. The Fly deployment keeps SQLite on a volume today, and the storage layer runs on either engine.

### Run the storage tests on both engines

The storage tests run on SQLite by default. An in-memory database needs no service and starts in about a millisecond. Set `TEST_POSTGRES_DSN` to run the same tests against Postgres:

```bash
docker run -d --rm --name mathua-pg -e POSTGRES_PASSWORD=pg \
  -e POSTGRES_DB=mathua -p 5432:5432 postgres:16-alpine
TEST_POSTGRES_DSN="postgres://postgres:pg@localhost:5432/mathua?sslmode=disable" \
  go test ./internal/storage/ -count=1
```

CI runs both engines. A method that behaves differently on Postgres fails the `Storage suite on Postgres` job.

The concept graph lives in per-domain files under `data/concepts/`. Generators live in `internal/generator/`. Here are the kinds of contributions that move the needle:

- New concepts, the single most impactful thing you can add. A pinned SymPy service checks equivalence for `polynomial`, `expression`, `symbolic` and `complex` grading. The service is `grading/sympy_service.py`, and its version is in `grading/requirements.txt`. The production image ships it. There is no pure-Go equivalence fallback. Run `go test ./...` and the grader tests with SymPy installed.
- New generators, make existing or new concepts produce better problems.
- Bug fixes in the scheduling engine or graders.
- Documentation and diagram improvements.

## Adding a concept, the 3-step summary

### 1. Define a new concept

Concepts live in per-domain files under `data/concepts/`. Each file is a JSON array, add your new concept as a new entry in the array for the matching domain:

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

Valid `grading_type` values: `numeric`, `polynomial`, `expression`, `multiple_choice`, `comparison`, `ordering`, `tuple`, `complex`, `symbolic` (9 total, checked by `validate` in `internal/concepts/loader.go`).

> The prerequisites list is the most important field. What must a student absolutely know before attempting this? If in doubt, add the prerequisite, the graph validator will catch cycles.

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
- Always return an `Explanation`, shown when a student asks to see the solution.
- Use the global `math/rand`, seeded per request (`Registry.GenerateContext` in `internal/generator/registry.go`). Do not hardcode problems. `ctx.Seed` is deterministic: `hashSeed(student|cid|attemptID)` (`internal/engine/engine.go`). Do not use `UnixNano`.

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

Run your tests before opening a PR (1000 local, 100 lightweight in CI `ci.yml:27` for speed, `make fuzz` is 1000):

```bash
go test ./... -count 1
go run scripts/validate_graph.go
python3 scripts/audit_lessons.py
go test ./internal/generator/... -run TestFuzz -count 1000  # local full; CI runs -count 100
make e2e                                                     # browsers; CI runs it sharded
```

`make e2e` is not optional if you touched a route, a button label or where a
component lives. The unit suite renders components in isolation and cannot see
any of those: a renamed button and a page that moved both reached `main` green
on `vitest` while six Playwright tests failed.

#### The browser suite never talks to the deployed API

The browser suite runs against a static export. CI builds it with
`NEXT_PUBLIC_API_URL` set to `http://192.0.2.1:9099`. That address is reserved
for documentation and never answers.

The export also cannot reach any other origin. `e2e/fixtures.ts` registers a net
on the browser context. A request to an origin that is not the export is
recorded, refused, and the test fails. The message names the method and the URL:

```
this test reached an origin outside the static export:
  POST http://192.0.2.1:9099/api/study/answer
Stub it in the spec with page.route, or add the endpoint to e2e/fixtures.ts if
every page needs it. Production must never answer a browser test.
```

The net is why a missing stub now fails the run. Before it existed, the build
pointed at the deployed API. Production answered a page that skipped a stub, so
the test passed on a fast reply and waited on a slow one.

`make e2e-hermetic` builds the export the same way CI does. Run it before you
change a stub:

```bash
make e2e-hermetic
```

Then run one test on its own, and then the whole shard:

```bash
npx playwright test e2e/admin.spec.ts:300 --project=chromium
PW_WORKERS=1 CI=1 npx playwright test --shard=1/2     # about three minutes
```

Fix an escape in the spec that caused it:

```ts
await page.route('**/api/study/answer', route => route.fulfill({ json: { correct: true } }))
```

Add the route to `e2e/fixtures.ts` only when every page asks for that endpoint.
The fixture holds the anonymous shell: no session, default settings, an empty
graph. A spec that asserts content registers its own route and wins, because its
route is later.

A stub must match the whole request URL. A pattern that ends in a name matches
that name only, and not the same path with a query string. Use a trailing `**`
for an endpoint that carries parameters.

Service workers are blocked for the same reason. The export registers
`public/sw.js`, and its fetch handler bypasses `page.route`. `playwright.config.ts`
sets `serviceWorkers: 'block'`, so no cached answer can hide a missing stub.

The diagnosis harness is opt-in. It measures typing latency and scenario
failures, and it does not assert, so it cannot fail a run. Run it when you need
it:

```bash
make e2e-diagnosis
```


### 4. (Optional) Write a lesson

Mathua includes a built-in lesson system. Each concept can have an associated lesson, a markdown file that teaches the material. Lessons live under `data/lessons/` and are registered in `data/lessons/lessons.json`, which maps concept IDs to file paths:

```json
  { "concept_id": "calc.deriv.power_rule", "source": "algebrica/polynomials/polynomials.md" }
```

Multiple concepts can share a single lesson file. Lesson content is written in Markdown with LaTeX via `\( ... \)` or `\[ ... \]` delimiters and rendered inside the app as a sidebar panel alongside practice problems.

### 5. (Recommended) Add knowledge-point shards

Every lesson is broken into exactly 3 **knowledge points** (KPs), each a subgoal-labeled worked example, the Math Academy scaffolding pattern. A KP shard for a concept is a JSON file at `data/lessons/kp/<concept_id>.json`:

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

- `section` names a heading **inside the lesson body** that serves as the worked example. It must match a heading exactly. Do not leave it empty. Since ADR-045, an empty `section` gives the learner a short placeholder that points at Study. It does not show the lesson body.
- `subgoals` are short labels shown before practice (subgoal-labeling effect). Keep 3 per KP.
- Every concept has exactly **3 KPs**. All 1250 shards in the corpus do, and `scripts/check_docs.py` fails if a shard has a different count.
- `scripts/gen_kp_shards.py` generates shard files from an inline spec and validates every section reference against the corpus, add your spec entries there and run `python3 scripts/gen_kp_shards.py <domain-prefix>`.

## The graph validator

A validator runs on every pull request. It checks invariants before any merge can happen:

1. **No cycles**, including `prerequisites` + `encompasses` edges (`internal/concepts/loader.go` Kahn’s)
2. **No orphans**, every `prerequisites`/`encompasses`/`key_prerequisites` must exist (`internal/concepts/loader.go`)
3. **No dup IDs, no empty label/domain, valid `mastery_threshold`** (`internal/concepts/loader.go`)
4. **No singleton `interference_group`** (`internal/concepts/loader.go`), valid `grading_type` enum (9 values `internal/concepts/loader.go`), `variants` `0.1–1.0` difficulty

A second audit guards the lesson corpus (`python3 scripts/audit_lessons.py:191`). It fails on these conditions:

- a DAG concept with no lesson
- a stale `lessons.json` id
- a lesson file that is missing on disk
- an orphaned source
- a KP shard section that does not resolve (all 1250 shard files × 3 = 3750 KPs)
- a diagram mapping that points at a missing asset or a non-existent concept (`engine.go:743`, 183 diagrams)
- a stale course target
- a `grading_type` or `threshold` that fails its sanity check
- incorrect `enrichment.json` wiring

## Submitting a pull request

1. Add the concept to the appropriate domain file in `data/concepts/` with correct prerequisites.
2. Write the generator in the appropriate domain subdirectory.
3. Write the fuzz test with 1 000 samples.
4. (Optional) Write a lesson and register it in `data/lessons/lessons.json`.
5. (Recommended) Add KP shards under `data/lessons/kp/`.
6. Run `go test ./...`, `go run scripts/validate_graph.go` and `python3 scripts/audit_lessons.py` locally.
7. Content waves change the corpus size, run `python3 scripts/counts.py --write` to refresh the counts in the docs (CI runs `--check` and fails on drift).
7. Open a PR. The CI pipeline runs the validator and all tests automatically.
8. A maintainer reviews the concept ordering, thresholds and generator quality.

Reviews usually happen within a few days. If a week passes with no response, ping the thread. We read every PR.

## Full guide

For the detailed walkthrough with field schemas, code examples and design conventions, see the **[Contributing docs](https://mathua.vercel.app/docs/contributing)** in the web app.
