/**
 * Two different intents that were sharing one helper, and a stale-progress bug.
 *
 * `oncePerSession` was introduced to stop React 18 StrictMode double-invoking effects
 * from doubling network traffic, and it did that — but it deleted its map entry only on
 * *rejection*, so a **successful** response was reused for the rest of the JavaScript
 * session. On `/graph` that meant a learner who answered questions and came back saw
 * pre-answer mastery bars, node colours and headline numbers: answering is a real server
 * write (`UpsertProgress` + `AddXP`) and none of it was reflected. The map was named
 * `inflight`, so the intent was in-flight deduplication all along.
 *
 * `dedupeInFlight` is that intent, honestly implemented: the entry goes when the promise
 * settles either way. StrictMode's two effect invocations still share the promise —
 * React runs both before the microtask queue drains, so nothing has settled yet — and a
 * later mount refetches.
 *
 * Which one to use is a question about the *data*, not about the caller:
 *
 * - `oncePerSession` for facts that cannot change while the tab is open: connectivity,
 *   and the DAG topology served by `GET /api/graph`, which is `dag.Order()` over the
 *   bundled corpus.
 * - `dedupeInFlight` for anything derived from the learner's history: progress, scores,
 *   weaknesses, lessons. These change the moment an answer is graded, and a stale copy is
 *   not a slow refresh, it is a wrong number.
 */

/** Runs `key` at most once per JavaScript session, sharing the promise with later callers. */
export function createSessionCache() {
  const settled = new Map<string, Promise<unknown>>()
  return function oncePerSession<T>(key: string, run: () => Promise<T>): Promise<T> {
    const existing = settled.get(key)
    if (existing) return existing as Promise<T>
    // A rejection is *not* kept: the caller has already been told it failed, and leaving a
    // rejected promise in the map would make every later attempt fail without trying.
    const p = run().catch(err => {
      settled.delete(key)
      throw err
    })
    settled.set(key, p)
    return p
  }
}

/** Shares one in-flight request between concurrent callers, and forgets it on settle. */
export function createInflightCache() {
  const inflight = new Map<string, Promise<unknown>>()
  return function dedupeInFlight<T>(key: string, run: () => Promise<T>): Promise<T> {
    const existing = inflight.get(key)
    if (existing) return existing as Promise<T>
    const p = run().finally(() => {
      inflight.delete(key)
    })
    inflight.set(key, p)
    return p
  }
}