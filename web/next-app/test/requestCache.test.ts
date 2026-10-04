import { describe, it, expect, vi } from 'vitest'
import { createSessionCache, createInflightCache } from '../lib/requestCache'

const tick = () => new Promise<void>(r => setTimeout(r, 0))

describe('createSessionCache', () => {
  it('shares one request between concurrent callers, so StrictMode does not double it', async () => {
    const once = createSessionCache()
    const run = vi.fn().mockResolvedValue('x')
    const [a, b] = await Promise.all([once('k', run), once('k', run)])
    expect(run).toHaveBeenCalledTimes(1)
    expect([a, b]).toEqual(['x', 'x'])
  })

  it('reuses a successful response for the rest of the session', async () => {
    const once = createSessionCache()
    const run = vi.fn().mockResolvedValue('x')
    await once('k', run)
    await once('k', run)
    expect(run).toHaveBeenCalledTimes(1)
  })

  it('does not cache a failure', async () => {
    const once = createSessionCache()
    const run = vi.fn().mockRejectedValueOnce(new Error('offline')).mockResolvedValue('back')
    await expect(once('k', run)).rejects.toThrow('offline')
    await expect(once('k', run)).resolves.toBe('back')
    expect(run).toHaveBeenCalledTimes(2)
  })
})

describe('createInflightCache', () => {
  it('shares one request between concurrent callers', async () => {
    const dedupe = createInflightCache()
    const run = vi.fn().mockResolvedValue('x')
    const [a, b] = await Promise.all([dedupe('k', run), dedupe('k', run)])
    expect(run).toHaveBeenCalledTimes(1)
    expect([a, b]).toEqual(['x', 'x'])
  })

  it('refetches after the request settles — the bug this exists to fix', async () => {
    const dedupe = createInflightCache()
    const run = vi.fn().mockResolvedValue('x')

    await dedupe('k', run)
    await tick()
    await dedupe('k', run)
    await tick()
    await dedupe('k', run)

    // The old helper kept a *successful* promise and deleted only on rejection, so a
    // learner who answered questions and revisited the page kept seeing the numbers they
    // had before answering. Progress, scores and weaknesses are written by the server the
    // moment an answer is graded; the copy on the client was the stale part.
    expect(run).toHaveBeenCalledTimes(3)
  })

  it('refetches after a failure, so the page recovers on revisit', async () => {
    const dedupe = createInflightCache()
    const run = vi.fn().mockRejectedValueOnce(new Error('offline')).mockResolvedValue('back')
    await expect(dedupe('k', run)).rejects.toThrow('offline')
    await tick()
    await expect(dedupe('k', run)).resolves.toBe('back')
    expect(run).toHaveBeenCalledTimes(2)
  })

  it('keys separately', async () => {
    const dedupe = createInflightCache()
    const run = vi.fn().mockResolvedValue('x')
    await Promise.all([dedupe('a', run), dedupe('b', run)])
    expect(run).toHaveBeenCalledTimes(2)
  })
})