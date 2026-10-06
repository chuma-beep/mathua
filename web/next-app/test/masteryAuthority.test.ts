// The authority boundary between the server's mastery ladder and the web client.
//
// CONTEXT.md ADR-039 makes the server ladder authoritative; ADR-040 removes the client's
// copy of the evidence model. This file is the guard on that.
//
// The boundary used to need defending rather than asserting. `lib/progression.ts` carried a
// second implementation of the mastery evidence score — recency-weighted accuracy ×
// difficulty × variety × time — with its own weights, its own thresholds and its own band
// vocabulary. It was wired into `LearnStepper` as though it decided something:
//
//   const advance = est.decision === 'advance' || next >= REQUIRED_IN_A_ROW
//
// where `decision === 'advance'` *required* the last two attempts correct — the same fact
// as `next >= REQUIRED_IN_A_ROW`. The first operand could never be the deciding half, which
// an exhaustive sweep over all 510 histories of length 1-8 confirmed (111 advance
// decisions, zero deciding alone). So the model had no authority and could not be reconciled
// with the server's because nothing was riding on it.
//
// It is gone. What is asserted here is that it cannot come back:
//
//   * `lib/progression.ts` no longer computes a score, and the client's copy of the
//     difficulty, variety and time terms is gone with it.
//   * No client module derives a score of its own.
//   * `LearnStepper` advances on the lesson rule alone, and renders the server's word
//     rather than a locally-computed one.

import { describe, it, expect } from 'vitest'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'

const read = (path: string) => readFileSync(path, 'utf8')

const stepper = read('components/LearnStepper.tsx')
const progression = read('lib/progression.ts')

/**
 * The source with its comments removed.
 *
 * Several assertions below are about the *shape of the code*, and several of the comments
 * in these files deliberately quote the code they replaced — a comment that names the old
 * disjunction would otherwise make "the old disjunction is gone" fail against its own
 * explanation. Assertions that are genuinely about prose keep using `read`.
 */
