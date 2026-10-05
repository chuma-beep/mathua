// The authority boundary between the server's mastery ladder and the client's
// `masteryEstimate`.
//
// CONTEXT.md ADR-039 makes the server ladder authoritative. This file exists so that
// boundary cannot quietly move, and so the ways the two models disagree are enumerated
// rather than discovered later.
//
// The short version, verified here rather than asserted in prose:
//
//   * `masteryEstimate`'s `decision` has **no authority at all**. `LearnStepper` reads it
//     once, as `est.decision === 'advance' || next >= REQUIRED_IN_A_ROW` — and the first
//     operand is dead. `decision === 'advance'` *requires* the last two attempts correct
//     (`streak2`), which is the same fact as `consecutive >= 2`, which already satisfies
//     the second operand. "advance never fires alone" proves that exhaustively over every
//     history of length 1-8.
//
//   * What survives is the `band`, rendered as text. So the client model is a label, and
//     the divergences below are differences in a label's wording and in a number nobody
//     acts on — not two definitions of "mastered" competing for the same decision.

import { describe, it, expect } from 'vitest'
import { readFileSync } from 'node:fs'
import { masteryEstimate, REQUIRED_IN_A_ROW, type Attempt } from '../lib/progression'

const stepperSrc = readFileSync('components/LearnStepper.tsx', 'utf8')

/** Mirrors the Learn loop's own accounting, which is three lines of state. */
function consecutiveRun(history: Attempt[]): number {
  return history.reduce((run, a) => (a.correct ? run + 1 : 0), 0)
}

/** Exactly the expression at components/LearnStepper.tsx:275. */
function learnLoopWouldAdvance(history: Attempt[], estAdvance: boolean): boolean {
  return estAdvance || consecutiveRun(history) + 1 >= REQUIRED_IN_A_ROW
}

/** Every history of `length` correctness bits, as attempts at one difficulty. */
function* allHistories(length: number, difficulty: number): Generator<Attempt[]> {
  for (let mask = 0; mask < 1 << length; mask++) {
    const history: Attempt[] = []
    for (let i = 0; i < length; i++) {
      history.push({
        correct: Boolean(mask & (1 << i)),
        difficulty,
        elapsed: 5,
        // Distinct per attempt, so the variety term saturates and can never be the reason
        // two histories diverge.
        variation: `q${i}`,
      })
    }
    yield history
  }
}

const attempt = (correct: boolean, difficulty: number, variation: string): Attempt => ({
  correct,
  difficulty,
  elapsed: 5,
  variation,
})

describe('the client estimate cannot decide anything', () => {
  it('advance never fires alone', () => {
    let advanceDecisions = 0
    let decidedAlone = 0

    for (let length = 1; length <= 8; length++) {
      for (const history of allHistories(length, 0.6)) {
        if (masteryEstimate(history).decision !== 'advance') continue
        advanceDecisions++

        // The claim: whenever the client says advance, the 2-in-a-row rule is *also*
        // satisfied, so deleting the client's opinion from line 275 changes nothing.
        const twoInARow = consecutiveRun(history) + 1 >= REQUIRED_IN_A_ROW
        if (!twoInARow) decidedAlone++
        expect(learnLoopWouldAdvance(history, true)).toBe(twoInARow)
      }
    }

    // Guard the guard: zero advance decisions would mean the test had stopped exercising
    // the path and would pass vacuously.
    expect(advanceDecisions).toBeGreaterThan(0)
    expect(decidedAlone).toBe(0)
  })

  it('has exactly one read site, and it is the dead one', () => {
    // Structural rather than behavioural: `est.decision` appears once in the app. Combined
    // with the test above, that single site cannot decide anything, so `intervene` and
    // `continue` are unreachable rather than merely unused.
    expect(stepperSrc.match(/est\.decision/g) ?? []).toHaveLength(1)
    expect(stepperSrc).toContain("est.decision === 'advance' || next >= REQUIRED_IN_A_ROW")
  })

  it('the score gates only the label it labels', () => {
    expect(stepperSrc.match(/est\.score/g) ?? []).toHaveLength(1)
    expect(stepperSrc).toContain('{est.score > 0 &&')
  })
})

describe('masteryEstimate is not imported anywhere but the Learn stepper', () => {
  it('has one production importer', () => {
    // The blast radius of this model is one component. `lib/progress.ts` and
    // `lib/nextUp.ts` are *consumers* of the server's status, not second models: they read
    // `ConceptProgress.status` and never compute a score.
    const importers = stepperSrc.includes("from '../lib/progression'")
    expect(importers).toBe(true)

    const progress = readFileSync('lib/progress.ts', 'utf8')
    const nextUp = readFileSync('lib/nextUp.ts', 'utf8')
    for (const src of [progress, nextUp]) {
      expect(src).not.toMatch(/masteryEstimate/)
      // Neither may invent a score of its own.
      expect(src).not.toMatch(/0\.6\s*\+\s*0\.4/)
    }
  })
})

