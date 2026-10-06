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
   condition for all three tiers as **"Same"**, while labeling the evidence "First
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
| Consistency | recent-weighted accuracy, weights 0.5 → 1.0 | yes — `attempts.correct` |
| Different problem instances | `new Set(attempts.map(a => a.variation)).size`, bonus up to 0.15 | yes — `attempts.question` — `variation` *is* the question text (`LearnStepper.tsx:256`) |
| Difficulty | correct answers weighted `0.6 + 0.4 × difficulty` | **not stored** |
| Time | median elapsed over correct answers, ×0.8 above 60s | yes — `attempts.elapsed_seconds` |

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

---

## What the ladder actually costs (measured)

The ladder is in `internal/mastery/evidence.go` and `Machine.Next` now reads evidence
rather than the streak. Measured end to end through `Engine.submitAnswerWithTask`, on the
two-concept fixture whose generator returns the *same* question text every attempt — the
worst possible case for the variety term:

```
attempt 1: UNSEEN      score 0.800  instances 1   ← LEARNING
attempt 2: LEARNING    score 0.800  instances 1   ← PRACTICING
attempt 3: PRACTICING  score 0.800  instances 1
attempt 4: PRACTICING  score 0.800  instances 1
attempt 5: PRACTICING  score 0.587  instances 1   ← a miss, and the score falls below
attempt 6: PRACTICING  score 0.640  instances 1   ← MASTERED, one good answer later
attempt 7: MASTERED    score 0.658  instances 1
```

**Six answers to master, from any starting `required_streak`, and the exit holds.** The
per-concept threshold no longer affects the cost at all: the score is already normalised
against the concept's own difficulty and time threshold, so concepts differ in what counts
as hard, not in how much evidence proves competence.

- `UNSEEN → LEARNING` takes **one** correct answer. Demanding ten first was the complaint.
- `PRACTICING → MASTERED` needs a score over `MasteryThreshold` across **at least six
  attempts** — a full evidence window.

Six is a floor on *evidence*, not a demand for six consecutive correct answers. It was set to
6 after measurement rather than assumed: `internal/mastery/measure_test.go` walks six
learner profiles (perfect, slowest, one-miss-in-eight, 80% correct, slow-but-right,
starts-rough) against all eight real `avg_time_seconds` values in the corpus, and every one
of them masters in exactly 6. Median, p25, p75, min and max are all 6.

### Does the score actually bind? Yes — and the first measurement said otherwise

**The first version of this measurement was wrong, and its error is worth more than its
result.** It swept accuracy from 40% to 100% and reported every level masterable, concluding
that `MasteryThreshold` was decorative. The profile was `i%20 >= 12` — twelve wrong, then
eight right — and a six-attempt window kept landing entirely inside the clean run. The
finding was an artefact of the probe, not a property of the model. Learners do not answer in
blocks of twenty, and a periodic test pattern will always manufacture a conclusion about
whatever it is probing if the window is short enough to fit inside one period.

With accuracy spread evenly (`correctAt(i, n, k)`, exactly `k` of every `n` correct, evenly
placed so every window sees the nominal rate), the threshold binds cleanly:

| accuracy | outcome |
|---|---|
| 100%, 80%, 60% | MASTERED in exactly 6 answers |
| 50%, 40%, 20%, 0% | never masters, at any length |

So `MasteryThreshold` stays at 0.6. It was right, and the measurement was wrong — which is
the opposite of the usual situation, and the reason the sweep is now
`TestMeasureAccuracyBandThatStillMasters` with a pattern whose property is even distribution
rather than an accident.

A few profiles for detail, across all eight real corpus `avg_time_seconds` values:

| profile | min | p25 | median | p75 | max |
|---|---|---|---|---|---|
| consistently correct | 6 | 6 | 6 | 6 | 6 |
| approximately 80% correct | 6 | 6 | 6 | 6 | 6 |
| approximately 60% correct | 6 | 6 | 6 | 6 | 6 |
| slow but accurate | 6 | 6 | 6 | 6 | 6 |
| correct on the hardest questions | 6 | 6 | 6 | 6 | 6 |
| correct on the easiest questions | 6 | 6 | 6 | 6 | 6 |
| starts rough, then steady | 8 | 8 | 8 | 8 | 8 |
| alternating correct/incorrect | never | | | | |
| approximately 40% correct | never | | | | |
| fast but inaccurate | never | | | | |

Note what does *not* appear: difficulty and speed are not in the cost. A learner who is right
every time masters in 6 answers whether the questions were the easiest or the hardest in the
corpus, and a slow-but-accurate learner masters in 6 exactly like a fast one. The difficulty
and time terms change the *score*, not the number of answers needed to clear a bar that these
profiles clear comfortably. They discriminate at the margin, which is where they matter — see
variety below.

### Variety is a weight, not a veto — and the difference is measurable

ADR-038 recorded that variety "scores but cannot veto", which was imprecise in a way worth
pinning down. The 0.15 bonus is large enough to flip the decision at the margin:

| accuracy | distinct instances | identical instances |
|---|---|---|
| 100%, 80% | masters | masters |
| 60% | masters | **never** |
| 40% | never | never |

So variety genuinely changes the outcome for a marginal learner, and that is correct: three
right answers out of five spread across three different problems is evidence of transfer, and
three right answers out of five on the *same* problem three times is evidence about one
problem. What must not happen is the failure the old `Instances >= 2` hard gate caused —
competence becoming *unreachable*. A learner who is right 80% of the time masters even on a
generator whose question text never varies. That is the asserted property
(`TestVarietyRaisesTheBarWithoutBlockingMastery`), rather than the looser claim that variety
"doesn't matter".

