// Topo rank for Study ordering (PR7): Kahn's topological order over the
// concept DAG, deterministic with lexicographic tiebreaks. Lessons within a
// domain sort by their earliest concept's rank, so browsing follows prereq
// order instead of title order. Cheap approximation — the scheduler (not
// this order) still decides what to practice.

export interface TopoNode {
  id: string
  prerequisites?: string[]
}

export function topoRank(catalog: TopoNode[]): Map<string, number> {
  const ids = new Set(catalog.map(c => c.id))
  const indegree = new Map<string, number>()
  const dependents = new Map<string, string[]>()
  for (const c of catalog) {
    indegree.set(c.id, 0)
    dependents.set(c.id, [])
  }
  for (const c of catalog) {
    const seen = new Set<string>()
    for (const p of c.prerequisites ?? []) {
      // Ignore prereqs outside the catalog (enrichment/variant edges) and
      // dedupe — they don't constrain the order we can compute here.
      if (!ids.has(p) || seen.has(p)) continue
      seen.add(p)
      indegree.set(c.id, (indegree.get(c.id) ?? 0) + 1)
      dependents.get(p)!.push(c.id)
    }
  }
  const queue: string[] = []
  for (const [id, d] of indegree) {
    if (d === 0) queue.push(id)
  }
  const rank = new Map<string, number>()
  while (queue.length > 0) {
    // Keep deterministic: always pop the lexicographically smallest ready id.
    queue.sort()
    const id = queue.shift()!
    rank.set(id, rank.size)
    for (const dep of dependents.get(id) ?? []) {
      const d = (indegree.get(dep) ?? 0) - 1
      indegree.set(dep, d)
      if (d === 0) queue.push(dep)
    }
  }
  // Cycles (rejected server-side by validate_graph) would leave nodes
  // unranked — append them in id order so every id has a rank.
  if (rank.size < ids.size) {
    for (const id of [...ids].sort()) {
      if (!rank.has(id)) rank.set(id, rank.size)
    }
  }
  return rank
}