describe('documented divergences from the authoritative model', () => {
  // Each of these is a way `masteryEstimate` and `internal/mastery.BuildEvidence` disagree
  // today. They are asserted so the disagreement stays visible: if one stops being true,
  // this fails and asks for the ADR to be updated, rather than the difference quietly
  // disappearing or quietly growing. None of them affects a decision — see above.

  it('difficulty: recency-weighted per attempt on the client, an unweighted mean on the server', () => {
    // Identical multiset of correctness and of difficulty; only the pairing differs. The
    // hard question sits on a correct attempt in one and on the missed attempt in the
    // other.
    const hardOnMiss = [
      attempt(true, 0.3, 'a'),
      attempt(true, 0.3, 'b'),
      attempt(false, 1.0, 'c'),
    ]
    const hardOnHit = [
      attempt(true, 1.0, 'a'),
      attempt(true, 0.3, 'b'),
      attempt(false, 0.3, 'c'),
    ]

    expect(masteryEstimate(hardOnMiss).score).toBeCloseTo(0.5, 6)
    expect(masteryEstimate(hardOnHit).score).toBeCloseTo(0.562222, 5)

    // The client is sensitive to *which* attempt was hard, because the difficulty factor is
    // multiplied by the recency weight attempt-by-attempt. The server averages the correct
    // attempts' difficulty first and applies it once, so it does not see that difference.
    // If this ever becomes false, the client has been aligned with the server and ADR-039
    // should be updated.
    expect(masteryEstimate(hardOnMiss).score).not.toBeCloseTo(masteryEstimate(hardOnHit).score, 6)
  })

  it('variety: the client counts every attempt in the window, the server only correct ones', () => {
    // Three correct on three instances plus a miss on a fourth. The client's
    // `new Set(recent.map(a => a.variation))` sees four distinct variations and awards the
    // full 0.15. `Evidence.Instances` is documented as spanning *the correct attempts*, so
    // the server sees three and awards 0.10.
    const history = [
      attempt(true, 0.6, 'a'),
      attempt(true, 0.6, 'b'),
      attempt(true, 0.6, 'c'),
      attempt(false, 0.6, 'd'),
    ]
    // Client: recency-weighted accuracy 0.56 x 0.84 = 0.47... precisely
    // (0.5 + 0.667 + 0.833) x 0.84 / 3.0 = 0.56, plus the full 0.15 for four distinct
    // variations = 0.71.
    expect(masteryEstimate(history).score).toBeCloseTo(0.71, 6)

    // The server's figure for the same evidence is 0.66: the same 0.56 accuracy-weighted
    // term, plus 0.10 because only three of the four instances were correct.
    //
    // Both land on 'strong' here, but the gap is enough to change the label on other
    // evidence — which is the part that matters, since the label is the only thing left of
    // this model that anybody reads.
    const missCorrectMiss = [
      attempt(false, 0.6, 'a'),
      attempt(true, 0.6, 'b'),
      attempt(false, 0.6, 'c'),
    ]
    expect(masteryEstimate(missCorrectMiss).score).toBeCloseTo(0.38, 6)
    expect(masteryEstimate(missCorrectMiss).band).toBe('developing')
    // The server scores the same three answers 0.28 — (0.75 / 2.25) x 0.84, plus no
    // variety bonus because a single correct attempt spans one instance — which reads
    // 'unseen'. Same evidence, a different word to the learner.
    //
    // That server figure is deliberately a comment and not an assertion. Asserting it here
    // would mean transcribing the authoritative formula into a second language, which is
    // the thing this whole pass is removing. If the server model changes, the right place
    // for that number to change is its own test.
  })

  it('time: a fixed 60s cliff on the client, the concept threshold on the server', () => {
    // The client cannot know the concept's `avg_time_seconds` — `masteryEstimate` receives
    // only the attempt history — so it hardcodes 60. The corpus thresholds run 10s to 150s,
    // with 162 of 657 concepts at 120s, so for most of the corpus the client penalises a
    // median the server would credit in full.
    const onTime = [attempt(true, 0.6, 'a'), attempt(true, 0.6, 'b'), attempt(true, 0.6, 'c')]
    const slow = onTime.map((a, i) => ({ ...a, elapsed: 90, variation: `q${i}` }))

    expect(masteryEstimate(onTime).score).toBeGreaterThan(masteryEstimate(slow).score)
    // 90s is over the client's 60s cliff. For a concept whose own threshold is 120s the
    // server applies no penalty at all, so the client's deduction is not one the
    // authoritative model would make.
    expect(masteryEstimate(slow).score / masteryEstimate(onTime).score).toBeCloseTo(0.8, 6)
  })
})

describe('the client cannot express an unknown difficulty', () => {
  it('Attempt.difficulty is required, so the NeutralDifficulty branch has no counterpart here', () => {
    // `masteryEstimate` has no missing-difficulty branch because the type forbids one. The
    // server does — `NeutralDifficulty` — which is precisely why ADR-039 had to fix the
    // plumbing rather than reconcile the two formulas.
    //
    // The consequence to keep in mind: a caller passing a difficulty it *requested* rather
    // than one the server *served* is silently reporting the wrong number, and the type
    // cannot catch it. `LearnStepper` passes its own `difficulty` state, which the server
    // honours on the generator path and ignores on the curated-questions path. That
    // divergence is dormant today because the questions table ships empty (0 rows) and
    // `scripts/seed-questions` is run by hand.
    // Worth noticing while we are here: a single correct answer scores 0.84 and so reads
    // 'strong'. That is not a divergence between the two models — they agree at 0.84 — it
    // is the band's vocabulary having no evidence floor, where the mastery ladder has one.
    // Another reason the label should not be read as a mastery claim.
    const est = masteryEstimate([attempt(true, 0.6, 'a')])
    expect(est.score).toBeCloseTo(0.84, 6)
    expect(est.band).toBe('strong')
  })
})