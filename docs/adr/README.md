# Architecture Decision Records

One file per decision. The domain language stays in [`CONTEXT.md`](../../CONTEXT.md).

| ADR | Decision | Date |
|-----|----------|------|
| [ADR-001](ADR-001-q1.md) | Q1 |  |
| [ADR-002](ADR-002-q2.md) | Q2 |  |
| [ADR-003](ADR-003-q3.md) | Q3 |  |
| [ADR-004](ADR-004-landing-graph.md) | Landing graph |  |
| [ADR-005](ADR-005-grading-contract.md) | Grading contract | 2026-09-27 |
| [ADR-006](ADR-006-grading-never-penalizes-an-infrastructure-fault.md) | Grading never penalizes an infrastructure fault | 2026-09-27 |
| [ADR-007](ADR-007-study-anchor-is-per-question.md) | Study anchor is per-question | 2026-09-27 |
| [ADR-008](ADR-008-placement-never-destroys-mastery.md) | Placement never destroys mastery | 2026-09-27 |
| [ADR-009](ADR-009-every-served-question-carries-grading-type-every.md) | Every served `Question` carries `grading_type`; every answer host shows the matching input (study, review, quiz, diagnostic). `Engine.GradeAnswer` is the one grader; `Result.Unavailable` is fail-soft everywhere (no attempt/streak/XP; client retries). `AvgResponseTime` includes incorrect attempts. `*.word` is `TaskMultistep 15` in both practice and study (implements ADR-002). Deferred |  |
| [ADR-010](ADR-010-grading-data-repair.md) | Grading data repair | 2026-09-27 |
| [ADR-011](ADR-011-structured-choice-input.md) | Structured choice input | 2026-09-27 |
| [ADR-012](ADR-012-postgres-is-a-first-class-store.md) | Postgres is a first-class store | 2026-09-27 |
| [ADR-013](ADR-013-corpus-counts-are-generated-not-hand-written.md) | Corpus counts are generated, not hand-written | 2026-09-27 |
| [ADR-014](ADR-014-one-diagnostic-route.md) | One diagnostic route | 2026-09-27 |
| [ADR-015](ADR-015-stub-lessons-fail-the-gate.md) | Stub lessons fail the gate | 2026-09-28 |
| [ADR-016](ADR-016-currency-vs-math-in-the-canonicalizer.md) | Currency-vs-math in the canonicalizer | 2026-09-28 |
| [ADR-017](ADR-017-adaptive-lesson-loop.md) | Adaptive lesson loop | 2026-09-29 |
| [ADR-018](ADR-018-attempt-first-feed-task-shelf.md) | Attempt-first feed + task shelf | 2026-09-29 |
| [ADR-019](ADR-019-destinations-honest-estimates.md) | Destinations + honest estimates | 2026-09-29 |
| [ADR-020](ADR-020-small-awards-economy.md) | Small-awards economy | 2026-09-29 |
| [ADR-021](ADR-021-study-is-reference-only-detail.md) | Study is reference-only detail | 2026-09-29 |
| [ADR-022](ADR-022-reset-clears-evidence-keeps-choices.md) | Reset clears evidence, keeps choices | 2026-09-29 |
| [ADR-023](ADR-023-single-next-head-over-a-ranked-shelf.md) | Single Next head over a ranked shelf | 2026-09-30 |
| [ADR-024](ADR-024-progress-moved-off-the-hub-2026-10.md) | Progress moved off the hub (2026-10-01), **partly superseded by ADR-025** |  |
| [ADR-025](ADR-025-the-report-lives-on-the-hub.md) | The report lives on the hub | 2026-10-01 |
| [ADR-026](ADR-026-a-wipe-must-be-unwritable-and-must.md) | A wipe must be unwritable, and must not lie about why (2026-10-03), amends ADR-022 |  |
| [ADR-027](ADR-027-the-self-check-is-a-simplification-not.md) | The self-check is a simplification, not an equivalence check | 2026-10-03 |
| [ADR-028](ADR-028-a-corpus-checker-must-state-what-it.md) | A corpus checker must state what it did not check | 2026-10-03 |
| [ADR-029](ADR-029-a-keycap-s-latex-is-an-insert.md) | A keycap's `latex` is an insert, not a command | 2026-10-03 |
| [ADR-030](ADR-030-the-keypad-may-only-be-raised-never.md) | The keypad may only be raised, never lowered | 2026-10-03 |
| [ADR-031](ADR-031-the-loading-state-is-the-editor-that.md) | The loading state is the editor that always works, not a spacer | 2026-10-03 |
| [ADR-032](ADR-032-verified-0-111-0-configuration-and-one.md) | Verified 0.111.0 configuration, and one doc that would have shipped broken | 2026-10-03 |
| [ADR-033](ADR-033-a-confirmation-phrase-is-compared-case-insensitively.md) | A confirmation phrase is compared case-insensitively, because phones capitalise it | 2026-10-03 |
| [ADR-034](ADR-034-mathlive-is-universal-and-the-keypad-answers.md) | MathLive is universal, and the keypad answers for the corpus rather than for arithmetic | 2026-10-03 |
| [ADR-035](ADR-035-the-spacebar-was-never-broken-it-was.md) | The spacebar was never broken, it was empty | 2026-10-04 |
| [ADR-036](ADR-036-one-box-square-corners-and-no-border.md) | One box, square corners, and no border outside the palette (2026-10-04), three independent commits, because the report bundled three unrelated defects |  |
| [ADR-037](ADR-037-mastery-is-an-exit-and-decay-is.md) | Mastery is an exit, and decay is not a return to class | 2026-10-04 |
| [ADR-038](ADR-038-mastery-is-decided-by-evidence-and-the.md) | Mastery is decided by evidence, and the ladder is three different questions | 2026-10-04 |
| [ADR-039](ADR-039-difficulty-belongs-to-the-question-not-to.md) | Difficulty belongs to the question, not to the host that graded it | 2026-10-05 |
| [ADR-040](ADR-040-the-client-renders-the-authoritative-score-instead.md) | The client renders the authoritative score instead of computing its own | 2026-10-05 |
| [ADR-041](ADR-041-a-finished-answer-with-nothing-after-it.md) | A finished answer with nothing after it is a bug, and it was not the learner's fault | 2026-10-05 |
| [ADR-042](ADR-042-the-engine-decides-what-comes-next-and.md) | The engine decides what comes next, and the browser is told | 2026-10-06 |
| [ADR-043](ADR-043-a-learner-is-told-where-they-are.md) | A learner is told where they are against what they can reach, not against the catalogue | 2026-10-06 |
| [ADR-044](ADR-044-learn-names-its-own-units-and-withholds.md) | `/learn` names its own units and withholds its own syllabus | 2026-10-06 |
| [ADR-045](ADR-045-learn-teaches-a-slice-and-never-the.md) | `/learn` teaches a slice, and never the article it was cut from | 2026-10-06 |
| [ADR-046](ADR-046-a-mastery-assessment-verifies-evidence-it-does.md) | A mastery assessment verifies evidence; it does not write mastery | 2026-10-06 |
| [ADR-047](ADR-047-study-is-closed-as-a-learner-workflow.md) | Study is closed as a learner workflow, and its content became a reference panel inside `/learn` | 2026-10-08 |
| [ADR-048](ADR-048-administration-is-a-role-on-the-account.md) | Administration is a role on the account, and it is re-read on every request | 2026-10-08 |
| [ADR-049](ADR-049-admin-v2-turns-administration-from-a-localhost.md) | Admin V2 turns administration from a localhost tool into a deployed, permissioned one | 2026-10-08 |
| [ADR-050](ADR-050-a-mastery-unit-is-a-named-partition.md) | A mastery unit is a named partition of the existing DAG, and a pass adds evidence without setting mastery | 2026-10-07 |
| [ADR-051](ADR-051-the-diagnostic-places-a-learner-it-is.md) | The diagnostic places a learner; it is never an active task, and placement is reachability without mastery | 2026-10-08 |
| [ADR-052](ADR-052-administration-reaches-a-phone-and-the-mobile.md) | Administration reaches a phone, and the mobile menu was dead while it said otherwise | 2026-10-08 |
| [ADR-053](ADR-053-the-admin-entry-was-missing-on-phones.md) | The Admin entry was missing on phones because the login path dropped the role, and no test walked that path | 2026-10-08 |
| [ADR-054](ADR-054-reference-access-is-not-learning-eligibility-and.md) | Reference access is not learning eligibility, and `/study` is restored as a context-bound reference page | 2026-10-11 |
| [ADR-055](ADR-055-ai-augmentation-is-a-registry-of-declared-policies.md) | AI augmentation is a registry of declared policies, not a trusted provider list | 2026-10-16 |
| [ADR-056](ADR-056-timed-quizzes-assess-fluency-and-spaced-review.md) | A timed quiz assesses fluency; spaced review owns the FSRS schedule | 2026-10-28 |
| [ADR-057](ADR-057-calculus-is-three-courses-and-three-domains.md) | Calculus is three courses, so it is three domains | 2026-10-30 |
