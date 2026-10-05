# Mastery semantics: what mastery costs, and what it should mean

Status: **investigation and design. No threshold was changed.**

This answers two questions raised after ADR-037:

- **P1 — is mastery too expensive, and is the cost intentional?**
- **P2 — what would it take to make mastery mean *demonstrated competence* rather
  than *thirty correct answers in a row*?**

Everything below is measured against the corpus and the code, not estimated.

---

## Summary

Mastery costs **3 × `mastery_threshold.streak` consecutive correct answers**, because
the engine resets the streak to 1 on every tier advance. For **459 of 657 concepts
(70%)** that is **30 in a row**. A single miss sets the streak to 0, so a run restarts
from the beginning: at a per-answer accuracy of 0.90 the chance of ever completing a
30-answer run on the first attempt is **4.2%**.

The cost is **not intentional and not documented**. The commit that introduced the
reset is titled *"fix mastery state machine — reset streak on level-up"* with no body,
`docs/system-design.md` states `required_streak` is *"typically 3"* when the corpus says
10 for 70% of concepts, and that document's own condition column records all three tier
transitions as **"Same"**.

The model the brief describes — independent demonstrations, varied problem instances,
increasing difficulty, consistent quality — **already exists in this repository**, at
`web/next-app/lib/progression.ts:masteryEstimate`, and its comment says *"Streak is one
input, not the decision."* It is client-side, it only produces a label, and the server
gates progression on the streak counter instead. **Three of its four inputs are already
persisted server-side.** The fourth, difficulty, is the only missing piece.

Recommendation: do not lower thresholds. Promote the evidence model to the server, add
the one missing persisted field, and keep the tier ladder as a *shape* rather than three
identical repetitions of the same test.

---

## P1 — What mastery costs today

### Answers required

Measured over `web/next-app/data/concepts.json` (657 concepts, all with a threshold):

| Answers for MASTERED | Concepts | Share |
|---|---|---|
| 9 | 171 | 26% |
| 15 | 21 | 3% |
| 21 | 5 | 1% |
| **30** | **459** | **70%** |
| 45 | 1 | 0.2% |

Mastering the whole corpus this way is **15,774 consecutive correct answers**.

### Time floor

Using each concept's `avg_time_seconds` — a floor, since practised answers are faster and
the threshold is a cap:

- **130 hours** of pure answer time for the corpus
- per concept: median **12 min**, p90 **18 min**, max **30 min**

`/learn` serves **3 questions per batch** (`LearnStepper.tsx:97`), so 30 consecutive is
ten rounds with no miss anywhere in them.

### The barrier is probabilistic, not just arithmetic

Each incorrect answer sets `Streak = 0` (`engine.go:1093`), so a run is all-or-nothing:

| Answers in a row | p = 0.95 | p = 0.90 | p = 0.85 |
|---|---|---|---|
| 9 | 63.0% | 38.7% | 23.2% |
| **30** | **21.5%** | **4.2%** | **0.8%** |

A learner who is right 90% of the time — strong, and the app's own `weakness_score`
treats 0.9 as weak — has a 96% chance of *not* completing a 30-answer run on the first
sitting, and every miss costs the whole run. That is the mechanism by which mastery can
become its own purgatory: the learner re-proving something they have demonstrated,
frequently, for a long time.

### Is it intentional? No.

Three pieces of evidence:

1. **The commit has no body.** `725fad9` — *"fix mastery state machine — reset streak on
   level-up, add interactive CLI loop (all tests pass)"*. The reset reads as a fix to a
   tier-advance problem, not as a decision to triple the cost of mastery.
2. **The documentation contradicts the corpus.** `docs/system-design.md` says
   *"Each concept has a `required_streak` (typically 3)"*. The actual distribution is 10
   for 459 concepts, 3 for 171. The documented world costs 9 answers; the real one costs
   30 for most of the corpus.
3. **The documentation records the redundancy itself.** Its transition table gives the
   condition for all three tiers as **"Same"**, while labelling the evidence "First
   correct", "Consistent correct", "Mastery achieved". Three tiers gated by one identical
   test is three repetitions of the same measurement, not three kinds of evidence. It
   also documents `streak` resetting to 0 on an incorrect answer and not at all on a tier
   advance.

And the strongest evidence is the state diagram itself. `docs/diagrams/student-model.svg`
labels the three transitions:

```
UNSEEN       --first correct-->  LEARNING
LEARNING     --streak reached--> PRACTICING
PRACTICING   --interval elapsed-> MASTERED
```

**The design intended three different kinds of evidence, the third of them a spacing
interval.** The implementation collapsed all three into `streak >= required AND
avg time <= threshold`. So spaced attainment is not a new wish — it is the third rung of
the ladder as drawn, and it is the one rung whose evidence the data model cannot record.

There is no ADR, no design note and no test asserting the 3× cost.

---

## P1 — The evidence model already exists

`web/next-app/lib/progression.ts:masteryEstimate`, over the last 6 attempts:

| Signal | How it is computed | Persisted server-side? |
|---|---|---|
| Consistency | recent-weighted accuracy, weights 0.5 → 1.0 | ✅ `attempts.correct` |
| Different problem instances | `new Set(attempts.map(a => a.variation)).size`, bonus up to 0.15 | ✅ `attempts.question` — `variation` *is* the question text (`LearnStepper.tsx:256`) |
| Difficulty | correct answers weighted `0.6 + 0.4 × difficulty` | ❌ **not stored** |
| Time | median elapsed over correct answers, ×0.8 above 60s | ✅ `attempts.elapsed_seconds` |

