import { describe, expect, it } from 'vitest'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { sinusoidSchema } from '../components/visual/schemas'
import { sinusoid, sinusoidValue } from '../components/visual/math/sinusoid'

const REPO_ROOT = join(__dirname, '../../..')
const readJson = (path: string): unknown => JSON.parse(readFileSync(join(REPO_ROOT, path), 'utf8'))
const shard = (id: string) => readJson(join('data', 'lessons', 'kp', `${id}.json`)) as Step[]

interface Step {
  label?: string
  visual?: {
    type: string
    title?: string
    caption?: string
    params?: Record<string, unknown>
    predict?: { question: string; options: string[]; answer: number; explain: string }
  }
}

const ids = [
  'trig.graph.sin',
  'trig.graph.cos',
  'trig.graph.amplitude',
  'trig.graph.period',
  'trig.graph.frequency',
  'trig.graph.phase_shift',
  'trig.graph.transformations',
  'trig.graph.sinusoidal_model',
  'trig.graph.harmonic_motion',
  'precalc.trig.graphs',
]

const attached = (id: string) => {
  const step = shard(id).find((item) => item.visual)
  if (!step?.visual) throw new Error(`${id} has no attached visual`)
  return step
}
const waveOf = (id: string) => sinusoid(sinusoidSchema.parse(attached(id).visual!.params))

