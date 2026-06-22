# Mathua System Design

Mathua is an open-source adaptive mathematics learning platform. A single Go binary with web delivery mode  --  Next.js/React static export  --  sharing one engine core.

---

## Table of Contents

1. [High-Level Architecture](#1-high-level-architecture)
2. [Technology Choices](#2-technology-choices)
3. [Database Design](#3-database-design)
4. [Core Engine](#4-core-engine)
5. [SM-2 Algorithm & Mastery](#5-sm-2-algorithm--mastery)
6. [Grading System](#6-grading-system)
7. [API Layer](#7-api-layer)
8. [Frontend Architecture](#8-frontend-architecture)
9. [End-to-End Data Flow](#9-end-to-end-data-flow)
10. [Computerized Adaptive Testing (CAT)](#10-computerized-adaptive-testing-cat)
11. [Key Design Decisions & Trade-offs](#11-key-design-decisions--trade-offs)

---

## 1. High-Level Architecture

Mathua follows a five-layer architecture. All layers are compiled into a single Go binary. The UI layer has a single implementation  --  the web frontend (React/Next.js served as a static export)  --  calling into the same engine through the REST API.

![System Architecture](diagrams/system-design-architecture-dark.svg)

### Layer Breakdown

| Layer | Technology | Responsibility |
|-------|-----------|----------------|
| **UI** | React/Next.js (static export) | Rendering, user interaction, client state |
| **API** | `net/http` (Go stdlib) | REST endpoints, JWT auth, rate limiting, CORS |
| **Core Engine** | Go | DAG loading, SM-2 scheduling, problem generation, mastery tracking, scoring, weakness propagation, CAT diagnostic |
| **Grading** | Go + SymPy (Python subprocess) | 6+ grading strategies dispatched by grading type |
| **Storage** | SQLite (local/dev), PostgreSQL (production) | Persistence via `Repository` interface |

The architecture is streamlined with a single UI delivery method through the REST API.

---

## 2. Technology Choices

| Component | Choice | Rationale |
|-----------|--------|-----------|
| **Language** | Go | Single binary, fast compilation, excellent concurrency for subprocess management |
| **Database** | SQLite / PostgreSQL | SQLite for zero-config offline use; PostgreSQL for concurrent web access |
| **ORM** | None (raw SQL with `database/sql`) | Maximum control over query performance; schema is small enough that an ORM adds unnecessary abstraction |
| **Problem Generation** | Go functions + SymPy | All problems are procedurally generated  --  no static question bank. SymPy handles symbolic math equivalence checking |
| **Frontend** | Next.js 14 (static export) | Server-rendered React, static files served by the Go binary itself |
| **Styling** | Tailwind CSS + CSS custom properties | Dark/light theme through CSS variables, no runtime JS for theming |
| **Auth** | JWT (HS256) with bcrypt | Stateless auth, 30-day tokens, rate-limited endpoints |
| **Diagrams** | D2 | Declarative, auto-layout, pre-rendered to SVG for use in docs and the web app |

---

## 3. Database Design

### Schema Overview

Five core tables: `students`, `concept_progress`, `sessions`, `attempts`, `questions`.

![Data Model](diagrams/system-design-data-model-dark.svg)

### Table Details

**`students`**  --  User accounts

| Column | Type | Notes |
|--------|------|-------|
| `id` | TEXT PK | UUID |
| `name` | TEXT | Display name |
| `username` | TEXT | Unique login identifier |
| `password_hash` | TEXT | bcrypt hash |
| `course_id` | TEXT | Active course filter |
| `xp_total` | INTEGER | Lifetime XP |
| `xp_today` | INTEGER | Resets via `xp_date` check |
| `settings` | TEXT (JSON) | Arbitrary key-value settings |
| `diagnostic_completed` | INTEGER | Boolean flag |

**`concept_progress`**  --  Per-concept student state with SM-2 spaced repetition

| Column | Type | Notes |
|--------|------|-------|
| `student_id` | TEXT PK FK | Composite with concept_id |
| `concept_id` | TEXT PK | Maps to DAG concept |
| `status` | TEXT | UNSEEN → LEARNING → PRACTICING → MASTERED |
| `streak` | INTEGER | Consecutive correct answers |
| `sm2_repetitions` | INTEGER | SM-2 repetition count |
| `sm2_interval` | INTEGER | Days until next review |
| `sm2_efactor` | REAL | Easiness factor (≥ 1.3) |
| `next_review_due` | TEXT | ISO 8601 date for spaced repetition |
| `weakness_score` | REAL | 0.0–1.0, propagated to dependents |

**`sessions`**  --  Practice sessions

| Column | Type | Notes |
|--------|------|-------|
| `id` | TEXT PK | UUID |
| `student_id` | TEXT FK | |
| `started_at` | TEXT | ISO 8601 timestamp |

**`attempts`**  --  Individual answer records

| Column | Type | Notes |
|--------|------|-------|
| `id` | INTEGER PK | Auto-increment |
| `session_id` | TEXT FK | |
| `concept_id` | TEXT | |
| `answer` | TEXT | Raw student input |
| `expected` | TEXT | Correct answer |
| `correct` | INTEGER | Boolean (0/1) |
| `elapsed_seconds` | REAL | Response time |
| `timestamp` | TEXT | |

**`questions`**  --  Optional pre-seeded bank

| Column | Type | Notes |
|--------|------|-------|
| `id` | INTEGER PK | |
| `concept_id` | TEXT | |
| `question` | TEXT | |
| `answer` | TEXT | |
| `explanation` | TEXT | |
| `difficulty` | REAL | 0.0–1.0 |

### Indexes

- `concept_progress(student_id)`  --  Fast progress lookups
- `sessions(student_id)`  --  Session history
- `attempts(session_id)`  --  Attempts per session
- `attempts(student_id, timestamp)`  --  Time-series queries
- `questions(concept_id)`  --  Question lookup

### SQLite Configuration

```sql
PRAGMA journal_mode=WAL;      -- Concurrent reads during writes
PRAGMA foreign_keys=ON;       -- Referential integrity
PRAGMA busy_timeout=5000;     -- 5s wait before locking error
```

### Migrations

Migrations are embedded directly in the application binary  --  no external migration tool. All DDL uses `CREATE TABLE IF NOT EXISTS`. Incremental migrations use `ALTER TABLE ... ADD COLUMN` guarded by error-checking for duplicate column names.

### Dual-Database Strategy

The `Repository` interface abstracts the storage layer:

- **SQLiteStore**  --  Fully implemented, used for local development
- **PostgresStore**  --  Stub (not yet implemented; every method returns "not implemented")

The codebase checks `DATABASE_URL` at startup and falls back to `mathua.db` with a warning if PostgreSQL is specified.

---

## 4. Core Engine

The engine is the central orchestrator. It holds references to every subsystem and coordinates the full practice loop.

### Modules

| Module | File | Responsibility |
|--------|------|----------------|
| **DAG Loader** | `internal/dag/` | Loads `data/concepts/*.json`, validates no cycles or orphaned prerequisites, produces topological sort via Kahn's algorithm |
| **Scheduler** | `internal/scheduler/` | Selects next concept using priority scoring (recency, mastery, weakness). Maintains 70/30 new/review balance |
| **Generators** | 18 domain registries | Every problem is generated on demand by a parameterized Go function  --  no static question bank |
| **Mastery Machine** | `internal/mastery/` | State machine: UNSEEN → LEARNING → PRACTICING → MASTERED (with DECAYING at read time) |
| **Scoring** | `internal/scoring/` | Two scores: lifetime topic score (permanent) and weekly score (resets Monday) |
| **Weakness Propagation** | In engine | If a concept's weakness > 0.3, propagates `w * 0.3` to dependents |
| **Diagnostic** | `internal/diagnostic/` | Binary search CAT on topologically sorted concept graph |
| **Planner** | `internal/planning/` | Course paths loaded from `data/courses.json` |
| **Lessons** | `internal/lessons/` | Markdown lesson content from `data/lessons/` |

### DAG (Concept Graph)

Concept relationships are defined as JSON files:

```json
{
  "id": "algebra.quadratic-formula",
  "label": "Quadratic Formula",
  "domain": "algebra",
  "subdomain": "equations",
  "grading_type": "numeric",
  "prerequisites": ["algebra.quadratic-equations", "algebra.square-roots"],
  "mastery_threshold": { "streak": 3, "avg_time_seconds": 0 }
}
```

The DAG loader:
1. Loads all JSON files from `data/concepts/`
2. Validates: no cycles (DFS cycle detection), no orphaned prerequisites, no duplicate IDs
3. Produces a topological sort (Kahn's algorithm)
4. Builds bidirectional maps for `prereqsOf` and `dependentsOf`

### Scheduler

The scheduler selects the next concept using a priority formula:

```
priority = 0.5 × daysSinceLastSeen
         + 0.2 × (1 - masteryScore)
         + 0.3 × weaknessScore
```

Decaying concepts (≥14 days since last review) get `+5` priority bonus. The scheduler then maintains a **70/30 new/review ratio** when selecting from the candidate pool  --  70% of questions should be new material, 30% should be spaced repetition reviews.

---

## 5. SM-2 Algorithm & Mastery

![SM-2 & Mastery](diagrams/system-design-sm2-mastery-dark.svg)

### Mastery State Machine

Four stored states, one computed state:

| State | Transition | Condition |
|-------|-----------|-----------|
| **UNSEEN** → LEARNING | First correct | Streak ≥ required AND avg time ≤ threshold |
| **LEARNING** → PRACTICING | Consistent correct | Same |
| **PRACTICING** → MASTERED | Mastery achieved | Same |
| **MASTERED** → DECAYING | 14+ days idle | Computed at read time, never stored |

The `streak` resets to 0 on any incorrect answer. Each concept has a `required_streak` (typically 3) and a `time_threshold` (typically 120s for generated, 60s for review).

### SM-2 Spaced Repetition

The SM-2 algorithm determines review intervals:

```
quality < 3:
  rep = 0, interval = 1, ef unchanged

quality >= 3:
  rep++
  ef += 0.1 - (5 - q) × (0.08 + (5 - q) × 0.02)
  ef = max(1.3, ef)
  interval:
    rep == 1 → 1 day
    rep == 2 → 6 days
    rep > 2 → round(prev_interval × ef)
```

Quality is derived from correctness and response time:

| Condition | Quality |
|-----------|---------|
| Streak met AND time ratio ≤ 0.5 | 5 |
| Streak met AND time ratio ≤ 1.0 | 4 |
| Streak met AND time ratio > 1.0 | 3 |
| Streak not met | 0 |

### XP Scoring

```
base = 10 (new) or 5 (review)
timeMultiplier = clamp(2.0 - elapsed / threshold, 0.5, 1.5)
streakMultiplier = 1.0 + min(streak, 10) × 0.1
XP = int(base × timeMultiplier × streakMultiplier)
```

### Levels

| Concepts Mastered | Title |
|-----------------|-------|
| 0 | Novice |
| 32 | Apprentice |
| 64 | Student |
| 96 | Scholar |
| 128 | Adept |
| 160 | Expert |
| 192 | Master |
| 224 | Grandmaster |
| 256 | Math Architect |

---

## 6. Grading System

The grading system uses a strategy pattern with a central Router that dispatches to the correct grader based on `grading_type`.

![Grading Pipeline](diagrams/system-design-grading-pipeline-dark.svg)

### Grader Strategies

| Grader | Grading Types | Implementation |
|--------|--------------|----------------|
| **Numeric** | `numeric` | Parses integers, decimals, fractions (`3/4`), mixed numbers (`1 1/2`), scientific notation (`1e2`). Strips commas and leading `+`. Uses `big.Rat` for exact rational arithmetic. Float tolerance `1e-9`. |
| **Multiple Choice** | `multiple_choice` | Case-insensitive equality + single-letter matching (e.g., `"B"` matches `"Option B"`). |
| **Comparison** | `comparison` | Operator matching: `>`, `<`, `=`, `>=`, `<=`, `!=`, `==`. |
| **Ordering** | `ordering` | Splits on whitespace/comma/semicolon/pipe. Positional exact match (order-sensitive). |
| **Tuple** | `tuple` | Strips parentheses and semicolons, splits on comma. Positional exact match. |
| **Complex** | `complex` | Preprocesses polar form `r(cos + i sin)`, handles `±`, converts `i` → `I`, delegates to SymPy. |
| **Symbolic (fallback)** | (various) | String equality after normalizing spaces and `**` → `^`. Used when SymPy is unavailable. |
| **SymPy** | `polynomial`, `expression` | Long-lived Python subprocess. Sends JSON via stdin, reads JSON from stdout. |

### SymPy Subprocess Integration

For symbolic math (polynomial simplification, expression equivalence), Mathua spawns a Python subprocess running `grading/sympy_service.py`:

1. **Input**: JSON line on stdin `{id, expected, answer}`
2. **Preprocessing**: Strip variable assignments (`y = ...`), strip `+ C` constants, normalize `e^x` → `exp(x)`, `ln` → `log`, `|x|` → `Abs(x)`, handle implicit multiplication (`2sin(x)` → `2*sin(x)`)
3. **Parsing**: `sympy.parsing.sympy_parser` with implicit multiplication support
4. **Symbolic equivalence**: `simplify(diff) == 0`, `expand(diff) == 0`, `trigsimp`, `radsimp`, `powsimp`, `factor`, `ee.equals(ea)`
5. **Numerical verification**: 10 random points checked within `1e-6` tolerance
6. **Output**: JSON line on stdout `{id, correct, feedback}`
7. **Safeguards**: 10-second `SIGALRM` timeout, 500-character input limit, character allowlist

If Python/SymPy are not installed, grading falls back gracefully to the pure-Go symbolic string normalizer.

---

## 7. API Layer

The web server exposes a REST API through Go's standard `net/http` package. No external HTTP framework.

### Middleware Stack

| Middleware | Purpose |
|-----------|---------|
| `cors` | Sets `Access-Control-Allow-Origin`, handles OPTIONS (204) |
| `logRequest` | Logs method and path |
| `authMiddleware` | Reads `Authorization: Bearer <token>`, validates JWT, injects `studentID` into context |
| `authLimiter` | Token-bucket rate limiter (5 rate, 10 burst, 1 min window) on auth endpoints |

### All Routes

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/api/auth/signup` | rate-limited | Create account → JWT |
| POST | `/api/auth/login` | rate-limited | Login → JWT |
| GET | `/api/auth/me` | Bearer | Current user + scores |
| POST | `/api/session` | optional | Start practice session |
| POST | `/api/answer` | optional | Submit answer → result + next |
| GET | `/api/progress/` | optional | All concept progress |
| GET | `/api/scores/` | optional | Aggregated scores |
| GET | `/api/config` | no | `{auth_enabled: bool}` |
| GET | `/api/graph` | no | Full DAG node list |
| GET | `/api/leaderboard` | no | Weekly leaderboard |
| GET | `/api/courses` | auth | Available courses |
| POST | `/api/courses/{id}/diagnostic` | auth | Set course + activate path |
| POST | `/api/diagnostic` | no | Start CAT diagnostic |
| POST | `/api/diagnostic/answer` | no | Submit diagnostic answer |
| GET | `/api/weaknesses` | auth | Weakness scores by domain |
| POST | `/api/goals/xp` | auth | Set daily XP goal |
| GET/PUT | `/api/settings` | auth | Student settings JSON |
| GET | `/api/reviews/due` | auth | Count due reviews |
| POST | `/api/reviews/session` | auth | Create review-only session |
| GET | `/api/lessons` | no | All lessons with progress |
| GET | `/api/concepts/{id}` | optional | Concept detail |
| GET | `/api/health` | no | `{"status":"ok"}` |

### Auth Flow

```
Signup: POST /auth/signup {name, username, password}
  → bcrypt hash password
  → CreateUser in SQLite
  → Generate HS256 JWT with 30-day expiry
  → Return {token, student_id, name}

Login: POST /auth/login {username, password}
  → FindByUsername
  → bcrypt.Compare
  → Generate JWT
  → Return token

Validate: Authorization: Bearer <token>
  → Parse HS256 JWT
  → Extract "sub" → studentID
  → Inject into context
```

JWT secret is from `JWT_SECRET` env var (hex-decoded 32 bytes or raw string ≥ 32 chars), or randomly generated each run.

---

## 8. Frontend Architecture

### Web App (Next.js 14  --  Static Export)

The frontend uses Next.js in static export mode  --  meaning the Go binary serves pre-built HTML/JS/CSS files with no Node.js runtime. This keeps deployment simple: a single Go binary that is both API server and file server.

**Directory structure:**

```
web/next-app/
├── app/                    # Next.js pages (App Router)
│   ├── page.tsx            # Home page
│   ├── diagnose/           # CAT diagnostic flow
│   ├── practice/           # Practice session
│   ├── progress/           # Progress visualization
│   ├── lessons/            # Lesson content
│   ├── leaderboard/        # Weekly scores
│   └── docs/               # Documentation
├── components/             # Shared React components
│   ├── Header.tsx          # Sticky nav bar
│   ├── D2Diagram.tsx       # Diagram renderer
│   ├── FormulaBlock.tsx    # KaTeX formula display
│   └── ...
├── diagrams/               # D2 source files
├── public/diagrams/        # Pre-rendered SVGs
└── tailwind.config.js      # Custom theme
```

**State management:** Each page manages its own state with React hooks (`useState`, `useEffect`). No global state library  --  the API is the source of truth.

**API client:** Axios-based, with a thin wrapper in `lib/`.

**Theming:** CSS custom properties for dark/light mode. Uses a warm palette: dark backgrounds are `#0b0f1a` (not pure black), light surfaces are `#fefcf4` (warm parchment). Gold accent (`#c8a96e`) marks mastery and structure.

---

## 9. End-to-End Data Flow

Here's the complete flow when a student submits an answer:

![Request Flow](diagrams/system-design-request-flow-dark.svg)

### Step-by-Step

1. **Client** sends `POST /api/answer` with `{session_id, answer, elapsed_seconds}`
2. **Auth middleware** validates JWT, injects `studentID` into request context
3. **Handler** calls `Engine.SubmitAnswer(sessionID, studentID, answer, elapsed)`
4. **Grading**  --  route to the correct grader based on `grading_type`: numeric graders parse the answer, SymPy handles symbolic equivalence, etc.
5. **Machine**  --  `mastery.Machine.Next()` determines if the student's performance qualifies for a state transition (UNSEEN → LEARNING → etc.)
6. **SM-2**  --  `scheduler.ComputeSM2()` computes new repetition count, interval, and easiness factor
7. **Upsert**  --  `repo.UpsertProgress()` saves all state (including SM-2 fields) with an SQLite `ON CONFLICT ... DO UPDATE` query
8. **Weakness**  --  if the concept's weakness score exceeds 0.3, propagate to dependent concepts at `w × 0.3`
9. **Attempt**  --  `repo.RecordAttempt()` inserts a row into the `attempts` table
10. **XP**  --  `computeXP()` calculates XP using base, time, and streak multipliers, then `repo.AddXP()` updates the student record
11. **Next**  --  `Engine.NextQuestion()` loads the next concept from the scheduler and asks the generator registry to produce a new problem
12. **Response**  --  JSON response with result, next question, and session state

Total latency target: under 100ms for numeric grading, under 500ms for SymPy-based grading (including subprocess round-trip).

---

## 10. Computerized Adaptive Testing (CAT)

Mathua implements a binary-search CAT to locate a student's knowledge frontier quickly. Instead of testing all ~284 concepts, the diagnostic requires approximately 20–35 questions.

![CAT Diagnostic](diagrams/system-design-cat-diagnostic-dark.svg)

### Algorithm

1. Load all concepts in topological order (from the DAG)
2. Set `Low = 0`, `High = N - 1`
3. `Position = Low + (High - Low) / 2`
4. Generate a question for the concept at `Position`
5. If the student answers **correctly within the time limit**: `Low = Position + 1` (move forward toward harder concepts)
6. If **incorrect or slow**: `High = Position - 1` (move backward toward foundational concepts)
7. After 3 consecutive correct answers at the boundary: **frontier is located**
8. Record a mastery estimate for every concept passed through

### Key Properties

- **Binary search efficiency**: O(log n) probes, ~20–35 questions vs 284 exhaustive
- **Optimistic merge**: Retaking the diagnostic does not delete existing progress. New estimates are merged with existing data, always preferring the more optimistic estimate.
- **Time-bounded**: Each question has a configurable time limit (typically 60s). Slow answers are treated as incorrect for diagnostic purposes.

### Goal-Oriented Diagnostic

A variant that targets a specific set of concept goals. Instead of starting at the DAG midpoint, it:
1. Collects all prerequisite chains for the target concepts
2. Runs binary search only within the prerequisite subspace
3. Produces a readiness assessment: which prerequisites are mastered and which need work

---

## 11. Key Design Decisions & Trade-offs

### Why Go over Python or Node.js?

| Factor | Go | Python | Node.js |
|--------|-------|--------|---------|
| Single binary deploy | ✅ Yes | ❌ Requires runtime | ❌ Requires runtime |
| SymPy subprocess | ✅ Natural with `os/exec` | ✅ Native | ⚠️ Indirect |
| Concurrency | ✅ Goroutines | ⚠️ GIL | ✅ Event loop |
| Compile-time safety | ✅ Strong | ❌ Runtime | ❌ Runtime |
| Startup time | ✅ ~5ms | ❌ ~500ms | ⚠️ ~100ms |

The biggest win: a single ~15MB binary that is both API server and static file server. Deploy with `scp mathua user@server: && ./mathua --serve`.

### Why SQLite + PostgreSQL (not just one)?

| Scenario | Database | Reason |
|----------|----------|--------|
| Student running locally | SQLite | Zero config, no server, file-based |
| Self-hosted / homelab | SQLite | Same binary, just point at a path |
| Multi-user web deployment | PostgreSQL | Concurrent writes, connection pooling |
| CI / Tests | SQLite (`:memory:`) | Fast, isolated, reproducible |

The `Repository` interface makes this transparent. The schema is identical across both databases.

### Why Generated Questions (not a static bank)?

- **Infinite variety**: Every session is unique
- **No cheating**: A static bank can be memorized and shared
- **Adaptive difficulty**: Parameters can be tuned in real-time based on student performance
- **No storage**: Questions are generated on-demand and never persisted

Trade-off: generation latency. Complex polynomial questions require a SymPy subprocess round-trip (typically 50–200ms). Simple numeric questions are pure Go and take <1ms.

### Why No ORM?

The schema has 5 tables with straightforward relationships. Raw SQL with `database/sql` gives:
- Full control over query optimization (especially the SM-2 upsert and leaderboard queries)
- No query generation overhead
- Transparent SQL for debugging
- Easy to port between SQLite and PostgreSQL

Trade-off: more boilerplate for CRUD operations. Worth it for the control over the SM-2 upsert query and leaderboard computation.

### Why D2 over alternatives?

The project previously used Mermaid. D2 provides:
- Superior auto-layout (especially for directed graphs with many nodes)
- Cleaner syntax  --  declarative, readable, version-control friendly
- First-class theme support (light/dark themes via `-t 0` / `-t 200`)
- Pre-rendered SVGs  --  no runtime rendering cost

### Why Static Next.js Export (not SSR)?

The Go binary serves the Next.js static export directly. This means:
- No Node.js server needed in production
- A single port for both API and static files
- Simpler deployment  --  just the Go binary plus the `out/` directory
- Trade-off: no server-side rendering, but for a learning app this is acceptable

---

## Summary

Mathua is an intentionally simple system. One language (Go), two databases (SQLite/PostgreSQL abstracted by an interface), one engine core with a single delivery method through the REST API. The complexity is in the algorithms  --  SM-2 spaced repetition, CAT binary-search diagnostics, 18 problem generators, 6+ grading strategies  --  not in the infrastructure. Every component can be understood by reading its source file start to finish.

