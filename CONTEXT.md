# Mathua — Domain Language

> Derived from `improve.md` gap analysis + Q1-Q3 locks (Lesson/Study/Profile/Seam/effort XP, 50 XP quiz). Used by `grill-with-docs`, `tdd`, `domain-modeling` to name seams and tests.

## Language

**Concept** — Atomic node in the DAG (`internal/concepts/concepts.go:8` `id/label/domain/subdomain/prerequisites/mastery_threshold`). Example: `frac.add.word` (`data/concepts/fractions.json:149`, 25.0s).  
_Avoid_: topic (MA), lesson (maps to corpus, not node)

**DAG** — Directed acyclic graph of Concepts. Hard rules: prereqs mastered before surface (`DAG.Available` in `internal/concepts/concepts.go` `Available`), Kahn's topo order, `scripts/validate_graph.go` rejects cycles/orphans.  
_Avoid_: graph alone, knowledge map

**Lesson (corpus / worked example)** — Inline worked example per KP, 10× granularity vs textbook (`improve.md:39` 300 topics ×3 KPs). Q1 locked: Lesson is corpus, not library. Rendered via `KatexContent` (`web/next-app/components/LearnStepper.tsx`) before practice, one knowledge point at a time — never the whole article (ADR-045, ADR-047). Source: `data/lessons/lessons.json:1` `teaching/*.md` + `algebrica/*.md` (`README.md:177`).  
_Avoid_: Study, chapter

**Reference layer** — Was *Study*: a browseable corpus library holding `Lesson` alongside `LessonQuiz`, which did not gate progression. ADR-047 closes that as a learner workflow, so there is no longer a surface for this term. The material is now `components/ConceptReference.tsx` — worked examples, report menus, prerequisites and the article fallback for one concept, rendered inside `/learn` with no link out of it. Two names survive and mean different things: **`/api/study/answer`** is Learn's answer seam (a name, not a destination), and **`/study`** is a 460 B forwarder to `/learn` that exists only for old bookmarks.  
_Avoid_: Practice, Learning Mode (removed `web/next-app/app/session/page.tsx:1` chooser)