describe('sinusoid shards agree with the wave engine', () => {
  it('attaches trig.sinusoid to step one only in ten lessons', () => {
    for (const id of ids) {
      const steps = shard(id)
      expect(steps[0].visual?.type).toBe('trig.sinusoid')
      expect(steps.filter((step) => step.visual)).toHaveLength(1)
    }
  })

  it('parses every attached param set and gives finite geometry', () => {
    for (const id of ids) {
      const params = sinusoidSchema.parse(attached(id).visual!.params)
      const wave = sinusoid(params)
      expect(Number.isFinite(wave.peak)).toBe(true)
      expect(Number.isFinite(wave.period)).toBe(true)
      expect(Number.isFinite(wave.max)).toBe(true)
      expect(Number.isFinite(wave.min)).toBe(true)
      expect(Number.isFinite(sinusoidValue(wave, 0.7))).toBe(true)
    }
  })

  it('derives the amplitude, period, shift and midline each lesson claims', () => {
    const sine = waveOf('trig.graph.sin')
    expect(sine.peak).toBeCloseTo(1, 12)
    expect(sinusoidValue(sine, Math.PI / 2)).toBeCloseTo(1, 12)

    const cosine = waveOf('trig.graph.cos')
    expect(cosine.form).toBe('cos')
    expect(sinusoidValue(cosine, 0)).toBeCloseTo(1, 12)

    const amplitude = waveOf('trig.graph.amplitude')
    expect(amplitude.peak).toBeCloseTo(3, 12)
    expect(amplitude.midline).toBeCloseTo(1, 12)
    expect(amplitude.max).toBeCloseTo(4, 12)
    expect(amplitude.min).toBeCloseTo(-2, 12)

    expect(waveOf('trig.graph.period').period).toBeCloseTo(Math.PI, 12)
    expect(waveOf('trig.graph.frequency').frequency).toBeCloseTo(3, 12)

    const phase = waveOf('trig.graph.phase_shift')
    expect(phase.phase).toBeCloseTo(Math.PI / 2, 12)
    expect(sinusoidValue(phase, Math.PI / 2)).toBeCloseTo(0, 12)

    const all = waveOf('trig.graph.transformations')
    expect(all.peak).toBeCloseTo(2, 12)
    expect(all.period).toBeCloseTo(Math.PI, 12)
    expect(all.phase).toBeCloseTo(Math.PI / 2, 12)
    expect(all.midline).toBeCloseTo(1, 12)
    expect(all.max).toBeCloseTo(3, 12)

    const model = waveOf('trig.graph.sinusoidal_model')
    expect(model.midline).toBeCloseTo(6, 12)
    expect(model.peak).toBeCloseTo(3, 12)
    expect(model.max).toBeCloseTo(9, 12)
    expect(model.min).toBeCloseTo(3, 12)

    const harmonic = waveOf('trig.graph.harmonic_motion')
    expect(harmonic.form).toBe('cos')
    expect(harmonic.period).toBeCloseTo(Math.PI, 12)

    const precalc = waveOf('precalc.trig.graphs')
    expect(precalc.peak).toBeCloseTo(2, 12)
    expect(precalc.period).toBeCloseTo((2 * Math.PI) / 3, 12)
    expect(precalc.phase).toBeCloseTo(-Math.PI / 4, 12)
    expect(precalc.midline).toBeCloseTo(-1, 12)
  })

  it('keeps every prediction, answer and explanation tied to the drawn wave', () => {
    const expected: Record<string, { question: string; options: string[]; answer: number; explain: string }> = {
      'trig.graph.sin': {
        question: 'For y = sin x, what is y at x = π/2?',
        options: ['1', '0', '−1'],
        answer: 0,
        explain: 'A quarter turn is π/2. At that angle the unit-circle point is (0, 1), so the height and therefore sin(π/2) is 1.',
      },
      'trig.graph.cos': {
        question: 'For y = cos x, what is y at x = 0?',
        options: ['0', '1', '−1'],
        answer: 1,
        explain: 'At angle 0 the unit-circle point is (1, 0), so the x-coordinate and cos 0 are 1.',
      },
      'trig.graph.amplitude': {
        question: 'For y = 3 sin x + 1, what are the midline and the amplitude?',
        options: ['Midline 1, amplitude 3', 'Midline 3, amplitude 1', 'Midline 0, amplitude 4'],
        answer: 0,
        explain: 'The added constant 1 moves the midline to y = 1. The coefficient 3 gives a peak 3 units above that midline, so the amplitude is 3 and the values run from −2 to 4.',
      },
      'trig.graph.period': {
        question: 'For y = sin(2x), what is the period?',
        options: ['4π', 'π', '2π'],
        answer: 1,
        explain: 'The period is 2π divided by |B|. With B = 2 the period is 2π/2 = π, so the wave repeats twice as often as sin x.',
      },
      'trig.graph.frequency': {
        question: 'How many full cycles does y = sin(3x) complete on 0 ≤ x ≤ 2π?',
        options: ['3', '1/3', '6'],
        answer: 0,
        explain: 'One cycle takes 2π/3, and three of those fit in 2π. The coefficient 3 multiplies the number of cycles, not the period directly.',
      },
      'trig.graph.phase_shift': {
        question: 'How does y = sin(x − π/2) move the graph of y = sin x?',
        options: ['Right by π/2', 'Left by π/2', 'Up by π/2'],
        answer: 0,
        explain: 'The inside expression is zero when x = π/2, so the cycle starts there. A subtraction inside the function shifts the graph to the right by π/2.',
      },
      'trig.graph.transformations': {
        question: 'For y = 2 sin(2(x − π/2)) + 1, which description is correct?',
        options: ['Amplitude 1, period 2π, shift left π/2', 'Amplitude 2, period 2π, shift right 1', 'Amplitude 2, period π, shift right π/2'],
        answer: 2,
        explain: 'The outside coefficient 2 is the amplitude. The inside coefficient 2 makes the period 2π/2 = π. Subtracting π/2 inside shifts the graph right by π/2, and the outside +1 moves the midline.',
      },
      'trig.graph.sinusoidal_model': {
        question: 'A sinusoidal wave runs from a minimum of 3 to a maximum of 9. What are its midline and amplitude?',
        options: ['Midline 6, amplitude 3', 'Midline 3, amplitude 6', 'Midline 12, amplitude 3'],
        answer: 0,
        explain: 'The midline is the average of the maximum and minimum: (9 + 3)/2 = 6. The amplitude is half the distance between them: (9 − 3)/2 = 3.',
      },
      'trig.graph.harmonic_motion': {
        question: 'For d = 2 cos(2t), what is the period?',
        options: ['π/2', 'π', '2π'],
        answer: 1,
        explain: 'The period is 2π divided by the angular frequency. With ω = 2 the period is 2π/2 = π, so the motion repeats every π time units.',
      },
      'precalc.trig.graphs': {
        question: 'For y = 2 sin(3(x + π/4)) − 1, what is the phase shift?',
        options: ['π/4 left', 'π/4 right', 'π/12 left'],
        answer: 0,
        explain: 'The inside is 3(x + π/4), which is 3(x − (−π/4)). The phase shift is C = −π/4, so the graph moves π/4 to the left; the 3 changes the period, not the size of the shift.',
      },
    }

    for (const id of ids) {
      const predict = attached(id).visual!.predict!
      expect(predict).toEqual(expected[id])
      expect(predict.answer).toBeGreaterThanOrEqual(0)
      expect(predict.answer).toBeLessThan(predict.options.length)
    }
  })

  it('keeps the coverage ledger at 31 topics and records the new template as supplementary', () => {
    const coverage = readJson('data/visualizations/coverage.json') as {
      topics: unknown[]
      supplementary: Array<{ id: string; reason: string }>
    }
    expect(coverage.topics).toHaveLength(31)
    const entry = coverage.supplementary.find((item) => item.id === 'trig.sinusoid')
    expect(entry).toBeTruthy()
    expect(entry!.reason).toContain('amplitude')
    expect(entry!.reason).toContain('unit-circle')
  })
})
