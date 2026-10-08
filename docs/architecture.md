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

## Administration — roles, not a second login

One authentication system serves everyone. A learner signs in at `/login`. An administrator
signs in at the same place with the same form. There is no admin login page and no admin
password. A role on the account separates the two.

### The roles

`students.role` holds one of four values. The table lists them from most to least privilege.

| Role | What it can do |
|---|---|
| `owner` | Everything. The only role that can grant another owner. |
| `admin` | Run administration and moderation. Cannot grant owner. |
| `moderator` | Work the moderation queue and inspect content. |
| `student` | Learn. No administrative access. |

The column is `NOT NULL DEFAULT 'student'`. Every account created before roles migrated to a
learner. A database from before roles holds no evidence about who administered it, because
nothing in that database permitted administration.

V1 held only `student` and `admin`. The migration moves each V1 `admin` to `owner`. V1 `admin`
was full control. The move keeps that control. Leaving the value as `admin` takes away the power to
manage staff.

The role never comes from XP, from mastery, from level or from an email domain. A learner can
change each of those. A request body cannot set the role.

### Where the role is read

`requirePermission` (`internal/server/adminrole.go`) is the boundary. It authenticates the
request. It identifies the account. Then it reads the current role from the database. It does
not trust a role carried in the token.

That is the decision the rest of the design follows from. A role embedded in a JWT is a claim
made at sign-in. The claim stays valid until the token expires. Expiry can be a day. An
administrator demoted at 09:00 keeps access until midnight under that scheme. That administrator
can use the access to promote themselves back.

One indexed read per request removes the window. A demotion then takes effect on the next call,
from any device, with no re-login.
`TestLastAdminDemotionTakesEffectImmediately` in `internal/server/adminrole_test.go` asserts
this.

### The permission model

`internal/admin/admin.go` holds the model. Each role maps to a set of permissions. A route names
one permission. The boundary checks the role against that permission.

The model is data, not scattered checks. A permission matrix spread across handlers is a matrix
nobody can review. The failure it produces is a route that forgot a check.

The matrix is small on purpose. `owner` holds every permission. `admin` holds all but
`roles.grant_owner`. `moderator` holds the moderation and content permissions. A learner holds
nothing.

`CanAssign` decides one role change. It applies three rules. An actor cannot grant a role above
their own. An actor cannot act on a peer or a superior, unless the actor is an owner. Granting
`owner` needs `roles.grant_owner`. `TestCanAssign` in `internal/admin/admin_test.go` states the
rules.

### Moderation under the same boundary

Reported content is moderated at `/api/admin/reports`. The role boundary guards it. A moderator
resolves a decision and the same transaction writes the audit row.

V1 gated content triage with a shared `ADMIN_PASSWORD`. V2 removed that mechanism. A shared
password cannot say which person used it. One moderation path replaces the two.

A report is never deleted. A decision sets the status, the reason, the actor and the time. The
row keeps the complaint, so the next moderator can read what was decided before.

### Staff access by invitation

An owner or admin adds a contributor from `/admin/contributors`. The server creates an
invitation. The invitation carries a token, an email and a role.

The server stores only a SHA-256 hash of the token. A database read therefore yields no usable
credential. The raw token returns once, in the response, and nowhere else.

The invitation is single-use, time-limited and bound to the invited address. Acceptance burns
the token. Acceptance checks the signed-in account against the invited email. A revoked,
expired, reused or mismatched invitation grants nothing. Each failure reads the same to the
caller, so a prober learns nothing about which near-miss they reached.

`ConsumeAdminInvitation` burns the token, grants the role and writes the audit row in one
transaction. No state exists in which someone holds a staff role with no record of how they got
it.

### The last owner

The system never reaches zero owners through an administrative action. The check sits inside the
UPDATE statement. It does not sit before the UPDATE.

```sql
UPDATE students SET role = ?
 WHERE id = ? AND role = 'owner'
   AND (SELECT COUNT(*) FROM students WHERE role = 'owner') > 1
```

A count read before a write loses a race. Two owners who demote each other at the same moment
both read two. Both write. The system reaches zero through the feature meant to stop it. One
statement means the second attempt matches no rows.

The invariant is on the owner, not the admin. An admin or a moderator can be removed freely. An
owner cannot be, because nothing but an owner can create another owner.

### Bootstrap

The operator creates the first owner on the machine that holds the database:

```bash
mathua admin promote <email>
```

The command requires an account that already exists. It refuses to create one. Shell access to
the machine already grants control of the deployment, so the machine is the correct place to
decide who administers.

The command grants `owner`, not `admin`. The first promotion must be able to create the rest.
The command calls the same `SetRoleAudited` as the HTTP handler, with `cli` as the actor. The
command and the web page therefore share one code path that writes the audit row.

An HTTP endpoint fails this test. Anyone who can reach the server could then ask for
administrator access. This repository stores no administrator identity. It stores the mechanism.

### The audit log

Every role change and every moderation decision writes one row to `admin_audit`. The insert and
the change share one transaction. A failed insert rolls the change back. A privileged change
therefore cannot exist without a record.
`TestFailedAuditInsertRollsBackPromotion` proves this.

The log has no foreign key on `actor_id`. A trail deleted together with the account it records
is not a trail. Neither store offers an update path. Neither store offers a delete path. No HTTP
route offers either. Each event names the actor, the action, the entity, the time, and the state
either side of the change.

`ListAdminAuditFor` reads the events about one entity. A role history and a report history are
queries over the trail, not columns anyone maintains by hand.

The before and after snapshots come from `storage.AdminUser`. That type has no field for a
password hash, a share token or learner progress. A credential cannot reach the log through it.

### What the frontend is not

`AdminNav` shows a section only when the caller holds the permission for it. The entry is a
convenience. A hidden link is not authorization.

Every administrative endpoint rejects a caller without the permission on its own. A learner who
types `/admin/users` reaches the page, and then the API refuses the request with a 403.

### Confirmation before a role change

`/admin/users` offers a role picker on each row. The picker opens `ConfirmDialog`. The dialog
names the account and the new role.

The dialog puts focus on `Cancel`. A stray Enter therefore does nothing. Escape, the backdrop
and `Cancel` all close the dialog and send no request.

The page shows two errors apart. One is the list failing to load. The other is the server
refusing a change. They are different sentences.

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
