import type { ConceptProgress } from './api'
import { countsAsMastered, prereqMet, type CatalogueConcept } from './progress'

export interface GoalConcept extends CatalogueConcept {
	label?: string
}

export interface GoalProgress {
	learned: number
	total: number
	pct: number
}

function buildIndex(concepts: Iterable<GoalConcept>): Map<string, GoalConcept> {
	const byId = new Map<string, GoalConcept>()
	for (const concept of concepts) byId.set(concept.id, concept)
	return byId
}

function isHandled(p: ConceptProgress | undefined): boolean {
	return countsAsMastered(p) || p?.completed === true || p?.placement_seeded === true
}

/**
 * The goal plus every transitive prerequisite. Unknown goals return an empty
 * list so a stale URL or account value cannot crash the graph.
 */
export function goalScope(goalId: string, concepts: Iterable<GoalConcept>): string[] {
	const byId = buildIndex(concepts)
	if (!byId.has(goalId)) return []

	const scope = new Set<string>()
	const queue: string[] = [goalId]
	for (let i = 0; i < queue.length; i++) {
		const id = queue[i]
		if (scope.has(id)) continue
		scope.add(id)
		const concept = byId.get(id)
		if (!concept) continue
		for (const prereq of concept.prerequisites ?? []) {
			if (!scope.has(prereq) && byId.has(prereq)) queue.push(prereq)
		}
	}
	return [...scope]
}

/**
 * Distance from each concept in the goal scope to the goal, walking backwards
 * along prerequisite edges. The goal is 0, its direct prerequisites are 1.
 */
export function goalDistances(goalId: string, concepts: Iterable<GoalConcept>): Map<string, number> {
	const byId = buildIndex(concepts)
	const distances = new Map<string, number>()
	if (!byId.has(goalId)) return distances

	distances.set(goalId, 0)
	const queue: string[] = [goalId]
	for (let i = 0; i < queue.length; i++) {
		const id = queue[i]
		const distance = distances.get(id) ?? 0
		const concept = byId.get(id)
		if (!concept) continue
		for (const prereq of concept.prerequisites ?? []) {
			if (!byId.has(prereq) || distances.has(prereq)) continue
			distances.set(prereq, distance + 1)
			queue.push(prereq)
		}
	}
	return distances
}

export function goalProgress(
	goalId: string,
	concepts: Iterable<GoalConcept>,
	progress: Record<string, ConceptProgress>,
): GoalProgress {
	const ids = goalScope(goalId, concepts)
	const learned = ids.filter(id => countsAsMastered(progress[id])).length
	const total = ids.length
	return {
		learned,
		total,
		pct: total === 0 ? 0 : Math.round((learned / total) * 100),
	}
}

/**
 * The unmastered concept in the goal scope that is ready now and closest to
 * the goal. Eligibility uses prereqMet, which also accepts completed and
 * placement-seeded prerequisites. The goal itself wins when it is ready.
 */
export function nextGoalStep(
	goalId: string,
	concepts: Iterable<GoalConcept>,
	progress: Record<string, ConceptProgress>,
): string | null {
	const byId = buildIndex(concepts)
	const distances = goalDistances(goalId, byId.values())
	let best: { id: string; distance: number; label: string } | null = null

	for (const id of goalScope(goalId, byId.values())) {
		const concept = byId.get(id)
		if (!concept) continue
		if (isHandled(progress[id])) continue
		if (!(concept.prerequisites ?? []).every(pid => prereqMet(progress, pid))) continue
		const distance = distances.get(id) ?? Number.MAX_SAFE_INTEGER
		const label = concept.label ?? id
		if (
			!best ||
			distance < best.distance ||
			(distance === best.distance && label.localeCompare(best.label) < 0)
		) {
			best = { id, distance, label }
		}
	}
	return best?.id ?? null
}
