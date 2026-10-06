import { describe, it, expect, vi, beforeEach } from 'vitest'
import { fetchShelf, recommendationResToShelf, recommendationToShelfItem } from '../lib/recommendations'
import type { RecommendationRes } from '../lib/api'

/**
 * The client's half of "the engine decides".
 *
 * Everything below is presentation. The decision — which concept, in what order, why, and with
 * what wording — is made in internal/scheduler and arrives over `/api/next`; the tests that
 * matter for it are `internal/scheduler/recommend_test.go` and `TestNext_*` in internal/server.
 *
 * What is asserted here is the contract the surfaces rely on, because that is the part a
 * refactor can silently break: the payload maps onto the shape the components render, the
 * server's wording passes through untouched, and a failure degrades to the library rather than
 * to an empty page.
 */

const rec = (over: Partial<RecommendationRes['primary']> = {}) => ({
  id: 'b',
  conceptId: 'b',
  conceptTitle: 'Beta',
  kind: 'learn' as const,
  reason: 'new' as const,
  priority: 1.5,
  xp: 2,
  estimatedMinutes: 1,
  action: { type: 'learn' as const, href: '/learn?concept=b' },
  badge: 'New',
  detail: 'Ready to learn',
  cta: 'Start →',
  ...over,
})

describe('recommendation mapping', () => {
  it("passes the server wording through untouched", () => {
    // The duplication defect was two components rewording one head: "Currently working towards
    // X / Answer it" beside "Next up / Continue: X". So nothing here may rewrite the copy.
    const r = rec({ badge: 'Practice', detail: 'Worth another go', cta: 'Keep practicing →' })
    const item = recommendationToShelfItem(r)
    expect(item.badge).toBe('Practice')
    expect(item.detail).toBe('Worth another go')
    expect(item.cta).toBe('Keep practicing →')
    expect(item.href).toBe('/learn?concept=b')
    expect(item.title).toBe('Beta')
    expect(item.xp).toBe(2)
  })

  it('maps each engine kind onto the surface vocabulary', () => {
    const kind = (k: 'learn' | 'practice' | 'review' | 'mastery_check') =>
      recommendationToShelfItem(rec({ kind: k })).kind
    expect(kind('review')).toBe('review')
    expect(kind('practice')).toBe('weakness')
    expect(kind('mastery_check')).toBe('diagnostic')
  })

  it('distinguishes resuming a concept from starting one', () => {
    // Both are kind "learn" on the wire; `reason` is what separates them, and the surfaces
    // render them differently ("Practice again" against "Start"). Mapping on reason is what
    // keeps that distinction from being lost.
    const fresh = recommendationToShelfItem(rec({ reason: 'new' }))
    const resuming = recommendationToShelfItem(rec({ reason: 'prerequisite' }))
    expect(fresh.kind).toBe('new')
    expect(resuming.kind).toBe('resume')
  })

  it('keeps the library fallback and the diagnostic identifiable', () => {
    // Both are ids rather than concepts, and both must not be mistaken for one.
    expect(recommendationToShelfItem(rec({ id: 'study', conceptId: '', kind: 'learn' })).kind).toBe('browse')
    expect(recommendationToShelfItem(rec({ id: 'diagnostic', kind: 'learn' })).kind).toBe('diagnostic')
  })

  it('never leaves alternatives undefined', () => {
    // The client iterates it without a null check, so a missing field must not become undefined.
    const shelf = recommendationResToShelf({ primary: rec(), generatedAt: '' } as RecommendationRes)
    expect(shelf.alternatives).toEqual([])
    expect(shelf.next.href).toBe('/learn?concept=b')
  })
})

describe('fetchShelf', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('asks once and adapts', async () => {
    const fetchMock = vi.fn(async () => ({
      ok: true,
      json: async () => ({ primary: rec(), alternatives: [], generatedAt: '2026-10-06T00:00:00Z' }),
    }))
    vi.stubGlobal('fetch', fetchMock)
    const shelf = await fetchShelf()
    expect(shelf.next.title).toBe('Beta')
    // One request, not five. The client used to assemble this from activity + progress +
    // weaknesses + reviews + scores plus a bundled copy of the corpus.
    expect(fetchMock).toHaveBeenCalledTimes(1)
  })

  it('passes the exclusion to the engine rather than filtering afterwards', async () => {
    // The concept just finished is unmastered by definition, so filtering client-side would
    // still leave it eligible server-side and it would head its own shelf.
    const fetchMock = vi.fn(async () => ({
      ok: true, json: async () => ({ primary: rec(), alternatives: [], generatedAt: '' }),
    }))
    vi.stubGlobal('fetch', fetchMock)
    await fetchShelf(['c'])
    const url = String((fetchMock.mock.calls as unknown as string[][])[0][0])
    expect(url).toContain('exclude=c')
  })

  it('degrades to the library rather than to nothing', async () => {
    // Every one of the five requests this replaced had a fallback, so a network blip already
    // degraded gracefully and must keep doing so. An empty home page would not.
    vi.stubGlobal('fetch', vi.fn(async () => ({ ok: false, status: 500 })))
    const shelf = await fetchShelf()
    expect(shelf.next.href).toBe('/study')
    expect(shelf.next.kind).toBe('browse')
  })
})
