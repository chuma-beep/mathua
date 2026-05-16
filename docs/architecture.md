# Mathua Architecture

Mathua is a single Go binary with two delivery modes. The engine core is identical — only the presentation layer differs.

## Five-layer architecture

```
┌──────────────────────────────────────────────┐
│  UI Layer                                     │
│  Web browser (React + Next.js)                │
│  Desktop TUI (Bubble Tea)                     │
├──────────────────────────────────────────────┤
│  API Layer                                     │
│  REST API — net/http                          │
│  Session · Graph · Leaderboard · Diagnostic   │
├──────────────────────────────────────────────┤
│  Core Engine                                   │
│  DAG Loader · Scheduler (SM-2) · Generators   │
│  Mastery Tracking · Scoring                   │
├──────────────────────────────────────────────┤
│  Grading                                       │
│  Numeric · Choice · Comparison · Ordering (Go)  │
│  Symbolic (Go) — fallback                       │
│  Polynomial & Expression → SymPy (Python)       │
├──────────────────────────────────────────────┤
│  Storage                                       │
│  SQLite (desktop) · PostgreSQL (web)           │
│  data/concepts/ — per-domain DAG files          │
└──────────────────────────────────────────────┘
```

### UI Layer — Two presentations, one engine

Mathua ships with two user interfaces that share the same engine core. The desktop interface requires no account, no network, and stores everything in SQLite. The web interface uses React with Next.js and connects to a shared Postgres database.

### API Layer — Four services, one router

The web server exposes a REST API through Go's standard `net/http` package. No external HTTP framework. Four services: Session (practice management), Graph (concept DAG visualisation), Leaderboard (weekly scoring), Diagnostic (adaptive testing). The TUI bypasses the API layer entirely — it calls the engine modules directly through Go function interfaces.

### Core Engine — Five modules

- **DAG Loader** — Loads `concepts.json`, validates no cycles and no orphaned prerequisites.
- **Scheduler (SM-2)** — Selects the next concept using priority scoring: recency, mastery, decay bonuses. Enforces prerequisite gating and 70/30 new/review balance.
- **Generators** — Every problem is generated on demand by a parameterised Go function. No static question bank.
- **Mastery & Spaced Repetition** — Tracks mastery via streak and response time. Simplified SM-2 algorithm schedules reviews.
- **Scoring** — Two scores: lifetime topic score (permanent) and weekly score (resets every Monday).

### Grading — Six strategies

Six grader types handled by a single `Router`:

- **Numeric** — integer, float, fraction, mixed number, scientific notation (pure Go).
- **Choice** — case-insensitive multiple-choice matching.
- **Comparison** — comparison operators (`>`, `<`, `=`, `>=`, `<=`).
- **Ordering** — ordered sequences separated by delimiters.
- **Symbolic** — string normalisation fallback for when SymPy is unavailable.

**Polynomial and Expression** — `polynomial` and `expression` grading types route through a long-lived Python subprocess running `grading/sympy_service.py`. The Go client sends a JSON expression pair over stdin; SymPy parses both into expression trees and tests equivalence via `simplify(expected - answer) == 0`. If Python/SymPy are not installed, grading falls back gracefully to the pure-Go `symbolic` string normaliser.

### Storage — SQLite and Postgres

SQLite for fully offline desktop operation. PostgreSQL for concurrent web access and global state. The data schema is identical across both databases, abstracted behind a repository interface.

## Diagnostic — Computerised Adaptive Testing

Locates a student's knowledge frontier using binary search on the topologically sorted concept graph. Reduces the assessment from 291 questions to approximately 20–35.

## Full documentation

For interactive D2 diagrams, internal module detail, and formatted code blocks, see the **[Architecture docs](https://mathua.vercel.app/docs/architecture)** in the web app.