It also produces a decision, not just a score: `advance` needs `score ≥ 0.6` **and** ≥3
attempts **and** 2-in-a-row at the end; `intervene` at 3 recent misses.

So three of four inputs already exist per attempt, and the fourth is one column.

**The gap is authority, not invention.** `masteryEstimate` renders a label in the
`/learn` sticky header (`· <band>`); `Machine.Next` decides everything else. Two models
of the same fact, one of them far better informed, and they do not talk.

---

## P2 — Design: mastery as demonstrated competence

### What the evidence should be

Per concept, over a bounded window of recent attempts:

```
score = (weightedAccuracy × difficultyWeight + varietyBonus) × timeFactor
```

with, for a first implementation:

- **weightedAccuracy** — recent-weighted correct rate, so a late miss still counts
- **difficultyWeight** — mean difficulty of the *correct* attempts, not a flat 1.0
- **varietyBonus** — distinct question instances among those attempts
- **timeFactor** — median elapsed against the concept's `avg_time_seconds`

Mastery at `PRACTICING → MASTERED` when the score clears a threshold **and** the evidence
spans more than one attempt and more than one question instance.

The tier ladder stops being three identical repetitions:

| Tier | Evidence it should add | Currently |
|---|---|---|
| UNSEEN → LEARNING | one correct answer | full streak |
| LEARNING → PRACTICING | accuracy at or above the time threshold | full streak |
| PRACTICING → MASTERED | score + variety + difficulty | full streak |

That is what makes the cost comparable to today while stopping it being the *same test
three times*: the early tiers get cheaper and easier (which is right — a learner should
not need 10 correct to be told they are learning), and the final tier carries the most
evidence.

### Acquisition vs retention

The brief's distinction is the right one and the current model conflates them:

```
ACQUIRE   — "can you solve this now?"      → gates leaving /learn
RETENTION — "can you retrieve it later?"   → schedules review
```

Mastery should be an **acquisition** verdict. Retention is already modelled by SM-2 and
`decayDays`; conflating them is why a learner cannot demonstrate retention before being
allowed to move forward. ADR-037 made the exit purely acquisition-based; the missing half
is that nothing *within* acquisition rewards a gap.

### Spacing, honestly

The ladder as drawn (`student-model.svg`) already ends in **"interval elapsed"** — a
spacing interval before `MASTERED`. That is the shape to restore, and the part the data
model cannot currently express:

- `concept_progress` has fourteen columns; none records when the current tier was attained.
- `mastered_at` is written once, for the final tier.
- `last_reviewed` refreshes on any answer meeting the threshold and is already the decay
  clock — it cannot be overloaded.

Making the three tiers require a minimum gap therefore needs a new column. It is a real
migration, not a threshold, and per this brief it should not be added speculatively.

**Recommendation: do not add the migration yet.** `sm2_repetitions` already grows with
each successful review and `next_review_due` already moves out, so durability improves
automatically as a learner passes reviews. What that cannot give is spaced evidence
*before* first mastery. If the acquisition model above lands and the 30-answer run is
gone, the remaining gap is far less costly than it looks today.

### Migration, if it is ever needed

- **SQLite** — `ALTER TABLE concept_progress ADD COLUMN tier_attained_at TEXT` in
  `internal/storage/migrate.go`, alongside the existing column additions.
- **Postgres** — the same statement in the Postgres migration path.
- **Parity** — ADR-012 requires an identical scenario against both stores; the existing
  `parity_test.go` covers the progress table, so a new column needs a case there or it is
  untested on one engine.
- **Backfill** — `NULL` means "unknown", and unknown must not read as satisfied *or* as
  unsatisfied. Existing `MASTERED` rows keep their status; the gate applies only to new
  advances, so no learner regresses.
- **Rollout** — behind the score model, not a flag. Two models of mastery at once is what
  ADR-037 spent a week untangling.

---

## What not to do

- **Do not lower the global threshold.** It lowers the bar uniformly and teaches the same
  thing 3× instead of 10×. The cost is not the number, it is that the number is the only
  evidence.
- **Do not add a `RETIRED` state.** `sm2_repetitions` and `next_review_due` already express
  durability.
- **Do not gate mastery on retention.** That is what forces long-term proof before
  progress, which is the purgatory in its other form.
- **Do not add the migration because it would be nice.** Decide what evidence justifies
  "Mathua can stop actively teaching this" first.

---

## Open questions for the product owner

1. **What is the intended number of answers to master a threshold-10 concept?** The docs
   say 9 (3 tiers × 3). The corpus says 30 (3 × 10). Which is the intent?
2. **Should the early tiers be cheaper?** Acquiring 10 correct before being told you are
   `LEARNING` is a lot to ask of the first attempt.
3. **Is difficulty worth persisting?** It is the only missing input, and one column. Its
   absence means every correct answer currently counts the same regardless of challenge.
4. **How long should a concept stay out of `/learn` after mastery?** ADR-037 makes it
   "until review says otherwise". Confirm that is the intent rather than an artefact of
   the fix.