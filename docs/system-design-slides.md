---
marp: true
theme: uncover
class:
  - lead
  - invert
paginate: true
---

# Mathua System Design

A single Go binary, two delivery modes, one engine core.

---

## 1. Five-Layer Architecture

![width:1000](https://raw.githubusercontent.com/chuma-beep/mathua/main/web/next-app/public/diagrams/system-design-architecture-dark.svg)

**UI** (React/Next.js + Bubble Tea) **API** (net/http) **Engine** (Go) **Grading** (Go + SymPy) **Storage** (SQLite/PostgreSQL)

---

## 2. Why Go?

| Factor | Go | Python | Node.js |
|--------|-----|--------|---------|
| Single binary | Yes | No | No |
| Startup time | ~5ms | ~500ms | ~100ms |
| Concurrency | Goroutines | GIL-bound | Event loop |
| SymPy bridge | os/exec | Native | Indirect |

**Result**: ~15MB binary that is both API server and static file server.

---

## 3. Database Schema

![width:1000](https://raw.githubusercontent.com/chuma-beep/mathua/main/web/next-app/public/diagrams/system-design-data-model-dark.svg)

5 tables: `students`, `concept_progress`, `sessions`, `attempts`, `questions`

---

## 4. SM-2 Spaced Repetition

![width:1000](https://raw.githubusercontent.com/chuma-beep/mathua/main/web/next-app/public/diagrams/system-design-sm2-mastery-dark.svg)

Mastery states: UNSEEN -> LEARNING -> PRACTICING -> MASTERED -> DECAYING

---

## 5. SM-2 Algorithm Detail

```
quality >= 3 (success):
  rep++, ef += 0.1 - (5-q)*(0.08 + (5-q)*0.02)
  interval: 0->1, 1->6, else round(prev * ef)

quality < 3 (failure):
  rep = 0, interval = 1

XP = base(10/5) * timeMult(0.5-1.5) * streakMult(1 + min(s,10)*0.1)
```

Quality from: correctness + response time ratio.

---

## 6. Grading Pipeline

![width:1000](https://raw.githubusercontent.com/chuma-beep/mathua/main/web/next-app/public/diagrams/system-design-grading-pipeline-dark.svg)

6+ grading strategies dispatched by grading type. SymPy for symbolic math.

---

## 7. Grading Strategies

| Grader | Types | Implementation |
|--------|-------|----------------|
| **Numeric** | int, float, fraction, mixed, sci-not | Go native, big.Rat, 1e-9 tolerance |
| **Multiple Choice** | multiple_choice | Case-insensitive, letter match |
| **Comparison** | comparison | > < = >= <= != |
| **Ordering** | ordering | Positional exact match |
| **Tuple** | tuple | Comma-separated, positional |
| **Complex** | complex | Polar form, delegates to SymPy |
| **Symbolic** | fallback | String normalization |

---

## 8. SymPy Subprocess

```
stdin:  {"id": "...", "expected": "x^2+2x+1", "answer": "(x+1)^2"}
               |
               v
    1. Preprocess (strip assignments, constants)
    2. Parse with sympy_parser (implicit multiplication)
    3. simplify(expected - answer) == 0
    4. expand, trigsimp, radsimp, powsimp, factor
    5. Numerical: 10 random points @ 1e-6 tolerance
               |
               v
stdout: {"id": "...", "correct": true, "feedback": ""}
```

10s timeout, 500 char limit. Falls back to Go symbolic grader.

---

## 9. API Layer (25 Routes)

| Group | Endpoints |
|-------|-----------|
| **Auth** | POST signup, login, GET me |
| **Session** | POST session, answer, reviews/* |
| **Progress** | GET progress, scores, weaknesses |
| **Diagnostic** | POST diagnostic, goal |
| **Data** | GET graph, leaderboard, lessons, courses, concepts |
| **Config** | GET config, health |

Middleware: CORS, logging, JWT auth, rate limiting (5/min on auth).

---

## 10. End-to-End Request Flow

![width:1000](https://raw.githubusercontent.com/chuma-beep/mathua/main/web/next-app/public/diagrams/system-design-request-flow-dark.svg)

12 steps from client submit to response. Target: <100ms numeric, <500ms SymPy.

---

## 11. Computerized Adaptive Testing

![width:1000](https://raw.githubusercontent.com/chuma-beep/mathua/main/web/next-app/public/diagrams/system-design-cat-diagnostic-dark.svg)

Binary search on topologically sorted DAG. ~20-35 questions vs 284.

---

## 12. DAG (Concept Graph)

Concepts defined in `data/concepts/*.json`:

```json
{
  "id": "algebra.quadratic-formula",
  "label": "Quadratic Formula",
  "prerequisites": ["algebra.quadratic-equations"],
  "grading_type": "numeric",
  "mastery_threshold": { "streak": 3 }
}
```

Validated: no cycles, no orphaned prereqs, no duplicates.
Topological sort via Kahn's algorithm.

---

## 13. Scheduler Priority

```
priority = 0.5 * daysSinceLastSeen
         + 0.2 * (1 - masteryScore)
         + 0.3 * weaknessScore
```

Decaying concepts (14+ days idle) get +5 bonus.

**70/30 balance**: 70% new material, 30% spaced repetition reviews.

---

## 14. Weakness Propagation

```
If concept.weakness > 0.3:
  For each dependent concept:
    dependent.weakness += concept.weakness * 0.3
```

Propagated at answer time. Creates a "weakness map" across the DAG.

---

## 15. Levels & Scoring

| Mastered | Title |
|----------|-------|
| 0 | Novice |
| 32 | Apprentice |
| 64 | Student |
| 96 | Scholar |
| 128 | Adept |
| 160 | Expert |
| 192 | Master |
| 224 | Grandmaster |
| 256 | Math Architect |

Weekly leaderboard: `weekly_mastered * 100`, resets Monday.

---

## 16. JWT Auth Flow

```
Signup:     POST /auth/signup -> bcrypt -> JWT (30d)
Login:      POST /auth/login  -> bcrypt compare -> JWT
Validate:   Authorization: Bearer <token> -> HS256 -> studentID
```

Rate-limited: 5 req/min per IP on auth endpoints.
Secret from `JWT_SECRET` env, or random per run.

---

## 17. Dual Database Strategy

```
Repository interface
  ├── SQLiteStore  (desktop/dev)
  │     PRAGMA journal_mode=WAL
  │     PRAGMA foreign_keys=ON
  │     PRAGMA busy_timeout=5000
  │
  └── PostgresStore (production web)
        Not yet implemented
```

Same schema, same interface. Transparent swap via `DATABASE_URL`.

---

## 18. Frontend (Next.js)

- **Static export** -- served by Go binary, no Node.js runtime
- **CSS custom properties** for dark/light theme
- **Warm palette**: `#0b0f1a` dark, `#fefcf4` light, `#c8a96e` gold accent
- **IBM Plex Serif** for text, **IBM Plex Mono** for code/labels
- **KaTeX** for math rendering
- **No global state** -- each page manages its own hooks

---

## 19. Key Trade-offs

| Decision | Chose | Instead of | Why |
|----------|-------|------------|-----|
| **Language** | Go | Python, Node | Single binary, fast startup |
| **Questions** | Generated | Static bank | Infinite variety, no cheating |
| **ORM** | Raw SQL | GORM, Prisma | Full query control, transparent |
| **Diagrams** | D2 + React Flow | Mermaid | Auto-layout, interactive |
| **Deploy** | Static export | SSR | One port, no Node runtime |

---

## 20. Summary

```
                    +------------------+
                    |   Web Browser    |
                    |  (Next.js/React) |
                    +--------+---------+
                             |
                    +--------+---------+
                    |   REST API       |
                    |  net/http        |
                    +--------+---------+
                             |
                    +--------+---------+
                    |   Core Engine    |
                    |  DAG + SM-2 + Gen|
                    +--------+---------+
                             |
                    +--------+---------+
                    |   Grading        |
                    |  Go + SymPy      |
                    +--------+---------+
                             |
                    +--------+---------+
                    |   SQLite/Postgres|
                    +------------------+
```

One language, two databases, two UIs, one engine.

---

## Questions?

mathua.app | github.com/chuma-beep/mathua