const codeOnly = (src: string) => src.replace(/\/\*[\s\S]*?\*\//g, '').replace(/^\s*\/\/.*$/gm, '')

describe('the client has no mastery model', () => {
  it('lib/progression.ts computes no score', () => {
    expect(codeOnly(progression)).not.toContain('masteryEstimate')
    // The four terms of the evidence model, gone from the client.
    expect(codeOnly(progression)).not.toMatch(/recent-weighted|variationBonus|timeFactor/)
    // The distinctive constant of the removed model.
    expect(progression).not.toMatch(/0\.6\s*\+\s*0\.4/)
  })

  it('nothing imports a score model any more', () => {
    expect(codeOnly(stepper)).not.toContain('masteryEstimate')
    expect(codeOnly(stepper)).not.toMatch(/from '\.\.\/lib\/progression'.*mastery/)
  })

  it('no client module derives a score of its own', () => {
    // The guard that would have caught the second model, and would catch the next one.
    // `lib/progress.ts` and `lib/nextUp.ts` are *consumers* of the server's status; neither
    // may compute anything from an attempt history.
    for (const path of ['lib/progress.ts', 'lib/nextUp.ts', 'components/LearnStepper.tsx']) {
      const src = codeOnly(read(path))
      expect(src, path).not.toMatch(/0\.6\s*\+\s*0\.4/)
      expect(src, path).not.toMatch(/0\.05\s*\*\s*Math\.max/)
      expect(src, path).not.toMatch(/0\.4\s*\*\s*a\.difficulty/)
    }
  })

  it('the client reads the progress word the server sent', () => {
    // The band is rendered, never computed. A score rendered as a number would be a fourth
    // place to disagree; a string cannot be.
    expect(stepper).toContain('res.evidence_band')
    expect(stepper).toMatch(/\{band &&/)
    expect(stepper).toMatch(/setBand\(res\.evidence_band\)/)
  })
})

describe('the lesson advance rule is the whole gate', () => {
  it('advances on 2-in-a-row alone', () => {
    expect(stepper).toContain('const advance = next >= REQUIRED_IN_A_ROW')
  })

  it('the old disjunction is gone, and cannot be reintroduced silently', () => {
    // Comments are stripped first: the explanatory comment above deliberately quotes the
    // old expression, so a raw substring check would fail against its own explanation.
    const assignments = codeOnly(stepper).match(/const advance = .*/g) ?? []
    expect(assignments).toHaveLength(1)
    expect(assignments[0]).toBe('const advance = next >= REQUIRED_IN_A_ROW')
  })

  it('REQUIRED_IN_A_ROW is still the lesson rule, not a mastery gate', () => {
    // It is 2-in-a-row from improve.md:39 — worked example then a couple of problems. It
    // moves between knowledge points inside a lesson. Mastery is the server's decision and
    // this constant does not participate in it, which is why the two can coexist without
    // either contradicting the other.
    expect(progression).toContain('export const REQUIRED_IN_A_ROW = 2')
    expect(progression).toContain('2-in-a-row correct to advance')
  })
})

describe('the score on the wire is the authoritative one', () => {
  it('StudyAnswerRes declares both, and marks the band as presentational', () => {
    const api = read('lib/api.ts')
    expect(api).toMatch(/evidence_score\?: number/)
    expect(api).toMatch(/evidence_band\?: string/)
    // The docs on the field must keep saying what it is, so nobody later reads it as a
    // mastery verdict.
    expect(api).toMatch(/claiming nothing about retention/)
  })
})

describe('the recommendation is decided by the engine, not the client', () => {
  const read = (f: string) => readFileSync(join(__dirname, '..', f), 'utf8')
  const strip = (s: string) => s.replace(/\/\*[\s\S]*?\*\//g, '').replace(/\/\/.*$/gm, '')
  const code = (f: string) => strip(read(f))

  it('lib/nextUp.ts computes no ranking', () => {
    const src = code('lib/nextUp.ts')
    // The client used to build candidates, weight them and choose the head here. Two copies of
    // the eligibility rules and of the XP award is how they came to disagree with
    // internal/scheduler — the client's still halved reviews after ADR-020 removed that from
    // the server. The decision is engine-side now, tested in internal/scheduler and by
    // TestNext_* in internal/server.
    for (const banned of [
      'buildCandidates',
      'selectShelfHead',
      'fixedWeightPolicy',
      'SelectionPolicy',
      'effortBase',
      'xpFor',
    ]) {
      expect(src, `lib/nextUp.ts still contains ${banned}`).not.toContain(banned)
    }
    // No sorting of candidates either: the engine orders them.
    expect(src).not.toMatch(/\.sort\(.*weakness/)
  })

  it('no surface ranks for itself', () => {
    // /learn, /profile and the Learn done card all read the same shelf. If one of them starts
    // sorting or filtering by score again, the two surfaces can disagree about what "next" is —
    // which is the duplication defect this replaced.
    for (const f of ['app/learn/page.tsx', 'app/profile/page.tsx', 'components/LearnStepper.tsx']) {
      const src = code(f)
      expect(src, `${f} must not compute a priority`).not.toMatch(/priority\s*[<>]=?/)
      expect(src, `${f} must not weight candidates`).not.toMatch(/weakness\s*\*\s*0\./)
    }
  })

  it('every surface asks the engine, and only the engine', () => {
    for (const f of ['app/learn/page.tsx', 'app/profile/page.tsx', 'components/LearnStepper.tsx']) {
      expect(read(f), `${f} should read the shelf from lib/recommendations`).toContain('fetchShelf')
      expect(code(f), `${f} should not import the ranking helpers`).not.toContain('selectShelfHead')
    }
    // code(), not read(): the module documents what it used to contain, and a comment
    // naming the deleted function is not the function coming back.
    expect(code('lib/nextUp.ts')).not.toContain('selectShelfHead')
  })

  it('the engine owns the wording, so two surfaces cannot describe one task differently', () => {
    // The duplication defect: /profile rendered "Currently working towards X / Answer it" and
    // "Next up / Continue: X" from the same object, sometimes with different labels for it. So
    // the copy arrives from the server and the mapper passes it through.
    const mapper = code('lib/recommendations.ts')
    expect(mapper).toContain('r.detail')
    expect(mapper).toContain('r.cta')
    expect(mapper).toContain('r.badge')
    const types = read('lib/api.ts')
    for (const field of ['badge: string', 'detail: string', 'cta: string']) {
      expect(types, `Recommendation should carry ${field} from the server`).toContain(field)
    }
  })
})
