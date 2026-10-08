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
│  Numeric · Choice · Comparison · Ordering ·    │
│  Tuple · Complex (pure Go)                     │
│  Symbolic · Polynomial · Expression → SymPy    │
│  (Python, shipped in the image)                │
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

### Grading — nine types, six Go graders + SymPy

The `grading_type` on each concept selects the grader. One `Router` dispatches nine enum values. Six are pure Go. Three route to SymPy:

- **Numeric** — integer, float, fraction, mixed number, scientific notation (pure Go).
- **Choice** — case-insensitive multiple-choice matching.
- **Comparison** — comparison operators (`>`, `<`, `=`, `>=`, `<=`).
- **Ordering** — ordered sequences separated by delimiters.
- **Tuple** — coordinate/n-tuple parsing.
- **Complex** — `a + bi` forms.

**Symbolic, Polynomial, Expression** — routed through a long-lived Python subprocess running `grading/sympy_service.py`. The Go client sends a JSON expression pair over stdin. SymPy parses both into expression trees. SymPy then tests equivalence (`simplify(expected - answer) == 0`, with expand/factor/trigsimp/radsimp/powsimp and numerical verification). Python and SymPy are **installed in the production image** and pinned in `grading/requirements.txt`.

The image build verifies the installation, so grading is available in production. There is **no pure-Go equivalence fallback** for these types. The router does not use the unused `symbolicGrader`.

### Storage — SQLite and Postgres

SQLite for local development (`mathua.db`, `WAL` `storage/migrate.go:23`). PostgreSQL for production web access via `DATABASE_URL=postgres://` (`internal/storage/postgres.go:7` `PostgresStore` via `pgx/v5/stdlib`, `pgSchema` `postgres.go:7`, `GREATEST` for XP floor, `STRING_AGG` for activity). `cmd/mathua` checks `DATABASE_URL` and opens `pgx` when `postgres://` is set, otherwise SQLite. The data schema is identical across both databases, abstracted behind a `Repository` interface (`storage/store.go:127`).

## Administration — role, not a second login

One authentication system serves everyone. A student signs in at `/login`. A Mathua
administrator signs in at the same place with the same form. There is no admin login page, no
admin password and no admin session. One column separates the two accounts.

### The role

`students.role` holds `student` or `admin`. The column is `NOT NULL DEFAULT 'student'`. Every
account created before the column therefore migrated to a learner.

That is the only safe direction. A database from before roles hold no evidence about who
administered it, because nothing in that database permitted administration.

The operator sets the role. Administrators set it too. The role never comes from XP, from
mastery, from level, from curriculum progress or from an email domain. A learner can change
each of those. A request body cannot set the role.

### Where the role is read

`requireAdmin` (`internal/server/adminrole.go`) is the boundary. It authenticates the request.
It identifies the account. Then it reads the current role from the database. It does not trust a
role carried in the token.

That is the decision the rest of the design follows from. A role embedded in a JWT is a claim
made at sign-in. The claim stays valid until the token expires. Expiry can be a day.

An administrator demoted at 09:00 keeps administrative access until midnight under that scheme.
That administrator can use the access to promote themselves back. A lookup per request costs one
indexed read. It removes the window entirely.

A demotion then takes effect on the next call, from any device, with no re-login.
`TestLastAdminDemotionTakesEffectImmediately` in `internal/server/adminrole_test.go` asserts this.

### Two trust models under `/admin`

Both live under `/admin`. They do not authorize each other.

| Path | Mechanism | Identifies an account |
|---|---|---|
| `/admin/reports` | Shared `ADMIN_PASSWORD` from the deployment environment | No |
| `/admin/users`, `/admin/audit` | Role on the account | Yes |

The triage password confers no account identity. It cannot promote anyone. It cannot be the
subject of a role check. `TestTriagePasswordDoesNotGrantRoleAdministration` asserts both
directions.

### The last administrator

The system never reaches zero administrators through an administrative action. The check sits
inside the UPDATE statement. It does not sit before the UPDATE.

