import { describe, it, expect } from 'vitest'
import {
  buildGraphMeta,
  buildGraphPayload,
  validateGraphPayload,
  GraphPayloadSchema,
  GRAPH_PAYLOAD_VERSION,
} from '../lib/graphPayload'
import type { ConceptRecord } from '../lib/conceptData'

const concepts: ConceptRecord[] = [
  { id: 'a', label: 'A', domain: 'arith', prerequisites: [] },
  { id: 'b', label: 'B', domain: 'arith', prerequisites: ['a'] },
  { id: 'c', label: 'C', domain: 'algebra', prerequisites: ['a', 'b'] },
  { id: 'd', label: 'D', domain: 'algebra', prerequisites: ['missing'] },
]

describe('buildGraphPayload', () => {
  it('indexes nodes and maps prerequisites to edge indices', () => {
    const payload = buildGraphPayload(concepts, { a: 3, b: 1, c: 0, d: 0 })
    expect(payload.v).toBe(GRAPH_PAYLOAD_VERSION)
    expect(payload.nodes).toEqual([
      ['a', 'A', 'arith', 3],
      ['b', 'B', 'arith', 1],
      ['c', 'C', 'algebra', 0],
      ['d', 'D', 'algebra', 0],
    ])
    expect(payload.edges).toEqual([
      [0, 1],
      [0, 2],
      [1, 2],
    ])
    expect(payload.domains).toEqual({ arith: 2, algebra: 2 })
  })

  it('counts domains and connections in the meta artifact', () => {
    const meta = buildGraphMeta(concepts)
    expect(meta.conceptCount).toBe(4)
    expect(meta.connectionCount).toBe(4)
    expect(meta.domainCount).toBe(2)
    expect(meta.domainCounts).toEqual({ arith: 2, algebra: 2 })
  })
})

describe('GraphPayloadSchema + validateGraphPayload', () => {
  it('accepts a well-formed payload', () => {
    const payload = buildGraphPayload(concepts, {})
    const parsed = GraphPayloadSchema.safeParse(JSON.parse(JSON.stringify(payload)))
    expect(parsed.success).toBe(true)
    if (parsed.success) {
      expect(validateGraphPayload(parsed.data)).not.toBeNull()
    }
  })

  it('rejects malformed payloads', () => {
    expect(GraphPayloadSchema.safeParse(null).success).toBe(false)
    expect(GraphPayloadSchema.safeParse({ v: 2, nodes: [['a', 'A', 'x', 1]], edges: [[0]] }).success).toBe(false)
    expect(GraphPayloadSchema.safeParse({ v: 2, nodes: 'nope', edges: [] }).success).toBe(false)
  })

  it('rejects edges that point outside the node array', () => {
    const payload = buildGraphPayload(concepts, {})
    expect(validateGraphPayload({ ...payload, edges: [[0, 9]] })).toBeNull()
    expect(validateGraphPayload({ ...payload, edges: [[-1, 0]] })).toBeNull()
  })
})
