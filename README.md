<div align="center">

```
███╗   ███╗ █████╗ ████████╗██╗  ██╗██╗   ██╗ █████╗
████╗ ████║██╔══██╗╚══██╔══╝██║  ██║██║   ██║██╔══██╗
██╔████╔██║███████║   ██║   ███████║██║   ██║███████║
██║╚██╔╝██║██╔══██║   ██║   ██╔══██║██║   ██║██╔══██║
██║ ╚═╝ ██║██║  ██║   ██║   ██║  ██║╚██████╔╝██║  ██║
╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝
```

**Math Understanding Agent**

*An open-source adaptive math learning system — counting to algebra, mastery-gated, locally-first.*

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![React](https://img.shields.io/badge/React-18+-61DAFB?style=flat-square&logo=react&logoColor=black)](https://react.dev)
[![License](https://img.shields.io/badge/license-MIT-c8a96e?style=flat-square)](LICENSE)
[![Build](https://img.shields.io/badge/build-passing-4db8a0?style=flat-square)](.github/workflows/ci.yml)
[![Concepts](https://img.shields.io/badge/concepts-81-a8a0f0?style=flat-square)](#content)

</div>

---

Mathua is a local-first adaptive math learning engine inspired by the mastery-gating philosophy of Math Academy. It guides a learner from counting numbers through algebra through a dependency graph of atomic concepts — never advancing until speed and accuracy thresholds are both met.

It runs as a single binary on your machine. No account required for the desktop version. No ads. No subscriptions. Fully open source.

---

## What makes it different

Most math platforms are content libraries. You pick a topic, watch a video, do some problems, and move on regardless of whether you actually understood it. Mathua works differently.

- **You cannot advance until you have mastered the prerequisite.** The concept graph enforces this — if you cannot reliably add fractions with different denominators, the system will not show you mixed number addition.
- **Speed matters, not just accuracy.** Mathua tracks how long each answer takes. A correct answer that took 45 seconds on a concept with a 12-second threshold counts as weak mastery. You need to be both right and fast.
- **Nothing is forgotten.** Concepts resurface automatically through spaced repetition. A concept mastered three weeks ago will reappear before it decays.
- **Problems are generated, not stored.** Every problem is produced on demand by a parameterised generator. The same concept gives you a different problem every time. There is nothing to memorise.

---

## Getting started

### Web — no install required

Visit **[mathua.dev](https://mathua.dev)** — create an account and start immediately in the browser. Progress is saved and syncs across devices.

### Desktop — offline, no account needed

Download the binary for your platform from the [latest release](https://github.com/chuma-beep/mathua/releases/latest):

| Platform | Binary |
|---|---|
| macOS (Apple Silicon) | `mathua-macos-arm64` |
| macOS (Intel) | `mathua-macos-amd64` |
| Linux | `mathua-linux-amd64` |
| Windows | `mathua-windows-amd64.exe` |

```bash
# macOS / Linux
chmod +x mathua-macos-arm64
./mathua-macos-arm64

# Windows
mathua-windows-amd64.exe
```

### Build from source

```bash
git clone https://github.com/chuma-beep/mathua.git
cd mathua
go build ./cmd/mathua

# Run the desktop TUI
./mathua

# Run the web server (requires Postgres)
DATABASE_URL=postgres://... ./mathua --serve --port 8080
```

> **Requirements:** Go 1.21+. For symbolic math grading (algebra concepts in v1.1): Python 3.10+ with SymPy — `pip install sympy`.

---

## How it works

### The concept graph

All math knowledge in Mathua is a directed acyclic graph. Each node is an atomic concept — something small enough to fail independently. Each edge is a prerequisite relationship.

```
count.objects → count.cardinality → count.number_line → count.compare
                                                               ↓
                                          arith.add.single ←──┘
                                                ↓
                                       arith.add.double → arith.add.carry
                                                ↓
                                       arith.sub.single → ...
```

The graph currently contains **81 concepts** across five domains. The scheduler reads your progress and the graph structure to decide what to show next — it will never surface a concept whose prerequisites you have not mastered.

### Mastery gating

Every concept has two thresholds: a **streak** (consecutive correct answers required) and a **time limit** (maximum acceptable average response time in seconds). Both must be met simultaneously.

A concept moves through five states:

```
UNSEEN → LEARNING → PRACTICING → MASTERED → DECAYING
```

`DECAYING` means you mastered the concept but have not seen it in over 14 days. The scheduler treats decaying concepts as urgent — they appear before new material until you demonstrate you still know them.

### The scheduler

```
priority = (0.7 × days_since_last_seen)
         + (0.3 × (1.0 - mastery_score))
         + 5.0  if status == DECAYING
         + 2.0  if newly unlocked
```

Three hard rules apply: never surface a concept whose prerequisites are not mastered, never show the same concept twice in a row, maintain roughly 70% new/practicing and 30% review.

### Spaced repetition

When a concept reaches MASTERED, Mathua schedules its next review using a simplified SM-2 algorithm. The interval grows with each successful review (1 day → 3 days → 1 week → ...) and resets on failure. Reviews are woven into normal sessions — there is no separate review mode.

---

## Content

### Current domains — v1.0

| Domain | Sub-domains | Concepts |
|---|---|---|
| Counting | Objects, cardinality, number line, comparison, skip counting | 7 |
| Arithmetic | Addition, subtraction, multiplication, division, factors, exponents, negatives, order of operations | 39 |
| Fractions | Concept, equivalence, addition, subtraction, multiplication, division, mixed numbers | 18 |
| Pre-Algebra | Decimals, percentages, ratios, intro variables | 17 |
| **Total** | | **81** |

### Roadmap — v1.1

Algebra domain: variables and expressions, linear equations (one-step through multi-step), inequalities, coordinate plane, slope, linear functions, systems of equations, polynomials, factoring, quadratic formula. Approximately 60 additional concepts, with symbolic grading via Python/SymPy.

---

## Ranking and progress

### Personal topic scores

Performance is broken down by domain and sub-domain — not a single opaque number.

```
sub_domain_score = (mastered_concepts × 100)
                 + speed_bonus per concept
                 + domain_completion_bonus (if 100% complete)
```

### Global weekly leaderboard

The web version includes a global leaderboard that resets every Monday.

```
weekly_score = (mastered_concepts × 100) + speed_bonus + (streak_days × 10)
```

Rank resets weekly. Your permanent **level** is based on all-time mastered concepts and never resets.

| Level | Name | Concepts mastered |
|---|---|---|
| 1 – 2 | Newcomer · Counter | 0 – 24 |
| 3 – 4 | Arithmetician · Calculator | 25 – 79 |
| 5 – 6 | Fraction Master · Algebraist | 80 – 149 |
| 7 – 8 | Equation Solver · Proof Reader | 150 – 199 |
| **9** | **Math Architect** | **200 — all mastered** |

---

## Architecture

Mathua is a single Go binary with two modes. The engine core is identical — only the delivery layer differs.

```
┌───────────────────────────────────────────────────────────┐
│  Web browser (React + KaTeX)  │  Desktop TUI (Bubble Tea) │
├───────────────────────────────────────────────────────────┤
│      REST API — net/http  (TUI calls engine directly)     │
├───────────────────────────────────────────────────────────┤
│  DAG loader │ Scheduler │ Mastery │ Generators │ SM-2     │
│  Path finder │ CAT engine │ Timer │ Leaderboard           │
├───────────────────────────────────────────────────────────┤
│  Numeric grader (Go)  │  Symbolic grader (SymPy, v1.1)   │
├───────────────────────────────────────────────────────────┤
│  SQLite (desktop)     │  Postgres (web)                   │
│  concepts.json — community-editable DAG                   │
└───────────────────────────────────────────────────────────┘
```

Full architecture documentation: [`docs/architecture.md`](docs/architecture.md)

---

## Contributing

The most impactful contributions are new concepts and improved generators. Adding a concept requires exactly three things.

**1. A JSON entry in `data/concepts.json`:**

```json
{
  "id":            "arith.add.single",
  "label":         "Single-digit addition",
  "domain":        "arithmetic",
  "subdomain":     "arithmetic.addition",
  "grading_type":  "numeric",
  "prerequisites": ["count.cardinality", "count.number_line"],
  "mastery_threshold": { "streak": 5, "avg_time_seconds": 8.0 }
}
```

**2. A generator in `internal/generator/[domain]/`:**

```go
func (g *AddSingleGen) Generate(difficulty float64) generator.Problem {
    max := int(5 + difficulty*4)
    a   := rand.Intn(max) + 2   // never 0 or 1
    b   := rand.Intn(max) + 2
    ans := a + b
    return generator.Problem{
        Question:    fmt.Sprintf("%d + %d = ?", a, b),
        Answer:      strconv.Itoa(ans),
        Explanation: fmt.Sprintf("%d + %d = %d", a, b, ans),
    }
}
```

**3. A fuzz test asserting 1000 valid samples:**

```go
func TestAddSingleGen(t *testing.T) {
    g := &AddSingleGen{}
    for i := 0; i < 1000; i++ {
        p := g.Generate(rand.Float64())
        assert.NotEmpty(t, p.Question)
        assert.NotEmpty(t, p.Answer)
        val, _ := strconv.Atoi(p.Answer)
        assert.Greater(t, val, 2)
    }
}
```

> The graph validator rejects cycles and orphaned nodes automatically on every PR. See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the full guide and concept schema reference.

### Other ways to contribute

- **Fix a generator** — play the app and note problems that feel wrong, ambiguous, or trivially easy. Open an issue or PR.
- **Improve mastery thresholds** — if a concept's time limit feels too strict or too lenient, propose a change with reasoning.
- **Report prerequisite errors** — if a concept unlocks too early or is blocked incorrectly, open an issue.
- **Translate content** — `concepts.json` supports a `label_[locale]` field for translated names.

---

## Running tests

```bash
# All tests
go test ./...

# Validate the concept graph (no cycles, no orphans)
go run scripts/validate_graph.go

# Generator fuzz — 1000 samples per generator
go test ./internal/generator/... -run TestFuzz -count 1000

# Integration tests
go test ./... -run TestIntegration -tags integration
```

---

## Self-hosting

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o mathua ./cmd/mathua

# Run
DATABASE_URL=postgres://user:pass@host/mathua \
JWT_SECRET=your-secret-here \
PORT=8080 \
./mathua --serve

# Or with Docker Compose
docker compose up
```

The server is available at `http://localhost:8080`. A `Dockerfile` and `docker-compose.yml` are included in the repo.

---

## Acknowledgements

Mathua's sequencing philosophy is directly inspired by [Math Academy](https://mathacademy.com) — a commercial product that demonstrated what mastery-gated, dependency-aware math learning can achieve. Mathua uses none of their content or code — it is an independent open implementation of similar ideas, built openly.

Content design references: [OpenStax](https://openstax.org) (open textbooks), [MIT OpenCourseWare](https://ocw.mit.edu) (curriculum structure), [Art of Problem Solving](https://artofproblemsolving.com) (problem quality).

---

<div align="center">

MIT License &nbsp;·&nbsp; Built in public by [Wisdom Chukwuemeka](https://portfolio3-seven-indol.vercel.app) &nbsp;·&nbsp; [GitHub](https://github.com/chuma-beep)

*Start from counting. End at algebra. Master everything in between.*

</div>