```sql
UPDATE students SET role = 'student'
 WHERE id = ? AND role = 'admin'
   AND (SELECT COUNT(*) FROM students WHERE role = 'admin') > 1
```

A count read before a write loses a race. Two administrators who demote each other at the same
moment both read two. Both write. The system reaches zero through the feature meant to stop it.

One statement means the second attempt matches no rows.
`TestConcurrentDemotionsCannotReachZeroAdmins` runs that race twenty times. It uses a
file-backed SQLite database, because `:memory:` opens a separate database per connection and
cannot express the concurrency.

### Bootstrap

The operator creates the first administrator on the machine that holds the database:

```bash
mathua admin promote <email>
```

The command requires an account that already exists. It refuses to create one. Shell access to the machine
already grants control of the deployment, so it is the correct place to decide who administers.

The command calls the same `SetRoleAudited` as the HTTP handler. It passes `cli` as the actor. The command and the
web page therefore share one code path. That path writes the audit row, and it applies the last-administrator rule.
A second way to change a role is a hole in the audit log, because the log cannot record what it never saw.

An HTTP endpoint fails this test. Anyone who can reach the server could then ask for
administrator access.

This repository stores no administrator identity. It stores the mechanism. The deployment decides
who holds the role.

### The audit log

Every role change writes one row to `admin_audit`. The insert and the role change share one
transaction. If the insert fails, the role change rolls back. A privileged change therefore
cannot exist without a record. `TestFailedAuditInsertRollsBackPromotion` proves this. That test makes
the insert fail.

The log has no foreign key on `actor_id`. A trail deleted together with the account it records
is not a trail. Neither store offers an update path for the table. Neither store offers a delete
path. No HTTP route offers either.

The before and after snapshots come from `storage.AdminUser`. That type has no field for a
password hash, a share token or learner progress. A credential cannot reach the log through it.

### What the frontend is not

The Admin navigation entry appears only for an account whose stored role is `admin`. The entry
is a convenience. A hidden link is not authorization.

Every administrative endpoint rejects a non-administrator on its own. A learner who types
`/admin/users` reaches the page, and then the API refuses the request.

### Confirmation before a role change

`/admin/users` names the action on each row. The two actions are `Make administrator` and
`Remove administrator`. Each opens `ConfirmDialog`. The dialog names the account.

The dialog puts focus on `Cancel`. It does not put focus on the action that changes the role. A
stray Enter therefore does nothing. Escape, the backdrop, and `Cancel` all close the dialog and
send no request.

The page shows two errors apart. One is the list failing to load. The other is the server
refusing a change. They are different sentences. One earlier version shared a single slot. When
the page reloaded its list, it dropped the server's refusal and displayed nothing.

### The deployed surface

The static export writes both `/admin.html` and a directory named `/admin`. `nextStaticFS` in
`cmd/mathua/main.go` probes the requested path and serves a file when one exists. It must also
require that the path is not a directory. Without that test the server answers `/admin` with a
redirect to `/admin/`, and `/admin/` returns a generated directory listing. Every route answers
`200` in that state, so no page-level test can see the fault.
`cmd/mathua/staticfs_test.go` asserts that no response body holds a directory listing.

## Diagnostic — Computerised Adaptive Testing

Locates a student's knowledge frontier using a compressed covering set + info-gain CAT (`internal/diagnostic/cat.go:1`, `internal/diagnostic/report.go:11`). The engine builds a minimal covering set of the DAG. It repeatedly picks the concept that gives the maximal entropy reduction. It propagates `+0.3` evidence to prerequisites on correct and `-0.3` to dependents on incorrect, and tracks per-concept `KnowledgeConfidence 0–1`. The frontier is the highest belief drop. A supplemental diagnostic runs when confidence `<0.7`.

Assessment is 25–45 adaptive questions, not 657 exhaustive ones. It targets `80%` difficulty via `engine.computeDifficulty` (`internal/engine/engine.go:157`, `weakness→difficulty` 0.3–1.0).

## Full documentation

For interactive D2 diagrams, internal module detail, and formatted code blocks, see the **[Architecture docs](https://mathua.vercel.app/docs/architecture)** in the web app.