**Profile (hub)** — Universal landing for guest + authed (`web/next-app/app/profile/page.tsx:44`). Guest gets ephemeral `mathua_guest_id` (`web/next-app/lib/auth.ts:4` `ensureGuestId`). The hub **and** the learner's own report: `ProfileStats` identity, `PositionBlock` (*where you are* — mastered / unlocked / locked plus the frontier the scheduler would serve next), the `NextUpSummary` task, the alerts (paused / reviews due / quiz gate / diagnostic), the new-user block, then the report — `Activity` (`ActivityHeatmap`) → `By domain` (`DomainProgress`) → `How you're doing` (the learner's own first/second-pass and attempts) → `StrugglesSection` → recently unlocked. Cohort figures (`total_students`, `retention_rate`, weekly first-pass) stay on `/docs/efficacy`, never here. Post-auth always → `/profile` (`login:165`, `onboard:186`, `goals:224`).  
_Avoid_: dashboard alone, account page, a separate progress page (there isn't one — ADR-025)

**History** — Every question the learner has answered (`web/next-app/app/history/page.tsx:1`), rendered by `AttemptList` (shared with the public token-based share report at `/share`). Defaults to the `Missed` filter; filterable by source; paginated. Was `/progress-card`, and the sidebar entry was called `Progress` while resolving here — a label pointing at a mistake log. Both are fixed: the rail reads `History` and the report lives on `/profile`.  
_Avoid_: progress, progress card, transcript (fine in prose, not in a route)

**Diagnostic (frontier)** — Adaptive CAT finding knowledge frontier (`internal/diagnostic/cat.go:1`, `improve.md:27` 30–45 min, covering set + info-gain). Entry: `/onboard` welcome (`onboard/page.tsx:56`), CTA on `Profile` (`profile:209`) labeled `Recommended` vs `Retake` (`profile:209`).  
_Avoid_: test, assessment, quiz (reserved)

**Seam** — Public boundary to test at (the `tdd` skill's term: the narrowest interface that makes the behavior observable). Pre-agreed seams for next bite: `Learn → SubmitStudyAnswer` (`POST /api/study/answer` `server.go:484` → `submitAnswerWithTask` `engine.go:1255`).

**TaskLesson / TaskReview / TaskMultistep / TaskQuiz** — XP task types (`internal/xp/policy.go` effort-derived bases). Q2 locked: `*.word` and plain share the effort calibration (thresholds already differ). Reviews and quizzes earn the same effort base as lessons; the task type does not change the number.
_Avoid_: flat per-question XP, word as separate prereq depth (temporary node, future `Variants` fold via `enrichment.json` `improve.md:114`)

**Quiz (50 XP gate)** — Q3 was MA-verbatim `150 XP` (`improve.md:27`); rescaled to `50 XP` with the small-awards economy (same ~5-lesson cadence, smaller numbers — user-tuned economy supersedes MA fidelity, see ADR-020). Timed closed-book, `80%` difficulty via `computeDifficulty`, diverse recents, immediate remedial + retake (`internal/quiz/quiz.go:1`). Daily goal `10` (drift fix shipped: versioned backfill in `sqlite.go`). Banner: `Study → Take Test` + `Profile` due-reviews.
_Avoid_: 1–2 lesson gate, 30 XP shortcut

**XP** — `Award` in `internal/xp/policy.go` (`EffortBase` clamped to 1-5, × `PerformanceMultiplier` 1.1 inside the time threshold and 1.0 outside; no time or streak multiplier): an effort-derived approximation, never exact time accounting. Miss 0, rush −5, diagnostic 0 XP. Daily goal `30` (`xp.DefaultDailyGoal`). Verified via `go vet` + `learn → answer → profile` `playwright`.

**GUEST_KEY** — `mathua_guest_id` persisted so guest `/profile` survives reload (`lib/auth.ts:4`).

**Landing graph (hero)** — WebGL concept map on `/` (`components/MathConceptGraph3D.tsx`), fed by build-time artifacts `public/data/concepts.graph.json` (nodes + indexed edges + domain counts) and `public/data/concepts.positions.json` (quantized `[x,y,z]`), generated by `scripts/build-graph.mjs` (`npm run graph:build`, also via `prebuild`). The corpus `data/concepts.json` stays out of First Load JS; only `data/graph.meta.json` (counts) is bundled. LOD tiers (far/mid/close) derive from world-units-per-pixel; labels render in one imperative DOM overlay; focus mode dims non-1-hop nodes/edges but never removes them.  
_Avoid_: bundling `concepts.json` into the landing bundle, per-node DOM/event listeners, client-side force layout

## Relationships

- DAG holds many Concepts; Concept has Prerequisites; `Available` filters by mastered
- Profile holds the Next task, the alerts, and the report (position/domains/activity/weakness); History holds attempted Questions; Study holds Lessons; Lesson holds Concepts + Problems
- Diagnostic compresses DAG → picks Concept via info-gain; Post-diagnostic → Profile; Profile holds the report, Profile → History for the attempts
- Study Lesson `→` TaskLesson (or TaskMultistep if `.word`) `→` SubmitAnswer `→` XP/Progress `→` Profile; at 50 XP `→` Quiz (TaskQuiz)

## Flagged ambiguities

- `Learning Mode` infinite loop removed (`session` now chooser `session/page.tsx:1` Diagnostic vs Study, then `Profile` hub); do not reintroduce as primary
- `Study` links that were wrong post-login (`login Skip → /study` etc) fixed to `→ /profile`; do not regress
- `skills/CONTEXT.md` is marketplace language (Issue tracker) — not Mathua; Mathua context lives here
- Stale-session submit after a wipe (follow-up, not an approved change): answering with a pre-wipe session id can fail with HTTP 500, because the referenced `sessions` row was deleted and `RecordAttempt` hits the foreign key. This fails closed — nothing is resurrected — but 500 is not a deliberate contract. Proposed contract, pending a separate decision: return 409 and require the client to re-serve or restart practice. Do not implement it silently inside another phase.

## Decisions (ADRs)

The 54 architecture decision records live in [`docs/adr/`](docs/adr/README.md), one file each,
indexed by number. They were inline in this file until it reached 140 KB across 124 lines with
no table of contents; the numbers are unchanged, so a citation to "ADR-048" in a review or a
commit message still resolves.
