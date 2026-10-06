# Mastery assessments: decision map

Status: **design. No migration written, no unit authored, no threshold chosen.**

This is the decision map for **ADR-046** (`CONTEXT.md`). ADR-046 settles the *shape* — a
mastery unit is a named partition of the existing concept DAG, an assessment contributes
evidence rather than writing mastery, two unit states are stored and five are derived, and
`Engine.ProgressionGate` is the single server-side answer to "can this learner progress".

It deliberately does **not** settle four things. They are below, as decisions with their
blocking edges, because each one is a product judgement rather than a technical fact and each
one changes the schema if guessed at.

Read this before writing any migration. Two of these (D1, D2) determine the table shape.

---

## Why this is a map and not a spec

The expensive part is not the code. It is:

- **D3** — authoring a partition over 657 concepts. Nothing in the corpus expresses block
  boundaries today: `data/concepts/*.json` carries only `id`, `label`, `domain`,
  `subdomain`, `prerequisites`, `grading_type`, `mastery_threshold`. `data/courses.json` has
  22 entries of 4 `targets` each — loose bundles, not ordered blocks — and
  `data/concepts/enrichment.json` is 10 heuristic overrides of which only 2 carry
  `encompasses`. So the partition is authored data, and it is the long pole by an order of
  magnitude over the Go work.
- **D4** — the pass threshold, which is a content judgement about what "can demonstrate this
  independently" means for this corpus, and which no test can tell us.

Everything downstream of those is mechanical. Sequencing them first would mean reworking the
table when D3 says a unit is a term's worth of concepts rather than a week's.

---

## Decisions

### D1 — What is stored for an assessment attempt?

**Blocks:** D5, and the table migration.

The spec says the result must be persisted and must not be computed by the frontend. It does
not say what a result *is*.

Open: is an attempt a single row (latest wins) or an append-only log of every attempt? The
retake policy (D5) needs the log; the dashboard needs only the latest.

- **Option A — one row per (unit, student), overwritten.** Simplest; retake policy has to
  infer history from a counter.
- **Option B — append-only `mastery_unit_attempts`.** Every attempt is a row; the current state
  is a query. Retake policy, analytics, and "attempts" in the spec all fall out.
- **Option C — a row per attempt *plus* a materialised current-state row.** B's history with
  A's cheap read, at the cost of two things that can disagree.

Recommended: **B**. The cost of an extra table is one insert; the cost of losing attempt
history is being unable to enforce or explain a retake rule. Note that `quiz.Session.Attempts`
already exists per session (`internal/quiz/quiz.go:19`) but is deleted on completion
(`server.go:2887`) and lives in memory, so it is not a history today.

Carry into the ADR: attempts are an event log; "is an assessment currently required" is derived
from it, not stored as a boolean that events must keep in sync.

### D2 — Is `assessment_required` stored, or derived?

**Blocks:** D5, the table migration, `ProgressionGate`.

ADR-046 stores `assessment_required` and `remediation_required` because they cannot be computed
from `concept_progress`. That is true for remediation (it needs to know an assessment *failed*,
which is not in `concept_progress`). It is less obviously true for `assessment_required`, which
may be derivable: "the active unit has all concepts learned and its assessment has not passed
since".

- **Option A — store both.** Simple queries; two state machines to keep consistent with the
  ladder.
- **Option B — store only `remediation_required`**, derive `assessment_required` from unit
  membership plus ladder state plus attempt history.
- **Option C — store neither**; both are queries over the attempt log.

Recommended: **B**, pending D1, because deriving the "required" case removes a state that must
agree with the ladder — and every disagreement between a stored status and the ladder is the
ADR-037 bug class again. If D1 lands on B (append-only log), C becomes cheap enough to prefer.

Carry into the ADR: the count of stored unit states, currently two.

### D3 — How are mastery units authored?

**Blocks:** D4 (you cannot set a pass threshold for an undefined unit), and everything.

Open: one unit per what? Options, in order of how much judgement each demands:

- **Option A — per subdomain.** `subdomain` already exists on every concept and is curated.
  Mechanical to generate, and almost certainly too coarse and too fine in places; also two
  adjacent subdomains may not share a prerequisite boundary, which the derivation rule in
  ADR-046 forbids.
- **Option B — hand-authored per curriculum block**, ~15–30 units, chosen so each is a coherent
  assessable group with a clean prerequisite boundary. This is what the spec describes and it
  is the only option that satisfies "do not invent arbitrary boundaries merely to satisfy the
  requirement".
- **Option C — per course** (`data/courses.json`). Reuses existing grouping, but a course is a
  grade-level target list of 4 concepts and cuts across prerequisite boundaries.

Recommended: **B**, generated from A as a *draft* and corrected by hand — with `validate_graph`
rejecting a unit whose members' DAG closure disagrees with its neighbours', which is the check
that forces the boundary to be real.

Open sub-question: **a unit of one concept**. The partition rules permit it; it should be a
content error. Decide whether that is a hard validator failure or a reported count, because a
hard failure makes authoring 657 concepts a one-pass job.

### D4 — What is the pass threshold?

**Blocks:** D5, and it is unanswerable before D3.

Open: PASS is currently undefined — there is no score at all, only `IsComplete`. Options:

