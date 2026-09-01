# Mathua Architecture

Mathua is a single Go binary with a web delivery mode. The engine core serves a single presentation layer.

## Five-layer architecture

```
┌──────────────────────────────────────────────┐
│  UI Layer                                     │
│  Web browser (React + Next.js)                │
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
│  SQLite (local) · PostgreSQL (web)            │
│  data/concepts/ — per-domain DAG files          │
└──────────────────────────────────────────────┘
```

### UI Layer — Web frontend, one engine

Mathua ships with a web frontend using React with Next.js, connecting to a shared Postgres database.

### API Layer — Four services, one router

The web server exposes a REST API through Go's standard `net/http` package. No external HTTP framework. Four services: Session (practice management), Graph (concept DAG visualisation), Leaderboard (weekly scoring), Diagnostic (adaptive testing).

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

SQLite for local development (`mathua.db`, `WAL` `storage/migrate.go:23`). PostgreSQL for production web access via `DATABASE_URL=postgres://` (`internal/storage/postgres.go:7` `PostgresStore` via `pgx/v5/stdlib`, `pgSchema` `postgres.go:7`, `GREATEST` for XP floor, `STRING_AGG` for activity). `cmd/mathua` checks `DATABASE_URL` and opens `pgx` when `postgres://` is set, otherwise SQLite. The data schema is identical across both databases, abstracted behind a `Repository` interface (`storage/store.go:127`).

## Diagnostic — Computerised Adaptive Testing

Locates a student's knowledge frontier using a compressed covering set + info-gain CAT (`internal/diagnostic/cat.go:1`, `internal/diagnostic/report.go:11`). The engine builds a minimal covering set of the DAG, repeatedly picks the concept with maximal entropy reduction, propagates `+0.3` evidence to prerequisites on correct and `-0.3` to dependents on incorrect, and tracks per-concept `KnowledgeConfidence 0–1`. The frontier is the highest belief drop; a supplemental diagnostic runs when confidence `<0.7`. Assessment is 25–45 adaptive questions (vs 580 exhaustive) with `80%` difficulty targeting via `engine.computeDifficulty` (`internal/engine/engine.go:157` `weakness→difficulty` 0.3–1.0).

## Full documentation

For interactive D2 diagrams, internal module detail, and formatted code blocks, see the **[Architecture docs](https://mathua.vercel.app/docs/architecture)** in the web app.
