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

*An open-source adaptive math learning engine — arithmetic to calculus, mastery-gated, locally-first.*

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![React](https://img.shields.io/badge/React-18+-61DAFB?style=flat-square&logo=react&logoColor=black)](https://react.dev)
[![License](https://img.shields.io/badge/license-MIT-c8a96e?style=flat-square)](LICENSE)
[![Concepts](https://img.shields.io/badge/concepts-580-a8a0f0?style=flat-square)](#content)
[![Domains](https://img.shields.io/badge/domains-17-c8a96e?style=flat-square)](#content)

</div>

---

Mathua is a local-first adaptive math learning engine inspired by the mastery-gating philosophy of Math Academy. It guides a learner from arithmetic through university math through a dependency graph of 580 atomic concepts and 21 courses — never advancing until speed and accuracy thresholds are both met. Knowledge is scaffolded into 1740 worked examples (580 KP shard files × 3 subgoals each) with subgoal labels and 115 dual-coded diagrams. It runs as a single Go binary with a web server delivery mode and React + Next.js for leagues, transcripts, and graph visualisation.

---

## What makes it different

- **You cannot advance until you have mastered the prerequisite.** The concept graph enforces this — if you cannot reliably add fractions with different denominators, the system will not show you mixed number addition.
- **Speed matters, not just accuracy.** A correct answer that took 45 seconds on a concept with a 12-second threshold counts as weak mastery. You need to be both right and fast.
- **Nothing is forgotten.** Concepts resurface automatically through per-topic spaced repetition (SM-2 scaled by your learning speed). A concept mastered three weeks ago will reappear before it decays.
- **Problems are generated, not stored.** Every problem is produced on demand by a parameterised generator. The same concept gives you a different problem every time. There is nothing to memorise.
- **Weekly leagues and shareable progress.** Bronze → Diamond leagues promote the top 2 each Monday, and any student can generate a read-only share link for a parent or teacher. Pauses, accommodated timing, and a 150 XP mastery-check quiz are built in.

---

## Getting started

### Build from source
```bash
git clone https://github.com/chuma-beep/mathua.git
cd mathua

go build -o bin/mathua ./cmd/mathua

# SQLite (zero-config, local)
./bin/mathua --serve --port 8080
# Or with DATABASE_URL explicitly
DATABASE_URL=sqlite:mathua.db ./bin/mathua --serve --port 8080
# Postgres for multi-user (via pgx/v5, implemented internal/storage/postgres.go:7)
DATABASE_URL=postgres://user:pass@host/mathua ./bin/mathua --serve --port 8080
# Or copy env file first
cp .env.example .env  # edit JWT_SECRET / CORS_ALLOWED_ORIGINS if needed
```

> **Requirements:** Go 1.21+. Python 3.8+ with `sympy` (optional — enables mathematical expression equivalence grading for algebra and beyond).
>
> **Env:** `DATABASE_URL` (`sqlite:mathua.db` default, or `postgres://` for production via `pgx` `postgres.go:7` / `main.go:57`), `PORT` (default `8080`), `JWT_SECRET` (hex 32 bytes or raw ≥32 chars, else random per run `internal/auth/auth.go:22`), `CORS_ALLOWED_ORIGINS` (comma-separated allowlist, empty = allow all dev `internal/server/server.go:32`). See `.env.example`.

---

## How it works

### The concept graph

All math knowledge in Mathua is a directed acyclic graph. Each node is an atomic concept — something small enough to fail independently. Each edge is a prerequisite relationship.

The graph contains **580 concepts** across **17 domains** (580 KP shard files × 3 = 1740 KPs). The scheduler reads your progress and the graph structure to decide what to show next — it will never surface a concept whose prerequisites you have not mastered.

### Mastery gating

Every concept has two thresholds: a **streak** (consecutive correct answers required) and a **time limit** (maximum acceptable average response time in seconds). Both must be met simultaneously.

A concept moves through five states:

```
UNSEEN → LEARNING → PRACTICING → MASTERED → DECAYING
```

`DECAYING` means you mastered the concept but have not seen it in over 14 days. The scheduler treats decaying concepts as urgent — they appear before new material until you demonstrate you still know them.

### The scheduler

```
priority = 0.5×days_since_last_seen + 0.2×(1-mastery) + 0.3×weakness
         + 5  if DECAYING
         + 2  if encompasses a weak prerequisite (layering)
         − 3  if sharing an InterferenceGroup with recent practice
interleaving: no two consecutive questions share a subdomain
top-3 dissimilar candidates → 70/30 new/review balance
```

Three hard rules apply: never surface a concept whose prerequisites are not mastered, never show the same concept twice in a row, maintain roughly 70% new/practicing and 30% review.

### Spaced repetition

When a concept reaches MASTERED, Mathua schedules its next review using SM-2 scaled by your per-topic learning speed (0.5 slow → 2.0 fast). The interval grows with each successful review (1 day → 3 days → 1 week → ...) and resets on failure. Reviews are woven into normal sessions; every 150 XP a mastery-check quiz surfaces at 80% difficulty. Diagnostic reports include a frontier placement, gaps by domain, and completion estimates for accreditation.

---

## Content

| Domain | Concepts |
|--------|----------|
| Arithmetic | 46 |
| Fractions | 26 |
| Pre-Algebra | 40 |
| Algebra | 84 |
| Geometry | 22 |
| Trigonometry | 22 |
| Precalculus | 23 |
| Calculus | 71 |
| Statistics | 36 |
| Linear Algebra | 25 |
| Discrete Math | 25 |
| Complex Numbers | 25 |
| Number Theory | 25 |
| Differential Equations | 25 |
| Abstract Algebra | 30 |
| Machine Learning | 25 |
| Topology | 30 |
| **Total** | **580** |

---

## Architecture

Mathua is a single Go binary with one delivery method. The engine core is identical — only the UI delivery layer differs.

```
Web (React + Next.js) ═══ REST API ═══ Core Engine ── Grading ── Storage
```

The five-layer architecture — UI, API, Core Engine, Grading, Storage — is fully documented at **[`docs/architecture.md`](docs/architecture.md)** and in the interactive **[web docs](/docs/architecture)** with Mermaid diagrams.

---

## Courses and transcripts

21 courses from 4th grade to university are wired through the DAG — including Calculus I/II, Linear Algebra, Discrete Math I/II, Probability & Statistics, Differential Equations, Abstract Algebra I/II, Topology, and Machine Learning. Each course shows mastered/total, percent, and an estimate of days remaining at your daily XP goal. Transcripts export as CSV and the read-only share link lets a parent or teacher follow along.

## Contributing

The most impactful contributions are new concepts and improved generators. Adding a concept requires exactly three things: a JSON entry in the concept graph, a Go generator function, and a fuzz test. Adding a lesson shard requires 3 KPs per concept (`data/lessons/kp/<id>.json`) with verified section refs. See **[`CONTRIBUTING.md`](CONTRIBUTING.md)** for the quick-start guide and the full **[contributing docs](/docs/contributing)** for the detailed walkthrough with code examples and field schemas.

---

## Running tests

```bash
# All tests
go test ./...

# Validate the concept graph (no cycles, no orphans)
go run scripts/validate_graph.go
python3 scripts/audit_lessons.py  # lessons + KP shards + diagrams + course targets (580 files ×3 =1740 KPs)

# Generator fuzz — 1000 local (make fuzz), 100 lightweight in CI (.github/workflows/ci.yml:27)
go test ./internal/generator/... -run TestFuzz -count 1000  # local full; CI runs -count 100 for speed
```

---

## Self-hosting

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o mathua ./cmd/mathua

# Run (SQLite — local)
./mathua --serve --port 8080
# Or Postgres for production (implemented via pgx internal/storage/postgres.go:7)
DATABASE_URL=postgres://user:pass@host/mathua PORT=8080 ./mathua --serve

# Or with Docker Compose (SQLite default, Postgres supported via DATABASE_URL)
docker compose up
# Or with env file
cp .env.example .env && docker compose up
```

The server is available at `http://localhost:8080`. A `Dockerfile` and `docker-compose.yml` are included in the repo.

---

## Acknowledgements

Mathua's sequencing philosophy is directly inspired by [Math Academy](https://mathacademy.com) — a commercial product that demonstrated what mastery-gated, dependency-aware math learning can achieve. Mathua uses none of their content or code — it is an independent open implementation of similar ideas, built openly.

Content design references: [OpenStax](https://openstax.org) (open textbooks), [MIT OpenCourseWare](https://ocw.mit.edu) (curriculum structure), [Art of Problem Solving](https://artofproblemsolving.com) (problem quality).

Lesson content and diagrams sourced from [Algebrica](https://algebrica.org) by Antonio Lupetti — used under [CC BY-NC 4.0](https://creativecommons.org/licenses/by-nc/4.0/). Algebrica provides a free, ad-free university-level mathematics knowledge base.

---

<div align="center">

MIT License

</div>