### Two things this change had to fix that were not in the brief

**A gate on a derived signal made mastery unreachable.** The first version of
`EnoughForMastery` also required `Instances >= 2` — distinct problem instances, the stand-in
for "does this transfer?". Wiring it into the engine immediately produced a concept that
could never be mastered, because a generator whose question text does not vary between
attempts holds `Instances` at 1 forever. There is no such generator in production
(`GenerateContext` errors rather than falling back), but the failure mode is the point: a
gate on a signal you do not fully control can lock a state. Variety now *scores* — and it is
worth being precise that "scores" still lets it change a marginal outcome, as measured above.
What it cannot do is make mastery unreachable.

**The decay clock was gated on the rule being replaced.** `LastReviewed` was set only when
`streak >= required_streak && avg time <= threshold` — a proxy for "a tier was just
advanced" back when that was what advancement meant. Under the ladder a concept can reach
MASTERED with a streak of 1, leaving `LastReviewed` nil, and `GetProgress` treats a nil
`LastReviewed` as fully decayed. The engine tests showed the ladder passing **through
DECAYING on its way up**: a concept simultaneously stored as MASTERED and read as DECAYING.
The decay clock now starts on tier attainment, which is what it was always a proxy for.

Also found while wiring this up: the transition logic existed **twice** —
`SubmitAnswer` (study, review) and `submitAnswerWithTask` (quiz, diagnostic) each had their
own copy. Both now read evidence, or a quiz answer would have graded differently from a
study answer, which is the class of divergence ADR-037 exists to prevent.

### Difficulty belongs to the question, not to the host that graded it

Worth its own section, because it was a correctness bug rather than a design gap, and because
the symptom appeared in the wrong place.

**The number.** `BuildEvidence` scored an attempt with no recorded difficulty as though the
question had been the hardest the engine could ask:

```go
difficultyTerm := 1.0        // when diffN == 0
```

The comment above it read "With no difficulty on record the term is dropped entirely rather
than assumed, so old attempts are neither rewarded nor punished for it." Dropping a
*multiplier* means multiplying by 1.0, which is the **top** of the range (`0.6 + 0.4 ×
difficulty` runs 0.72 to 1.0). So a missing value was maximally rewarded, and the comment
described the opposite of what the code did. At 50% accuracy:

| difficulty on the attempt | score | clears 0.6? |
|---|---|---|
| unknown | 0.65 | yes |
| 0.6 | 0.57 | no |

Same learner, same answers, different verdict — decided by which screen they used.

**Why unknown difficulty was common.** Difficulty went *into* `GeneratorContext` and stopped
there. `generator.Problem` did not carry it back out, so every surface that generated a
question and later graded it had to thread the value back by hand:

| host | difficulty source | before |
|---|---|---|
| `SubmitAnswer` (`/api/answer`) | `activeSession.questionDifficulty` | real |
| `submitAnswerWithTask` — **Learn loop** (`/api/study/answer`) | none | `0` → unknown |
| `submitAnswerWithTask` — quizzes (`/api/quiz/answer`) | none | `0` → unknown |
| diagnostic (`/api/goal/diagnostic/answer`) | none | not recorded at all |

The main teaching loop was the host with no difficulty, not the quiz — the opposite of where
this was reported from. The diagnostic wrote its attempt row directly in `server.go`, bypassing
the shared authority entirely, and those rows still feed `GetRecentAttemptsForConcept`, so they
landed in the evidence window as unknown-difficulty attempts.

**The fix**, in the order the architecture implies:

1. `generator.Problem.Difficulty *float64`, stamped by `Registry.GenerateContext` and
   `BatchGenerateContext`. It is now structurally impossible to generate a question without
   its difficulty travelling alongside it.
2. `engine.Anchor` carries it too, because the anchor is already the durable record of what
   was served and is authoritative over the expected answer for the same reason. The study
   serve fills it from the Problem (generator path) or from `questions.difficulty`
   (curated path — `NOT NULL DEFAULT 0.5`, so real).
3. Quiz passes `sess.LastProblem.Difficulty`; the diagnostic stamps the row it writes.
4. `gradedAttempt.difficulty` became `*float64`, so "unknown" is a value a caller must choose
   rather than a `0` that gets reinterpreted downstream. That overloading is what hid the bug.
5. Unknown is credited at `mastery.NeutralDifficulty = 0.55`, which is not chosen here — it
   is what `difficultyFromWeakness` serves when it knows nothing about the learner (weakness
   defaults to 0.5, so `0.3 + 0.5 × 0.5`). `TestNeutralDifficultyMatchesEngineDefault` in the
   engine package fails if the scoring side and the serving side ever drift apart.

**What the shared authority deliberately does not do:** guess. Re-deriving difficulty from the
learner's *current* weakness would credit an attempt with the difficulty a **future** question
would get, which is a different claim from the one the learner answered. `applyGradedAttempt`
receives the value and does not reinterpret it.

**Known remaining divergence, stated rather than hidden:** `web/next-app/lib/progression.ts`
is a second copy of this model and does not match it — it weights difficulty per attempt
rather than by the mean of correct attempts, uses a fixed 60s time threshold rather than the
concept's own, and has no unknown-difficulty case because its `difficulty` field is a required
number. So this change creates no new browser/server divergence, but the pre-existing one
remains: the two disagree at the margin, and the browser's copy feeds `LearnStepper`'s
`advance` decision and header label. The server ladder is authoritative for progression, so
that is a display-consistency problem, not a correctness one.