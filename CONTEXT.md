# Mathua — Domain Language

> Derived from `improve.md` gap analysis + Q1-Q3 locks (Lesson/Study/Profile/Seam/XP 10/5/15/20, 150 XP quiz). Used by `grill-with-docs`, `tdd`, `domain-modeling` to name seams and tests.

## Language

**Concept** — Atomic node in the DAG (`internal/concepts/concepts.go:8` `id/label/domain/subdomain/prerequisites/mastery_threshold`). Example: `frac.add.word` (`data/concepts/fractions.json:149`, 25.0s).  
_Avoid_: topic (MA), lesson (maps to corpus, not node)

**DAG** — Directed acyclic graph of Concepts. Hard rules: prereqs mastered before surface (`internal/concepts/concepts.go:73` `Available`), Kahn's topo order, `scripts/validate_graph.go` rejects cycles/orphans.  
_Avoid_: graph alone, knowledge map

**Lesson (corpus / worked example)** — Inline worked example per KP, 10× granularity vs textbook (`improve.md:39` 300 topics ×3 KPs). Q1 locked: Lesson is corpus, not library. Rendered via `KatexContent` (`web/next-app/app/study/page.tsx:592`) before practice. Source: `data/lessons/lessons.json:1` `teaching/*.md` + `algebrica/*.md` (`README.md:177`).  
_Avoid_: Study, chapter

**Study (library)** — Browseable corpus library (`web/next-app/app/study/page.tsx:1` `LessonDetail/DomainOverview`). Q1 locked: Study stays library label, not landing. Holds `Lesson` + `LessonQuiz` (`study/page.tsx:596`) but does not gate progression.  
_Avoid_: Practice, Learning Mode (removed `web/next-app/app/session/page.tsx:1` chooser)

**Profile (hub)** — Universal landing for guest + authed (`web/next-app/app/profile/page.tsx:44`). Guest gets ephemeral `mathua_guest_id` (`web/next-app/lib/auth.ts:4` `ensureGuestId`). Shows `ProfileStats`, `DomainProgress`, `StrugglesSection`, `ActivityHeatmap`, and `Diagnostic CTA`. Post-auth always → `/profile` (`login:165`, `onboard:186`, `goals:224`).  
_Avoid_: dashboard alone, account page

**Diagnostic (frontier)** — Adaptive CAT finding knowledge frontier (`internal/diagnostic/cat.go:1`, `improve.md:27` 30–45 min, covering set + info-gain). Entry: `/onboard` welcome (`onboard/page.tsx:56`), CTA on `Profile` (`profile:209`) labelled `Recommended` vs `Retake` (`profile:209`).  
_Avoid_: test, assessment, quiz (reserved)

**Seam** — Public boundary to test at (`tdd` `skills/skills/engineering/tdd/SKILL.md:1`). Pre-agreed seams for next bite: `LessonQuiz → SubmitAnswer` (`internal/engine/engine.go:398` `POST /api/answer` `server.go:286`).

**TaskLesson / TaskReview / TaskMultistep / TaskQuiz** — XP task types (`internal/engine/engine.go:580` `taskBaseXP` `10/5/15/20`). Q2 locked: `*.word` (`frac.add.word`, `arith.*.word` `internal/generator/fractions/generators.go:272`) → `TaskMultistep 15`; plain → `TaskLesson 10`; `isReview` → `TaskReview 5`; `Quiz` → `20`.  
_Avoid_: lesson XP flat, word as separate prereq depth (temporary node, future `Variants` fold via `enrichment.json` `improve.md:114`)

**Quiz (150 XP gate)** — Q3 locked to MA verbatim `improve.md:27`: every `150 XP`, timed closed-book, `80%` difficulty via `computeDifficulty` (`internal/engine/engine.go:157` `weakness→difficulty`), diverse recents, immediate remedial + retake (`internal/quiz/quiz.go:1` `improve.md:133`). Daily goal `30` (drift fix shipped: versioned backfill in `sqlite.go`). Banner: `Study → Take Test` + `Profile` due-reviews.  
_Avoid_: 1–2 lesson gate, 30 XP shortcut

**XP** — `computeXPForTask` (`engine.go:607` `10/5/15/20 × timeMult 0.5–1.5 × streakMult 1+0.1*streak`). MA `~1 min` per XP (`improve.md:42`), daily goal `30`. Verified via `go vet` + `study → answer → profile` `playwright`.

**GUEST_KEY** — `mathua_guest_id` persisted so guest `/profile` survives reload (`lib/auth.ts:4`).

## Relationships

- DAG holds many Concepts; Concept has Prerequisites; `Available` filters by mastered
- Profile holds Progress/Scores/Weakness/Activity; Study holds Lessons; Lesson holds Concepts + Problems
- Diagnostic compresses DAG → picks Concept via info-gain; Post-diagnostic → Profile
- Study Lesson `→` TaskLesson 10 (or TaskMultistep 15 if `.word`) `→` SubmitAnswer `→` XP/Progress `→` Profile; at 150 XP `→` Quiz (TaskQuiz 20)

## Flagged ambiguities

- `Learning Mode` infinite loop removed (`session` now chooser `session/page.tsx:1` Diagnostic vs Study, then `Profile` hub); do not reintroduce as primary
- `Study` links that were wrong post-login (`login Skip → /study` etc) fixed to `→ /profile`; do not regress
- `skills/CONTEXT.md` is marketplace language (Issue tracker) — not Mathua; Mathua context lives here

## Decisions (ADRs, inline)

- ADR-001 Q1: Lock `Lesson=corpus`, `Study=library`, `Profile=hub` (2026-08-29) — prevents catalog-as-path relapse
- ADR-002 Q2: `*.word` as `TaskMultistep 15` seam, keep separate node pending `Variants` fold (2026-08-29)
- ADR-003 Q3: Quiz gate `150 XP` MA verbatim, not lesson-count shortcut (2026-08-29) — daily goal `30` makes `150/30=5d` honest
