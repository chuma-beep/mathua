import { describe, expect, it } from 'vitest'

import type { ConceptProgress } from '../lib/api'
import {
	goalDistances,
	goalProgress,
	goalScope,
	nextGoalStep,
	type GoalConcept,
} from '../lib/goalScope'

const concepts: GoalConcept[] = [
	{ id: 'a', label: 'A', domain: 'd', prerequisites: [] },
	{ id: 'b', label: 'B', domain: 'd', prerequisites: ['a'] },
	{ id: 'c', label: 'C', domain: 'd', prerequisites: ['b'] },
	{ id: 'd', label: 'D', domain: 'd', prerequisites: ['a'] },
	{ id: 'e', label: 'E', domain: 'd', prerequisites: ['c', 'd'] },
	{ id: 'unrelated', label: 'Unrelated', domain: 'd', prerequisites: [] },
]

function progress(status: string, extra: Partial<ConceptProgress> = {}): ConceptProgress {
	return { status, streak: 0, ...extra }
}

describe('goalScope', () => {
	it('returns the goal and its transitive prerequisites only', () => {
		expect(goalScope('e', concepts).sort()).toEqual(['a', 'b', 'c', 'd', 'e'])
	})

	it('does not include branches that are not prerequisites', () => {
		expect(goalScope('b', concepts).sort()).toEqual(['a', 'b'])
	})

	it('returns an empty scope for an unknown goal', () => {
		expect(goalScope('missing', concepts)).toEqual([])
	})

	it('skips unknown prerequisite references', () => {
		const withMissing = [...concepts, { id: 'g', domain: 'd', prerequisites: ['nope'] }]
		expect(goalScope('g', withMissing)).toEqual(['g'])
	})
})

describe('goalDistances', () => {
	it('measures reverse prerequisite distance from the goal', () => {
		const distances = goalDistances('e', concepts)
		expect(distances.get('e')).toBe(0)
		expect(distances.get('c')).toBe(1)
		expect(distances.get('d')).toBe(1)
		expect(distances.get('b')).toBe(2)
		expect(distances.get('a')).toBe(2)
		expect(distances.has('unrelated')).toBe(false)
	})
})

describe('goalProgress', () => {
	it('counts mastered and decaying concepts as learned', () => {
		const result = goalProgress('e', concepts, {
			a: progress('MASTERED'),
			b: progress('DECAYING'),
		})
		expect(result).toEqual({ learned: 2, total: 5, pct: 40 })
	})

	it('does not count completed-only concepts as learned', () => {
		const result = goalProgress('b', concepts, {
			a: progress('LEARNING', { completed: true }),
		})
		expect(result).toEqual({ learned: 0, total: 2, pct: 0 })
	})

	it('returns zero progress for an unknown goal', () => {
		expect(goalProgress('missing', concepts, {})).toEqual({ learned: 0, total: 0, pct: 0 })
	})
})

describe('nextGoalStep', () => {
	it('chooses the eligible unlearned concept closest to the goal', () => {
		expect(
			nextGoalStep('e', concepts, {
				a: progress('MASTERED'),
			}),
		).toBe('d')
	})

	it('treats completed and placement-seeded concepts as eligible', () => {
		expect(
			nextGoalStep('e', concepts, {
				a: progress('MASTERED'),
				d: progress('LEARNING', { completed: true }),
				b: progress('LEARNING', { placement_seeded: true }),
			}),
		).toBe('c')
	})

	it('returns the goal when the goal itself is eligible', () => {
		expect(
			nextGoalStep('e', concepts, {
				a: progress('MASTERED'),
				b: progress('MASTERED'),
				c: progress('MASTERED'),
				d: progress('MASTERED'),
			}),
		).toBe('e')
	})

	it('returns null when no concept in scope is ready', () => {
		const blocked: GoalConcept[] = [
			{ id: 'goal', label: 'Goal', domain: 'd', prerequisites: ['missing'] },
		]
		expect(nextGoalStep('goal', blocked, {})).toBeNull()
	})
})
