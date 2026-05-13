<div align="center">

<pre style="font-family: monospace; font-size: 11px; line-height: 1.2;">
███╗   ███╗ █████╗ ████████╗██╗  ██╗██╗   ██╗ █████╗
████╗ ████║██╔══██╗╚══██╔══╝██║  ██║██║   ██║██╔══██╗
██╔████╔██║███████║   ██║   ███████║██║   ██║███████║
██║╚██╔╝██║██╔══██║   ██║   ██╔══██║██║   ██║██╔══██║
██║ ╚═╝ ██║██║  ██║   ██║   ██║  ██║╚██████╔╝██║  ██║
╚═╝     ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝
</pre>

**Math Understanding Agent**

*An open-source adaptive math learning engine — counting to calculus, mastery-gated, locally-first.*

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![React](https://img.shields.io/badge/React-18+-61DAFB?style=flat-square&logo=react&logoColor=black)](https://react.dev)
[![License](https://img.shields.io/badge/license-MIT-c8a96e?style=flat-square)](LICENSE)
[![Concepts](https://img.shields.io/badge/concepts-291-a8a0f0?style=flat-square)](#content)
[![Domains](https://img.shields.io/badge/domains-16-c8a96e?style=flat-square)](#content)

</div>

---

Mathua is a local-first adaptive math learning engine inspired by the mastery-gating philosophy of Math Academy. It guides a learner from counting through calculus through a dependency graph of 291 atomic concepts — never advancing until speed and accuracy thresholds are both met. It runs as a single Go binary: desktop TUI for solo practice, or a web server with React + Next.js for global leaderboards and graph visualisation.

---

## What makes it different

- **You cannot advance until you have mastered the prerequisite.** The concept graph enforces this — if you cannot reliably add fractions with different denominators, the system will not show you mixed number addition.
- **Speed matters, not just accuracy.** A correct answer that took 45 seconds on a concept with a 12-second threshold counts as weak mastery. You need to be both right and fast.
- **Nothing is forgotten.** Concepts resurface automatically through spaced repetition. A concept mastered three weeks ago will reappear before it decays.
- **Problems are generated, not stored.** Every problem is produced on demand by a parameterised generator. The same concept gives you a different problem every time. There is nothing to memorise.

---

## Getting started

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

> **Requirements:** Go 1.21+.

---

## How it works

### The concept graph

All math knowledge in Mathua is a directed acyclic graph. Each node is an atomic concept — something small enough to fail independently. Each edge is a prerequisite relationship.

The graph contains **291 concepts** across **17 domains**. The scheduler reads your progress and the graph structure to decide what to show next — it will never surface a concept whose prerequisites you have not mastered.

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

| Domain | Concepts |
|---|---|---|
| Counting | 10 |
| Arithmetic | 46 |
| Fractions | 21 |
| Pre-Algebra | 37 |
| Algebra | 43 |
| Geometry | 20 |
| Trigonometry | 14 |
| Precalculus | 0 |
| Calculus | 26 |
| Statistics | 12 |
| Linear Algebra | 16 |
| Discrete Math | 14 |
| Complex Numbers | 8 |
| Number Theory | 8 |
| Differential Equations | 8 |
| Abstract Algebra | 5 |
| Topology | 3 |
| **Total** | **291** |

---

## Architecture

Mathua is a single Go binary with two modes. The engine core is identical — only the delivery layer differs.

```
Web (React + Next.js) ═══ REST API ═══┐
                                       ├── Core Engine ── Grading ── Storage
Desktop TUI (Bubble Tea) ═══ direct ───┘    (Scheduler · Generators · SM-2   Numeric · Polynomial   SQLite · Postgres)
```

The five-layer architecture — UI, API, Core Engine, Grading, Storage — is fully documented at **[`docs/architecture.md`](docs/architecture.md)** and in the interactive **[web docs](/docs/architecture)** with D2 diagrams.

---

## Contributing

The most impactful contributions are new concepts and improved generators. Adding a concept requires exactly three things: a JSON entry in the concept graph, a Go generator function, and a fuzz test. See **[`CONTRIBUTING.md`](CONTRIBUTING.md)** for the quick-start guide and the full **[contributing docs](/docs/contributing)** for the detailed walkthrough with code examples and field schemas.

---

## Running tests

```bash
# All tests
go test ./...

# Validate the concept graph (no cycles, no orphans)
go run scripts/validate_graph.go

# Generator fuzz — 1000 samples per generator
go test ./internal/generator/... -run TestFuzz -count 1000
```

---

## Self-hosting

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o mathua ./cmd/mathua

# Run
DATABASE_URL=postgres://user:pass@host/mathua \
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

MIT License

</div>
