// Fetch-only graph access for page-bundle and idle-loaded code.
// No JSON imports here: the graph payload must stay out of First Load JS.
// The static poster and the WebGL scene share one in-flight request each.
import { GraphPayloadSchema, validateGraphPayload, type GraphPayload } from './graphPayload'

type Vec3 = [number, number, number]
type JsonValue = string | number | boolean | null | Array<JsonValue> | { [id: string]: JsonValue }

function isPositionEntry(entry: unknown): entry is [string, Vec3] {
  if (!Array.isArray(entry) || entry.length !== 2) return false
  const [id, v] = entry
  return (
    typeof id === 'string' &&
    Array.isArray(v) &&
    v.length === 3 &&
    v.every(n => typeof n === 'number')
  )
}

function narrowFetched(json: { [id: string]: JsonValue }): Array<[string, Vec3]> {
  return Object.entries(json).filter(isPositionEntry)
}

let inflightPositions: Promise<Array<[string, Vec3]>> | null = null
let inflightPayload: Promise<GraphPayload | null> | null = null

/** Fetched positions for the static poster. Single flight, HTTP-cached. */
export function loadPositionEntries(): Promise<Array<[string, Vec3]>> {
  if (!inflightPositions) {
    inflightPositions = fetch('/data/concepts.positions.json')
      .then(r => (r.ok ? (r.json() as Promise<{ [id: string]: JsonValue }>) : null))
      .then(json => (json === null ? [] : narrowFetched(json)))
      .catch(() => [])
  }
  return inflightPositions
}

/** Fetched nodes + indexed edges for the WebGL scene. Null on failure. */
export function loadGraphPayload(): Promise<GraphPayload | null> {
  if (!inflightPayload) {
    inflightPayload = fetch('/data/concepts.graph.json')
      .then(r => (r.ok ? (r.json() as Promise<JsonValue>) : null))
      .then(json => {
        if (json === null) return null
        const parsed = GraphPayloadSchema.safeParse(json)
        return parsed.success ? validateGraphPayload(parsed.data) : null
      })
      .catch(() => null)
  }
  return inflightPayload
}