- **Option A — proportion correct on the unit's assessment.** Needs a number.
- **Option B — the evidence model.** A pass means the assessment's own answers cleared
  `mastery.MasteryThreshold` (0.6) for every concept tested, rather than a single number. This
  is the most consistent with ADR-038, and it means "pass" is defined by the same model that
  will judge mastery — which is either elegant or circular, and that is the thing to decide.
- **Option C — per-unit thresholds.** Most faithful to "can demonstrate this independently",
  least maintainable.

Also unresolved, and related: **the assessment-readiness heuristic** — when a unit becomes
`assessment_ready`. ADR-046 permits the existing `xp.QuizGateXP` (50) as a starting proxy but
requires it stay separate from the boundary, because ADR-046's whole premise is that XP is
evidence of effort rather than of competence.

### D5 — What clears a failed assessment?

**Blocks:** the retake policy, `ProgressionGate`, and the remediation UI.

Spec requires `FAIL → remediation → sufficient new evidence → available again`, and forbids
immediate retakes. Open: what counts as "sufficient new evidence"?

- **Option A — XP since the failure** past some fraction of `QuizGateXP`. Cheap, and it is the
  metric that triggered the requirement in the first place.
- **Option B — attempts on the failed concepts** past a floor, using `mastery.MasteryEvidenceFloor`
  (6). Consistent with how mastery is measured; a learner can pass the floor on concepts that
  are not the ones they failed.
- **Option C — the failed concepts reach `mastery` again.** Strictest, and the only version that
  cannot be satisfied by unrelated work. Also the most frustrating, and it interacts with
  ADR-037: failing does not un-learn a concept, so "reach mastery again" may already hold at
  the moment of failure.

Option C has a real problem worth naming: `TestMasteryIsNeverLostToBadEvidence` means a failed
assessment does **not** demote a concept, so "recover mastery on the failed concepts" can be
satisfied before the learner does anything. C needs a *since-failure* qualifier, which makes it
A or B with extra steps.

Recommended: **B**, per failed concept, with the count stored on the attempt so it is not
re-derived from a moving window.

### D6 — Where is the boundary enforced?

**Blocks:** nothing, but it determines the test surface.

ADR-046 names `Engine.ProgressionGate` as the single implementation and says it is called by
`RecommendNext`, `NextQuestion` and the practice endpoint. Open: what does "blocked" refuse?

The spec is explicit that it must **not** be a lock on the app — Study, Graph, History,
Settings, and remediation all stay reachable. The remaining question is whether `/study`'s
embedded practice and `/learn?concept=<a concept beyond the boundary>` are refusals or
warnings, given ADR-021 made Study a reference-only surface and ADR-017 made `/learn` a soft
dismissible banner rather than an auto-redirect.

---

## Invariants and where each is tested

Per the spec's §18. ADR-046 places these; this is the index.

| # | Invariant | Test |
|---|---|---|
| 1 | A required assessment cannot have an unrelated new concept as `primary` | scheduler unit test, constructed state |
| 2 | A refresh does not clear it | `ProgressionGate` against a **fresh engine instance** |
| 3 | It survives navigation | same, plus no client state in the path |
| 4 | Passing clears the requirement | table transition test |
| 5 | Failing does not clear it | table transition test |
| 6 | Failing yields remediation, not immediate retake | table transition test + `ProgressionGate` |
| 7 | The server decides, not the frontend | assert the three call sites route through the gate |
| 8 | A constructed URL cannot bypass it | HTTP-level test on `/api/lessons/{id}/practice` — a client-side test proves nothing |
| 9 | Study/reference stays reachable | assert `ProgressionGate` is not called from Study/Graph/History/Settings |
| 10 | Remediation is reachable while blocked | practice endpoint serves the failed concepts |

Invariants 2 and 3 are the pair most likely to pass vacuously: a requirement held in client
state survives both tests. Assert against a new engine instance so the only thing carrying it is
the stored row.

---

## Order of work

```
D3 (units) ──→ D4 (pass + readiness)
                    │
                    ├─→ D1 (attempt log) ──→ D2 (stored states) ──→ D5 (clear-on-fail)
                    │                                                      │
                    └──────────────────────────────────────────────────────┘
                                          ↓
                          concepts.Build partition checks (D3 data lands)
                                          ↓
                          mastery_unit_assempts table, both stores + ADR-012 parity
                                          ↓
                          Engine.ProgressionGate ──→ RecommendNext / NextQuestion / practice
                                          ↓
                          Recommendation gains Status + Reason: mastery_required
                                          ↓
                          dashboard: Required / Practice required / Recommended check
```

Content first, because D4 is unanswerable without it and every migration is rework if it moves.

---

## Not decided here

- **Unit count and size.** D3 decides; this document does not guess.
- **Whether the diagnostic placement interacts.** ADR-024 deferred persisting placement so
  `FrontierLabel` survives the session; a unit partition is a plausible place to hang that, and
  it should not be adopted by accident.
- **Retention across the assessment boundary.** Still ADR-037's open `tier_attained_at`. An
  assessment is a spacing event, so this feature makes that gap more visible, not smaller.
- **Whether assessment evidence is weighted differently from practice evidence.** ADR-038's
  `BuildEvidence` has no notion of question provenance beyond difficulty; a passing assessment
  currently contributes exactly what a guided practice answer contributes, plus the `TaskQuiz`
  premium. If an assessment should count for more, that is a change to the evidence model and
  belongs in its own ADR.